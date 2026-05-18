package graph

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Edge struct {
	to_node int64
	weight  float64
}

type Graph struct {
	graph  map[int64][]Edge
	Coords map[int64][2]float64
}

func (gr *Graph) LoadGraph(ctx context.Context, pool *pgxpool.Pool, mode string) (*Graph, error) {
	gr.graph = make(map[int64][]Edge)
	gr.Coords = make(map[int64][2]float64)
	var query string
	if mode == "car" {
		query = `SELECT from_node, to_node, weight_road, oneway FROM roads WHERE highway IN ('motorway','primary','primary_link','secondary','secondary_link','tertiary','tertiary_link','residential','living_street','service')`
	} else {
		query = `SELECT from_node, to_node, weight_road, oneway FROM roads`
	}
	rows, err := pool.Query(ctx, query)
	if err != nil {
		fmt.Println("Ошибка SELECT из roads: ", err)
		return nil, err
	}
	for rows.Next() {
		var from_node, to_node int64
		var weight_road float64
		var oneway *bool
		rows.Scan(&from_node, &to_node, &weight_road, &oneway)
		if oneway != nil && *oneway {
			gr.graph[from_node] = append(gr.graph[from_node], Edge{to_node: to_node, weight: weight_road})
		} else {
			gr.graph[from_node] = append(gr.graph[from_node], Edge{to_node: to_node, weight: weight_road})
			gr.graph[to_node] = append(gr.graph[to_node], Edge{to_node: from_node, weight: weight_road})
		}
	}
	if rows.Err() != nil {
		fmt.Println("Ошибка построения графа: ", rows.Err())
		return nil, rows.Err()
	}

	query2 := `SELECT id, ST_X(geom), ST_Y(geom) FROM nodes`
	rows2, err := pool.Query(ctx, query2)
	if err != nil {
		fmt.Println("Ошибка SELECT из nodes: ", err)
		return nil, err
	}
	for rows2.Next() {
		var id int64
		var x, y float64
		rows2.Scan(&id, &x, &y)
		gr.Coords[id] = [2]float64{x, y}
	}
	if rows2.Err() != nil {
		fmt.Println("Ошибка занесения координат: ", rows2.Err())
		return nil, rows2.Err()
	}

	fmt.Println("рёбер в графе от узла 0:", len(gr.graph[0]))
	fmt.Println("рёбер в графе от узла 142242:", len(gr.graph[142242]))
	return gr, nil
}
