package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type URL struct {
	ID         int       `json:"id"`
	URL        string    `json:"url"`
	Domain     string    `json:"domain"`
	Created_At time.Time `json:"date"`
}

var (
	DB *sql.DB
)

func main() {
	createDBConnection()
	defer DB.Close()
	r := gin.Default()
	setupRouters(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func setupRouters(r *gin.Engine) {

	r.POST("/url", PostUrl)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "MESSAGE": "page not found"})

	})

}

// POST
func PostUrl(c *gin.Context) {

	reqBody := URL{}
	err := c.Bind(&reqBody)
	if err != nil {

		res := gin.H{
			"error": "invalid request body",
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, res)

		return
	}

	parsedURL, err := url.Parse(reqBody.URL)
	if err != nil {
		// Handle invalid URL error
		// ...
	}
	reqBody.Domain = strings.TrimPrefix(parsedURL.Hostname(), "www.")
	reqBody.Created_At = time.Now()

	fmt.Println(reqBody)
	res, err := DB.Exec(`INSERT INTO "url" ("url", "domain_name", "created_at")
	VALUES ($1, $2, $3)`, reqBody.URL, reqBody.Domain, reqBody.Created_At)
	if err != nil {
		fmt.Println("err inserting data: ", err)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return

	}
	lastInsID, err := res.LastInsertId()
	fmt.Println("errr: ", err)
	reqBody.ID = int(lastInsID)
	fmt.Println("res: ", lastInsID)
	c.JSON(http.StatusOK, reqBody)
	c.Writer.Header().Set("Content-Type", "application/jason")
}
