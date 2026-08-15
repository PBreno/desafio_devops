package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type response struct {
	Name string `json:"nome"`
	Time string `json:"horario"`
}

var data = []response{

	{Name: "Projeto Korp", Time: time.Now().String()[11:19]},
}

func getData(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, data)

}
func main() {

	router := gin.Default()
	router.GET("/", getData)

	router.Run("0.0.0.0:8080")
}
