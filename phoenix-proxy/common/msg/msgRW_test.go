package msg

import (
	"bytes"
	"testing"
)

func TestMsgReaderWriter_Read(t *testing.T) {
	var network bytes.Buffer

	request := RegisterProxy{
		ProxyName:  "proxy-test",
		ProxyType:  1,
		RemotePort: "6000",
	}
	if err := Write(&network, request); err != nil {
		t.Log(err.Error())
		return
	}

	response, err := Read(&network)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Logf(" %#v \n", response)
	switch m := response.(type) {
	case *Login:
		t.Logf("Login %#v \n", m)
	case *RegisterProxy:
		t.Logf("RegisterProxy %#v \n", m)
	default:
		t.Logf(" %#v \n", m)
	}

}

func TestMsgReaderWriter_HeartBeat(t *testing.T) {
	var network bytes.Buffer

	request := HeartBeat{}
	if err := Write(&network, request); err != nil {
		t.Log(err.Error())
		return
	}

	response, err := Read(&network)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Logf(" %#v \n", response)
	switch m := response.(type) {
	case *Login:
		t.Logf("Login %#v \n", m)
	case *RegisterProxy:
		t.Logf("RegisterProxy %#v \n", m)
	case *HeartBeat:
		t.Logf("HeartBeat %#v \n", m)

	default:
		t.Logf(" %#v \n", m)
	}
}

func TestMsgReaderWriter_ReadInto(t *testing.T) {
	var network bytes.Buffer
	login := Login{
		HostName:  "k8s",
		Os:        "linux",
		Arch:      "amd64",
		ControlID: 0x123,
		ProxyName: "ssh",
	}
	if err := Write(&network, login); err != nil {
		t.Log(err.Error())
		return
	}
	var m Login
	if err := ReadInto(&network, &m); err != nil {
		t.Log(err.Error())
		return
	}
	t.Logf("%#v \n", m)

}
