package cistern

import "errors"

const (
	IsRed       bool = true
	IsBlack     bool = false
	LeftRotate  bool = true
	RightRotate bool = false
)

// 定义树，默认只包含一个根节点
type Tree struct {
	root *TreeNode
}

func (t *Tree) insert(node *TreeNode) {

}

// 定义节点
type TreeNode struct {
	value               []byte
	key                 string
	color               bool
	left, right, parent *TreeNode
}

func (tn *TreeNode) getParent() *TreeNode {
	return tn.parent
}

// 树旋转，如果有根节点变动则返回变动后的根节点
func (tn *TreeNode) rotate(leftRotate bool) (*TreeNode, error) {
	var root *TreeNode
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
