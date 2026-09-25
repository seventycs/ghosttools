package hwid

import (
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"ghosttools/internal/ui"
)

// ---------- helpers ----------

func runPS(script string) (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func randHex(n int) string {
	const chars = "0123456789ABCDEF"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func randUUID() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		randHex(8), randHex(4), randHex(4), randHex(4), randHex(12))
}

func randMAC() string {
	return fmt.Sprintf("%02X-%02X-%02X-%02X-%02X-%02X",
		rand.Intn(256), rand.Intn(256), rand.Intn(256),
		rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func randStr(n int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// ---------- show ----------

// ShowCurrent dumps current HWID-identifying values.
func ShowCurrent() {
	ui.Cyan("reading current hardware identifiers...")

	fields := []struct {
		name string
		ps   string
	}{
		{"MachineGuid", `(Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Cryptography').MachineGuid`},
		{"ComputerName", `$env:COMPUTERNAME`},
		{"ProductId", `(Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion').ProductId`},
		{"VolumeSerial", `cmd /c vol C:`},
		{"BIOS Serial", `(Get-CimInstance Win32_BIOS).SerialNumber`},
		{"Baseboard Serial", `(Get-CimInstance Win32_BaseBoard).SerialNumber`},
		{"Disk Serial", `(Get-CimInstance Win32_DiskDrive | Select-Object -First 1).SerialNumber`},
		{"GPU", `(Get-CimInstance Win32_VideoController | Select-Object -First 1).Name`},
		{"CPU ID", `(Get-CimInstance Win32_Processor | Select-Object -First 1).ProcessorId`},
		{"UUID (SMBIOS)", `(Get-CimInstance Win32_ComputerSystemProduct).UUID`},
	}
	for _, f := range fields {
		out, _ := runPS(f.ps)
		fmt.Printf("  %-18s %s\n", ui.Cyan(f.name+":"), ui.White(out))
	}

	fmt.Println()
	ui.Cyan("network adapters:")
	out, _ := runPS(`Get-NetAdapter | Select-Object Name,MacAddress,Status | Format-Table -AutoSize | Out-String`)
	fmt.Println(out)

	fmt.Println()
	ui.Yellow("changing these to random values = spoof. admin required.")
}

// ---------- spoof ----------

// SpoofMachineGuid replaces HKLM\SOFTWARE\Microsoft\Cryptography\MachineGuid.
func SpoofMachineGuid() {
	ui.Cyan("spoofing MachineGuid...")
	newGUID := randUUID()
	script := fmt.Sprintf(`Set-ItemProperty -Path 'HKLM:\SOFTWARE\Microsoft\Cryptography' -Name MachineGuid -Value '%s'`, newGUID)
	if _, err := runPS(script); err != nil {
		ui.Red("failed (need admin): " + err.Error())
		return
	}
	ui.Green("new MachineGuid: " + newGUID)
}

// SpoofProductID replaces the Windows product ID.
func SpoofProductID() {
	ui.Cyan("spoofing Windows ProductId...")
	newID := randStr(5) + "-" + randStr(5) + "-" + randStr(5) + "-" + randStr(5)
	script := fmt.Sprintf(`Set-ItemProperty -Path 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion' -Name ProductId -Value '%s'`, newID)
	if _, err := runPS(script); err != nil {
		ui.Red("failed (need admin): " + err.Error())
		return
	}
	ui.Green("new ProductId: " + newID)
}

// SpoofComputerName renames the PC and needs a reboot.
func SpoofComputerName() {
	ui.Cyan("spoofing computer name...")
	newName := "DESKTOP-" + randStr(7)
	if _, err := runPS(fmt.Sprintf(`Rename-Computer -NewName '%s' -Force`, newName)); err != nil {
		ui.Red("failed (need admin): " + err.Error())
		return
	}
	ui.Green("new name: " + newName + " (reboot required)")
}

// SpoofMAC changes MAC on every enabled network adapter.
func SpoofMAC() {
	ui.Cyan("spoofing MAC on all enabled adapters...")

	// get all enabled adapters
	out, err := runPS(`Get-NetAdapter | Where-Object {$_.Status -eq 'Up'} | Select-Object -ExpandProperty Name`)
	if err != nil || out == "" {
		ui.Red("couldn't list adapters: " + err.Error())
		return
	}
	adapters := strings.Split(out, "\n")
	for _, a := range adapters {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		mac := randMAC()
		if _, err := runPS(fmt.Sprintf(`Set-NetAdapter -Name '%s' -MacAddress '%s' -Confirm:$false`, a, mac)); err != nil {
			ui.Red(fmt.Sprintf("failed on %s: %v", a, err))
			continue
		}
		ui.Green(fmt.Sprintf("%s -> %s", a, mac))
	}
}

// SpoofMACRegistry changes MAC via registry (persists across reboots).
func SpoofMACRegistry() {
	ui.Cyan("spoofing MAC via registry (persists after reboot)...")
	base := `HKLM:\SYSTEM\CurrentControlSet\Control\Class\{4d36e972-e325-11ce-bfc1-08002be10318}`
	script := fmt.Sprintf(`
Get-ChildItem '%s' -ErrorAction SilentlyContinue | ForEach-Object {
	$name = $_.PSChildName
	if ($name -match '^\d{4}$') {
		$mac = (1..6 | ForEach-Object { '{0:X2}' -f (Get-Random -Maximum 256) }) -join ''
		Set-ItemProperty -Path $_.PSPath -Name 'NetworkAddress' -Value $mac -ErrorAction SilentlyContinue
		Write-Host "[+] $name -> $mac"
	}
}`, base)
	out, _ := runPS(script)
	fmt.Println(out)
	ui.Green("done. disable/re-enable adapter to apply.")
}

// SpoofVolumeSerial changes the NTFS volume serial of the C: drive.
// requires SysInternals volumeid.exe in C:\Windows\System32.
func SpoofVolumeSerial() {
	ui.Cyan("spoofing C: volume serial...")

	// check if volumeid exists
	if _, err := runPS(`Test-Path 'C:\Windows\System32\volumeid.exe'`); err != nil {
		ui.Red("error checking volumeid: " + err.Error())
		return
	}
	check, _ := runPS(`Test-Path 'C:\Windows\System32\volumeid.exe'`)
	if strings.TrimSpace(check) != "True" {
		ui.Yellow("SysInternals VolumeID not found.")
		ui.Yellow("download volumeid.exe from live.sysinternals.com and place at C:\\Windows\\System32\\volumeid.exe")
		return
	}

	newSerial := randHex(4) + "-" + randHex(4)
	cmd := exec.Command("C:\\Windows\\System32\\volumeid.exe", "C:", newSerial)
	out, err := cmd.CombinedOutput()
	if err != nil {
		ui.Red("failed: " + string(out))
		return
	}
	fmt.Println(string(out))
	ui.Green("done (reboot required).")
}

// SpoofDisplayAdapter renames the GPU in the registry.
func SpoofDisplayAdapter() {
	ui.Cyan("spoofing display adapter name...")
	script := `
$newNames = @('NVIDIA GeForce RTX 4090','AMD Radeon RX 7900 XTX','Intel Arc A770','NVIDIA RTX 4080')
$base = 'HKLM:\SYSTEM\CurrentControlSet\Control\Class\{4d36e968-e325-11ce-bfc1-08002be10318}'
Get-ChildItem $base -ErrorAction SilentlyContinue | ForEach-Object {
	$name = $_.PSChildName
	if ($name -match '^\d{4}$') {
		$newName = $newNames | Get-Random
		Set-ItemProperty -Path $_.PSPath -Name 'DriverDesc' -Value $newName -ErrorAction SilentlyContinue
		Write-Host "[+] $name -> $newName"
	}
}`
	out, _ := runPS(script)
	fmt.Println(out)
	ui.Green("done (reboot required).")
}

// SpoofCPUString changes the processor string in the registry.
func SpoofCPUString() {
	ui.Cyan("spoofing CPU string...")
	newCPU := []string{
		"AMD Ryzen 9 7950X 16-Core Processor",
		"Intel(R) Core(TM) i9-14900K",
		"AMD Ryzen 7 7800X3D 8-Core Processor",
		"Intel(R) Xeon(R) W-2495X",
	}[rand.Intn(4)]
	script := fmt.Sprintf(`Set-ItemProperty -Path 'HKLM:\HARDWARE\DESCRIPTION\System\CentralProcessor\0' -Name ProcessorNameString -Value '%s'`, newCPU)
	if _, err := runPS(script); err != nil {
		ui.Red("failed (need admin): " + err.Error())
		return
	}
	ui.Green("new CPU: " + newCPU)
}

// SpoofBIOSSerial attempts to change SMBIOS serial (needs admin + AMI tooling).
func SpoofBIOSSerial() {
	ui.Yellow("BIOS serial spoofing requires AMIDEWINx64.exe (AMI tools).")
	ui.Yellow("Place AMIDEWINx64.exe in C:\\Windows\\System32, then run this.")
	check, _ := runPS(`Test-Path 'C:\Windows\System32\AMIDEWINx64.exe'`)
	if strings.TrimSpace(check) != "True" {
		ui.Red("AMIDEWINx64.exe not found.")
		return
	}
	newSerial := randStr(10)
	cmd := exec.Command("C:\\Windows\\System32\\AMIDEWINx64.exe", "/SS", newSerial)
	out, err := cmd.CombinedOutput()
	if err != nil {
		ui.Red("failed: " + string(out))
		return
	}
	ui.Green("done: " + newSerial)
}

// SpoofSMBIOSUUID attempts to change the SMBIOS UUID.
func SpoofSMBIOSUUID() {
	ui.Yellow("SMBIOS UUID spoofing requires AMIDEWINx64.exe.")
	check, _ := runPS(`Test-Path 'C:\Windows\System32\AMIDEWINx64.exe'`)
	if strings.TrimSpace(check) != "True" {
		ui.Red("AMIDEWINx64.exe not found.")
		return
	}
	newUUID := randUUID()
	cmd := exec.Command("C:\\Windows\\System32\\AMIDEWINx64.exe", "/SU", "auto")
	out, _ := cmd.CombinedOutput()
	fmt.Println(string(out))
	ui.Green("SMBIOS UUID spoofed (verify with ShowCurrent). new uuid = " + newUUID)
}

// ---------- meta ----------

// SpoofAll attempts every spoof in sequence.
func SpoofAll() {
	ui.Cyan("=== spoofing every hardware identifier ===")
	ui.Cyan("admin required for most operations")
	fmt.Println()
	fmt.Println(ui.NeonGreen("— MachineGuid —"))
	SpoofMachineGuid()
	time.Sleep(300 * time.Millisecond)

	fmt.Println()
	fmt.Println(ui.NeonGreen("— ProductId —"))
	SpoofProductID()
	time.Sleep(300 * time.Millisecond)

	fmt.Println()
	fmt.Println(ui.NeonGreen("— ComputerName —"))
	SpoofComputerName()
	time.Sleep(300 * time.Millisecond)

	fmt.Println()
	fmt.Println(ui.NeonGreen("— MAC (live) —"))
	SpoofMAC()
	time.Sleep(300 * time.Millisecond)

	fmt.Println()
	fmt.Println(ui.NeonGreen("— CPU String —"))
	SpoofCPUString()
	time.Sleep(300 * time.Millisecond)

	fmt.Println()
	fmt.Println(ui.NeonGreen("— GPU Name —"))
	SpoofDisplayAdapter()

	fmt.Println()
	ui.Green("=== spoof complete — reboot required for everything to take effect ===")
	ui.Yellow("note: some IDs (BIOS serial, SMBIOS UUID, VolumeID) need external tools.")
}

// RestoreWarning is not a real restore — changing IDs isn't reversible.
func RestoreWarning() {
	ui.Red("hardware IDs cannot be easily restored after spoofing.")
	ui.Red("if you have a backup of the original values, restore manually.")
	ui.Yellow("values we changed:")
	ui.Yellow("  - MachineGuid (registry)")
	ui.Yellow("  - ProductId (registry)")
	ui.Yellow("  - ComputerName")
	ui.Yellow("  - MAC address (live + registry)")
	ui.Yellow("  - GPU DriverDesc")
	ui.Yellow("  - CPU ProcessorNameString")
	ui.Yellow("reboot restores nothing — these persist until changed again.")
}

// ---------- extra ----------

// WipeHWIDTrace clears things that log HWID-history.
func WipeHWIDTrace() {
	ui.Cyan("wiping HWID traces...")
	cmds := []string{
		`wevtutil cl System`,
		`wevtutil cl Application`,
		`wevtutil cl Security`,
		`wevtutil cl Setup`,
		`del /f /q /s C:\Windows\Prefetch\* >nul 2>&1`,
		`del /f /q C:\Windows\Temp\* >nul 2>&1`,
	}
	for _, c := range cmds {
		exec.Command("cmd", "/c", c).Run()
	}
	ui.Green("event logs + prefetch wiped.")
}

// DetectSpoof shows which IDs look spoofed (all-zero, FF, or random format).
func DetectSpoof() {
	ui.Cyan("checking for common spoof signatures...")
	out, _ := runPS(`(Get-ItemProperty 'HKLM:\SOFTWARE\Microsoft\Cryptography').MachineGuid`)
	if !regexp.MustCompile(`^[0-9a-fA-F]{8}-`).MatchString(out) {
		ui.Red("MachineGuid looks malformed or spoofed: " + out)
	} else {
		ui.Green("MachineGuid format looks normal: " + out)
	}
	bios, _ := runPS(`(Get-CimInstance Win32_BIOS).SerialNumber`)
	if strings.Contains(strings.ToLower(bios), "to be filled") || strings.Contains(strings.ToLower(bios), "0") {
		ui.Red("BIOS serial looks fake: " + bios)
	} else {
		ui.Green("BIOS serial: " + bios)
	}
}