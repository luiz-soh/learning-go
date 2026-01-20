package main

// @title DevBook API
// @version 1.0
// @description API para rede social DevBook
// @host localhost:5000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"api/src/config"
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Carregar()
	r := router.Gerar()

	fmt.Printf("Rodando API na porta %d", config.Porta)
	porta := fmt.Sprintf("localhost:%d", config.Porta)
	if config.DockerRun {
		porta = fmt.Sprintf("0.0.0.0:%d", config.Porta)
	}
	log.Fatal(http.ListenAndServe(porta, r))
}
