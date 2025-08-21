package smtp

import (
  "errors"
  "io"

  "github.com/emersion/go-smtp"
  "github.com/justinnamilee/discord-smtp-relay/discord"
)

type Backend struct {
  discord *discord.Session
  username, password string
}

type Session struct {
  backend *Backend
  webhook, from string
}

func New(discord *discord.Session, username, password string) (*Backend, error) {
  return &Backend{
    discord: discord,
    username: username,
    password: password,
  }, nil
}


// --- Important Discord Handling Part ---
func (s *Session) Data(r io.Reader) error {
  return s.backend.discord.Message(r)
}


// generic SMTP wrappers for stuff

func (b *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
  if username != b.username || password != b.password {
    return nil, errors.New("Invalid username or password")
  }

  return &Session{
    backend: b,
  }, nil
}

func (b *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
  return nil, smtp.ErrAuthRequired
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
  s.from = from
  return nil
}

func (s *Session) Rcpt(to string) error {
  return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
  return nil
}
