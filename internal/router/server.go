package router

import (
	"context"
	"log"
	"net/http"

	"github.com/payments-core/logs"
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
	ctx := context.Background()
	logs.CtxInfo(ctx, "[Server] [Start] opening db , driver:%v dsn:%v", cfg.DBDriver, cfg.DBDSN)
	conn, err := db.Open(cfg)
	if err != nil {
		logs.CtxInfo(ctx, "[Server] [Start] error opening db , err:%v", err)
		log.Fatal(err)
	}
	logs.CtxInfo(ctx, "[Server] [Start] db open success")
	logs.CtxInfo(ctx, "[Server] [Start] migrating db")
	if err := db.Migrate(cfg, conn); err != nil {
		logs.CtxInfo(ctx, "[Server] [Start] error migrating db , err:%v", err)
		log.Fatal(err)
	}
	logs.CtxInfo(ctx, "[Server] [Start] migrate success")
	srv := NewServer(cfg, conn)
	logs.CtxInfo(ctx, "[Server] [Start] http listen on addr:%v", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, srv.Router()); err != nil {
		logs.CtxInfo(ctx, "[Server] [Start] http server error , err:%v", err)
		log.Fatal(err)
	}
}
