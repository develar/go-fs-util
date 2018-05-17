package fsutil

import (
	"os"
	"path/filepath"
	"github.com/develar/errors"
	"io"
)

// Creates the named file and parent directories if need
func CreateFile(name string) (*os.File, error) {
	return open(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

func open(name string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err == nil {
		return file, nil
	}

	if !os.IsNotExist(err) {
		return nil, errors.WithStack(err)
	}

	err = os.MkdirAll(filepath.Dir(name), 0777)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	file, err = os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}

func CopyFile(from string, to string, fromInfo os.FileInfo) error {
	sourceFile, err := os.Open(from)
	if err != nil {
		return errors.WithStack(err)
	}

	defer sourceFile.Close()
	return WriteFile(sourceFile, to, fromInfo, make([]byte, 32*1024))
}

func WriteFile(source io.Reader, to string, fromInfo os.FileInfo, buffer []byte) error {
	// cannot use file mode as is because of *** *** *** umask
	destinationFile, err := open(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = io.CopyBuffer(destinationFile, source, buffer)
	if err != nil {
		destinationFile.Close()
		return errors.WithStack(err)
	}

	perm := fromInfo.Mode().Perm()
	if perm != 0644 {
		err = os.Chmod(to, perm)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	err = destinationFile.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func EnsureEmptyDir(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.WithStack(os.MkdirAll(dirPath, 0777))
		} else {
			return errors.WithStack(err)
		}
	}

	defer dir.Close()

	files, err := dir.Readdirnames(0)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, name := range files {
		err = os.RemoveAll(filepath.Join(dirPath, name))
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func ReadDirContent(dirPath string) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	files, err := dir.Readdirnames(0)
	return files, CloseAndCheckError(err, dir)
}

func CloseAndCheckError(err error, closable io.Closer) error {
	closeErr := closable.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	if closeErr != nil && closeErr != os.ErrClosed {
		return errors.WithStack(closeErr)
	}
	return nil
}