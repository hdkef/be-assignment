package main

import (
	"context"
	"fmt"

	"github.com/hdkef/be-assignment/services/transaction/config"
	"github.com/hdkef/be-assignment/services/transaction/internal/repository"
	"github.com/hdkef/be-assignment/services/transaction/internal/service"
	"github.com/hdkef/be-assignment/services/transaction/internal/usecase"
)

func main() {

	db := config.InitDB()
	defer db.Close()

	rbmqConn := config.InitRBMQ()
	defer rbmqConn.Close()

	accBalanceRepo := repository.NewAccountBalanceRepo(db)
	trxRepo := repository.NewTransactionLogsRepo(db)
	scheduleRepo := repository.NewScheduleRepo(db)
	queueRepo := repository.NewQueueRepo(db)

	publisher := service.NewPublisher(rbmqConn)

	trxUC := usecase.TransactionUC{
		UoW: repository.UnitOfWorkImplementor{
			Db: db,
		},
		AccBalanceRepo: accBalanceRepo,
		TrxLogsRepo:    trxRepo,
		ScheduleRepo:   scheduleRepo,
		Publisher:      publisher,
		QueueRepo:      queueRepo,
	}

	// run usecase immediately
	ctx := context.Background()
	err := trxUC.ProcessAutodebetDaily(ctx)
	if err != nil {
		fmt.Println("Error running usecase:", err)
	}
	err = trxUC.ProcessQueue(ctx)
	if err != nil {
		fmt.Println("Error running usecase:", err)
	}

	fmt.Println("CRON FINISHED")
}
