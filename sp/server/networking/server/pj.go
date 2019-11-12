package server

//
//import "syscall"
//
//func (s *Server) Start(md MessageDetector, messageChan chan RawClientMessage) {
//	log.Info("Accepting connections")
//
//	buffs := make(map[UID][]byte)
//	buff := make([]byte, BuffLen)
//	rfds := &syscall.FdSet{}
//
//	for {
//		FD_ZERO(rfds)
//		FD_SET(rfds, s.TCP.FD)
//		maxFd := s.TCP.FD
//
//		for _, client := range s.Clients {
//			FD_SET(rfds, client.TCP.FD)
//
//			if client.TCP.FD > maxFd {
//				maxFd = client.TCP.FD
//			}
//		}
//
//		activeFd, err := syscall.Select(maxFd+1, rfds, nil, nil, nil)
//
//		if err != nil {
//			log.Errorln("Select error:", err)
//			continue
//		}
//
//		if activeFd < 0 {
//			log.Errorln("Negative activeFd")
//			continue
//		}
//
//		if FD_ISSET(rfds, s.TCP.FD) {
//			clientFd, sockaddr, err := syscall.Accept(s.TCP.FD)
//
//			if err != nil {
//				log.Errorln("Accept error:", err)
//				continue
//			}
//
//			client := newClient(clientFd, sockaddr)
//			s.Clients[client.TCP.FD] = &client
//			buffs[client.UID] = make([]byte, 0)
//
//			log.Infof("New FromClient[FD %v]: %v:%v",
//				client.TCP.FD,
//				client.TCP.Sockaddr.(*syscall.SockaddrInet4).Addr,
//				client.TCP.Sockaddr.(*syscall.SockaddrInet4).Port,
//			)
//
//			client.Send("Aloha!\n")
//		} else {
//			for _, client := range s.Clients {
//				if FD_ISSET(rfds, client.TCP.FD) {
//					n, err := syscall.Read(client.TCP.FD, buff)
//
//					if err != nil {
//						log.Errorln("Read error:", err)
//						break
//					}
//
//					if n == 0 {
//						log.Infof("FromClient disconnected [FD %v]: %v:%v",
//							client.TCP.FD,
//							client.TCP.Sockaddr.(*syscall.SockaddrInet4).Addr,
//							client.TCP.Sockaddr.(*syscall.SockaddrInet4).Port,
//						)
//
//						delete(s.Clients, client.TCP.FD)
//					} else {
//						// FromClient wanna talk
//						buffs[client.UID] = append(buffs[client.UID], buff[:n]...)
//
//						result := md(buffs[client.UID])
//
//						if result.ClearBuff {
//							buffs[client.UID] = make([]byte, 0)
//						}
//						if result.Err == nil {
//							messageChan <- RawClientMessage{
//								Data:       result.Data,
//								FromClient: client,
//							}
//						}
//					}
//				}
//
//				if len(buffs[client.UID]) > MaxClientMessageLen {
//					// Kick disobedient FromClient
//
//					delete(s.Clients, client.TCP.FD)
//					client.TCP.Close()
//				}
//			}
//		}
//	}
//}
