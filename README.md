# Repository-based File System Wrapper

<!-- > 2018-04-03T11:44:19+0800 -->

The object-oriented file system wrapped as repository in Go packaged as 'fs'.

## Motivations

A folder/directory in an OS works like a repository where you can put and fetch stuff by id(filename), hence needed to be wrapped as a package(library).

## Installation

```bash
go get github.com/zhanbei/fs
```

## Usage

Codes are in [examples/main.go](examples/main.go).

```go

package main

import (
	"os"
	"fmt"
	"time"

	"github.com/zhanbei/fs"
)

var TestRepository *fs.Repository

// Your default file and directory permissions.
const (
	// -rw-rw-r--
	FilePermission = 0664
	// -rwxrwxr-x
	DirPermission = 0775
)

func init() {
	// Make sure these directories exists
	err := os.MkdirAll("test-dir-name/sub-dir", DirPermission)
	if err != nil {
		panic("Creating directories failed")
	}

	// Initialize repositories.
	TestRepository = InitRepositoryAndPanic("test-dir-name")
}

// Init a repository and panic if error encountered.
func InitRepositoryAndPanic(dirPath string) *fs.Repository {
	repository, err := fs.NewRepositoryWithPermissions(dirPath, FilePermission, DirPermission)
	if err != nil {
		// panic only can be used in 'init()' func.
		panic("Expected directory: [" + dirPath + "] didn't exist or is not a directory: " + err.Error())
	}
	return repository
}

func main() {
	testFile := "sub-dir/test.txt"

	// Create a file.
	err := TestRepository.WriteFile(testFile, []byte("This is a file for test."))
	if err != nil {
		panic("Failed to write file:" + err.Error())
	}
	fmt.Println("Wrote files [" + testFile + "] successfully.")

	time.Sleep(5 * time.Second)

	// Remove a file.
	err = TestRepository.RemoveFile(testFile)
	if err != nil {
		panic("Failed to remove file:" + err.Error())
	}
	fmt.Println("Removed files [" + testFile + "] successfully.")

	time.Sleep(3 * time.Second)

	// Remove cdn folder(empty filename refers the repository folder.).
	err = TestRepository.RemoveFile("")
	if err != nil {
		fmt.Println("Cannot remove repository folder: " + err.Error())
		return
	}
	panic("Unexpectedly removed folder: " + TestRepository.Path)
}

```

## Investigate the Detail of the Project

Currently the project, containing two files, [fs.go](fs.go) and [repository.go](repository.go), as shown below, is simple and easy to understand.

```go
// fs.go
package fs

import (
	"os"
	"errors"
)

var (
	ErrNotExist = os.ErrNotExist

	ErrNotDir = errors.New("file is not a directory")
)

// Default permissions for files/directories being created.
// @see http://permissions-calculator.org/
const (
	// -rw-rw-r--
	DefaultFilePermission = 0664
	// -rwxrwxr-x
	DefaultDirPermission = 0775
)

// Create a `Repository` instance with default permissions set.
func NewRepository(dirPath string) (*Repository, error) {
	return NewRepositoryWithPermissions(dirPath, DefaultFilePermission, DefaultDirPermission)
}

// Create a `Repository` instance with dirPath and permissions.
// Error will be returned if dirPath does not exist or is not a directory.
// @see https://www.thegeekstuff.com/2010/04/unix-file-and-directory-permissions/
func NewRepositoryWithPermissions(dirPath string, newFilePerm, newDirPerm os.FileMode) (*Repository, error) {
	stat, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return nil, ErrNotExist
	}
	if !stat.IsDir() {
		return nil, ErrNotDir
	}
	return &Repository{
		dirPath,
		newFilePerm,
		newDirPerm,
	}, nil
}
```

```go
// repository.go
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
```
