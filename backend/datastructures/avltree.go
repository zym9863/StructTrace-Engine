package datastructures

import "fmt"

// AVLNode represents a node in the AVL Tree
type AVLNode struct {
	ID     int
	Value  int
	Height int
	Left   *AVLNode
	Right  *AVLNode
}

// AVLTree represents an AVL Tree with step tracking
type AVLTree struct {
	Root   *AVLNode
	nextID int
	steps  []Step
}

// NewAVLTree creates a new AVL Tree
func NewAVLTree() *AVLTree {
	return &AVLTree{
		Root:   nil,
		nextID: 0,
		steps:  make([]Step, 0),
	}
}

func (t *AVLTree) clearSteps() {
	t.steps = make([]Step, 0)
}

func (t *AVLTree) addStep(stepType StepType, desc string, nodeID *int, extra ...interface{}) {
	step := Step{
		Type:        stepType,
		Description: desc,
		NodeID:      nodeID,
		TreeState:   t.getTreeSnapshot(),
	}
	if len(extra) > 0 {
		if highlights, ok := extra[0].([]int); ok {
			step.Highlight = highlights
		}
	}
	t.steps = append(t.steps, step)
}

func (t *AVLTree) getTreeSnapshot() []TreeNodeSnapshot {
	var nodes []TreeNodeSnapshot
	t.inorderSnapshot(t.Root, &nodes, 0, 0, 800)
	return nodes
}

func (t *AVLTree) inorderSnapshot(node *AVLNode, nodes *[]TreeNodeSnapshot, depth int, xMin, xMax float64) {
	if node == nil {
		return
	}

	x := (xMin + xMax) / 2
	y := float64(depth*80 + 50)

	snapshot := TreeNodeSnapshot{
		ID:     node.ID,
		Value:  node.Value,
		Height: node.Height,
		X:      x,
		Y:      y,
	}

	if node.Left != nil {
		leftID := node.Left.ID
		snapshot.LeftID = &leftID
	}
	if node.Right != nil {
		rightID := node.Right.ID
		snapshot.RightID = &rightID
	}

	*nodes = append(*nodes, snapshot)

	t.inorderSnapshot(node.Left, nodes, depth+1, xMin, x)
	t.inorderSnapshot(node.Right, nodes, depth+1, x, xMax)
}

