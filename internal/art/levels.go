package art

import (
	"bytes"
)

// Level represents a level of a tree, that needs to be drawn.
type Level struct {
	elements []Element
	indents  []int
}

// Ascii returns the ascii art representation of a level.
func (l *Level) Ascii() []byte {
	var buf []byte
	for {
		var isNewLine bool
		var newLine []byte
		var startIndent int
		for i := 0; i < len(l.elements); i++ {
			newLine = append(newLine, bytes.Repeat([]byte{' '}, l.indents[i]-startIndent)...)
			startIndent = l.indents[i] + l.elements[i].Width()
			line, isValid := l.elements[i].Next()
			if isValid {
				isNewLine = true
			} else {
				line = bytes.Repeat([]byte{' '}, l.elements[i].Width())
			}
			newLine = append(newLine, line...)
		}
		if isNewLine {
			buf = append(buf, newLine...)
			buf = append(buf, '\n')
			continue
		}
		break
	}
	return buf
}

// LevelConfig is config needed for last level
type LevelConfig struct {
	NodeWidth, NodeHeight int
	Distance              int // Distance between two sibling nodes.
	Sep                   int // Distance between two consecutive nodes on the same level, that don't have the same parent.
}

// LastLevel builds the last level of a tree.
func LastLevel(config LevelConfig, elements []Element) *Level {
	indents := make([]int, len(elements))
	for i := 1; i < len(indents); i++ {
		if i%2 == 1 {
			indents[i] = indents[i-1] + elements[i-1].Width() + config.Distance
		} else {
			indents[i] = indents[i-1] + elements[i-1].Width() + config.Sep
		}
	}
	return &Level{
		indents:  indents,
		elements: elements,
	}
}

// ParentLevel builds the parent level of a given level.
func ParentLevel(child *Level, elements []Element) *Level {
	indents := make([]int, len(child.elements)/2)
	for i := 0; i < len(indents); i++ {
		distance := child.indents[2*i+1] - child.indents[2*i] - child.elements[2*i].Width()
		indents[i] = child.indents[2*i] + child.elements[2*i].Width()/2 + distance/2
	}
	return &Level{
		indents:  indents,
		elements: elements,
	}
}

// EdgeLevel builds the dotEdge level between given nodes levels.
func EdgeLevel(parent, child *Level, edgeHeight int, labels []string) *Level {
	indents := make([]int, len(child.elements))
	elements := make([]Element, len(child.elements))
	for i := 0; i < len(parent.elements); i++ {
		//left dotEdge
		indents[2*i] = child.indents[2*i] + child.elements[2*i].Width()/2
		leftEdge := NewLeftEdge(parent.indents[i]-child.indents[2*i]+
			parent.elements[i].Width()/2-child.elements[2*i].Width()/2,
			edgeHeight, labels[2*i], child.elements[2*i].IsInvisible())
		elements[2*i] = leftEdge
		// right dotEdge
		indents[2*i+1] = parent.indents[i] + parent.elements[i].Width()/2
		rightEdge := NewRightEdge(child.indents[2*i+1]-parent.indents[i]+
			child.elements[2*i+1].Width()/2-parent.elements[i].Width()/2,
			edgeHeight, labels[2*i+1], child.elements[2*i+1].IsInvisible())
		elements[2*i+1] = rightEdge
	}
	return &Level{
		elements: elements,
		indents:  indents,
	}
}
