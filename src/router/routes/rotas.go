package routes

import (
	"api/src/middlewares"
	"net/http"

	_ "api/docs" // importa os docs gerados

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Rota representa a struct para as rotas da API
type Rota struct {
	URI                string
	Metodo             string
	Funcao             func(http.ResponseWriter, *http.Request)
	RequerAutenticacao bool
}

// Configurar insere todas as rotas  no mux Router
func Configurar(r *mux.Router) *mux.Router {
	rotas := rotasUsuarios
	rotas = append(rotas, loginRoute)
	rotas = append(rotas, rotasPublicacoes...) //Os 3 pontos em sequencia é para informar que ta passando uma lista como append

	for _, rota := range rotas {
		if rota.RequerAutenticacao {
			r.HandleFunc(rota.URI,
				middlewares.Logger(middlewares.Autenticar(rota.Funcao))).Methods(rota.Metodo)
		} else {
			r.HandleFunc(rota.URI, middlewares.Logger(rota.Funcao)).Methods(rota.Metodo)
		}
	}

	// rota para acessar o swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
