package smtp

import (
	"io"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/justinnamilee/discord-smtp-relay/discord"
)

type Backend struct {
	discord            *discord.Session
	username, password string
}

type Session struct {
	backend       *Backend
	authenticated bool
}

func New(discord *discord.Session, username, password string) (*Backend, error) {
	return &Backend{
		discord:  discord,
		username: username,
		password: password,
	}, nil
}

func (b *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{backend: b}, nil
}

func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

func (s *Session) Auth(mechanism string) (sasl.Server, error) {
	if mechanism != sasl.Plain {
		return nil, smtp.ErrAuthUnknownMechanism
	}

	return sasl.NewPlainServer(func(identity, username, password string) error {
		if identity != "" && identity != username {
			return smtp.ErrAuthFailed
		}
		if s.backend.username == "" || username != s.backend.username || password != s.backend.password {
			return smtp.ErrAuthFailed
		}

		s.authenticated = true
		return nil
	}), nil
}

func (s *Session) Mail(_ string, _ *smtp.MailOptions) error {
	if !s.authenticated && (s.backend.username != "" || s.backend.password != "") {
		return smtp.ErrAuthRequired
	}
	return nil
}

func (s *Session) Rcpt(_ string, _ *smtp.RcptOptions) error {
	return nil
}

func (s *Session) Data(r io.Reader) error {
	return s.backend.discord.Message(r)
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}
