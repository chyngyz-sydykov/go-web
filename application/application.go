package application

import (
	"log"
	"time"

	"github.com/chyngyz-sydykov/go-web/application/handlers"
	"github.com/chyngyz-sydykov/go-web/infrastructure/config"
	"github.com/chyngyz-sydykov/go-web/infrastructure/db"
	"github.com/chyngyz-sydykov/go-web/infrastructure/logger"
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

	db := initializeDatabase()

	ratingClient := initializeRatingGrpcClient()

	logger := logger.NewLogger()

	commonHandler := handlers.NewCommonHandler(logger)

	ratingService := rating.NewRatingService(ratingClient, time.Duration(config.GrpcTimeoutDuration)*time.Second)
	ratingHandler := handlers.NewRatingHandler(*commonHandler)

	bookService := book.NewBookService(db, ratingService)
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

func initializeDatabase() *gorm.DB {
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
