package http

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"tiktok-uploader/log"
)

func Open(f *multipart.FileHeader) (multipart.File, error) {
	file, err := f.Open()
	if file != nil {
		defer Close(file)
	}
	return file, err
}

func ReadAll(f *multipart.FileHeader) ([]byte, error) {
	file, err := Open(f)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func Close(file io.Closer)  {
	err := file.Close()
	if err != nil {
		log.Error(err)
	}
}