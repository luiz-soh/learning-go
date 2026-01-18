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
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", config.Porta), r))

}
