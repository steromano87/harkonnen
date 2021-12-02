package compiler

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

type File struct {
	Path               string
	Hash               string
	MnemonicName       string
	CompiledObjectPath string
}

func NewCompilableFile(path string, mnemonicName string) (*File, error) {
	script := new(File)

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	script.Path = path
	script.MnemonicName = mnemonicName
	hash, err := script.calculateHash()
	if err != nil {
		return nil, err
	}

	script.Hash = hash

	return script, nil
}

func (script File) HasChanged() (bool, error) {
	newHash, err := script.calculateHash()
	if err != nil {
		return true, err
	}

	return newHash != script.Hash, nil
}

func (script File) calculateHash() (string, error) {
	file, err := os.Open(script.Path)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = file.Close()
	}()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes), nil
}
