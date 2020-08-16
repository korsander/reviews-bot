package cfg

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	SlackToken        string
	VerificationToken string
	CertChain         string
	CertPrivate       string
	EventsAddr        string
	CIAddr            string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return Config{
		SlackToken:        os.Getenv("SLACK_TOKEN"),
		VerificationToken: os.Getenv("VERIFICATION_TOKEN"),
		CertChain:         os.Getenv("CERT_FULL_CHAIN"),
		CertPrivate:       os.Getenv("CERT_PRIVATE_KEY"),
		EventsAddr:        os.Getenv("EVENTS_ADDR"),
		CIAddr:            os.Getenv("CI_ADDR"),
	}
}
