package handlers

import (
	"net/http"

	"gin/datastructures"

	"github.com/gin-gonic/gin"
)

// OperationRequest represents a request to perform an operation on a data structure
type OperationRequest struct {
	Structure string                 `json:"structure" binding:"required"`
	Operation string                 `json:"operation" binding:"required"`
	Params    map[string]interface{} `json:"params"`
}

// Stateful data structures for persistence
var (
	rbTree  = datastructures.NewRedBlackTree()
	avlTree = datastructures.NewAVLTree()
	graph   = datastructures.CreateSampleGraph()
)

// HandleOperation handles data structure operation requests
func HandleOperation(c *gin.Context) {
	var req OperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request: " + err.Error(),
		})
		return
	}

	var result datastructures.OperationResult

	switch req.Structure {
	case "rbtree":
		result = handleRBTreeOperation(req)
	case "avltree":
		result = handleAVLTreeOperation(req)
	case "graph":
		result = handleGraphOperation(req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Unknown structure: " + req.Structure,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func handleRBTreeOperation(req OperationRequest) datastructures.OperationResult {
	switch req.Operation {
	case "insert":
		value := getIntParam(req.Params, "value", 0)
		return rbTree.Insert(value)
	case "search":
		value := getIntParam(req.Params, "value", 0)
		return rbTree.Search(value)
	case "delete":
		value := getIntParam(req.Params, "value", 0)
		return rbTree.Delete(value)
	case "reset":
		rbTree = datastructures.NewRedBlackTree()
		return datastructures.OperationResult{
			Success: true,
			Message: "Red-Black Tree 已重置",
			Steps:   []datastructures.Step{},
		}
	default:
		return datastructures.OperationResult{
			Success: false,
			Message: "Unknown operation: " + req.Operation,
		}
	}
}

func handleAVLTreeOperation(req OperationRequest) datastructures.OperationResult {
	switch req.Operation {
	case "insert":
		value := getIntParam(req.Params, "value", 0)
		return avlTree.Insert(value)
	case "search":
		value := getIntParam(req.Params, "value", 0)
		return avlTree.Search(value)
	case "delete":
		value := getIntParam(req.Params, "value", 0)
		return avlTree.Delete(value)
	case "reset":
		avlTree = datastructures.NewAVLTree()
		return datastructures.OperationResult{
			Success: true,
			Message: "AVL Tree 已重置",
			Steps:   []datastructures.Step{},
		}
	default:
		return datastructures.OperationResult{
			Success: false,
			Message: "Unknown operation: " + req.Operation,
		}
	}
}

func handleGraphOperation(req OperationRequest) datastructures.OperationResult {
	switch req.Operation {
	case "shortest_path":
		start := getStringParam(req.Params, "start", "A")
		end := getStringParam(req.Params, "end", "F")
		return graph.Dijkstra(start, end)
	case "reset":
		graph = datastructures.CreateSampleGraph()
		return datastructures.OperationResult{
			Success: true,
			Message: "Graph 已重置",
			Steps:   []datastructures.Step{},
		}
	default:
		return datastructures.OperationResult{
			Success: false,
			Message: "Unknown operation: " + req.Operation,
		}
	}
}

func getIntParam(params map[string]interface{}, key string, defaultVal int) int {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		}
	}
	return defaultVal
}

func getStringParam(params map[string]interface{}, key string, defaultVal string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultVal
}

// HandleReset resets all data structures
func HandleReset(c *gin.Context) {
	rbTree = datastructures.NewRedBlackTree()
	avlTree = datastructures.NewAVLTree()
	graph = datastructures.CreateSampleGraph()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "All data structures have been reset",
	})
}
