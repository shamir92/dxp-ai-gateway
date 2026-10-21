package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
	TLS      string // starttls | tls | none
}

type SMTPMailer struct {
	cfg  SMTPConfig
	tmpl *Templates
}

func NewSMTPMailer(cfg SMTPConfig, tmpl *Templates) *SMTPMailer {
	return &SMTPMailer{cfg: cfg, tmpl: tmpl}
}

func (m *SMTPMailer) SendMagicCode(ctx context.Context, msg MagicCodeMessage) error {
	subject, textBody, htmlBody, err := m.tmpl.Render(msg)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	mime := buildMIME(m.headerFrom(), msg.To, subject, textBody, htmlBody)
	return m.send(msg.To, mime)
}

func (m *SMTPMailer) SendLoginAlert(ctx context.Context, msg LoginAlertMessage) error {
	subject, textBody, htmlBody, err := m.tmpl.RenderLoginAlert(msg)
	if err != nil {
		return fmt.Errorf("render login alert: %w", err)
	}
	mime := buildMIME(m.headerFrom(), msg.To, subject, textBody, htmlBody)
	return m.send(msg.To, mime)
}

// headerFrom is the RFC 5322 From header value: display name plus addr-spec.
func (m *SMTPMailer) headerFrom() string {
	if m.cfg.FromName != "" {
		return fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.From)
	}
	return m.cfg.From
}

func (m *SMTPMailer) send(to, mime string) error {
	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	// The envelope sender must be a bare reverse-path: no display name, no
	// angle brackets. "Name <a@b>" passes it verbatim to MAIL FROM and the
	// server rejects it with 501.
	envelopeFrom := address(m.cfg.From)

	switch m.cfg.TLS {
	case "tls":
		return sendTLS(addr, auth, envelopeFrom, to, mime)
	case "none":
		return smtp.SendMail(addr, auth, envelopeFrom, []string{to}, []byte(mime))
	default: // starttls
		return smtp.SendMail(addr, auth, envelopeFrom, []string{to}, []byte(mime))
	}
}

// address extracts the bare addr-spec from a possibly display-name decorated
// address: "DXP AI Gateway <noreply@example.com>" -> "noreply@example.com".
func address(from string) string {
	start := strings.IndexByte(from, '<')
	if start < 0 {
		return strings.TrimSpace(from)
	}
	end := strings.IndexByte(from[start:], '>')
	if end < 0 {
		return strings.TrimSpace(from[start+1:])
	}
	return strings.TrimSpace(from[start+1 : start+end])
}

func sendTLS(addr string, auth smtp.Auth, from, to, mime string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		return err
	}
	defer conn.Close()
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Quit()
	if err := c.Auth(auth); err != nil {
		return err
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(mime)); err != nil {
		return err
	}
	return w.Close()
}

func buildMIME(from, to, subject, textBody, htmlBody string) string {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	boundary := "dxp-magic-code-boundary"
	b.WriteString("Content-Type: multipart/alternative; boundary=" + boundary + "\r\n\r\n")
	b.WriteString("--" + boundary + "\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	b.WriteString(textBody)
	b.WriteString("\r\n--" + boundary + "\r\n")
	b.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	b.WriteString(htmlBody)
	b.WriteString("\r\n--" + boundary + "--\r\n")
	return b.String()
}

// ConsoleMailer prints messages — tests and local dry-run only.
type ConsoleMailer struct {
	tmpl *Templates
}

func NewConsoleMailer(tmpl *Templates) *ConsoleMailer {
	return &ConsoleMailer{tmpl: tmpl}
}

func (m *ConsoleMailer) SendMagicCode(_ context.Context, msg MagicCodeMessage) error {
	subject, textBody, _, err := m.tmpl.Render(msg)
	if err != nil {
		return err
	}
	fmt.Printf("[mail:console] to=%s subject=%q\n%s\n", msg.To, subject, textBody)
	return nil
}

func (m *ConsoleMailer) SendLoginAlert(_ context.Context, msg LoginAlertMessage) error {
	subject, textBody, _, err := m.tmpl.RenderLoginAlert(msg)
	if err != nil {
		return err
	}
	fmt.Printf("[mail:console] to=%s subject=%q\n%s\n", msg.To, subject, textBody)
	return nil
}
