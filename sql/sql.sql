CREATE DATABASE IF NOT EXISTS devbook;
USE devbook;

DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios(
    Id int auto_increment primary key,
    nome varchar(50) not null,
    nick varchar(50) not null unique,
    email varchar(50) not null unique,
    senha varchar(100) not null,
    criadoEm timestamp default current_timestamp()
);

DROP TABLE IF EXISTS seguidores;

CREATE TABLE seguidores(
    usuario_id int not null,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    seguidor_id int not null,
    FOREIGN KEY (seguidor_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    PRIMARY KEY(usuario_id, seguidor_id)
)

DROP TABLE IF EXISTS publicacoes;

CREATE TABLE publicacoes(
    Id int auto_increment primary key,
    titulo VARCHAR(50) NOT NULL,
    conteudo VARCHAR(300) NOT NULL,
    autor_id int not null,
    FOREIGN KEY (autor_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    criadaEm timestamp default current_timestamp()
)

DROP TABLE IF EXISTS curtidas;

CREATE TABLE curtidas(
    usuario_id int not null,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    publicacao_id int not null,
    FOREIGN KEY (publicacao_id) REFERENCES publicacoes(id) ON DELETE CASCADE,
    PRIMARY KEY(usuario_id, publicacao_id)
)