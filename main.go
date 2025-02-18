package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // or docker service name if running
	port     = 5432         // default PostgreSQL port
	database = "mydatabase" // as defined in docker-compose.yml
	username = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			// IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			// ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful: true, // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Failed to connect database.")

	}

	fmt.Printf("Connect successful.")
	fmt.Print(db)

	db.AutoMigrate(&Book{})

	app := fiber.New()
	app.Get("/Books", func(c fiber.Ctx) error {
		return c.JSON(getBooks(db))
	})
	app.Get("/Book/:id", func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		book, err := getBook(db, id)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		// if book == nil {
		// 	return c.SendStatus(fiber.StatusNotFound)
		// }

		return c.JSON(book)

	})
	app.Post("/Book", func(c fiber.Ctx) error {
		var book Book
		if err := c.Bind().JSON(&book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		if err := createBook(db, &book); err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.JSON(fiber.Map{
			"message": "Create successful",
		})

	})
	app.Put("/Book/:id", func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// var book Book
		updatedBook, err := getBook(db, id)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if err := c.Bind().JSON(&updatedBook); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		if err := updateBook(db, updatedBook); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.JSON(fiber.Map{
			"message": "update successful.",
		})
	})
	app.Delete("/book/:id", func(c fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if _, err := getBook(db, id); err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if err := deleteBook(db, id); err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.JSON(fiber.Map{
			"message": "Delete successful",
		})
	})

	app.Listen(":8080")

}

func getBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		log.Fatalf("Error get books: %v", result.Error)
	}
	return books
}
func getBook(db *gorm.DB, id int) (*Book, error) {
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		//log.Fatalf("Error get book: %v", result.Error)
		return nil, result.Error
	}

	return &book, nil
}
func createBook(db *gorm.DB, book *Book) error {
	result := db.Create(&book)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func updateBook(db *gorm.DB, book *Book) error {
	result := db.Save(&book)
	//result := db.Model(&book).Updates(book)
	if result != nil {
		return result.Error
	}

	return nil
}
func deleteBook(db *gorm.DB, id int) error {
	book := new(Book)
	result := db.Delete(&book, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
