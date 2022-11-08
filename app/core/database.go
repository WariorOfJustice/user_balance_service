package core

import "github.com/jackc/pgx/v5"

var Db *pgx.Conn

func SetDb(db_new *pgx.Conn) {
	Db = db_new
}
