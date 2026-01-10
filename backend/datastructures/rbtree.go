package datastructures

import "fmt"

// RBNode represents a node in the Red-Black Tree
type RBNode struct {
	ID     int
	Value  int
	Color  NodeColor
	Left   *RBNode
	Right  *RBNode
	Parent *RBNode
}

// RedBlackTree represents a Red-Black Tree with step tracking
type RedBlackTree struct {
	Root   *RBNode
	NIL    *RBNode
	nextID int
	steps  []Step
}

// NewRedBlackTree creates a new Red-Black Tree
func NewRedBlackTree() *RedBlackTree {
	nil := &RBNode{Color: Black, ID: -1}
	return &RedBlackTree{
		Root:   nil,
		NIL:    nil,
		nextID: 0,
		steps:  make([]Step, 0),
	}
}

// clearSteps resets the step tracking
func (t *RedBlackTree) clearSteps() {
	t.steps = make([]Step, 0)
}

// addStep records a step in the algorithm
func (t *RedBlackTree) addStep(stepType StepType, desc string, nodeID *int, extra ...interface{}) {
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

// getTreeSnapshot creates a snapshot of the current tree state
func (t *RedBlackTree) getTreeSnapshot() []TreeNodeSnapshot {
	var nodes []TreeNodeSnapshot
	t.inorderSnapshot(t.Root, &nodes, 0, 0, 800)
	return nodes
}

func (t *RedBlackTree) inorderSnapshot(node *RBNode, nodes *[]TreeNodeSnapshot, depth int, xMin, xMax float64) {
	if node == t.NIL || node == nil {
		return
	}

	x := (xMin + xMax) / 2
	y := float64(depth*80 + 50)

	snapshot := TreeNodeSnapshot{
		ID:    node.ID,
		Value: node.Value,
		Color: node.Color,
		X:     x,
		Y:     y,
	}

	if node.Left != t.NIL && node.Left != nil {
		leftID := node.Left.ID
		snapshot.LeftID = &leftID
	}
	if node.Right != t.NIL && node.Right != nil {
		rightID := node.Right.ID
		snapshot.RightID = &rightID
	}
	if node.Parent != t.NIL && node.Parent != nil {
		parentID := node.Parent.ID
		snapshot.ParentID = &parentID
	}

	*nodes = append(*nodes, snapshot)

	t.inorderSnapshot(node.Left, nodes, depth+1, xMin, x)
	t.inorderSnapshot(node.Right, nodes, depth+1, x, xMax)
}

// leftRotate performs a left rotation
func (t *RedBlackTree) leftRotate(x *RBNode) {
	y := x.Right
	x.Right = y.Left
	if y.Left != t.NIL {
		y.Left.Parent = x
	}
	y.Parent = x.Parent
	if x.Parent == t.NIL {
		t.Root = y
	} else if x == x.Parent.Left {
		x.Parent.Left = y
	} else {
		x.Parent.Right = y
	}
	y.Left = x
	x.Parent = y

	t.addStep(StepRotateLeft, fmt.Sprintf("对节点 %d 进行左旋", x.Value), &x.ID, []int{x.ID, y.ID})
}

// rightRotate performs a right rotation
func (t *RedBlackTree) rightRotate(y *RBNode) {
	x := y.Left
	y.Left = x.Right
	if x.Right != t.NIL {
		x.Right.Parent = y
	}
	x.Parent = y.Parent
	if y.Parent == t.NIL {
		t.Root = x
	} else if y == y.Parent.Left {
		y.Parent.Left = x
	} else {
		y.Parent.Right = x
	}
	x.Right = y
	y.Parent = x

	t.addStep(StepRotateRight, fmt.Sprintf("对节点 %d 进行右旋", y.Value), &y.ID, []int{x.ID, y.ID})
}

// Insert inserts a value into the Red-Black Tree
func (t *RedBlackTree) Insert(value int) OperationResult {
	t.clearSteps()

	// Create new node
	z := &RBNode{
		ID:     t.nextID,
		Value:  value,
		Color:  Red,
		Left:   t.NIL,
		Right:  t.NIL,
		Parent: t.NIL,
	}
	t.nextID++

	t.addStep(StepInsert, fmt.Sprintf("创建新节点 %d (红色)", value), &z.ID)

	// BST insert
	var y *RBNode = t.NIL
	x := t.Root

	for x != t.NIL {
		y = x
		t.addStep(StepCompare, fmt.Sprintf("比较 %d 与节点 %d", value, x.Value), &x.ID, []int{x.ID})
		if z.Value < x.Value {
			x = x.Left
		} else {
			x = x.Right
		}
	}

	z.Parent = y
	if y == t.NIL {
		t.Root = z
		t.addStep(StepInsert, fmt.Sprintf("节点 %d 成为根节点", value), &z.ID)
	} else if z.Value < y.Value {
		y.Left = z
		t.addStep(StepInsert, fmt.Sprintf("节点 %d 作为 %d 的左子节点", value, y.Value), &z.ID, []int{y.ID, z.ID})
	} else {
		y.Right = z
		t.addStep(StepInsert, fmt.Sprintf("节点 %d 作为 %d 的右子节点", value, y.Value), &z.ID, []int{y.ID, z.ID})
	}

	// Fix Red-Black properties
	t.insertFixup(z)

	t.addStep(StepComplete, "插入完成", nil)

	return OperationResult{
		Success:   true,
		Steps:     t.steps,
		FinalTree: t.getTreeSnapshot(),
	}
}

// insertFixup fixes Red-Black Tree properties after insertion
func (t *RedBlackTree) insertFixup(z *RBNode) {
	for z.Parent != t.NIL && z.Parent.Color == Red {
		if z.Parent == z.Parent.Parent.Left {
			y := z.Parent.Parent.Right // uncle
			if y != t.NIL && y.Color == Red {
				// Case 1: Uncle is red
				t.addStep(StepRebalance, "情况1: 叔节点为红色，重新着色", &z.ID, []int{z.ID, z.Parent.ID, y.ID})
				z.Parent.Color = Black
				y.Color = Black
				z.Parent.Parent.Color = Red
				t.addStep(StepColorChange, fmt.Sprintf("节点 %d, %d 变黑，%d 变红",
					z.Parent.Value, y.Value, z.Parent.Parent.Value), &z.Parent.Parent.ID)
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Right {
					// Case 2: Uncle is black, z is right child
					z = z.Parent
					t.addStep(StepRebalance, "情况2: 叔节点为黑色，当前节点是右子节点", &z.ID)
					t.leftRotate(z)
				}
				// Case 3: Uncle is black, z is left child
				t.addStep(StepRebalance, "情况3: 叔节点为黑色，当前节点是左子节点", &z.ID)
				z.Parent.Color = Black
				z.Parent.Parent.Color = Red
				t.addStep(StepColorChange, fmt.Sprintf("节点 %d 变黑，%d 变红",
					z.Parent.Value, z.Parent.Parent.Value), &z.Parent.ID)
				t.rightRotate(z.Parent.Parent)
			}
		} else {
			// Mirror cases
			y := z.Parent.Parent.Left // uncle
			if y != t.NIL && y.Color == Red {
				t.addStep(StepRebalance, "情况1(镜像): 叔节点为红色，重新着色", &z.ID, []int{z.ID, z.Parent.ID, y.ID})
				z.Parent.Color = Black
				y.Color = Black
				z.Parent.Parent.Color = Red
				t.addStep(StepColorChange, fmt.Sprintf("节点 %d, %d 变黑，%d 变红",
					z.Parent.Value, y.Value, z.Parent.Parent.Value), &z.Parent.Parent.ID)
				z = z.Parent.Parent
			} else {
				if z == z.Parent.Left {
					z = z.Parent
					t.addStep(StepRebalance, "情况2(镜像): 叔节点为黑色，当前节点是左子节点", &z.ID)
					t.rightRotate(z)
				}
				t.addStep(StepRebalance, "情况3(镜像): 叔节点为黑色，当前节点是右子节点", &z.ID)
				z.Parent.Color = Black
				z.Parent.Parent.Color = Red
				t.addStep(StepColorChange, fmt.Sprintf("节点 %d 变黑，%d 变红",
					z.Parent.Value, z.Parent.Parent.Value), &z.Parent.ID)
				t.leftRotate(z.Parent.Parent)
			}
		}
	}
	if t.Root.Color == Red {
		t.Root.Color = Black
		t.addStep(StepColorChange, "根节点变黑", &t.Root.ID)
	}
}

// Search searches for a value in the Red-Black Tree
func (t *RedBlackTree) Search(value int) OperationResult {
	t.clearSteps()

	x := t.Root
	for x != t.NIL {
		t.addStep(StepCompare, fmt.Sprintf("比较 %d 与节点 %d", value, x.Value), &x.ID, []int{x.ID})
		if value == x.Value {
			t.addStep(StepFound, fmt.Sprintf("找到节点 %d", value), &x.ID, []int{x.ID})
			return OperationResult{
				Success:   true,
				Message:   fmt.Sprintf("找到值 %d", value),
				Steps:     t.steps,
				FinalTree: t.getTreeSnapshot(),
			}
		} else if value < x.Value {
			x = x.Left
		} else {
			x = x.Right
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

// transplant replaces subtree rooted at u with subtree rooted at v
func (t *RedBlackTree) transplant(u, v *RBNode) {
	if u.Parent == t.NIL {
		t.Root = v
	} else if u == u.Parent.Left {
		u.Parent.Left = v
	} else {
		u.Parent.Right = v
	}
	v.Parent = u.Parent
}

// minimum finds the minimum node in a subtree
func (t *RedBlackTree) minimum(node *RBNode) *RBNode {
	for node.Left != t.NIL {
		node = node.Left
	}
	return node
}

// searchNode searches for a node with given value
func (t *RedBlackTree) searchNode(value int) *RBNode {
	x := t.Root
	for x != t.NIL {
		if value == x.Value {
			return x
		} else if value < x.Value {
			x = x.Left
		} else {
			x = x.Right
		}
	}
	return t.NIL
}

// Delete deletes a value from the Red-Black Tree
func (t *RedBlackTree) Delete(value int) OperationResult {
	t.clearSteps()

	// Search for the node to delete
	z := t.searchNode(value)

	if z == t.NIL {
		t.addStep(StepNotFound, fmt.Sprintf("值 %d 不存在于树中，无法删除", value), nil)
		return OperationResult{
			Success:   false,
			Message:   fmt.Sprintf("值 %d 不存在，无法删除", value),
			Steps:     t.steps,
			FinalTree: t.getTreeSnapshot(),
		}
	}

	t.addStep(StepDelete, fmt.Sprintf("找到要删除的节点 %d", value), &z.ID, []int{z.ID})

	y := z
	yOriginalColor := y.Color
	var x *RBNode

	if z.Left == t.NIL {
		// Case 1: No left child
		t.addStep(StepDelete, fmt.Sprintf("节点 %d 没有左子节点，用右子节点替换", z.Value), &z.ID)
		x = z.Right
		t.transplant(z, z.Right)
	} else if z.Right == t.NIL {
		// Case 2: No right child
		t.addStep(StepDelete, fmt.Sprintf("节点 %d 没有右子节点，用左子节点替换", z.Value), &z.ID)
		x = z.Left
		t.transplant(z, z.Left)
	} else {
		// Case 3: Both children exist
		y = t.minimum(z.Right)
		yOriginalColor = y.Color
		x = y.Right
		t.addStep(StepDelete, fmt.Sprintf("节点 %d 有两个子节点，找到后继节点 %d", z.Value, y.Value), &y.ID, []int{z.ID, y.ID})

		if y.Parent == z {
			x.Parent = y
		} else {
			t.transplant(y, y.Right)
			y.Right = z.Right
			y.Right.Parent = y
		}
		t.transplant(z, y)
		y.Left = z.Left
		y.Left.Parent = y
		y.Color = z.Color
		t.addStep(StepDelete, fmt.Sprintf("用后继节点 %d 替换被删除节点", y.Value), &y.ID)
	}

	// Fix Red-Black Tree properties if needed
	if yOriginalColor == Black {
		t.addStep(StepRebalance, "删除了黑色节点，需要修复红黑树性质", nil)
		t.deleteFixup(x)
	}

	t.addStep(StepComplete, fmt.Sprintf("删除节点 %d 完成", value), nil)

	return OperationResult{
		Success:   true,
		Message:   fmt.Sprintf("成功删除值 %d", value),
		Steps:     t.steps,
		FinalTree: t.getTreeSnapshot(),
	}
}

// deleteFixup fixes Red-Black Tree properties after deletion
func (t *RedBlackTree) deleteFixup(x *RBNode) {
	for x != t.Root && x.Color == Black {
		if x == x.Parent.Left {
			w := x.Parent.Right // sibling
			if w.Color == Red {
				// Case 1: Sibling is red
				t.addStep(StepRebalance, "情况1: 兄弟节点为红色", &w.ID, []int{x.ID, w.ID})
				w.Color = Black
				x.Parent.Color = Red
				t.addStep(StepColorChange, fmt.Sprintf("兄弟节点 %d 变黑，父节点 %d 变红", w.Value, x.Parent.Value), &x.Parent.ID)
				t.leftRotate(x.Parent)
				w = x.Parent.Right
			}
			if w.Left.Color == Black && w.Right.Color == Black {
				// Case 2: Sibling is black with two black children
				t.addStep(StepRebalance, "情况2: 兄弟节点为黑色，其两个子节点均为黑色", &w.ID, []int{w.ID})
				w.Color = Red
				t.addStep(StepColorChange, fmt.Sprintf("兄弟节点 %d 变红", w.Value), &w.ID)
				x = x.Parent
			} else {
				if w.Right.Color == Black {
					// Case 3: Sibling is black, left child is red, right child is black
					t.addStep(StepRebalance, "情况3: 兄弟节点为黑色，左子为红，右子为黑", &w.ID)
					w.Left.Color = Black
					w.Color = Red
					t.addStep(StepColorChange, fmt.Sprintf("兄弟左子节点变黑，兄弟 %d 变红", w.Value), &w.ID)
					t.rightRotate(w)
					w = x.Parent.Right
				}
				// Case 4: Sibling is black with red right child
				t.addStep(StepRebalance, "情况4: 兄弟节点为黑色，右子为红色", &w.ID)
				w.Color = x.Parent.Color
				x.Parent.Color = Black
				w.Right.Color = Black
				t.addStep(StepColorChange, fmt.Sprintf("重新着色完成"), &w.ID)
				t.leftRotate(x.Parent)
				x = t.Root
			}
		} else {
			// Mirror cases
			w := x.Parent.Left // sibling
			if w.Color == Red {
				t.addStep(StepRebalance, "情况1(镜像): 兄弟节点为红色", &w.ID, []int{x.ID, w.ID})
				w.Color = Black
				x.Parent.Color = Red
				t.addStep(StepColorChange, fmt.Sprintf("兄弟节点 %d 变黑，父节点 %d 变红", w.Value, x.Parent.Value), &x.Parent.ID)
				t.rightRotate(x.Parent)
				w = x.Parent.Left
			}
			if w.Right.Color == Black && w.Left.Color == Black {
				t.addStep(StepRebalance, "情况2(镜像): 兄弟节点为黑色，其两个子节点均为黑色", &w.ID, []int{w.ID})
				w.Color = Red
				t.addStep(StepColorChange, fmt.Sprintf("兄弟节点 %d 变红", w.Value), &w.ID)
				x = x.Parent
			} else {
				if w.Left.Color == Black {
					t.addStep(StepRebalance, "情况3(镜像): 兄弟节点为黑色，右子为红，左子为黑", &w.ID)
					w.Right.Color = Black
					w.Color = Red
					t.addStep(StepColorChange, fmt.Sprintf("兄弟右子节点变黑，兄弟 %d 变红", w.Value), &w.ID)
					t.leftRotate(w)
					w = x.Parent.Left
				}
				t.addStep(StepRebalance, "情况4(镜像): 兄弟节点为黑色，左子为红色", &w.ID)
				w.Color = x.Parent.Color
				x.Parent.Color = Black
				w.Left.Color = Black
				t.addStep(StepColorChange, fmt.Sprintf("重新着色完成"), &w.ID)
				t.rightRotate(x.Parent)
				x = t.Root
			}
		}
	}
	if x.Color == Red {
		x.Color = Black
		t.addStep(StepColorChange, "将当前节点变黑以完成修复", &x.ID)
	}
}
