package fsutil

import (
	"os"
	"path/filepath"
	"github.com/develar/errors"
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




