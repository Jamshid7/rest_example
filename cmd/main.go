package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

var db *gorm.DB
var err error

type Card struct {
	ID         uint   `json:"id"`
	HolderName string `json:"holdername"`
	ExpDate    string `json:"expdate"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("HOST")

	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.AutoMigrate(&Card{})

	r := gin.Default()
	r.GET("/cards/", GetCards)
	r.GET("/cards/:id", GetCard)
	r.POST("/cards", CreateCard)
	r.PUT("/cards/:id", UpdateCard)
	r.DELETE("/cards/:id", DeleteCard)

	r.Run(host)
}

func DeleteCard(c *gin.Context) {
	id := c.Params.ByName("id")
	var card Card

	d := db.Where("id = ?", id).Delete(&card)
	fmt.Println(d)

	c.JSON(200, gin.H{"id #" + id: "Deleted"})
}

func UpdateCard(c *gin.Context) {
	var card Card
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&card).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	c.BindJSON(&card)
	db.Save(&card)
	c.JSON(200, card)
}

func CreateCard(c *gin.Context) {
	var card Card
	c.BindJSON(&card)

	db.Create(&card)
	c.JSON(200, card)
}

func GetCard(c *gin.Context) {
	id := c.Params.ByName("id")
	var card Card
	if err := db.Where("id= ?", id).First(&card).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, card)
	}
}

func GetCards(c *gin.Context) {
	var cards []Card
	if err := db.Find(&cards).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, cards)
	}

}
