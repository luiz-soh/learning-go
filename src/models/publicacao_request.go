package models

type PublicacaoRequest struct {
	Titulo   string `json:"titulo,omitempty"`
	Conteudo string `json:"conteudo,omitempty"`
}
