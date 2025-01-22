package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/chyngyz-sydykov/go-web/application"
	"github.com/chyngyz-sydykov/go-web/infrastructure/config"
	"github.com/chyngyz-sydykov/go-web/infrastructure/db/models"
	"github.com/chyngyz-sydykov/go-web/internal/book"
	"github.com/stretchr/testify/mock"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (suite *IntegrationSuite) TestBookUpdating_ShouldSuccessfullyPublishToRabbitMQ() {

	// Arrange
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}

	db := application.InitializeDatabase()
	publisher := application.InitializeRabbitMqPublisher(config, &LoggerMock{})

	ratingServiceMock := &RatingServiceMock{}

	bookService := book.NewBookService(db, publisher, ratingServiceMock)

	testAuthor := models.Author{Firstname: "John", Lastname: "Doe"}
	suite.db.Create(&testAuthor)

	publishedAt := time.Now()
	testBook := models.Book{Title: "test book", ICBN: "123423ASDF", PublishedAt: &publishedAt, AuthorId: int64(testAuthor.ID)}
	suite.db.Create(&testBook)

	// Act
	var updatedBook *models.Book
	updatedBook, err = bookService.Update(int(testBook.ID), models.Book{
		Title: "updated title",
	})
	if err != nil {
		log.Fatalf("Could not update book: %v", err)
	}

	// consume message from rabbitmq -> confirm it is correct one
	rabbitMQURL := "amqp://" + config.RabbitMqUser + ":" + config.RabbitMqPassword + "@" + config.RabbitMqContainerName + ":5672/"
	conn, ch := newConsumer(rabbitMQURL, config.RabbitMqQueueName)
	defer closeConsumer(conn, ch)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		fmt.Println("Waiting for 2 seconds")
		time.Sleep(2 * time.Second)
		closeConsumer(conn, ch)
		publisher.Close()
		defer wg.Done()
	}()

	msgs, err := ch.Consume(
		config.RabbitMqQueueName, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		var bookMessage book.BookMessage
		if err := json.Unmarshal(msg.Body, &bookMessage); err != nil {
			log.Fatalf("failed to unmarshal message: %v", err)
		}

		fmt.Println("bookMessage: ", bookMessage)

		// Assert
		suite.Suite.Assert().Equal("bookUpdated", bookMessage.Event)
		suite.Suite.Assert().Equal(testBook.ID, updatedBook.ID)
	}

	wg.Wait()

	suite.db.Unscoped().Delete(&models.Book{}, testBook.ID)
	suite.db.Unscoped().Delete(&models.Author{}, testAuthor.ID)
}
func (suite *IntegrationSuite) TestRabbitMQ_ShouldLogErrorMessageIfCannotConnectToRabbitMQ() {
	// Arrange
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not config: %v", err)
	}
	config.RabbitMqContainerName = "invalid"
	loggerMock := &LoggerMock{}

	// Assert
	loggerMock.On("LogError",
		mock.Anything,
		mock.MatchedBy(func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "failed to initialize message publisher")
		}),
	).Once()

	//act
	_ = application.InitializeRabbitMqPublisher(config, loggerMock)
}

func newConsumer(url, queueName string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}

	// Declare the queue
	_, err = ch.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}
	return conn, ch
}

func closeConsumer(conn *amqp.Connection, ch *amqp.Channel) {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
}

type LoggerMock struct {
	mock.Mock
}

func (m *LoggerMock) LogError(statusCode int, err error) {
	m.Called(statusCode, err)
}
