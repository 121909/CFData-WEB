//go:build windows

package main

import (
	"os/exec"
)

func openBrowserOnWindows(targetURL string) {
	if targetURL == "" {
		return
	}
	go func() {
		cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", targetURL)
		if err := cmd.Start(); err == nil && cmd.Process != nil {
			_ = cmd.Process.Release()
		}
	}()
}
