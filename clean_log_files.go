package main

import (
	"io/ioutil"
	"log"
	"net/smtp"
	"path"
	"os"
	"fmt"
)

const LOG_DIR = "/var/log/"
const LOG_FILE_SIZE = 1000000000
const FROM_EMAIL = "spiraltest@vnext.vn"

func main() {
	var validFileNames []string
	var deletedFiles []string
	files, err := ioutil.ReadDir(LOG_DIR)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.Size() > LOG_FILE_SIZE {
			validFileNames = append(validFileNames, path.Join(LOG_DIR, f.Name()))
		}
	}
	for i := range (validFileNames) {
		deletedFile := removeFile(validFileNames[i])
		if len((deletedFile)) > 0 {
			newFile := createFile(deletedFile)
			if len(newFile) > 0 {
				deletedFiles = append(deletedFiles, newFile)
			}
		}
	}

	if len(deletedFiles) > 0 {
		sendMail(deletedFiles)
	} else {
		fmt.Println("Has no file has been deleted!")
	}
	fmt.Println("Done job!")
}

func sendMail(fileNames []string) {
	from := FROM_EMAIL
	pass := "iujebqtqslortspq"
	to := "hongtt@vnext.vn"
	msg := generateMsg(from, to, "List log files have cleaned", fileNames)
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
	if err != nil {
		log.Fatalf("smtp errors: %s", err)
		return
	}
	log.Printf("There is a mail sent to address: %s, please check it!", to)
}

func generateMsg(from, to, subject string, fileNames []string) string {
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n"
	body := ""

	for _, fileName := range fileNames {
		body = body + fileName + "\n"
	}
	return msg + body
}

func removeFile(path string) string {
	result := path
	if !isDir(path) {
		err := os.Remove(path)
		if (isError(err)) {
			result = "";
		}
	}

	return result
}

func createFile(path string) string {
	result := path
	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if (isError(err)) {
			result = "";
		}
		defer file.Close()
	}
	return result
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err)
	}

	return (err != nil)
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		return fi.IsDir()
	}
	return false
}
