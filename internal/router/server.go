package router

import (
	"github.com/payments-core/internal/app"
	"github.com/payments-core/internal/config"
	"github.com/payments-core/internal/repo"
	"github.com/payments-core/internal/repo/postgres"
	sqliterepo "github.com/payments-core/internal/repo/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	cfg config.Config
	db  *gorm.DB
	svc *app.TransferService
}

func NewServer(cfg config.Config, db *gorm.DB) *Server {
	var ar repo.AccountRepository
	var tr repo.TransactionRepository

	switch cfg.DBDriver {
	case "postgres":
		ar = postgres.NewAccountRepo()
		tr = postgres.NewTransactionRepo()
	case "sqlite":
		ar = sqliterepo.NewAccountRepo()
		tr = sqliterepo.NewTransactionRepo()
	default:
		ar = postgres.NewAccountRepo()
		tr = postgres.NewTransactionRepo()
	}

	svc := app.NewTransferService(db, cfg.DBDriver, ar, tr)
	return &Server{cfg: cfg, db: db, svc: svc}
}
