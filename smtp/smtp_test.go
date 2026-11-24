package smtp_test

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/smtp"
	"net/textproto"
	"strings"
	"testing"
	"time"

	gosmtp "github.com/emersion/go-smtp"
	"github.com/justinnamilee/discord-smtp-relay/discord"
	relaysmtp "github.com/justinnamilee/discord-smtp-relay/smtp"
)

func TestRelayMail(t *testing.T) {
	const (
		from    = "sender@example.com"
		to      = "relay@example.com"
		subject = "Test alert"
		body    = "The relay received this message."
		message = "From: " + from + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + body + "\r\n"
	)

	tests := []struct {
		name      string
		anonymous bool
		auth      smtp.Auth
		wantCode  int
	}{
		{name: "anonymous", anonymous: true},
		{name: "authenticated", auth: smtp.PlainAuth("", "relay", "secret", "127.0.0.1")},
		{name: "missing credentials", wantCode: 502},
		{name: "invalid password", auth: smtp.PlainAuth("", "relay", "wrong", "127.0.0.1"), wantCode: 535},
		{name: "invalid identity", auth: smtp.PlainAuth("other", "relay", "secret", "127.0.0.1"), wantCode: 535},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type payload struct {
				Embeds []struct {
					Title       string
					Description string
					Fields      []struct{ Name, Value string }
				}
			}
			requests := make(chan payload, 1)
			webhook := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/webhook" || r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("webhook request = %s %s (%s)", r.Method, r.URL.Path, r.Header.Get("Content-Type"))
				}

				var p payload
				if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
					t.Errorf("decode webhook payload: %v", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				select {
				case requests <- p:
				default:
					t.Error("unexpected additional webhook request")
				}
				w.WriteHeader(http.StatusNoContent)
			}))
			defer webhook.Close()

			client := newRelayClient(t, webhook.URL+"/webhook", tt.anonymous)
			switch {
			case tt.wantCode != 0 && tt.auth != nil:
				checkSMTPCode(t, client.Auth(tt.auth), tt.wantCode)
			case tt.wantCode != 0:
				checkSMTPCode(t, client.Mail(from), tt.wantCode)
			default:
				if tt.auth != nil {
					if err := client.Auth(tt.auth); err != nil {
						t.Fatalf("AUTH PLAIN: %v", err)
					}
					if err := client.Reset(); err != nil {
						t.Fatalf("RSET after AUTH: %v", err)
					}
				}
				if err := client.Mail(from); err != nil {
					t.Fatalf("MAIL FROM: %v", err)
				}
				if err := client.Rcpt(to); err != nil {
					t.Fatalf("RCPT TO: %v", err)
				}
				data, err := client.Data()
				if err != nil {
					t.Fatalf("DATA: %v", err)
				}
				if _, err := io.WriteString(data, message); err != nil {
					t.Fatalf("write message: %v", err)
				}
				if err := data.Close(); err != nil {
					t.Fatalf("finish DATA: %v", err)
				}
			}

			// The SMTP reply to DATA is sent after the webhook request completes.
			select {
			case got := <-requests:
				if tt.wantCode != 0 || len(got.Embeds) != 1 {
					t.Fatalf("unexpected webhook payload: %+v", got)
				}
				embed := got.Embeds[0]
				if embed.Title != "📨 "+subject || strings.TrimRight(embed.Description, "\r\n") != body {
					t.Errorf("embed title = %q, description = %q", embed.Title, embed.Description)
				}
				fields := make(map[string]string)
				for _, field := range embed.Fields {
					fields[field.Name] = field.Value
				}
				if fields["From"] != from || fields["To"] != to {
					t.Errorf("embed addresses = %q and %q, want %q and %q", fields["From"], fields["To"], from, to)
				}
			default:
				if tt.wantCode == 0 {
					t.Fatal("expected a webhook request")
				}
			}
		})
	}
}

func checkSMTPCode(t *testing.T, err error, want int) {
	t.Helper()
	var reply *textproto.Error
	if !errors.As(err, &reply) || reply.Code != want {
		t.Fatalf("SMTP error = %v, want reply code %d", err, want)
	}
}

func newRelayClient(t *testing.T, webhook string, anonymous bool) *smtp.Client {
	t.Helper()

	discordSession, err := discord.New(webhook, "../etc/template.example.json")
	if err != nil {
		t.Fatal(err)
	}
	username, password := "relay", "secret"
	if anonymous {
		username, password = "", ""
	}
	backend, err := relaysmtp.New(discordSession, username, password)
	if err != nil {
		t.Fatal(err)
	}

	const timeout = 5 * time.Second
	const host = "127.0.0.1"
	server := gosmtp.NewServer(backend)
	server.Domain = "localhost"
	server.AllowInsecureAuth = true
	server.ReadTimeout = timeout
	server.WriteTimeout = timeout

	listener, err := net.Listen("tcp", host+":0")
	if err != nil {
		t.Fatal(err)
	}
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Serve(listener) }()
	t.Cleanup(func() {
		listener.Close()
		if err := <-serverDone; err != nil && !errors.Is(err, net.ErrClosed) {
			t.Errorf("SMTP server: %v", err)
		}
	})

	conn, err := net.DialTimeout("tcp", listener.Addr().String(), timeout)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		t.Fatal(err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		t.Fatal(err)
	}
	return client
}
