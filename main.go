package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Satellite struct {
	Name     string   `json:"name"`
	Distance float64  `json:"distance"`
	Message  []string `json:"message"`
	Position Point    `json:"position"`
}

var fleets = []Satellite{
	{
		Name:     "kenobi",
		Distance: 0,
		Message:  []string{},
		Position: Point{X: -500.0, Y: -200.0},
	},
	{
		Name:     "skywalker",
		Distance: 0,
		Message:  []string{},
		Position: Point{X: 100.0, Y: -100.0},
	},
	{
		Name:     "sato",
		Distance: 0,
		Message:  []string{},
		Position: Point{X: 500.0, Y: 100.0},
	},
}
var mode = "flex"

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"https://fernandopriet1997.github.io"},
		AllowMethods:  []string{"GET", "POST"},
		AllowHeaders:  []string{"content-type"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	router.POST("/topsecret", postTopSecret)
	router.GET("/topsecret", getTopSecret)
	router.GET("/topsecret_split", getTopSecretSplit)
	router.POST("/topsecret_split/:name", postTopSecretSplit)
	router.GET("/topsecret_split/:name", getSatellite)
	router.POST("/config-mode", setMode)
	router.GET("/config-mode", getMode)
	router.POST("/config/:name", postSetPosition)
	return router
}
func main() {
	r := setupRouter()
	r.Run(":8080")
}
func getTopSecret(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, fleets)
}
func postTopSecret(c *gin.Context) {

	type RequestBody struct {
		Satellites []Satellite `json:"satellites"`
	}
	var satsConf RequestBody

	if err := c.ShouldBindJSON(&satsConf); err != nil {
		fmt.Println(err)
		return
	}

	reslocation, err1 := GetLocation(satsConf.Satellites)
	resMessage, err2 := GetMessage(satsConf.Satellites, mode)

	if err1 == nil && err2 == nil {
		type Response struct {
			Position Point  `json:"position"`
			Message  string `json:"message"`
		}
		response := Response{
			Position: reslocation,
			Message:  resMessage,
		}

		c.IndentedJSON(http.StatusOK, response)
	} else {
		c.IndentedJSON(http.StatusNotFound, "")
	}
}
func postTopSecretSplit(c *gin.Context) {
	var satellite Satellite
	if err := c.ShouldBindJSON(&satellite); err != nil {
		fmt.Println(err)
		return
	}
	val := false
	for index, fleet := range fleets {
		if fleet.Name == c.Param("name") {
			fleets[index].Distance = satellite.Distance
			fleets[index].Message = satellite.Message
			val = true
		}
	}
	if val {
		c.IndentedJSON(http.StatusOK, "Satellite configurado correctamente")
	} else {
		c.IndentedJSON(http.StatusNotFound, "Satellite No encontrado")
	}
}
func getSatellite(c *gin.Context) {
	for _, fleet := range fleets {
		if fleet.Name == c.Param("name") {
			c.IndentedJSON(http.StatusOK, fleet)
		}
	}
	//c.IndentedJSON(http.StatusNotFound, "")
}
func getTopSecretSplit(c *gin.Context) {
	var satsConf []Satellite
	for _, fleet := range fleets {
		satsConf = append(satsConf, Satellite{
			Name:     fleet.Name,
			Distance: fleet.Distance,
			Message:  fleet.Message,
		})
	}
	reslocation, err1 := GetLocation(satsConf)
	resMessage, err2 := GetMessage(satsConf, "flex")
	if err1 == nil && err2 == nil {
		type Response struct {
			Position Point  `json:"position"`
			Message  string `json:"message"`
		}
		response := Response{
			Position: reslocation,
			Message:  resMessage,
		}

		c.IndentedJSON(http.StatusOK, response)
	} else {
		c.IndentedJSON(http.StatusNotFound, "")
	}
}
func setMode(c *gin.Context) {
	type RequestBody struct {
		Mode string `json:"mode"`
	}
	var newMode RequestBody
	if err := c.ShouldBindJSON(&newMode); err != nil {
		fmt.Println(err)
		return
	}
	if newMode.Mode == "strict" {
		mode = newMode.Mode
		c.IndentedJSON(http.StatusOK, "Configurado modo: 'strict'")
	} else if newMode.Mode == "flex" {
		mode = newMode.Mode
		c.IndentedJSON(http.StatusOK, "Configurado modo: 'flex'")
	} else {
		c.IndentedJSON(http.StatusNotFound, "El modo indicado no es aceptable")
	}
}
func getMode(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, mode)
}
func postSetPosition(c *gin.Context) {
	var newPosition Point
	if err := c.ShouldBindJSON(&newPosition); err != nil {
		fmt.Println(err)
		return
	}
	val := false
	for index, fleet := range fleets {
		if fleet.Name == c.Param("name") {
			fleets[index].Position = newPosition
			val = true
		}
	}
	if val {
		c.IndentedJSON(http.StatusOK, "Satellite configurado correctamente")
	} else {
		c.IndentedJSON(http.StatusNotFound, "Satellite No encontrado")
	}
}
