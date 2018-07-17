package archive 

import (
	"archive/tar"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/system"
	"github.com/sirupsen/logrus"
)


var (
	ErrNotDirectory      = errors.New("not a directory")
	ErrDirNotExists      = errors.New("no such directory")
	ErrCannotCopyDir     = errors.New("cannot copy directory")
	ErrInvalidCopySource = errors.New("invalid copy source content")
)








func PreserveTrailingDotOrSeparator(cleanedPath string, originalPath string, sep byte) string {
	
	cleanedPath = strings.Replace(cleanedPath, "/", string(sep), -1)
	originalPath = strings.Replace(originalPath, "/", string(sep), -1)

	if !specifiesCurrentDir(cleanedPath) && specifiesCurrentDir(originalPath) {
		if !hasTrailingPathSeparator(cleanedPath, sep) {
			
			
			cleanedPath += string(sep)
		}
		cleanedPath += "."
	}

	if !hasTrailingPathSeparator(cleanedPath, sep) && hasTrailingPathSeparator(originalPath, sep) {
		cleanedPath += string(sep)
	}

	return cleanedPath
}




func assertsDirectory(path string, sep byte) bool {
	return hasTrailingPathSeparator(path, sep) || specifiesCurrentDir(path)
}



func hasTrailingPathSeparator(path string, sep byte) bool {
	return len(path) > 0 && path[len(path)-1] == sep
}



func specifiesCurrentDir(path string) bool {
	return filepath.Base(path) == "."
}




func SplitPathDirEntry(path string) (dir, base string) {
	cleanedPath := filepath.Clean(filepath.FromSlash(path))

	if specifiesCurrentDir(path) {
		cleanedPath += string(os.PathSeparator) + "."
	}

	return filepath.Dir(cleanedPath), filepath.Base(cleanedPath)
}








func TarResource(sourceInfo CopyInfo) (content io.ReadCloser, err error) {
	return TarResourceRebase(sourceInfo.Path, sourceInfo.RebaseName)
}



func TarResourceRebase(sourcePath, rebaseName string) (content io.ReadCloser, err error) {
	sourcePath = normalizePath(sourcePath)
	if _, err = os.Lstat(sourcePath); err != nil {
		
		
		
		return
	}

	
	
	sourceDir, sourceBase := SplitPathDirEntry(sourcePath)
	opts := TarResourceRebaseOpts(sourceBase, rebaseName)

	logrus.Debugf("copying %q from %q", sourceBase, sourceDir)
	return TarWithOptions(sourceDir, opts)
}



func TarResourceRebaseOpts(sourceBase string, rebaseName string) *TarOptions {
	filter := []string{sourceBase}
	return &TarOptions{
		Compression:      Uncompressed,
		IncludeFiles:     filter,
		IncludeSourceDir: true,
		RebaseNames: map[string]string{
			sourceBase: rebaseName,
		},
	}
}



type CopyInfo struct {
	Path       string
	Exists     bool
	IsDir      bool
	RebaseName string
}






func CopyInfoSourcePath(path string, followLink bool) (CopyInfo, error) {
	
	
	
	path = normalizePath(path)

	resolvedPath, rebaseName, err := ResolveHostSourcePath(path, followLink)
	if err != nil {
		return CopyInfo{}, err
	}

	stat, err := os.Lstat(resolvedPath)
	if err != nil {
		return CopyInfo{}, err
	}

	return CopyInfo{
		Path:       resolvedPath,
		Exists:     true,
		IsDir:      stat.IsDir(),
		RebaseName: rebaseName,
	}, nil
}




func CopyInfoDestinationPath(path string) (info CopyInfo, err error) {
	maxSymlinkIter := 10 
	path = normalizePath(path)
	originalPath := path

	stat, err := os.Lstat(path)

	if err == nil && stat.Mode()&os.ModeSymlink == 0 {
		
		return CopyInfo{
			Path:   path,
			Exists: true,
			IsDir:  stat.IsDir(),
		}, nil
	}

	
	for n := 0; err == nil && stat.Mode()&os.ModeSymlink != 0; n++ {
		if n > maxSymlinkIter {
			
			return CopyInfo{}, errors.New("too many symlinks in " + originalPath)
		}

		
		
		
		
		
		
		var linkTarget string

		linkTarget, err = os.Readlink(path)
		if err != nil {
			return CopyInfo{}, err
		}

		if !system.IsAbs(linkTarget) {
			
			dstParent, _ := SplitPathDirEntry(path)
			linkTarget = filepath.Join(dstParent, linkTarget)
		}

		path = linkTarget
		stat, err = os.Lstat(path)
	}

	if err != nil {
		
		
		if !os.IsNotExist(err) {
			return CopyInfo{}, err
		}

		
		dstParent, _ := SplitPathDirEntry(path)

		parentDirStat, err := os.Stat(dstParent)
		if err != nil {
			return CopyInfo{}, err
		}
		if !parentDirStat.IsDir() {
			return CopyInfo{}, ErrNotDirectory
		}

		return CopyInfo{Path: path}, nil
	}

	
	return CopyInfo{
		Path:   path,
		Exists: true,
		IsDir:  stat.IsDir(),
	}, nil
}





