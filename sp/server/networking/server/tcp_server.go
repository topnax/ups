package server

/*
#include <sys/select.h>
void fdclr(int fd, fd_set *set) {
	FD_CLR(fd, set);
}
int fdisset(int fd, fd_set *set) {
	return FD_ISSET(fd, set);
}
void fdset(int fd, fd_set *set) {
	FD_SET(fd, set);
}
void fdzero(fd_set *set) {
	FD_ZERO(set);
}
*/
import (
	"C"
)

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"syscall"
	"unsafe"
	"ups/sp/server/protocol/def"
	"ups/sp/server/utils"
)

const (
	MaxClients = 20
	FdBits     = int(unsafe.Sizeof(0) * 8)
)

type TcpServer struct {
	UID     int
	Fd      int
	Port    int
	Clients map[int]Client
}

// creates a new server based no the given IP address
func NewServer(addr syscall.SockaddrInet4) (TcpServer, error) {
	server := TcpServer{}

	server.Clients = make(map[int]Client)

	log.Debugln("Starting server at address", addr.Addr, "at port", addr.Port)

	serverFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)

	if err != nil {
		return server, errors.New("syscall.Socket has failed")
	}

	log.Debugln("TcpServer was given FD of", serverFd)

	server.Port = addr.Port
	server.Fd = serverFd

	// reuse the addr
	_ = syscall.SetsockoptInt(server.Fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	// bind the address to the fd
	err = syscall.Bind(server.Fd, &addr)

	if err != nil {
		return server, errors.New("syscall.Bind has failed")
	}

	err = syscall.Listen(server.Fd, MaxClients)

	if err != nil {
		return server, errors.New("syscall.Listen has failed")
	}

	return server, nil
}

// sends the message to the given FD
func (server *TcpServer) Send(message string, fd int) {
	log.Debugf("TcpServer writing to socket=%d '%s'", fd, message)
	_, ok := server.Clients[fd]
	if ok {
		server.Clients[fd].Send(message)
	}
}

// adds a client
func (server *TcpServer) addClient(fd int) Client {
	client := Client{
		Fd:  fd,
		UID: server.UID,
	}
	server.Clients[fd] = client
	server.UID++
	return client
}

// closes a fd
func (server *TcpServer) CloseFd(fd int) {
	log.Warnf("Removing FD=%d", fd)
	delete(server.Clients, fd)
	_ = syscall.Shutdown(fd, syscall.SHUT_RDWR)
	_ = syscall.Close(fd)
}

// starts a server. This function is blocking as it starts the select loop
func (server *TcpServer) Start(receiver def.TcpMessageReceiver) {
	readfds := syscall.FdSet{}

	buff := make([]byte, 100)
	buffs := make(map[int][]byte)

	FD_ZERO(&readfds)

	FD_SET(&readfds, server.Fd)
	FD_SET(&readfds, 0)

	for {
		FD_ZERO(&readfds)
		FD_SET(&readfds, server.Fd)
		FD_SET(&readfds, 0)

		maxFd := server.Fd
		for _, clientFd := range server.Clients {
			FD_SET(&readfds, clientFd.Fd)
			maxFd = utils.Max(clientFd.Fd, maxFd)
		}

		activeFd, err := syscall.Select(maxFd+1, &readfds, nil, nil, nil)

		if err != nil {
			log.Errorln("Select error:", err)
			continue
		}

		if activeFd < 0 {
			log.Errorln("Negative activeFd")
		}

		if FD_ISSET(&readfds, server.Fd) {
			client, err := server.acceptClient()
			if err == nil {
				log.Infof("Client [%d] of ID %d has joined", client.Fd, client.UID)
			}
		} else {
			for _, client := range server.Clients {
				if FD_ISSET(&readfds, client.Fd) {
					n, err := syscall.Read(client.Fd, buff)

					if err != nil {
						log.Errorln("Read error:", err)
						break
					}

					if n == 0 {
						log.Debugf("Client %d of ID %d disconnected on select level",
							client.Fd,
							client.UID,
						)
						server.removeClient(client.Fd)
					} else {
						buffs[client.Fd] = append(buffs[client.Fd], buff[:n]...)
						result := string(buff[:n])
						//json.Un
						log.Debugf("Received '%s' from %d of length %d", result, client.Fd, n)
						receiver.Receive(client.Fd, buff, n)
					}
				}
			}
		}
	}
}

// removes a client from the server map
func (server *TcpServer) removeClient(fd int) {
	delete(server.Clients, fd)
}

// accepts a client
func (server *TcpServer) acceptClient() (Client, error) {
	clientSocket, _, err := syscall.Accept(server.Fd)
	if err != nil {
		log.Errorln("Failed to accept:", err)
		return Client{}, errors.New("Failed to accept")
	} else {
		log.Debugln("New fd accepted:", clientSocket)
		return server.addClient(clientSocket), nil
	}
}

func FD_SET(p *syscall.FdSet, fd int) {
	//C.fdset(C.int(fd), (*C.fd_set)(unsafe.Pointer(p)))
	p.Bits[fd/FdBits] |= int64(uint(1) << (uint(fd) % uint(FdBits)))
}

func FD_ISSET(p *syscall.FdSet, fd int) bool {
	return (p.Bits[fd/FdBits] & int64(uint(1)<<(uint(fd)%uint(FdBits)))) != 0
	//return C.fdisset(C.int(fd), (*C.fd_set)(unsafe.Pointer(p))) != 0
}

func FD_ZERO(p *syscall.FdSet) {
	for i := range p.Bits {
		p.Bits[i] = 0
	}
	//C.fdzero((*C.fd_set)(unsafe.Pointer(p)))
}
