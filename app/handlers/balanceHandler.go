package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"user_balance_service/app/core"
	"user_balance_service/app/models"
	"user_balance_service/app/utils"
	"user_balance_service/app/utils/enums"
)

func GetUserBalance(gContext *gin.Context) {
	var getUserBalance models.GetUserBalance
	if err := gContext.ShouldBindJSON(&getUserBalance); err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, "Неправильный json формат")
		return
	}

	var returnUserBalance models.ReturnUserBalance
	returnUserBalance.UserBalance = GetUserBalanceFromDb(getUserBalance.IdUser)
	gContext.JSON(enums.HttpStatusOK, returnUserBalance)
}

func AddUserBalance(gContext *gin.Context) {
	var addUserBalance models.AddUserBalance

	if err := gContext.ShouldBindJSON(&addUserBalance); err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, "Неправильный json формат")
		return
	}

	if addUserBalance.AddableBalance <= 0 {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest,
			"Баланс должен пополняться на положительную сумму")
		return
	}

	err := addUserBalanceInDb(addUserBalance.IdUser, addUserBalance.AddableBalance)
	if err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(gContext, enums.HttpStatusOK)
}

//TODO: перенести ниженаходящиеся функции в отдельное место, т к они не хендлеры

func GetUserBalanceFromDb(idUser int) float32 {
	var balance float32
	sql := "select balance from users_balance where id_user = $1"
	core.Db.QueryRow(context.Background(), sql, idUser).Scan(&balance)
	return balance
}

func UpdateUserBalanceInDbWithTx(tx pgx.Tx, idUser int, newBalance float32) error {
	sql := `
		update users_balance
		set balance = $1
		where id_user = $2
	`
	_, err := tx.Exec(context.Background(), sql, newBalance, idUser)
	return err
}

func addUserBalanceInDb(IdUser int, AddableBalance float32) error {
	sql := `
		insert into users_balance (id_user, balance)
		values ($1, $2)
		on conflict (id_user)
		do update set balance = users_balance.balance + $2;
	`
	_, err := core.Db.Exec(context.Background(), sql, IdUser, AddableBalance)
	return err
}
