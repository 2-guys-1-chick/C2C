package utils

import (
	"io"
	"net"
	"os"

	"syscall"
)

func IsDisconnectError(err error) bool {
	if err == io.EOF {
		return true
	}

	if netErr, ok := err.(*net.OpError); ok {
		if syscallErr, ok := netErr.Err.(*os.SyscallError); ok {
			if syscallErr.Err.Error() == syscall.ECONNRESET.Error() {
				return true
			}
		}
	}

	return false
}

