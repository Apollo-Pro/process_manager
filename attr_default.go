// +build !windows,!plan9

package process_manager

import (
	"log"
	"syscall"
	"unsafe"
)

func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}

func Kill(pid int) error {
	if pid == 0 {
		return nil
	}
	return syscall.Kill(pid, syscall.SIGKILL)
}

func PidExist(pid int) bool {
	if pid == 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	if err == nil {
		return true
	}
	log.Printf(err.Error())
	return false
}

func SetProcessName(name string) error {
	bytes := append([]byte(name), 0)
	ptr := unsafe.Pointer(&bytes[0])
	if _, _, errno := syscall.RawSyscall6(syscall.SYS_PRCTL, syscall.PR_SET_NAME, uintptr(ptr), 0, 0, 0, 0); errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}
