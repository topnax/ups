package server

///*
//#include <sys/select.h>
//void fdclr(int fd, fd_set *set) {
//	FD_CLR(fd, set);
//}
//int fdisset(int fd, fd_set *set) {
//	return FD_ISSET(fd, set);
//}
//void fdset(int fd, fd_set *set) {
//	FD_SET(fd, set);
//}
//void fdzero(fd_set *set) {
//	FD_ZERO(set);
//}
//*/
//import (
//	"C"
//)

//import "C"
import (
	"errors"
	log "github.com/sirupsen/logrus"
	"syscall"
	"unsafe"
	"ups/sp/server/encoding"
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

func (server *Server) addClient(fd int) Client {
	client := Client{
		Fd:  fd,
		UID: server.UID,
	}
	server.Clients = append(server.Clients, client)
	server.UID++
	return client
}

func (server *Server) Start(reader encoding.MessageReader) {

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
			log.Debugln("Fd setting:", clientFd)
			FD_SET(&readfds, clientFd.Fd)
			maxFd = clientFd.Fd
		}

		log.Debugln("Readfds for fd", server.Fd, ": ", readfds)

		activeFd, err := syscall.Select(maxFd+1, &readfds, nil, nil, nil)

		log.Debugln("Selected for fd", server.Fd, ": ", readfds)

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
			client.Send("Welcome to this amazing server :)\n")
		} else {
			for i, client := range server.Clients {
				if FD_ISSET(&readfds, client.Fd) {
					n, err := syscall.Read(client.Fd, buff)

					if err != nil {
						log.Errorln("Read error:", err)
						break
					}

					if n == 0 {
						log.Infof("Client %d of id %d disconnected ",
							client.Fd,
							client.UID,
						)
						server.removeClient(i)
					} else {
						buffs[client.Fd] = append(buffs[client.Fd], buff[:n]...)
						result := string(buff[:n])
						//json.Un
						log.Debugln("Received '%s' from %d of length %d", result, client.Fd, n)
						reader.Receive(client.Fd, buff, n)
					}
				}
			}
		}
	}
}

func (server *Server) removeClient(i int) {
	server.Clients = append(server.Clients[:i], server.Clients[i+1:]...)
}

func (server *Server) acceptClient() (Client, error) {
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
