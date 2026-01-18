package controllers

import (
	"api/src/authentication"
	"api/src/database"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// @Summary		Criar Publicação
// @Description Cria uma nova publicação
// @Tags 	publicacoes
// @Accept	json
// @Produce	json
// @Param publicacao body models.PublicacaoRequest true "Dados da publicação"
// @Success	201 {object} models.Publicacao
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /publicacoes [post]
func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
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

	var publicacaoRequest models.PublicacaoRequest

	if erro = json.Unmarshal(bodyRequest, &publicacaoRequest); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	publicacao := models.Publicacao{
		Titulo:   publicacaoRequest.Titulo,
		Conteudo: publicacaoRequest.Conteudo,
		AutorId:  usuarioId,
	}

	if erro := publicacao.Preparar(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewPublicacoesRepo(db)
	publicacao.Id, erro = repositorio.Criar(usuarioId, publicacao)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, publicacao)
}

// @Summary		Buscar Publicações
// @Description Busca publicações do usuário autenticado e seus seguidores
// @Tags 	publicacoes
// @Accept	json
// @Produce	json
// @Success	200 {array} models.Publicacao
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /publicacoes [get]
func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {
	usuarioId, erro := authentication.ExtrairUsuarioId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, nil)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewPublicacoesRepo(db)
	publicacoes, erro := repositorio.BuscarPublicacoes(usuarioId)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publicacoes)
}

// @Summary		Buscar Publicação
// @Description Busca uma publicação por ID
// @Tags 	publicacoes
// @Accept	json
// @Produce	json
// @Param id path int true "ID da publicação"
// @Success	200 {object} models.Publicacao
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /publicacoes/{id} [get]
func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 64)

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

	repositorio := repositories.NewPublicacoesRepo(db)
	publicacao, erro := repositorio.BuscarPorId(ID)
	if erro != nil {
		if erro == sql.ErrNoRows {
			responses.Erro(w, http.StatusNotFound, erro)
			return
		}

		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publicacao)
}

// @Summary		Atualizar Publicação
// @Description Atualiza uma publicação do usuário autenticado
// @Tags 	publicacoes
// @Accept	json
// @Produce	json
// @Param id path int true "ID da publicação"
// @Param publicacao body models.PublicacaoRequest true "Dados atualizados da publicação"
// @Success	204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /publicacoes/{id} [put]
func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	publicacaoId, erro := strconv.ParseUint(parametros["id"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

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

	var publicacaoRequest models.PublicacaoRequest

	if erro = json.Unmarshal(bodyRequest, &publicacaoRequest); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	publicacao := models.Publicacao{
		Titulo:   publicacaoRequest.Titulo,
		Conteudo: publicacaoRequest.Conteudo,
		AutorId:  usuarioId,
	}

	if erro := publicacao.Preparar(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewPublicacoesRepo(db)
	publicacaoExiste, erro := repositorio.PublicacaoUsuarioExiste(publicacaoId, usuarioId)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if !publicacaoExiste {
		responses.Erro(w, http.StatusNotFound, nil)
		return
	}
	publicacao.Id = publicacaoId
	if erro = repositorio.Atualizar(publicacao); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// @Summary		Deletar Publicação
// @Description Deleta uma publicação do usuário autenticado
// @Tags 	publicacoes
// @Accept	json
// @Produce	json
// @Param id path int true "ID da publicação"
// @Success	204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /publicacoes/{id} [delete]
func DeletarPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	publicacaoId, erro := strconv.ParseUint(parametros["id"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioId, erro := authentication.ExtrairUsuarioId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, nil)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewPublicacoesRepo(db)
	publicacaoExiste, erro := repositorio.PublicacaoUsuarioExiste(publicacaoId, usuarioId)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if !publicacaoExiste {
		responses.Erro(w, http.StatusNotFound, nil)
		return
	}

	if erro = repositorio.Deletar(publicacaoId); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

// @Summary		Buscar Publicações do Usuário
// @Description Busca publicações de um usuário específico
// @Tags 	publicacoes
// @Accept	json
// @Produce	json
// @Param usuarioId path int true "ID do usuário"
// @Success	200 {array} models.Publicacao
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /usuarios/{usuarioId}/publicacoes [get]
func BuscarPublicacoesUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	usuarioId, erro := strconv.ParseUint(parametros["usuarioId"], 10, 64)

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

	repositorio := repositories.NewPublicacoesRepo(db)
	publicacoes, erro := repositorio.BuscarPublicacoesUsuario(usuarioId)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publicacoes)
}

// @Summary		Curtir Publicação
// @Description Curte ou descurte uma publicação
// @Tags 	publicacoes
// @Accept	json
// @Produce	json
// @Param id path int true "ID da publicação"
// @Success	204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /publicacoes/{id}/curtir [post]
func CurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	publicacaoId, erro := strconv.ParseUint(parametros["id"], 10, 64)

	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioId, erro := authentication.ExtrairUsuarioId(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, nil)
		return
	}

	db, erro := database.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NewPublicacoesRepo(db)
	if erro = repositorio.CurtiPublicacao(publicacaoId, usuarioId); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
