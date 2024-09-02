package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/sevlyar/go-daemon"
)

func main() {

	cntxt := &daemon.Context{
		PidFileName: "usdt-monitor.pid",
		PidFilePerm: 0644,
		LogFileName: "usdt-monitor.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"go usdt-monitor"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("Monitor starting...")

	// Start signal handler in a separate goroutine
	go handleSignals(cntxt)

	// Start the monitoring loop
	startMonitoring()

}

func handleSignals(cntxt *daemon.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)

	for {
		sig := <-sigChan
		if sig == syscall.SIGHUP {
			log.Print("Received SIGHUP, restarting...")
			restart(cntxt)
		}
	}
}

func restart(cntxt *daemon.Context) {
	cntxt.Release()
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	cmd := exec.Command(exe)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to restart process: %v", err)
	}

	os.Exit(0)
}
