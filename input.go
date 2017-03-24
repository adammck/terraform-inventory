package main

import (
	"github.com/adammck/venv"
	"github.com/blang/vfs"
)

func GetInputPath(fs vfs.Filesystem, env venv.Env) string {

	var fn string

	fn = env.Getenv("TF_STATE")
	if fn != "" {
		return fn
	}

	fn = env.Getenv("TI_TFSTATE")
	if fn != "" {
		return fn
	}

	fn = "terraform.tfstate"
	_, err := fs.Stat(fn)
	if err == nil {
		return fn
	}

	return "."
}
