































package generator

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/generator/internal/remap"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)





const generatedCodeVersion = 2



type Plugin interface {
	
	Name() string
	
	
	Init(g *Generator)
	
	
	Generate(file *FileDescriptor)
	
	
	GenerateImports(file *FileDescriptor)
}

var plugins []Plugin



func RegisterPlugin(p Plugin) {
	plugins = append(plugins, p)
}


type GoImportPath string

func (p GoImportPath) String() string { return strconv.Quote(string(p)) }


type GoPackageName string








type common struct {
	file *FileDescriptor 
}


func (c *common) GoImportPath() GoImportPath {
	return c.file.importPath
}

func (c *common) File() *FileDescriptor { return c.file }

func fileIsProto3(file *descriptor.FileDescriptorProto) bool {
	return file.GetSyntax() == "proto3"
}

func (c *common) proto3() bool { return fileIsProto3(c.file.FileDescriptorProto) }


type Descriptor struct {
	common
	*descriptor.DescriptorProto
	parent   *Descriptor            
	nested   []*Descriptor          
	enums    []*EnumDescriptor      
	ext      []*ExtensionDescriptor 
	typename []string               
	index    int                    
	path     string                 
	group    bool
}



func (d *Descriptor) TypeName() []string {
	if d.typename != nil {
		return d.typename
	}
	n := 0
	for parent := d; parent != nil; parent = parent.parent {
		n++
	}
	s := make([]string, n)
	for parent := d; parent != nil; parent = parent.parent {
		n--
		s[n] = parent.GetName()
	}
	d.typename = s
	return s
}



type EnumDescriptor struct {
	common
	*descriptor.EnumDescriptorProto
	parent   *Descriptor 
	typename []string    
	index    int         
	path     string      
}



func (e *EnumDescriptor) TypeName() (s []string) {
	if e.typename != nil {
		return e.typename
	}
	name := e.GetName()
	if e.parent == nil {
		s = make([]string, 1)
	} else {
		pname := e.parent.TypeName()
		s = make([]string, len(pname)+1)
		copy(s, pname)
	}
	s[len(s)-1] = name
	e.typename = s
	return s
}



func (e *EnumDescriptor) prefix() string {
	if e.parent == nil {
		
		return CamelCase(*e.Name) + "_"
	}
	typeName := e.TypeName()
	return CamelCaseSlice(typeName[0:len(typeName)-1]) + "_"
}


func (e *EnumDescriptor) integerValueAsString(name string) string {
	for _, c := range e.Value {
		if c.GetName() == name {
			return fmt.Sprint(c.GetNumber())
		}
	}
	log.Fatal("cannot find value for enum constant")
	return ""
}



type ExtensionDescriptor struct {
	common
	*descriptor.FieldDescriptorProto
	parent *Descriptor 
}



func (e *ExtensionDescriptor) TypeName() (s []string) {
	name := e.GetName()
	if e.parent == nil {
		
		s = make([]string, 1)
	} else {
		pname := e.parent.TypeName()
		s = make([]string, len(pname)+1)
		copy(s, pname)
	}
	s[len(s)-1] = name
	return s
}


func (e *ExtensionDescriptor) DescName() string {
	
	typeName := e.TypeName()
	
	for i, s := range typeName {
		typeName[i] = CamelCase(s)
	}
	return "E_" + strings.Join(typeName, "_")
}


type ImportedDescriptor struct {
	common
	o Object
}

func (id *ImportedDescriptor) TypeName() []string { return id.o.TypeName() }




type FileDescriptor struct {
	*descriptor.FileDescriptorProto
	desc []*Descriptor          
	enum []*EnumDescriptor      
	ext  []*ExtensionDescriptor 
	imp  []*ImportedDescriptor  

	
	comments map[string]*descriptor.SourceCodeInfo_Location

	
	
	
	exported map[Object][]symbol

	fingerprint string        
	importPath  GoImportPath  
	packageName GoPackageName 

	proto3 bool 
}




func (d *FileDescriptor) VarName() string {
	name := strings.Map(badToUnderscore, baseName(d.GetName()))
	return fmt.Sprintf("fileDescriptor_%s_%s", name, d.fingerprint)
}





func (d *FileDescriptor) goPackageOption() (impPath GoImportPath, pkg GoPackageName, ok bool) {
	opt := d.GetOptions().GetGoPackage()
	if opt == "" {
		return "", "", false
	}
	
	sc := strings.Index(opt, ";")
	if sc >= 0 {
		return GoImportPath(opt[:sc]), cleanPackageName(opt[sc+1:]), true
	}
	
	slash := strings.LastIndex(opt, "/")
	if slash >= 0 {
		return GoImportPath(opt), cleanPackageName(opt[slash+1:]), true
	}
	return "", cleanPackageName(opt), true
}


func (d *FileDescriptor) goFileName(pathType pathType) string {
	name := *d.Name
	if ext := path.Ext(name); ext == ".proto" || ext == ".protodevel" {
		name = name[:len(name)-len(ext)]
	}
	name += ".pb.go"

	if pathType == pathTypeSourceRelative {
		return name
	}

	
	
	if impPath, _, ok := d.goPackageOption(); ok && impPath != "" {
		
		_, name = path.Split(name)
		name = path.Join(string(impPath), name)
		return name
	}

	return name
}

func (d *FileDescriptor) addExport(obj Object, sym symbol) {
	d.exported[obj] = append(d.exported[obj], sym)
}


type symbol interface {
	
	
	GenerateAlias(g *Generator, pkg GoPackageName)
}

type messageSymbol struct {
	sym                         string
	hasExtensions, isMessageSet bool
	oneofTypes                  []string
}

type getterSymbol struct {
	name     string
	typ      string
	typeName string 
	genType  bool   
}

func (ms *messageSymbol) GenerateAlias(g *Generator, pkg GoPackageName) {
	g.P("type ", ms.sym, " = ", pkg, ".", ms.sym)
	for _, name := range ms.oneofTypes {
		g.P("type ", name, " = ", pkg, ".", name)
	}
}

type enumSymbol struct {
	name   string
	proto3 bool 
}

func (es enumSymbol) GenerateAlias(g *Generator, pkg GoPackageName) {
	s := es.name
	g.P("type ", s, " = ", pkg, ".", s)
	g.P("var ", s, "_name = ", pkg, ".", s, "_name")
	g.P("var ", s, "_value = ", pkg, ".", s, "_value")
}

type constOrVarSymbol struct {
	sym  string
	typ  string 
	cast string 
}

func (cs constOrVarSymbol) GenerateAlias(g *Generator, pkg GoPackageName) {
	v := string(pkg) + "." + cs.sym
	if cs.cast != "" {
		v = cs.cast + "(" + v + ")"
	}
	g.P(cs.typ, " ", cs.sym, " = ", v)
}


type Object interface {
	GoImportPath() GoImportPath
	TypeName() []string
	File() *FileDescriptor
}


type Generator struct {
	*bytes.Buffer

	Request  *plugin.CodeGeneratorRequest  
	Response *plugin.CodeGeneratorResponse 

	Param             map[string]string 
	PackageImportPath string            
	ImportPrefix      string            
	ImportMap         map[string]string 

	Pkg map[string]string 

	outputImportPath GoImportPath                   
	allFiles         []*FileDescriptor              
	allFilesByName   map[string]*FileDescriptor     
	genFiles         []*FileDescriptor              
	file             *FileDescriptor                
	packageNames     map[GoImportPath]GoPackageName 
	usedPackages     map[GoImportPath]bool          
	usedPackageNames map[GoPackageName]bool         
	typeNameToObject map[string]Object              
	init             []string                       
	indent           string
	pathType         pathType 
	writeOutput      bool
	annotateCode     bool                                       
	annotations      []*descriptor.GeneratedCodeInfo_Annotation 
}

type pathType int

const (
	pathTypeImport pathType = iota
	pathTypeSourceRelative
)


func New() *Generator {
	g := new(Generator)
	g.Buffer = new(bytes.Buffer)
	g.Request = new(plugin.CodeGeneratorRequest)
	g.Response = new(plugin.CodeGeneratorResponse)
	return g
}


func (g *Generator) Error(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ":" + err.Error()
	log.Print("protoc-gen-go: error:", s)
	os.Exit(1)
}


func (g *Generator) Fail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Print("protoc-gen-go: error:", s)
	os.Exit(1)
}




func (g *Generator) CommandLineParameters(parameter string) {
	g.Param = make(map[string]string)
	for _, p := range strings.Split(parameter, ",") {
		if i := strings.Index(p, "="); i < 0 {
			g.Param[p] = ""
		} else {
			g.Param[p[0:i]] = p[i+1:]
		}
	}

	g.ImportMap = make(map[string]string)
	pluginList := "none" 
	for k, v := range g.Param {
		switch k {
		case "import_prefix":
			g.ImportPrefix = v
		case "import_path":
			g.PackageImportPath = v
		case "paths":
			switch v {
			case "import":
				g.pathType = pathTypeImport
			case "source_relative":
				g.pathType = pathTypeSourceRelative
			default:
				g.Fail(fmt.Sprintf(`Unknown path type %q: want "import" or "source_relative".`, v))
			}
		case "plugins":
			pluginList = v
		case "annotate_code":
			if v == "true" {
				g.annotateCode = true
			}
		default:
			if len(k) > 0 && k[0] == 'M' {
				g.ImportMap[k[1:]] = v
			}
		}
	}
	if pluginList != "" {
		
		enabled := make(map[string]bool)
		for _, name := range strings.Split(pluginList, "+") {
			enabled[name] = true
		}
		var nplugins []Plugin
		for _, p := range plugins {
			if enabled[p.Name()] {
				nplugins = append(nplugins, p)
			}
		}
		plugins = nplugins
	}
}




