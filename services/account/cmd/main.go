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
	"github.com/hdkef/be-assignment/services/account/config"
	deliveryConsumer "github.com/hdkef/be-assignment/services/account/internal/delivery/consumer"
	deliveryhttp "github.com/hdkef/be-assignment/services/account/internal/delivery/http"
	"github.com/hdkef/be-assignment/services/account/internal/repository"
	"github.com/hdkef/be-assignment/services/account/internal/service"
	"github.com/hdkef/be-assignment/services/account/internal/usecase"
	"github.com/sirupsen/logrus"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {

	// init repo
	db := config.InitDB()
	defer db.Close()

	cfg := config.InitAccountConfig()
	rbmqConn := config.InitRBMQ()
	defer rbmqConn.Close()

	userRepo := repository.NewUserRepo(db)
	userAddressRepo := repository.NewUserAddressRepo(db)
	accRepo := repository.NewAccountRepo(db)
	historyRepo := repository.NewHistoryRepo(db)

	// init service
	publisher := service.NewPublisher(rbmqConn)

	// init usecase
	userUC := usecase.UserUC{
		UoW: repository.UnitOfWorkImplementor{
			Db: db,
		},
		Publisher:       publisher,
		UserRepo:        userRepo,
		UserAddressRepo: userAddressRepo,
		AccountRepo:     accRepo,
	}

	accountUC := usecase.AccountUC{
		UoW: repository.UnitOfWorkImplementor{
			Db: db,
		},
		HistoryRepo: historyRepo,
		AccountRepo: accRepo,
		Publisher:   publisher,
	}

	// init delivery

	ch, err := rbmqConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	handlerHttp := &deliveryhttp.HttpHandler{
		UserUsecase: &userUC,
		AccUsecase:  &accountUC,
	}
	handlerConsumer := &deliveryConsumer.ConsumerDelivery{
		Conn:      rbmqConn,
		UserUC:    &userUC,
		AccountUC: &accountUC,
	}

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
			emailpassword.Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: handlerHttp.SuperTokenSignUp,
				},
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "name",
						},
						{
							ID: "dateOfBirth",
						},
						{
							ID: "job",
						},
						{
							ID: "address",
						},
						{
							ID: "district",
						},
						{
							ID: "city",
						},
						{
							ID: "province",
						},
						{
							ID: "country",
						},
						{
							ID: "accCurrency",
						},
						{
							ID: "accDesc",
						},
						{
							ID: "zip",
						},
					},
				},
			}),
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

	// delivery http
	if cfg.DEBUGMODE == "Y" {
		router.GET("/history", middleware.VerifySessionMiddleware(nil), handlerHttp.GetHistory)
		router.POST("/additional-account", middleware.VerifySessionMiddleware(nil), handlerHttp.CreateAccount)
	} else {
		router.GET("/history", middleware.VerifySessionMiddleware(nil), handlerHttp.GetHistory)
		router.POST("/additional-account", middleware.VerifySessionMiddleware(nil), handlerHttp.CreateAccount)
	}

	// delivery consumer
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go handlerConsumer.TransactionCreatedEvent()

	err = router.Run(fmt.Sprintf(":%s", cfg.AppPort))
	if err != nil {
		panic(err)
	}

	<-sigs
}
