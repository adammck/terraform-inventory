package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/adammck/venv"
	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/stretchr/testify/assert"
)

func TestGetInputPath(t *testing.T) {
	assert.Equal(t, ".", GetInputPath(memfs.Create(), venv.Mock()))
	assert.Equal(t, "aaa", GetInputPath(memfs.Create(), envWith(map[string]string{"TF_STATE": "aaa"})))
	assert.Equal(t, "bbb", GetInputPath(memfs.Create(), envWith(map[string]string{"TI_TFSTATE": "bbb"})))
	assert.Equal(t, "terraform.tfstate", GetInputPath(fsWithFiles([]string{"terraform.tfstate"}), venv.Mock()))
	assert.Equal(t, ".", GetInputPath(fsWithFiles([]string{".terraform/terraform.tfstate"}), venv.Mock()))
	assert.Equal(t, "terraform", GetInputPath(fsWithDirs([]string{"terraform"}), envWith(map[string]string{"TF_STATE": "terraform"})))
}

func envWith(env map[string]string) venv.Env {
	e := venv.Mock()

	for k, v := range env {
		e.Setenv(k, v)
	}

	return e
}

func fsWithFiles(filenames []string) vfs.Filesystem {
	fs := memfs.Create()
	var err error

	for _, fn := range filenames {

		path := filepath.Dir(fn)
		if path != "" {
			err = vfs.MkdirAll(fs, path, 0700)
			if err != nil {
				panic(err)
			}
		}

		err = touchFile(fs, fn)
		if err != nil {
			panic(err)
		}
	}

	return fs
}

func fsWithDirs(dirs []string) vfs.Filesystem {
	fs := memfs.Create()

	var err error

	for _, fp := range dirs {
		err = vfs.MkdirAll(fs, fp, 0700)
		if err != nil {
			panic(err)
		}
	}

	return fs
}

// TODO: Upgrade this later with file contents.
func touchFile(fs vfs.Filesystem, filename string) error {
	return writeFile(fs, filename, []byte{}, 0600)
}

// port of ioutil.Writefile for vfs
func writeFile(fs vfs.Filesystem, filename string, data []byte, perm os.FileMode) error {
	f, err := fs.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
