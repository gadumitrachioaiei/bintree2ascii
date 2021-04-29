package bintree2ascii

import (
	"errors"

	"gonum.org/v1/gonum/graph/formats/dot"
	"gonum.org/v1/gonum/graph/formats/dot/ast"
)

// dotTree represents a tree read from a dot format.
type dotTree struct {
	adj map[string]*[2]*dotEdge
}

// dotEdge represents an edge in the dotTree.
type dotEdge struct {
	to    string
	label string
}

// direction represents direction of an edge from the bin tree.
type direction string

const (
	left  direction = "left"
	right           = "right"
)

// newDotTree reads a tree from a dot format.
func newDotTree(s string) (dotTree, error) {
	f, err := dot.ParseString(s)
	if err != nil {
		return dotTree{}, err
	}
	if len(f.Graphs) != 1 {
		return dotTree{}, errors.New("we need exactly one graph")
	}
	graph := f.Graphs[0]
	// we want to make the adjacency matrix
	adj := make(map[string]*[2]*dotEdge)
	for _, stmt := range graph.Stmts {
		edgeStmt, ok := stmt.(*ast.EdgeStmt)
		if !ok {
			continue
		}
		fromNode, ok := edgeStmt.From.(*ast.Node)
		if !ok {
			continue
		}
		toNode, ok := edgeStmt.To.Vertex.(*ast.Node)
		if !ok {
			continue
		}
		var edgeIndex int
		var label string
		for _, attr := range edgeStmt.Attrs {
			if attr.Key == "direction" && attr.Val == right {
				edgeIndex = 1
			} else if attr.Key == "label" {
				label = attr.Val
			}
		}
		if adj[fromNode.ID] == nil {
			adj[fromNode.ID] = &[2]*dotEdge{}
		}
		if _, ok := adj[toNode.ID]; !ok {
			adj[toNode.ID] = nil
		}
		adj[fromNode.ID][edgeIndex] = &dotEdge{to: toNode.ID, label: label}
	}
	return dotTree{adj: adj}, nil
}

// dotTreeToInterface returns a tree that implements the Interface, so we can draw it as ascii art.
func (dt dotTree) dotTreeToInterface() jsonTree {
	adj := dt.adj
	var toTreeNode func(string) *jsonNode
	toTreeNode = func(node string) *jsonNode {
		if adj[node] == nil {
			return &jsonNode{
				Name: node,
			}
		}
		var left, right *jsonNode
		var leftEdge, rightEdge string
		if adj[node][0] != nil {
			left = toTreeNode(adj[node][0].to)
			leftEdge = adj[node][0].label
		}
		if adj[node][1] != nil {
			right = toTreeNode(adj[node][1].to)
			rightEdge = adj[node][1].label
		}
		return &jsonNode{
			Name:      node,
			Left:      left,
			Right:     right,
			LeftEdge:  leftEdge,
			RightEdge: rightEdge,
		}
	}
	return jsonTree{toTreeNode(dt.getRoot())}
}

// getRoot returns the root node.
func (dt dotTree) getRoot() string {
	adj := dt.adj
	var root string
	for start := range adj {
		isRoot := true
		for _, edges := range adj {
			if edges == nil {
				continue
			}
			for _, edge := range *edges {
				if edge != nil && edge.to == start {
					isRoot = false
					break
				}
			}
			if !isRoot {
				break
			}
		}
		if isRoot {
			root = start
			break
		}
	}
	return root
}
