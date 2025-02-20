package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"osy_msg"
	"regexp"
	"strings"
)

var clt_id uint32
var names map[uint32]string
var inputReader *bufio.Reader

const server_name string = "localhost:8080"

func main() {
	inputReader = bufio.NewReader(os.Stdin)
	names = make(map[uint32]string, 0)
	conn, err := net.Dial("tcp", server_name)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Host " + server_name + " Is Connected")

	InputName(&conn)
	go Receiver(&conn)
	for {
		msg := GetText()
		osy_msg.SendMessage(&conn, &msg)
	}
}

func InputName(conn_ptr *net.Conn) {
	fmt.Println("Waiting for name:")
	var tempstr string
	fmt.Scan(&tempstr)
	matched, err := regexp.Match("^[a-zA-Z0-9_]*", []byte(tempstr))
	if !matched || err != nil {
		fmt.Println("Illegal Input")
		InputName(conn_ptr)
		return
	}
	msg := osy_msg.SetMessage(osy_msg.NameHash(tempstr), 0, 1, uint32(len(tempstr)), []byte(tempstr))
	osy_msg.SendMessage(conn_ptr, &msg)
}

func Receiver(conn_ptr *net.Conn) {
	var recv osy_msg.OsyMessage
	for osy_msg.RecvMessage(conn_ptr, &recv) {
		switch recv.Message_type {
		case 0:
			fmt.Println(names[recv.Sender_id] + ": " + string(recv.Message_body[:recv.Message_size]))
		case 1:
			clt_id = recv.Receiver_id
		case 2:
			names[recv.Sender_id] = string(recv.Message_body[:recv.Message_size])
			fmt.Println(names[recv.Sender_id] + " is online.")
		case 1023:
			fmt.Println(names[recv.Sender_id] + " is offline.")
			delete(names, recv.Sender_id)
		}
	}

}

func GetText() osy_msg.OsyMessage {
	tempstr, _ := inputReader.ReadString('\n')
	tempstr = tempstr[:len(tempstr)-1]
	text := regexp.MustCompile("^to [A-Za-z0-9_]*:").FindStringSubmatch(tempstr)
	if len(text) == 0 {
		return osy_msg.SetMessage(clt_id, 0, 0, uint32(len(tempstr)), []byte(tempstr))
	}
	tempstr = strings.TrimLeft(tempstr[len(text[0]):], " \t")
	return osy_msg.SetMessage(clt_id, osy_msg.NameHash(text[0][3:(len(text[0])-1)]), 0, uint32(len(tempstr)), []byte(tempstr))
}
