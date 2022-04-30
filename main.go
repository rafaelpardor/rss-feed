package main

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/mmcdole/gofeed"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func buildMessage(mail Mail) string {
	var msg string

	mail.Body += getKubernetesFeed(mail.Body)

	msg += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func getKubernetesFeed(msg string) string {
	URL := "https://kubernetes.io/feed.xml"
	html := `<a href='%s' target=_blank><p>%s</p></a>` + "\n"
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(URL)
	msg += feed.Title + "\n"
	for i := 0; i < len(feed.Items); i++ {
		feedi := feed.Items[i]
		msg += fmt.Sprintf(html, feedi.Link, feedi.Title)
	}

	return msg
}

func main() {

	from := "rafaelpardo98@gmail.com"
	password := ""
	smtpAddr := "smtp.gmail.com"
	smtpHost := "smtp.gmail.com:587"

	to := []string{
		"rafaelpardo98@gmail.com",
	}

	request := Mail{
		Sender:  from,
		To:      to,
		Subject: "Weekly Update",
		Body:    "",
	}
	message := buildMessage(request)

	auth := smtp.PlainAuth("", from, password, smtpAddr)
	err := smtp.SendMail(smtpHost, auth, from, to, []byte(message))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Email Sent Successfully!")
}
