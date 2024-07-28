package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums []album

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected")

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", addAlbum)
	router.GET("/albums/:id", getAlbumByID)

	router.Run("localhost:8080")
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	var alb album

	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Could not find album"})
			return
		}
	}
	c.IndentedJSON(http.StatusOK, alb)
}

func addAlbum(c *gin.Context) {
	var alb album

	if err := c.ShouldBindBodyWithJSON(&alb); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", &alb.Title, &alb.Artist, &alb.Price)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	id, err := result.LastInsertId()

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"received": id})
}

// func getAlbumBySlug(c *gin.Context) {
// 	slug := c.Param("slug")
// 	rows, err := db.Query("SELECT * FROM album where slug = ?", slug)

// 	if err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, err)
// 		return
// 	}

// 	defer rows.Close()

// 	// loop through rows, using Scan to assign column data to struct fields.
// 	for rows.Next() {
// 		var alb album
// 		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
// 			c.IndentedJSON(http.StatusInternalServerError, err)
// 			return
// 		}
// 		albums = append(albums, alb)
// 	}
// 	if err := rows.Err(); err != nil {
// 		c.IndentedJSON(http.StatusInternalServerError, err)
// 		return
// 	}
// 	c.IndentedJSON(http.StatusFound, albums)
// }

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// func getAlbumByID(c *gin.Context) {
// 	id := c.Param("id")

// 	for _, a := range albums {
// 		if a.ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
// }
