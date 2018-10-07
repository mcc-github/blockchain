

































package grpc

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)





const generatedCodeVersion = 4



const (
	contextPkgPath = "golang.org/x/net/context"
	grpcPkgPath    = "google.golang.org/grpc"
)

func init() {
	generator.RegisterPlugin(new(grpc))
}



type grpc struct {
	gen *generator.Generator
}


func (g *grpc) Name() string {
	return "grpc"
}




var (
	contextPkg string
	grpcPkg    string
)


func (g *grpc) Init(gen *generator.Generator) {
	g.gen = gen
	contextPkg = generator.RegisterUniquePackageName("context", nil)
	grpcPkg = generator.RegisterUniquePackageName("grpc", nil)
}



func (g *grpc) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}


func (g *grpc) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}


func (g *grpc) P(args ...interface{}) { g.gen.P(args...) }


func (g *grpc) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	g.P("
	g.P("var _ ", contextPkg, ".Context")
	g.P("var _ ", grpcPkg, ".ClientConn")
	g.P()

	
	g.P("
	g.P("
	g.P("const _ = ", grpcPkg, ".SupportPackageIsVersion", generatedCodeVersion)
	g.P()

	for i, service := range file.FileDescriptorProto.Service {
		g.generateService(file, service, i)
	}
}


func (g *grpc) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	g.P("import (")
	g.P(contextPkg, " ", generator.GoImportPath(path.Join(string(g.gen.ImportPrefix), contextPkgPath)))
	g.P(grpcPkg, " ", generator.GoImportPath(path.Join(string(g.gen.ImportPrefix), grpcPkgPath)))
	g.P(")")
	g.P()
}


var reservedClientName = map[string]bool{
	
}

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }



var deprecationComment = "


