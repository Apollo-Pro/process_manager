package main

import (
	ProcessManager "github.com/Apollo-Pro/process_manager"
	"log"
	"time"
)

func main() {

	processManager := ProcessManager.New(true, "test", "./")
	processManager.Daemon = false
	processManager.MaxCount = 2

	processManager.Start(func(pid int) {
		server(pid)
	})

	processManager.Check(func(pid int) {
		log.Printf("启动成功:%d", pid)
	})
}

func server(pid int) {
	log.Println(pid, "start...")
	time.Sleep(time.Second * 10)
	log.Println(pid, "end")
}
