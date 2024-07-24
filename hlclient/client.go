package hlclient

import (
	"encoding/binary"
	"net"
)

type MessageResp struct {
	Code byte
	Data []byte
}

type Client struct {
	addr      string
	conn      net.Conn
	sendCount int
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) SendMessage(msgType byte, msg []byte, reply bool) (*MessageResp, error) {
	if c.sendCount >= 1 {
		return nil, ErrorSendMessageOnce
	}
	c.sendCount++

	data := make([]byte, len(msg)+7)
	bccK := len(msg) + 5
	_bccK := uint16(bccK - 4)
	data[0] = 0x20
	data[1] = msgType
	data[2] = byte(_bccK)
	data[3] = byte(_bccK >> 8)
	for i := 0; i < len(msg); i++ {
		data[i+4] = msg[i]
	}
	bcc := byte(255)
	bccStart := 1
	bccEnd := bccK - 1
	_bcc := byte(0)
	for bccEnd > 0 {
		_bcc ^= data[bccStart]
		bccStart++
		bccEnd--
	}
	bcc = ^_bcc
	data[bccK] = bcc

	_, err := c.conn.Write(data)
	if err != nil {
		return nil, err
	}

	if !reply {
		return nil, nil
	} else {
		return c.RecvMessage()
	}
}

func (c *Client) RecvMessage() (*MessageResp, error) {
	header := make([]byte, 8)
	_, err := c.conn.Read(header)
	if err != nil {
		return nil, err
	}
	if header[0] != 0x20 {
		return nil, ErrorRecvMessageInvalid
	}

	code := header[1]
	dataLen := int32(binary.LittleEndian.Uint32(header[2:6]))
	if dataLen <= 0 {
		return &MessageResp{Code: code, Data: []byte{}}, nil
	}
	respData := make([]byte, dataLen)
	_, err = c.conn.Read(respData)
	if err != nil {
		return nil, err
	}

	return &MessageResp{Code: code, Data: respData}, nil
}

func (c *Client) SendData(data []byte) (int, error) {
	return c.conn.Write(data)
}

func (c *Client) RecvData(data []byte) (n int, err error) {
	return c.conn.Read(data)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
