package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"fyne.io/systray"
)

func systrayQuit() {
	systray.Quit()
}

func runTray(httpAddr string, quit func()) {
	dashURL := "http://" + httpAddr
	systray.Run(func() {
		if runtime.GOOS == "windows" {
			systray.SetIcon(iconICO)
		} else {
			systray.SetIcon(iconPNG)
		}
		systray.SetTitle("Sector One")
		systray.SetTooltip("Sector One — " + dashURL)

		openItem := systray.AddMenuItem("Open dashboard", "Open the dash in the browser")
		systray.AddSeparator()
		quitItem := systray.AddMenuItem("Quit", "Stop the engine")

		go func() {
			for {
				select {
				case <-openItem.ClickedCh:
					openBrowser(dashURL)
				case <-quitItem.ClickedCh:
					quit()
					systray.Quit()
					return
				}
			}
		}()

		openBrowser(dashURL)
	}, nil)
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("open dashboard:", err)
	}
}
