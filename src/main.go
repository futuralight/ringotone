package main

import (
	"container/ring"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"ringotone/src/errorhandling"
	"ringotone/src/logging"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//Rings - map of rings
var Rings map[string]*ring.Ring

func main() {
	Rings = make(map[string]*ring.Ring)
	err := loadEnv()
	errorhandling.HandleError(err)
	logging.LoadLogFile(os.Getenv("LOGGING_FILE"))
	loadServer()
	fmt.Println("ring ring ring")
}

func loadServer() {
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	routes(r)
	r.Run(":" + os.Getenv("SERVER_PORT"))
}

func routes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/next/:ring", nextRingItem)
	r.GET("/prev/:ring", prevRingItem)
	r.GET("/add/:ring/:value", addRingItem)
	r.POST("/add/:ring", postRingItem)
}

func nextRingItem(c *gin.Context) {
	r, exst := Rings[strings.ToLower(c.Param("ring"))]
	if !exst {
		c.JSON(200, gin.H{
			"ring":    c.Param("ring"),
			"success": false,
		})
		return
	}
	r = r.Next()
	c.JSON(200, gin.H{
		"ring":  c.Param("ring"),
		"value": r.Value,
	})
	Rings[strings.ToLower(c.Param("ring"))] = r
}

func prevRingItem(c *gin.Context) {
	r, exst := Rings[strings.ToLower(c.Param("ring"))]
	if !exst {
		c.JSON(200, gin.H{
			"ring":    c.Param("ring"),
			"success": false,
		})
		return
	}
	r = r.Prev()
	c.JSON(200, gin.H{
		"ring":  c.Param("ring"),
		"value": r.Value,
	})
	Rings[strings.ToLower(c.Param("ring"))] = r
}

func addRingItem(c *gin.Context) {
	r, exst := Rings[c.Param("ring")]
	if !exst {
		r = ring.New(1)
		r.Value = c.Param("value")
		Rings[strings.ToLower(c.Param("ring"))] = r
	} else {
		n := r.Len()
		newR := ring.New(n + 1)
		for i := 0; i < n; i++ {
			newR.Value = r.Value
			r = r.Next()
			newR = newR.Next()
		}
		newR.Value = c.Param("value")
		Rings[strings.ToLower(c.Param("ring"))] = newR
		r = newR
	}
	c.JSON(200, gin.H{
		"ring":    c.Param("ring"),
		"success": true,
		"value":   c.Param("value"),
		"length":  r.Len(),
	})
}

func postRingItem(c *gin.Context) {
	r, exst := Rings[strings.ToLower(c.Param("ring"))]
	if !exst {
		r = ring.New(1)
		x, _ := ioutil.ReadAll(c.Request.Body)
		r.Value = string(x)
		Rings[strings.ToLower(c.Param("ring"))] = r
	} else {
		n := r.Len()
		newR := ring.New(n + 1)
		for i := 0; i < n; i++ {
			newR.Value = r.Value
			r = r.Next()
			newR = newR.Next()
		}
		x, _ := ioutil.ReadAll(c.Request.Body)
		newR.Value = string(x)
		Rings[strings.ToLower(c.Param("ring"))] = newR
		r = newR
	}
	c.JSON(200, gin.H{
		"ring":    c.Param("ring"),
		"success": true,
		"value":   r.Value,
		"length":  r.Len(),
	})
}

func loadEnv() error {
	godotenv.Load()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	err = godotenv.Load(dir + "/.env") //Загрузка .env файла
	return nil
}
