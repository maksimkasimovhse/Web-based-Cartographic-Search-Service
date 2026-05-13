package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Place struct {
	name     string
	category string
	osmID    int64
	x, y     float64
}

var amenityMap = map[string]string{
	"cafe":       "Кафе",
	"bank":       "Банк",
	"atm":        "Банкомат",
	"fast_food":  "Фастфуд",
	"restaurant": "Ресторан",
	"fuel":       "Заправка",
	"bar":        "Бар",
	"parking":    "Паркинг",
	"clinic":     "Клиника",
}

var shopMap = map[string]string{
	"supermarket": "Супермаркет",
	"alcohol":     "Алкомаркет",
	"clothes":     "Магазин одежды",
	"car_repair":  "Ремонт автомобиля",
	"hairdresser": "Парикмахерская",
	"electronics": "Магазин электроники",
}

func importPlaces(ctx context.Context, conn *pgx.Conn, query string) error {
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка SELECT: %w", err)
	}
	defer rows.Close()

	var places []Place
	for rows.Next() {
		var name, amenity, shop *string
		var osmID int64
		var x, y float64
		if err := rows.Scan(&name, &amenity, &shop, &osmID, &x, &y); err != nil {
			return fmt.Errorf("ошибка Scan: %w", err)
		}

		var nameVal, amenityVal, shopVal string
		if name != nil {
			nameVal = *name
		}
		if amenity != nil {
			amenityVal = *amenity
		}
		if shop != nil {
			shopVal = *shop
		}

		if cat, ok := amenityMap[amenityVal]; ok {
			places = append(places, Place{nameVal, cat, osmID, x, y})
		} else if cat, ok := shopMap[shopVal]; ok {
			places = append(places, Place{nameVal, cat, osmID, x, y})
		}
	}
	if rows.Err() != nil {
		return fmt.Errorf("ошибка чтения строк: %w", rows.Err())
	}

	fmt.Println("Собрано объектов:", len(places))

	for _, p := range places {
		_, err := conn.Exec(ctx,
			"INSERT INTO places(name, category, osm_id, geom) VALUES($1, $2, $3, ST_SetSRID(ST_Point($4, $5), 4326))",
			p.name, p.category, p.osmID, p.x, p.y,
		)
		if err != nil {
			return fmt.Errorf("ошибка INSERT: %w", err)
		}
	}
	return nil
}

func main() {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Println("DATABASE_URL не задана")
		os.Exit(1)
	}
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	pointQuery := `SELECT name, amenity, shop, osm_id,
		ST_X(ST_Transform(way, 4326)), ST_Y(ST_Transform(way, 4326))
		FROM planet_osm_point WHERE amenity IS NOT NULL OR shop IS NOT NULL`

	polygonQuery := `SELECT name, amenity, shop, osm_id,
		ST_X(ST_Transform(ST_Centroid(way), 4326)), ST_Y(ST_Transform(ST_Centroid(way), 4326))
		FROM planet_osm_polygon WHERE amenity IS NOT NULL OR shop IS NOT NULL`

	if err := importPlaces(ctx, conn, pointQuery); err != nil {
		fmt.Println("Ошибка импорта точек:", err)
		os.Exit(1)
	}

	if err := importPlaces(ctx, conn, polygonQuery); err != nil {
		fmt.Println("Ошибка импорта полигонов:", err)
		os.Exit(1)
	}

	fmt.Println("Импорт завершён")
}