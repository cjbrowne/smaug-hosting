package email

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
	"html/template"
	"os"
	"strings"
)

// todo: i18n
const plainTextTempl = `
	Welcome to Smaug Hosting!  

	Please verify your email address by copying the following URL into your browser:
    {.VerificationUrl}

    Yours Sincerely,

    Smaug Hosting
`

const htmlEmailTempl = `
<html>
<body>
	<p>
		Welcome to Smaug Hosting!
	</p>
	<p>
		Please verify your email by clicking the following link:<br/>
		<a href="{{.VerificationUrl}}">Verify Email Address</a><br/>
	</p>
	<p>
    	Yours Sincerely,
	</p>
	<p>
    	Smaug Hosting
	</p>
</body>
</html>
`

func SendVerificationEmail(address string, verificationCode string) error {
	from := mail.NewEmail("Smaug Hosting", "no-reply@smaug-hosting.co.uk")
	subject := "Welcome to Smaug Hosting!"
	to := mail.NewEmail("Smaug Hosting User", address)

	plainText, err := template.New("plainTextTemplate").Parse(plainTextTempl)
	html, err := template.New("htmlTemplate").Parse(htmlEmailTempl)

	var plainTextContent bytes.Buffer
	var htmlContent bytes.Buffer

	frontendBaseUrl := strings.TrimRight(os.Getenv("FRONTEND_BASE_URL"), "/")

	templateVars := struct{
		VerificationUrl string
	}{
		VerificationUrl: fmt.Sprintf("%s/verify?code=%s", frontendBaseUrl, verificationCode),
	}

	err = plainText.Execute(&plainTextContent, templateVars)
	if err != nil {
		return err
	}

	err = html.Execute(&htmlContent, templateVars)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject, to, plainTextContent.String(), htmlContent.String())
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err == nil && (response.StatusCode < 200 || response.StatusCode > 299) {
		logrus.Errorf("Got response %d from sendgrid: %s", response.StatusCode, response.Body)
		return errors.New("sendgrid response not 200")
	}
	return err
}
