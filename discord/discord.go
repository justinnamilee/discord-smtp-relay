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
	webhook string
}

type TemplateParams struct {
	Date, From, To, Subject, Body string
}

func NewSession(discordWebhookUri string) (*Session, error) {
	return &Session{
		webhook: discordWebhookUri,
	}, nil
}

func (s *Session) Send(r io.Reader) error {

	msg, err := s.ParseTemplate(r)
	if err != nil {
		return err
	}

	err = s.SendToDiscord(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) ParseTemplate(r io.Reader) (string, error) {
	m, err := mail.ReadMessage(r)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		return "", err
	}

	header := m.Header
	ttParams := TemplateParams{
		header.Get("Date"),
		header.Get("From"),
		header.Get("To"),
		header.Get("Subject"),
		string(body),
	}

	// Template
	tdat, err := ioutil.ReadFile("/message.tt")
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("message").Parse(string(tdat)))

	buf := new(bytes.Buffer)
	err = t.Execute(buf, ttParams)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *Session) SendToDiscord(content string) error {
	reqBody, err := json.Marshal(
		map[string]string{
			"content": content,
		},
	)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		s.webhook,
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}