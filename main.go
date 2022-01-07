package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pooriaghaedi/Go-Gin-Library/config"

	// "gorm.io/driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

var err error

type book struct {
	gorm.Model
	Name      string `json:"name"`
	IBSN      string `json:"ibsn"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Category  string `json:"category"`
	Photo     string `json:"photo"`
}

// var books = []book{
// 	{ID: 1, Name: "Networking and Kubernetes", IBSN: "978-1492081654", Author: "James Strong", Publisher: "O'Rielly", Category: "DevOps"},
// }

var deletedBooks = []book{}

// postBooks adds a book from JSON received in the request body.
func postBooks(c *gin.Context) {
	var Books book
	c.BindJSON(&Books)

	db.Create(&Books)
	c.JSON(200, Books)
}

// deleteBook is used to delete a book with its ID.
func deleteBook(c *gin.Context) {
	id := c.Param("id")
	var Books book
	if err := db.Where("id = ?", id).First(&Books).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		db.Where("id = ?", id).Delete(&Books)
		c.JSON(200, gin.H{"id #" + id: "deleted"})
	}

}

// getBookByID is used to get a book with its ID.
func getBookByID(c *gin.Context) {
	id := c.Param("id")
	var getbookbyid book
	if err := db.Where("id = ?", id).First(&getbookbyid).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.IndentedJSON(http.StatusOK, getbookbyid)
	}
}

// UpdateBooks is used to update an existing book.
func UpdateBooks(c *gin.Context) {

	var Book book
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&Book).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&Book)

	db.Save(&Book)
	c.JSON(200, Book)

}

// getBooks responds with the list of all Books as JSON.
func getBooks(c *gin.Context) {
	var Books []book
	if err := db.Find(&Books).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, Books)
	}
}

func main() {

	router := gin.Default()
	router.GET("/v1/books/", getBooks)
	router.POST("/v1/books/", postBooks)
	router.GET("/v1/books/:id", getBookByID)
	router.DELETE("/v1/books/:id", deleteBook)
	router.PUT("/v1/books/:id", UpdateBooks)

	router.Run("0.0.0.0:8080")
}

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&book{})
}
