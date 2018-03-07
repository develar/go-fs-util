package fsutil

import (
	"os"
	"path/filepath"
	"github.com/develar/errors"
	"io"
)

// Creates the named file and parent directories if need
func CreateFile(name string) (*os.File, error) {
	file, err := os.Create(name)
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

	file, err = os.Create(name)
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
	return WriteFile(sourceFile, to, fromInfo)
}

func WriteFile(source io.Reader, to string, fromInfo os.FileInfo) error {
	// cannot use file mode as is because of *** *** *** umask
	destinationFile, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = io.Copy(destinationFile, source)
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
