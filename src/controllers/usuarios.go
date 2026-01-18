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
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	bodyRequest, erro := io.ReadAll(r.Body)

	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario

	if erro = json.Unmarshal(bodyRequest, &usuario); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro := usuario.Preparar(true); erro != nil {
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
	usuario.Id, erro = repositorio.Criar(usuario)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, usuario)
}

func ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, erro := database.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewUsuariosRepo(db)
	usuarios, erro := repositorio.Buscar(nomeOuNick)

	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, usuarios)
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)

	if erro != nil {
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
	usuario, erro := repositorio.BuscarPorId(ID)
	if erro != nil {
		if erro == sql.ErrNoRows {
			responses.Erro(w, http.StatusNotFound, erro)
			return
		}

		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, usuario)
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioId, erro := strconv.ParseUint(parametros["id"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIdToken, erro := authentication.ExtrairUsuarioId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioId != usuarioIdToken {
		responses.Erro(w, http.StatusForbidden, nil)
		return
	}

	bodyRequest, erro := io.ReadAll(r.Body)

	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario

	if erro = json.Unmarshal(bodyRequest, &usuario); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro := usuario.Preparar(false); erro != nil {
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
	if erro = repositorio.Atualizar(usuarioId, usuario); erro != nil {
		if erro == sql.ErrNoRows {
			responses.Erro(w, http.StatusNotFound, erro)
			return
		}
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioId, erro := strconv.ParseUint(parametros["id"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIdToken, erro := authentication.ExtrairUsuarioId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioId != usuarioIdToken {
		responses.Erro(w, http.StatusForbidden, nil)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewUsuariosRepo(db)
	if erro = repositorio.Deletar(usuarioId); erro != nil {
		if erro == sql.ErrNoRows {
			responses.Erro(w, http.StatusNotFound, erro)
			return
		}
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func AlternarSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorId, erro := authentication.ExtrairUsuarioId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, nil)
		return
	}

	parametros := mux.Vars(r)
	usuarioId, erro := strconv.ParseUint(parametros["id"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioId == seguidorId {
		responses.Erro(w, http.StatusBadRequest, errors.New("Você não pode seguir ou deixar de seguir a você mesmo"))
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewUsuariosRepo(db)
	if erro = repositorio.AlternarSeguir(usuarioId, seguidorId); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}

func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioId, erro := strconv.ParseUint(parametros["id"], 10, 64)
	if erro != nil {
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
	usuarios, erro := repositorio.BuscarSeguidores(usuarioId)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, usuarios)
}

func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioId, erro := strconv.ParseUint(parametros["id"], 10, 64)
	if erro != nil {
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
	usuarios, erro := repositorio.BuscarSeguindo(usuarioId)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, usuarios)
}

func AlterarSenha(w http.ResponseWriter, r *http.Request) {
	usuarioId, erro := authentication.ExtrairUsuarioId(r)

	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, nil)
		return
	}
	bodyRequest, erro := io.ReadAll(r.Body)

	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var alterarSenha models.AlterarSenha

	if erro = json.Unmarshal(bodyRequest, &alterarSenha); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro := alterarSenha.Validar(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	novaSenha, erro := security.Hash(alterarSenha.Nova)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
	}

	db, erro := database.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewUsuariosRepo(db)
	senhaAtual, erro := repositorio.BuscarSenhaPorId(usuarioId)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, nil)
		return
	}

	if erro := security.VerificarSenha(alterarSenha.Atual, senhaAtual); erro != nil {
		responses.Erro(w, http.StatusUnauthorized, errors.New("Senha atual incorreta"))
		return
	}

	if erro := repositorio.AlterarSenha(usuarioId, string(novaSenha)); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}
