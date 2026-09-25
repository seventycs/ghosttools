#!/usr/bin/env python3
"""
ghosttools selfbot v1
prefix: +
"""
import discord
from discord.ext import commands
import asyncio
import aiohttp
import json
import os
import random
import time
import io
import subprocess
from datetime import datetime, timedelta

PREFIX = "+"
TOKEN_FILE = os.path.join(os.environ.get("TEMP", "/tmp"), ".selfbot_token")

bot = commands.Bot(command_prefix=PREFIX, self_bot=True, help_command=None)
START_TIME = time.time()

SNIPE_CACHE = {}
EDITSNIPE_CACHE = {}
REACTIONSNIPE_CACHE = {}
AUTOREPLIES = {}
MESSAGE_LOG = []
AFK_STATE = {"afk": False, "reason": ""}

# ==================================================================
# TOKEN
# ==================================================================
def load_token():
    if os.path.exists(TOKEN_FILE):
        with open(TOKEN_FILE) as f:
            return f.read().strip()
    return None

def save_token(t):
    with open(TOKEN_FILE, "w") as f:
        f.write(t)

# ==================================================================
# HELPERS
# ==================================================================
def make_embed(title, description=None, color=0x00ffcc, fields=None, footer=None, thumb=None, image=None):
    e = discord.Embed(title=title, description=description, color=color)
    if fields:
        for name, val, inline in fields:
            e.add_field(name=name, value=str(val)[:1024], inline=inline)
    if thumb:
        e.set_thumbnail(url=thumb)
    if image:
        e.set_image(url=image)
    e.set_footer(text=footer or "ghosttools v1")
    return e

async def send_embed(ctx, e, delete_after=None):
    """flatten embeds into plain text — discord.py-self embeds are broken on this version"""
    try:
        await ctx.message.delete()
    except Exception:
        pass

    lines = []
    d = e.to_dict()
    title = d.get("title")
    desc = d.get("description")
    if title:
        lines.append(f"**{title}**")
    if desc:
        lines.append(desc)
    for f in (d.get("fields") or []):
        name = f.get("name", "")
        val = f.get("value", "")
        lines.append(f"**{name}**\n{val}")
    footer = (d.get("footer") or {}).get("text")
    if footer:
        lines.append(f"_— {footer}_")

    text = "\n\n".join(lines) if lines else "_(empty)_"
    if len(text) > 1990:
        text = text[:1990] + "…"

    try:
        if delete_after:
            return await ctx.send(text, delete_after=delete_after)
        return await ctx.send(text)
    except Exception as ex:
        try:
            return await ctx.send(f"`err: {ex}`", delete_after=5)
        except Exception:
            return None

async def ghost_reply(ctx, text, delete=5):
    try:
        await ctx.message.delete()
    except Exception:
        pass
    try:
        return await ctx.send(text, delete_after=delete)
    except Exception:
        return None

async def safe_delete(msg, delay=3):
    try:
        await asyncio.sleep(delay)
        await msg.delete()
    except Exception:
        pass

# ==================================================================
# EVENTS
# ==================================================================
@bot.event
async def on_ready():
    print(f"[ghosttools selfbot v1] {bot.user}")
    print(f"  id:      {bot.user.id}")
    print(f"  guilds:  {len(bot.guilds)}")
    try:
        await bot.change_presence(status=discord.Status.online, activity=discord.Game(name="+help"))
    except Exception:
        pass

@bot.event
async def on_command_error(ctx, error):
    try:
        await ctx.message.delete()
    except Exception:
        pass
    try:
        await ctx.send(f"`err: {error}`", delete_after=5)
    except Exception:
        pass

@bot.event
async def on_message_delete(message):
    if message.author.bot:
        return
    SNIPE_CACHE[message.channel.id] = {
        "content": message.content,
        "author": str(message.author),
        "time": datetime.utcnow(),
    }

@bot.event
async def on_message_edit(before, after):
    if before.content == after.content or before.author.bot:
        return
    EDITSNIPE_CACHE[before.channel.id] = {
        "before": before.content,
        "after": after.content,
        "author": str(before.author),
    }

@bot.event
async def on_reaction_add(reaction, user):
    if user.bot:
        return
    REACTIONSNIPE_CACHE[reaction.message.channel.id] = {
        "emoji": str(reaction.emoji),
        "author": str(user),
        "message": reaction.message.content,
    }

@bot.event
async def on_message(message):
    MESSAGE_LOG.append({
        "guild": message.guild.id if message.guild else None,
        "channel": message.channel.id,
        "author": str(message.author),
        "author_id": message.author.id,
        "content": message.content,
        "time": time.time(),
    })
    if len(MESSAGE_LOG) > 500:
        MESSAGE_LOG.pop(0)
    if message.author.id != (bot.user.id if bot.user else 0):
        content_lower = message.content.lower()
        for trigger, reply in AUTOREPLIES.items():
            if trigger in content_lower:
                try:
                    await message.channel.send(reply)
                except Exception:
                    pass
                break
    await bot.process_commands(message)

