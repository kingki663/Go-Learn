package codec

import (
	"Go_Learn/RpcDemo/protocol"
	"Go_Learn/RpcDemo/serializer"
	"bufio"
	"io"
	"log"
)

type ServerCodec struct {
	conn       io.ReadWriteCloser
	serializer serializer.JsonSerializer
	w          *bufio.Writer
	r          *bufio.Reader
}

func NewServerCodec(conn io.ReadWriteCloser) *ServerCodec {
	return &ServerCodec{
		conn: conn,
		w:    bufio.NewWriter(conn),
		r:    bufio.NewReader(conn),
	}
}

func (c *ServerCodec) ReadRequest() (*protocol.Request, error) {
	req := new(protocol.Request)
	reqBytes, err := recvFrame(c.r)
	if err != nil {
		log.Println("rpc server: serverCodec ReadRequest: " + err.Error())
		return nil, err
	}
	log.Println("rpc server: ReadRequest JSON:", string(reqBytes))

	err = c.serializer.Unmarshal(reqBytes, req)
	if err != nil {
		log.Println("rpc server: serverCodec ReadRequest: " + err.Error())
	}
	return req, nil
}

func (c *ServerCodec) WriteResponse(errMsg error, replies []interface{}) {
	defer func() {
		err := c.w.Flush() // 将所有的缓存数据写入底层的IO接口
		if err != nil {
			_ = c.Close() // 发生错误则关闭
		}
	}()

	resp := new(protocol.Response)
	resp.Replies = replies

	if errMsg == nil {
		resp.Err = ""
	} else {
		resp.Err = errMsg.Error()
	}
	respBytes, err := c.serializer.Marshal(&resp)
	if err != nil {
		log.Println("rpc server: serverCodec WriteResponse: " + err.Error())
		return
	}

	err = sendFrame(c.w, respBytes)
	if err != nil {
		log.Println("rpc server: serverCodec WriteResponse: " + err.Error())
		return
	}
}

func (c *ServerCodec) Close() error {
	return c.conn.Close()
}
