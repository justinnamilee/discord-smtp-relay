# 📬 discord-smtp-relay

A tiny Go SMTP server that turns your incoming emails into **fancy Discord embeds**...Because plain text is so last decade. 🎉

[![Go ≥1.15](https://img.shields.io/badge/go-%3E%3D1.15-blue.svg)](https://golang.org) [![MIT License](https://img.shields.io/github/license/justinnamilee/discord-smtp-relay.svg)](LICENSE)

---

## 🚀 Features

* **SMTP listener** on `0.0.0.0:1025` (or whatever you choose)
* **Discord embeds** out of the box (*no more wall-of-text “content” messages*)
* **Env-driven** config — perfect for PM2 or Docker
* **Basic auth** support (SMTP `AUTH PLAIN`)

---

## 🛠️ Getting Started

### Prerequisites

* Go 1.15+ (modules enabled)
* A Discord **Incoming Webhook URL**
* Common sense

### Clone & Build

```bash
git clone https://github.com/justinnamilee/discord-smtp-relay.git
cd discord-smtp-relay
go build -o discord-smtp-relay main.go

# ...
```

### Configuration

Set these environment variables before 🚀 launch:

| Variable   | Required | Default     | Description                       |
| :--------- | :------: | :---------- | :-------------------------------- |
| `WEBHOOK`  |    ✅    | —           | Your Discord incoming webhook URL |
| `TEMPLATE` |    ✅    | —           | Path to your JSON embed template  |
| `USERNAME` |    ❌    | `discord`   | SMTP AUTH username                |
| `PASSWORD` |    ❌    | `discord`   | SMTP AUTH password                |
| `HOST`     |    ❌    | `0.0.0.0`   | Listen interface                  |
| `PORT`     |    ❌    | `1025`      | Listen port                       |
| `DOMAIN`   |    ❌    | `localhost` | EHLO/HELO domain                  |
| `READ`     |    ❌    | `10`        | Read timeout (seconds)            |
| `WRITE`    |    ❌    | `10`        | Write timeout (seconds)           |
| `SIZE`     |    ❌    | `1024`      | Max message size (KB)             |

### Run it!

```bash
# ...

env \
  WEBHOOK="https://discord.com/api/webhooks/..." \
  TEMPLATE="etc/template.example.json" \
  ./discord-smtp-relay
```

---

## 🔧 PM2 Installion or Docker Usage


### PM2 Installation

```bash
# ...assuming you built in last step

cp ecosystem.config.example.js ecosystem.config.js
# edit the ecosystem file as desired (set WEBHOOK, TEMPLATE)
pm2 start ecosystem.config.js
pm2 save
```

### Docker Usage

_**No need to download the repo or anything!**_

```bash
wget 'https://github.com/justinnamilee/discord-smtp-relay/blob/main/compose.example.yml' -O compose.yml
# edit the compose file as desired (set WEBHOOK, TEMPLATE)
docker compose up -d
```

---

## 📤 Send an Email

Use your favorite SMTP client, or good old `telnet`:

```bash
$ telnet localhost 1025

EHLO localhost
AUTH PLAIN AGRpc2NvcmQAZGlzY29yZA==
MAIL FROM:<you@example.com>
RCPT TO:<bot@example.com>
DATA
From: Your Pretty Name Here
To: Bots Pretty Name Here
Subject: Testing 1, 2, 3
Date: Tue, 24 Jun 2025 21:37:13 +0000

Hello Discord! This is an embedded email.
.
```

Watch your Discord channel light up with a nice embed instead of a boring text dump. ✨

**Also**, *please*, **PLEASE** make sure you update the default credentials.  To generate a new `AUTH PLAIN` line for testing, use:

```bash
printf '\0username\0password' | base64
```

---

## 🧩 Template Example

This is the file that's at `etc/template.example.json`, the rest of that directory is in `.gitignore` so add as many new ones as you want...  For more details on this please see [text/template](https://pkg.go.dev/text/template) and [Discord Embeds](https://discordjs.guide/popular-topics/embeds.html).  Feel free to submit ideas for new fields to support or new default templates to provide via issues or PR or w/e people normally do.

<details>
<summary>Click to expand</summary>

```json
{
  "embeds": [
    {
      "title": "📨 {{ .Subject }}",
      "description": {{ printf "%q" .Body }},
      "color": 5814783,
      "fields": [
        {
          "name": "From",
          "value": {{ printf "%q" .From }},
          "inline": true
        },
        {
          "name": "To",
          "value": {{ printf "%q" .To }},
          "inline": true
        },
        {
          "name": "Sent",
          "value": {{ printf "%q" .Date }},
          "inline": false
        }
      ],
      "footer": {
        "text": "Received on {{ .DateGet }}"
      }
    }
  ]
}
```

</details>

The fields supported right now are:
- Subject
- From
- To
- Date (*generated from SMTP sender*)
- DateGet (*generated when we get the message*)

---

## 📄 License

This project is licensed under the **MIT** License – see [LICENSE](LICENSE) for details.

---

> **Pro tip #1:** Markdown in embeds is supported!

> **Pro tip #2:** Pull requests, feature requests, and emoji suggestions are very welcome. 🎈

> **Pro tip #3:** Chew your food before swallowing.

> **Pro tip #4:** If you use PM2, there's an example ecosystem.config.js file for you.
