package process_manager

import (
	"fmt"
	"os"
	"os/exec"
)

//启动程序
func startProc(args, env []string) (*exec.Cmd, error) {

	cmd := &exec.Cmd{
		Path:        args[0],
		Args:        args,
		Env:         env,
		SysProcAttr: NewSysProcAttr(),
	}

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (c *Manager) isDaemonProcess() bool {
	return c.getProcessType() == PTYPE_DAEMON
}

func (c *Manager) isStartProcess() bool {
	return c.getProcessType() == PTYPE_START
}

func (c *Manager) isMainProcess() bool {
	return c.getProcessType() == PTYPE_MAIN
}

func (c *Manager) getProcessType() string {
	pType := os.Getenv("PROCESS_TYPE")
	if pType == "" {
		pType = PTYPE_START
	}
	return pType
}

// 后台运行
func (c *Manager) background(processType string) (*exec.Cmd, error) {
	//设置子进程环境变量
	env := os.Environ()
	env = append(env, fmt.Sprintf("PROCESS_TYPE=%s", processType))
	//启动子进程
	cmd, err := startProc(os.Args, env)
	if err != nil {
		c.logger.Printf("%d,启动子进程失败:%s", os.Getpid(), err.Error())
		return nil, err
	}
	return cmd, nil
}

func (c *Manager) startProcWait() (*exec.Cmd, error) {
	cmd, err := c.background(PTYPE_MAIN)
	if err != nil { //启动失败
		c.logger.Println("子进程启动失败;", "err:", err)
		return cmd, err
	}

	c.logger.Printf("主进程,开始运行...")

	//父进程: 等待子进程退出
	err = cmd.Wait()

	c.delPidFile(c.MainPidFile)

	c.logger.Printf("程序退出")

	return cmd, err
}

//守护进程启动一个子进程, 并循环监视
func (c *Manager) daemon() {
	if c.Daemon {
		c.listenerProcess()
	} else {
		c.startProcWait()
	}
	c.delPidFile(c.DaemonPidFile)
	os.Exit(0)
}
