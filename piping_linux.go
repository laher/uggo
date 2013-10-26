// +build linux

package uggo

import "syscall"

const ioctlReadTermios = syscall.TCGETS
