package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	_ "embed"

	"github.com/getlantern/systray"
)

var (
	//go:embed gamepad.png
	icon []byte

	// Sub-process handler
	xboxdrv *exec.Cmd

	// Evdev symlink, vendor and OS dependent
	evdev string
)

func init() {
	flag.StringVar(&evdev, "j", "/dev/input/by-id/usb-SHANWAN_PS3_PC_Gamepad-event-joystick", "Joysitk event `DEVICE`")
}

func main() {
	flag.Parse()
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon)

	// Add exit menu and exit menu handler to stop the app
	exitMenu := systray.AddMenuItem("Exit", "Terminates the controller emulation")
	go func() {
		<-exitMenu.ClickedCh
		systray.Quit()
	}()

	// Launch the sub-process
	xboxdrv = exec.Command("xboxdrv", "--evdev", evdev, "--evdev-debug", "--evdev-no-grab", "--mimic-xpad")
	xboxdrv.Stdout = os.Stdout
	xboxdrv.Stderr = os.Stderr
	if err := xboxdrv.Start(); err != nil {
		log.Printf("Cound not initialize: %v", err)
		systray.Quit()
	}

	msg := fmt.Sprintf("Running with PID %v", xboxdrv.Process.Pid)
	log.Println(msg)
	systray.SetTitle(msg)
	systray.SetTooltip(msg)
}

func onExit() {
	if xboxdrv != nil {
		if err := xboxdrv.Process.Kill(); err != nil {
			log.Printf("Error killing process: %v", err)
		}
	}
}
