package models

type ReserveUserBalance struct {
	IdUser        int     `json:"id_user" binding:"required"`
	IdService     int     `json:"id_service" binding:"required"`
	IdOrder       int     `json:"id_order" binding:"required"`
	ReservedMoney float32 `json:"reserved_money" binding:"required"`
}

type ReturnReserveInfo struct {
	IdReservation int    `json:"id_reservation"`
	Result        string `json:"result"`
}

type AddRevenueToCompany struct {
	IdReservation int `json:"id_reservation" binding:"required"`
}
