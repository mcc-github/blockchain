



package archive

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)

const (
	
	Uncompressed Compression = iota
	
	Bzip2
	
	Gzip
	
	Xz
)

const (
	modeISDIR  = 040000  
	modeISFIFO = 010000  
	modeISREG  = 0100000 
	modeISLNK  = 0120000 
	modeISBLK  = 060000  
	modeISCHR  = 020000  
	modeISSOCK = 0140000 
)


type Compression int


func (compression *Compression) Extension() string {
	switch *compression {
	case Uncompressed:
		return "tar"
	case Bzip2:
		return "tar.bz2"
	case Gzip:
		return "tar.gz"
	case Xz:
		return "tar.xz"
	}
	return ""
}


type WhiteoutFormat int


type TarOptions struct {
	IncludeFiles     []string
	ExcludePatterns  []string
	Compression      Compression
	NoLchown         bool
	UIDMaps          []idtools.IDMap
	GIDMaps          []idtools.IDMap
	ChownOpts        *idtools.Identity
	IncludeSourceDir bool
	
	
	
	WhiteoutFormat WhiteoutFormat
	
	
	NoOverwriteDirNonDir bool
	
	
	RebaseNames map[string]string
	InUserNS    bool
}



func TarWithOptions(srcPath string, options *TarOptions) (io.ReadCloser, error) {

	
	
	srcPath = fixVolumePathPrefix(srcPath)

	pm, err := fileutils.NewPatternMatcher(options.ExcludePatterns)
	if err != nil {
		return nil, err
	}

	pipeReader, pipeWriter := io.Pipe()

	compressWriter, err := CompressStream(pipeWriter, options.Compression)
	if err != nil {
		return nil, err
	}

	go func() {
		ta := newTarAppender(
			idtools.NewIDMappingsFromMaps(options.UIDMaps, options.GIDMaps),
			compressWriter,
			options.ChownOpts,
		)
		ta.WhiteoutConverter = getWhiteoutConverter(options.WhiteoutFormat)

		defer func() {
			
			if err := ta.TarWriter.Close(); err != nil {
				logrus.Errorf("Can't close tar writer: %s", err)
			}
			if err := compressWriter.Close(); err != nil {
				logrus.Errorf("Can't close compress writer: %s", err)
			}
			if err := pipeWriter.Close(); err != nil {
				logrus.Errorf("Can't close pipe writer: %s", err)
			}
		}()

		
		defer pools.BufioWriter32KPool.Put(ta.Buffer)

		
		
		
		

		stat, err := os.Lstat(srcPath)
		if err != nil {
			return
		}

		if !stat.IsDir() {
			
			
			
			
			if len(options.IncludeFiles) > 0 {
				logrus.Warn("Tar: Can't archive a file with includes")
			}

			dir, base := SplitPathDirEntry(srcPath)
			srcPath = dir
			options.IncludeFiles = []string{base}
		}

		if len(options.IncludeFiles) == 0 {
			options.IncludeFiles = []string{"."}
		}

		seen := make(map[string]bool)

		for _, include := range options.IncludeFiles {
			rebaseName := options.RebaseNames[include]

			walkRoot := getWalkRoot(srcPath, include)
			filepath.Walk(walkRoot, func(filePath string, f os.FileInfo, err error) error {
				if err != nil {
					logrus.Errorf("Tar: Can't stat file %s to tar: %s", srcPath, err)
					return nil
				}

				relFilePath, err := filepath.Rel(srcPath, filePath)
				if err != nil || (!options.IncludeSourceDir && relFilePath == "." && f.IsDir()) {
					
					
					return nil
				}

				if options.IncludeSourceDir && include == "." && relFilePath != "." {
					relFilePath = strings.Join([]string{".", relFilePath}, string(filepath.Separator))
				}

				skip := false

				
				
				
				
				
				if include != relFilePath {
					skip, err = pm.Matches(relFilePath)
					if err != nil {
						logrus.Errorf("Error matching %s: %v", relFilePath, err)
						return err
					}
				}

				if skip {
					
					
					
					

					
					if !f.IsDir() {
						return nil
					}

					
					if !pm.Exclusions() {
						return filepath.SkipDir
					}

					dirSlash := relFilePath + string(filepath.Separator)

					for _, pat := range pm.Patterns() {
						if !pat.Exclusion() {
							continue
						}
						if strings.HasPrefix(pat.String()+string(filepath.Separator), dirSlash) {
							
							return nil
						}
					}

					
					return filepath.SkipDir
				}

				if seen[relFilePath] {
					return nil
				}
				seen[relFilePath] = true

				
				if rebaseName != "" {
					var replacement string
					if rebaseName != string(filepath.Separator) {
						
						
						
						replacement = rebaseName
					}

					relFilePath = strings.Replace(relFilePath, include, replacement, 1)
				}

				if err := ta.addTarFile(filePath, relFilePath); err != nil {
					logrus.Errorf("Can't add file %s to tar: %s", filePath, err)
					
					if err == io.ErrClosedPipe {
						return err
					}
				}
				return nil
			})
		}
	}()

	return pipeReader, nil
}


func CompressStream(dest io.Writer, compression Compression) (io.WriteCloser, error) {
	p := pools.BufioWriter32KPool
	buf := p.Get(dest)
	switch compression {
	case Uncompressed:
		writeBufWrapper := p.NewWriteCloserWrapper(buf, buf)
		return writeBufWrapper, nil
	case Gzip:
		gzWriter := gzip.NewWriter(dest)
		writeBufWrapper := p.NewWriteCloserWrapper(buf, gzWriter)
		return writeBufWrapper, nil
	case Bzip2, Xz:
		
		
		return nil, fmt.Errorf("Unsupported compression format %s", (&compression).Extension())
	default:
		return nil, fmt.Errorf("Unsupported compression format %s", (&compression).Extension())
	}
}

