package main

import (
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"strconv"
	"syscall"
	"ups/sp/server/kris_kros_server"
	"ups/sp/server/networking/server"
	"ups/sp/server/protocol/impl"
)

const (
	defaultClientLimit = 20
	defaultPort        = 10000
	defaultLobbyLimit  = 10
)

func main() {
	host := [4]byte{byte(127), byte(0), byte(0), byte(1)}
	port := defaultPort
	clientLimit := defaultClientLimit
	lobbyLimit := defaultLobbyLimit

	if len(os.Args) >= 3 {
		argHost := net.ParseIP(os.Args[1])
		argPort, err := strconv.Atoi(os.Args[2])

		if len(argHost) == 16 {
			if err == nil {
				host = [4]byte{argHost[12], argHost[13], argHost[14], argHost[15]}
				port = argPort
			} else {
				log.Warnf("Ignoring server config, because of invalid port %s", os.Args[2])
			}
		} else {
			log.Warnf("Ignoring server config, because of invalid host %s", os.Args[1])
		}

		if len(os.Args) >= 4 {
			argClientLimit, err := strconv.Atoi(os.Args[3])
			if err != nil {
				log.Warnf("Ignoring client limit, because input '%s' is invalid", os.Args[3])
			} else {
				log.Infof("Read client limit argument of %d", argClientLimit)
				clientLimit = argClientLimit
			}
		}

		if len(os.Args) >= 5 {
			argLobbyLimit, err := strconv.Atoi(os.Args[4])
			if err != nil {
				log.Warnf("Ignoring client limit, because input '%s' is invalid", os.Args[4])
			} else {
				log.Infof("Read lobby limit argument of %d", argLobbyLimit)
				lobbyLimit = argLobbyLimit
			}
		}
	}
	log.SetLevel(log.InfoLevel)

	log.Infof("Starting server at %d.%d.%d.%d:%d", host[0], host[1], host[2], host[3], port)

	tcpServer, err := server.NewServer(syscall.SockaddrInet4{
		Addr: host,
		Port: port,
	}, clientLimit)

	if err != nil {
		log.Errorln(err)
		return
	}

	tcpReceiver := impl.SimpleTcpMessageReceiver{}
	krisKrosServer := kris_kros_server.NewKrisKrosServer(&tcpReceiver, lobbyLimit)
	messageReader := impl.NewSimpleMessageReader(krisKrosServer, krisKrosServer.Router.Handlers)

	tcpReceiver.SetMessageReader(&messageReader)
	tcpReceiver.SetOutput(&tcpServer)
	tcpReceiver.SetSocketCloser(&tcpServer)
	krisKrosServer.Router.SetSocketCloser(&tcpServer)

	tcpServer.Start(&tcpReceiver)
}
