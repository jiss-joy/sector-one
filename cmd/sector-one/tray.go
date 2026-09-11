package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"fyne.io/systray"
)

func runTray(httpAddr string, quit func()) {
	dashboardURL := "http://" + httpAddr
	systray.Run(func() { setupTray(dashboardURL, quit) }, nil)
}

func setupTray(dashboardURL string, quit func()) {
	if runtime.GOOS == "windows" {
		systray.SetIcon(iconICO)
	} else {
		systray.SetIcon(iconPNG)
	}
	systray.SetTitle("Sector One")
	systray.SetTooltip("Sector One - " + dashboardURL)
	openItem := systray.AddMenuItem("Open dashboard", "Open the dashboard in the browser")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Stop the engine")

	go func() {
		for {
			select {
			case <-openItem.ClickedCh:
				openBrowser(dashboardURL)
			case <-quitItem.ClickedCh:
				quit()
				systray.Quit()
				return
			}
		}
	}()
}

func systrayQuit() {
	systray.Quit()
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
	err := cmd.Start()
	if err != nil {
		fmt.Println("open dashboard", err)
	}
}
