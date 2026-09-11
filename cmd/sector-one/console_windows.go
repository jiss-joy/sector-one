//go:build windows

package main

import "golang.org/x/sys/windows"

func hideConsole() {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	user32 := windows.NewLazySystemDLL("user32.dll")
	windowHandle, _, _ := kernel32.NewProc("GetConsoleWindow").Call()
	if windowHandle == 0 {
		return
	}
	user32.NewProc("ShowWindow").Call(windowHandle, 0) // SW_HIDE
}
