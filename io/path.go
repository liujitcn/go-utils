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

// GetWorkingDirPath 获取当前工作目录路径。
func GetWorkingDirPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// GetExePath 获取当前可执行文件所在目录。
func GetExePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	return exePath
}

// GetAbsPath 获取当前启动程序的绝对目录路径。
func GetAbsPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

// GetFileList 递归获取目录下的所有文件路径。
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

// GetFolderNameList 获取目录下的一级文件夹名称列表。
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

// MatchPath 判断路径是否匹配指定模式。
func MatchPath(pattern string, path string) bool {
	if g, err := glob.Compile(pattern); err == nil {
		return g.Match(path)
	}

	return false
}

// ExpandUser 展开路径中的用户主目录标记。
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

// IsNonemptyExecutableFile 判断路径是否为非空可执行文件。
func IsNonemptyExecutableFile(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.Size() > 0 && (stat.Mode().Perm()&0111) != 0 {
		return true
	}

	return false
}

// IsNonemptyFile 判断路径是否为非空普通文件。
func IsNonemptyFile(path string) bool {
	if FileExists(path) {
		if stat, err := os.Stat(path); err == nil && stat.Size() > 0 {
			return true
		}
	}

	return false
}

// IsNonemptyDir 判断路径是否为非空目录。
func IsNonemptyDir(path string) bool {
	if DirExists(path) {
		if entries, err := ioutil.ReadDir(path); err == nil && len(entries) > 0 {
			return true
		}
	}

	return false
}

// Exists 判断路径是否存在。
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

// LinkExists 判断路径是否对应符号链接。
func LinkExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return IsSymlink(stat.Mode())
	}

	return false
}

// FileExists 判断路径是否为普通文件。
func FileExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.Mode().IsRegular()
	}

	return false
}

// DirExists 判断路径是否为目录。
func DirExists(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}

	return false
}

// PathExist 判断路径是否存在。
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

// IsSymlink 判断文件模式是否表示符号链接。
func IsSymlink(mode os.FileMode) bool {
	return mode&os.ModeSymlink != 0
}

// IsDevice 判断文件模式是否表示设备文件。
func IsDevice(mode os.FileMode) bool {
	return mode&os.ModeDevice != 0
}

// IsCharDevice 判断文件模式是否表示字符设备。
func IsCharDevice(mode os.FileMode) bool {
	return mode&os.ModeCharDevice != 0
}

// IsNamedPipe 判断文件模式是否表示命名管道。
func IsNamedPipe(mode os.FileMode) bool {
	return mode&os.ModeNamedPipe != 0
}

// IsSocket 判断文件模式是否表示套接字。
func IsSocket(mode os.FileMode) bool {
	return mode&os.ModeSocket != 0
}

// IsSticky 判断文件模式是否包含粘滞位。
func IsSticky(mode os.FileMode) bool {
	return mode&os.ModeSticky != 0
}

// IsSetuid 判断文件模式是否包含 setuid 位。
func IsSetuid(mode os.FileMode) bool {
	return mode&os.ModeSetuid != 0
}

// IsSetgid 判断文件模式是否包含 setgid 位。
func IsSetgid(mode os.FileMode) bool {
	return mode&os.ModeSetgid != 0
}

// IsTemporary 判断文件模式是否表示临时文件。
func IsTemporary(mode os.FileMode) bool {
	return mode&os.ModeTemporary != 0
}

// IsExclusive 判断文件模式是否包含独占位。
func IsExclusive(mode os.FileMode) bool {
	return mode&os.ModeExclusive != 0
}

// IsAppend 判断文件模式是否包含追加写标记。
func IsAppend(mode os.FileMode) bool {
	return mode&os.ModeAppend != 0
}

// IsReadable 判断文件是否可读。
func IsReadable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_RDONLY, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}

// IsWritable 判断文件是否可写。
func IsWritable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_WRONLY, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}

// IsAppendable 判断文件是否可追加写。
func IsAppendable(filename string) bool {
	if f, err := os.OpenFile(filename, os.O_APPEND, 0); err == nil {
		defer f.Close()
		return true
	} else {
		return false
	}
}
