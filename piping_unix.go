// +build !windows

package uggo

import(
	"syscall"
	"unsafe"
)

func IsPipingStdin() bool {
	return !isTerminal(0)
}

//IsTerminal is only implemented in non-Windows so far, so only used by IsPipingStdin
func isTerminal(fd int) bool {
        var termios syscall.Termios
        _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
        return err == 0
}
