package graph

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SnapToGraph(ctx context.Context, pool *pgxpool.Pool, lat, lon float64) (*int64, error) {
	var id int64
	err := pool.QueryRow(ctx, `SELECT id FROM nodes ORDER BY ST_Transform(ST_SetSRID(geom, 3857), 4326) <-> ST_SetSRID(ST_MakePoint($1, $2), 4326) LIMIT 1`, lon, lat).Scan(&id)
	if err != nil {
		fmt.Println("Ошибка KNN/SELECT из nodes: ", err)
		return nil, err
	}
	return &id, nil
}
