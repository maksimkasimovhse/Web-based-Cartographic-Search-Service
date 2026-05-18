package graph

import (
	"container/heap"
	"fmt"
	"math"
)

type Node struct {
	node_id int64
	dist    float64
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }

func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Node)) }

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func (gr *Graph) Dijkstra(node_from int64, node_to int64) ([]int64, bool) {
	dist := make(map[int64]float64)
	prev := make(map[int64]int64)

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	dist[node_from] = 0.0
	dist[node_to] = math.Inf(1)
	prev[node_from] = -1

	heap.Push(&pq, &Node{node_id: node_from, dist: 0.0})

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Node)
		if item.dist > dist[item.node_id] {
			continue
		}
		for _, edge := range gr.graph[item.node_id] {
			if d, ok := dist[edge.to_node]; !ok || item.dist+edge.weight < d {
				dist[edge.to_node] = item.dist + edge.weight
				prev[edge.to_node] = item.node_id
				heap.Push(&pq, &Node{node_id: edge.to_node, dist: dist[edge.to_node]})
			}
		}
	}

	fmt.Println("dist[node_to]:", dist[node_to])
	if dist[node_to] == math.Inf(1) {
		return nil, false
	}

	var path []int64
	for a := node_to; a != -1; a = prev[a] {
		path = append(path, a)
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, true
}
