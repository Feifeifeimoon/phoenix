package msg

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"reflect"
)

/*
	Read Message From End-Point
	Return Interface Message and err if err occur
*/
func Read(in io.Reader) (m Message, err error) {
	return read(in, nil)
}

func ReadInto(in io.Reader, m Message) error {
	_, err := read(in, m)
	return err
}

func read(in io.Reader, msg Message) (m Message, err error) {
	// read first byte get message type
	buffer := make([]byte, 1)
	_, err = in.Read(buffer)
	if err != nil {
		return
	}
	typeByte := buffer[0]

	var length int64
	err = binary.Read(in, binary.BigEndian, &length)
	if err != nil {
		return
	}

	buffer = make([]byte, length)
	// ReadFull reads exactly len(buf) bytes from in into buffer.
	n, err := io.ReadFull(in, buffer)
	if err != nil {
		return nil, err
	}
	if int64(n) != length {
		return nil, errors.New("Error Length")
	}

	if msg == nil {
		t, ok := byte2type[typeByte]
		if !ok {
			return nil, errors.New("Error Type")
		}
		m = reflect.New(t).Interface().(Message)
	} else {
		m = msg
	}
	err = json.Unmarshal(buffer, m)
	return
}

/*
	Write Message to End-Point
	Return if err occur
*/
func Write(out io.Writer, m Message) (err error) {
	typeByte, ok := type2byte[reflect.TypeOf(m)]
	if !ok {
		return errors.New("Error Msg Type")
	}

	ret, err := json.Marshal(m)
	if err != nil {
		return
	}

	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte(typeByte)
	_ = binary.Write(buffer, binary.BigEndian, int64(len(ret)))
	buffer.Write(ret)

	_, err = out.Write(buffer.Bytes())
	return
}
