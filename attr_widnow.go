// +build windows

package process_manager

import (
	"bytes"
	"fmt"
	"strings"
	//"github.com/axgle/mahonia"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: true,
	}
}

func Kill(pid int) error {
	if pid == 0 {
		return nil
	}
	cmd := exec.Command("taskkill.exe", "/F", "/PID", fmt.Sprintf("%d", pid))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	//enc := mahonia.NewEncoder("GBK")
	//log.Info(enc.ConvertString(out.String()))
	return err
}

func PidExist(pid int) bool {
	if pid == 0 {
		return false
	}
	param := fmt.Sprintf(`PID eq %d`, pid)
	cmd := exec.Command("tasklist.exe", "/NH", "/fi", param)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return strings.Contains(string(out), fmt.Sprintf(`%d`, pid))
}

func SetProcessName(name string) error {
	return nil
}
