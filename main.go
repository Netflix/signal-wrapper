package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func runShutdownScript(ctx context.Context, shutdownScript string) {
	cmd := exec.CommandContext(ctx, shutdownScript)
	// Pass through stderr / stdout directly
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.WithField("error", err).Error("Error running shutdown script")
	}
}

func signalWatcher(ctx context.Context, cmd *exec.Cmd, shutdownScript string) {
	signalChan := make(chan os.Signal, 100)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	signal := <-signalChan
	log.Info("Received first signal: ", signal)
	log.WithField("shutdownScript", shutdownScript).Info("Running shutdown script")
	runShutdownScript(ctx, shutdownScript)
	log.WithField("shutdownScript", shutdownScript).Info("Shutdown script complete")
	log.Info("Forwarding signal: ", signal)
	cmd.Process.Signal(signal)

	for signal = range signalChan {
		log.Info("Forwarding signal: ", signal)
		cmd.Process.Signal(signal)
	}
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s [shutdown command] [cmd] [args]\n", os.Args[0])
		os.Exit(1)
	}
	shutdownScript := os.Args[1]

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	go signalWatcher(ctx, cmd, shutdownScript)
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
