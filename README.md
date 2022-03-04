感谢[@zh-five](https://github.com/zh-five)，在这里借鉴了他的项目代码
 
```go
package main

import (
	"flag"
	ProcessManager "github.com/Apollo-Pro/process_manager"
	"log"
	"time"
)

func main() {
	d := flag.Bool("d", false, "是否后台守护进程方式运行")
	action := flag.String("a", "", "操作 start|stop|restart")
	flag.Parse()
	//后台运行
	processManager := ProcessManager.New(*d, "test", "./")
	processManager.Daemon = true //守护进程
	processManager.MaxCount = 2

	switch *action {
	case "start":
		//启动进程
		processManager.Start(func(pid int) {
			server(pid)
		})
		break
	case "restart":
		//重启进程
		processManager.Restart(func() {
		}, func(pid int) {
			log.Printf("重启成功:%d", pid)
			server(pid)
		})
		break
	case "stop":
		//关闭进程
		processManager.Stop(func() {
			log.Printf("关闭成功")
		})
		break
	case "status":
		//关闭进程
		running, pid :=processManager.MainProcessIsRunning()
		log.Printf("进程状态:%t,pid:%d", running,pid)
		break
	}
	//操作结束
	processManager.Check(func(pid int) {
		log.Printf("启动成功:%d", pid)
	})
}

func server(pid int) {
	log.Println(pid, "start...")
	time.Sleep(time.Second * 10)
	log.Println(pid, "end")
}

```

