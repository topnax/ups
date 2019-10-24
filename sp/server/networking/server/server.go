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
)

const (
	MAX_CLIENTS = 20
	FD_BITS     = int(unsafe.Sizeof(0) * 8)
)

type Server struct {
	UID     int
	Fd      int
	Port    int
	Clients []Client
}

func (server *Server) Init(addr syscall.SockaddrInet4) error {
	log.Debugln("Starting server at address", addr.Addr, "at port", addr.Port)

	serverFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)

	if err != nil {
		return errors.New("syscall.Socket has failed")
	}

	log.Debugln("Server was given fd of", serverFd)

	server.Port = addr.Port
	server.Fd = serverFd

	// reuse the addr
	_ = syscall.SetsockoptInt(server.Fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	// bind the address to the fd
	err = syscall.Bind(server.Fd, &addr)

	if err != nil {
		return errors.New("syscall.Bind has failed")
	}

	err = syscall.Listen(server.Fd, MAX_CLIENTS)

	if err != nil {
		return errors.New("syscall.Listen has failed")
	}

	return nil
}

func (server *Server) addClient(fd int) {
	server.Clients = append(server.Clients, Client{
		Fd:  fd,
		UID: server.UID,
	})
	server.UID++
}

func (server *Server) Start() {

	readfds := syscall.FdSet{}

	buff := make([]byte, 100)
	buffs := make(map[int][]byte)

	clientSocket := 0

	FD_ZERO(&readfds)

	FD_SET(&readfds, server.Fd)
	FD_SET(&readfds, 0)

	for {
		FD_ZERO(&readfds)
		FD_SET(&readfds, server.Fd)
		FD_SET(&readfds, 0)

		maxFd := server.Fd
		for _, clientFd := range server.Clients {
			log.Debugln("fd setting:", clientFd)
			FD_SET(&readfds, clientFd.Fd)
			maxFd = clientFd.Fd
		}

		log.Infoln("readfds: ", readfds)

		activeFd, err := syscall.Select(maxFd+1, &readfds, nil, nil, nil)

		if err != nil {
			log.Errorln("select error:", err)
			continue
		}

		if activeFd < 0 {
			log.Errorln("Negative activeFd")
		}

		if FD_ISSET(&readfds, server.Fd) {
			clientSocket, _, err = syscall.Accept(server.Fd)
			if err != nil {
				log.Infoln("new fd accepted:", clientSocket)
				server.addClient(clientSocket)
			}
		} else {
			for i, client := range server.Clients {
				if FD_ISSET(&readfds, client.Fd) {
					n, err := syscall.Read(client.Fd, buff)

					if err != nil {
						log.Errorln("Read error:", err)
						break
					}

					if n == 0 {
						log.Infof("client %d disconnected ",
							client.Fd,
						)
						server.removeClient(i)
					} else {
						buffs[client.Fd] = append(buffs[client.Fd], buff[:n]...)
						result := string(buff[:n])
						log.Infof("Received '%s' from %d", result, client.Fd)
					}
				}
			}
		}
	}
}

func (server *Server) removeClient(i int) {
	server.Clients = append(server.Clients[:i], server.Clients[i+1:]...)
}

func FD_SET(p *syscall.FdSet, fd int) {
	C.fdset(C.int(fd), (*C.fd_set)(unsafe.Pointer(p)))
	p.Bits[fd/FD_BITS] |= int64(uint(1) << (uint(fd) % uint(FD_BITS)))
}

func FD_ISSET(p *syscall.FdSet, fd int) bool {
	return (p.Bits[fd/FD_BITS] & int64(uint(1)<<(uint(fd)%uint(FD_BITS)))) != 0
	//return C.fdisset(C.int(i), (*C.fd_set)(unsafe.Pointer(p))) != 0
}

func FD_ZERO(p *syscall.FdSet) {
	for i := range p.Bits {
		p.Bits[i] = 0
	}
	//C.fdzero((*C.fd_set)(unsafe.Pointer(p)))
}
