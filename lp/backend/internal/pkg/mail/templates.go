package mail

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
	texttemplate "text/template"
)

// embedded holds the templates so the binary is self-contained: it no longer
// needs the working directory (or a Docker COPY) to render mail.
//
//go:embed templates/*.tmpl
var embedded embed.FS

type templateData struct {
	Code             string
	Purpose          string
	ExpiresInMinutes int
	ProductName      string
	SupportEmail     string
	ActionLabel      string
}

type Templates struct {
	html              *template.Template
	text              *texttemplate.Template
	subject           *texttemplate.Template
	loginAlertSubject *texttemplate.Template
	loginAlertText    *texttemplate.Template
	loginAlertHTML    *template.Template
}

// LoadTemplates parses the templates embedded in the binary.
func LoadTemplates() (*Templates, error) {
	htmlT, err := template.ParseFS(embedded, "templates/magic_code.html.tmpl")
	if err != nil {
		return nil, fmt.Errorf("html template: %w", err)
	}
	textT, err := texttemplate.ParseFS(embedded, "templates/magic_code.txt.tmpl")
	if err != nil {
		return nil, fmt.Errorf("text template: %w", err)
	}
	subjT, err := texttemplate.ParseFS(embedded, "templates/magic_code.subject.tmpl")
	if err != nil {
		return nil, fmt.Errorf("subject template: %w", err)
	}
	laSubj, err := texttemplate.ParseFS(embedded, "templates/login_alert.subject.tmpl")
	if err != nil {
		return nil, fmt.Errorf("login alert subject template: %w", err)
	}
	laText, err := texttemplate.ParseFS(embedded, "templates/login_alert.txt.tmpl")
	if err != nil {
		return nil, fmt.Errorf("login alert text template: %w", err)
	}
	laHTML, err := template.ParseFS(embedded, "templates/login_alert.html.tmpl")
	if err != nil {
		return nil, fmt.Errorf("login alert html template: %w", err)
	}
	return &Templates{
		html:              htmlT,
		text:              textT,
		subject:           subjT,
		loginAlertSubject: laSubj,
		loginAlertText:    laText,
		loginAlertHTML:    laHTML,
	}, nil
}

func (t *Templates) Render(msg MagicCodeMessage) (subject, textBody, htmlBody string, err error) {
	data := templateData{
		Code:             msg.Code,
		Purpose:          msg.Purpose,
		ExpiresInMinutes: msg.ExpiresInMinutes,
		ProductName:      "DXP AI Gateway",
		SupportEmail:     "hello@fromdxp.com",
		ActionLabel:      "sign in",
	}
	if msg.Purpose == "register" {
		data.ActionLabel = "finish signing up"
	}
	var subjBuf, textBuf, htmlBuf bytes.Buffer
	if err := t.subject.Execute(&subjBuf, data); err != nil {
		return "", "", "", err
	}
	if err := t.text.Execute(&textBuf, data); err != nil {
		return "", "", "", err
	}
	if err := t.html.Execute(&htmlBuf, data); err != nil {
		return "", "", "", err
	}
	return strings.TrimSpace(subjBuf.String()), textBuf.String(), htmlBuf.String(), nil
}

type loginAlertData struct {
	Device       string
	IP           string
	Location     string
	LoggedInAt   string
	ProductName  string
	SupportEmail string
}

func (t *Templates) RenderLoginAlert(msg LoginAlertMessage) (subject, textBody, htmlBody string, err error) {
	data := loginAlertData{
		Device:       msg.Device,
		IP:           msg.IP,
		Location:     msg.Location,
		LoggedInAt:   msg.LoggedInAt,
		ProductName:  "DXP AI Gateway",
		SupportEmail: "hello@fromdxp.com",
	}
	var subjBuf, textBuf, htmlBuf bytes.Buffer
	if err := t.loginAlertSubject.Execute(&subjBuf, data); err != nil {
		return "", "", "", err
	}
	if err := t.loginAlertText.Execute(&textBuf, data); err != nil {
		return "", "", "", err
	}
	if err := t.loginAlertHTML.Execute(&htmlBuf, data); err != nil {
		return "", "", "", err
	}
	return strings.TrimSpace(subjBuf.String()), textBuf.String(), htmlBuf.String(), nil
}
