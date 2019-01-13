package main

import (
	"io"
	"log"
	"net"
)

type Iec104Client struct {
	Host string
	Port string
	conn *net.TCPConn
	//command chan []byte
}

type receiverHandler func([]byte)

func (c *Iec104Client) Connect() {
	srvAddr := c.Host + ":" + c.Port
	tcpAddr, err := net.ResolveTCPAddr("tcp", srvAddr)
	if err != nil {
		log.Fatal("ResolveTCPAddr failed:", err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("Dial failed:", err.Error())
	}

	c.conn = conn
}

func (c *Iec104Client) addReceiverHandler(handle receiverHandler) {
	go func() {
		for {
			buff := make([]byte, 1024)
			size, err := c.conn.Read(buff)
			if err != nil && err != io.EOF {
				log.Println("Read error: ", err)
			}
			if size > 0 {
				handle(buff[:size])
			}
		}
	}()
}

func (c *Iec104Client) SendCommand(cmd []byte) {
	sendCommand := []byte{0x68, byte(len(cmd))}
	sendCommand = append(sendCommand, cmd...)
	log.Printf("Send command: %x", sendCommand)
	if _, err := c.conn.Write(sendCommand); err != nil {
		log.Println("Write command error: ", err)
	}
}

func (c *Iec104Client) Close() {
	c.conn.Close()
}
