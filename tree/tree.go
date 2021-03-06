package tree

import "errors"

const (
	IsRed       bool = true
	IsBlack     bool = false
	LeftRotate  bool = true
	RightRotate bool = false
)

// 定义树，默认只包含一个根节点
type RBTree struct {
	root *RBNode
}

func (t *RBTree) insert(current *RBNode, key string, value []byte) {
	hash := hash(key)
	if hash < current.hash {
		if current.left == nil {
			node := NewNode(key, value, hash)
			current.left = node
			node.parent = current
			t.check(node)
		} else {
			t.insert(current.left, key, value)
		}
	} else if hash > current.hash {
		if current.right == nil {
			node := NewNode(key, value, hash)
			current.right = node
			node.parent = current
			t.check(node)
		} else {
			t.insert(current.right, key, value)
		}
	} else {
		current.value = value
	}
}

// 根据红黑树规则 校验和调整树形结构
func (t *RBTree) check(node *RBNode) {

}

// 定义节点
type RBNode struct {
	hash                uint32
	key                 string
	value               []byte
	color               bool
	left, right, parent *RBNode
}

func NewNode(key string, value []byte, hash uint32) (tNode *RBNode) {
	tNode = new(RBNode)
	tNode.key = key
	tNode.value = value
	tNode.hash = hash
	return
}

func (tn *RBNode) getParent() *RBNode {
	return tn.parent
}

// 树旋转，如果有根节点变动则返回变动后的根节点
func (tn *RBNode) rotate(leftRotate bool) (*RBNode, error) {
	var root *RBNode
	if tn == nil {
		return root, nil
	}
	if leftRotate && tn.right == nil {
		return root, errors.New("right node must not be nil")
	}
	if !leftRotate && tn.left == nil {
		return root, errors.New("left node must not be nil")
	}
	parent := tn.parent
	var isLeft bool
	if parent != nil {
		isLeft = parent.left == tn
	}
	if leftRotate { // 左选
		grandson := tn.right.left
		tn.parent = tn.right
		tn.parent.left = tn
		tn.right = grandson
	} else { // 右旋
		grandson := tn.left.right
		tn.parent = tn.left
		tn.parent.right = tn
		tn.left = grandson
	}
	if parent == nil {
		tn.parent.parent = nil
		root = tn.parent
	} else {
		if isLeft {
			parent.left = tn.parent
		} else {
			parent.right = tn.parent
		}
		tn.parent.parent = parent

	}
	return root, nil
}

func hash(s string) uint32 {
	return 0 // TODO
}
