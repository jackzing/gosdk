// Hyperchain License
// Copyright (C) 2016 The Hyperchain Authors.

package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// ExitWithLockFile helps Hyperchain stop successfully together with creating a lock
// file which is used to prevent Hyperchain daemon(hocMonitor) restarting Hyperchain
// repeatedly.
// Normal flag indicates stop Hyperchain normally or not.
func ExitWithLockFile(fileName string, normal bool, reason string) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Open lock file failed, reason: %s\n", err)
		return
	}
	defer closeFile(file)

	reason = fmt.Sprintf("%s\n"+
		"Please make sure that you have deleted the lock file before restart Hyperchain.", reason)

	_, err = file.WriteString(reason)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Write reason to lock file failed, reason: %s\n", err)
		os.Exit(1)
	}

	err = file.Sync()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Sync lock file failed, reason: %s\n", err)
		os.Exit(1)
	}

	_, _ = fmt.Fprintln(os.Stderr, reason)

	if normal {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// FileExist checks a file exists or not
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// GetGoPath gets the GOPATH in this environment
func GetGoPath() string {
	env := os.Getenv("GOPATH")
	l := strings.Split(env, ":")
	if len(l) > 1 {
		return l[len(l)-1]
	}
	return l[0]
}

// GetPath get complete path for namespace level file
func GetPath(namespace, shortPath string) string {
	if len(namespace) == 0 {
		return shortPath
	}
	return path.Join("namespaces", namespace, shortPath)
}

//GetLongNSPath get complete path for short path level file in namespace
func GetLongNSPath(nsPrefix, namespace, shortPath string) string {
	if len(nsPrefix) == 0 {
		return shortPath
	}
	return path.Join(nsPrefix, namespace, shortPath)
}

// SeekAndAppend seek item by pattern in file whose path is filepath and than append content after that item
func SeekAndAppend(item, filePath, appendContent string) error {

	//1. check validity of arguments.
	if len(filePath) == 0 {
		return fmt.Errorf("invalid filePath, file path is empty")
	}

	if len(item) == 0 {
		return fmt.Errorf("invalid iterm, item is empty")
	}

	if len(appendContent) == 0 {
		return fmt.Errorf("invalid iterm, appenContent is empty")
	}

	//2. create a temp file to store new file content.

	newFilePath := filePath + ".tmp"

	newFile, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	defer closeFile(newFile)

	//3. copy and append appendContent

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer closeFile(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if _, err = newFile.WriteString(line + "\n"); err != nil {
			return err
		}
		if strings.Contains(line, item) {
			if _, err = newFile.WriteString(appendContent + "\n"); err != nil {
				return err
			}
		}
	}

	//4. rename .tmp file to original file name
	if err = os.Rename(newFilePath, filePath); err != nil {
		return err
	}
	return nil
}

//closeFile close file and print if have error
func closeFile(file *os.File) {
	if file == nil {
		return
	}
	err := file.Close()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Close file failed, reason: %s\n", err)
	}
}

//GetProjectPath get complete path for project
func GetProjectPath() string {
	goPath := os.Getenv("GOPATH")
	goRoot := os.Getenv("GOROOT")
	if strings.Contains(goPath, ":") {
		goPathes := strings.Split(goPath, ":")
		var realPathes []string
		for _, p := range goPathes {
			if p == goRoot {
				continue
			}
			realPathes = append(realPathes, p)
		}
		goPath = realPathes[0]
	}
	return goPath + "/src/github.com/hyperchain/hyperchain"
}

// CopyDir copys all files under srcDir to dstDir
func CopyDir(dstDir, srcDir string) error {
	fileInfos, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	dstInfos, oerr := ioutil.ReadDir(dstDir)
	if len(dstInfos) > 0 {
		return errors.New("now dstDir must be empty, please check:" + dstDir)
	}
	if os.IsNotExist(oerr) {
		if err := os.Mkdir(dstDir, os.ModePerm); err != nil {
			return err
		}
	}

	for i := 0; i < len(fileInfos); i++ {
		dst := path.Join(dstDir, fileInfos[i].Name())
		src := path.Join(srcDir, fileInfos[i].Name())
		if fileInfos[i].IsDir() {
			if err := CopyDir(dst, src); err != nil {
				return err
			}
		} else {
			if err := CopyFile(dst, src); err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile copys a single file
// will return error if there are files with same relative path and name in src and dst
func CopyFile(dst, src string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = dstFile.Close()
	}()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}