# ==================================================================
# HELP
# ==================================================================
@bot.command(name="help")
async def help_cmd(ctx):
    fields = [
        ("utility",
         "`+ping` `+uptime` `+userinfo [@u]` `+avatar [@u]` `+banner [@u]` "
         "`+serverinfo` `+roles` `+emojis` `+invites` `+tokeninfo` `+whois <id>` "
         "`+created <id>` `+joined [@u]` `+permissions [@u]` `+nick [@u]`"),
        ("messages",
         "`+say <txt>` `+embed <title> | <desc>` `+purge <n>` `+spam <n> <msg>` "
         "`+massdm <msg>` `+pin` `+unpin` `+react <emoji>` `+edit <id> <txt>` "
         "`+snipe` `+editsnipe` `+reactionsnipe` `+clearsnipe`"),
        ("moderation",
         "`+kick <@u>` `+ban <@u>` `+unban <id>` `+mute <@u>` `+unmute <@u>` "
         "`+deafen <@u>` `+undeafen <@u>` `+timeout <@u> [mins]` `+untimeout <@u>` "
         "`+nickname <@u> <name>` `+role add|remove <@u> <role>`"),
        ("status",
         "`+status <online|idle|dnd|invisible>` `+playing <txt>` `+streaming <txt>` "
         "`+listening <txt>` `+watching <txt>` `+custom <txt>` `+clearstatus`"),
        ("voice",
         "`+joinvc [id]` `+leavevc` `+vcspam <id>` `+deafen_self` `+mute_self`"),
        ("ghost",
         "`+ghost <secs>` `+ghostdm <@u> <msg>` `+clear` `+ghostall <secs>` `+ghostchannel`"),
        ("recon",
         "`+track <@u>` `+messages <@u>` `+lastseen <@u>` `+search <query>` "
         "`+history [n]` `+typing` `+afk [reason]` `+back`"),
        ("mass actions",
         "`+massban <ids>` `+masskick` `+massrole <role> add|remove` `+massinvite`"),
        ("fun",
         "`+8ball <q>` `+coinflip` `+dice [n]` `+meme` `+joke` `+roast [@u]` "
         "`+compliment [@u]` `+rate <thing>` `+gay [@u]` `+ship <@u1> <@u2>`"),
        ("server mgmt",
         "`+createchannel <n>` `+deletechannel <id>` `+createrole <n>` "
         "`+deleterole <id>` `+createemoji <n> <url>` `+webhooks` "
         "`+createwebhook <n>` `+delwebhook <id>`"),
        ("backup",
         "`+backup` `+backup_roles` `+backup_channels` `+backup_emojis` "
         "`+dumpmembers` `+dumpinvites`"),
        ("automation",
         "`+autoreply <trigger> | <response>` `+clearautoreply`"),
        ("owner",
         "`+eval <code>` `+exec <shell>` `+restart` `+shutdown` `+log`"),
    ]
    e = make_embed("ghosttools selfbot v1", f"prefix `{PREFIX}` · {len(bot.guilds)} guilds")
    for n, v in fields:
        e.add_field(name=n, value=v, inline=False)
    await send_embed(ctx, e)

# ==================================================================
# UTILITY
# ==================================================================
@bot.command(name="ping")
async def ping_cmd(ctx):
    await ghost_reply(ctx, f"`{round(bot.latency*1000)}ms`")

@bot.command(name="uptime")
async def uptime_cmd(ctx):
    s = int(time.time() - START_TIME)
    h, m, sec = s//3600, (s%3600)//60, s%60
    await ghost_reply(ctx, f"`{h}h {m}m {sec}s`")

@bot.command(name="userinfo", aliases=["ui"])
async def userinfo_cmd(ctx, user: discord.User = None):
    user = user or ctx.author
    e = make_embed(f"{user}", thumb=user.display_avatar.url)
    e.add_field(name="id", value=str(user.id), inline=True)
    e.add_field(name="bot", value=str(user.bot), inline=True)
    e.add_field(name="created", value=user.created_at.strftime("%Y-%m-%d"), inline=True)
    if ctx.guild and isinstance(user, discord.Member):
        e.add_field(name="roles", value=str(len(user.roles)-1), inline=True)
        if user.nick:
            e.add_field(name="nick", value=user.nick, inline=True)
    await send_embed(ctx, e)

@bot.command(name="avatar", aliases=["av"])
async def avatar_cmd(ctx, user: discord.User = None):
    user = user or ctx.author
    try:
        await ctx.message.delete()
    except Exception:
        pass
    await ctx.send(user.display_avatar.url)

@bot.command(name="banner")
async def banner_cmd(ctx, user: discord.User = None):
    user = user or ctx.author
    try:
        u = await bot.fetch_user(user.id)
    except Exception:
        u = user
    if not getattr(u, "banner", None):
        return await ghost_reply(ctx, "`no banner`")
    await send_embed(ctx, make_embed("banner", image=u.banner.url))

@bot.command(name="serverinfo")
async def serverinfo_cmd(ctx):
    g = ctx.guild
    if not g:
        return await ghost_reply(ctx, "`not in guild`")
    e = make_embed(f"{g.name}", thumb=g.icon.url if g.icon else None)
    e.add_field(name="id", value=str(g.id), inline=True)
    e.add_field(name="owner", value=str(g.owner), inline=True)
    e.add_field(name="members", value=str(g.member_count), inline=True)
    e.add_field(name="channels", value=str(len(g.channels)), inline=True)
    e.add_field(name="roles", value=str(len(g.roles)), inline=True)
    e.add_field(name="emojis", value=str(len(g.emojis)), inline=True)
    await send_embed(ctx, e)

@bot.command(name="roles")
async def roles_cmd(ctx):
    if not ctx.guild:
        return
    txt = ", ".join([f"`{r.name}`" for r in ctx.guild.roles][:50])
    await send_embed(ctx, make_embed(f"roles ({len(ctx.guild.roles)})", txt[:4000]))

