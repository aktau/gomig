package main

import (
	"io"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CopyFile(src, dst string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer df.Close()

	return io.Copy(df, sf)
}
