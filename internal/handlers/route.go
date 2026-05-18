package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maksimkasimovhse/Web-based-Cartographic-Search-Service/internal/graph"
)

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func RouteHandler(grWalk *graph.Graph, grCar *graph.Graph, pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		mode := c.DefaultQuery("mode", "walk")
		gr := grWalk
		if mode == "car" {
			gr = grCar
		}

		from_lat := c.Query("latFrom")
		from_lon := c.Query("lonFrom")
		to_lat := c.Query("latTo")
		to_lon := c.Query("lonTo")

		fromLat, e1 := strconv.ParseFloat(from_lat, 64)
		fromLon, e2 := strconv.ParseFloat(from_lon, 64)
		toLat, e3 := strconv.ParseFloat(to_lat, 64)
		toLon, e4 := strconv.ParseFloat(to_lon, 64)

		if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
			fmt.Println("Ошибка парасинга")
			c.JSON(400, gin.H{"error": "invalid coordinates"})
			return
		}

		idFROM, err := graph.SnapToGraph(c.Request.Context(), pool, fromLat, fromLon)
		fmt.Println("idFROM:", *idFROM)
		if err != nil {
			fmt.Println("Ошибка SnapToGraph: ", err)
			c.JSON(400, gin.H{"error": "invalid coordinates"})
			return
		}

		idTO, err := graph.SnapToGraph(c.Request.Context(), pool, toLat, toLon)
		fmt.Println("idTO:", *idTO)
		if err != nil {
			fmt.Println("Ошибка SnapToGraph: ", err)
			c.JSON(400, gin.H{"error": "invalid coordinates"})
			return
		}

		way, distance, ok := gr.Dijkstra(*idFROM, *idTO)
		if !ok {
			fmt.Println("Ошибка Дейкстры: ", err)
			c.JSON(404, gin.H{"error": "route not found"})
			return
		}

		var g [][2]float64
		for _, id := range way {
			g = append(g, gr.Coords[id])
		}

		c.JSON(200, gin.H{"type": "Feature", "distance": distance, "geometry": gin.H{"type": "LineString", "coordinates": g}})

	}
}
