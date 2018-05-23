package main

import "fmt"
import "os/exec"
import "os"
import "os/signal"
import "syscall"
import "time"

var cmd *exec.Cmd

func signalWatcher() {
	signalCount := 0
	signalChan := make(chan os.Signal, 10)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for signal := range signalChan {
		fmt.Println("Retransmitting signal")
		if signalCount == 0 {
			fmt.Println("Delaying first signal by 10s")
			time.Sleep(10 * time.Second)
			signalCount = 1
		}
		cmd.Process.Signal(signal)
	}
}
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [cmd] [args]\n", os.Args[1])
		os.Exit(1)
	}
	go signalWatcher()
	cmd = exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Setpgid = true
	err := cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		} else {
			os.Exit(1)
		}
	}
}
