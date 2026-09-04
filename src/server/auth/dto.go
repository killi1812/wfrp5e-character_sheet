package auth

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenDto struct {
	AccessToken string `json:"accessToken"`
}
