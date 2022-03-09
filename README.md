感谢[@zh-five](https://github.com/zh-five)，在这里借鉴了他的项目代码
 
```go
package main

import (
	"flag"
	ProcessManager "github.com/Apollo-Pro/process_manager"
	"log"
    "fmt"
	"time"
)

func main() {
	d := flag.Bool("d", false, "是否后台守护进程方式运行")
	action := flag.String("a", "", "操作 start|stop|restart")
	flag.Parse()
	
	ProcessName := "test"      //进程名称(暂时没用)
	runtimePath := "./runtime" //运行目录，存放日志以及PID文件
	processManager := ProcessManager.New(*d, ProcessName, runtimePath)
	processManager.Daemon = true //守护进程
	processManager.MaxCount = 2  //最大重启次数
	
 	startFunc := func(pid int) {
		//程序入口代码
		server(pid)
	}

	switch *action {
	case "start":
		//启动进程
		processManager.Start(startFunc)
		break
	case "restart":
		//重启进程
		processManager.Restart(func() {
			log.Printf("进程已退出")
		}, startFunc)
		break
	case "stop":
		//关闭进程
		processManager.Stop(func() {
			log.Printf("关闭成功")
		})
		break
	case "status":
		//关闭进程
		running, pid := processManager.MainProcessIsRunning()
		log.Printf("进程状态:%t,pid:%d", running, pid)
		break
	}
	//操作结束
	processManager.Check(func(pid int) {
		log.Printf("启动成功:%d", pid)
	})

	//kill 掉指定文件中记录的pid
    killRet, err := processManager.KillProcess("./app.pid")
    fmt.Println(killRet, err)

    //判断pid是否存活
    isExist := ProcessManager.PidExist(88888)
    fmt.Println(isExist)

    //干掉指定pid
    killErr := ProcessManager.Kill(88888)
    fmt.Println(killErr)

}

func server(pid int) {
	log.Println(pid, "start...")
	time.Sleep(time.Second * 10)
	log.Println(pid, "end")
}

```

