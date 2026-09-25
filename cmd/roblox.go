package cmd

import (
    "ghosttools/internal/roblox"
    "ghosttools/internal/ui"
)

func init() {
    ui.RegisterCategory(ui.Category{
        Name: "Roblox",
        Items: []ui.Item{
            {"Account Info", "retrieve detailed roblox account info", roblox.UserInfo},
            {"Cookie Info", "validate and check cookie data", roblox.CookieInfo},
            {"Cookie Login", "login via roblox cookie", roblox.CookieLogin},
            {"Group Lookup", "analyze details of target groups", roblox.GroupInfo},
            {"Asset Pull", "download game assets and textures", roblox.AssetDownload},
            {"Name Log", "track and display past usernames", roblox.NameHistory},
            {"Username Lookup", "check availability of usernames", roblox.UsernameChecker},
            {"Cookie Refresh", "generate new cookie from existing", roblox.CookieRefresher},
            {"Game Lookup", "fetch game details and stats", roblox.GameInfo},
            {"Inventory Pull", "dump user inventory items", roblox.InventoryDump},
            {"Trade Scanner", "scan user trades for value", roblox.TradeScanner},
            {"Limited Sniper", "snipe limited items on sale", roblox.LimitedSniper},
            {"Group Copier", "clone a roblox group structure", roblox.GroupCloner},
            {"Follower Bot", "mass follow a user with bots", roblox.BotFollower},
        },
    })
}