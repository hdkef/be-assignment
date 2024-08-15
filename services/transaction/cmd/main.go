package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hdkef/be-assignment/pkg/logger"
	"github.com/hdkef/be-assignment/pkg/middleware"
	"github.com/hdkef/be-assignment/services/transaction/config"
	deliveryConsumer "github.com/hdkef/be-assignment/services/transaction/internal/delivery/consumer"
	deliveryhttp "github.com/hdkef/be-assignment/services/transaction/internal/delivery/http"
	"github.com/sirupsen/logrus"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"

	"github.com/hdkef/be-assignment/services/transaction/internal/repository"
	"github.com/hdkef/be-assignment/services/transaction/internal/usecase"

	"github.com/hdkef/be-assignment/services/transaction/internal/service"
)

func main() {

	db := config.InitDB()
	defer db.Close()
	cfg := config.InitTransactionConfig()

	rbmqConn := config.InitRBMQ()
	defer rbmqConn.Close()

	accBalanceRepo := repository.NewAccountBalanceRepo(db)
	trxRepo := repository.NewTransactionLogsRepo(db)
	scheduleRepo := repository.NewScheduleRepo(db)
	queueRepo := repository.NewQueueRepo(db)

	publisher := service.NewPublisher(rbmqConn)

	accBalanceUC := usecase.AccountBalanceUC{
		UoW: repository.UnitOfWorkImplementor{
			Db: db,
		},
		AccBalanceRepo: accBalanceRepo,
	}

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

	ch, err := rbmqConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	apiBasePath := "/auth"
	websiteBasePath := "/auth"
	err = supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			// These are the connection details of the app you created on supertokens.com
			ConnectionURI: cfg.SuperTokenUrl,
			APIKey:        cfg.SuperTokenAPIKey,
		},
		AppInfo: supertokens.AppInfo{
			AppName:         cfg.AppName,
			APIDomain:       cfg.APIDomain,
			WebsiteDomain:   cfg.WebDomain,
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			session.Init(nil),
		},
	})

	if err != nil {
		panic(err.Error())
	}

	router := gin.New()

	var log = logrus.New()
	router.Use(logger.LoggingMiddleware(log))

	corsOrigin := os.Getenv("CORS_ALLOW_ORIGIN")
	corsHeader := os.Getenv("CORS_ALLOW_HEADER")

	debugMode := os.Getenv("DEBUG_MODE") == "Y" || os.Getenv("DEBUG_MODE") == ""
	if debugMode {
		corsOrigin = "*"
		corsHeader = "*"
	}

	router.Use(middleware.CORSMiddleware(corsOrigin, corsHeader))

	// Adding the SuperTokens middleware
	router.Use(func(c *gin.Context) {
		supertokens.Middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				c.Next()
			})).ServeHTTP(c.Writer, c.Request)
		// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
		c.Abort()
	})

	handlerConsumer := &deliveryConsumer.ConsumerDelivery{
		Conn:         rbmqConn,
		AccBalanceUC: &accBalanceUC,
	}
	handlerHttp := &deliveryhttp.HttpHandler{
		TransactionUC: &trxUC,
	}

	if cfg.DEBUGMODE == "Y" {
		router.POST("/withdraw", handlerHttp.Withdraw)
		router.POST("/send", handlerHttp.Send)
		router.POST("/autodebet", handlerHttp.SetAutodebet)
	} else {
		router.POST("/withdraw", middleware.VerifySessionMiddleware(nil), handlerHttp.Withdraw)
		router.POST("/send", middleware.VerifySessionMiddleware(nil), handlerHttp.SetAutodebet)
		router.POST("/autodebet", middleware.VerifySessionMiddleware(nil), handlerHttp.SetAutodebet)

	}

	// delivery consumer
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go handlerConsumer.SignUpEvent()
	go handlerConsumer.AccountCreatedEvent()

	err = router.Run(fmt.Sprintf(":%s", cfg.AppPort))
	if err != nil {
		panic(err)
	}

	<-sigs
}