func PrepareArchiveCopy(srcContent io.Reader, srcInfo, dstInfo CopyInfo) (dstDir string, content io.ReadCloser, err error) {
	
	srcInfo.Path = normalizePath(srcInfo.Path)
	dstInfo.Path = normalizePath(dstInfo.Path)

	
	
	dstDir, dstBase := SplitPathDirEntry(dstInfo.Path)
	_, srcBase := SplitPathDirEntry(srcInfo.Path)

	switch {
	case dstInfo.Exists && dstInfo.IsDir:
		
		
		
		return dstInfo.Path, ioutil.NopCloser(srcContent), nil
	case dstInfo.Exists && srcInfo.IsDir:
		
		
		
		return "", nil, ErrCannotCopyDir
	case dstInfo.Exists:
		
		
		
		if len(srcInfo.RebaseName) != 0 {
			srcBase = srcInfo.RebaseName
		}
		return dstDir, RebaseArchiveEntries(srcContent, srcBase, dstBase), nil
	case srcInfo.IsDir:
		
		
		
		
		
		
		if len(srcInfo.RebaseName) != 0 {
			srcBase = srcInfo.RebaseName
		}
		return dstDir, RebaseArchiveEntries(srcContent, srcBase, dstBase), nil
	case assertsDirectory(dstInfo.Path, os.PathSeparator):
		
		
		
		
		return "", nil, ErrDirNotExists
	default:
		
		
		
		
		
		
		if len(srcInfo.RebaseName) != 0 {
			srcBase = srcInfo.RebaseName
		}
		return dstDir, RebaseArchiveEntries(srcContent, srcBase, dstBase), nil
	}

}



func RebaseArchiveEntries(srcContent io.Reader, oldBase, newBase string) io.ReadCloser {
	if oldBase == string(os.PathSeparator) {
		
		
		
		oldBase = ""
	}

	rebased, w := io.Pipe()

	go func() {
		srcTar := tar.NewReader(srcContent)
		rebasedTar := tar.NewWriter(w)

		for {
			hdr, err := srcTar.Next()
			if err == io.EOF {
				
				rebasedTar.Close()
				w.Close()
				return
			}
			if err != nil {
				w.CloseWithError(err)
				return
			}

			hdr.Name = strings.Replace(hdr.Name, oldBase, newBase, 1)
			if hdr.Typeflag == tar.TypeLink {
				hdr.Linkname = strings.Replace(hdr.Linkname, oldBase, newBase, 1)
			}

			if err = rebasedTar.WriteHeader(hdr); err != nil {
				w.CloseWithError(err)
				return
			}

			if _, err = io.Copy(rebasedTar, srcTar); err != nil {
				w.CloseWithError(err)
				return
			}
		}
	}()

	return rebased
}







func CopyResource(srcPath, dstPath string, followLink bool) error {
	var (
		srcInfo CopyInfo
		err     error
	)

	
	srcPath = normalizePath(srcPath)
	dstPath = normalizePath(dstPath)

	
	srcPath = PreserveTrailingDotOrSeparator(filepath.Clean(srcPath), srcPath, os.PathSeparator)
	dstPath = PreserveTrailingDotOrSeparator(filepath.Clean(dstPath), dstPath, os.PathSeparator)

	if srcInfo, err = CopyInfoSourcePath(srcPath, followLink); err != nil {
		return err
	}

	content, err := TarResource(srcInfo)
	if err != nil {
		return err
	}
	defer content.Close()

	return CopyTo(content, srcInfo, dstPath)
}



func CopyTo(content io.Reader, srcInfo CopyInfo, dstPath string) error {
	
	
	dstInfo, err := CopyInfoDestinationPath(normalizePath(dstPath))
	if err != nil {
		return err
	}

	dstDir, copyArchive, err := PrepareArchiveCopy(content, srcInfo, dstInfo)
	if err != nil {
		return err
	}
	defer copyArchive.Close()

	options := &TarOptions{
		NoLchown:             true,
		NoOverwriteDirNonDir: true,
	}

	return Untar(copyArchive, dstDir, options)
}





func ResolveHostSourcePath(path string, followLink bool) (resolvedPath, rebaseName string, err error) {
	if followLink {
		resolvedPath, err = filepath.EvalSymlinks(path)
		if err != nil {
			return
		}

		resolvedPath, rebaseName = GetRebaseName(path, resolvedPath)
	} else {
		dirPath, basePath := filepath.Split(path)

		
		var resolvedDirPath string
		resolvedDirPath, err = filepath.EvalSymlinks(dirPath)
		if err != nil {
			return
		}
		
		
		resolvedPath = resolvedDirPath + string(filepath.Separator) + basePath
		if hasTrailingPathSeparator(path, os.PathSeparator) &&
			filepath.Base(path) != filepath.Base(resolvedPath) {
			rebaseName = filepath.Base(path)
		}
	}
	return resolvedPath, rebaseName, nil
}



func GetRebaseName(path, resolvedPath string) (string, string) {
	
	
	var rebaseName string
	if specifiesCurrentDir(path) &&
		!specifiesCurrentDir(resolvedPath) {
		resolvedPath += string(filepath.Separator) + "."
	}

	if hasTrailingPathSeparator(path, os.PathSeparator) &&
		!hasTrailingPathSeparator(resolvedPath, os.PathSeparator) {
		resolvedPath += string(filepath.Separator)
	}

	if filepath.Base(path) != filepath.Base(resolvedPath) {
		
		
		
		
		rebaseName = filepath.Base(path)
	}
	return resolvedPath, rebaseName
}
