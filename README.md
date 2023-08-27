# Discord-SMTP Server
A simple relay that accepts SMTP messages and forwards them to a Discord webhook.

## Usage

### Local

```
env DISCORD_WEBHOOK_URI=xxxxxxxxxxxx SMTP_USERNAME=username SMTP_PASSWORD=password go run main.go
```

### Docker

#### Run

```
docker run -t discord-smtp \
           -e PORT=25 \
           -e DISCORD_WEBHOOK_URI=xxxxxxxxxxxx \
           -e SMTP_USERNAME=username \
           -e SMTP_PASSWORD=password \
           nullcosmos/discord-smtp-server
```

#### Compose

```
discord-smtp:
  image: nullcosmos/discord-smtp-server
  container_name: discord-smtp
  env:
    - PORT=25
    - DISCORD_WEBHOOK_URI=xxxxxxxxxxxx
    - SMTP_USERNAME=username
    - SMTP_PASSWORD=password
  restart: always
```

#### Testing

```
$ telnet localhost 1025
```

```
EHLO localhost
AUTH PLAIN
AHVzZXJuYW1lAHBhc3N3b3Jk
MAIL FROM:<test@test.com>
RCPT TO:<smtp@alert.karenplankton>
DATA
Hey
.
```

## Features

* SMTP Authentication