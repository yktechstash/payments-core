package router

import (
	"log"
	"net/http"

	"github.com/payments-core/internal/app"
	"github.com/payments-core/internal/config"
	"github.com/payments-core/internal/db"
	"github.com/payments-core/internal/db/dal"
	"github.com/payments-core/internal/db/postgres"
	"github.com/payments-core/internal/db/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	cfg config.Config
	db  *gorm.DB
	svc *app.TransferService
}

func NewServer(cfg config.Config, db *gorm.DB) *Server {
	var ar dal.AccountsLedger
	var tr dal.TransactionsLedger

	switch cfg.DBDriver {
case "postgres":
	ar = postgres.NewAccountsLedger(db)
	tr = postgres.NewTransactionsLedger(db)
case "sqlite":
	ar = sqlite.NewAccountsLedger(db)
	tr = sqlite.NewTransactionsLedger(db)
default:
	ar = postgres.NewAccountsLedger(db)
	tr = postgres.NewTransactionsLedger(db)
}

	svc := app.NewTransferService(db, cfg.DBDriver, ar, tr)
	return &Server{cfg: cfg, db: db, svc: svc}
}

func Start() {
	cfg := config.Conf
	conn, err := db.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Migrate(cfg, conn); err != nil {
		log.Fatal(err)
	}
	srv := NewServer(cfg, conn)
	if err := http.ListenAndServe(cfg.Addr, srv.Router()); err != nil {
		log.Fatal(err)
	}
}
