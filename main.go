package main

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
