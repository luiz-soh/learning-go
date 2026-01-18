package router

import (
	"api/src/router/routes"

	"github.com/gorilla/mux"
)

// Gerar deve retornar as rotas configuradas
func Gerar() *mux.Router {
	r := mux.NewRouter()
	return routes.Configurar(r)
}
