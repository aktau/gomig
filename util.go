package main

import (
	"fmt"
	"io"
	"os"
	"strings"
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

func LoadConfigOrDie(file string) *Config {
	conf, err := LoadConfig(file)
	if err != nil {
		fmt.Printf("error while loading config file: '%v'\n", err)
		fmt.Println("to generate a sample config file use the generate-config command")
		os.Exit(1)
	}

	return conf
}

/* indents every line in str with indent */
func IndentWith(str string, indent string) string {
	nl := "\n"
	spl := strings.Split(str, nl)
	if spl[len(spl)-1] == "" {
		spl = spl[:len(spl)-1]
	}
	for idx, line := range spl {
		spl[idx] = indent + line
	}
	return strings.Join(spl, nl)
}
