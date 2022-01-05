package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type book struct {
	ID        uint       `json:"id"`
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

var books = []book{
	{ID: 1, Name: "Networking and Kubernetes", IBSN: "978-1492081654", Author: "James Strong", Publisher: "O'Rielly", Category: "DevOps"},
}

var deletedBooks = []book{}

// postusers adds an user from JSON received in the request body.
func postBooks(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}
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
	if err := db.Where("id = ?", id).First(&books).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.IndentedJSON(http.StatusCreated, getbookbyid)
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
}

// getUsers responds with the list of all Users as JSON.
func getBooks(c *gin.Context) {
	db.Create(&books)
	c.IndentedJSON(http.StatusOK, books)
}

func main() {
	// db, _ = gorm.Open("mysql", "root:'PaS$Wd'@tcp(127.0.0.1:3306)/library?charset=utf8&parseTime=True&loc=Local")
	dsn := "lib:test123321@tcp(10.19.0.8:3306)/library?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&books)

	db.AutoMigrate(&book{})
	router := gin.Default()
	router.GET("/v1/books", getBooks)
	router.POST("/v1/books", postBooks)
	router.GET("/v1/books/:id", getBookByID)
	// router.DELETE("/v1/books/:id", deleteBook)

	router.Run("0.0.0.0:8080")
}
