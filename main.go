package main

import (
  "log"
  "os"
  "strconv"
  "strings"
  "time"
  "net/url"

  gosmtp "github.com/emersion/go-smtp"
  "github.com/justinnamilee/discord-smtp-relay/discord"
  "github.com/justinnamilee/discord-smtp-relay/smtp"
)

func getEnv(key string, required bool) string {
  val := strings.TrimSpace(os.Getenv(key))
  if required && val == "" {
    log.Fatalf("Environment variable %s is required but not set or empty", key)
  }
  return val
}

func main() {
  webhookURL := getEnv("WEBHOOK", true)
  templatePath := getEnv("TEMPLATE", true)
  username := getEnv("USERNAME", true)
  password := getEnv("PASSWORD", true)

  port := getEnv("PORT", false)
  if port == "" {
    port = "1025"
  }
  host := getEnv("HOST", false)
  if host == "" {
    host = "localhost"
  }

  if _, err := strconv.Atoi(port); err != nil {
    log.Fatalf("Invalid PORT %q: %v", port, err)
  }

  if u, err := url.Parse(webhookURL); err != nil || u.Scheme == "" || u.Host == "" {
    log.Fatalf("WEBHOOK is not a valid URL: %q", webhookURL)
  }

  if info, err := os.Stat(templatePath); err != nil || info.IsDir() {
    log.Fatalf("TEMPLATE file %q does not exist or is a directory", templatePath)
  }

  discordSess, err := discord.NewSession(webhookURL, templatePath)
  if err != nil {
    log.Fatal(err)
  }

  backend, err := smtp.NewBackend(discordSess, username, password)
  if err != nil {
    log.Fatal(err)
  }

  server := gosmtp.NewServer(backend)
  server.Addr = ":" + port
  server.Domain = host
  server.ReadTimeout = 10 * time.Second
  server.WriteTimeout = 10 * time.Second
  server.MaxMessageBytes = 1024 * 1024
  server.MaxRecipients = 50
  server.AllowInsecureAuth = true

  log.Println("Starting server at", server.Addr)
  if err := server.ListenAndServe(); err != nil {
    log.Fatal(err)
  }
}
