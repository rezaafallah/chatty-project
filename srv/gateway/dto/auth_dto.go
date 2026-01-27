package dto

// RegisterReq
type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginReq
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRes
type LoginRes struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

// UserRes
type UserRes struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}