package bintree2ascii

// jsonTree wraps a jsonNode object so it implements Interface, so we can easily generate ascii art for it.
type jsonTree struct {
	JsonNode *jsonNode
}

// jsonNode represents the deserialization of a tree in json format.
type jsonNode struct {
	Name      string    `jsonTree:"name"`
	Left      *jsonNode `jsonTree:"left"`
	Right     *jsonNode `jsonTree:"right"`
	LeftEdge  string    `jsonTree:"leftEdge"`
	RightEdge string    `jsonTree:"rightEdge"`
}

func (jt jsonTree) Left() Interface {
	if jt.JsonNode.Left == nil {
		return nil
	}
	return jsonTree{jt.JsonNode.Left}
}

func (jt jsonTree) Right() Interface {
	if jt.JsonNode.Right == nil {
		return nil
	}
	return jsonTree{jt.JsonNode.Right}
}

func (jt jsonTree) Key() string {
	return jt.JsonNode.Name
}

func (jt jsonTree) LeftEdge() string {
	return jt.JsonNode.LeftEdge
}

func (jt jsonTree) RightEdge() string {
	return jt.JsonNode.RightEdge
}
