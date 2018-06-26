package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/asdine/storm"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Product type
type Product struct {
	ID          int    `form:"id" json:"id" storm:"id,increment"`
	Title       string `form:"title" json:"title" binding:"required"`
	Description string `form:"description" json:"description" binding:"required"`
	Category    string `form:"category" json:"category" binding:"required"`
	Status      string `form:"status" json:"status"`
	Images      string `form:"images" json:"images"`
	Tags        string `form:"tags" json:"tags"`
	CreatedAt   string `form:"createdAt" json:"createdAt"` //date
}

//Database middleware handler
func Database(db *storm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

//ProductCollectionHandler collection handler
func ProductCollectionHandler(c *gin.Context) {
	db := c.MustGet("DB").(*storm.DB)
	var products []Product
	err := db.All(&products, storm.Limit(100), storm.Reverse())
	if err != nil {
		log.Fatal(err)
	}

	c.Header("Content-Type", "application/hal+json")
	c.JSON(http.StatusOK, gin.H{
		"_links": gin.H{
			"self": gin.H{
				"href": "/products",
			},
		},
		"_embedded": gin.H{
			"product": products,
		},
	})
}

//ProductCreateHandler create product handler
func ProductCreateHandler(c *gin.Context) {
	db := c.MustGet("DB").(*storm.DB)
	var form Product

	if err := c.Bind(&form); err == nil {
		product := Product{Title: form.Title,
			Description: form.Description,
			Category:    form.Category,
			Status:      "pending",
			Images:      form.Images,
			Tags:        form.Tags,
			CreatedAt:   ""}
		err := db.Save(&product)
		if err != nil {
			c.Header("Content-Type", "application/problem+json")
			c.JSON(http.StatusOK, gin.H{"type": "error"})
		} else {
			c.Header("Content-Type", "application/hal+json")
			c.JSON(http.StatusOK, product)
		}
	} else {
		c.Header("Content-Type", "application/problem+json")
		fmt.Println(form)
		c.JSON(http.StatusBadRequest, gin.H{"title": "Bad Request", "detail": err.Error()})
	}
}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(http.StatusOK, struct{}{})
}

func main() {

	db, err := storm.Open("./db.json")
	if err != nil {
		fmt.Println("Error with local db : ")
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(Database(db))
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	router.OPTIONS("/", preflight)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	product := router.Group("/products")
	{
		//Retrieve products
		product.GET("/", ProductCollectionHandler)

		//Retrieve products
		product.POST("/", ProductCreateHandler)
		//TODO
	}

	router.Run(":5000")
}
