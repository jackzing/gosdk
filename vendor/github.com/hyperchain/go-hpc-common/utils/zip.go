package utils

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

// Compress compress the input files to dest
func Compress(files []*os.File, dest string) error {
	d, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() { _ = d.Close() }()
	w := zip.NewWriter(d)
	defer func() { _ = w.Close() }()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			_ = os.Remove(dest)
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = path.Join(prefix, info.Name())
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(path.Join(file.Name(), fi.Name()))
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = path.Join(prefix, header.Name)
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		_ = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// DeCompress the file named zipFile to dest
func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		filename := path.Join(dest, file.Name)
		err = os.MkdirAll(path.Dir(filename), 0755)
		if err != nil {
			_ = rc.Close()
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			_ = rc.Close()
			return err
		}
		_, err = io.Copy(w, rc)
		if err != nil {
			_ = rc.Close()
			_ = w.Close()
			return err
		}
		_ = rc.Close()
		_ = w.Close()
	}
	return nil
}
