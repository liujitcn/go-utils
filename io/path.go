package io

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

func GetWorkingDirPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func GetExePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	return exePath
}

func GetAbsPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func GetFileList(root string) []string {
	var files []string

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	}); err != nil {
		return nil
	}

	return files
}

func GetFolderNameList(root string) []string {
	var names []string
	fs, _ := os.ReadDir(root)
	for _, file := range fs {
		if file.IsDir() {
			names = append(names, file.Name())
		}
	}
	return names
}
func MatchPath(pattern string, path string) bool {
	if g, err := glob.Compile(pattern); err == nil {
		return g.Match(path)
	}

	return false
}

func ExpandUser(path string) (string, error) {
	if u, err := user.Current(); err == nil {
		fullTilde := fmt.Sprintf("~%s", u.Name)

		if strings.HasPrefix(path, `~/`) || path == `~` {
			return strings.Replace(path, `~`, u.HomeDir, 1), nil
		}

		if strings.HasPrefix(path, fullTilde+`/`) || path == fullTilde {
			return strings.Replace(path, fullTilde, u.HomeDir, 1), nil
		}

		return path, nil
	} else {
		return path, err
	}
}

func IsNonemptyExecutableFile(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.Size() > 0 && (stat.Mode().Perm()&0111) != 0 {
		return true
	}

	return false
}

func IsNonemptyFile(path string) bool {
	if FileExists(path) {
		if stat, err := os.Stat(path); err == nil && stat.Size() > 0 {
			return true
		}
	}

	return false
}

func IsNonemptyDir(path string) bool {
	if DirExists(path) {
		if entries, err := ioutil.ReadDir(path); err == nil && len(entries) > 0 {
			return true
		}
	}

	return false
}

func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func LinkExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return IsSymlink(stat.Mode())
	}

	return false
}

func FileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.Mode().IsRegular()
	}

	return false
}

func DirExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}

	return false
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsSymlink(mode os.FileMode) bool {
	return mode&os.ModeSymlink != 0
}

func IsDevice(mode os.FileMode) bool {
	return mode&os.ModeDevice != 0
}

func IsCharDevice(mode os.FileMode) bool {
	return mode&os.ModeCharDevice != 0
}

func IsNamedPipe(mode os.FileMode) bool {
	return mode&os.ModeNamedPipe != 0
}

func IsSocket(mode os.FileMode) bool {
	return mode&os.ModeSocket != 0
}

func IsSticky(mode os.FileMode) bool {
	return mode&os.ModeSticky != 0
}

func IsSetuid(mode os.FileMode) bool {
	return mode&os.ModeSetuid != 0
}

func IsSetgid(mode os.FileMode) bool {
	return mode&os.ModeSetgid != 0
}

func IsTemporary(mode os.FileMode) bool {
	return mode&os.ModeTemporary != 0
}

func IsExclusive(mode os.FileMode) bool {
	return mode&os.ModeExclusive != 0
}

func IsAppend(mode os.FileMode) bool {
	return mode&os.ModeAppend != 0
}

func IsReadable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_RDONLY, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}

func IsWritable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_WRONLY, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}

func IsAppendable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_APPEND, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}
