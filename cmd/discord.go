package cmd

import (

    "ghosttools/internal/discord"
    "ghosttools/internal/ui"
)

func init() {
    ui.RegisterCategory(ui.Category{
        Name: "Discord",
        Items: []ui.Item{
            {"Webhook Spammer", "spam or delete discord webhooks", discord.WebhookTools},
            {"Token Panel", "account nuker, token login, info", discord.TokenTools},
            {"Guild Lookup", "retrieve guild information", discord.ServerInfo},
            {"Invite Builder", "generate admin bot invite links", discord.BotInviteGen},
            {"Self Bot", "launch advanced selfbot", discord.Selfbot},
            {"Server Copier", "clone a discord server using a token", discord.ServerCloner},
            {"Server Nuker", "advanced server destruction console", discord.NukeBot},
            {"Username Lookup", "check availability of usernames", discord.UsernameChecker},
            {"Token Checker", "validate tokens in bulk", discord.TokenChecker},
            {"Server Backup", "backup full server structure", discord.GuildBackup},
            {"Chat Logger", "log messages from a channel", discord.MessageLogger},
            {"Voice Raid", "join and spam voice channels", discord.VoiceSpam},
        },
    })
}