







package unix

import "unsafe"


func cmsgAlignOf(salen int) int {
	salign := sizeofPtr
	
	
	
	if darwin64Bit || dragonfly64Bit || solaris64Bit {
		salign = 4
	}
	return (salen + salign - 1) & ^(salign - 1)
}



func CmsgLen(datalen int) int {
	return cmsgAlignOf(SizeofCmsghdr) + datalen
}



func CmsgSpace(datalen int) int {
	return cmsgAlignOf(SizeofCmsghdr) + cmsgAlignOf(datalen)
}

func cmsgData(h *Cmsghdr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(h)) + uintptr(cmsgAlignOf(SizeofCmsghdr)))
}


type SocketControlMessage struct {
	Header Cmsghdr
	Data   []byte
}



func ParseSocketControlMessage(b []byte) ([]SocketControlMessage, error) {
	var msgs []SocketControlMessage
	i := 0
	for i+CmsgLen(0) <= len(b) {
		h, dbuf, err := socketControlMessageHeaderAndData(b[i:])
		if err != nil {
			return nil, err
		}
		m := SocketControlMessage{Header: *h, Data: dbuf}
		msgs = append(msgs, m)
		i += cmsgAlignOf(int(h.Len))
	}
	return msgs, nil
}

func socketControlMessageHeaderAndData(b []byte) (*Cmsghdr, []byte, error) {
	h := (*Cmsghdr)(unsafe.Pointer(&b[0]))
	if h.Len < SizeofCmsghdr || uint64(h.Len) > uint64(len(b)) {
		return nil, nil, EINVAL
	}
	return h, b[cmsgAlignOf(SizeofCmsghdr):h.Len], nil
}



func UnixRights(fds ...int) []byte {
	datalen := len(fds) * 4
	b := make([]byte, CmsgSpace(datalen))
	h := (*Cmsghdr)(unsafe.Pointer(&b[0]))
	h.Level = SOL_SOCKET
	h.Type = SCM_RIGHTS
	h.SetLen(CmsgLen(datalen))
	data := cmsgData(h)
	for _, fd := range fds {
		*(*int32)(data) = int32(fd)
		data = unsafe.Pointer(uintptr(data) + 4)
	}
	return b
}



func ParseUnixRights(m *SocketControlMessage) ([]int, error) {
	if m.Header.Level != SOL_SOCKET {
		return nil, EINVAL
	}
	if m.Header.Type != SCM_RIGHTS {
		return nil, EINVAL
	}
	fds := make([]int, len(m.Data)>>2)
	for i, j := 0, 0; i < len(m.Data); i += 4 {
		fds[j] = int(*(*int32)(unsafe.Pointer(&m.Data[i])))
		j++
	}
	return fds, nil
}
