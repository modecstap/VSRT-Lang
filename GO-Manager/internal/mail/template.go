package mail

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed password_reset.tmpl
var passwordResetTemplate string

type Message struct {
	ContentType string
	Subject     string
	Body        string
}

type templateData struct {
	Link string
}

func RenderPasswordReset(link string) (Message, error) {
	headers, bodyTmpl, err := splitTemplate(passwordResetTemplate)
	if err != nil {
		return Message{}, err
	}

	tmpl, err := template.New("password_reset_body").Parse(bodyTmpl)
	if err != nil {
		return Message{}, err
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, templateData{Link: link}); err != nil {
		return Message{}, err
	}

	subjectTmpl := headers["Subject"]
	if subjectTmpl == "" {
		return Message{}, fmt.Errorf("password reset template missing Subject")
	}
	st, err := template.New("password_reset_subject").Parse(subjectTmpl)
	if err != nil {
		return Message{}, err
	}
	var subject bytes.Buffer
	if err := st.Execute(&subject, templateData{Link: link}); err != nil {
		return Message{}, err
	}

	contentType := headers["Content-Type"]
	if contentType == "" {
		contentType = "text/plain"
	}

	return Message{
		ContentType: contentType,
		Subject:     subject.String(),
		Body:        strings.TrimSuffix(body.String(), "\n"),
	}, nil
}

func splitTemplate(raw string) (map[string]string, string, error) {
	parts := strings.SplitN(raw, "\n\n", 2)
	if len(parts) != 2 {
		return nil, "", fmt.Errorf("password reset template missing body separator")
	}
	headers := map[string]string{}
	for _, line := range strings.Split(parts[0], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return nil, "", fmt.Errorf("invalid header line: %q", line)
		}
		headers[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return headers, parts[1], nil
}
