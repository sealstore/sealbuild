package utils

import (
	"archive/tar"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//压缩单个文件
func tarFile(filesource string, sfileInfo os.FileInfo, tarwriter *tar.Writer) error {
	sfile, err := os.Open(filesource)
	if err != nil {
		return err
	}
	defer sfile.Close()
	header, err := tar.FileInfoHeader(sfileInfo, "")
	if err != nil {
		return err
	}
	header.Name = filesource
	err = tarwriter.WriteHeader(header)
	if err != nil {
		return err
	}
	if _, err = io.Copy(tarwriter, sfile); err != nil {
		return err
	}
	return nil
}

//压缩文件夹
func tarFolder(directory string, tarwriter *tar.Writer) error {
	return filepath.Walk(directory, func(targetpath string, file os.FileInfo, err error) error {
		//read the file failure
		if file == nil {
			return err
		}
		if file.IsDir() {
			if directory == targetpath {
				return nil
			}
			header, err := tar.FileInfoHeader(file, "")
			if err != nil {
				return err
			}
			header.Name = filepath.Join(directory, strings.TrimPrefix(targetpath, directory))
			if err = tarwriter.WriteHeader(header); err != nil {
				return err
			}
			os.Mkdir(strings.TrimPrefix(directory, file.Name()), os.ModeDir)
			//如果压缩的目录里面还有目录，则递归压缩
			return tarFolder(targetpath, tarwriter)
		}
		return tarFile(targetpath, file, tarwriter)
	})
}

//untarFile 解压
func untarFile(tarFile string, untarPath string) error {
	//打开要解包的文件，tarFile是要解包的 .tar 文件的路径
	fr, er := os.Open(tarFile)
	if er != nil {
		return er
	}
	defer fr.Close()
	// 创建 tar.Reader，准备执行解包操作
	tr := tar.NewReader(fr)
	//用 tr.Next() 来遍历包中的文件，然后将文件的数据保存到磁盘中
	for hdr, er := tr.Next(); er != io.EOF; hdr, er = tr.Next() {
		if er != nil {
			return er
		}
		//先创建目录
		fileName := untarPath + "/" + hdr.Name
		dir := path.Dir(fileName)
		_, err := os.Stat(dir)
		//如果err 为空说明文件夹已经存在，就不用创建
		if err != nil {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
		//获取文件信息
		fi := hdr.FileInfo()
		//创建空文件，准备写入解压后的数据
		fw, er := os.Create(fileName)
		if er != nil {
			return er
		}
		defer fw.Close()
		// 写入解压后的数据
		_, er = io.Copy(fw, tr)
		if er != nil {
			return er
		}
		// 设置文件权限
		os.Chmod(fileName, fi.Mode().Perm())
	}
	return nil
}
