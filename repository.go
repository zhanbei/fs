package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// The Directory struct but named as `Repository`.
type Repository struct {
	// Path of the repository.
	Path string `json:"path"`
	// Permission for a file being created.
	FilePermission os.FileMode
	// Permission for a directory being created.
	DirPermission os.FileMode
}

func (m *Repository) GetFilePath(filename string) string {
	return filepath.Join(m.Path, filename)
}

// @see https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
func (m *Repository) IsFileExists(filename string) bool {
	_, err := os.Stat(m.GetFilePath(filename))
	return !os.IsNotExist(err)
}

func (m *Repository) CreateDirectory(dirPath string) error {
	return os.MkdirAll(m.GetFilePath(dirPath), m.DirPermission)
}

func (m *Repository) RemoveDirectory(dirPath string) error {
	return m.RemoveFile(dirPath)
}

func (m *Repository) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(m.GetFilePath(filename))
}

func (m *Repository) WriteFile(filename string, data []byte) error {
	return ioutil.WriteFile(m.GetFilePath(filename), data, m.FilePermission)
}

func (m *Repository) AppendFile(filename string, data []byte) error {
	f, err := os.OpenFile(m.GetFilePath(filename), os.O_APPEND|os.O_WRONLY, m.FilePermission)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (m *Repository) RenameFile(filename, newFilename string) error {
	return os.Rename(m.GetFilePath(filename), m.GetFilePath(newFilename))
}

func (m *Repository) RemoveFile(filename string) error {
	return os.Remove(m.GetFilePath(filename))
}
