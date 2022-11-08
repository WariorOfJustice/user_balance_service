package models

type GetUserBalance struct {
	IdUser int `json:"id_user" binding:"required"`
}

type ReturnUserBalance struct {
	UserBalance float32 `json:"user_balance"`
}

type AddUserBalance struct {
	IdUser         int     `json:"id_user" binding:"required"`
	AddableBalance float32 `json:"addable_balance" binding:"required"`
}
