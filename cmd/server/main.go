package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/maksimkasimovhse/Web-based-Cartographic-Search-Service/internal/db"
	"github.com/maksimkasimovhse/Web-based-Cartographic-Search-Service/internal/graph"
	"github.com/maksimkasimovhse/Web-based-Cartographic-Search-Service/internal/handlers"
)

func corsMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		fmt.Println("Ошибка подключения к pgx: ", err)
		os.Exit(1)
	}

	router := gin.Default()
	router.Use(corsMiddleware())
	router.GET("/places/nearby", handlers.NearbyPlaces(conn))

	gr := &graph.Graph{}
	gr.LoadGraph(ctx, conn)
	fmt.Println("граф загружен, узлов:", len(gr.Coords))
	router.GET("/route", handlers.RouteHandler(gr, conn))
	router.Run(":8080")
}
