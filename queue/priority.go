package queue

import (
	"errors"

	"github.com/dwayn/gofirst/stats"
)

type itemNode struct {
	itemId string
	prev   *itemNode
	next   *itemNode
}

type scoreTreeNode struct {
	left  *scoreTreeNode
	right *scoreTreeNode
	head  *itemNode
	tail  *itemNode
	score int
}

type itemTreeNode struct {
	left  *itemTreeNode
	right *itemTreeNode
	item  *itemNode
	score int
}

type PriorityQueue struct {
	scoreRoot *scoreTreeNode
	itemRoot  *itemTreeNode
	Metrics   chan stats.Metric
}

func (q *PriorityQueue) PeekNext() (string, int) {
	snode := q.findMaxScore(q.scoreRoot)
	if snode == nil {
		return "", QueueEmpty
	}
	return snode.head.itemId, Ok
}

func (q *PriorityQueue) GetNext() (string, int) {
	snode := q.findMaxScore(q.scoreRoot)
	if snode == nil {
		return "", QueueEmpty
	}
	rval := snode.head.itemId
	itNode := q.findItem(snode.head.itemId, q.itemRoot)
	stillPopulated := q.removeItemNode(snode, itNode.item)
	if !stillPopulated {
		q.scoreRoot = q.deleteScoreTreeNode(q.scoreRoot, snode)
	}
	itNode.score = 0
	// tmp := itNode.item
	q.itemRoot = q.deleteItemTreeNode(q.itemRoot, itNode)
	// ######################################################################### stats.items--
	q.Metrics <- stats.Metric{Metric: "items", Op: stats.Sub, Value: 1}
	return rval, Ok
}

func (q *PriorityQueue) GetScore(itemId string) (int, int) {
	itNode := q.findItem(itemId, q.itemRoot)
	if itNode == nil {
		return -1, NotFound
	}
	return itNode.score, Ok
}

// returns 1 on adding item, 0 on  update of item
func (q *PriorityQueue) Update(itemId string, score int) (int, int) {
	rval := 0
	newScore := score
	itNode := q.findItem(itemId, q.itemRoot)
	if itNode == nil {
		rval = 1
		itNode = &itemTreeNode{item: &itemNode{itemId: itemId}}
		// ########################################################################### stats.items++
		q.Metrics <- stats.Metric{Metric: "items", Op: stats.Add, Value: 1}
		itNode.score = 0
		q.itemRoot = q.addItemTreeNode(q.itemRoot, itNode)
	}

	var sNode *scoreTreeNode
	if itNode.score > 0 {
		sNode = q.findScore(itNode.score, q.scoreRoot)
		newScore = score + itNode.score
		populated := q.removeItemNode(sNode, itNode.item)
		if !populated {
			q.scoreRoot = q.deleteScoreTreeNode(q.scoreRoot, sNode)
		}
		itNode.score = 0
	}
	sNode = q.findScore(newScore, q.scoreRoot)
	if sNode == nil {
		sNode = &scoreTreeNode{score: newScore}
		q.scoreRoot = q.addScoreTreeNode(q.scoreRoot, sNode)
	}
	q.addItemNode(sNode, itNode.item)
	itNode.score = sNode.score
	// ############################################################################### stats.updates++
	q.Metrics <- stats.Metric{Metric: "updates", Op: stats.Add, Value: 1}
	// TODO: is there error cases here that should be returned, or is the only error that can happen in here a panic due to unable to allocate memory
	return rval, Ok
}

func (q *PriorityQueue) findItem(itemId string, tree *itemTreeNode) *itemTreeNode {
	if tree == nil {
		return nil
	}
	switch {
	case itemId < tree.item.itemId:
		return q.findItem(itemId, tree.left)
	case itemId > tree.item.itemId:
		return q.findItem(itemId, tree.right)
	default:
		return tree
	}
}

// returns true if there is still items in the list, 0 if empty
// assumes that the itemNode is actually in the scoreTreeNode's list (does not validate the existence here)
func (q *PriorityQueue) removeItemNode(list *scoreTreeNode, iNode *itemNode) bool {
	// if item is last thing in list
	if list.head == iNode && list.tail == iNode {
		list.head = nil
		list.tail = nil
		return false
	}

	// if item is the head node in list
	if list.head == iNode {
		list.head = list.head.next
		list.head.prev = nil
		iNode.next = nil
		return true
	}

	// if item is tail node in list
	if list.tail == iNode {
		list.tail = list.tail.prev
		list.tail.next = nil
		iNode.prev = nil
		return true
	}

	// node is somewhere in the middle of the list
	iNode.next.prev = iNode.prev
	iNode.prev.next = iNode.next
	iNode.prev = nil
	iNode.next = nil
	return true
}

