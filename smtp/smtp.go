package smtp

import (
	"errors"
	"io"

	"github.com/emersion/go-smtp"
	"github.com/nullcosmos/discord-smtp-server/discord"
)

type Backend struct {
	discord       *discord.Session
	username      string
	password      string
}

func NewBackend(discord *discord.Session, username, password string) (*Backend, error) {
	return &Backend{
		discord:       discord,
		username:      username,
		password:      password,
	}, nil
}

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

type Session struct {
	backend *Backend
	webhook string
	from    string
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	return nil
}

func (s *Session) Data(r io.Reader) error {
	return s.backend.discord.Send(r)
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
