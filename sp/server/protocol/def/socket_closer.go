package def

// defines a socket closer that closes the given socket
type SocketCloser interface {
	CloseFd(socket int)
}
