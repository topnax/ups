package main

import (
	log "github.com/sirupsen/logrus"
	"syscall"
	"ups/sp/server/kris_kros_server"
	"ups/sp/server/networking/server"
	"ups/sp/server/protocol/impl"
)

func main() {
	log.SetLevel(log.InfoLevel)
	//37
	//log.SetOutput(os.Stdout)

	tcpServer, err := server.NewServer(syscall.SockaddrInet4{
		Addr: [4]byte{byte(127), byte(0), byte(0), byte(1)},
		Port: 10000,
	})

	if err != nil {
		log.Errorln(err)
		return
	}

	tcpReceiver := impl.SimpleTcpMessageReceiver{}
	krisKrosServer := kris_kros_server.NewKrisKrosServer(&tcpReceiver)
	messageReader := impl.NewSimpleMessageReader(krisKrosServer, krisKrosServer.Router.Handlers)

	tcpReceiver.SetMessageReader(&messageReader)
	tcpReceiver.SetOutput(&tcpServer)
	tcpReceiver.SetFileDescriptorRemover(&tcpServer)
	tcpServer.SetOnClientDisconnectedListener(&krisKrosServer)
	tcpServer.SetOnClientConnectedListener(&krisKrosServer)

	tcpServer.Start(&tcpReceiver)

}
