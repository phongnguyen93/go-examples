package main

import (
	"net/http"
	"strconv"

	"example/database"

	"github.com/gin-gonic/gin"
)

func getAlbums(c *gin.Context) {
	queryAlbums, err := database.GetAllAlbums()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, queryAlbums)
}

func postAlbums(c *gin.Context) {
	var newAlbum database.Album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	insertRow, err := database.InsertNewAlbum(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	c.IndentedJSON(http.StatusCreated, insertRow)
}

func getAlbumByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}
	queryAlbum, err := database.GetAlbumByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	}
	c.IndentedJSON(http.StatusOK, queryAlbum)
}

func main() {
	database.Init()
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.Run("localhost:8080")
}
