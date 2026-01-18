package models

import (
	"api/src/security"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

// Usuario representa o usuario utilizado na plataforma
type Usuario struct {
	Id       uint64    `json:"id,omitempty"`
	Nome     string    `json:"nome,omitempty"`
	Nick     string    `json:"nick,omitempty"`
	Email    string    `json:"email,omitempty"`
	Senha    string    `json:"senha,omitempty"`
	CriadoEm time.Time `json:"CriadoEm,omitzero"`
}

func (usuario *Usuario) Preparar(cadastro bool) error {
	if erro := usuario.validar(cadastro); erro != nil {
		return erro
	}

	if erro := usuario.formatar(cadastro); erro != nil {
		return erro
	}
	return nil
}

func (usuario *Usuario) validar(cadastro bool) error {

	if usuario.Nome == "" {
		return errors.New("Nome do usuário é obrigatório")
	}
	if usuario.Nick == "" {
		return errors.New("Nick do usuário é obrigatório")
	}
	if usuario.Email == "" {
		return errors.New("Email do usuário é obrigatório")
	}

	if erro := checkmail.ValidateFormat(usuario.Email); erro != nil {
		return errors.New("E-mail inválido")
	}

	if cadastro {
		if usuario.Senha == "" {
			return errors.New("Senha do usuário é obrigatória")
		}
	}
	return nil
}

func (usuario *Usuario) formatar(cadastro bool) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)

	if cadastro {
		senhaComHash, erro := security.Hash(usuario.Senha)
		if erro != nil {
			return erro
		}
		usuario.Senha = string(senhaComHash)
	}
	return nil
}
