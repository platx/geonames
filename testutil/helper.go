package testutil

import "io/fs"

func MustOpen(fs fs.FS, name string) fs.File {
	f, err := fs.Open(name)
	if err != nil {
		panic(err)
	}

	return f
}
