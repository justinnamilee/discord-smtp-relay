package discord

import (
  "bytes"
  "io"
  "io/ioutil"
  "net/http"
  "net/mail"
  "strings"
  "text/template"
  "time"
)

type Session struct {
  webhook, template string
  cached *template.Template
}

type TemplateParams struct {
  Date, From, To, Subject, Body, DateGet string
}

func New(discordWebhookURI, discordTemplatePath string) (*Session, error) {
  raw, err := ioutil.ReadFile(discordTemplatePath)
  if err != nil {
    return nil, err
  }

  parsed, err := template.New("embed").Parse(string(raw))
  if err != nil {
    return nil, err
  }

  return &Session{
    webhook: discordWebhookURI,
    template: discordTemplatePath,
    cached: parsed,
  }, nil
}

func (s *Session) Message(r io.Reader) error {
  msg, err := s.parseTemplate(r)
  if err != nil {
    return err
  }

  return s.sendToDiscord(msg)
}

func (s *Session) parseTemplate(r io.Reader) (string, error) {
  m, err := mail.ReadMessage(r)
  if err != nil {
    return "", err
  }

  bodyBytes, err := ioutil.ReadAll(m.Body)
  if err != nil {
    return "", err
  }

  params := TemplateParams{
    Date: m.Header.Get("Date"),
    DateGet: time.Now().Format(time.RFC1123Z),
    From: m.Header.Get("From"),
    To: m.Header.Get("To"),
    Subject: m.Header.Get("Subject"),
    Body: string(bodyBytes),
  }

  buf := new(bytes.Buffer)
  if err := s.cached.Execute(buf, params); err != nil {
    return "", err
  }

  return buf.String(), nil
}

func (s *Session) sendToDiscord(payload string) error {
  req, err := http.NewRequest("POST", s.webhook, strings.NewReader(payload))
  if err != nil {
    return err
  }

  req.Header.Set("Content-Type", "application/json")
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  return nil
}
