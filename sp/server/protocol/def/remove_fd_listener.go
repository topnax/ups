package def

type FileDescriptorRemover interface {
	RemoveFd(fd int)
}
