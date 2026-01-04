package datastructures

import (
	"container/heap"
	"fmt"
	"math"
)

// Edge represents an edge in the graph
type Edge struct {
	To     string
	Weight int
}

// Graph represents a weighted graph with step tracking
type Graph struct {
	Nodes      map[string][]Edge
	NodeCoords map[string][2]float64
	steps      []Step
}

// NewGraph creates a new Graph
func NewGraph() *Graph {
	return &Graph{
		Nodes:      make(map[string][]Edge),
		NodeCoords: make(map[string][2]float64),
		steps:      make([]Step, 0),
	}
}

func (g *Graph) clearSteps() {
	g.steps = make([]Step, 0)
}

func (g *Graph) addStep(stepType StepType, desc string, distances map[string]int, visited map[string]bool, path []string, currentEdge *[2]string) {
	nodes := make([]GraphNodeSnapshot, 0)
	for id := range g.Nodes {
		dist := distances[id]
		var distPtr *int
		if dist != math.MaxInt32 {
			distPtr = &dist
		}
		coords := g.NodeCoords[id]
		inPath := false
		for _, p := range path {
			if p == id {
				inPath = true
				break
			}
		}
		nodes = append(nodes, GraphNodeSnapshot{
			ID:       id,
			Label:    id,
			X:        coords[0],
			Y:        coords[1],
			Distance: distPtr,
			Visited:  visited[id],
			InPath:   inPath,
		})
	}

	edges := make([]GraphEdgeSnapshot, 0)
	for from, neighbors := range g.Nodes {
		for _, e := range neighbors {
			inPath := false
			for i := 0; i < len(path)-1; i++ {
				if (path[i] == from && path[i+1] == e.To) || (path[i] == e.To && path[i+1] == from) {
					inPath = true
					break
				}
			}
			selected := false
			if currentEdge != nil && ((currentEdge[0] == from && currentEdge[1] == e.To) || (currentEdge[0] == e.To && currentEdge[1] == from)) {
				selected = true
			}
			edges = append(edges, GraphEdgeSnapshot{
				From:     from,
				To:       e.To,
				Weight:   e.Weight,
				InPath:   inPath,
				Selected: selected,
			})
		}
	}

	step := Step{
		Type:        stepType,
		Description: desc,
		GraphNodes:  nodes,
		GraphEdges:  edges,
	}
	g.steps = append(g.steps, step)
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(id string, x, y float64) {
	if _, exists := g.Nodes[id]; !exists {
		g.Nodes[id] = make([]Edge, 0)
	}
	g.NodeCoords[id] = [2]float64{x, y}
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(from, to string, weight int) {
	g.Nodes[from] = append(g.Nodes[from], Edge{To: to, Weight: weight})
	g.Nodes[to] = append(g.Nodes[to], Edge{To: from, Weight: weight}) // Undirected
}

// PriorityQueueItem for Dijkstra
type PriorityQueueItem struct {
	node     string
	priority int
	index    int
}

type PriorityQueue []*PriorityQueueItem

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*PriorityQueueItem)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Dijkstra finds the shortest path from start to end
func (g *Graph) Dijkstra(start, end string) OperationResult {
	g.clearSteps()

	distances := make(map[string]int)
	previous := make(map[string]string)
	visited := make(map[string]bool)

	for node := range g.Nodes {
		distances[node] = math.MaxInt32
	}
	distances[start] = 0

	g.addStep(StepVisit, fmt.Sprintf("初始化：起点 %s 距离设为 0", start), distances, visited, nil, nil)

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &PriorityQueueItem{node: start, priority: 0})

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*PriorityQueueItem)

		if visited[current.node] {
			continue
		}
		visited[current.node] = true

		g.addStep(StepSelectNode, fmt.Sprintf("选择距离最小的未访问节点: %s (距离: %d)", current.node, distances[current.node]), distances, visited, nil, nil)

		if current.node == end {
			// Reconstruct path
			path := make([]string, 0)
			for at := end; at != ""; at = previous[at] {
				path = append([]string{at}, path...)
			}
			g.addStep(StepComplete, fmt.Sprintf("找到最短路径: %v, 总距离: %d", path, distances[end]), distances, visited, path, nil)

			return OperationResult{
				Success: true,
				Message: fmt.Sprintf("最短路径距离: %d", distances[end]),
				Steps:   g.steps,
				FinalGraph: &struct {
					Nodes []GraphNodeSnapshot `json:"nodes"`
					Edges []GraphEdgeSnapshot `json:"edges"`
				}{
					Nodes: g.steps[len(g.steps)-1].GraphNodes,
					Edges: g.steps[len(g.steps)-1].GraphEdges,
				},
			}
		}

		for _, edge := range g.Nodes[current.node] {
			if visited[edge.To] {
				continue
			}

			newDist := distances[current.node] + edge.Weight
			edgePtr := &[2]string{current.node, edge.To}

			if newDist < distances[edge.To] {
				distances[edge.To] = newDist
				previous[edge.To] = current.node
				heap.Push(&pq, &PriorityQueueItem{node: edge.To, priority: newDist})
				g.addStep(StepUpdateDist, fmt.Sprintf("更新节点 %s 距离: %d → %d (通过 %s)", edge.To, distances[edge.To], newDist, current.node), distances, visited, nil, edgePtr)
			} else {
				g.addStep(StepCompare, fmt.Sprintf("边 %s→%s: 新距离 %d >= 当前距离 %d，不更新", current.node, edge.To, newDist, distances[edge.To]), distances, visited, nil, edgePtr)
			}
		}
	}

	g.addStep(StepNotFound, fmt.Sprintf("无法从 %s 到达 %s", start, end), distances, visited, nil, nil)
	return OperationResult{
		Success: false,
		Message: "无法到达目标节点",
		Steps:   g.steps,
	}
}

// CreateSampleGraph creates a sample graph for demonstration
func CreateSampleGraph() *Graph {
	g := NewGraph()

	// Add nodes with positions
	g.AddNode("A", 100, 150)
	g.AddNode("B", 250, 50)
	g.AddNode("C", 250, 250)
	g.AddNode("D", 400, 100)
	g.AddNode("E", 400, 200)
	g.AddNode("F", 550, 150)

	// Add edges with weights
	g.AddEdge("A", "B", 4)
	g.AddEdge("A", "C", 2)
	g.AddEdge("B", "C", 1)
	g.AddEdge("B", "D", 5)
	g.AddEdge("C", "D", 8)
	g.AddEdge("C", "E", 10)
	g.AddEdge("D", "E", 2)
	g.AddEdge("D", "F", 6)
	g.AddEdge("E", "F", 3)

	return g
}
