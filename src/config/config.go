package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ConnectionString = ""
	Porta            = 5000
	SecretKey        []byte
	DockerRun        = false
)

// Carregar vai preencher as variaveis de ambiente
func Carregar() {
	var erro error

	_ = godotenv.Load()

	Porta, erro = strconv.Atoi(os.Getenv("API_PORT"))
	if erro != nil {
		Porta = 9000
	}

	DockerRun, erro = strconv.ParseBool(os.Getenv("DOCKER"))
	if erro != nil {
		DockerRun = false
	}

	ConnectionString = os.Getenv("ConnectionString")

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
