package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/elgs/gosqljson"
)

const maxUploadSize = 2 * 1024 * 1024 // 2 mb
const uploadPath = "./files"

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			t, _ := template.ParseFiles("html/upload.html")
			t.Execute(res, nil)
			return
		}
		if err := req.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(res, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}

		stringId := req.FormValue("userID")
		subject := req.FormValue("subject")

		userID, err := strconv.Atoi(stringId)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		auth := userExists(userID)

		fmt.Println(auth) //esse cara tem que vir do cookie do front, falta implementar

		if auth != nil {
			renderError(res, "User not logged in", http.StatusInternalServerError)
			return
		}

		// parse and validate file and post parameters
		file, fileHeader, err := req.FormFile("uploadFile")
		if err != nil {
			renderError(res, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		fmt.Println(fileHeader.Filename)
		defer file.Close()
		// Get and print out file size
		fileSize := fileHeader.Size
		fmt.Printf("File size (bytes): %v\n", fileSize)
		// validate file size
		if fileSize > maxUploadSize {
			renderError(res, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(res, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		detectedFileType := http.DetectContentType(fileBytes)
		switch detectedFileType {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
		case "application/pdf":
			break
		default:
			renderError(res, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileNameToken := randToken(12)
		fileEndings, err := mime.ExtensionsByType(detectedFileType)
		if err != nil {
			renderError(res, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
		fileName := filepath.Join(uploadPath, fileNameToken+fileEndings[0])
		fmt.Printf("FileType: %s, File: %s\n", detectedFileType, fileName)

		const dns = "http://localhost:3000/%s"
		url := fmt.Sprintf(dns, fileName)

		erro := saveFileMetadata(fileHeader.Filename, url, userID, subject)
		if erro != nil {
			renderError(res, "CANT_SAVE_METADATA", http.StatusInternalServerError)
			return
		}

		// write file
		newFile, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
			renderError(res, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(res, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		res.Write([]byte("SUCCESS"))
	})
}

func renderError(res http.ResponseWriter, message string, statusCode int) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func saveFileMetadata(fileName string, codedName string, userID int, subject string) error {

	_, err := dbMaster.Exec(
		`INSERT INTO 
			newt.file_metadata(
				user_id,
				name,
				url,
				subject
			) 
		VALUES (?,?,?,?)
	`, userID, fileName, codedName, subject)

	if err != nil {
		return err
	}
	return nil
}

func download(res http.ResponseWriter, req *http.Request) {
	const theCase = "original"

	const query = `
	SELECT
		*
	FROM 
		newt.file_metadata 
	`

	links, err := gosqljson.QueryDbToMap(dbMaster, theCase, query)
	if err != nil {
		log.Fatal(err)
	}

	generateJSON := map[string]interface{}{
		"links": links,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}
