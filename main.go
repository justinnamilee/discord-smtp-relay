package main

import (
  "log"
  "net"
  "os"
  "strconv"
  "strings"
  "time"
  "net/url"

  gosmtp "github.com/emersion/go-smtp"
  "github.com/justinnamilee/discord-smtp-relay/discord"
  "github.com/justinnamilee/discord-smtp-relay/smtp"
)

func getEnv(key string, required bool, def string) string {
  v := strings.TrimSpace(os.Getenv(key))
  if v == "" {
    if required {
      log.Fatalf("Environment variable %s is required but not set or empty", key)
    }
    return def
  }
  return v
}

func main() {
  log.Print("Starting up...")

  //--- core settings
  webhookURL := getEnv("WEBHOOK",  true,  "")
  templatePath := getEnv("TEMPLATE", true,  "")
  username := getEnv("USERNAME", false, "discord")
  password := getEnv("PASSWORD", false, "discord")

  //--- bind address
  bindHost := getEnv("HOST", false, "0.0.0.0")
  bindPort := getEnv("PORT", false, "1025")
  if _, err := strconv.Atoi(bindPort); err != nil {
    log.Fatalf("PORT must be numeric: %v", err)
  }
  listenAddr := net.JoinHostPort(bindHost, bindPort)

  //--- SMTP EHLO/HELO banner name
  smtpDomain := getEnv("DOMAIN", false, "localhost")
  if u, err := url.Parse("smtp://" + smtpDomain); err != nil || u.Hostname() == "" {
    log.Fatalf("DOMAIN isn’t a valid hostname: %q", smtpDomain)
  }

  readEnv := getEnv("READ", false, "10")
  readSecs, err := strconv.Atoi(readEnv)
  if err != nil {
    log.Fatalf("Invalid READ value %q: %v", readEnv, err)
  }

  writeEnv := getEnv("WRITE", false, "10")
  writeSecs, err := strconv.Atoi(writeEnv)
  if err != nil {
    log.Fatalf("Invalid WRITE value %q: %v", writeEnv, err)
  }

  sizeEnv := getEnv("SIZE", false, "1024")
  sizeKB, err := strconv.Atoi(sizeEnv)
  if err != nil {
    log.Fatalf("Invalid SIZE value %q: %v", sizeEnv, err)
  }

  //--- build Discord relay & SMTP backend
  discordSess, err := discord.New(webhookURL, templatePath)
  if err != nil {
    log.Fatal(err)
  }
  backend, err := smtp.New(discordSess, username, password)
  if err != nil {
    log.Fatal(err)
  }

  //--- configure and start SMTP server
  server := gosmtp.NewServer(backend)

  server.Addr = listenAddr
  server.Domain = smtpDomain
  server.ReadTimeout = time.Duration(readSecs) * time.Second
  server.WriteTimeout = time.Duration(writeSecs) * time.Second
  server.MaxMessageBytes = sizeKB * 1024
  server.MaxRecipients = 1
  server.AllowInsecureAuth = true

  log.Printf(
    "SMTP listening on %s, advertising as %s (read=%ds write=%ds size=%dKB)",
    listenAddr, smtpDomain, readSecs, writeSecs, sizeKB,
  )
  log.Fatal(server.ListenAndServe())
}