func (g *grpc) generateService(file *generator.FileDescriptor, service *pb.ServiceDescriptorProto, index int) {
	path := fmt.Sprintf("6,%d", index) 

	origServName := service.GetName()
	fullServName := origServName
	if pkg := file.GetPackage(); pkg != "" {
		fullServName = pkg + "." + fullServName
	}
	servName := generator.CamelCase(origServName)
	deprecated := service.GetOptions().GetDeprecated()

	g.P()
	g.P(fmt.Sprintf(`



	
	if deprecated {
		g.P("
		g.P(deprecationComment)
	}
	g.P("type ", servName, "Client interface {")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) 
		g.P(g.generateClientSignature(servName, method))
	}
	g.P("}")
	g.P()

	
	g.P("type ", unexport(servName), "Client struct {")
	g.P("cc *", grpcPkg, ".ClientConn")
	g.P("}")
	g.P()

	
	if deprecated {
		g.P(deprecationComment)
	}
	g.P("func New", servName, "Client (cc *", grpcPkg, ".ClientConn) ", servName, "Client {")
	g.P("return &", unexport(servName), "Client{cc}")
	g.P("}")
	g.P()

	var methodIndex, streamIndex int
	serviceDescVar := "_" + servName + "_serviceDesc"
	
	for _, method := range service.Method {
		var descExpr string
		if !method.GetServerStreaming() && !method.GetClientStreaming() {
			
			descExpr = fmt.Sprintf("&%s.Methods[%d]", serviceDescVar, methodIndex)
			methodIndex++
		} else {
			
			descExpr = fmt.Sprintf("&%s.Streams[%d]", serviceDescVar, streamIndex)
			streamIndex++
		}
		g.generateClientMethod(servName, fullServName, serviceDescVar, method, descExpr)
	}

	
	serverType := servName + "Server"
	g.P("
	if deprecated {
		g.P("
		g.P(deprecationComment)
	}
	g.P("type ", serverType, " interface {")
	for i, method := range service.Method {
		g.gen.PrintComments(fmt.Sprintf("%s,2,%d", path, i)) 
		g.P(g.generateServerSignature(servName, method))
	}
	g.P("}")
	g.P()

	
	if deprecated {
		g.P(deprecationComment)
	}
	g.P("func Register", servName, "Server(s *", grpcPkg, ".Server, srv ", serverType, ") {")
	g.P("s.RegisterService(&", serviceDescVar, `, srv)`)
	g.P("}")
	g.P()

	
	var handlerNames []string
	for _, method := range service.Method {
		hname := g.generateServerMethod(servName, fullServName, method)
		handlerNames = append(handlerNames, hname)
	}

	
	g.P("var ", serviceDescVar, " = ", grpcPkg, ".ServiceDesc {")
	g.P("ServiceName: ", strconv.Quote(fullServName), ",")
	g.P("HandlerType: (*", serverType, ")(nil),")
	g.P("Methods: []", grpcPkg, ".MethodDesc{")
	for i, method := range service.Method {
		if method.GetServerStreaming() || method.GetClientStreaming() {
			continue
		}
		g.P("{")
		g.P("MethodName: ", strconv.Quote(method.GetName()), ",")
		g.P("Handler: ", handlerNames[i], ",")
		g.P("},")
	}
	g.P("},")
	g.P("Streams: []", grpcPkg, ".StreamDesc{")
	for i, method := range service.Method {
		if !method.GetServerStreaming() && !method.GetClientStreaming() {
			continue
		}
		g.P("{")
		g.P("StreamName: ", strconv.Quote(method.GetName()), ",")
		g.P("Handler: ", handlerNames[i], ",")
		if method.GetServerStreaming() {
			g.P("ServerStreams: true,")
		}
		if method.GetClientStreaming() {
			g.P("ClientStreams: true,")
		}
		g.P("},")
	}
	g.P("},")
	g.P("Metadata: \"", file.GetName(), "\",")
	g.P("}")
	g.P()
}


func (g *grpc) generateClientSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}
	reqArg := ", in *" + g.typeName(method.GetInputType())
	if method.GetClientStreaming() {
		reqArg = ""
	}
	respName := "*" + g.typeName(method.GetOutputType())
	if method.GetServerStreaming() || method.GetClientStreaming() {
		respName = servName + "_" + generator.CamelCase(origMethName) + "Client"
	}
	return fmt.Sprintf("%s(ctx %s.Context%s, opts ...%s.CallOption) (%s, error)", methName, contextPkg, reqArg, grpcPkg, respName)
}

func (g *grpc) generateClientMethod(servName, fullServName, serviceDescVar string, method *pb.MethodDescriptorProto, descExpr string) {
	sname := fmt.Sprintf("/%s/%s", fullServName, method.GetName())
	methName := generator.CamelCase(method.GetName())
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())

	if method.GetOptions().GetDeprecated() {
		g.P(deprecationComment)
	}
	g.P("func (c *", unexport(servName), "Client) ", g.generateClientSignature(servName, method), "{")
	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		g.P("out := new(", outType, ")")
		
		g.P(`err := c.cc.Invoke(ctx, "`, sname, `", in, out, opts...)`)
		g.P("if err != nil { return nil, err }")
		g.P("return out, nil")
		g.P("}")
		g.P()
		return
	}
	streamType := unexport(servName) + methName + "Client"
	g.P("stream, err := c.cc.NewStream(ctx, ", descExpr, `, "`, sname, `", opts...)`)
	g.P("if err != nil { return nil, err }")
	g.P("x := &", streamType, "{stream}")
	if !method.GetClientStreaming() {
		g.P("if err := x.ClientStream.SendMsg(in); err != nil { return nil, err }")
		g.P("if err := x.ClientStream.CloseSend(); err != nil { return nil, err }")
	}
	g.P("return x, nil")
	g.P("}")
	g.P()

	genSend := method.GetClientStreaming()
	genRecv := method.GetServerStreaming()
	genCloseAndRecv := !method.GetServerStreaming()

	
	g.P("type ", servName, "_", methName, "Client interface {")
	if genSend {
		g.P("Send(*", inType, ") error")
	}
	if genRecv {
		g.P("Recv() (*", outType, ", error)")
	}
	if genCloseAndRecv {
		g.P("CloseAndRecv() (*", outType, ", error)")
	}
	g.P(grpcPkg, ".ClientStream")
	g.P("}")
	g.P()

	g.P("type ", streamType, " struct {")
	g.P(grpcPkg, ".ClientStream")
	g.P("}")
	g.P()

	if genSend {
		g.P("func (x *", streamType, ") Send(m *", inType, ") error {")
		g.P("return x.ClientStream.SendMsg(m)")
		g.P("}")
		g.P()
	}
	if genRecv {
		g.P("func (x *", streamType, ") Recv() (*", outType, ", error) {")
		g.P("m := new(", outType, ")")
		g.P("if err := x.ClientStream.RecvMsg(m); err != nil { return nil, err }")
		g.P("return m, nil")
		g.P("}")
		g.P()
	}
	if genCloseAndRecv {
		g.P("func (x *", streamType, ") CloseAndRecv() (*", outType, ", error) {")
		g.P("if err := x.ClientStream.CloseSend(); err != nil { return nil, err }")
		g.P("m := new(", outType, ")")
		g.P("if err := x.ClientStream.RecvMsg(m); err != nil { return nil, err }")
		g.P("return m, nil")
		g.P("}")
		g.P()
	}
}


func (g *grpc) generateServerSignature(servName string, method *pb.MethodDescriptorProto) string {
	origMethName := method.GetName()
	methName := generator.CamelCase(origMethName)
	if reservedClientName[methName] {
		methName += "_"
	}

	var reqArgs []string
	ret := "error"
	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		reqArgs = append(reqArgs, contextPkg+".Context")
		ret = "(*" + g.typeName(method.GetOutputType()) + ", error)"
	}
	if !method.GetClientStreaming() {
		reqArgs = append(reqArgs, "*"+g.typeName(method.GetInputType()))
	}
	if method.GetServerStreaming() || method.GetClientStreaming() {
		reqArgs = append(reqArgs, servName+"_"+generator.CamelCase(origMethName)+"Server")
	}

	return methName + "(" + strings.Join(reqArgs, ", ") + ") " + ret
}

func (g *grpc) generateServerMethod(servName, fullServName string, method *pb.MethodDescriptorProto) string {
	methName := generator.CamelCase(method.GetName())
	hname := fmt.Sprintf("_%s_%s_Handler", servName, methName)
	inType := g.typeName(method.GetInputType())
	outType := g.typeName(method.GetOutputType())

	if !method.GetServerStreaming() && !method.GetClientStreaming() {
		g.P("func ", hname, "(srv interface{}, ctx ", contextPkg, ".Context, dec func(interface{}) error, interceptor ", grpcPkg, ".UnaryServerInterceptor) (interface{}, error) {")
		g.P("in := new(", inType, ")")
		g.P("if err := dec(in); err != nil { return nil, err }")
		g.P("if interceptor == nil { return srv.(", servName, "Server).", methName, "(ctx, in) }")
		g.P("info := &", grpcPkg, ".UnaryServerInfo{")
		g.P("Server: srv,")
		g.P("FullMethod: ", strconv.Quote(fmt.Sprintf("/%s/%s", fullServName, methName)), ",")
		g.P("}")
		g.P("handler := func(ctx ", contextPkg, ".Context, req interface{}) (interface{}, error) {")
		g.P("return srv.(", servName, "Server).", methName, "(ctx, req.(*", inType, "))")
		g.P("}")
		g.P("return interceptor(ctx, in, info, handler)")
		g.P("}")
		g.P()
		return hname
	}
	streamType := unexport(servName) + methName + "Server"
	g.P("func ", hname, "(srv interface{}, stream ", grpcPkg, ".ServerStream) error {")
	if !method.GetClientStreaming() {
		g.P("m := new(", inType, ")")
		g.P("if err := stream.RecvMsg(m); err != nil { return err }")
		g.P("return srv.(", servName, "Server).", methName, "(m, &", streamType, "{stream})")
	} else {
		g.P("return srv.(", servName, "Server).", methName, "(&", streamType, "{stream})")
	}
	g.P("}")
	g.P()

	genSend := method.GetServerStreaming()
	genSendAndClose := !method.GetServerStreaming()
	genRecv := method.GetClientStreaming()

	
	g.P("type ", servName, "_", methName, "Server interface {")
	if genSend {
		g.P("Send(*", outType, ") error")
	}
	if genSendAndClose {
		g.P("SendAndClose(*", outType, ") error")
	}
	if genRecv {
		g.P("Recv() (*", inType, ", error)")
	}
	g.P(grpcPkg, ".ServerStream")
	g.P("}")
	g.P()

	g.P("type ", streamType, " struct {")
	g.P(grpcPkg, ".ServerStream")
	g.P("}")
	g.P()

	if genSend {
		g.P("func (x *", streamType, ") Send(m *", outType, ") error {")
		g.P("return x.ServerStream.SendMsg(m)")
		g.P("}")
		g.P()
	}
	if genSendAndClose {
		g.P("func (x *", streamType, ") SendAndClose(m *", outType, ") error {")
		g.P("return x.ServerStream.SendMsg(m)")
		g.P("}")
		g.P()
	}
	if genRecv {
		g.P("func (x *", streamType, ") Recv() (*", inType, ", error) {")
		g.P("m := new(", inType, ")")
		g.P("if err := x.ServerStream.RecvMsg(m); err != nil { return nil, err }")
		g.P("return m, nil")
		g.P("}")
		g.P()
	}

	return hname
}
