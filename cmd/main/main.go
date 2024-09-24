package main

import (
	"bufio"
	"bytes"
	"cmd-gram-cli/models"
	"cmd-gram-cli/view"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

var ip = flag.String("ip", "0.0.0.0:8080", "Input an ip to connect")

var client = &http.Client{}
var u = &models.User{}

var r = bufio.NewReader(os.Stdin)
var scanner = bufio.NewScanner(r)

func main() {
	flag.Parse()

	scanner.Split(bufio.ScanLines)
	for {

		var url string
		var meth string
		scanner.Scan()
		line := scanner.Text()

		cmdS := strings.Split(line, " ")
		fmt.Println(cmdS)
		switch cmdS[0] {
		case "/login":
			if u.Email != "" {
				fmt.Println("you already in your account")
				continue
			}
			url = fmt.Sprintf("http://%s/api/login", *ip)
			meth = http.MethodPost

			jsonData, err := login()
			if err != nil {
				fmt.Println(err)
				continue
			}

			resp, err := sendReq(jsonData, meth, url)
			if err != nil {
				fmt.Println(err)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Can`t read body due to error: ", err)
				continue
			}
			if resp.Status[0] != byte('2') {
				get := models.Error{}
				err = json.Unmarshal(body, &get)
				if err != nil {
					fmt.Println("cannot unmarshal due to error: ", err)
					continue
				}
				fmt.Println(get.Text)
			} else {
				err = json.Unmarshal(body, &u)
				if err != nil {
					fmt.Println("cannot unmarshal due to error: ", err)
					continue
				}
				fmt.Printf("Welcome %s\n", u.Email)
			}

		case "/signin":
			if u.Email != "" {
				fmt.Println("You already in your account!")
				continue
			}
			url = fmt.Sprintf("http://%s/api/sign-in", *ip)
			meth = http.MethodPost

			jsonData, err := signin()
			if err != nil {
				fmt.Println(err)
				continue
			}

			resp, err := sendReq(jsonData, meth, url)
			if err != nil {
				fmt.Println(err)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Can`t read body due to error: ", err)
				continue
			}
			if resp.Status[0] != byte('2') {
				get := models.Error{}
				err = json.Unmarshal(body, &get)
				if err != nil {
					fmt.Println("cannot unmarshal due to error: ", err)
					continue
				}
				fmt.Println(get.Text)
			} else {
				err = json.Unmarshal(body, &u)
				if err != nil {
					fmt.Println("cannot unmarshal due to error: ", err)
					continue
				}
				fmt.Printf("Welcome %s\n", u.Email)
			}

		case "/new-chat":
			if u.Email == "" {
				fmt.Println("Please log into your account")
				continue
			}

			url = fmt.Sprintf("http://%s/api/new-chat", *ip)
			meth = http.MethodPost

			jsonData, err := newChat()
			if err != nil {
				fmt.Println(err)
				continue
			}
			resp, err := sendReq(jsonData, meth, url)
			if err != nil {
				fmt.Println(err)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Can`t read body due to error: ", err)
				continue
			}

			if resp.Status[0] != byte('2') {
				get := models.Error{}
				err = json.Unmarshal(body, &get)
				if err != nil {
					fmt.Println("cannot unmarshal due to error: ", err)
					continue
				}
				fmt.Println(get.Text)
			} else {
				chat := &models.Chat{}
				err = json.Unmarshal(body, &chat)
				if err != nil {
					fmt.Println("cannot unmarshal due to error: ", err)
					continue
				}
				fmt.Printf("Chat %s created\n", chat.Name)
			}
		case "/all-chats":
			if u.Email == "" {
				fmt.Println("Please log into your account")
				continue
			}

			url = fmt.Sprintf("http://%s/api/%d/chats", *ip, u.ID)
			meth = http.MethodGet

			resp, err := sendReq([]byte{}, meth, url)
			if err != nil {
				fmt.Println("cannot unmarshal due to error: ", err)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Can`t read body due to error: ", err)
				continue
			}

			if resp.Status[0] != byte('2') {
				get := models.Error{}
				err = json.Unmarshal(body, &get)
				if err != nil {
					fmt.Println("cannot unmarshal due to error: ", err)
					continue
				}
				fmt.Println(get.Text)
			} else {
				chats := map[string][]models.Chat{}
				err = json.Unmarshal(body, &chats)
				if err != nil {
					fmt.Println("cannot unmarshal due to error: ", err)
					continue
				}
				fmt.Println(chats)
				for _, v := range chats["chat"] {
					fmt.Println(v.ID, v.Name)
				}
			}
		case "/open-chat":
			if u.Email == "" {
				fmt.Println("Please log into your account")
			}

			n, err := strconv.Atoi(cmdS[1])
			if err != nil {
				fmt.Println("Invalid argument")
				continue
			}

			origin := "http://localhost/"
			url = fmt.Sprintf("ws://%s/api/%d/chats/%d", *ip, u.ID, n)
			conn, err := websocket.Dial(url, "tcp", origin)
			if err != nil {
				fmt.Println(err)
				continue
			}

			msgs := []*models.MessageDTO{}
			// _, err = conn.Read(body)
			// if err != nil {
			// 	fmt.Println("Can`t read body due to error: ", err)
			// 	continue
			// }
			// var msgs map[string][]models.MessageDTO
			// err = websocket.JSON.Receive(conn, &msgs)
			err = websocket.JSON.Receive(conn, &msgs)
			fmt.Println(msgs[0].Body, err)
			if err != nil {
				fmt.Println("cannot unmarshal due to error: ", err)
				err = startMessaging(n, conn)
				if err != nil {
					fmt.Println(err)
				}
				continue
			}
			for _, v := range msgs {
				view.Messages(v, u)
			}
			err = startMessaging(n, conn)
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println(cmdS[0], "/open-chat")
			fmt.Println("Undefined command")
		}

	}
}

func sendReq(jsonData []byte, meth string, url string) (*http.Response, error) {
	if jsonData == nil && meth != http.MethodGet {
		return nil, fmt.Errorf("data is nil")
	}
	req, err := http.NewRequest(meth, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("somthing went wrong, check your connection %v", err)

	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can`t do req due to error: %v", err)
	}

	return resp, nil
}

func login() ([]byte, error) {
	gotU := &models.UserDTO{}

	fmt.Println("email: ")
	fmt.Scan(&gotU.Email)
	fmt.Println("password: ")
	fmt.Scan(&gotU.EncryptedPassword)

	jsonData, err := json.Marshal(gotU)
	if err != nil {
		return nil, err
	}
	return jsonData, err
}

func signin() ([]byte, error) {
	gotU := &models.UserDTO{}

	fmt.Println("email: ")
	fmt.Scan(&gotU.Email)
	fmt.Println("password: ")
	fmt.Scan(&gotU.EncryptedPassword)

	jsonData, err := json.Marshal(gotU)
	if err != nil {
		return nil, err
	}

	return jsonData, err
}

func newChat() ([]byte, error) {
	var u1 = &models.UserDTO{Email: u.Email}
	var u2 = &models.UserDTO{}

	fmt.Println("Intput the user")
	fmt.Scan(&u2.Email)

	jsonData, err := json.Marshal([]*models.UserDTO{u1, u2})
	if err != nil {
		return nil, err
	}

	return jsonData, err
}

func startMessaging(cid int, conn *websocket.Conn) error {
	go func(c *websocket.Conn) {
		msg := &models.MessageDTO{}
		for {
			if !c.IsServerConn() {
				return
			}
			websocket.JSON.Receive(conn, msg)

			view.Messages(msg, u)
		}
	}(conn)

	for {
		msg := &models.MessageDTO{UserID: u.ID, ChatID: cid, Time: time.Now()}

		scanner.Scan()

		b := scanner.Text()
		if b == "" {
			continue
		}
		msg.Body = b

		jsonData, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		_, err = conn.Write(jsonData)
		if err != nil {
			return err
		}

		if b == "/exit" {
			break
		}
	}
	return nil
}