func (g *Generator) DefaultPackageName(obj Object) string {
	importPath := obj.GoImportPath()
	if importPath == g.outputImportPath {
		return ""
	}
	return string(g.GoPackageName(importPath)) + "."
}


func (g *Generator) GoPackageName(importPath GoImportPath) GoPackageName {
	if name, ok := g.packageNames[importPath]; ok {
		return name
	}
	name := cleanPackageName(baseName(string(importPath)))
	for i, orig := 1, name; g.usedPackageNames[name]; i++ {
		name = orig + GoPackageName(strconv.Itoa(i))
	}
	g.packageNames[importPath] = name
	g.usedPackageNames[name] = true
	return name
}

var globalPackageNames = map[GoPackageName]bool{
	"fmt":   true,
	"math":  true,
	"proto": true,
}



func RegisterUniquePackageName(pkg string, f *FileDescriptor) string {
	name := cleanPackageName(pkg)
	for i, orig := 1, name; globalPackageNames[name]; i++ {
		name = orig + GoPackageName(strconv.Itoa(i))
	}
	globalPackageNames[name] = true
	return string(name)
}

var isGoKeyword = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"else":        true,
	"defer":       true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
}

func cleanPackageName(name string) GoPackageName {
	name = strings.Map(badToUnderscore, name)
	
	if isGoKeyword[name] {
		name = "_" + name
	}
	
	if r, _ := utf8.DecodeRuneInString(name); unicode.IsDigit(r) {
		name = "_" + name
	}
	return GoPackageName(name)
}



func (g *Generator) defaultGoPackage() GoPackageName {
	p := g.PackageImportPath
	if i := strings.LastIndex(p, "/"); i >= 0 {
		p = p[i+1:]
	}
	return cleanPackageName(p)
}




func (g *Generator) SetPackageNames() {
	g.outputImportPath = g.genFiles[0].importPath

	defaultPackageNames := make(map[GoImportPath]GoPackageName)
	for _, f := range g.genFiles {
		if _, p, ok := f.goPackageOption(); ok {
			defaultPackageNames[f.importPath] = p
		}
	}
	for _, f := range g.genFiles {
		if _, p, ok := f.goPackageOption(); ok {
			
			f.packageName = p
		} else if p, ok := defaultPackageNames[f.importPath]; ok {
			
			
			
			
			
			f.packageName = p
		} else if p := g.defaultGoPackage(); p != "" {
			
			
			
			
			f.packageName = p
		} else if p := f.GetPackage(); p != "" {
			
			f.packageName = cleanPackageName(p)
		} else {
			
			f.packageName = cleanPackageName(baseName(f.GetName()))
		}
	}

	
	for _, f := range g.genFiles[1:] {
		if a, b := g.genFiles[0].importPath, f.importPath; a != b {
			g.Fail(fmt.Sprintf("inconsistent package import paths: %v, %v", a, b))
		}
		if a, b := g.genFiles[0].packageName, f.packageName; a != b {
			g.Fail(fmt.Sprintf("inconsistent package names: %v, %v", a, b))
		}
	}

	
	
	g.Pkg = map[string]string{
		"fmt":   "fmt",
		"math":  "math",
		"proto": "proto",
	}
}




func (g *Generator) WrapTypes() {
	g.allFiles = make([]*FileDescriptor, 0, len(g.Request.ProtoFile))
	g.allFilesByName = make(map[string]*FileDescriptor, len(g.allFiles))
	genFileNames := make(map[string]bool)
	for _, n := range g.Request.FileToGenerate {
		genFileNames[n] = true
	}
	for _, f := range g.Request.ProtoFile {
		fd := &FileDescriptor{
			FileDescriptorProto: f,
			exported:            make(map[Object][]symbol),
			proto3:              fileIsProto3(f),
		}
		
		if substitution, ok := g.ImportMap[f.GetName()]; ok {
			
			
			
			fd.importPath = GoImportPath(substitution)
		} else if genFileNames[f.GetName()] && g.PackageImportPath != "" {
			
			
			
			
			fd.importPath = GoImportPath(g.PackageImportPath)
		} else if p, _, _ := fd.goPackageOption(); p != "" {
			
			
			
			fd.importPath = p
		} else {
			
			
			
			fd.importPath = GoImportPath(path.Dir(f.GetName()))
		}
		
		fd.desc = wrapDescriptors(fd)
		g.buildNestedDescriptors(fd.desc)
		fd.enum = wrapEnumDescriptors(fd, fd.desc)
		g.buildNestedEnums(fd.desc, fd.enum)
		fd.ext = wrapExtensions(fd)
		extractComments(fd)
		g.allFiles = append(g.allFiles, fd)
		g.allFilesByName[f.GetName()] = fd
	}
	for _, fd := range g.allFiles {
		fd.imp = wrapImported(fd, g)
	}

	g.genFiles = make([]*FileDescriptor, 0, len(g.Request.FileToGenerate))
	for _, fileName := range g.Request.FileToGenerate {
		fd := g.allFilesByName[fileName]
		if fd == nil {
			g.Fail("could not find file named", fileName)
		}
		fingerprint, err := fingerprintProto(fd.FileDescriptorProto)
		if err != nil {
			g.Error(err)
		}
		fd.fingerprint = fingerprint
		g.genFiles = append(g.genFiles, fd)
	}
}




func fingerprintProto(m proto.Message) (string, error) {
	b, err := proto.Marshal(m)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:8]), nil
}


func (g *Generator) buildNestedDescriptors(descs []*Descriptor) {
	for _, desc := range descs {
		if len(desc.NestedType) != 0 {
			for _, nest := range descs {
				if nest.parent == desc {
					desc.nested = append(desc.nested, nest)
				}
			}
			if len(desc.nested) != len(desc.NestedType) {
				g.Fail("internal error: nesting failure for", desc.GetName())
			}
		}
	}
}

func (g *Generator) buildNestedEnums(descs []*Descriptor, enums []*EnumDescriptor) {
	for _, desc := range descs {
		if len(desc.EnumType) != 0 {
			for _, enum := range enums {
				if enum.parent == desc {
					desc.enums = append(desc.enums, enum)
				}
			}
			if len(desc.enums) != len(desc.EnumType) {
				g.Fail("internal error: enum nesting failure for", desc.GetName())
			}
		}
	}
}


func newDescriptor(desc *descriptor.DescriptorProto, parent *Descriptor, file *FileDescriptor, index int) *Descriptor {
	d := &Descriptor{
		common:          common{file},
		DescriptorProto: desc,
		parent:          parent,
		index:           index,
	}
	if parent == nil {
		d.path = fmt.Sprintf("%d,%d", messagePath, index)
	} else {
		d.path = fmt.Sprintf("%s,%d,%d", parent.path, messageMessagePath, index)
	}

	
	
	if parent != nil {
		parts := d.TypeName()
		if file.Package != nil {
			parts = append([]string{*file.Package}, parts...)
		}
		exp := "." + strings.Join(parts, ".")
		for _, field := range parent.Field {
			if field.GetType() == descriptor.FieldDescriptorProto_TYPE_GROUP && field.GetTypeName() == exp {
				d.group = true
				break
			}
		}
	}

	for _, field := range desc.Extension {
		d.ext = append(d.ext, &ExtensionDescriptor{common{file}, field, d})
	}

	return d
}


func wrapDescriptors(file *FileDescriptor) []*Descriptor {
	sl := make([]*Descriptor, 0, len(file.MessageType)+10)
	for i, desc := range file.MessageType {
		sl = wrapThisDescriptor(sl, desc, nil, file, i)
	}
	return sl
}


func wrapThisDescriptor(sl []*Descriptor, desc *descriptor.DescriptorProto, parent *Descriptor, file *FileDescriptor, index int) []*Descriptor {
	sl = append(sl, newDescriptor(desc, parent, file, index))
	me := sl[len(sl)-1]
	for i, nested := range desc.NestedType {
		sl = wrapThisDescriptor(sl, nested, me, file, i)
	}
	return sl
}


func newEnumDescriptor(desc *descriptor.EnumDescriptorProto, parent *Descriptor, file *FileDescriptor, index int) *EnumDescriptor {
	ed := &EnumDescriptor{
		common:              common{file},
		EnumDescriptorProto: desc,
		parent:              parent,
		index:               index,
	}
	if parent == nil {
		ed.path = fmt.Sprintf("%d,%d", enumPath, index)
	} else {
		ed.path = fmt.Sprintf("%s,%d,%d", parent.path, messageEnumPath, index)
	}
	return ed
}


func wrapEnumDescriptors(file *FileDescriptor, descs []*Descriptor) []*EnumDescriptor {
	sl := make([]*EnumDescriptor, 0, len(file.EnumType)+10)
	
	for i, enum := range file.EnumType {
		sl = append(sl, newEnumDescriptor(enum, nil, file, i))
	}
	
	for _, nested := range descs {
		for i, enum := range nested.EnumType {
			sl = append(sl, newEnumDescriptor(enum, nested, file, i))
		}
	}
	return sl
}


func wrapExtensions(file *FileDescriptor) []*ExtensionDescriptor {
	var sl []*ExtensionDescriptor
	for _, field := range file.Extension {
		sl = append(sl, &ExtensionDescriptor{common{file}, field, nil})
	}
	return sl
}


func wrapImported(file *FileDescriptor, g *Generator) (sl []*ImportedDescriptor) {
	for _, index := range file.PublicDependency {
		df := g.fileByName(file.Dependency[index])
		for _, d := range df.desc {
			if d.GetOptions().GetMapEntry() {
				continue
			}
			sl = append(sl, &ImportedDescriptor{common{file}, d})
		}
		for _, e := range df.enum {
			sl = append(sl, &ImportedDescriptor{common{file}, e})
		}
		for _, ext := range df.ext {
			sl = append(sl, &ImportedDescriptor{common{file}, ext})
		}
	}
	return
}

