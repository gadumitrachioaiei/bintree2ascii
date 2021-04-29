// Package bintree2ascii generates ascii art for binary trees, especially bsts.
//
// If we will ever want to use fancy looks, look here:
//
// http://efanzh.org/tree-graph-generator/
// https://en.wikipedia.org/wiki/Box-drawing_character
package bintree2ascii

import (
	"encoding/json"

	"github.com/gadumitrachioaiei/bintree2ascii/internal/art"
)

// AsciiArt generates ascii representation for a binary tree.
type AsciiArt struct {
	Config

	levels []*art.Level
}

// Config represents configuration for the AsciiArt.
type Config struct {
	NodeWidth  int // width of a node's content
	NodeHeight int // height of a node's content
	EdgeHeight int // height of a dotEdge
	Distance   int // Distance between two sibling nodes
	Sep        int // Distance between two consecutive nodes that are not siblings
}

// NewAsciiArt returns a new configured AsciiArt.
func NewAsciiArt(config Config) *AsciiArt {
	return &AsciiArt{
		Config: config,
	}
}

// Interface needs to be implemented by your tree object in ordered to be represented as an ascii art.
type Interface interface {
	Left() Interface
	Right() Interface
	Key() string
	LeftEdge() string
	RightEdge() string
}

// FromInterface returns the ascii representation of the interface.
func (at *AsciiArt) FromInterface(i Interface) string {
	treeLevels := levels(i)
	depth := len(treeLevels)
	asciiLevels := make([]*art.Level, 2*depth-1)
	asciiLevels[len(asciiLevels)-1] = art.LastLevel(art.LevelConfig{
		NodeWidth:  at.NodeWidth,
		NodeHeight: at.NodeHeight,
		Distance:   at.Distance,
		Sep:        at.Sep,
	}, at.asciiLevel(treeLevels[depth-1]))
	for i := len(asciiLevels) - 1; i >= 2; i -= 2 {
		// current tree level for this ascii level
		treeLevel := treeLevels[(i-2)/2]
		// node level
		asciiLevels[i-2] = art.ParentLevel(asciiLevels[i], at.asciiLevel(treeLevel))
		// dotEdge level
		asciiLevels[i-1] = art.EdgeLevel(asciiLevels[i-2], asciiLevels[i], at.EdgeHeight, edgeLabels(treeLevel))
	}
	at.levels = asciiLevels
	return string(at.ascii())
}

// FromJson returns the ascii representation of a jsonTree serialization of a bin tree.
//
// See the jsonNode type for naming conventions of jsonTree representation.
func (at *AsciiArt) FromJson(s []byte) (string, error) {
	var tree jsonNode
	if err := json.Unmarshal(s, &tree); err != nil {
		return "", err
	}
	return at.FromInterface(jsonTree{&tree}), nil
}

// FromDot returns the ascii representation of a tree given its adjacency data, in graphviz dot format.
// A dotEdge should look like this: 1 -> 2[label="labelEdge", direction="right or left"].
func (at *AsciiArt) FromDot(s string) (string, error) {
	dt, err := newDotTree(s)
	if err != nil {
		return "", err
	}
	return at.FromInterface(dt.dotTreeToInterface()), nil
}

// ascii returns the ascii representation of the stored tree.
func (at *AsciiArt) ascii() []byte {
	var result []byte
	for i := 0; i < len(at.levels); i++ {
		result = append(result, at.levels[i].Ascii()...)
	}
	return result
}

// asciiLevel converts a level of tree nodes into a list of nodes that we can draw as ascii.
func (at *AsciiArt) asciiLevel(treeLevel []Interface) []art.Element {
	var nodes []art.Element
	for i := 0; i < len(treeLevel); i++ {
		var node *art.Node
		if treeLevel[i] == nil {
			node = art.NewInvisibleNode(at.NodeWidth, at.NodeHeight)
		} else {
			node = art.NewNode(treeLevel[i].Key(), at.NodeWidth, at.NodeHeight)
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// edgeLabels returns the edges' labels of the given nodes.
func edgeLabels(level []Interface) []string {
	labels := make([]string, 2*len(level))
	for i := 0; i < len(level); i++ {
		if level[i] != nil {
			labels[2*i] = level[i].LeftEdge()
			labels[2*i+1] = level[i].RightEdge()
		}
	}
	return labels
}

// levels returns the list of levels of the tree.
func levels(i Interface) [][]Interface {
	var levels [][]Interface
	parentLevel := []Interface{i}
	levels = append(levels, parentLevel)
	for {
		isLastLevel := true
		var childLevel []Interface
		for i := 0; i < len(parentLevel); i++ {
			if parentLevel[i] == nil {
				childLevel = append(childLevel, nil, nil)
			} else {
				left, right := parentLevel[i].Left(), parentLevel[i].Right()
				if left != nil || right != nil {
					isLastLevel = false
				}
				childLevel = append(childLevel, left, right)
			}
		}
		if isLastLevel {
			break
		}
		levels = append(levels, childLevel)
		parentLevel = childLevel
	}
	return levels
}
