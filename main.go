package main

import (
	"log"
	"os"
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
		Args:        []string{"[go-daemon usdt-monitor]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("Monitor started")

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
			cntxt.Release()
			os.Exit(0)
		}
	}
}
