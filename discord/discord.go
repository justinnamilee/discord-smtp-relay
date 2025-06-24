package discord

import (
  "bytes"
  "encoding/json"
  "io"
  "io/ioutil"
  "net/http"
  "net/mail"
  "text/template"
)

type Session struct {
  webhook, template string
  cached *template.Template
}

type TemplateParams struct {
  Date, From, To, Subject, Body string
}

func New(discordWebhookUri, discordTemplatePath string) (*Session, error) {
  raw, err := ioutil.ReadFile(discordTemplatePath)
  if err != nil {
    return nil, err
  }

  parsed, err := template.New("message").Parse(string(raw))
  if err != nil {
    return nil, err
  }

  return &Session{
    webhook: discordWebhookUri,
    template: discordTemplatePath,
    cached: parsed,
  }, nil
}

func (s *Session) message(r io.Reader) error {
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

func (s *Session) sendToDiscord(content string) error {
  payload, err := json.Marshal(map[string]string{
    "content": content,
  })
  if err != nil {
    return err
  }

  resp, err := http.Post(s.webhook, "application/json", bytes.NewBuffer(payload))
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  return nil
}
