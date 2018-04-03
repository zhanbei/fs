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
