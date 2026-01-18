package repositories

import (
	"api/src/models"
	"database/sql"
)

type Publicacoes struct {
	db *sql.DB
}

// Cria instancia de publicacoes com banco para realizar as funções
// Pode ser passado qualqer banco
func NewPublicacoesRepo(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

func (repository Publicacoes) Criar(usuarioId uint64, publicacao models.Publicacao) (uint64, error) {
	statement, erro := repository.db.Prepare("insert into publicacoes (titulo, conteudo, autor_id) values (?,?,?)")

	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	insercao, erro := statement.Exec(publicacao.Titulo, publicacao.Conteudo, usuarioId)

	if erro != nil {
		return 0, erro
	}

	idInserido, erro := insercao.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(idInserido), nil
}

func (repository Publicacoes) BuscarPorId(id uint64) (models.Publicacao, error) {
	publicacao := models.Publicacao{}
	linha, erro := repository.db.Query(`SELECT p.id, p.titulo, p.conteudo, p.autor_id, p.criadaEm,
       COUNT(c.publicacao_id) AS curtidas,
       u.nick FROM publicacoes p
	   INNER JOIN usuarios u ON u.id = p.autor_id
	   LEFT JOIN curtidas c ON c.publicacao_id = p.id
	   WHERE p.id = ?
	   GROUP BY p.id, p.titulo, p.conteudo, p.autor_id, p.criadaEm, u.nick`, id)
	if erro != nil {
		return publicacao, erro
	}

	// 	linha, erro := repository.db.Query(`select p.*, u.nick from publicacoes p
	// inner join usuarios u on u.id = p.autor_id where p.id = ? order by 1 desc`, id)

	defer linha.Close()

	if linha.Next() {
		if erro := linha.Scan(&publicacao.Id, &publicacao.Titulo, &publicacao.Conteudo,
			&publicacao.AutorId, &publicacao.CriadaEm, &publicacao.Curtidas,
			&publicacao.AutorNick); erro != nil {
			return publicacao, erro
		}
		return publicacao, nil
	}

	return publicacao, sql.ErrNoRows

}

func (repository Publicacoes) BuscarPublicacoes(usuarioId uint64) ([]models.Publicacao, error) {
	publicacoes := []models.Publicacao{}
	linhas, erro := repository.db.Query(`SELECT p.id, p.titulo, p.conteudo, p.autor_id, p.criadaEm,
       COUNT(c.publicacao_id) AS curtidas,
       u.nick FROM publicacoes p
	   INNER JOIN usuarios u ON u.id = p.autor_id
	   LEFT JOIN seguidores s on p.autor_id = s.usuario_id
	   LEFT JOIN curtidas c ON c.publicacao_id = p.id
	   where u.id = ? or s.seguidor_id = ?
	   GROUP BY p.id, p.titulo, p.conteudo, p.autor_id, p.criadaEm, u.nick
	   ORDER BY 1 DESC
	`, usuarioId, usuarioId)

	if erro != nil {
		return nil, erro
	}

	for linhas.Next() {
		var publicacao models.Publicacao
		if erro := linhas.Scan(&publicacao.Id, &publicacao.Titulo, &publicacao.Conteudo,
			&publicacao.AutorId, &publicacao.CriadaEm, &publicacao.Curtidas,
			&publicacao.AutorNick); erro != nil {
			return nil, erro
		}
		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

func (repository Publicacoes) PublicacaoUsuarioExiste(publicacaoId uint64, usuarioId uint64) (bool, error) {
	var exists bool
	err := repository.db.QueryRow("select 1 from publicacoes where id = ? and autor_id = ?",
		publicacaoId, usuarioId).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil // não existe, mas não é erro
	}
	if err != nil {
		return false, err // erro de banco
	}
	return true, nil // existe

}

func (repository Publicacoes) Atualizar(publicacao models.Publicacao) error {
	statement, erro := repository.db.Prepare("update publicacoes set titulo = ?, conteudo = ? where id = ?")

	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacao.Id); erro != nil {
		return erro
	}
	return nil
}

func (repository Publicacoes) Deletar(publicacaoId uint64) error {
	statement, erro := repository.db.Prepare("delete from publicacoes where id = ?")

	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro := statement.Exec(publicacaoId); erro != nil {
		return erro
	}
	return nil
}

func (repository Publicacoes) BuscarPublicacoesUsuario(usuarioId uint64) ([]models.Publicacao, error) {
	publicacoes := []models.Publicacao{}
	linhas, erro := repository.db.Query(`SELECT p.id, p.titulo, p.conteudo, p.autor_id, p.criadaEm,
       COUNT(c.publicacao_id) AS curtidas,
       u.nick FROM publicacoes p
	   INNER JOIN usuarios u ON u.id = p.autor_id
	   LEFT JOIN curtidas c ON c.publicacao_id = p.id
	   WHERE u.id = ?
	   GROUP BY p.id, p.titulo, p.conteudo, p.autor_id, p.criadaEm, u.nick
	   ORDER BY 1 DESC`, usuarioId)
	if erro != nil {
		return nil, erro
	}

	for linhas.Next() {
		var publicacao models.Publicacao
		if erro := linhas.Scan(&publicacao.Id, &publicacao.Titulo, &publicacao.Conteudo,
			&publicacao.AutorId, &publicacao.CriadaEm, &publicacao.Curtidas,
			&publicacao.AutorNick); erro != nil {
			return nil, erro
		}
		publicacoes = append(publicacoes, publicacao)
	}

	return publicacoes, nil
}

func (repository Publicacoes) CurtiPublicacao(publicacaoId uint64, usuarioId uint64) error {
	resultado, erro := repository.db.Exec(
		"delete from curtidas where usuario_id = ? and publicacao_id = ?",
		usuarioId,
		publicacaoId,
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
			"insert into curtidas (usuario_id, publicacao_id) values (?, ?)",
			usuarioId,
			publicacaoId,
		)
		return erro
	}

	return nil
}