func extractComments(file *FileDescriptor) {
	file.comments = make(map[string]*descriptor.SourceCodeInfo_Location)
	for _, loc := range file.GetSourceCodeInfo().GetLocation() {
		if loc.LeadingComments == nil {
			continue
		}
		var p []string
		for _, n := range loc.Path {
			p = append(p, strconv.Itoa(int(n)))
		}
		file.comments[strings.Join(p, ",")] = loc
	}
}




func (g *Generator) BuildTypeNameMap() {
	g.typeNameToObject = make(map[string]Object)
	for _, f := range g.allFiles {
		
		
		
		dottedPkg := "." + f.GetPackage()
		if dottedPkg != "." {
			dottedPkg += "."
		}
		for _, enum := range f.enum {
			name := dottedPkg + dottedSlice(enum.TypeName())
			g.typeNameToObject[name] = enum
		}
		for _, desc := range f.desc {
			name := dottedPkg + dottedSlice(desc.TypeName())
			g.typeNameToObject[name] = desc
		}
	}
}



func (g *Generator) ObjectNamed(typeName string) Object {
	o, ok := g.typeNameToObject[typeName]
	if !ok {
		g.Fail("can't find object with type", typeName)
	}

	
	
	
	
	direct := *o.File().Name == *g.file.Name
	if !direct {
		for _, dep := range g.file.Dependency {
			if *g.fileByName(dep).Name == *o.File().Name {
				direct = true
				break
			}
		}
	}
	if !direct {
		found := false
	Loop:
		for _, dep := range g.file.Dependency {
			df := g.fileByName(*g.fileByName(dep).Name)
			for _, td := range df.imp {
				if td.o == o {
					
					o = td
					found = true
					break Loop
				}
			}
		}
		if !found {
			log.Printf("protoc-gen-go: WARNING: failed finding publicly imported dependency for %v, used in %v", typeName, *g.file.Name)
		}
	}

	return o
}


type AnnotatedAtoms struct {
	source string
	path   string
	atoms  []interface{}
}



func Annotate(file *FileDescriptor, path string, atoms ...interface{}) *AnnotatedAtoms {
	return &AnnotatedAtoms{source: *file.Name, path: path, atoms: atoms}
}


func (g *Generator) printAtom(v interface{}) {
	switch v := v.(type) {
	case string:
		g.WriteString(v)
	case *string:
		g.WriteString(*v)
	case bool:
		fmt.Fprint(g, v)
	case *bool:
		fmt.Fprint(g, *v)
	case int:
		fmt.Fprint(g, v)
	case *int32:
		fmt.Fprint(g, *v)
	case *int64:
		fmt.Fprint(g, *v)
	case float64:
		fmt.Fprint(g, v)
	case *float64:
		fmt.Fprint(g, *v)
	case GoPackageName:
		g.WriteString(string(v))
	case GoImportPath:
		g.WriteString(strconv.Quote(string(v)))
	default:
		g.Fail(fmt.Sprintf("unknown type in printer: %T", v))
	}
}





func (g *Generator) P(str ...interface{}) {
	if !g.writeOutput {
		return
	}
	g.WriteString(g.indent)
	for _, v := range str {
		switch v := v.(type) {
		case *AnnotatedAtoms:
			begin := int32(g.Len())
			for _, v := range v.atoms {
				g.printAtom(v)
			}
			if g.annotateCode {
				end := int32(g.Len())
				var path []int32
				for _, token := range strings.Split(v.path, ",") {
					val, err := strconv.ParseInt(token, 10, 32)
					if err != nil {
						g.Fail("could not parse proto AST path: ", err.Error())
					}
					path = append(path, int32(val))
				}
				g.annotations = append(g.annotations, &descriptor.GeneratedCodeInfo_Annotation{
					Path:       path,
					SourceFile: &v.source,
					Begin:      &begin,
					End:        &end,
				})
			}
		default:
			g.printAtom(v)
		}
	}
	g.WriteByte('\n')
}



func (g *Generator) addInitf(stmt string, a ...interface{}) {
	g.init = append(g.init, fmt.Sprintf(stmt, a...))
}


func (g *Generator) In() { g.indent += "\t" }


func (g *Generator) Out() {
	if len(g.indent) > 0 {
		g.indent = g.indent[1:]
	}
}


func (g *Generator) GenerateAllFiles() {
	
	for _, p := range plugins {
		p.Init(g)
	}
	
	
	
	genFileMap := make(map[*FileDescriptor]bool, len(g.genFiles))
	for _, file := range g.genFiles {
		genFileMap[file] = true
	}
	for _, file := range g.allFiles {
		g.Reset()
		g.annotations = nil
		g.writeOutput = genFileMap[file]
		g.generate(file)
		if !g.writeOutput {
			continue
		}
		fname := file.goFileName(g.pathType)
		g.Response.File = append(g.Response.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fname),
			Content: proto.String(g.String()),
		})
		if g.annotateCode {
			
			
			g.Response.File = append(g.Response.File, &plugin.CodeGeneratorResponse_File{
				Name:    proto.String(file.goFileName(g.pathType) + ".meta"),
				Content: proto.String(proto.CompactTextString(&descriptor.GeneratedCodeInfo{Annotation: g.annotations})),
			})
		}
	}
}


func (g *Generator) runPlugins(file *FileDescriptor) {
	for _, p := range plugins {
		p.Generate(file)
	}
}



