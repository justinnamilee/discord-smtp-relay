# Discord SMTP Relay

[![Go Reference](https://pkg.go.dev/badge/github.com/justinnamilee/discord-smtp-relay.svg)](https://pkg.go.dev/github.com/justinnamilee/discord-smtp-relay)

A minimal SMTP server that forwards incoming emails into a Discord channel via webhook, using customizable Go templates.

---

## Features

* 📨 Receive email over SMTP
* 🔗 Relay parsed email (from, to, subject, body) to Discord
* 🛠️ Fully configurable via environment variables
* 📄 Customizable message templates with Go’s `text/template`

---

## Prerequisites

* **Go** 1.14+
* A **Discord Incoming Webhook URL**

---

## Installation

```bash
git clone https://github.com/justinnamilee/discord-smtp-relay.git
cd discord-smtp-relay

# Build the binary
go build -o discord-smtp-relay
```

---

## Configuration

All settings are controlled via environment variables. Required ones are noted accordingly.

| Name       | Default     | Required | Description                                        |
| ---------- | ----------- | -------- | -------------------------------------------------- |
| `WEBHOOK`  | —           | Yes      | Discord incoming webhook URL                       |
| `TEMPLATE` | —           | Yes      | Path to your Go `text/template` file for messaging |
| `USERNAME` | `discord`   | No       | SMTP AUTH username                                 |
| `PASSWORD` | `discord`   | No       | SMTP AUTH password                                 |
| `HOST`     | `0.0.0.0`   | No       | Host/interface to bind to                          |
| `PORT`     | `1025`      | No       | Port to listen on                                  |
| `DOMAIN`   | `localhost` | No       | EHLO/HELO banner domain                            |
| `READ`     | `10`        | No       | Read timeout (seconds)                             |
| `WRITE`    | `10`        | No       | Write timeout (seconds)                            |
| `SIZE`     | `1024`      | No       | Max message size (in KB)                           |

### Template file

Your `TEMPLATE` is a Go `text/template` that gets these fields:

* `.From`
* `.To`
* `.Subject`
* `.Body`

Example snippet (`template.tmpl`):

```gotemplate
**New email received!**

**From:** {{ .From }}
**To:** {{ .To }}
**Subject:** {{ .Subject }}

{{ .Body }}
```

---

## Usage

1. **Export your settings:**

   ```bash
   export WEBHOOK="https://discord.com/api/webhooks/…"
   export TEMPLATE="./template.tmpl"
   # (optional) other env vars…
   ```

2. **Run the server:**

   ```bash
   ./discord-smtp-relay
   ```

   You should see:

   ```
   SMTP listening on 0.0.0.0:1025, advertising as localhost (read=10s write=10s size=1024KB)
   ```

3. **Send a test email:**

   ```bash
   echo -e "Subject: Hello Discord\n\nThis is a test!" \
     | sendmail -S localhost:1025 -au discord -ap discord you@example.com
   ```

   You’ll get a nicely formatted message in your Discord channel.

---

## Docker

Here’s a quick Docker setup if you’d rather run it in a container:

```dockerfile
# Dockerfile
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o discord-smtp-relay

FROM alpine
COPY --from=builder /app/discord-smtp-relay /usr/local/bin/
ENTRYPOINT ["discord-smtp-relay"]
```

Build & run:

```bash
docker build -t discord-smtp-relay .
docker run -e WEBHOOK="…" -e TEMPLATE="/path/to/template.tmpl" -p 1025:1025 discord-smtp-relay
```

---

## Contributing

Contributions, issues, and feature requests are welcome! Please open an issue or PR on GitHub.

---

## License

This project is licensed under the MIT License © 2021 Kyle Lucas-Rodriguez, 2025 Justin Lee.
See the full text in [LICENSE](./LICENSE).
