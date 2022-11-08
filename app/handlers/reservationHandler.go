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

func ReserveUserBalance(gContext *gin.Context) {
	var reserveUserBalance models.ReserveUserBalance

	if err := gContext.ShouldBindJSON(&reserveUserBalance); err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, "Неправильный json формат")
		return
	}

	if reserveUserBalance.ReservedMoney <= 0 {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest,
			"Нельзя зарезервировать на отрицательную сумму")
		return
	}

	var currentBalance = GetUserBalanceFromDb(reserveUserBalance.IdUser)
	var newBalance = currentBalance - reserveUserBalance.ReservedMoney
	if newBalance < 0 {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, "На счету пользователя недостаточно денежных средств")
		return
	}

	tx, err := core.Db.Begin(context.Background())
	if err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, err.Error())
		return
	}

	err = UpdateUserBalanceInDbWithTx(tx, reserveUserBalance.IdUser, newBalance)
	if err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, err.Error())
		return
	}

	idReservation, err := reserveUserBalanceInDbWithTx(tx, reserveUserBalance)
	if err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, err.Error())
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, err.Error())
		return
	}

	var ReturnReserveInfo models.ReturnReserveInfo
	ReturnReserveInfo.IdReservation = idReservation
	ReturnReserveInfo.Result = "success"
	gContext.JSON(enums.HttpStatusOK, ReturnReserveInfo)
}

func AddRevenueToCompany(gContext *gin.Context) {
	var addRevenueToCompany models.AddRevenueToCompany

	if err := gContext.ShouldBindJSON(&addRevenueToCompany); err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, "Неправильный json формат")
		return
	}

	tx, err := core.Db.Begin(context.Background())
	if err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, err.Error())
		return
	}

	reserveUserBalance, err := deleteReservationAndReturnRevenueWithTx(tx, addRevenueToCompany.IdReservation)
	if err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, "Передан неверный id_reservation")
		return
	}

	err = addRevenueToCompanyInDb(tx, reserveUserBalance)
	if err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, "Передан неверный id_reservation")
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		utils.SendErrResponse(gContext, enums.HttpStatusBadRequest, err.Error())
		return
	}

	utils.SendSuccessResponse(gContext, enums.HttpStatusOK)
}

// TODO: перенести ниженаходящиеся функции в отдельное место, т к они не хендлеры

func reserveUserBalanceInDbWithTx(tx pgx.Tx, reserveUserBalance models.ReserveUserBalance) (int, error) {
	sql := `
		insert into reservation_money (id_user, id_service, id_order, reserved_money)
		values ($1, $2, $3, $4)
		returning id_reservation
	`
	var idReservation int
	err := tx.QueryRow(context.Background(), sql,
		reserveUserBalance.IdUser,
		reserveUserBalance.IdService,
		reserveUserBalance.IdOrder,
		reserveUserBalance.ReservedMoney).Scan(&idReservation)
	return idReservation, err
}

func deleteReservationAndReturnRevenueWithTx(tx pgx.Tx, idReservation int) (models.ReserveUserBalance, error) {
	sql := `
		delete from reservation_money
		where id_reservation = $1
		returning id_user, id_service, id_order, reserved_money
	`
	var reserveUserBalance models.ReserveUserBalance
	err := tx.QueryRow(context.Background(), sql, idReservation).Scan(
		&reserveUserBalance.IdUser,
		&reserveUserBalance.IdService,
		&reserveUserBalance.IdOrder,
		&reserveUserBalance.ReservedMoney)
	return reserveUserBalance, err
}

func addRevenueToCompanyInDb(tx pgx.Tx, reserveUserBalance models.ReserveUserBalance) error {
	sql := `
		insert into company_revenue (id_user, id_service, id_order, revenue)
		values ($1, $2, $3, $4)
	`
	_, err := tx.Exec(context.Background(), sql,
		reserveUserBalance.IdUser,
		reserveUserBalance.IdService,
		reserveUserBalance.IdOrder,
		reserveUserBalance.ReservedMoney)
	return err
}