func (g *Generator) generate(file *FileDescriptor) {
	g.file = file
	g.usedPackages = make(map[GoImportPath]bool)
	g.packageNames = make(map[GoImportPath]GoPackageName)
	g.usedPackageNames = make(map[GoPackageName]bool)
	for name := range globalPackageNames {
		g.usedPackageNames[name] = true
	}

	g.P("
	g.P("
	g.P("
	g.P("
	g.P("const _ = ", g.Pkg["proto"], ".ProtoPackageIsVersion", generatedCodeVersion, " 
	g.P()

	for _, td := range g.file.imp {
		g.generateImported(td)
	}
	for _, enum := range g.file.enum {
		g.generateEnum(enum)
	}
	for _, desc := range g.file.desc {
		
		if desc.GetOptions().GetMapEntry() {
			continue
		}
		g.generateMessage(desc)
	}
	for _, ext := range g.file.ext {
		g.generateExtension(ext)
	}
	g.generateInitFunction()

	
	g.runPlugins(file)

	g.generateFileDescriptor(file)

	
	rem := g.Buffer
	remAnno := g.annotations
	g.Buffer = new(bytes.Buffer)
	g.annotations = nil
	g.generateHeader()
	g.generateImports()
	if !g.writeOutput {
		return
	}
	
	for _, anno := range remAnno {
		*anno.Begin += int32(g.Len())
		*anno.End += int32(g.Len())
		g.annotations = append(g.annotations, anno)
	}
	g.Write(rem.Bytes())

	
	fset := token.NewFileSet()
	original := g.Bytes()
	if g.annotateCode {
		
		original = append([]byte(nil), original...)
	}
	ast, err := parser.ParseFile(fset, "", original, parser.ParseComments)
	if err != nil {
		
		
		
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(original))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		g.Fail("bad Go source code was generated:", err.Error(), "\n"+src.String())
	}
	g.Reset()
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(g, fset, ast)
	if err != nil {
		g.Fail("generated Go source code could not be reformatted:", err.Error())
	}
	if g.annotateCode {
		m, err := remap.Compute(original, g.Bytes())
		if err != nil {
			g.Fail("formatted generated Go source code could not be mapped back to the original code:", err.Error())
		}
		for _, anno := range g.annotations {
			new, ok := m.Find(int(*anno.Begin), int(*anno.End))
			if !ok {
				g.Fail("span in formatted generated Go source code could not be mapped back to the original code")
			}
			*anno.Begin = int32(new.Pos)
			*anno.End = int32(new.End)
		}
	}
}


func (g *Generator) generateHeader() {
	g.P("
	if g.file.GetOptions().GetDeprecated() {
		g.P("
	} else {
		g.P("
	}
	g.P()

	importPath, _, _ := g.file.goPackageOption()
	if importPath == "" {
		g.P("package ", g.file.packageName)
	} else {
		g.P("package ", g.file.packageName, " 
	}
	g.P()

	if loc, ok := g.file.comments[strconv.Itoa(packagePath)]; ok {
		g.P(" comment block.
		text := strings.TrimSuffix(loc.GetLeadingComments(), "\n")
		for _, line := range strings.Split(text, "\n") {
			line = strings.TrimPrefix(line, " ")
			
			line = strings.Replace(line, "*/", "* /", -1)
			g.P(line)
		}
		g.P("*/")
		g.P()
	}
}



var deprecationComment = "





func (g *Generator) PrintComments(path string) bool {
	if !g.writeOutput {
		return false
	}
	if c, ok := g.makeComments(path); ok {
		g.P(c)
		return true
	}
	return false
}


func (g *Generator) makeComments(path string) (string, bool) {
	loc, ok := g.file.comments[path]
	if !ok {
		return "", false
	}
	w := new(bytes.Buffer)
	nl := ""
	for _, line := range strings.Split(strings.TrimSuffix(loc.GetLeadingComments(), "\n"), "\n") {
		fmt.Fprintf(w, "%s
		nl = "\n"
	}
	return w.String(), true
}

func (g *Generator) fileByName(filename string) *FileDescriptor {
	return g.allFilesByName[filename]
}


func (g *Generator) weak(i int32) bool {
	for _, j := range g.file.WeakDependency {
		if j == i {
			return true
		}
	}
	return false
}


func (g *Generator) generateImports() {
	
	
	
	g.P("import "+g.Pkg["proto"]+" ", GoImportPath(g.ImportPrefix)+"github.com/golang/protobuf/proto")
	g.P("import " + g.Pkg["fmt"] + ` "fmt"`)
	g.P("import " + g.Pkg["math"] + ` "math"`)
	var (
		imports       = make(map[GoImportPath]bool)
		strongImports = make(map[GoImportPath]bool)
		importPaths   []string
	)
	for i, s := range g.file.Dependency {
		fd := g.fileByName(s)
		importPath := fd.importPath
		
		if importPath == g.file.importPath {
			continue
		}
		if !imports[importPath] {
			importPaths = append(importPaths, string(importPath))
		}
		imports[importPath] = true
		if !g.weak(int32(i)) {
			strongImports[importPath] = true
		}
	}
	sort.Strings(importPaths)
	for i := range importPaths {
		importPath := GoImportPath(importPaths[i])
		packageName := g.GoPackageName(importPath)
		fullPath := GoImportPath(g.ImportPrefix) + importPath
		
		if !strongImports[importPath] {
			g.P("
			continue
		}
		
		
		
		if _, ok := g.usedPackages[importPath]; !ok {
			packageName = "_"
		}
		g.P("import ", packageName, " ", fullPath)
	}
	g.P()
	
	for _, p := range plugins {
		p.GenerateImports(g.file)
		g.P()
	}
	g.P("
	g.P("var _ = ", g.Pkg["proto"], ".Marshal")
	g.P("var _ = ", g.Pkg["fmt"], ".Errorf")
	g.P("var _ = ", g.Pkg["math"], ".Inf")
	g.P()
}

func (g *Generator) generateImported(id *ImportedDescriptor) {
	tn := id.TypeName()
	sn := tn[len(tn)-1]
	df := id.o.File()
	filename := *df.Name
	if df.importPath == g.file.importPath {
		
		g.P("
		g.P()
		return
	}
	if !supportTypeAliases {
		g.Fail(fmt.Sprintf("%s: public imports require at least go1.9", filename))
	}
	g.P("
	g.usedPackages[df.importPath] = true

	for _, sym := range df.exported[id.o] {
		sym.GenerateAlias(g, g.GoPackageName(df.importPath))
	}

	g.P()
}


func (g *Generator) generateEnum(enum *EnumDescriptor) {
	
	typeName := enum.TypeName()
	
	ccTypeName := CamelCaseSlice(typeName)
	ccPrefix := enum.prefix()

	deprecatedEnum := ""
	if enum.GetOptions().GetDeprecated() {
		deprecatedEnum = deprecationComment
	}
	g.PrintComments(enum.path)
	g.P("type ", Annotate(enum.file, enum.path, ccTypeName), " int32", deprecatedEnum)
	g.file.addExport(enum, enumSymbol{ccTypeName, enum.proto3()})
	g.P("const (")
	for i, e := range enum.Value {
		etorPath := fmt.Sprintf("%s,%d,%d", enum.path, enumValuePath, i)
		g.PrintComments(etorPath)

		deprecatedValue := ""
		if e.GetOptions().GetDeprecated() {
			deprecatedValue = deprecationComment
		}

		name := ccPrefix + *e.Name
		g.P(Annotate(enum.file, etorPath, name), " ", ccTypeName, " = ", e.Number, " ", deprecatedValue)
		g.file.addExport(enum, constOrVarSymbol{name, "const", ccTypeName})
	}
	g.P(")")
	g.P("var ", ccTypeName, "_name = map[int32]string{")
	generated := make(map[int32]bool) 
	for _, e := range enum.Value {
		duplicate := ""
		if _, present := generated[*e.Number]; present {
			duplicate = "
		}
		g.P(duplicate, e.Number, ": ", strconv.Quote(*e.Name), ",")
		generated[*e.Number] = true
	}
	g.P("}")
	g.P("var ", ccTypeName, "_value = map[string]int32{")
	for _, e := range enum.Value {
		g.P(strconv.Quote(*e.Name), ": ", e.Number, ",")
	}
	g.P("}")

	if !enum.proto3() {
		g.P("func (x ", ccTypeName, ") Enum() *", ccTypeName, " {")
		g.P("p := new(", ccTypeName, ")")
		g.P("*p = x")
		g.P("return p")
		g.P("}")
	}

	g.P("func (x ", ccTypeName, ") String() string {")
	g.P("return ", g.Pkg["proto"], ".EnumName(", ccTypeName, "_name, int32(x))")
	g.P("}")

	if !enum.proto3() {
		g.P("func (x *", ccTypeName, ") UnmarshalJSON(data []byte) error {")
		g.P("value, err := ", g.Pkg["proto"], ".UnmarshalJSONEnum(", ccTypeName, `_value, data, "`, ccTypeName, `")`)
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P("*x = ", ccTypeName, "(value)")
		g.P("return nil")
		g.P("}")
	}

	var indexes []string
	for m := enum.parent; m != nil; m = m.parent {
		
		indexes = append([]string{strconv.Itoa(m.index)}, indexes...)
	}
	indexes = append(indexes, strconv.Itoa(enum.index))
	g.P("func (", ccTypeName, ") EnumDescriptor() ([]byte, []int) {")
	g.P("return ", g.file.VarName(), ", []int{", strings.Join(indexes, ", "), "}")
	g.P("}")
	if enum.file.GetPackage() == "google.protobuf" && enum.GetName() == "NullValue" {
		g.P("func (", ccTypeName, `) XXX_WellKnownType() string { return "`, enum.GetName(), `" }`)
	}

	g.P()
}














func (g *Generator) goTag(message *Descriptor, field *descriptor.FieldDescriptorProto, wiretype string) string {
	optrepreq := ""
	switch {
	case isOptional(field):
		optrepreq = "opt"
	case isRequired(field):
		optrepreq = "req"
	case isRepeated(field):
		optrepreq = "rep"
	}
	var defaultValue string
	if dv := field.DefaultValue; dv != nil { 
		defaultValue = *dv
		
		switch *field.Type {
		case descriptor.FieldDescriptorProto_TYPE_BOOL:
			if defaultValue == "true" {
				defaultValue = "1"
			} else {
				defaultValue = "0"
			}
		case descriptor.FieldDescriptorProto_TYPE_STRING,
			descriptor.FieldDescriptorProto_TYPE_BYTES:
			
		case descriptor.FieldDescriptorProto_TYPE_ENUM:
			
			obj := g.ObjectNamed(field.GetTypeName())
			if id, ok := obj.(*ImportedDescriptor); ok {
				
				
				obj = id.o
			}
			enum, ok := obj.(*EnumDescriptor)
			if !ok {
				log.Printf("obj is a %T", obj)
				if id, ok := obj.(*ImportedDescriptor); ok {
					log.Printf("id.o is a %T", id.o)
				}
				g.Fail("unknown enum type", CamelCaseSlice(obj.TypeName()))
			}
			defaultValue = enum.integerValueAsString(defaultValue)
		}
		defaultValue = ",def=" + defaultValue
	}
	enum := ""
	if *field.Type == descriptor.FieldDescriptorProto_TYPE_ENUM {
		
		
		obj := g.ObjectNamed(field.GetTypeName())
		if id, ok := obj.(*ImportedDescriptor); ok {
			obj = id.o
		}
		enum = ",enum="
		if pkg := obj.File().GetPackage(); pkg != "" {
			enum += pkg + "."
		}
		enum += CamelCaseSlice(obj.TypeName())
	}
	packed := ""
	if (field.Options != nil && field.Options.GetPacked()) ||
		
		
		(message.proto3() && (field.Options == nil || field.Options.Packed == nil) &&
			isRepeated(field) && isScalar(field)) {
		packed = ",packed"
	}
	fieldName := field.GetName()
	name := fieldName
	if *field.Type == descriptor.FieldDescriptorProto_TYPE_GROUP {
		
		
		
		
		name = *field.TypeName
		if i := strings.LastIndex(name, "."); i >= 0 {
			name = name[i+1:]
		}
	}
	if json := field.GetJsonName(); json != "" && json != name {
		
		
		name += ",json=" + json
	}
	name = ",name=" + name
	if message.proto3() {
		name += ",proto3"
	}
	oneof := ""
	if field.OneofIndex != nil {
		oneof = ",oneof"
	}
	return strconv.Quote(fmt.Sprintf("%s,%d,%s%s%s%s%s%s",
		wiretype,
		field.GetNumber(),
		optrepreq,
		packed,
		name,
		enum,
		oneof,
		defaultValue))
}

func needsStar(typ descriptor.FieldDescriptorProto_Type) bool {
	switch typ {
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return false
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return false
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return false
	}
	return true
}






func (g *Generator) TypeName(obj Object) string {
	return g.DefaultPackageName(obj) + CamelCaseSlice(obj.TypeName())
}


func (g *Generator) GoType(message *Descriptor, field *descriptor.FieldDescriptorProto) (typ string, wire string) {
	
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		typ, wire = "float64", "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		typ, wire = "float32", "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		typ, wire = "int64", "varint"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		typ, wire = "uint64", "varint"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		typ, wire = "int32", "varint"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		typ, wire = "uint32", "varint"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		typ, wire = "uint64", "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		typ, wire = "uint32", "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		typ, wire = "bool", "varint"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		typ, wire = "string", "bytes"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		desc := g.ObjectNamed(field.GetTypeName())
		typ, wire = "*"+g.TypeName(desc), "group"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		desc := g.ObjectNamed(field.GetTypeName())
		typ, wire = "*"+g.TypeName(desc), "bytes"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		typ, wire = "[]byte", "bytes"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		desc := g.ObjectNamed(field.GetTypeName())
		typ, wire = g.TypeName(desc), "varint"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		typ, wire = "int32", "fixed32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		typ, wire = "int64", "fixed64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		typ, wire = "int32", "zigzag32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		typ, wire = "int64", "zigzag64"
	default:
		g.Fail("unknown type for", field.GetName())
	}
	if isRepeated(field) {
		typ = "[]" + typ
	} else if message != nil && message.proto3() {
		return
	} else if field.OneofIndex != nil && message != nil {
		return
	} else if needsStar(*field.Type) {
		typ = "*" + typ
	}
	return
}

func (g *Generator) RecordTypeUse(t string) {
	if _, ok := g.typeNameToObject[t]; ok {
		
		obj := g.ObjectNamed(t)
		g.usedPackages[obj.GoImportPath()] = true
	}
}




var methodNames = [...]string{
	"Reset",
	"String",
	"ProtoMessage",
	"Marshal",
	"Unmarshal",
	"ExtensionRangeArray",
	"ExtensionMap",
	"Descriptor",
}



var wellKnownTypes = map[string]bool{
	"Any":       true,
	"Duration":  true,
	"Empty":     true,
	"Struct":    true,
	"Timestamp": true,

	"Value":       true,
	"ListValue":   true,
	"DoubleValue": true,
	"FloatValue":  true,
	"Int64Value":  true,
	"UInt64Value": true,
	"Int32Value":  true,
	"UInt32Value": true,
	"BoolValue":   true,
	"StringValue": true,
	"BytesValue":  true,
}



func (g *Generator) getterDefault(field *descriptor.FieldDescriptorProto, goMessageType string) string {
	if isRepeated(field) {
		return "nil"
	}
	if def := field.GetDefaultValue(); def != "" {
		defaultConstant := g.defaultConstantName(goMessageType, field.GetName())
		if *field.Type != descriptor.FieldDescriptorProto_TYPE_BYTES {
			return defaultConstant
		}
		return "append([]byte(nil), " + defaultConstant + "...)"
	}
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "false"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return `""`
	case descriptor.FieldDescriptorProto_TYPE_GROUP, descriptor.FieldDescriptorProto_TYPE_MESSAGE, descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "nil"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		obj := g.ObjectNamed(field.GetTypeName())
		var enum *EnumDescriptor
		if id, ok := obj.(*ImportedDescriptor); ok {
			
			enum, _ = id.o.(*EnumDescriptor)
		} else {
			enum, _ = obj.(*EnumDescriptor)
		}
		if enum == nil {
			log.Printf("don't know how to generate getter for %s", field.GetName())
			return "nil"
		}
		if len(enum.Value) == 0 {
			return "0 
		}
		first := enum.Value[0].GetName()
		return g.DefaultPackageName(obj) + enum.prefix() + first
	default:
		return "0"
	}
}



func (g *Generator) defaultConstantName(goMessageType, protoFieldName string) string {
	return "Default_" + goMessageType + "_" + CamelCase(protoFieldName)
}














type msgCtx struct {
	goName  string      
	message *Descriptor 
}


type fieldCommon struct {
	goName     string 
	protoName  string 
	getterName string 
	goType     string 
	tags       string 
	fullPath   string 
}


func (f *fieldCommon) getProtoName() string {
	return f.protoName
}


func (f *fieldCommon) getGoType() string {
	return f.goType
}


type simpleField struct {
	fieldCommon
	protoTypeName string                               
	protoType     descriptor.FieldDescriptorProto_Type 
	deprecated    string                               
	getterDef     string                               
	protoDef      string                               
	comment       string                               
}


func (f *simpleField) decl(g *Generator, mc *msgCtx) {
	g.P(f.comment, Annotate(mc.message.file, f.fullPath, f.goName), "\t", f.goType, "\t`", f.tags, "`", f.deprecated)
}


func (f *simpleField) getter(g *Generator, mc *msgCtx) {
	star := ""
	tname := f.goType
	if needsStar(f.protoType) && tname[0] == '*' {
		tname = tname[1:]
		star = "*"
	}
	if f.deprecated != "" {
		g.P(f.deprecated)
	}
	g.P("func (m *", mc.goName, ") ", Annotate(mc.message.file, f.fullPath, f.getterName), "() "+tname+" {")
	if f.getterDef == "nil" { 
		g.P("if m != nil {")
		g.P("return m." + f.goName)
		g.P("}")
		g.P("return nil")
		g.P("}")
		g.P()
		return
	}
	if mc.message.proto3() {
		g.P("if m != nil {")
	} else {
		g.P("if m != nil && m." + f.goName + " != nil {")
	}
	g.P("return " + star + "m." + f.goName)
	g.P("}")
	g.P("return ", f.getterDef)
	g.P("}")
	g.P()
}


func (f *simpleField) setter(g *Generator, mc *msgCtx) {
	
}


func (f *simpleField) getProtoDef() string {
	return f.protoDef
}


func (f *simpleField) getProtoTypeName() string {
	return f.protoTypeName
}


func (f *simpleField) getProtoType() descriptor.FieldDescriptorProto_Type {
	return f.protoType
}


type oneofSubField struct {
	fieldCommon
	protoTypeName string                               
	protoType     descriptor.FieldDescriptorProto_Type 
	oneofTypeName string                               
	fieldNumber   int                                  
	getterDef     string                               
	protoDef      string                               
}


func (f *oneofSubField) wireTypeName() string {
	switch f.protoType {
	case descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return "WireFixed64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return "WireFixed32"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		return "WireStartGroup"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE,
		descriptor.FieldDescriptorProto_TYPE_STRING,
		descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "WireBytes"
	default: 
		return "WireVarint"
	}
}



func (f *oneofSubField) typedNil(g *Generator) {
	g.P("(*", f.oneofTypeName, ")(nil),")
}


func (f *oneofSubField) marshalCase(g *Generator) {
	g.P("case *", f.oneofTypeName, ":")
	wire := f.wireTypeName()
	var pre, post string
	val := "x." + f.goName 
	switch f.protoType {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		pre = "b.EncodeFixed64(" + g.Pkg["math"] + ".Float64bits("
		post = "))"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		pre = "b.EncodeFixed32(uint64(" + g.Pkg["math"] + ".Float32bits("
		post = ")))"
	case descriptor.FieldDescriptorProto_TYPE_INT64, descriptor.FieldDescriptorProto_TYPE_UINT64:
		pre, post = "b.EncodeVarint(uint64(", "))"
	case descriptor.FieldDescriptorProto_TYPE_INT32, descriptor.FieldDescriptorProto_TYPE_UINT32, descriptor.FieldDescriptorProto_TYPE_ENUM:
		pre, post = "b.EncodeVarint(uint64(", "))"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64, descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		pre, post = "b.EncodeFixed64(uint64(", "))"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32, descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		pre, post = "b.EncodeFixed32(uint64(", "))"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		g.P("t := uint64(0)")
		g.P("if ", val, " { t = 1 }")
		val = "t"
		pre, post = "b.EncodeVarint(", ")"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		pre, post = "b.EncodeStringBytes(", ")"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		pre, post = "b.Marshal(", ")"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		pre, post = "b.EncodeMessage(", ")"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		pre, post = "b.EncodeRawBytes(", ")"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		pre, post = "b.EncodeZigzag32(uint64(", "))"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		pre, post = "b.EncodeZigzag64(uint64(", "))"
	default:
		g.Fail("unhandled oneof field type ", f.protoType.String())
	}
	g.P("b.EncodeVarint(", f.fieldNumber, "<<3|", g.Pkg["proto"], ".", wire, ")")
	if t := f.protoType; t != descriptor.FieldDescriptorProto_TYPE_GROUP && t != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		g.P(pre, val, post)
	} else {
		g.P("if err := ", pre, val, post, "; err != nil {")
		g.P("return err")
		g.P("}")
	}
	if f.protoType == descriptor.FieldDescriptorProto_TYPE_GROUP {
		g.P("b.EncodeVarint(", f.fieldNumber, "<<3|", g.Pkg["proto"], ".WireEndGroup)")
	}
}


func (f *oneofSubField) unmarshalCase(g *Generator, origOneofName string, oneofName string) {
	g.P("case ", f.fieldNumber, ": 
	g.P("if wire != ", g.Pkg["proto"], ".", f.wireTypeName(), " {")
	g.P("return true, ", g.Pkg["proto"], ".ErrInternalBadWireType")
	g.P("}")
	lhs := "x, err" 
	var dec, cast, cast2 string
	switch f.protoType {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		dec, cast = "b.DecodeFixed64()", g.Pkg["math"]+".Float64frombits"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		dec, cast, cast2 = "b.DecodeFixed32()", "uint32", g.Pkg["math"]+".Float32frombits"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		dec, cast = "b.DecodeVarint()", "int64"
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		dec = "b.DecodeVarint()"
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		dec, cast = "b.DecodeVarint()", "int32"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		dec = "b.DecodeFixed64()"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		dec, cast = "b.DecodeFixed32()", "uint32"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		dec = "b.DecodeVarint()"
		
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		dec = "b.DecodeStringBytes()"
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		g.P("msg := new(", f.goType[1:], ")") 
		lhs = "err"
		dec = "b.DecodeGroup(msg)"
		
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		g.P("msg := new(", f.goType[1:], ")") 
		lhs = "err"
		dec = "b.DecodeMessage(msg)"
		
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		dec = "b.DecodeRawBytes(true)"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		dec, cast = "b.DecodeVarint()", "uint32"
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		dec, cast = "b.DecodeVarint()", f.goType
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		dec, cast = "b.DecodeFixed32()", "int32"
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		dec, cast = "b.DecodeFixed64()", "int64"
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		dec, cast = "b.DecodeZigzag32()", "int32"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		dec, cast = "b.DecodeZigzag64()", "int64"
	default:
		g.Fail("unhandled oneof field type ", f.protoType.String())
	}
	g.P(lhs, " := ", dec)
	val := "x"
	if cast != "" {
		val = cast + "(" + val + ")"
	}
	if cast2 != "" {
		val = cast2 + "(" + val + ")"
	}
	switch f.protoType {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		val += " != 0"
	case descriptor.FieldDescriptorProto_TYPE_GROUP,
		descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		val = "msg"
	}
	g.P("m.", oneofName, " = &", f.oneofTypeName, "{", val, "}")
	g.P("return true, err")
}


func (f *oneofSubField) sizerCase(g *Generator) {
	g.P("case *", f.oneofTypeName, ":")
	val := "x." + f.goName
	var varint, fixed string
	switch f.protoType {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		fixed = "8"
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		fixed = "4"
	case descriptor.FieldDescriptorProto_TYPE_INT64, descriptor.FieldDescriptorProto_TYPE_UINT64, descriptor.FieldDescriptorProto_TYPE_INT32, descriptor.FieldDescriptorProto_TYPE_UINT32, descriptor.FieldDescriptorProto_TYPE_ENUM:
		varint = val
	case descriptor.FieldDescriptorProto_TYPE_FIXED64, descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		fixed = "8"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32, descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		fixed = "4"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		fixed = "1"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		fixed = "len(" + val + ")"
		varint = fixed
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		fixed = g.Pkg["proto"] + ".Size(" + val + ")"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		g.P("s := ", g.Pkg["proto"], ".Size(", val, ")")
		fixed = "s"
		varint = fixed
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		fixed = "len(" + val + ")"
		varint = fixed
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		varint = "(uint32(" + val + ") << 1) ^ uint32((int32(" + val + ") >> 31))"
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		varint = "uint64(" + val + " << 1) ^ uint64((int64(" + val + ") >> 63))"
	default:
		g.Fail("unhandled oneof field type ", f.protoType.String())
	}
	
	
	tagAndWireSize := proto.SizeVarint(uint64(f.fieldNumber << 3)) 
	g.P("n += ", tagAndWireSize, " 
	if varint != "" {
		g.P("n += ", g.Pkg["proto"], ".SizeVarint(uint64(", varint, "))")
	}
	if fixed != "" {
		g.P("n += ", fixed)
	}
	if f.protoType == descriptor.FieldDescriptorProto_TYPE_GROUP {
		g.P("n += ", tagAndWireSize, " 
	}
}


func (f *oneofSubField) getProtoDef() string {
	return f.protoDef
}


func (f *oneofSubField) getProtoTypeName() string {
	return f.protoTypeName
}


func (f *oneofSubField) getProtoType() descriptor.FieldDescriptorProto_Type {
	return f.protoType
}



type oneofField struct {
	fieldCommon
	subFields []*oneofSubField 
	comment   string           
}


func (f *oneofField) decl(g *Generator, mc *msgCtx) {
	comment := f.comment
	for _, sf := range f.subFields {
		comment += "
	}
	g.P(comment, Annotate(mc.message.file, f.fullPath, f.goName), " ", f.goType, " `", f.tags, "`")
}



func (f *oneofField) getter(g *Generator, mc *msgCtx) {
	
	g.P("type ", f.goType, " interface {")
	g.P(f.goType, "()")
	g.P("}")
	g.P()
	
	for _, sf := range f.subFields {
		g.P("type ", Annotate(mc.message.file, sf.fullPath, sf.oneofTypeName), " struct {")
		g.P(Annotate(mc.message.file, sf.fullPath, sf.goName), " ", sf.goType, " `", sf.tags, "`")
		g.P("}")
		g.P()
	}
	for _, sf := range f.subFields {
		g.P("func (*", sf.oneofTypeName, ") ", f.goType, "() {}")
		g.P()
	}
	
	g.P("func (m *", mc.goName, ") ", Annotate(mc.message.file, f.fullPath, f.getterName), "() ", f.goType, " {")
	g.P("if m != nil { return m.", f.goName, " }")
	g.P("return nil")
	g.P("}")
	g.P()
	
	for _, sf := range f.subFields {
		g.P("func (m *", mc.goName, ") ", Annotate(mc.message.file, sf.fullPath, sf.getterName), "() "+sf.goType+" {")
		g.P("if x, ok := m.", f.getterName, "().(*", sf.oneofTypeName, "); ok {")
		g.P("return x.", sf.goName)
		g.P("}")
		g.P("return ", sf.getterDef)
		g.P("}")
		g.P()
	}
}


func (f *oneofField) setter(g *Generator, mc *msgCtx) {
	
}


type topLevelField interface {
	decl(g *Generator, mc *msgCtx)   
	getter(g *Generator, mc *msgCtx) 
	setter(g *Generator, mc *msgCtx) 
}


type defField interface {
	getProtoDef() string                                
	getProtoName() string                               
	getGoType() string                                  
	getProtoTypeName() string                           
	getProtoType() descriptor.FieldDescriptorProto_Type 
}



func (g *Generator) generateDefaultConstants(mc *msgCtx, topLevelFields []topLevelField) {
	
	dFields := []defField{}
	for _, pf := range topLevelFields {
		if f, ok := pf.(*oneofField); ok {
			for _, osf := range f.subFields {
				dFields = append(dFields, osf)
			}
			continue
		}
		dFields = append(dFields, pf.(defField))
	}
	for _, df := range dFields {
		def := df.getProtoDef()
		if def == "" {
			continue
		}
		fieldname := g.defaultConstantName(mc.goName, df.getProtoName())
		typename := df.getGoType()
		if typename[0] == '*' {
			typename = typename[1:]
		}
		kind := "const "
		switch {
		case typename == "bool":
		case typename == "string":
			def = strconv.Quote(def)
		case typename == "[]byte":
			def = "[]byte(" + strconv.Quote(unescape(def)) + ")"
			kind = "var "
		case def == "inf", def == "-inf", def == "nan":
			
			switch def {
			case "inf":
				def = "math.Inf(1)"
			case "-inf":
				def = "math.Inf(-1)"
			case "nan":
				def = "math.NaN()"
			}
			if df.getProtoType() == descriptor.FieldDescriptorProto_TYPE_FLOAT {
				def = "float32(" + def + ")"
			}
			kind = "var "
		case df.getProtoType() == descriptor.FieldDescriptorProto_TYPE_ENUM:
			
			obj := g.ObjectNamed(df.getProtoTypeName())
			var enum *EnumDescriptor
			if id, ok := obj.(*ImportedDescriptor); ok {
				
				enum, _ = id.o.(*EnumDescriptor)
			} else {
				enum, _ = obj.(*EnumDescriptor)
			}
			if enum == nil {
				log.Printf("don't know how to generate constant for %s", fieldname)
				continue
			}
			def = g.DefaultPackageName(obj) + enum.prefix() + def
		}
		g.P(kind, fieldname, " ", typename, " = ", def)
		g.file.addExport(mc.message, constOrVarSymbol{fieldname, kind, ""})
	}
	g.P()
}


func (g *Generator) generateInternalStructFields(mc *msgCtx, topLevelFields []topLevelField) {
	g.P("XXX_NoUnkeyedLiteral\tstruct{} `json:\"-\"`") 
	if len(mc.message.ExtensionRange) > 0 {
		messageset := ""
		if opts := mc.message.Options; opts != nil && opts.GetMessageSetWireFormat() {
			messageset = "protobuf_messageset:\"1\" "
		}
		g.P(g.Pkg["proto"], ".XXX_InternalExtensions `", messageset, "json:\"-\"`")
	}
	g.P("XXX_unrecognized\t[]byte `json:\"-\"`")
	g.P("XXX_sizecache\tint32 `json:\"-\"`")

}


func (g *Generator) generateOneofFuncs(mc *msgCtx, topLevelFields []topLevelField) {
	ofields := []*oneofField{}
	for _, f := range topLevelFields {
		if o, ok := f.(*oneofField); ok {
			ofields = append(ofields, o)
		}
	}
	if len(ofields) == 0 {
		return
	}
	enc := "_" + mc.goName + "_OneofMarshaler"
	dec := "_" + mc.goName + "_OneofUnmarshaler"
	size := "_" + mc.goName + "_OneofSizer"
	encSig := "(msg " + g.Pkg["proto"] + ".Message, b *" + g.Pkg["proto"] + ".Buffer) error"
	decSig := "(msg " + g.Pkg["proto"] + ".Message, tag, wire int, b *" + g.Pkg["proto"] + ".Buffer) (bool, error)"
	sizeSig := "(msg " + g.Pkg["proto"] + ".Message) (n int)"

	
	g.P("
	g.P("func (*", mc.goName, ") XXX_OneofFuncs() (func", encSig, ", func", decSig, ", func", sizeSig, ", []interface{}) {")
	g.P("return ", enc, ", ", dec, ", ", size, ", []interface{}{")
	for _, of := range ofields {
		for _, sf := range of.subFields {
			sf.typedNil(g)
		}
	}
	g.P("}")
	g.P("}")
	g.P()

	
	g.P("func ", enc, encSig, " {")
	g.P("m := msg.(*", mc.goName, ")")
	for _, of := range ofields {
		g.P("
		g.P("switch x := m.", of.goName, ".(type) {")
		for _, sf := range of.subFields {
			
			sf.marshalCase(g)
		}
		g.P("case nil:")
		g.P("default:")
		g.P(" return ", g.Pkg["fmt"], `.Errorf("`, mc.goName, ".", of.goName, ` has unexpected type %T", x)`)
		g.P("}")
	}
	g.P("return nil")
	g.P("}")
	g.P()

	
	g.P("func ", dec, decSig, " {")
	g.P("m := msg.(*", mc.goName, ")")
	g.P("switch tag {")
	for _, of := range ofields {
		for _, sf := range of.subFields {
			sf.unmarshalCase(g, of.getProtoName(), of.goName)
		}
	}
	g.P("default:")
	g.P("return false, nil")
	g.P("}")
	g.P("}")
	g.P()

	
	g.P("func ", size, sizeSig, " {")
	g.P("m := msg.(*", mc.goName, ")")
	for _, of := range ofields {
		g.P("
		g.P("switch x := m.", of.goName, ".(type) {")
		for _, sf := range of.subFields {
			
			sf.sizerCase(g)
		}
		g.P("case nil:")
		g.P("default:")
		g.P("panic(", g.Pkg["fmt"], ".Sprintf(\"proto: unexpected type %T in oneof\", x))")
		g.P("}")
	}
	g.P("return n")
	g.P("}")
	g.P()
}


func (g *Generator) generateMessageStruct(mc *msgCtx, topLevelFields []topLevelField) {
	comments := g.PrintComments(mc.message.path)

	
	if mc.message.GetOptions().GetDeprecated() {
		if comments {
			
			
			g.P("
		}
		g.P(deprecationComment)
	}

	g.P("type ", Annotate(mc.message.file, mc.message.path, mc.goName), " struct {")
	for _, pf := range topLevelFields {
		pf.decl(g, mc)
	}
	g.generateInternalStructFields(mc, topLevelFields)
	g.P("}")
}


func (g *Generator) generateGetters(mc *msgCtx, topLevelFields []topLevelField) {
	for _, pf := range topLevelFields {
		pf.getter(g, mc)
	}
}


func (g *Generator) generateSetters(mc *msgCtx, topLevelFields []topLevelField) {
	for _, pf := range topLevelFields {
		pf.setter(g, mc)
	}
}


func (g *Generator) generateCommonMethods(mc *msgCtx) {
	
	g.P("func (m *", mc.goName, ") Reset() { *m = ", mc.goName, "{} }")
	g.P("func (m *", mc.goName, ") String() string { return ", g.Pkg["proto"], ".CompactTextString(m) }")
	g.P("func (*", mc.goName, ") ProtoMessage() {}")
	var indexes []string
	for m := mc.message; m != nil; m = m.parent {
		indexes = append([]string{strconv.Itoa(m.index)}, indexes...)
	}
	g.P("func (*", mc.goName, ") Descriptor() ([]byte, []int) {")
	g.P("return ", g.file.VarName(), ", []int{", strings.Join(indexes, ", "), "}")
	g.P("}")
	
	
	if mc.message.file.GetPackage() == "google.protobuf" && wellKnownTypes[mc.message.GetName()] {
		g.P("func (*", mc.goName, `) XXX_WellKnownType() string { return "`, mc.message.GetName(), `" }`)
	}

	
	if len(mc.message.ExtensionRange) > 0 {
		
		if opts := mc.message.Options; opts != nil && opts.GetMessageSetWireFormat() {
			g.P()
			g.P("func (m *", mc.goName, ") MarshalJSON() ([]byte, error) {")
			g.P("return ", g.Pkg["proto"], ".MarshalMessageSetJSON(&m.XXX_InternalExtensions)")
			g.P("}")
			g.P("func (m *", mc.goName, ") UnmarshalJSON(buf []byte) error {")
			g.P("return ", g.Pkg["proto"], ".UnmarshalMessageSetJSON(buf, &m.XXX_InternalExtensions)")
			g.P("}")
		}

		g.P()
		g.P("var extRange_", mc.goName, " = []", g.Pkg["proto"], ".ExtensionRange{")
		for _, r := range mc.message.ExtensionRange {
			end := fmt.Sprint(*r.End - 1) 
			g.P("{Start: ", r.Start, ", End: ", end, "},")
		}
		g.P("}")
		g.P("func (*", mc.goName, ") ExtensionRangeArray() []", g.Pkg["proto"], ".ExtensionRange {")
		g.P("return extRange_", mc.goName)
		g.P("}")
	}

	
	
	
	
	

	
	g.P("func (m *", mc.goName, ") XXX_Unmarshal(b []byte) error {")
	g.P("return xxx_messageInfo_", mc.goName, ".Unmarshal(m, b)")
	g.P("}")

	g.P("func (m *", mc.goName, ") XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {")
	g.P("return xxx_messageInfo_", mc.goName, ".Marshal(b, m, deterministic)")
	g.P("}")

	g.P("func (dst *", mc.goName, ") XXX_Merge(src ", g.Pkg["proto"], ".Message) {")
	g.P("xxx_messageInfo_", mc.goName, ".Merge(dst, src)")
	g.P("}")

	g.P("func (m *", mc.goName, ") XXX_Size() int {") 
	g.P("return xxx_messageInfo_", mc.goName, ".Size(m)")
	g.P("}")

	g.P("func (m *", mc.goName, ") XXX_DiscardUnknown() {")
	g.P("xxx_messageInfo_", mc.goName, ".DiscardUnknown(m)")
	g.P("}")

	g.P("var xxx_messageInfo_", mc.goName, " ", g.Pkg["proto"], ".InternalMessageInfo")
	g.P()
}


func (g *Generator) generateMessage(message *Descriptor) {
	topLevelFields := []topLevelField{}
	oFields := make(map[int32]*oneofField)
	
	typeName := message.TypeName()
	
	goTypeName := CamelCaseSlice(typeName)

	usedNames := make(map[string]bool)
	for _, n := range methodNames {
		usedNames[n] = true
	}

	
	
	
	allocNames := func(ns ...string) []string {
	Loop:
		for {
			for _, n := range ns {
				if usedNames[n] {
					for i := range ns {
						ns[i] += "_"
					}
					continue Loop
				}
			}
			for _, n := range ns {
				usedNames[n] = true
			}
			return ns
		}
	}

	mapFieldTypes := make(map[*descriptor.FieldDescriptorProto]string) 

	
	for i, field := range message.Field {
		
		
		
		
		
		base := CamelCase(*field.Name)
		ns := allocNames(base, "Get"+base)
		fieldName, fieldGetterName := ns[0], ns[1]
		typename, wiretype := g.GoType(message, field)
		jsonName := *field.Name
		tag := fmt.Sprintf("protobuf:%s json:%q", g.goTag(message, field, wiretype), jsonName+",omitempty")

		oneof := field.OneofIndex != nil
		if oneof && oFields[*field.OneofIndex] == nil {
			odp := message.OneofDecl[int(*field.OneofIndex)]
			base := CamelCase(odp.GetName())
			names := allocNames(base, "Get"+base)
			fname, gname := names[0], names[1]

			
			
			oneofFullPath := fmt.Sprintf("%s,%d,%d", message.path, messageOneofPath, *field.OneofIndex)
			c, ok := g.makeComments(oneofFullPath)
			if ok {
				c += "\n
			}
			c += "
			
			

			dname := "is" + goTypeName + "_" + fname
			tag := `protobuf_oneof:"` + odp.GetName() + `"`
			of := oneofField{
				fieldCommon: fieldCommon{
					goName:     fname,
					getterName: gname,
					goType:     dname,
					tags:       tag,
					protoName:  odp.GetName(),
					fullPath:   oneofFullPath,
				},
				comment: c,
			}
			topLevelFields = append(topLevelFields, &of)
			oFields[*field.OneofIndex] = &of
		}

		if *field.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			desc := g.ObjectNamed(field.GetTypeName())
			if d, ok := desc.(*Descriptor); ok && d.GetOptions().GetMapEntry() {
				
				keyField, valField := d.Field[0], d.Field[1]
				keyType, keyWire := g.GoType(d, keyField)
				valType, valWire := g.GoType(d, valField)
				keyTag, valTag := g.goTag(d, keyField, keyWire), g.goTag(d, valField, valWire)

				
				
				
				keyType = strings.TrimPrefix(keyType, "*")
				switch *valField.Type {
				case descriptor.FieldDescriptorProto_TYPE_ENUM:
					valType = strings.TrimPrefix(valType, "*")
					g.RecordTypeUse(valField.GetTypeName())
				case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
					g.RecordTypeUse(valField.GetTypeName())
				default:
					valType = strings.TrimPrefix(valType, "*")
				}

				typename = fmt.Sprintf("map[%s]%s", keyType, valType)
				mapFieldTypes[field] = typename 

				tag += fmt.Sprintf(" protobuf_key:%s protobuf_val:%s", keyTag, valTag)
			}
		}

		dvalue := g.getterDefault(field, goTypeName)
		if oneof {
			tname := goTypeName + "_" + fieldName
			
			
			for {
				ok := true
				for _, desc := range message.nested {
					if CamelCaseSlice(desc.TypeName()) == tname {
						ok = false
						break
					}
				}
				for _, enum := range message.enums {
					if CamelCaseSlice(enum.TypeName()) == tname {
						ok = false
						break
					}
				}
				if !ok {
					tname += "_"
					continue
				}
				break
			}

			oneofField := oFields[*field.OneofIndex]
			tag := "protobuf:" + g.goTag(message, field, wiretype)
			sf := oneofSubField{
				fieldCommon: fieldCommon{
					goName:     fieldName,
					getterName: fieldGetterName,
					goType:     typename,
					tags:       tag,
					protoName:  field.GetName(),
					fullPath:   fmt.Sprintf("%s,%d,%d", message.path, messageFieldPath, i),
				},
				protoTypeName: field.GetTypeName(),
				fieldNumber:   int(*field.Number),
				protoType:     *field.Type,
				getterDef:     dvalue,
				protoDef:      field.GetDefaultValue(),
				oneofTypeName: tname,
			}
			oneofField.subFields = append(oneofField.subFields, &sf)
			g.RecordTypeUse(field.GetTypeName())
			continue
		}

		fieldDeprecated := ""
		if field.GetOptions().GetDeprecated() {
			fieldDeprecated = deprecationComment
		}

		fieldFullPath := fmt.Sprintf("%s,%d,%d", message.path, messageFieldPath, i)
		c, ok := g.makeComments(fieldFullPath)
		if ok {
			c += "\n"
		}
		rf := simpleField{
			fieldCommon: fieldCommon{
				goName:     fieldName,
				getterName: fieldGetterName,
				goType:     typename,
				tags:       tag,
				protoName:  field.GetName(),
				fullPath:   fieldFullPath,
			},
			protoTypeName: field.GetTypeName(),
			protoType:     *field.Type,
			deprecated:    fieldDeprecated,
			getterDef:     dvalue,
			protoDef:      field.GetDefaultValue(),
			comment:       c,
		}
		var pf topLevelField = &rf

		topLevelFields = append(topLevelFields, pf)
		g.RecordTypeUse(field.GetTypeName())
	}

	mc := &msgCtx{
		goName:  goTypeName,
		message: message,
	}

	g.generateMessageStruct(mc, topLevelFields)
	g.P()
	g.generateCommonMethods(mc)
	g.P()
	g.generateDefaultConstants(mc, topLevelFields)
	g.P()
	g.generateGetters(mc, topLevelFields)
	g.P()
	g.generateSetters(mc, topLevelFields)
	g.P()
	g.generateOneofFuncs(mc, topLevelFields)
	g.P()

	if !message.group {

		var oneofTypes []string
		for _, f := range topLevelFields {
			if of, ok := f.(*oneofField); ok {
				for _, osf := range of.subFields {
					oneofTypes = append(oneofTypes, osf.oneofTypeName)
				}
			}
		}

		opts := message.Options
		ms := &messageSymbol{
			sym:           goTypeName,
			hasExtensions: len(message.ExtensionRange) > 0,
			isMessageSet:  opts != nil && opts.GetMessageSetWireFormat(),
			oneofTypes:    oneofTypes,
		}
		g.file.addExport(message, ms)
	}

	for _, ext := range message.ext {
		g.generateExtension(ext)
	}

	fullName := strings.Join(message.TypeName(), ".")
	if g.file.Package != nil {
		fullName = *g.file.Package + "." + fullName
	}

	g.addInitf("%s.RegisterType((*%s)(nil), %q)", g.Pkg["proto"], goTypeName, fullName)
	
	for _, k := range mapFieldKeys(mapFieldTypes) {
		fullName := strings.TrimPrefix(*k.TypeName, ".")
		g.addInitf("%s.RegisterMapType((%s)(nil), %q)", g.Pkg["proto"], mapFieldTypes[k], fullName)
	}

}

type byTypeName []*descriptor.FieldDescriptorProto

func (a byTypeName) Len() int           { return len(a) }
func (a byTypeName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTypeName) Less(i, j int) bool { return *a[i].TypeName < *a[j].TypeName }


func mapFieldKeys(m map[*descriptor.FieldDescriptorProto]string) []*descriptor.FieldDescriptorProto {
	keys := make([]*descriptor.FieldDescriptorProto, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(byTypeName(keys))
	return keys
}

var escapeChars = [256]byte{
	'a': '\a', 'b': '\b', 'f': '\f', 'n': '\n', 'r': '\r', 't': '\t', 'v': '\v', '\\': '\\', '"': '"', '\'': '\'', '?': '?',
}




func unescape(s string) string {
	
	
	

	var out []byte
	for len(s) > 0 {
		
		if s[0] != '\\' || len(s) < 2 {
			out = append(out, s[0])
			s = s[1:]
		} else if c := escapeChars[s[1]]; c != 0 {
			
			out = append(out, c)
			s = s[2:]
		} else if s[1] == 'x' || s[1] == 'X' {
			
			if len(s) < 4 {
				
				out = append(out, s[:2]...)
				s = s[2:]
				continue
			}
			v, err := strconv.ParseUint(s[2:4], 16, 8)
			if err != nil {
				out = append(out, s[:4]...)
			} else {
				out = append(out, byte(v))
			}
			s = s[4:]
		} else if '0' <= s[1] && s[1] <= '7' {
			
			
			n := len(s[1:]) - len(strings.TrimLeft(s[1:], "01234567"))
			if n > 3 {
				n = 3
			}
			v, err := strconv.ParseUint(s[1:1+n], 8, 8)
			if err != nil {
				out = append(out, s[:1+n]...)
			} else {
				out = append(out, byte(v))
			}
			s = s[1+n:]
		} else {
			
			out = append(out, s[0])
			s = s[1:]
		}
	}

	return string(out)
}

func (g *Generator) generateExtension(ext *ExtensionDescriptor) {
	ccTypeName := ext.DescName()

	extObj := g.ObjectNamed(*ext.Extendee)
	var extDesc *Descriptor
	if id, ok := extObj.(*ImportedDescriptor); ok {
		
		
		extDesc = id.o.(*Descriptor)
	} else {
		extDesc = extObj.(*Descriptor)
	}
	extendedType := "*" + g.TypeName(extObj) 
	field := ext.FieldDescriptorProto
	fieldType, wireType := g.GoType(ext.parent, field)
	tag := g.goTag(extDesc, field, wireType)
	g.RecordTypeUse(*ext.Extendee)
	if n := ext.FieldDescriptorProto.TypeName; n != nil {
		
		g.RecordTypeUse(*n)
	}

	typeName := ext.TypeName()

	
	
	
	
	
	
	
	
	mset := false
	if extDesc.GetOptions().GetMessageSetWireFormat() && typeName[len(typeName)-1] == "message_set_extension" {
		typeName = typeName[:len(typeName)-1]
		mset = true
	}

	
	
	extName := strings.Join(typeName, ".")
	if g.file.Package != nil {
		extName = *g.file.Package + "." + extName
	}

	g.P("var ", ccTypeName, " = &", g.Pkg["proto"], ".ExtensionDesc{")
	g.P("ExtendedType: (", extendedType, ")(nil),")
	g.P("ExtensionType: (", fieldType, ")(nil),")
	g.P("Field: ", field.Number, ",")
	g.P(`Name: "`, extName, `",`)
	g.P("Tag: ", tag, ",")
	g.P(`Filename: "`, g.file.GetName(), `",`)

	g.P("}")
	g.P()

	if mset {
		
		g.addInitf("%s.RegisterMessageSetType((%s)(nil), %d, %q)", g.Pkg["proto"], fieldType, *field.Number, extName)
	}

	g.file.addExport(ext, constOrVarSymbol{ccTypeName, "var", ""})
}

func (g *Generator) generateInitFunction() {
	for _, enum := range g.file.enum {
		g.generateEnumRegistration(enum)
	}
	for _, d := range g.file.desc {
		for _, ext := range d.ext {
			g.generateExtensionRegistration(ext)
		}
	}
	for _, ext := range g.file.ext {
		g.generateExtensionRegistration(ext)
	}
	if len(g.init) == 0 {
		return
	}
	g.P("func init() {")
	for _, l := range g.init {
		g.P(l)
	}
	g.P("}")
	g.init = nil
}

func (g *Generator) generateFileDescriptor(file *FileDescriptor) {
	
	
	pb := proto.Clone(file.FileDescriptorProto).(*descriptor.FileDescriptorProto)
	pb.SourceCodeInfo = nil

	b, err := proto.Marshal(pb)
	if err != nil {
		g.Fail(err.Error())
	}

	var buf bytes.Buffer
	w, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	w.Write(b)
	w.Close()
	b = buf.Bytes()

	v := file.VarName()
	g.P()
	g.P("func init() { ", g.Pkg["proto"], ".RegisterFile(", strconv.Quote(*file.Name), ", ", v, ") }")
	g.P("var ", v, " = []byte{")
	g.P("
	for len(b) > 0 {
		n := 16
		if n > len(b) {
			n = len(b)
		}

		s := ""
		for _, c := range b[:n] {
			s += fmt.Sprintf("0x%02x,", c)
		}
		g.P(s)

		b = b[n:]
	}
	g.P("}")
}

func (g *Generator) generateEnumRegistration(enum *EnumDescriptor) {
	
	pkg := enum.File().GetPackage()
	if pkg != "" {
		pkg += "."
	}
	
	typeName := enum.TypeName()
	
	ccTypeName := CamelCaseSlice(typeName)
	g.addInitf("%s.RegisterEnum(%q, %[3]s_name, %[3]s_value)", g.Pkg["proto"], pkg+ccTypeName, ccTypeName)
}

func (g *Generator) generateExtensionRegistration(ext *ExtensionDescriptor) {
	g.addInitf("%s.RegisterExtension(%s)", g.Pkg["proto"], ext.DescName())
}




func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}


func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}









func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		
		t = append(t, 'X')
		i++
	}
	
	
	
	
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue 
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		
		
		if isASCIILower(c) {
			c ^= ' ' 
		}
		t = append(t, c) 
		
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}



func CamelCaseSlice(elem []string) string { return CamelCase(strings.Join(elem, "_")) }


func dottedSlice(elem []string) string { return strings.Join(elem, ".") }


func isOptional(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_OPTIONAL
}


func isRequired(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REQUIRED
}


func isRepeated(field *descriptor.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED
}


func isScalar(field *descriptor.FieldDescriptorProto) bool {
	if field.Type == nil {
		return false
	}
	switch *field.Type {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_BOOL,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_ENUM,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		return true
	default:
		return false
	}
}




func badToUnderscore(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		return r
	}
	return '_'
}


func baseName(name string) string {
	
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[0:i]
	}
	return name
}









const (
	
	packagePath = 2 
	messagePath = 4 
	enumPath    = 5 
	
	messageFieldPath   = 2 
	messageMessagePath = 3 
	messageEnumPath    = 4 
	messageOneofPath   = 8 
	
	enumValuePath = 2 
)

var supportTypeAliases bool

func init() {
	for _, tag := range build.Default.ReleaseTags {
		if tag == "go1.9" {
			supportTypeAliases = true
			return
		}
	}
}
