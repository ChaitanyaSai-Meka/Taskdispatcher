package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync/atomic"

	"github.com/ChaitanyaSai-Meka/Taskdispatcher/internal/dispatcher"
	"github.com/ChaitanyaSai-Meka/Taskdispatcher/models"
)

var nextID int64 

func StartTCP(addr string, d *dispatcher.Dispatcher) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Println("TCP server listening on", addr)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println("accept error:", err)
				continue
			}
			go handleConn(conn, d)
		}
	}()

	return nil
}

func handleConn(conn net.Conn, d *dispatcher.Dispatcher) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	line := strings.TrimSpace(scanner.Text())

	fields := strings.Fields(line)
	if len(fields) < 2 {
		fmt.Fprintln(conn, "error: expected '<class> <name>', e.g. 'class1 mytask'")
		return
	}

	class, ok := parseClass(fields[0])
	if !ok {
		fmt.Fprintln(conn, "error: unknown class", fields[0])
		return
	}
	name := strings.Join(fields[1:], " ")

	id := atomic.AddInt64(&nextID, 1)
	task := models.Task{
		ID:     int(id),
		TaskName:   name,
		Class:  class,
		Status: models.StatusQueued,
	}

	d.Submit(task)

	fmt.Fprintf(conn, "task_id: %d\n", task.ID)
}

func parseClass(s string) (models.JobClass, bool) {
	switch s {
	case "class1":
		return models.Class1, true
	case "class2":
		return models.Class2, true
	case "class3":
		return models.Class3, true
	default:
		return 0, false
	}
}