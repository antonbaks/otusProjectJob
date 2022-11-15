package storage

import (
	"crypto/sha512"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

type Storage struct {
	uploadDir string
}

func NewStorage(uploadDir string) Storage {
	return Storage{
		uploadDir: uploadDir,
	}
}

func (s Storage) CreateUploadDir() error {
	if _, err := os.Stat(s.uploadDir); !os.IsNotExist(err) {
		return nil
	}

	if err := os.Mkdir(s.uploadDir, 0o777); err != nil {
		return err
	}

	return nil
}

func (s Storage) Create(url, width, height string) (*os.File, error) {
	return os.Create(s.pathToFile(url, width, height))
}

func (s Storage) Open(url, width, height string) (*os.File, error) {
	return os.OpenFile(s.pathToFile(url, width, height), os.O_RDWR|os.O_APPEND, os.ModeAppend)
}

func (s Storage) Delete(filename string) error {
	return os.Remove(s.generatePathByFileName(filename))
}

func (s Storage) FileName(url, width, height string) string {
	return s.generateFileName(url, width, height)
}

func (s Storage) pathToFile(url, width, height string) string {
	return s.generatePathByFileName(s.generateFileName(url, width, height))
}

func (s Storage) generatePathByFileName(filename string) string {
	return filepath.Join(s.uploadDir, filename)
}

func (s Storage) generateFileName(url, width, height string) string {
	stringBuilder := strings.Builder{}
	hash := sha512.Sum512_224([]byte(url))

	stringBuilder.WriteString(hex.EncodeToString(hash[:]))
	stringBuilder.WriteString("_")
	stringBuilder.WriteString(width)
	stringBuilder.WriteString("_")
	stringBuilder.WriteString(height)
	stringBuilder.WriteString(".jpg")

	return stringBuilder.String()
}
