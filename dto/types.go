package dto

type AuthUserRegisterReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

type AuthUserLoginReq struct {
	LoginOrEmail string `json:"loginOrEmail" validate:"required"`
	Password     string `json:"password" validate:"required"`
}
