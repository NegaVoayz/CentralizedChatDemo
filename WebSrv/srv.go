package main

import (
	"fmt"
	"net"
	"osy_msg"
)

var names map[uint32]string
var message map[uint32]chan osy_msg.OsyMessage

func main() {
	// init mdb
	names = make(map[uint32]string, 0)
	message = make(map[uint32]chan osy_msg.OsyMessage, 0)
	names[0] = "Everyone"

	// init socket
	listener, _ := net.Listen("tcp", "localhost:8080")
	fmt.Println("Server started.")

	// prepare for safe closure
	defer listener.Close()

	// main loop
	for {
		conn, _ := listener.Accept()
		go SrvConnProcAnsw(&conn)
	}
}

func SrvConnProcAnsw(conn_ptr *net.Conn) {
	var clt_id uint32 = 0
	defer (*conn_ptr).Close()
	var recv osy_msg.OsyMessage
	for osy_msg.RecvMessage(conn_ptr, &recv) {
		if recv.Sender_id == 0 {
			continue
		}
		fmt.Println(
			"Received:\n    from",
			names[recv.Sender_id],
			" to ",
			names[recv.Receiver_id],
			":\n        ",
			string(recv.Message_body[:recv.Message_size]))
		switch recv.Message_type {
		case 0:

			if recv.Receiver_id == 0 {
				for id, channel := range message {
					if id == recv.Sender_id {
						continue
					}
					channel <- recv
				}
			} else {
				message[recv.Receiver_id] <- recv
			}
		case 1:
			if clt_id != 0 {
				close(message[clt_id])
				delete(message, clt_id)
				delete(names, clt_id)
			}
			name := string(recv.Message_body[:recv.Message_size])
			clt_id = recv.Sender_id
			names[clt_id] = name
			message[clt_id] = make(chan osy_msg.OsyMessage, 32)
			message[clt_id] <- osy_msg.SetMessage(0, clt_id, 1, 0, make([]byte, 0))
			recv.Message_type = 2
			for id, channel := range message {
				// could be removed, exists only for the standard
				recv.Receiver_id = id

				channel <- recv
			}
			go SrvConnProcPush(conn_ptr, clt_id)
		}
	}
	msg := osy_msg.SetMessage(clt_id, 0, 1023, 0, make([]byte, 0))
	for id, channel := range message {
		if id == clt_id {
			continue
		}

		// could be removed, exists only for the standard
		msg.Receiver_id = id

		channel <- msg
	}
	fmt.Println("Connection to "+names[clt_id]+"(id=", clt_id, ") was closed")
	close(message[clt_id])
	delete(message, clt_id)
	delete(names, clt_id)
}

func SrvConnProcPush(conn_ptr *net.Conn, clt_id uint32) {
	// sync loop
	for id, name := range names {
		if id == 0 || id == clt_id {
			continue
		}
		send_buf := osy_msg.SetMessage(id, clt_id, 2, uint32(len(name)), []byte(name))
		osy_msg.SendMessage(conn_ptr, &send_buf)
	}

	// sender loop
	for send_buf := range message[clt_id] {
		osy_msg.SendMessage(conn_ptr, &send_buf)
		fmt.Println("sent ", send_buf.Message_type, " to ", clt_id)
	}

	fmt.Println("Pusher to ", clt_id, " was closed")
}
