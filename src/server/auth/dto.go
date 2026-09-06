package auth

type LoginDto struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
}

type TokenDto struct {
	AccessToken string `json:"accessToken"`
}
