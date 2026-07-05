package domain

type Token struct {
	accessToken  string
	refreshToken string
}

func NewToken(accessToken, refreshToken string) *Token {
	return &Token{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

func (t *Token) AccessToken() string  { return t.accessToken }
func (t *Token) RefreshToken() string { return t.refreshToken }