func (q *PriorityQueue) deleteScoreTreeNode(tree *scoreTreeNode, node *scoreTreeNode) *scoreTreeNode {
	if !(node.head == nil && node.tail == nil) {
		panic(errors.New("Tried to free scoreTreeNode that still has itemNodes in list"))
	}

	if tree == nil {
		panic(errors.New("Attempted to delete node from tree that is not in the tree"))
	}

	switch {
	case node.score < tree.score:
		tree.left = q.deleteScoreTreeNode(tree.left, node)
	case node.score > tree.score:
		tree.right = q.deleteScoreTreeNode(tree.right, node)
	case tree.left != nil && tree.right != nil:
		tmpCell := q.findMinScore(tree.right)
		tree.head = tmpCell.head
		tree.tail = tmpCell.tail
		tmpscore := tree.score
		tree.score = tmpCell.score
		tmpCell.score = tmpscore
		tmpCell.head = nil
		tmpCell.tail = nil
		tree.right = q.deleteScoreTreeNode(tree.right, tmpCell)
	default:
		switch {
		case tree.left == nil:
			tree = tree.right
		case tree.right == nil:
			tree = tree.left
		}
		// ############################################################################ stats.pools--
		q.Metrics <- stats.Metric{Metric: "pools", Op: stats.Sub, Value: 1}
	}
	return tree
}

func (q *PriorityQueue) deleteItemTreeNode(tree *itemTreeNode, node *itemTreeNode) *itemTreeNode {
	if tree == nil {
		panic(errors.New("Attempt to delete node from empty tree"))
	}
	switch {
	case node.item.itemId < tree.item.itemId:
		tree.left = q.deleteItemTreeNode(tree.left, node)
	case node.item.itemId > tree.item.itemId:
		tree.right = q.deleteItemTreeNode(tree.right, node)
	case tree.left != nil && tree.right != nil:
		tmpCell := q.findMinItem(tree.right)
		tmpItem := tree.item
		tree.item = tmpCell.item
		tree.score = tmpCell.score
		tmpCell.item = tmpItem
		tmpCell.score = 0
		tree.right = q.deleteItemTreeNode(tree.right, tmpCell)
	default:
		// tmpCell := tree
		switch {
		case tree.left == nil:
			tree = tree.right
		case tree.right == nil:
			tree = tree.left
		}
	}
	return tree
}
func (q *PriorityQueue) findScore(score int, tree *scoreTreeNode) *scoreTreeNode {
	if tree == nil {
		return nil
	}
	switch {
	case score < tree.score:
		return q.findScore(score, tree.left)
	case score > tree.score:
		return q.findScore(score, tree.right)
	}
	return tree
}
func (q *PriorityQueue) addItemTreeNode(tree *itemTreeNode, node *itemTreeNode) *itemTreeNode {
	if tree == nil {
		tree = node
	} else {
		switch {
		case node.item.itemId < tree.item.itemId:
			tree.left = q.addItemTreeNode(tree.left, node)
		case node.item.itemId > tree.item.itemId:
			tree.right = q.addItemTreeNode(tree.right, node)
		}
	}
	return tree
}

func (q *PriorityQueue) addScoreTreeNode(tree *scoreTreeNode, node *scoreTreeNode) *scoreTreeNode {
	if tree == nil {
		tree = node
		// ########################################################################## stats.pools++
		q.Metrics <- stats.Metric{Metric: "pools", Op: stats.Add, Value: 1}
	} else {
		switch {
		case node.score < tree.score:
			tree.left = q.addScoreTreeNode(tree.left, node)
		case node.score > tree.score:
			tree.right = q.addScoreTreeNode(tree.right, node)
		}
	}
	return tree
}

func (q *PriorityQueue) addItemNode(sNode *scoreTreeNode, iNode *itemNode) {
	if sNode.head == nil {
		sNode.head = iNode
		sNode.tail = iNode
		return
	}
	sNode.tail.next = iNode
	iNode.prev = sNode.tail
	sNode.tail = iNode
}

func (q *PriorityQueue) findMaxScore(node *scoreTreeNode) *scoreTreeNode {
	if node == nil {
		return nil
	}
	tmpCell := node
	for {
		if tmpCell.right == nil {
			break
		}
		tmpCell = tmpCell.right
	}
	return tmpCell
}

func (q *PriorityQueue) findMinScore(node *scoreTreeNode) *scoreTreeNode {
	if node == nil {
		return nil
	}
	tmpCell := node
	for {
		if tmpCell.left == nil {
			break
		}
		tmpCell = tmpCell.left
	}
	return tmpCell
}

func (q *PriorityQueue) findMinItem(node *itemTreeNode) *itemTreeNode {
	if node == nil {
		return nil
	}
	tmpCell := node
	for {
		if tmpCell.left == nil {
			break
		}
		tmpCell = tmpCell.left
	}
	return tmpCell
}

func (q *PriorityQueue) findMaxItem(node *itemTreeNode) *itemTreeNode {
	if node == nil {
		return nil
	}
	tmpCell := node
	for {
		if tmpCell.right == nil {
			break
		}
		tmpCell = tmpCell.right
	}
	return tmpCell
}