func height(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (t *AVLTree) getBalance(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return height(node.Left) - height(node.Right)
}

func (t *AVLTree) rightRotate(y *AVLNode) *AVLNode {
	x := y.Left
	T2 := x.Right

	x.Right = y
	y.Left = T2

	y.Height = max(height(y.Left), height(y.Right)) + 1
	x.Height = max(height(x.Left), height(x.Right)) + 1

	t.addStep(StepRotateRight, fmt.Sprintf("对节点 %d 进行右旋", y.Value), &y.ID, []int{x.ID, y.ID})

	return x
}

func (t *AVLTree) leftRotate(x *AVLNode) *AVLNode {
	y := x.Right
	T2 := y.Left

	y.Left = x
	x.Right = T2

	x.Height = max(height(x.Left), height(x.Right)) + 1
	y.Height = max(height(y.Left), height(y.Right)) + 1

	t.addStep(StepRotateLeft, fmt.Sprintf("对节点 %d 进行左旋", x.Value), &x.ID, []int{x.ID, y.ID})

	return y
}

func (t *AVLTree) insert(node *AVLNode, value int) *AVLNode {
	if node == nil {
		newNode := &AVLNode{
			ID:     t.nextID,
			Value:  value,
			Height: 1,
		}
		t.nextID++
		t.addStep(StepInsert, fmt.Sprintf("插入节点 %d", value), &newNode.ID, []int{newNode.ID})
		return newNode
	}

	t.addStep(StepCompare, fmt.Sprintf("比较 %d 与节点 %d", value, node.Value), &node.ID, []int{node.ID})

	if value < node.Value {
		node.Left = t.insert(node.Left, value)
	} else if value > node.Value {
		node.Right = t.insert(node.Right, value)
	} else {
		return node // Duplicate values not allowed
	}

	node.Height = 1 + max(height(node.Left), height(node.Right))

	balance := t.getBalance(node)

	// Left Left Case
	if balance > 1 && value < node.Left.Value {
		t.addStep(StepRebalance, "LL情况：需要右旋", &node.ID)
		return t.rightRotate(node)
	}

	// Right Right Case
	if balance < -1 && value > node.Right.Value {
		t.addStep(StepRebalance, "RR情况：需要左旋", &node.ID)
		return t.leftRotate(node)
	}

	// Left Right Case
	if balance > 1 && value > node.Left.Value {
		t.addStep(StepRebalance, "LR情况：先左旋后右旋", &node.ID)
		node.Left = t.leftRotate(node.Left)
		return t.rightRotate(node)
	}

	// Right Left Case
	if balance < -1 && value < node.Right.Value {
		t.addStep(StepRebalance, "RL情况：先右旋后左旋", &node.ID)
		node.Right = t.rightRotate(node.Right)
		return t.leftRotate(node)
	}

	return node
}

// Insert inserts a value into the AVL Tree
func (t *AVLTree) Insert(value int) OperationResult {
	t.clearSteps()
	t.addStep(StepInsert, fmt.Sprintf("开始插入值 %d", value), nil)
	t.Root = t.insert(t.Root, value)
	t.addStep(StepComplete, "插入完成", nil)

	return OperationResult{
		Success:   true,
		Steps:     t.steps,
		FinalTree: t.getTreeSnapshot(),
	}
}

// Search searches for a value in the AVL Tree
func (t *AVLTree) Search(value int) OperationResult {
	t.clearSteps()

	current := t.Root
	for current != nil {
		t.addStep(StepCompare, fmt.Sprintf("比较 %d 与节点 %d", value, current.Value), &current.ID, []int{current.ID})
		if value == current.Value {
			t.addStep(StepFound, fmt.Sprintf("找到节点 %d", value), &current.ID, []int{current.ID})
			return OperationResult{
				Success:   true,
				Message:   fmt.Sprintf("找到值 %d", value),
				Steps:     t.steps,
				FinalTree: t.getTreeSnapshot(),
			}
		} else if value < current.Value {
			current = current.Left
		} else {
			current = current.Right
		}
	}

	t.addStep(StepNotFound, fmt.Sprintf("值 %d 不存在于树中", value), nil)
	return OperationResult{
		Success:   false,
		Message:   fmt.Sprintf("值 %d 不存在", value),
		Steps:     t.steps,
		FinalTree: t.getTreeSnapshot(),
	}
}

// minValueNode finds the node with minimum value in a subtree
func (t *AVLTree) minValueNode(node *AVLNode) *AVLNode {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

// delete deletes a node with given value from the subtree
func (t *AVLTree) delete(node *AVLNode, value int) *AVLNode {
	if node == nil {
		t.addStep(StepNotFound, fmt.Sprintf("值 %d 不存在于树中", value), nil)
		return node
	}

	t.addStep(StepCompare, fmt.Sprintf("比较 %d 与节点 %d", value, node.Value), &node.ID, []int{node.ID})

	if value < node.Value {
		node.Left = t.delete(node.Left, value)
	} else if value > node.Value {
		node.Right = t.delete(node.Right, value)
	} else {
		// Node to be deleted found
		t.addStep(StepDelete, fmt.Sprintf("找到要删除的节点 %d", value), &node.ID, []int{node.ID})

		// Node with only one child or no child
		if node.Left == nil {
			t.addStep(StepDelete, fmt.Sprintf("节点 %d 没有左子节点，用右子节点替换", node.Value), &node.ID)
			return node.Right
		} else if node.Right == nil {
			t.addStep(StepDelete, fmt.Sprintf("节点 %d 没有右子节点，用左子节点替换", node.Value), &node.ID)
			return node.Left
		}

		// Node with two children: Get the inorder successor (smallest in right subtree)
		successor := t.minValueNode(node.Right)
		t.addStep(StepDelete, fmt.Sprintf("节点 %d 有两个子节点，找到后继节点 %d", node.Value, successor.Value), &successor.ID, []int{node.ID, successor.ID})

		// Copy the inorder successor's value to this node
		node.Value = successor.Value
		node.ID = successor.ID
		t.addStep(StepDelete, fmt.Sprintf("用后继节点 %d 替换被删除节点", successor.Value), &node.ID)

		// Delete the inorder successor
		node.Right = t.delete(node.Right, successor.Value)
	}

	// Update height
	node.Height = 1 + max(height(node.Left), height(node.Right))

	// Get balance factor
	balance := t.getBalance(node)

	// Left Left Case
	if balance > 1 && t.getBalance(node.Left) >= 0 {
		t.addStep(StepRebalance, "LL情况：需要右旋", &node.ID, []int{node.ID})
		return t.rightRotate(node)
	}

	// Left Right Case
	if balance > 1 && t.getBalance(node.Left) < 0 {
		t.addStep(StepRebalance, "LR情况：先左旋后右旋", &node.ID, []int{node.ID})
		node.Left = t.leftRotate(node.Left)
		return t.rightRotate(node)
	}

	// Right Right Case
	if balance < -1 && t.getBalance(node.Right) <= 0 {
		t.addStep(StepRebalance, "RR情况：需要左旋", &node.ID, []int{node.ID})
		return t.leftRotate(node)
	}

	// Right Left Case
	if balance < -1 && t.getBalance(node.Right) > 0 {
		t.addStep(StepRebalance, "RL情况：先右旋后左旋", &node.ID, []int{node.ID})
		node.Right = t.rightRotate(node.Right)
		return t.leftRotate(node)
	}

	return node
}

// Delete deletes a value from the AVL Tree
func (t *AVLTree) Delete(value int) OperationResult {
	t.clearSteps()
	t.addStep(StepDelete, fmt.Sprintf("开始删除值 %d", value), nil)

	// Check if value exists
	found := false
	current := t.Root
	for current != nil {
		if value == current.Value {
			found = true
			break
		} else if value < current.Value {
			current = current.Left
		} else {
			current = current.Right
		}
	}

	if !found {
		t.addStep(StepNotFound, fmt.Sprintf("值 %d 不存在于树中，无法删除", value), nil)
		return OperationResult{
			Success:   false,
			Message:   fmt.Sprintf("值 %d 不存在，无法删除", value),
			Steps:     t.steps,
			FinalTree: t.getTreeSnapshot(),
		}
	}

	t.Root = t.delete(t.Root, value)
	t.addStep(StepComplete, fmt.Sprintf("删除节点 %d 完成", value), nil)

	return OperationResult{
		Success:   true,
		Message:   fmt.Sprintf("成功删除值 %d", value),
		Steps:     t.steps,
		FinalTree: t.getTreeSnapshot(),
	}
}
