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
