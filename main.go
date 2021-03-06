package main

import (
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

func main() {
	wac, err := whatsapp.NewConn(10 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	err = login(wac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return
	}

	// <-time.After(3 * time.Second)
	mensagemBytes, _ := ioutil.ReadFile("./message.txt")
	mensagem := string(mensagemBytes)

	contactsFile, _ := os.Open("./contacts.txt")
	defer contactsFile.Close()
	reader := csv.NewReader(contactsFile)

	for {

		rec, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		// get float value
		contactNumber := rec[0]
		contactNumber = fixPhoneNumber(contactNumber)

		retyingCount := 0

		for {

			err = sendMensage(wac, contactNumber, mensagem)
			if err == nil {
				retyingCount = 0
				break
			}

			if retyingCount < 4 {
				retyingCount++
				fmt.Fprintf(os.Stderr, "error sending message to : %v | Retrying %d\n", contactNumber, retyingCount)
				continue
			} else {
				fmt.Fprintf(os.Stderr, "error sending message to : %v | Already retried %d times. Exiting program\n", err, retyingCount)
				os.Exit(1)
				break
			}

		}

	}
}

func sendMensage(wac *whatsapp.Conn, number, text string) error {

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: number + "@s.whatsapp.net",
		},
		Text: text,
	}
	_, err := wac.Send(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v\n", err)
		return err
	} else {
		fmt.Println("Message Sent -> Number : " + number)
	}
	return nil
}

func login(wac *whatsapp.Conn) error {
	//load saved session
	// session, err := readSession()
	// if err == nil {
	// 	//restore session
	// 	session, err = wac.RestoreWithSession(session)
	// 	if err != nil {
	// 		return fmt.Errorf("restoring failed: %v\n", err)
	// 	}
	// } else {
	// 	//no saved session -> regular login
	// 	qr := make(chan string)
	// 	go func() {
	// 		terminal := qrcodeTerminal.New()
	// 		terminal.Get(<-qr).Print()
	// 	}()
	// 	session, err = wac.Login(qr)
	// 	if err != nil {
	// 		return fmt.Errorf("error during login: %v\n", err)
	// 	}
	// }
	qr := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New()
		terminal.Get(<-qr).Print()
	}()
	session, err := wac.Login(qr)
	if err != nil {
		return fmt.Errorf("error during login: %v\n", err)
	}
	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}

func processCSV(rc io.Reader) (ch chan []string) {
	ch = make(chan []string, 10)
	go func() {
		r := csv.NewReader(rc)
		if _, err := r.Read(); err != nil { //read header
			log.Fatal(err)
		}
		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)

			}
			ch <- rec
		}
	}()
	return
}

func fixPhoneNumber(number string) string {
	number = strings.ReplaceAll(number, " ", "")
	number = strings.ReplaceAll(number, "-", "")
	number = strings.ReplaceAll(number, "+", "")
	number = strings.ReplaceAll(number, ".", "")
	number = strings.ReplaceAll(number, "(", "")
	number = strings.ReplaceAll(number, ")", "")
	number = strings.ReplaceAll(number, "/", "")
	number = strings.TrimPrefix(number, "0")

	size := len(number)
	result := number

	if size == 9 {
		result = "5521" + number
	}

	if size == 11 {
		result = "55" + number
	}

	return result
}
