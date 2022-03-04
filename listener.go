package process_manager

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (c *Manager) listenerProcess() {
	var t int64
	count := 0
	for {

		//daemon 信息描述
		dInfo := fmt.Sprintf("守护进程(pid:%d; 重启次数:%d/%d):",
			os.Getpid(), count, c.MaxCount)

		c.logger.Println(dInfo, "执行启动")

		if c.MaxCount > 0 && count > c.MaxCount {
			c.logger.Println(dInfo, "重启次数太多退出")
			break
		}
		count++
		t = time.Now().Unix() //启动时间戳
		_, err := c.startProcWait()

		if err != nil { //启动失败
			c.logger.Println(err.Error())
		}

		dat := time.Now().Unix() - t //子进程运行秒数

		c.logger.Printf("%s 监视到子进程退出, 共运行了%d秒: %v\n", dInfo, dat, err)
	}
}

//监听退出信号
func (c *Manager) listenerExit(callback func()) {
	channel := make(chan os.Signal)
	signal.Notify(channel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range channel {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				callback()
				os.Exit(0)
				break
			default:
				break
			}
		}
	}()
}
