package application

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chyngyz-sydykov/go-web/application/handlers"
	"github.com/chyngyz-sydykov/go-web/infrastructure/config"
	"github.com/chyngyz-sydykov/go-web/infrastructure/db"
	"github.com/chyngyz-sydykov/go-web/infrastructure/logger"
	"github.com/chyngyz-sydykov/go-web/infrastructure/messagebroker"
	"github.com/chyngyz-sydykov/go-web/internal/book"
	"github.com/chyngyz-sydykov/go-web/internal/rating"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"

	pb "github.com/chyngyz-sydykov/go-web/proto/rating"
)

type App struct {
	BookHandler   handlers.BookHandler
	RatingHandler handlers.RatingHandler
}

func InitializeApplication() *App {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}

	logger := logger.NewLogger()
	db := InitializeDatabase()
	//defer db.Close()

	rabbitMqPublisher := InitializeRabbitMqPublisher(config, logger)
	//defer rabbitMqPublisher.Close()

	commonHandler := handlers.NewCommonHandler(logger)

	ratingClient := initializeRatingGrpcClient()
	ratingService := rating.NewRatingService(ratingClient, time.Duration(config.GrpcTimeoutDuration)*time.Second)
	ratingHandler := handlers.NewRatingHandler(ratingService, *commonHandler)

	bookService := book.NewBookService(db, rabbitMqPublisher, ratingService)
	bookHandler := handlers.NewBookHandler(*bookService, *commonHandler)

	app := &App{
		BookHandler:   *bookHandler,
		RatingHandler: *ratingHandler,
	}
	return app
}

func initializeRatingGrpcClient() pb.RatingServiceClient {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}

	conn, err := grpc.NewClient(config.RatingServiceServer+":"+config.RatingServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	//defer conn.Close()

	// Create the gRPC client for the RatingService.
	ratingClient := pb.NewRatingServiceClient(conn)
	return ratingClient
}
func InitializeRabbitMqPublisher(config *config.Config, logger logger.LoggerInterface) messagebroker.MessageBrokerInterface {
	rabbitMQURL := "amqp://" + config.RabbitMqUser + ":" + config.RabbitMqPassword + "@" + config.RabbitMqContainerName + ":5672/"
	publisher, err := messagebroker.NewRabbitMQPublisher(rabbitMQURL, config.RabbitMqQueueName)

	if err != nil {
		err = fmt.Errorf("failed to initialize message publisher: %v", err)
		logger.LogError(http.StatusInternalServerError, err)
	} else {
		publisher.InitializeMessageBroker()
	}
	return publisher
}

func InitializeDatabase() *gorm.DB {
	dbConfig, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("Could not load database config: %v", err)
	}
	dbInstance, err := db.InitializeDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Coult not initialize db connection %v", err)
	}
	db.Migrate()
	return dbInstance

}
