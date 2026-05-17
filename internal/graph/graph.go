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
	graph map[int64][]Edge
}

func (gr *Graph) LoadGraph(ctx context.Context, pool *pgxpool.Pool) (*Graph, error) {
	gr.graph = make(map[int64][]Edge)
	query := `SELECT from_node, to_node, weight_road, oneway FROM roads`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		fmt.Println("Ошибка SELECT из roads")
		return nil, err
	}
	for rows.Next() {
		var from_node, to_node int64
		var weight_road float64
		var oneway string
		rows.Scan(&from_node, &to_node, &weight_road, &oneway)
		if oneway == "yes" {
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

	return gr, nil
}
