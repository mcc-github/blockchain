package archive 

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/docker/pkg/fileutils"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)

var unpigzPath string

func init() {
	if path, err := exec.LookPath("unpigz"); err != nil {
		logrus.Debug("unpigz binary not found in PATH, falling back to go gzip library")
	} else {
		logrus.Debugf("Using unpigz binary found at path %s", path)
		unpigzPath = path
	}
}

type (
	
	Compression int
	
	WhiteoutFormat int

	
	TarOptions struct {
		IncludeFiles     []string
		ExcludePatterns  []string
		Compression      Compression
		NoLchown         bool
		UIDMaps          []idtools.IDMap
		GIDMaps          []idtools.IDMap
		ChownOpts        *idtools.IDPair
		IncludeSourceDir bool
		
		
		
		WhiteoutFormat WhiteoutFormat
		
		
		NoOverwriteDirNonDir bool
		
		
		RebaseNames map[string]string
		InUserNS    bool
	}
)




type Archiver struct {
	Untar         func(io.Reader, string, *TarOptions) error
	IDMappingsVar *idtools.IDMappings
}


func NewDefaultArchiver() *Archiver {
	return &Archiver{Untar: Untar, IDMappingsVar: &idtools.IDMappings{}}
}




type breakoutError error

const (
	
	Uncompressed Compression = iota
	
	Bzip2
	
	Gzip
	
	Xz
)

const (
	
	AUFSWhiteoutFormat WhiteoutFormat = iota
	
	
	OverlayWhiteoutFormat
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



func IsArchivePath(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	rdr, err := DecompressStream(file)
	if err != nil {
		return false
	}
	r := tar.NewReader(rdr)
	_, err = r.Next()
	return err == nil
}


func DetectCompression(source []byte) Compression {
	for compression, m := range map[Compression][]byte{
		Bzip2: {0x42, 0x5A, 0x68},
		Gzip:  {0x1F, 0x8B, 0x08},
		Xz:    {0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00},
	} {
		if len(source) < len(m) {
			logrus.Debug("Len too short")
			continue
		}
		if bytes.Equal(m, source[:len(m)]) {
			return compression
		}
	}
	return Uncompressed
}

func xzDecompress(ctx context.Context, archive io.Reader) (io.ReadCloser, error) {
	args := []string{"xz", "-d", "-c", "-q"}

	return cmdStream(exec.CommandContext(ctx, args[0], args[1:]...), archive)
}

func gzDecompress(ctx context.Context, buf io.Reader) (io.ReadCloser, error) {
	if unpigzPath == "" {
		return gzip.NewReader(buf)
	}

	disablePigzEnv := os.Getenv("MOBY_DISABLE_PIGZ")
	if disablePigzEnv != "" {
		if disablePigz, err := strconv.ParseBool(disablePigzEnv); err != nil {
			return nil, err
		} else if disablePigz {
			return gzip.NewReader(buf)
		}
	}

	return cmdStream(exec.CommandContext(ctx, unpigzPath, "-d", "-c"), buf)
}

func wrapReadCloser(readBuf io.ReadCloser, cancel context.CancelFunc) io.ReadCloser {
	return ioutils.NewReadCloserWrapper(readBuf, func() error {
		cancel()
		return readBuf.Close()
	})
}


func DecompressStream(archive io.Reader) (io.ReadCloser, error) {
	p := pools.BufioReader32KPool
	buf := p.Get(archive)
	bs, err := buf.Peek(10)
	if err != nil && err != io.EOF {
		
		
		
		
		
		
		return nil, err
	}

	compression := DetectCompression(bs)
	switch compression {
	case Uncompressed:
		readBufWrapper := p.NewReadCloserWrapper(buf, buf)
		return readBufWrapper, nil
	case Gzip:
		ctx, cancel := context.WithCancel(context.Background())

		gzReader, err := gzDecompress(ctx, buf)
		if err != nil {
			cancel()
			return nil, err
		}
		readBufWrapper := p.NewReadCloserWrapper(buf, gzReader)
		return wrapReadCloser(readBufWrapper, cancel), nil
	case Bzip2:
		bz2Reader := bzip2.NewReader(buf)
		readBufWrapper := p.NewReadCloserWrapper(buf, bz2Reader)
		return readBufWrapper, nil
	case Xz:
		ctx, cancel := context.WithCancel(context.Background())

		xzReader, err := xzDecompress(ctx, buf)
		if err != nil {
			cancel()
			return nil, err
		}
		readBufWrapper := p.NewReadCloserWrapper(buf, xzReader)
		return wrapReadCloser(readBufWrapper, cancel), nil
	default:
		return nil, fmt.Errorf("Unsupported compression format %s", (&compression).Extension())
	}
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






type TarModifierFunc func(path string, header *tar.Header, content io.Reader) (*tar.Header, []byte, error)



func ReplaceFileTarWrapper(inputTarStream io.ReadCloser, mods map[string]TarModifierFunc) io.ReadCloser {
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		tarReader := tar.NewReader(inputTarStream)
		tarWriter := tar.NewWriter(pipeWriter)
		defer inputTarStream.Close()
		defer tarWriter.Close()

		modify := func(name string, original *tar.Header, modifier TarModifierFunc, tarReader io.Reader) error {
			header, data, err := modifier(name, original, tarReader)
			switch {
			case err != nil:
				return err
			case header == nil:
				return nil
			}

			header.Name = name
			header.Size = int64(len(data))
			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}
			if len(data) != 0 {
				if _, err := tarWriter.Write(data); err != nil {
					return err
				}
			}
			return nil
		}

		var err error
		var originalHeader *tar.Header
		for {
			originalHeader, err = tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				pipeWriter.CloseWithError(err)
				return
			}

			modifier, ok := mods[originalHeader.Name]
			if !ok {
				
				if err := tarWriter.WriteHeader(originalHeader); err != nil {
					pipeWriter.CloseWithError(err)
					return
				}
				if _, err := pools.Copy(tarWriter, tarReader); err != nil {
					pipeWriter.CloseWithError(err)
					return
				}
				continue
			}
			delete(mods, originalHeader.Name)

			if err := modify(originalHeader.Name, originalHeader, modifier, tarReader); err != nil {
				pipeWriter.CloseWithError(err)
				return
			}
		}

		
		for name, modifier := range mods {
			if err := modify(name, nil, modifier, nil); err != nil {
				pipeWriter.CloseWithError(err)
				return
			}
		}

		pipeWriter.Close()

	}()
	return pipeReader
}


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





