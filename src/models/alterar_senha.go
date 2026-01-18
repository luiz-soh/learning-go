package models

import "errors"

type AlterarSenha struct {
	Nova  string `json:"nova,omitempty"`
	Atual string `json:"atual,omitempty"`
}

func (senha *AlterarSenha) Validar() error {

	if senha.Atual == "" {
		return errors.New("É obrigatório informar a senha atual")
	}
	if senha.Nova == "" {
		return errors.New("É obrigatório informar a nova senha")
	}
	return nil
}
