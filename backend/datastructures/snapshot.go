package datastructures

// NodeColor represents the color of a node in Red-Black Tree
type NodeColor string

const (
	Red   NodeColor = "red"
	Black NodeColor = "black"
)

// StepType represents the type of operation step
type StepType string

const (
	StepInsert       StepType = "insert"
	StepDelete       StepType = "delete"
	StepRotateLeft   StepType = "rotate_left"
	StepRotateRight  StepType = "rotate_right"
	StepColorChange  StepType = "color_change"
	StepCompare      StepType = "compare"
	StepVisit        StepType = "visit"
	StepFound        StepType = "found"
	StepNotFound     StepType = "not_found"
	StepUpdateDist   StepType = "update_distance"
	StepSelectNode   StepType = "select_node"
	StepMarkVisited  StepType = "mark_visited"
	StepRebalance    StepType = "rebalance"
	StepComplete     StepType = "complete"
)

// TreeNodeSnapshot represents a snapshot of a tree node
type TreeNodeSnapshot struct {
	ID       int       `json:"id"`
	Value    int       `json:"value"`
	Color    NodeColor `json:"color,omitempty"`
	LeftID   *int      `json:"leftId,omitempty"`
	RightID  *int      `json:"rightId,omitempty"`
	ParentID *int      `json:"parentId,omitempty"`
	Height   int       `json:"height,omitempty"`
	X        float64   `json:"x,omitempty"`
	Y        float64   `json:"y,omitempty"`
}

// GraphNodeSnapshot represents a snapshot of a graph node
type GraphNodeSnapshot struct {
	ID       string  `json:"id"`
	Label    string  `json:"label"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Distance *int    `json:"distance,omitempty"`
	Visited  bool    `json:"visited"`
	InPath   bool    `json:"inPath"`
}

// GraphEdgeSnapshot represents a snapshot of a graph edge
type GraphEdgeSnapshot struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Weight   int    `json:"weight"`
	InPath   bool   `json:"inPath"`
	Selected bool   `json:"selected"`
}

// Step represents a single step in the algorithm execution
type Step struct {
	Type        StepType            `json:"type"`
	Description string              `json:"description"`
	NodeID      *int                `json:"nodeId,omitempty"`
	TargetID    *int                `json:"targetId,omitempty"`
	Value       *int                `json:"value,omitempty"`
	OldColor    NodeColor           `json:"oldColor,omitempty"`
	NewColor    NodeColor           `json:"newColor,omitempty"`
	TreeState   []TreeNodeSnapshot  `json:"treeState,omitempty"`
	GraphNodes  []GraphNodeSnapshot `json:"graphNodes,omitempty"`
	GraphEdges  []GraphEdgeSnapshot `json:"graphEdges,omitempty"`
	Highlight   []int               `json:"highlight,omitempty"`
}

// OperationResult represents the result of a data structure operation
type OperationResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Steps     []Step `json:"steps"`
	FinalTree []TreeNodeSnapshot  `json:"finalTree,omitempty"`
	FinalGraph *struct {
		Nodes []GraphNodeSnapshot `json:"nodes"`
		Edges []GraphEdgeSnapshot `json:"edges"`
	} `json:"finalGraph,omitempty"`
}
