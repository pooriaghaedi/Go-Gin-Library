package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pooriaghaedi/Go-Gin-Library/config"

	// "gorm.io/driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

var err error

type book struct {
	gorm.Model
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	IBSN      string    `json:"ibsn"`
	Author    string    `json:"author"`
	Publisher string    `json:"publisher"`
	Category  string    `json:"category"`
	Photo     string    `json:"photo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// var books = []book{
// 	{ID: 1, Name: "Networking and Kubernetes", IBSN: "978-1492081654", Author: "James Strong", Publisher: "O'Rielly", Category: "DevOps"},
// }

var deletedBooks = []book{}

// postusers adds an user from JSON received in the request body.
func postBooks(c *gin.Context) {
	var newBook book

	// if err := c.BindJSON(&newBook); err != nil {
	// 	return
	// }
	c.BindJSON(&newBook)

	db.Create(&newBook)

	c.IndentedJSON(http.StatusCreated, newBook)
}

// func deleteBook(c *gin.Context) {
// 	id := c.Param("id")
// 	var deleteBook book
// 	currentTime := time.Now()

// 	// if err := c.BindJSON(&deleteBook); err != nil {
// 	// 	fmt.Printf("error")
// 	// 	return

// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
// }
func getBookByID(c *gin.Context) {
	id := c.Param("id")
	var getbookbyid book
	if err := db.Where("id = ?", id).First(getbookbyid).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.IndentedJSON(http.StatusCreated, getbookbyid)
	}

	// c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// getUsers responds with the list of all Users as JSON.
func getBooks(c *gin.Context) {
	var Books []book
	db.Find(&Books)
	c.IndentedJSON(http.StatusAccepted, Books)

}

func main() {

	router := gin.Default()
	router.GET("/v1/books", getBooks)
	router.POST("/v1/books", postBooks)
	router.GET("/v1/books/:id", getBookByID)
	// router.DELETE("/v1/books/:id", deleteBook)

	router.Run("0.0.0.0:8080")
}

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&book{})
}
