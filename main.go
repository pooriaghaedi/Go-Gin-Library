package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pooriaghaedi/Go-Gin-Library/config"
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
type File struct {
	Name string `uri:"name" binding:"required"`
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

func UploadBookcover(c *gin.Context) {
	var Book book
	id := c.Params.ByName("id")
	// Source

	if err := db.Where("id = ?", id).First(&Book).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}
		if err := c.SaveUploadedFile(file, "public/"+file.Filename); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		} else {
			if err := db.Model(&Book).Where("id = ?", id).Update("photo", file.Filename).Error; err != nil {
				fmt.Println(err)
			}
			c.String(http.StatusOK, "File %s uploaded successfully with fields name=%s and id=%s.", file.Filename, id)
		}

	}
}

func Download(n string) (string, []byte, error) {
	dst := fmt.Sprintf("%s/%s", "public", n)
	b, err := ioutil.ReadFile(dst)
	if err != nil {
		return "", nil, err
	}
	m := http.DetectContentType(b[:512])

	return m, b, nil
}

func GetBookcover(c *gin.Context) {
	id := c.Params.ByName("id")
	var Book book
	if err := db.Where("id = ?", id).First(&Book).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	fmt.Println(Book.Photo)
	// var f File
	if err := c.ShouldBindUri(Book.Photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	m, cn, err := Download(Book.Photo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+Book.Photo)
	c.Data(http.StatusOK, m, cn)

}
func main() {

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.GET("/v1/books/", getBooks)
	router.POST("/v1/books/", postBooks)
	router.GET("/v1/books/:id", getBookByID)
	router.DELETE("/v1/books/:id", deleteBook)
	router.PUT("/v1/books/:id", UpdateBooks)
	router.PUT("/v1/upload/:id", UploadBookcover)
	router.GET("/v1/upload/:id", GetBookcover)

	router.Run("0.0.0.0:8080")
}

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&book{})
}