@bot.command(name="emojis")
async def emojis_cmd(ctx):
    if not ctx.guild:
        return
    txt = " ".join(str(e) for e in ctx.guild.emojis[:50])
    await send_embed(ctx, make_embed(f"emojis ({len(ctx.guild.emojis)})", txt or "none"))

@bot.command(name="invites")
async def invites_cmd(ctx):
    if not ctx.guild:
        return
    try:
        invites = await ctx.guild.invites()
    except Exception:
        invites = []
    lines = [f"`{i.code}` — {i.uses} uses" for i in invites[:20]]
    await send_embed(ctx, make_embed(f"invites ({len(invites)})", "\n".join(lines) or "none"))

@bot.command(name="tokeninfo")
async def tokeninfo_cmd(ctx):
    async with aiohttp.ClientSession() as s:
        async with s.get("https://discord.com/api/v10/users/@me",
                         headers={"Authorization": bot.http.token}) as r:
            data = await r.json()
    e = make_embed("token info")
    for k, v in data.items():
        e.add_field(name=k, value=str(v)[:1024], inline=False)
    await send_embed(ctx, e)

@bot.command(name="whois")
async def whois_cmd(ctx, user_id: int = None):
    if not user_id:
        return await ghost_reply(ctx, "`usage: +whois <id>`")
    try:
        u = await bot.fetch_user(user_id)
    except Exception as e:
        return await ghost_reply(ctx, f"`err: {e}`")
    e = make_embed(f"{u}", thumb=u.display_avatar.url)
    e.add_field(name="id", value=str(u.id), inline=True)
    e.add_field(name="created", value=u.created_at.strftime("%Y-%m-%d"), inline=True)
    await send_embed(ctx, e)

@bot.command(name="created")
async def created_cmd(ctx, user_id: int = None):
    if not user_id:
        return await ghost_reply(ctx, "`usage: +created <id>`")
    ts = ((user_id >> 22) + 1420070400000) / 1000
    dt = datetime.utcfromtimestamp(ts)
    await ghost_reply(ctx, f"`{dt.strftime('%Y-%m-%d %H:%M:%S')} UTC`")

@bot.command(name="joined")
async def joined_cmd(ctx, user: discord.Member = None):
    user = user or ctx.author
    if not ctx.guild or not user.joined_at:
        return await ghost_reply(ctx, "`no joined_at`")
    await ghost_reply(ctx, f"`{user.joined_at.strftime('%Y-%m-%d %H:%M')}`")

@bot.command(name="permissions", aliases=["perms"])
async def permissions_cmd(ctx, user: discord.Member = None):
    user = user or ctx.author
    if not ctx.guild:
        return
    perms = [p for p, v in user.guild_permissions if v]
    txt = ", ".join(f"`{p}`" for p in perms[:40])
    await send_embed(ctx, make_embed(f"perms · {user}", txt or "none"))

@bot.command(name="nick")
async def nick_cmd(ctx, user: discord.Member = None):
    user = user or ctx.author
    await ghost_reply(ctx, f"`{user.nick or user.name}`")

# ==================================================================
# MESSAGES
# ==================================================================
@bot.command(name="say")
async def say_cmd(ctx, *, text: str = ""):
    try:
        await ctx.message.delete()
    except Exception:
        pass
    if text:
        await ctx.send(text)

@bot.command(name="embed")
async def embed_cmd(ctx, *, args: str = ""):
    if "|" in args:
        title, desc = args.split("|", 1)
    else:
        title, desc = args, None
    await send_embed(ctx, make_embed(title.strip(), desc.strip() if desc else None))

@bot.command(name="purge")
async def purge_cmd(ctx, n: int = 10):
    try:
        await ctx.message.delete()
    except Exception:
        pass
    try:
        deleted = await ctx.channel.purge(limit=n)
        await safe_delete(await ctx.send(f"`purged {len(deleted)}`"), 3)
    except Exception:
        pass

@bot.command(name="spam")
async def spam_cmd(ctx, n: int = 5, *, msg: str = "ghosttools"):
    try:
        await ctx.message.delete()
    except Exception:
        pass
    for _ in range(n):
        try:
            await ctx.send(msg)
        except Exception:
            break
        await asyncio.sleep(0.6)

@bot.command(name="massdm")
async def massdm_cmd(ctx, *, msg: str = "ghosttools"):
    try:
        await ctx.message.delete()
    except Exception:
        pass
    if not ctx.guild:
        return
    sent = 0
    for m in ctx.guild.members:
        if m.bot or m == bot.user:
            continue
        try:
            await m.send(msg)
            sent += 1
            await asyncio.sleep(1.2)
        except Exception:
            pass
    await safe_delete(await ctx.send(f"`sent to {sent}`"), 3)

@bot.command(name="pin")
async def pin_cmd(ctx):
    try:
        if ctx.message.reference:
            msg = await ctx.channel.fetch_message(ctx.message.reference.message_id)
            await msg.pin()
    except Exception:
        pass
    try:
        await ctx.message.delete()
    except Exception:
        pass

@bot.command(name="unpin")
async def unpin_cmd(ctx):
    try:
        if ctx.message.reference:
            msg = await ctx.channel.fetch_message(ctx.message.reference.message_id)
            await msg.unpin()
    except Exception:
        pass
    try:
        await ctx.message.delete()
    except Exception:
        pass

