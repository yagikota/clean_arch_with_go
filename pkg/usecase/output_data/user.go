package outputdata

type UserCreateResponse struct {
	Token string `json:"token"`
}

type UserGetResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	HighScore int32  `json:"highScore"`
	Coin      int32  `json:"coin"`
}
