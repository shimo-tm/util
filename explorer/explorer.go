// 文件资源管理器

package explorer

import (
	"io"
	"os"
	"path"
)

// IsExist 验证文件or目录是否存在
func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// IsDir 是否目录
func IsDir(filePath string) bool {
	fsInfo, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return fsInfo.IsDir()
}

// CopyFile  复制文件
func CopyFile(src, dest string) error {
	if !IsExist(src) {
		return os.ErrNotExist
	}

	srcFp, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFp.Close()

	destFp, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFp.Close()

	if _, err = io.Copy(destFp, srcFp); err != nil {
		return err
	}
	return nil
}

// CopyDir 复制目录
func CopyDir(src, dest string) error {
	if !IsDir(src) {
		return os.ErrNotExist
	}
	// 遍历目录
	fsInfoList, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if !IsDir(dest) {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return err
		}
	}
	for _, fsInfo := range fsInfoList {
		if fsInfo.IsDir() {
			// 递归
			err := CopyDir(path.Join(src, fsInfo.Name()), path.Join(dest, fsInfo.Name()))
			if err != nil {
				return err
			}
			continue
		}
		// 文件
		CopyFile(path.Join(src, fsInfo.Name()), path.Join(dest, fsInfo.Name()))
	}
	return nil
}

// MoveFile 移动文件
func MoveFile(src, dest string) error {
	if !IsExist(src) {
		return os.ErrNotExist
	}
	srcFp, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFp.Close()

	destFp, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFp.Close()

	_, err = io.Copy(destFp, srcFp)
	if err != nil {
		return err
	}

	err = os.Remove(src)
	if err != nil {
		return err
	}
	return nil
}

// MoveDir 移动目录
func MoveDir(src, dest string) error {
	if !IsDir(src) {
		return os.ErrNotExist
	}
	// 遍历目录
	fsInfoList, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if !IsDir(dest) {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			return err
		}
	}
	for _, fsInfo := range fsInfoList {
		if fsInfo.IsDir() {
			// 递归
			err := MoveDir(path.Join(src, fsInfo.Name()), path.Join(dest, fsInfo.Name()))
			if err != nil {
				return err
			}
			continue
		}
		// 文件
		CopyFile(path.Join(src, fsInfo.Name()), path.Join(dest, fsInfo.Name()))
	}
	err = os.RemoveAll(src)
	if err != nil {
		return err
	}
	return nil
}
