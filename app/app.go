package app

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/mailgun/mailgun-go/v3"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

// Config stores our app's config
type Config struct {
	MailgunAPIKey       string `envconfig:"MAILGUN_API_KEY" required:"true"`
	MailgunSenderDomain string `envconfig:"MAILGUN_SENDER_DOMAIN" required:"true"`
	SenderName          string `envconfig:"SENDER_NAME" required:"true"`
	SenderEmail         string `envconfig:"SENDER_EMAIL", required:"true"`
	ReplyToName         string `envconfig:"REPLY_TO_NAME" required:"true"`
	ReplyToEmail        string `envconfig:"REPLY_TO_EMAIL", required:"true"`
	EmailSubject        string `envconfig:"EMAIL_SUBJECT", required:"true"`
	MessageSignOff      string `envconfig:"MESSAGE_SIGN_OFF", required:"true"`
}

// MustParseConfig returns an inflated Config object from the provided file path
// or fails on error
func MustParseConfig(path string) *Config {
	if err := godotenv.Load(path); err != nil {
		log.Fatal(err)
	}

	var cfg Config
	envconfig.MustProcess("", &cfg)

	return &cfg
}

// MustParseTemplate returns a parsed template object by walking the provided directory path
// or fails on error
func MustParseTemplate(path string) *template.Template {
	tpl := template.New("emails")

	if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		tpl, err = tpl.ParseFiles(path)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	return tpl
}

// NewMailgunClient returns a new Mailgun client
func NewMailgunClient(senderDomain, apiKey string) *mailgun.MailgunImpl {
	mg := mailgun.NewMailgun(senderDomain, apiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)
	return mg
}

// Container defines the interface for our app's Dependencies container
type Container interface {
	ConfigInjector
	TemplateInjector
	MailgunInjector
}

// ConfigInjector defines the behaviour for injecting our Config
type ConfigInjector interface {
	Config() *Config
}

// TemplateInjector defines the behaviour for injecting our Template
type TemplateInjector interface {
	Template() *template.Template
}

// MailgunInjector defines the behaviour for injecting our Mailgun client
type MailgunInjector interface {
	Mailgun() *mailgun.MailgunImpl
}
