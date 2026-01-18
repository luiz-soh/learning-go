package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

type Usuarios struct {
	db *sql.DB
}

// Cria instancia de usuarios com banco para realizar as funções
// Pode ser passado qualqer banco
func NewUsuariosRepo(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

func (repository Usuarios) Criar(usuario models.Usuario) (uint64, error) {
	statement, erro := repository.db.Prepare("insert into usuarios (nome, nick, email, senha) values (?,?,?,?)")

	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	insercao, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)

	if erro != nil {
		return 0, erro
	}

	idInserido, erro := insercao.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(idInserido), nil
}

func (repository Usuarios) Buscar(nomeOuNick string) ([]models.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick)
	linhas, erro := repository.db.Query("select id, nome, email, nick, criadoEm from usuarios where nome like ? or nick like ?", nomeOuNick, nomeOuNick)

	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	var usuarios []models.Usuario

	for linhas.Next() {
		var usuario models.Usuario
		if erro := linhas.Scan(&usuario.Id, &usuario.Nome,
			&usuario.Email, &usuario.Nick, &usuario.CriadoEm); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (repository Usuarios) BuscarPorId(id uint64) (models.Usuario, error) {
	usuario := models.Usuario{}
	linha, erro := repository.db.Query("select id, nome, email,nick, criadoEm from usuarios where id = ?", id)
	if erro != nil {
		return usuario, erro
	}

	defer linha.Close()

	if linha.Next() {
		if erro := linha.Scan(&usuario.Id, &usuario.Nome,
			&usuario.Email, &usuario.Nick, &usuario.CriadoEm); erro != nil {
			return usuario, erro
		}
		return usuario, nil
	}

	return usuario, sql.ErrNoRows
}

func (repository Usuarios) Atualizar(usuarioId uint64, usuario models.Usuario) error {
	existe, erro := usuarioExiste(repository.db, usuarioId)

	if erro != nil {
		return erro
	}

	if !existe {
		return sql.ErrNoRows
	}

	statement, erro := repository.db.Prepare("update usuarios set nome = ?, nick = ?, email = ? where id = ?")

	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuarioId); erro != nil {
		return erro
	}
	return nil
}

func (repository Usuarios) Deletar(usuarioId uint64) error {
	existe, erro := usuarioExiste(repository.db, usuarioId)

	if erro != nil {
		return erro
	}

	if !existe {
		return sql.ErrNoRows
	}

	statement, erro := repository.db.Prepare("delete from usuarios where id = ?")

	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuarioId); erro != nil {
		return erro
	}
	return nil
}

func (repository Usuarios) BuscarPorEmail(email string) (models.Usuario, error) {
	usuario := models.Usuario{}
	linha, erro := repository.db.Query("select id, senha from usuarios where email = ?", email)
	if erro != nil {
		return usuario, erro
	}

	defer linha.Close()

	if linha.Next() {
		if erro := linha.Scan(&usuario.Id, &usuario.Senha); erro != nil {
			return usuario, erro
		}
		return usuario, nil
	}

	return usuario, sql.ErrNoRows
}

func (repository Usuarios) AlternarSeguir(usuarioId uint64, seguidorId uint64) error {
	resultado, erro := repository.db.Exec(
		"delete from seguidores where usuario_id = ? and seguidor_id = ?",
		usuarioId,
		seguidorId,
	)
	if erro != nil {
		return erro
	}

	linhasAfetadas, erro := resultado.RowsAffected()
	if erro != nil {
		return erro
	}

	if linhasAfetadas == 0 {
		_, erro = repository.db.Exec(
			"insert into seguidores (usuario_id, seguidor_id) values (?, ?)",
			usuarioId,
			seguidorId,
		)
		return erro
	}

	return nil
}

func (repository Usuarios) BuscarSeguidores(usuarioId uint64) ([]models.Usuario, error) {
	linhas, erro := repository.db.Query(`
		select u.nome, u.email, u.nick from usuarios u
		inner join seguidores s on u.id = s.seguidor_id
		where s.usuario_id = ?
	`, usuarioId)

	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	usuarios := []models.Usuario{} //Inicializo aqui para ter um slice vazio ao invés de nil

	for linhas.Next() {
		var usuario models.Usuario
		if erro := linhas.Scan(&usuario.Nome, &usuario.Email, &usuario.Nick); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (repository Usuarios) BuscarSeguindo(usuarioId uint64) ([]models.Usuario, error) {
	linhas, erro := repository.db.Query(`
		select u.nome, u.email, u.nick from usuarios u
		inner join seguidores s on u.id = s.usuario_id
		where s.seguidor_id = ?
	`, usuarioId)

	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	usuarios := []models.Usuario{} //Inicializo aqui para ter um slice vazio ao invés de nil

	for linhas.Next() {
		var usuario models.Usuario
		if erro := linhas.Scan(&usuario.Nome, &usuario.Email, &usuario.Nick); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (repository Usuarios) BuscarSenhaPorId(usuarioId uint64) (string, error) {
	var senha string
	linha, erro := repository.db.Query("select senha from usuarios where id = ?", usuarioId)
	if erro != nil {
		return "", erro
	}

	defer linha.Close()

	if linha.Next() {
		if erro := linha.Scan(&senha); erro != nil {
			return "", erro
		}
		return senha, nil
	}

	return "", sql.ErrNoRows
}

func (repository Usuarios) AlterarSenha(usuarioId uint64, novaSenha string) error {
	statement, erro := repository.db.Prepare("update usuarios set senha = ? where id = ?")

	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(novaSenha, usuarioId); erro != nil {
		return erro
	}

	return nil
}

func usuarioExiste(db *sql.DB, id uint64) (bool, error) {
	var exists bool
	err := db.QueryRow("select 1 from usuarios where id = ?", id).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil // não existe, mas não é erro
	}
	if err != nil {
		return false, err // erro de banco
	}
	return true, nil // existe
}
