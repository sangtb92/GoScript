package main

import (
	"github.com/rjeczalik/notify"
	"log"
	"fmt"
	"net/smtp"
)

const (
	WORK_DIR        = "/home/sangnd/..."
	WORK_DIR_COMMON = "/home/sangnd/Desktop/..."
	buffer          = 8192
	from            = ""
	password        = ""
	to              = ""
	subject         = "Report file has changed!"
)

func main() {

	c := make(chan notify.EventInfo, buffer)
	done := make(chan bool)
	//go watchUsr(c)
	go watchCommon(c)
	go watchResult(c)
	<-done
	defer notify.Stop(c)
	var input string
	fmt.Scanln(&input)
}

//set watcher
func watchUsr(c chan notify.EventInfo) {
	for {
		if err := notify.Watch(WORK_DIR, c, notify.Rename); err != nil {
			log.Fatal(err)
		}
	}
}

//set watcher
func watchCommon(c chan notify.EventInfo) {
	for {
		if err := notify.Watch(WORK_DIR_COMMON, c, notify.Rename); err != nil {
			log.Fatal(err)
		}
	}
}

func watchResult(c chan notify.EventInfo) {
	for {
		event := <-c
		sendMailRs(event.Path())
	}
}

func sendMailRs(filePath string) {
	msg := genMsg(filePath, subject)
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from,
		[]string{to},
		[]byte(msg))
	if err != nil {
		log.Fatalf("smtp errors: %s", err)
		return
	}
	log.Printf("There is a mail sent to address: %s, please check it!", to)
}

func genMsg(filePath, subject string) string {
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" + filePath + "\n"
	return msg
}