@bot.command(name="react")
async def react_cmd(ctx, emoji: str = "👻"):
    try:
        if ctx.message.reference:
            msg = await ctx.channel.fetch_message(ctx.message.reference.message_id)
            await msg.add_reaction(emoji)
    except Exception:
        pass
    try:
        await ctx.message.delete()
    except Exception:
        pass

@bot.command(name="edit")
async def edit_cmd(ctx, msg_id: int = None, *, text: str = ""):
    if not msg_id:
        return await ghost_reply(ctx, "`usage: +edit <id> <text>`")
    try:
        m = await ctx.channel.fetch_message(msg_id)
        if m.author.id != bot.user.id:
            return await ghost_reply(ctx, "`can only edit own`")
        await m.edit(content=text)
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="snipe")
async def snipe_cmd(ctx):
    data = SNIPE_CACHE.get(ctx.channel.id)
    if not data:
        return await ghost_reply(ctx, "`nothing`")
    e = make_embed("sniped", data["content"] or "(no text)", color=0xff6b6b)
    e.add_field(name="author", value=data["author"], inline=True)
    e.add_field(name="at", value=data["time"].strftime("%H:%M:%S"), inline=True)
    await send_embed(ctx, e)

@bot.command(name="editsnipe", aliases=["esnipe"])
async def editsnipe_cmd(ctx):
    data = EDITSNIPE_CACHE.get(ctx.channel.id)
    if not data:
        return await ghost_reply(ctx, "`nothing`")
    e = make_embed("editsnipe",
                   f"before: {data['before'][:500]}\nafter: {data['after'][:500]}",
                   color=0xffaa00)
    e.add_field(name="author", value=data["author"], inline=True)
    await send_embed(ctx, e)

@bot.command(name="reactionsnipe", aliases=["rsnipe"])
async def reactionsnipe_cmd(ctx):
    data = REACTIONSNIPE_CACHE.get(ctx.channel.id)
    if not data:
        return await ghost_reply(ctx, "`nothing`")
    e = make_embed("reactionsnipe", f"{data['author']} reacted {data['emoji']}")
    e.add_field(name="to", value=data["message"][:500], inline=False)
    await send_embed(ctx, e)

@bot.command(name="clearsnipe")
async def clearsnipe_cmd(ctx):
    SNIPE_CACHE.pop(ctx.channel.id, None)
    EDITSNIPE_CACHE.pop(ctx.channel.id, None)
    REACTIONSNIPE_CACHE.pop(ctx.channel.id, None)
    await ghost_reply(ctx, "`cleared`")

