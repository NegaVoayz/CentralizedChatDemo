package osy_msg

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

type OsyMessage struct {
	Timestamp    time.Time
	Sender_id    uint32
	Receiver_id  uint32
	Message_type uint32
	Message_size uint32
	Message_body []byte
}

func NameHash(str string) uint32 {
	var ans uint32 = 0
	for _, elem := range str {
		ans = ans*127 + uint32(elem)
	}
	return ans
}

func SetMessage(
	sender_id uint32,
	receiver_id uint32,
	message_type uint32,
	message_size uint32,
	message_body []byte) OsyMessage {
	var ans OsyMessage
	ans.Timestamp = time.Now()
	ans.Sender_id = sender_id
	ans.Receiver_id = receiver_id
	ans.Message_type = message_type
	ans.Message_size = message_size
	ans.Message_body = message_body
	return ans
}

func SendMessage(conn_ptr *net.Conn, msg_ptr *OsyMessage) bool {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, (*msg_ptr).Timestamp.Unix()) // 64
	binary.Write(buf, binary.LittleEndian, (*msg_ptr).Sender_id)        // 32
	binary.Write(buf, binary.LittleEndian, (*msg_ptr).Receiver_id)      // 32
	binary.Write(buf, binary.LittleEndian, (*msg_ptr).Message_type)     // 32
	binary.Write(buf, binary.LittleEndian, (*msg_ptr).Message_size)     // 32
	binary.Write(buf, binary.LittleEndian, []byte((*msg_ptr).Message_body))

	(*conn_ptr).Write(buf.Bytes())
	return true
}

func RecvFull(conn_ptr *net.Conn, length int32) (buf_return []byte, done bool) {
	done = false
	buf_return = make([]byte, 0)

	buf_temp := make([]byte, length)
	n, err := (*conn_ptr).Read(buf_temp)
	if n == 0 || err != nil {
		return
	}

	for n != int(length) {
		buf_return = append(buf_return, buf_temp[:n]...)
		length -= int32(n)
		buf_temp = make([]byte, length)
		n, err = (*conn_ptr).Read(buf_temp)
		if n == 0 || err != nil {
			return
		}
	}
	buf_return = append(buf_return, buf_temp[:n]...)
	done = true
	return
}

func RecvMessage(conn_ptr *net.Conn, msg_ptr *OsyMessage) bool {
	buf_data, done := RecvFull(conn_ptr, 24)
	if !done {
		return false
	}
	buf := bytes.NewReader(buf_data)

	var temp int64
	binary.Read(buf, binary.LittleEndian, &temp)
	(*msg_ptr).Timestamp = time.Unix(temp, 0)
	binary.Read(buf, binary.LittleEndian, &((*msg_ptr).Sender_id))
	binary.Read(buf, binary.LittleEndian, &((*msg_ptr).Receiver_id))
	binary.Read(buf, binary.LittleEndian, &((*msg_ptr).Message_type))
	binary.Read(buf, binary.LittleEndian, &((*msg_ptr).Message_size))

	if (*msg_ptr).Message_size == 0 {
		return true
	}

	(*msg_ptr).Message_body, done = RecvFull(conn_ptr, int32((*msg_ptr).Message_size))
	return done
}
