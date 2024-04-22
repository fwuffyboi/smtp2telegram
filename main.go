package main

import (
	"bufio"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

const (
	// VERSION Version is the current version of the program
	VERSION = "0.1.0"

	// Username to authenticate to the SMTP server
	username = "user"
	// Password to authenticate to the SMTP server
	password = "password"
	// Port to listen on
	port = 25
	// Host to listen on
	host = "0.0.0.0"
	// Telegram bot token
	telegramBotToken = "6652952174:AAH73dKp7N78SikhicHe3FIljOHXoKRU2f0"
	// Telegram chat id
	telegramChatID = 6291710804
	// Telegram message template
	telegramMessageTemplate = "From: %%from%%\n" +
		"To:%%to%%\n" +
		"Subject: %%sub%%\n" +
		"Body: %%body%%\n"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	conn.Write([]byte("220 Welcome\r\n"))

	inData := false
	var body strings.Builder
	var from, to string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from connection: %s", err)
			break
		}

		fmt.Println("Received:", line)

		// Handle SMTP commands
		if strings.HasPrefix(line, "HELO") {
			conn.Write([]byte("250 Hello, pleased to meet you\r\n"))
		} else if strings.HasPrefix(line, "MAIL FROM:") {
			conn.Write([]byte("250 Ok\r\n"))
			from = strings.TrimPrefix(line, "MAIL FROM:")
		} else if strings.HasPrefix(line, "RCPT TO:") {
			conn.Write([]byte("250 Ok\r\n"))
			to = strings.TrimPrefix(line, "RCPT TO:")
		} else if strings.HasPrefix(line, "DATA") {
			conn.Write([]byte("354 End data with <CR><LF>.<CR><LF>\r\n"))
			inData = true
		} else if inData && strings.HasPrefix(line, ".") {
			conn.Write([]byte("250 Ok\r\n"))
			inData = false
			fmt.Println("From:", from)
			fmt.Println("To:", to)
			fmt.Println("Body:", body.String())
			body.Reset()
		} else if inData {
			body.WriteString(line)
		} else if strings.HasPrefix(line, "QUIT") {
			conn.Write([]byte("221 Bye\r\n"))
			break
		} else if !inData {
			conn.Write([]byte("500 Error: command not recognized\r\n"))
		}
	}
}

func main() {
	// Start listening for incoming connection

	log.Infof("Starting smtp2telegram...")
	log.Infof("Application version: %s", VERSION)

	var listenaddr = fmt.Sprintf("%s:%d", host, port)

	ln, err := net.Listen("tcp", listenaddr)
	if err != nil {
		log.Errorf("Error listening: %s", err)
		return
	}
	defer ln.Close()

	log.Infof("Fake mail server listening successfully on port %s", "25")

	// start a thread to send a message to myself every 10s
	go func() {
		for {
			sendmail()
			time.Sleep(10 * time.Second)
		}
	}()

	// Accept incoming connections and handle them
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Errorf("Error accepting connection: %s", err)
			continue
		}
		go handleConnection(conn)
	}
}

func sendTelegramMessage(message string) {
	// Send the message to Telegram

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false // assume we are in production

	log.Infof("Authorized on account %s", bot.Self.UserName)

	msg := tgbotapi.NewMessage(telegramChatID, message)
	bot.Send(msg)

}
