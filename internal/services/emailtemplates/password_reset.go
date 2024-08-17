package emailtemplates

import (
	"strings"
	"text/template"
)

var templatePasswordReset = template.Must(template.New("password-reset").
	Parse(`Hi,
We received a request to reset your password for your {{ .CompanyName }} account. If you didn’t make this request, you can ignore this email.

To reset your password, click the link below:

{{ .BaseURL }}/reset-password?token={{ .Token }}

This link will expire in 24 hours for your security.

Best regards,
{{ .CompanyName }} Team`))

func PasswordReset(company, baseurl, token string) string {
	bldr := &strings.Builder{}

	_ = templatePasswordReset.Execute(bldr, map[string]any{
		"CompanyName": company,
		"Token":       token,
		"BaseURL":     baseurl,
	})

	return bldr.String()
}
