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

// @Summary		Criar Usuário
// @Description Cria um usuário
// @Tags 	usuarios
// @Accept	json
// @Produce	json
// @Param usuario body models.CreateUsuarioRequest true "Dados do usuário"
// @Success	201 {object} models.Usuario
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /usuarios [post]
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	bodyRequest, erro := io.ReadAll(r.Body)

	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var request models.CreateUsuarioRequest

	if erro = json.Unmarshal(bodyRequest, &request); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuario := models.Usuario{
		Nome:  request.Nome,
		Nick:  request.Nick,
		Senha: request.Senha,
		Email: request.Email,
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

// @Summary		Listar Usuários
// @Description Lista usuários, opcionalmente filtrando por nome ou nick
// @Tags 	usuarios
// @Accept	json
// @Produce	json
// @Param usuario query string false "Nome ou nick do usuário"
// @Success	200 {array} models.Usuario
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /usuarios [get]
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

// @Summary		Buscar Usuário
// @Description Busca um usuário por ID
// @Tags 	usuarios
// @Accept	json
// @Produce	json
// @Param id path int true "ID do usuário"
// @Success	200 {object} models.Usuario
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /usuarios/{id} [get]
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

// @Summary		Atualizar Usuário
// @Description Atualiza os dados de um usuário
// @Tags 	usuarios
// @Accept	json
// @Produce	json
// @Param id path int true "ID do usuário"
// @Param usuario body models.UpdateUsuarioRequest true "Dados atualizados do usuário"
// @Success	204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /usuarios/{id} [put]
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

	var request models.UpdateUsuarioRequest

	if erro = json.Unmarshal(bodyRequest, &request); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuario := models.Usuario{
		Nome:  request.Nome,
		Nick:  request.Nick,
		Email: request.Email,
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

// @Summary		Deletar Usuário
// @Description Deleta um usuário
// @Tags 	usuarios
// @Accept	json
// @Produce	json
// @Param id path int true "ID do usuário"
// @Success	204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /usuarios/{id} [delete]
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

// @Summary		Alternar Seguir Usuário
// @Description Segue ou deixa de seguir um usuário
// @Tags 	usuarios
// @Accept	json
// @Produce	json
// @Param id path int true "ID do usuário"
// @Success	204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /usuarios/{id}/seguir [post]
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

// @Summary		Buscar Seguidores
// @Description Busca os seguidores de um usuário
// @Tags 	usuarios
// @Accept	json
// @Produce	json
// @Param id path int true "ID do usuário"
// @Success	200 {array} models.Usuario
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /usuarios/{id}/seguidores [get]
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

// @Summary		Buscar Seguindo
// @Description Busca os usuários que um usuário está seguindo
// @Tags 	usuarios
// @Accept	json
// @Produce	json
// @Param id path int true "ID do usuário"
// @Success	200 {array} models.Usuario
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /usuarios/{id}/seguindo [get]
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

// @Summary		Alterar Senha
// @Description Altera a senha do usuário autenticado
// @Tags 	usuarios
// @Accept	json
// @Produce	json
// @Param alterarSenha body models.AlterarSenha true "Senhas atual e nova"
// @Success	204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /usuarios/alterar_senha [post]
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
