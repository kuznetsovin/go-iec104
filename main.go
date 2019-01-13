package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/signal"
)

func handlerIecFrame(apdu []byte) {
	log.Println(apdu)
	sAcpi := apdu[2:6]

	switch {
	case sAcpi[0]&1 == 0:
		// если 1 бит = 1 у первого котрольного байта I frame
		// передаваемый порядковый номер N(S), смещение для того чтобы исключить контрольный бит описанный выше
		frameSendSequenceNumber := binary.LittleEndian.Uint16(sAcpi[:2]) >> 1

		// принимаемый порядковый номер N(R)
		frameRecvSequenceNumber := binary.LittleEndian.Uint16(sAcpi[2:]) >> 1

		log.Printf("Received frame I(%v, %v)", frameSendSequenceNumber, frameRecvSequenceNumber)

		asdu := apdu[6:]
		switch asdu[0] {
		case M_ME_NB_1:
			ioa, val := Parse_M_ME_NB(asdu)
			log.Println("IOA: ", ioa, " Value: ", val)
		default:
			log.Println("Uknown format")
		}

	case sAcpi[0]&3 == 1:
		// если 1 бит = 1, 2 бит = 0 у первого котрольного байта S frame

		// принимаемый порядковый номер N(R)
		frameRecvSequenceNumber := binary.LittleEndian.Uint16(sAcpi[3:4]) >> 1
		log.Printf("Received frame (%v)", frameRecvSequenceNumber)

	case sAcpi[0]&3 == 3:
		// если 1 бит = 1, 2 бит = 1 у первого котрольного байта S frame
		log.Println("Received U-FRAME")

		// проверка установки соединения
		if bytes.Equal(sAcpi, STARTDT_CON) {
			log.Println("Connection established")
		}
	}
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)

	iecClient := Iec104Client{Host: "localhost", Port: "2404"}

	iecClient.Connect()

	iecClient.addReceiverHandler(handlerIecFrame)

	iecClient.SendCommand(STARTDT_ACT)

	select {
	case <-signals:
		iecClient.Close()
	}
}