# ==================================================================
# MODERATION
# ==================================================================
@bot.command(name="kick")
async def kick_cmd(ctx, member: discord.Member = None, *, reason: str = "ghosttools"):
    if not member or not ctx.guild:
        return
    try:
        await member.kick(reason=reason)
        await ghost_reply(ctx, f"`kicked {member}`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="ban")
async def ban_cmd(ctx, member: discord.Member = None, *, reason: str = "ghosttools"):
    if not member or not ctx.guild:
        return
    try:
        await member.ban(reason=reason)
        await ghost_reply(ctx, f"`banned {member}`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="unban")
async def unban_cmd(ctx, user_id: int = None):
    if not user_id or not ctx.guild:
        return
    try:
        user = await bot.fetch_user(user_id)
        await ctx.guild.unban(user)
        await ghost_reply(ctx, f"`unbanned {user}`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="mute")
async def mute_cmd(ctx, member: discord.Member = None):
    if not member or not ctx.guild:
        return
    try:
        await member.edit(mute=True)
        await ghost_reply(ctx, "`muted`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="unmute")
async def unmute_cmd(ctx, member: discord.Member = None):
    if not member or not ctx.guild:
        return
    try:
        await member.edit(mute=False)
        await ghost_reply(ctx, "`unmuted`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="deafen")
async def deafen_cmd(ctx, member: discord.Member = None):
    if not member or not ctx.guild:
        return
    try:
        await member.edit(deafen=True)
        await ghost_reply(ctx, "`deafened`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="undeafen")
async def undeafen_cmd(ctx, member: discord.Member = None):
    if not member or not ctx.guild:
        return
    try:
        await member.edit(deafen=False)
        await ghost_reply(ctx, "`undeafened`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="timeout")
async def timeout_cmd(ctx, member: discord.Member = None, mins: int = 5, *, reason: str = ""):
    if not member or not ctx.guild:
        return
    try:
        await member.timeout(timedelta(minutes=mins), reason=reason)
        await ghost_reply(ctx, f"`timed out {mins}m`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="untimeout")
async def untimeout_cmd(ctx, member: discord.Member = None):
    if not member or not ctx.guild:
        return
    try:
        await member.timeout(None)
        await ghost_reply(ctx, "`untimeout`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="nickname", aliases=["setnick"])
async def nickname_cmd(ctx, member: discord.Member = None, *, name: str = None):
    if not member or not ctx.guild:
        return
    try:
        await member.edit(nick=name)
        await ghost_reply(ctx, "`nick set`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="role")
async def role_cmd(ctx, action: str = "", member: discord.Member = None, *, role: discord.Role = None):
    if not action or not member or not role or not ctx.guild:
        return await ghost_reply(ctx, "`usage: +role add|remove @u <role>`")
    try:
        if action == "add":
            await member.add_roles(role)
        elif action == "remove":
            await member.remove_roles(role)
        else:
            return await ghost_reply(ctx, "`action add|remove`")
        await ghost_reply(ctx, f"`{action}ed {role.name}`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

# ==================================================================
# STATUS
# ==================================================================
@bot.command(name="status")
async def status_cmd(ctx, status: str = "online"):
    statuses = {"online": discord.Status.online, "idle": discord.Status.idle,
                "dnd": discord.Status.dnd, "invisible": discord.Status.invisible}
    await bot.change_presence(status=statuses.get(status, discord.Status.online))
    await ghost_reply(ctx, f"`-> {status}`")

@bot.command(name="playing")
async def playing_cmd(ctx, *, text: str = "ghosttools"):
    await bot.change_presence(activity=discord.Game(name=text))
    await ghost_reply(ctx, "`ok`")

@bot.command(name="streaming")
async def streaming_cmd(ctx, *, text: str = "ghosttools"):
    await bot.change_presence(activity=discord.Streaming(name=text, url="https://twitch.tv/ghosttools"))
    await ghost_reply(ctx, "`ok`")

@bot.command(name="listening")
async def listening_cmd(ctx, *, text: str = "ghosttools"):
    await bot.change_presence(activity=discord.Activity(type=discord.ActivityType.listening, name=text))
    await ghost_reply(ctx, "`ok`")

@bot.command(name="watching")
async def watching_cmd(ctx, *, text: str = "ghosttools"):
    await bot.change_presence(activity=discord.Activity(type=discord.ActivityType.watching, name=text))
    await ghost_reply(ctx, "`ok`")

@bot.command(name="custom")
async def custom_cmd(ctx, *, text: str = "ghosttools"):
    await bot.change_presence(activity=discord.CustomActivity(name=text))
    await ghost_reply(ctx, "`ok`")

@bot.command(name="clearstatus")
async def clearstatus_cmd(ctx):
    await bot.change_presence(activity=None, status=discord.Status.online)
    await ghost_reply(ctx, "`cleared`")

# ==================================================================
# VOICE
# ==================================================================
@bot.command(name="joinvc")
async def joinvc_cmd(ctx, channel_id: int = 0):
    ch = bot.get_channel(channel_id) if channel_id else (ctx.author.voice.channel if ctx.author.voice else None)
    if not ch:
        return await ghost_reply(ctx, "`no channel`")
    try:
        await ch.connect()
        await ghost_reply(ctx, "`joined`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="leavevc")
async def leavevc_cmd(ctx):
    try:
        if ctx.voice_client:
            await ctx.voice_client.disconnect()
        await ghost_reply(ctx, "`left`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="vcspam")
async def vcspam_cmd(ctx, channel_id: int = 0):
    if not channel_id or not ctx.guild:
        return
    try:
        for _ in range(10):
            vc = await bot.get_channel(channel_id).connect()
            await vc.disconnect()
            await asyncio.sleep(0.4)
        await ghost_reply(ctx, "`done`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="deafen_self")
async def deafen_self_cmd(ctx):
    try:
        if ctx.voice_client:
            await ctx.voice_client.guild.change_voice_state(channel=ctx.voice_client.channel, self_deaf=True)
        await ghost_reply(ctx, "`ok`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="mute_self")
async def mute_self_cmd(ctx):
    try:
        if ctx.voice_client:
            await ctx.voice_client.guild.change_voice_state(channel=ctx.voice_client.channel, self_mute=True)
        await ghost_reply(ctx, "`ok`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

# ==================================================================
# GHOST
# ==================================================================
@bot.command(name="ghost")
async def ghost_cmd(ctx, seconds: int = 10):
    try:
        await ctx.message.delete()
    except Exception:
        pass
    cutoff = datetime.utcnow().timestamp() - seconds
    deleted = 0
    try:
        async for m in ctx.channel.history(limit=200):
            if m.author == bot.user and m.created_at.timestamp() > cutoff:
                try:
                    await m.delete()
                    deleted += 1
                except Exception:
                    pass
    except Exception:
        pass
    await safe_delete(await ctx.send(f"`ghosted {deleted}`"), 3)

@bot.command(name="ghostdm")
async def ghostdm_cmd(ctx, user: discord.User = None, *, msg: str = ""):
    if not user:
        return
    try:
        m = await user.send(msg)
        await asyncio.sleep(5)
        await m.delete()
    except Exception:
        pass
    try:
        await ctx.message.delete()
    except Exception:
        pass

@bot.command(name="clear")
async def clear_cmd(ctx):
    deleted = 0
    try:
        async for m in ctx.channel.history(limit=500):
            if m.author == bot.user:
                try:
                    await m.delete()
                    deleted += 1
                except Exception:
                    pass
    except Exception:
        pass
    await safe_delete(await ctx.send(f"`cleared {deleted}`"), 3)

@bot.command(name="ghostall")
async def ghostall_cmd(ctx, seconds: int = 10):
    total = 0
    if not ctx.guild:
        return
    cutoff = datetime.utcnow().timestamp() - seconds
    for ch in ctx.guild.text_channels:
        try:
            async for m in ch.history(limit=50):
                if m.author == bot.user and m.created_at.timestamp() > cutoff:
                    try:
                        await m.delete()
                        total += 1
                    except Exception:
                        pass
        except Exception:
            pass
    await ghost_reply(ctx, f"`ghosted {total}`")

@bot.command(name="ghostchannel")
async def ghostchannel_cmd(ctx):
    deleted = 0
    try:
        async for m in ctx.channel.history(limit=1000):
            if m.author == bot.user:
                try:
                    await m.delete()
                    deleted += 1
                except Exception:
                    pass
    except Exception:
        pass
    await safe_delete(await ctx.send(f"`ghosted {deleted}`"), 3)

# ==================================================================
# RECON
# ==================================================================
@bot.command(name="track")
async def track_cmd(ctx, user: discord.User = None):
    if not user:
        return await ghost_reply(ctx, "`usage: +track @user`")
    count = sum(1 for m in MESSAGE_LOG if m["author_id"] == user.id)
    e = make_embed(f"tracking {user}")
    e.add_field(name="messages seen", value=str(count), inline=True)
    await send_embed(ctx, e)

@bot.command(name="messages")
async def messages_cmd(ctx, user: discord.User = None):
    if not user:
        return await ghost_reply(ctx, "`usage: +messages @user`")
    msgs = [m for m in MESSAGE_LOG if m["author_id"] == user.id][-10:]
    if not msgs:
        return await ghost_reply(ctx, "`no cached`")
    txt = "\n".join(f"[{datetime.utcfromtimestamp(m['time']).strftime('%H:%M')}] {m['content'][:80]}" for m in msgs)
    await send_embed(ctx, make_embed(f"messages {user}", f"```{txt}```"))

@bot.command(name="lastseen")
async def lastseen_cmd(ctx, user: discord.User = None):
    if not user:
        return
    msgs = [m for m in MESSAGE_LOG if m["author_id"] == user.id]
    if not msgs:
        return await ghost_reply(ctx, "`never seen`")
    last = msgs[-1]
    await ghost_reply(ctx, f"`{datetime.utcfromtimestamp(last['time']).strftime('%Y-%m-%d %H:%M:%S')}` in <#{last['channel']}>")

@bot.command(name="search")
async def search_cmd(ctx, *, query: str = ""):
    if not query:
        return
    matches = [m for m in MESSAGE_LOG if query.lower() in m["content"].lower()][-15:]
    if not matches:
        return await ghost_reply(ctx, "`no matches`")
    txt = "\n".join(f"[{m['author']}] {m['content'][:100]}" for m in matches)
    await send_embed(ctx, make_embed(f"search: {query}", f"```{txt}```"))

@bot.command(name="history")
async def history_cmd(ctx, n: int = 20):
    msgs = [m for m in MESSAGE_LOG if m["channel"] == ctx.channel.id][-n:]
    if not msgs:
        return await ghost_reply(ctx, "`no cached`")
    txt = "\n".join(f"[{m['author']}] {m['content'][:100]}" for m in msgs)
    await send_embed(ctx, make_embed(f"last {len(msgs)}", f"```{txt[:3900]}```"))

@bot.command(name="typing")
async def typing_cmd(ctx):
    try:
        async with ctx.typing():
            await asyncio.sleep(5)
    except Exception:
        pass

@bot.command(name="afk")
async def afk_cmd(ctx, *, reason: str = "afk"):
    AFK_STATE["afk"] = True
    AFK_STATE["reason"] = reason
    await ghost_reply(ctx, f"`afk: {reason}`")

@bot.command(name="back")
async def back_cmd(ctx):
    AFK_STATE["afk"] = False
    await ghost_reply(ctx, "`back`")

# ==================================================================
# MASS
# ==================================================================
@bot.command(name="massban")
async def massban_cmd(ctx, *, ids: str = ""):
    if not ids or not ctx.guild:
        return
    id_list = [i.strip() for i in ids.split(",") if i.strip().isdigit()]
    n = 0
    for uid in id_list:
        try:
            u = await bot.fetch_user(int(uid))
            await ctx.guild.ban(u)
            n += 1
            await asyncio.sleep(0.5)
        except Exception:
            pass
    await ghost_reply(ctx, f"`banned {n}`")

@bot.command(name="masskick")
async def masskick_cmd(ctx):
    if not ctx.guild:
        return
    try:
        await ctx.message.delete()
    except Exception:
        pass
    n = 0
    for m in list(ctx.guild.members):
        if m.bot or m.id == bot.user.id:
            continue
        try:
            await m.kick(reason="ghosttools")
            n += 1
            await asyncio.sleep(0.35)
        except Exception:
            pass
    await safe_delete(await ctx.send(f"`kicked {n}`"), 3)

@bot.command(name="massrole")
async def massrole_cmd(ctx, role: discord.Role = None, action: str = "add"):
    if not role or not ctx.guild:
        return
    try:
        await ctx.message.delete()
    except Exception:
        pass
    n = 0
    for m in ctx.guild.members:
        try:
            if action == "add":
                await m.add_roles(role)
            else:
                await m.remove_roles(role)
            n += 1
            await asyncio.sleep(0.2)
        except Exception:
            pass
    await safe_delete(await ctx.send(f"`{action}ed to {n}`"), 3)

@bot.command(name="massinvite")
async def massinvite_cmd(ctx):
    if not ctx.guild:
        return
    try:
        ch = ctx.guild.text_channels[0] if ctx.guild.text_channels else None
        if not ch:
            return await ghost_reply(ctx, "`no channels`")
        inv = await ch.create_invite(max_age=86400, max_uses=0)
        await ghost_reply(ctx, f"`{inv.url}`", 30)
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

# ==================================================================
# FUN
# ==================================================================
@bot.command(name="8ball")
async def eightball_cmd(ctx, *, q: str = ""):
    answers = ["yes","no","maybe","definitely","absolutely not","ask later",
               "signs point to yes","my sources say no","without a doubt","don't count on it"]
    await ghost_reply(ctx, f"`{random.choice(answers)}`")

@bot.command(name="coinflip")
async def coinflip_cmd(ctx):
    await ghost_reply(ctx, f"`{random.choice(['heads','tails'])}`")

@bot.command(name="dice")
async def dice_cmd(ctx, sides: int = 6):
    await ghost_reply(ctx, f"`{random.randint(1, sides)}`")

@bot.command(name="meme")
async def meme_cmd(ctx):
    try:
        async with aiohttp.ClientSession() as s:
            async with s.get("https://meme-api.com/gimme") as r:
                data = await r.json()
        await ctx.send(data.get("url", "no meme"))
        try:
            await ctx.message.delete()
        except Exception:
            pass
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="joke")
async def joke_cmd(ctx):
    try:
        async with aiohttp.ClientSession() as s:
            async with s.get("https://official-joke-api.appspot.com/random_joke") as r:
                data = await r.json()
        await ghost_reply(ctx, f"`{data.get('setup','')} — {data.get('punchline','')}`", 15)
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="rate")
async def rate_cmd(ctx, *, thing: str = ""):
    if not thing:
        return
    seed = sum(ord(c) for c in thing)
    random.seed(seed)
    score = random.randint(1, 10)
    random.seed()
    await ghost_reply(ctx, f"`{thing}: {score}/10`")

@bot.command(name="gay")
async def gay_cmd(ctx, user: discord.User = None):
    user = user or ctx.author
    await ghost_reply(ctx, f"`{user.name} is {user.id % 101}% gay`")

@bot.command(name="ship")
async def ship_cmd(ctx, u1: discord.User = None, u2: discord.User = None):
    if not u1 or not u2:
        return await ghost_reply(ctx, "`usage: +ship @a @b`")
    seed = (u1.id + u2.id) % 101
    bar = "█" * (seed // 5) + "░" * (20 - seed // 5)
    await ghost_reply(ctx, f"`{u1.name} + {u2.name}\n{bar} {seed}%`", 15)

@bot.command(name="roast")
async def roast_cmd(ctx, user: discord.User = None):
    user = user or ctx.author
    roasts = [
        f"{user.name} is the reason the gene pool needs a lifeguard.",
        f"{user.name} brings everyone joy — when they leave the room.",
        f"{user.name} isn't stupid, they just have bad luck when thinking.",
        f"{user.name}'s brain has too many tabs open and one is playing music.",
    ]
    await ghost_reply(ctx, f"`{random.choice(roasts)}`", 10)

@bot.command(name="compliment")
async def compliment_cmd(ctx, user: discord.User = None):
    user = user or ctx.author
    comps = [
        f"{user.name} makes the room brighter.",
        f"{user.name} is genuinely kind.",
        f"{user.name} has great taste.",
        f"{user.name} is smarter than they let on.",
    ]
    await ghost_reply(ctx, f"`{random.choice(comps)}`", 10)

# ==================================================================
# SERVER MGMT
# ==================================================================
@bot.command(name="createchannel")
async def createchannel_cmd(ctx, *, name: str = "new-channel"):
    if not ctx.guild:
        return
    try:
        ch = await ctx.guild.create_text_channel(name=name)
        await ghost_reply(ctx, f"`created #{ch.name}`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="deletechannel")
async def deletechannel_cmd(ctx, channel_id: int = 0):
    if not ctx.guild or not channel_id:
        return
    try:
        ch = ctx.guild.get_channel(channel_id)
        if ch:
            await ch.delete()
            await ghost_reply(ctx, "`deleted`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="createrole")
async def createrole_cmd(ctx, *, name: str = "new-role"):
    if not ctx.guild:
        return
    try:
        r = await ctx.guild.create_role(name=name, color=discord.Color.random())
        await ghost_reply(ctx, f"`created {r.name}`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="deleterole")
async def deleterole_cmd(ctx, role: discord.Role = None):
    if not ctx.guild or not role:
        return
    try:
        await role.delete()
        await ghost_reply(ctx, "`deleted`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="createemoji")
async def createemoji_cmd(ctx, name: str = "emoji", url: str = None):
    if not ctx.guild or not url:
        return await ghost_reply(ctx, "`usage: +createemoji <n> <url>`")
    try:
        async with aiohttp.ClientSession() as s:
            async with s.get(url) as r:
                data = await r.read()
        em = await ctx.guild.create_custom_emoji(name=name, image=data)
        await ghost_reply(ctx, f"`{em}`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="webhooks")
async def webhooks_cmd(ctx):
    if not ctx.guild:
        return
    try:
        hooks = await ctx.channel.webhooks()
        txt = "\n".join(f"`{h.name}` — {h.id}" for h in hooks[:15])
        await send_embed(ctx, make_embed(f"webhooks ({len(hooks)})", txt or "none"))
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="createwebhook")
async def createwebhook_cmd(ctx, *, name: str = "ghosttools"):
    if not ctx.guild:
        return
    try:
        h = await ctx.channel.create_webhook(name=name)
        await ghost_reply(ctx, f"`{h.url}`", 15)
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

@bot.command(name="delwebhook")
async def delwebhook_cmd(ctx, webhook_id: int = 0):
    if not ctx.guild or not webhook_id:
        return
    try:
        hooks = await ctx.channel.webhooks()
        for h in hooks:
            if h.id == webhook_id:
                await h.delete()
                return await ghost_reply(ctx, "`deleted`")
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

# ==================================================================
# BACKUP
# ==================================================================
@bot.command(name="backup")
async def backup_cmd(ctx):
    if not ctx.guild:
        return
    g = ctx.guild
    data = {
        "name": g.name, "id": g.id,
        "roles": [{"name": r.name, "color": r.color.value, "permissions": r.permissions.value} for r in g.roles],
        "channels": [{"name": c.name, "type": str(c.type), "position": c.position} for c in g.channels],
        "emojis": [{"name": e.name, "url": str(e.url)} for e in g.emojis],
    }
    f = io.BytesIO(json.dumps(data, indent=2).encode())
    await send_embed(ctx, make_embed("backup complete", f"roles: {len(data['roles'])}\nchannels: {len(data['channels'])}\nemojis: {len(data['emojis'])}"))
    await ctx.send(file=discord.File(f, filename=f"{g.name}-backup.json"))

@bot.command(name="backup_roles")
async def backup_roles_cmd(ctx):
    if not ctx.guild:
        return
    data = [{"name": r.name, "color": r.color.value, "permissions": r.permissions.value} for r in ctx.guild.roles]
    f = io.BytesIO(json.dumps(data, indent=2).encode())
    await ctx.send(file=discord.File(f, filename="roles.json"))

@bot.command(name="backup_channels")
async def backup_channels_cmd(ctx):
    if not ctx.guild:
        return
    data = [{"name": c.name, "type": str(c.type), "position": c.position} for c in ctx.guild.channels]
    f = io.BytesIO(json.dumps(data, indent=2).encode())
    await ctx.send(file=discord.File(f, filename="channels.json"))

@bot.command(name="backup_emojis")
async def backup_emojis_cmd(ctx):
    if not ctx.guild:
        return
    data = [{"name": e.name, "url": str(e.url), "animated": e.animated} for e in ctx.guild.emojis]
    f = io.BytesIO(json.dumps(data, indent=2).encode())
    await ctx.send(file=discord.File(f, filename="emojis.json"))

@bot.command(name="dumpmembers")
async def dumpmembers_cmd(ctx):
    if not ctx.guild:
        return
    data = [{"id": m.id, "name": str(m), "bot": m.bot, "roles": [r.name for r in m.roles]} for m in ctx.guild.members]
    f = io.BytesIO(json.dumps(data, indent=2).encode())
    await send_embed(ctx, make_embed(f"members ({len(data)})", "dumping"))
    await ctx.send(file=discord.File(f, filename="members.json"))

@bot.command(name="dumpinvites")
async def dumpinvites_cmd(ctx):
    if not ctx.guild:
        return
    try:
        invites = await ctx.guild.invites()
        data = [{"code": i.code, "uses": i.uses, "channel": i.channel.name if i.channel else None} for i in invites]
        f = io.BytesIO(json.dumps(data, indent=2).encode())
        await ctx.send(file=discord.File(f, filename="invites.json"))
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`")

# ==================================================================
# AUTOMATION
# ==================================================================
@bot.command(name="autoreply")
async def autoreply_cmd(ctx, *, args: str = ""):
    if "|" not in args:
        return await ghost_reply(ctx, "`usage: +autoreply <trigger> | <response>`")
    trigger, response = args.split("|", 1)
    AUTOREPLIES[trigger.strip().lower()] = response.strip()
    await ghost_reply(ctx, f"`+ {trigger.strip()}`")

@bot.command(name="clearautoreply")
async def clearautoreply_cmd(ctx):
    AUTOREPLIES.clear()
    await ghost_reply(ctx, "`cleared`")

# ==================================================================
# OWNER
# ==================================================================
@bot.command(name="eval")
async def eval_cmd(ctx, *, code: str = ""):
    if not code:
        return
    try:
        result = eval(code, {"bot": bot, "ctx": ctx, "discord": discord, "asyncio": asyncio})
        if asyncio.iscoroutine(result):
            result = await result
        await ghost_reply(ctx, f"`{str(result)[:1900]}`", 10)
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`", 10)

@bot.command(name="exec")
async def exec_cmd(ctx, *, cmd: str = ""):
    if not cmd:
        return
    try:
        out = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=30)
        txt = (out.stdout or "") + (out.stderr or "")
        await ghost_reply(ctx, f"`{txt[:1900] or '(no output)'}`", 10)
    except Exception as e:
        await ghost_reply(ctx, f"`err: {e}`", 10)

@bot.command(name="restart")
async def restart_cmd(ctx):
    await ghost_reply(ctx, "`restarting…`", 2)
    os.execv(sys.executable, [sys.executable] + sys.argv)

@bot.command(name="shutdown")
async def shutdown_cmd(ctx):
    await ghost_reply(ctx, "`bye`", 2)
    await bot.close()

@bot.command(name="log")
async def log_cmd(ctx, n: int = 20):
    txt = "\n".join(f"[{m['author']}] {m['content'][:80]}" for m in MESSAGE_LOG[-n:])
    await send_embed(ctx, make_embed(f"log ({min(n,len(MESSAGE_LOG))})", f"```{txt[:3900]}```"))

# ==================================================================
# ENTRY
# ==================================================================
if __name__ == "__main__":
    tok = load_token()
    if not tok:
        try:
            tok = input("enter discord token: ").strip()
            save_token(tok)
        except (KeyboardInterrupt, EOFError):
            print("cancelled.")
            raise SystemExit
    try:
        bot.run(tok)
    except Exception as e:
        print(f"selfbot crashed: {e}")
        input("press enter...")