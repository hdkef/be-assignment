package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hdkef/be-assignment/pkg/logger"
	"github.com/hdkef/be-assignment/services/transaction/config"
	deliveryConsumer "github.com/hdkef/be-assignment/services/transaction/internal/delivery/consumer"
	deliveryhttp "github.com/hdkef/be-assignment/services/transaction/internal/delivery/http"
	"github.com/sirupsen/logrus"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
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
		Publisher:      publisher,
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

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{cfg.WebDomain},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders: append([]string{"content-type"},
			supertokens.GetAllCORSHeaders()...),
		AllowCredentials: true,
	}))

	var log = logrus.New()
	router.Use(logger.LoggingMiddleware(log))

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
	} else {
		router.POST("/withdraw", verifySessionMiddleware(nil), handlerHttp.Withdraw)
		router.POST("/send", verifySessionMiddleware(nil), handlerHttp.Send)
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

// verifySessionMiddleware adapts the Supertoken VerifySession function to work as a Gin middleware
func verifySessionMiddleware(options *sessmodels.VerifySessionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sessionVerified bool
		var err error

		session.VerifySession(options, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Session verified successfully
			sessionVerified = true
			// Update the Gin context with the new request that has the session information
			c.Request = c.Request.WithContext(r.Context())
		})).ServeHTTP(c.Writer, c.Request)

		if !sessionVerified {
			// Session verification failed
			err = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			if err != nil {
				// If there's an error sending the response, log it
				fmt.Printf("Error sending unauthorized response: %v\n", err)
			}
			return
		}

		// Continue with the next middleware/handler
		c.Next()
	}
}
