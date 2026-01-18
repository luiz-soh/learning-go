package models

type CreateUsuarioRequest struct {
	Nome  string `json:"nome,omitempty"`
	Nick  string `json:"nick,omitempty"`
	Email string `json:"email,omitempty"`
	Senha string `json:"senha,omitempty"`
}

type UpdateUsuarioRequest struct {
	Nome  string `json:"nome,omitempty"`
	Nick  string `json:"nick,omitempty"`
	Email string `json:"email,omitempty"`
}
