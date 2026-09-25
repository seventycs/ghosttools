package sysgen

import (
	"fmt"
	"os/exec"
	"runtime"

	"ghosttools/internal/ui"
)

func Debloat() {
	if runtime.GOOS != "windows" {
		ui.Red("windows only")
		return
	}
	ui.Cyan("running windows debloater...")

	type task struct {
		name string
		cmd  []string
	}
	tasks := []task{
		{"removing xbox", []string{"powershell", "-Command", "Get-AppxPackage *xbox* | Remove-AppxPackage -ErrorAction SilentlyContinue"}},
		{"removing bing", []string{"powershell", "-Command", "Get-AppxPackage *bing* | Remove-AppxPackage -ErrorAction SilentlyContinue"}},
		{"removing skype", []string{"powershell", "-Command", "Get-AppxPackage *skype* | Remove-AppxPackage -ErrorAction SilentlyContinue"}},
		{"removing zune", []string{"powershell", "-Command", "Get-AppxPackage *zune* | Remove-AppxPackage -ErrorAction SilentlyContinue"}},
		{"removing solitaire", []string{"powershell", "-Command", "Get-AppxPackage *solitaire* | Remove-AppxPackage -ErrorAction SilentlyContinue"}},
		{"removing onedrive", []string{"powershell", "-Command", "Get-AppxPackage *OneDrive* | Remove-AppxPackage -ErrorAction SilentlyContinue"}},
		{"removing cortana", []string{"powershell", "-Command", "Get-AppxPackage *Cortana* | Remove-AppxPackage -ErrorAction SilentlyContinue"}},
		{"disabling telemetry", []string{"powershell", "-Command", "Set-ItemProperty -Path 'HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows\\DataCollection' -Name AllowTelemetry -Value 0 -ErrorAction SilentlyContinue"}},
		{"disabling cortana policy", []string{"powershell", "-Command", "New-Item -Path 'HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows\\Windows Search' -Force | Out-Null; Set-ItemProperty -Path 'HKLM:\\SOFTWARE\\Policies\\Microsoft\\Windows\\Windows Search' -Name AllowCortana -Value 0"}},
		{"disabling app suggestions", []string{"powershell", "-Command", "New-Item -Path 'HKCU:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\ContentDeliveryManager' -Force | Out-Null; Set-ItemProperty -Path 'HKCU:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\ContentDeliveryManager' -Name SilentInstalledAppsEnabled -Value 0"}},
		{"disabling scheduled telemetry", []string{"powershell", "-Command", "Disable-ScheduledTask -TaskName 'Microsoft\\Windows\\Application Experience\\Microsoft Compatibility Appraiser' -ErrorAction SilentlyContinue"}},
		{"disabling feedback", []string{"powershell", "-Command", "Disable-ScheduledTask -TaskName 'Microsoft\\Windows\\Feedback\\Siuf\\DmClient' -ErrorAction SilentlyContinue"}},
		{"removing edge desktop icon", []string{"powershell", "-Command", "Remove-Item 'C:\\Users\\Public\\Desktop\\Microsoft Edge.lnk' -ErrorAction SilentlyContinue"}},
		{"removing onedrive startup", []string{"powershell", "-Command", "Remove-Item 'HKCU:\\Software\\Microsoft\\Windows\\CurrentVersion\\Run\\OneDrive' -ErrorAction SilentlyContinue"}},
	}

	for _, t := range tasks {
		fmt.Println(ui.Cyan("[*] ") + t.name)
		exec.Command(t.cmd[0], t.cmd[1:]...).Run()
	}
	ui.Green("debloat complete")
}