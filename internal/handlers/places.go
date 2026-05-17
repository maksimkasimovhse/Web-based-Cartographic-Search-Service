package handlers

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Feature struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Properties struct {
	OsmID    int64  `json:"osm_id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func NearbyPlaces(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		lat := c.Query("lat")
		lon := c.Query("lon")
		radius := c.Query("radius")

		latF, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid lat"})
			return
		}
		lonF, err := strconv.ParseFloat(lon, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid lon"})
			return
		}
		radiusI, err := strconv.ParseFloat(radius, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid radius"})
			return
		}

		category := c.Query("category")

		ctx := c.Request.Context()

		rows, err := pool.Query(ctx, `SELECT osm_id, name, category, ST_X(geom), ST_Y(geom) FROM places WHERE (category = $1 OR $1 = '') AND ST_DWithin(geom::geography, ST_SetSRID(ST_Point($2, $3), 4326)::geography, $4)`, category, lonF, latF, radiusI)
		if err != nil {
			log.Println("Ошибка SELECT-запроса", err)
			c.JSON(500, gin.H{"error": "internal server wrong"})
			return
		}
		defer rows.Close()

		collection := FeatureCollection{
			Type:     "FeatureCollection",
			Features: []Feature{},
		}

		for rows.Next() {
			var osm_id int64
			var name string
			var category string
			var x, y float64
			if err := rows.Scan(&osm_id, &name, &category, &x, &y); err != nil {
				log.Println("Ошибка чтения строки", err)
				continue
			}

			feature := Feature{
				Type: "Feature",
				Properties: Properties{
					OsmID:    osm_id,
					Name:     name,
					Category: category,
				},
				Geometry: Geometry{
					Type:        "Point",
					Coordinates: []float64{x, y},
				},
			}

			collection.Features = append(collection.Features, feature)

		}

		if err := rows.Err(); err != nil {
			log.Println("rows error:", err)
		}

		c.JSON(200, collection)
	}
}