type tarWhiteoutConverter interface {
	ConvertWrite(*tar.Header, string, os.FileInfo) (*tar.Header, error)
	ConvertRead(*tar.Header, string) (bool, error)
}

type tarAppender struct {
	TarWriter *tar.Writer
	Buffer    *bufio.Writer

	
	SeenFiles       map[uint64]string
	IdentityMapping *idtools.IdentityMapping
	ChownOpts       *idtools.Identity

	
	
	
	
	WhiteoutConverter tarWhiteoutConverter
}

func newTarAppender(idMapping *idtools.IdentityMapping, writer io.Writer, chownOpts *idtools.Identity) *tarAppender {
	return &tarAppender{
		SeenFiles:       make(map[uint64]string),
		TarWriter:       tar.NewWriter(writer),
		Buffer:          pools.BufioWriter32KPool.Get(nil),
		IdentityMapping: idMapping,
		ChownOpts:       chownOpts,
	}
}


func (ta *tarAppender) addTarFile(path, name string) error {
	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}

	var link string
	if fi.Mode()&os.ModeSymlink != 0 {
		var err error
		link, err = os.Readlink(path)
		if err != nil {
			return err
		}
	}

	hdr, err := FileInfoHeader(name, fi, link)
	if err != nil {
		return err
	}
	if err := ReadSecurityXattrToTarHeader(path, hdr); err != nil {
		return err
	}

	
	
	if !fi.IsDir() && hasHardlinks(fi) {
		inode, err := getInodeFromStat(fi.Sys())
		if err != nil {
			return err
		}
		
		
		if oldpath, ok := ta.SeenFiles[inode]; ok {
			hdr.Typeflag = tar.TypeLink
			hdr.Linkname = oldpath
			hdr.Size = 0 
		} else {
			ta.SeenFiles[inode] = name
		}
	}

	
	
	isOverlayWhiteout := fi.Mode()&os.ModeCharDevice != 0 && hdr.Devmajor == 0 && hdr.Devminor == 0

	
	
	
	if !isOverlayWhiteout &&
		!strings.HasPrefix(filepath.Base(hdr.Name), WhiteoutPrefix) &&
		!ta.IdentityMapping.Empty() {
		fileIdentity, err := getFileIdentity(fi.Sys())
		if err != nil {
			return err
		}
		hdr.Uid, hdr.Gid, err = ta.IdentityMapping.ToContainer(fileIdentity)
		if err != nil {
			return err
		}
	}

	
	if ta.ChownOpts != nil {
		hdr.Uid = ta.ChownOpts.UID
		hdr.Gid = ta.ChownOpts.GID
	}

	if ta.WhiteoutConverter != nil {
		wo, err := ta.WhiteoutConverter.ConvertWrite(hdr, path, fi)
		if err != nil {
			return err
		}

		
		
		
		
		
		if wo != nil {
			if err := ta.TarWriter.WriteHeader(hdr); err != nil {
				return err
			}
			if hdr.Typeflag == tar.TypeReg && hdr.Size > 0 {
				return fmt.Errorf("tar: cannot use whiteout for non-empty file")
			}
			hdr = wo
		}
	}

	if err := ta.TarWriter.WriteHeader(hdr); err != nil {
		return err
	}

	if hdr.Typeflag == tar.TypeReg && hdr.Size > 0 {
		
		
		
		file, err := system.OpenSequential(path)
		if err != nil {
			return err
		}

		ta.Buffer.Reset(ta.TarWriter)
		defer ta.Buffer.Reset(nil)
		_, err = io.Copy(ta.Buffer, file)
		file.Close()
		if err != nil {
			return err
		}
		err = ta.Buffer.Flush()
		if err != nil {
			return err
		}
	}

	return nil
}



func ReadSecurityXattrToTarHeader(path string, hdr *tar.Header) error {
	capability, _ := system.Lgetxattr(path, "security.capability")
	if capability != nil {
		hdr.Xattrs = make(map[string]string)
		hdr.Xattrs["security.capability"] = string(capability)
	}
	return nil
}





func FileInfoHeader(name string, fi os.FileInfo, link string) (*tar.Header, error) {
	hdr, err := tar.FileInfoHeader(fi, link)
	if err != nil {
		return nil, err
	}
	hdr.Mode = fillGo18FileTypeBits(int64(chmodTarEntry(os.FileMode(hdr.Mode))), fi)
	name, err = canonicalTarName(name, fi.IsDir())
	if err != nil {
		return nil, fmt.Errorf("tar: cannot canonicalize path: %v", err)
	}
	hdr.Name = name
	if err := setHeaderForSpecialDevice(hdr, name, fi.Sys()); err != nil {
		return nil, err
	}
	return hdr, nil
}



func fillGo18FileTypeBits(mode int64, fi os.FileInfo) int64 {
	fm := fi.Mode()
	switch {
	case fm.IsRegular():
		mode |= modeISREG
	case fi.IsDir():
		mode |= modeISDIR
	case fm&os.ModeSymlink != 0:
		mode |= modeISLNK
	case fm&os.ModeDevice != 0:
		if fm&os.ModeCharDevice != 0 {
			mode |= modeISCHR
		} else {
			mode |= modeISBLK
		}
	case fm&os.ModeNamedPipe != 0:
		mode |= modeISFIFO
	case fm&os.ModeSocket != 0:
		mode |= modeISSOCK
	}
	return mode
}



func canonicalTarName(name string, isDir bool) (string, error) {
	name, err := CanonicalTarNameForPath(name)
	if err != nil {
		return "", err
	}

	
	if isDir && !strings.HasSuffix(name, "/") {
		name += "/"
	}
	return name, nil
}
