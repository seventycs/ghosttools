package cmd

import (
	"ghosttools/internal/hwid"
	"ghosttools/internal/ui"
)

func init() {
	ui.RegisterCategory(ui.Category{
		Name: "HWID",
		Items: []ui.Item{
			{"Show IDs",         "display all current hardware identifiers",  hwid.ShowCurrent},
			{"Spoof MachineGuid","change the windows machine guid",           hwid.SpoofMachineGuid},
			{"Spoof ProductId",  "change the windows product id",             hwid.SpoofProductID},
			{"Spoof PC Name",    "rename the computer",                       hwid.SpoofComputerName},
			{"Spoof MAC (live)", "change mac on all active adapters",         hwid.SpoofMAC},
			{"Spoof MAC (reg)",  "change mac via registry (persists reboot)", hwid.SpoofMACRegistry},
			{"Spoof Volume ID",  "change c: ntfs serial (needs volumeid.exe)", hwid.SpoofVolumeSerial},
			{"Spoof GPU Name",   "rename the graphics adapter",               hwid.SpoofDisplayAdapter},
			{"Spoof CPU String", "change the processor string",               hwid.SpoofCPUString},
			{"Spoof BIOS Serial","change the smbios serial (needs AMI tool)", hwid.SpoofBIOSSerial},
			{"Spoof SMBIOS UUID","change the smbios uuid (needs AMI tool)",   hwid.SpoofSMBIOSUUID},
			{"Spoof All",        "run every spoof in sequence",               hwid.SpoofAll},
			{"Wipe HWID Trace",  "clear event logs + prefetch",               hwid.WipeHWIDTrace},
			{"Detect Spoof",     "check for common spoof signatures",         hwid.DetectSpoof},
			{"Restore Warning",  "info about restoring ids",                  hwid.RestoreWarning},
		},
	})
}