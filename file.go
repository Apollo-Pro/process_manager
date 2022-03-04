package process_manager

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

//保存主进程id
func (c *Manager) saveMainPid(pid int) {
	c.savePid(pid, c.MainPidFile)
}

//保存守护进程id
func (c *Manager) saveDaemonPid(pid int) {
	c.savePid(pid, c.DaemonPidFile)
}

//保存pid
func (c *Manager) savePid(pid int, pidFile string) {
	file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		c.logger.Fatalf("open file error: %s", err.Error())
	}
	_, writeErr := fmt.Fprint(file, pid)
	file.Close()
	if writeErr != nil {
		Kill(pid)
		c.logger.Fatalf("Save Pid error: %s", writeErr.Error())
	}
}

/*
   判断文件或文件夹是否存在
   如果返回的错误为nil,说明文件或文件夹存在
   如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
   如果返回的错误为其它类型,则不确定是否在存在
*/
func pathExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//获取PID
func (c *Manager) getPid(pidFile string) int {
	isExists, _ := pathExists(pidFile)
	if isExists == false {
		//c.logger.Printf("Pid 文件不存在:%s", pidFile)
		return 0
	}

	data, err := ioutil.ReadFile(pidFile)
	if err != nil {
		c.logger.Fatalf("File reading error %s", err.Error())
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		c.logger.Fatalf("Pid 文件:%s,pid:%d,%s", pidFile, pid, err.Error())
	}

	return pid
}

//删除PID文件
func (c *Manager) delPidFile(pidFile string) {
	err := os.Remove(pidFile)
	if err != nil {
		c.logger.Fatalf("删除文件:%s,失败 %s", pidFile, err.Error())
	}
}

//初始化日志
func (c *Manager) Logger(logPath string) *log.Logger {
	if len(logPath) > 0 && (c.isDaemonProcess() || c.isMainProcess()) {
		name := "daemo.log"
		switch c.getProcessType() {
		case PTYPE_MAIN:
			name = "main.log"
			break
		case PTYPE_DAEMON:
			name = "daemo.log"
			break
		}

		logFile := path.Join(logPath, name)

		if _, err := os.Stat(logFile); err != nil {
			if _, err := os.Create(logFile); err != nil {
				log.Fatalf(err.Error())
			}
		}

		f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		return log.New(io.MultiWriter(f), fmt.Sprintf("[%d] ", os.Getpid()), log.Ldate|log.Ltime)
	}

	return log.New(os.Stderr, fmt.Sprintf("[%d]  ", os.Getpid()), log.Ldate|log.Ltime)
}
