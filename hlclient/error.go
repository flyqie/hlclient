package hlclient

import "errors"

var (
	ErrorSendMessageOnce    = errors.New("sendMessage called more than once")
	ErrorRecvMessageInvalid = errors.New("recvMessage invalid")
)
