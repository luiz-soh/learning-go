package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ConnectionString = ""
	Porta            = 5000
	SecretKey        []byte
)

// Carregar vai preencher as variaveis de ambiente
func Carregar() {
	var erro error

	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	Porta, erro = strconv.Atoi(os.Getenv("API_PORT"))
	if erro != nil {
		Porta = 9000
	}

	ConnectionString = os.Getenv("ConnectionString")

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
