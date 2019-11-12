package server

import (
	log "github.com/sirupsen/logrus"
	"syscall"
)

type Client struct {
	Fd  int
	UID int
}

func (client Client) Send(message string) {
	_, err := syscall.Write(client.Fd, []byte(message))
	if err != nil {
		log.Errorln("could not send message to client", client.Fd, "of id", client.UID)
	}
}

func (client Client) Close() {
	_ = syscall.Close(client.Fd)
}
