package controllers

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/security"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	requestBody, erro := io.ReadAll(r.Body)

	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario
	if erro = json.Unmarshal(requestBody, &usuario); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewUsuariosRepo(db)
	usuarioComSenha, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		if erro == sql.ErrNoRows {
			responses.Erro(w, http.StatusUnauthorized, errors.New("usuário ou senha incorretos"))
			return
		}
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = security.VerificarSenha(usuario.Senha, usuarioComSenha.Senha); erro != nil {
		responses.Erro(w, http.StatusUnauthorized, errors.New("usuário ou senha incorretos"))
		return
	}

	token, erro := authentication.GenerateToken(usuarioComSenha.Id)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, token)
}
