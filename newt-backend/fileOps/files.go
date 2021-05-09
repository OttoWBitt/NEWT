package fileOps

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const maxUploadSize = 4 * 1024 * 1024 // 4 mb
const UploadPath = "./files"

func UploadFileHandler(res http.ResponseWriter, req *http.Request) (*string, *string, error, int) {
	if err := req.ParseMultipartForm(maxUploadSize); err != nil {
		erro := fmt.Sprintf("Could not parse multipart form: %v\n", err)
		return nil, nil, errors.New(erro), http.StatusInternalServerError
	}

	// parse and validate file and post parameters
	file, fileHeader, err := req.FormFile("artifactFile")
	if err != nil {
		erro := fmt.Sprintf("invalid_file -  %v\n", err)
		return nil, nil, errors.New(erro), http.StatusBadRequest
	}

	defer file.Close()

	// Get and print out file size
	fileSize := fileHeader.Size

	//fmt.Printf("File size (bytes): %v\n", fileSize)
	// validate file size
	if fileSize > maxUploadSize {
		return nil, nil, errors.New("file_too_big"), http.StatusBadRequest
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		erro := fmt.Sprintf("invalid_file -  %v\n", err)
		return nil, nil, errors.New(erro), http.StatusBadRequest
	}

	// check file type, detectcontenttype only needs the first 512 bytes
	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "image/jpeg", "image/jpg":
	case "image/gif", "image/png":
	case "application/pdf":
		break
	default:
		return nil, nil, errors.New("invalid_file_extension"), http.StatusBadRequest
	}

	fileNameToken := randToken(12)
	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		erro := fmt.Sprintf("can't_read_file_type -  %v\n", err)
		return nil, nil, errors.New(erro), http.StatusInternalServerError
	}

	fileName := filepath.Join(UploadPath, fileNameToken+fileEndings[0])
	fmt.Printf("FileType: %s, File: %s\n", detectedFileType, fileName)

	const dns = "http://newt.ottobittencourt.com:3001/api/%s"
	url := fmt.Sprintf(dns, fileName)

	// write file
	newFile, err := os.Create(fileName)
	if err != nil {
		erro := fmt.Sprintf("can't_write_file -  %v\n", err)
		return nil, nil, errors.New(erro), http.StatusInternalServerError
	}

	defer newFile.Close() // idempotent, okay to call twice
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		erro := fmt.Sprintf("can't_write_file -  %v\n", err)
		return nil, nil, errors.New(erro), http.StatusInternalServerError
	}

	return &fileHeader.Filename, &url, nil, http.StatusOK
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
