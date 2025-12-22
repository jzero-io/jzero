package filex

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// FileExists check file exist
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return info.IsDir()
}

// copyDir 递归复制目录，将 src 目录下的所有文件复制到 dst 目录，已存在的文件会被覆盖
func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// 计算目标路径
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// 创建目录
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 复制文件
		return CopyFile(path, dstPath)
	})
}

// copyFile 复制单个文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 确保目标目录存在
	if err = os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
