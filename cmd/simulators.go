package cmd

import (
    "ghosttools/internal/simulators"
    "ghosttools/internal/ui"
)

func init() {
    ui.RegisterCategory(ui.Category{
        Name: "Simulators",
        Items: []ui.Item{
            {"Fake Token", "generate fake discord tokens", simulators.FakeTokenGen},
            {"Fake Mail", "generate fake email addresses", simulators.FakeMailGen},
            {"Fake ID", "generate a full Fake ID", simulators.FakeIdentity},
            {"Fake CC", "generate fake CC numbers", simulators.FakeCreditCard},
            {"Social Bot", "simulate social media growth", simulators.SocialBotter},
            {"Fake PayPal", "simulate paypal OTP page", simulators.FakePaypalOTP},
            {"Fake Fortnite", "simulate fortnite account check", simulators.FakeFortnite},
            {"Fake Exodu", "simulate exodus wallet page", simulators.FakeExodus},
            {"Hacker Terminal", "fake hacker terminal for pranks", simulators.HackerTerminal},
            {"Fake Brute", "simulate a bruteforce attack", simulators.FakeBruteforcer},
            {"QR Code", "generate standard and fake QR", simulators.QRCodeGen},
            {"Fake Discord", "simulate discord login page", simulators.FakeDiscord},
            {"Fake Steam", "simulate steam login page", simulators.FakeSteam},
            {"Fake Instagram", "simulate instagram login page", simulators.FakeInstagram},
        },
    })
}