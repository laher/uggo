//TODO test with freebsd,netbsd

// +build darwin freebsd netbsd

package uggo

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA
