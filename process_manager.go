package process_manager

import (
	"fmt"
	"os"

	"time"
	//"syscall"
	"log"
)

//版本号
const VERSION = "1.0.1"

const (
	PTYPE_START  = "start_process"  //启动进程
	PTYPE_DAEMON = "daemon_process" //守护进程
	PTYPE_MAIN   = "main_process"   //主进程
)

type Manager struct {
	Daemon        bool   //是否守护进程
	Background    bool   //是否后台运行
	MaxCount      int    //循环重启最大次数, 若为0则无限重启
	MainPidFile   string //主进程PID文件
	DaemonPidFile string //守护进程PID文件
	isStartAction bool   //是否为启动操作
	ServerName    string //进程名称
	runtimePath   string //工作目录
	logger        *log.Logger
}

func New(background bool, serverName string, runtimePath string) *Manager {
	manager := &Manager{
		Background:  background,
		runtimePath: runtimePath,
		ServerName:  serverName}
	manager.init()
	return manager
}

//初始化
func (c *Manager) init() {
	if err := os.MkdirAll(c.runtimePath, 0666); err != nil {
		log.Fatalf(err.Error())
	}
	c.logger = c.Logger(c.runtimePath)
	c.MainPidFile = fmt.Sprintf("%s/main.pid", c.runtimePath)
	c.DaemonPidFile = fmt.Sprintf("%s/daemon.pid", c.runtimePath)
}

//启动进程
func (c *Manager) Start(callback func(pid int)) {
	c.isStartAction = true

	isRunning, runningPid := c.DaemonProcessIsRunning()
	if c.isStartProcess() && isRunning == true {
		c.logger.Fatalf("守护进程运行中，请先执行Stop, Pid:%d", runningPid)
	}

	isMRunning, runningMPid := c.MainProcessIsRunning()
	if c.isStartProcess() && isMRunning == true {
		c.logger.Fatalf("主进程运行中，请先执行Stop, Pid:%d", runningMPid)
	}

	if c.Background && c.isStartProcess() { //先开一个守护进程
		c.background(PTYPE_DAEMON)
	} else if c.Background && c.isDaemonProcess() { //在守护进程中开主进程
		c.saveDaemonPid(os.Getpid())
		c.daemon()
	} else {
		c.saveMainPid(os.Getpid())
		c.listenerExit(func() {
			c.delPidFile(c.MainPidFile)
		})
		callback(os.Getpid())
	}
}

//判断主进程是否运行中
func (c *Manager) MainProcessIsRunning() (bool, int) {
	return IsRunning(c.getPid(c.MainPidFile))
}

//判断守护进程是否运行中
func (c *Manager) DaemonProcessIsRunning() (bool, int) {
	return IsRunning(c.getPid(c.DaemonPidFile))
}

//判断主进程是否运行中
func IsRunning(pid int) (bool, int) {
	if PidExist(pid) {
		return true, pid
	}
	return false, 0
}

//重启进程
func (c *Manager) Restart(stopCallback func(), startCallback func(pid int)) {
	if c.isStartProcess() {
		c.Stop(stopCallback)
		time.Sleep(1 * time.Second)
	}
	c.Start(startCallback)
}

//关闭进程
func (c *Manager) Stop(callback func()) error {
	//关闭守护进程
	if dRunning, killDaemonErr := c.KillProcess(c.DaemonPidFile); killDaemonErr != nil && dRunning {
		return killDaemonErr
	}

	//关闭主进程
	if mRunning, killMainErr := c.KillProcess(c.MainPidFile); killMainErr != nil && mRunning {
		return killMainErr
	}

	callback()

	return nil
}

func (c *Manager) KillProcess(pidFile string) (bool, error) {
	isRunning, pid := IsRunning(c.getPid(pidFile))
	if isRunning {
		if killErr := Kill(pid); killErr != nil {
			c.logger.Printf("PID:%d,关闭进程错误! %s", pid, killErr.Error())
			return true, killErr
		}
		c.logger.Printf("%s,Pid:%d, 退出成功!", c.ServerName, pid)
		c.delPidFile(pidFile)
		return true, nil
	}
	return false, nil
}

//启动进程执行完毕，检查主进程是否正常运行
func (c *Manager) Check(callback func(pid int)) {
	if c.isStartProcess() && c.isStartAction {
		for i := 30; i > 0; i-- {
			time.Sleep(time.Second * 1)
			if isRunning, pid := c.MainProcessIsRunning(); isRunning {
				callback(pid)
				os.Exit(0)
			}
			c.logger.Printf(" %2d...", i)
		}
		c.logger.Fatalf("启动失败")
	}
}