func FileInfoHeader(name string, fi os.FileInfo, link string) (*tar.Header, error) {
	hdr, err := tar.FileInfoHeader(fi, link)
	if err != nil {
		return nil, err
	}
	hdr.Format = tar.FormatPAX
	hdr.ModTime = hdr.ModTime.Truncate(time.Second)
	hdr.AccessTime = time.Time{}
	hdr.ChangeTime = time.Time{}
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



func ReadSecurityXattrToTarHeader(path string, hdr *tar.Header) error {
	capability, _ := system.Lgetxattr(path, "security.capability")
	if capability != nil {
		hdr.Xattrs = make(map[string]string)
		hdr.Xattrs["security.capability"] = string(capability)
	}
	return nil
}

type tarWhiteoutConverter interface {
	ConvertWrite(*tar.Header, string, os.FileInfo) (*tar.Header, error)
	ConvertRead(*tar.Header, string) (bool, error)
}

type tarAppender struct {
	TarWriter *tar.Writer
	Buffer    *bufio.Writer

	
	SeenFiles  map[uint64]string
	IDMappings *idtools.IDMappings
	ChownOpts  *idtools.IDPair

	
	
	
	
	WhiteoutConverter tarWhiteoutConverter
}

func newTarAppender(idMapping *idtools.IDMappings, writer io.Writer, chownOpts *idtools.IDPair) *tarAppender {
	return &tarAppender{
		SeenFiles:  make(map[uint64]string),
		TarWriter:  tar.NewWriter(writer),
		Buffer:     pools.BufioWriter32KPool.Get(nil),
		IDMappings: idMapping,
		ChownOpts:  chownOpts,
	}
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
		!ta.IDMappings.Empty() {
		fileIDPair, err := getFileUIDGID(fi.Sys())
		if err != nil {
			return err
		}
		hdr.Uid, hdr.Gid, err = ta.IDMappings.ToContainer(fileIDPair)
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

func createTarFile(path, extractDir string, hdr *tar.Header, reader io.Reader, Lchown bool, chownOpts *idtools.IDPair, inUserns bool) error {
	
	
	
	hdrInfo := hdr.FileInfo()

	switch hdr.Typeflag {
	case tar.TypeDir:
		
		
		if fi, err := os.Lstat(path); !(err == nil && fi.IsDir()) {
			if err := os.Mkdir(path, hdrInfo.Mode()); err != nil {
				return err
			}
		}

	case tar.TypeReg, tar.TypeRegA:
		
		
		
		file, err := system.OpenFileSequential(path, os.O_CREATE|os.O_WRONLY, hdrInfo.Mode())
		if err != nil {
			return err
		}
		if _, err := io.Copy(file, reader); err != nil {
			file.Close()
			return err
		}
		file.Close()

	case tar.TypeBlock, tar.TypeChar:
		if inUserns { 
			return nil
		}
		
		if err := handleTarTypeBlockCharFifo(hdr, path); err != nil {
			return err
		}

	case tar.TypeFifo:
		
		if err := handleTarTypeBlockCharFifo(hdr, path); err != nil {
			return err
		}

	case tar.TypeLink:
		targetPath := filepath.Join(extractDir, hdr.Linkname)
		
		if !strings.HasPrefix(targetPath, extractDir) {
			return breakoutError(fmt.Errorf("invalid hardlink %q -> %q", targetPath, hdr.Linkname))
		}
		if err := os.Link(targetPath, path); err != nil {
			return err
		}

	case tar.TypeSymlink:
		
		
		targetPath := filepath.Join(filepath.Dir(path), hdr.Linkname)

		
		
		if !strings.HasPrefix(targetPath, extractDir) {
			return breakoutError(fmt.Errorf("invalid symlink %q -> %q", path, hdr.Linkname))
		}
		if err := os.Symlink(hdr.Linkname, path); err != nil {
			return err
		}

	case tar.TypeXGlobalHeader:
		logrus.Debug("PAX Global Extended Headers found and ignored")
		return nil

	default:
		return fmt.Errorf("unhandled tar header type %d", hdr.Typeflag)
	}

	
	if Lchown && runtime.GOOS != "windows" {
		if chownOpts == nil {
			chownOpts = &idtools.IDPair{UID: hdr.Uid, GID: hdr.Gid}
		}
		if err := os.Lchown(path, chownOpts.UID, chownOpts.GID); err != nil {
			return err
		}
	}

	var errors []string
	for key, value := range hdr.Xattrs {
		if err := system.Lsetxattr(path, key, []byte(value), 0); err != nil {
			if err == syscall.ENOTSUP {
				
				
				
				
				errors = append(errors, err.Error())
				continue
			}
			return err
		}

	}

	if len(errors) > 0 {
		logrus.WithFields(logrus.Fields{
			"errors": errors,
		}).Warn("ignored xattrs in archive: underlying filesystem doesn't support them")
	}

	
	
	if err := handleLChmod(hdr, path, hdrInfo); err != nil {
		return err
	}

	aTime := hdr.AccessTime
	if aTime.Before(hdr.ModTime) {
		
		aTime = hdr.ModTime
	}

	
	if hdr.Typeflag == tar.TypeLink {
		if fi, err := os.Lstat(hdr.Linkname); err == nil && (fi.Mode()&os.ModeSymlink == 0) {
			if err := system.Chtimes(path, aTime, hdr.ModTime); err != nil {
				return err
			}
		}
	} else if hdr.Typeflag != tar.TypeSymlink {
		if err := system.Chtimes(path, aTime, hdr.ModTime); err != nil {
			return err
		}
	} else {
		ts := []syscall.Timespec{timeToTimespec(aTime), timeToTimespec(hdr.ModTime)}
		if err := system.LUtimesNano(path, ts); err != nil && err != system.ErrNotSupportedPlatform {
			return err
		}
	}
	return nil
}



func Tar(path string, compression Compression) (io.ReadCloser, error) {
	return TarWithOptions(path, &TarOptions{Compression: compression})
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


func Unpack(decompressedArchive io.Reader, dest string, options *TarOptions) error {
	tr := tar.NewReader(decompressedArchive)
	trBuf := pools.BufioReader32KPool.Get(nil)
	defer pools.BufioReader32KPool.Put(trBuf)

	var dirs []*tar.Header
	idMappings := idtools.NewIDMappingsFromMaps(options.UIDMaps, options.GIDMaps)
	rootIDs := idMappings.RootPair()
	whiteoutConverter := getWhiteoutConverter(options.WhiteoutFormat)

	
loop:
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			
			break
		}
		if err != nil {
			return err
		}

		
		
		
		hdr.Name = filepath.Clean(hdr.Name)

		for _, exclude := range options.ExcludePatterns {
			if strings.HasPrefix(hdr.Name, exclude) {
				continue loop
			}
		}

		
		
		
		if !strings.HasSuffix(hdr.Name, string(os.PathSeparator)) {
			
			parent := filepath.Dir(hdr.Name)
			parentPath := filepath.Join(dest, parent)
			if _, err := os.Lstat(parentPath); err != nil && os.IsNotExist(err) {
				err = idtools.MkdirAllAndChownNew(parentPath, 0777, rootIDs)
				if err != nil {
					return err
				}
			}
		}

		path := filepath.Join(dest, hdr.Name)
		rel, err := filepath.Rel(dest, path)
		if err != nil {
			return err
		}
		if strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return breakoutError(fmt.Errorf("%q is outside of %q", hdr.Name, dest))
		}

		
		
		
		
		if fi, err := os.Lstat(path); err == nil {
			if options.NoOverwriteDirNonDir && fi.IsDir() && hdr.Typeflag != tar.TypeDir {
				
				
				return fmt.Errorf("cannot overwrite directory %q with non-directory %q", path, dest)
			}

			if options.NoOverwriteDirNonDir && !fi.IsDir() && hdr.Typeflag == tar.TypeDir {
				
				
				return fmt.Errorf("cannot overwrite non-directory %q with directory %q", path, dest)
			}

			if fi.IsDir() && hdr.Name == "." {
				continue
			}

			if !(fi.IsDir() && hdr.Typeflag == tar.TypeDir) {
				if err := os.RemoveAll(path); err != nil {
					return err
				}
			}
		}
		trBuf.Reset(tr)

		if err := remapIDs(idMappings, hdr); err != nil {
			return err
		}

		if whiteoutConverter != nil {
			writeFile, err := whiteoutConverter.ConvertRead(hdr, path)
			if err != nil {
				return err
			}
			if !writeFile {
				continue
			}
		}

		if err := createTarFile(path, dest, hdr, trBuf, !options.NoLchown, options.ChownOpts, options.InUserNS); err != nil {
			return err
		}

		
		
		if hdr.Typeflag == tar.TypeDir {
			dirs = append(dirs, hdr)
		}
	}

	for _, hdr := range dirs {
		path := filepath.Join(dest, hdr.Name)

		if err := system.Chtimes(path, hdr.AccessTime, hdr.ModTime); err != nil {
			return err
		}
	}
	return nil
}






func Untar(tarArchive io.Reader, dest string, options *TarOptions) error {
	return untarHandler(tarArchive, dest, options, true)
}




func UntarUncompressed(tarArchive io.Reader, dest string, options *TarOptions) error {
	return untarHandler(tarArchive, dest, options, false)
}


func untarHandler(tarArchive io.Reader, dest string, options *TarOptions, decompress bool) error {
	if tarArchive == nil {
		return fmt.Errorf("Empty archive")
	}
	dest = filepath.Clean(dest)
	if options == nil {
		options = &TarOptions{}
	}
	if options.ExcludePatterns == nil {
		options.ExcludePatterns = []string{}
	}

	r := tarArchive
	if decompress {
		decompressedArchive, err := DecompressStream(tarArchive)
		if err != nil {
			return err
		}
		defer decompressedArchive.Close()
		r = decompressedArchive
	}

	return Unpack(r, dest, options)
}



func (archiver *Archiver) TarUntar(src, dst string) error {
	logrus.Debugf("TarUntar(%s %s)", src, dst)
	archive, err := TarWithOptions(src, &TarOptions{Compression: Uncompressed})
	if err != nil {
		return err
	}
	defer archive.Close()
	options := &TarOptions{
		UIDMaps: archiver.IDMappingsVar.UIDs(),
		GIDMaps: archiver.IDMappingsVar.GIDs(),
	}
	return archiver.Untar(archive, dst, options)
}


func (archiver *Archiver) UntarPath(src, dst string) error {
	archive, err := os.Open(src)
	if err != nil {
		return err
	}
	defer archive.Close()
	options := &TarOptions{
		UIDMaps: archiver.IDMappingsVar.UIDs(),
		GIDMaps: archiver.IDMappingsVar.GIDs(),
	}
	return archiver.Untar(archive, dst, options)
}





func (archiver *Archiver) CopyWithTar(src, dst string) error {
	srcSt, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcSt.IsDir() {
		return archiver.CopyFileWithTar(src, dst)
	}

	
	
	
	rootIDs := archiver.IDMappingsVar.RootPair()
	
	logrus.Debugf("Creating dest directory: %s", dst)
	if err := idtools.MkdirAllAndChownNew(dst, 0755, rootIDs); err != nil {
		return err
	}
	logrus.Debugf("Calling TarUntar(%s, %s)", src, dst)
	return archiver.TarUntar(src, dst)
}




func (archiver *Archiver) CopyFileWithTar(src, dst string) (err error) {
	logrus.Debugf("CopyFileWithTar(%s, %s)", src, dst)
	srcSt, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcSt.IsDir() {
		return fmt.Errorf("Can't copy a directory")
	}

	
	
	if dst[len(dst)-1] == os.PathSeparator {
		dst = filepath.Join(dst, filepath.Base(src))
	}
	
	if err := system.MkdirAll(filepath.Dir(dst), 0700, ""); err != nil {
		return err
	}

	r, w := io.Pipe()
	errC := make(chan error, 1)

	go func() {
		defer close(errC)

		errC <- func() error {
			defer w.Close()

			srcF, err := os.Open(src)
			if err != nil {
				return err
			}
			defer srcF.Close()

			hdr, err := tar.FileInfoHeader(srcSt, "")
			if err != nil {
				return err
			}
			hdr.Format = tar.FormatPAX
			hdr.ModTime = hdr.ModTime.Truncate(time.Second)
			hdr.AccessTime = time.Time{}
			hdr.ChangeTime = time.Time{}
			hdr.Name = filepath.Base(dst)
			hdr.Mode = int64(chmodTarEntry(os.FileMode(hdr.Mode)))

			if err := remapIDs(archiver.IDMappingsVar, hdr); err != nil {
				return err
			}

			tw := tar.NewWriter(w)
			defer tw.Close()
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
			if _, err := io.Copy(tw, srcF); err != nil {
				return err
			}
			return nil
		}()
	}()
	defer func() {
		if er := <-errC; err == nil && er != nil {
			err = er
		}
	}()

	err = archiver.Untar(r, filepath.Dir(dst), nil)
	if err != nil {
		r.CloseWithError(err)
	}
	return err
}


func (archiver *Archiver) IDMappings() *idtools.IDMappings {
	return archiver.IDMappingsVar
}

func remapIDs(idMappings *idtools.IDMappings, hdr *tar.Header) error {
	ids, err := idMappings.ToHost(idtools.IDPair{UID: hdr.Uid, GID: hdr.Gid})
	hdr.Uid, hdr.Gid = ids.UID, ids.GID
	return err
}




func cmdStream(cmd *exec.Cmd, input io.Reader) (io.ReadCloser, error) {
	cmd.Stdin = input
	pipeR, pipeW := io.Pipe()
	cmd.Stdout = pipeW
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	
	go func() {
		if err := cmd.Wait(); err != nil {
			pipeW.CloseWithError(fmt.Errorf("%s: %s", err, errBuf.String()))
		} else {
			pipeW.Close()
		}
	}()

	return pipeR, nil
}




func NewTempArchive(src io.Reader, dir string) (*TempArchive, error) {
	f, err := ioutil.TempFile(dir, "")
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(f, src); err != nil {
		return nil, err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := st.Size()
	return &TempArchive{File: f, Size: size}, nil
}



type TempArchive struct {
	*os.File
	Size   int64 
	read   int64
	closed bool
}



func (archive *TempArchive) Close() error {
	if archive.closed {
		return nil
	}

	archive.closed = true

	return archive.File.Close()
}

func (archive *TempArchive) Read(data []byte) (int, error) {
	n, err := archive.File.Read(data)
	archive.read += int64(n)
	if err != nil || archive.read == archive.Size {
		archive.Close()
		os.Remove(archive.File.Name())
	}
	return n, err
}
