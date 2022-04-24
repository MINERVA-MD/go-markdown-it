package pkg

import . "go-markdown-it/pkg/types"

// AttrIndex Searches for attribute index by name.
// It returns the index if found, otherwise -1;
func AttrIndex(token *Token, name string) int {
	if token.Attrs == nil {
		return -1
	}

	for i, s := range token.Attrs {
		if s.Name == name {
			return i
		}
	}

	return -1
}

// AttrPush adds attribute to list. Initializes if necessary.
func AttrPush(token *Token, attribute Attribute) {
	if token.Attrs != nil {
		token.Attrs = append(token.Attrs, attribute)
	} else {
		token.Attrs = []Attribute{attribute}
	}
}

// AttrSet sets attribute.
// Overrides old value, if exists.
func AttrSet(token *Token, attribute Attribute) {
	var idx = AttrIndex(token, attribute.Name)

	if idx < 0 {
		token.Attrs = append(token.Attrs, attribute)
	} else {
		token.Attrs[idx] = attribute
	}
}

// AttrGet gets the value of attribute `name`.
// Returns nil if it does not exist, otherwise the value of Attribute
func AttrGet(token *Token, name string) string {
	var value string
	var idx = AttrIndex(token, name)

	if idx >= 0 {
		value = token.Attrs[idx].Value
	}

	return value
}

// AttrJoin joins attribute.Value to existing attribute via space.
// Or create new attribute if not exists.
func AttrJoin(token *Token, attribute Attribute) {
	var idx int = AttrIndex(token, attribute.Name)

	if idx < 0 {
		token.Attrs = append(token.Attrs, attribute)
	} else {
		token.Attrs[idx].Value += " " + attribute.Value
	}
}
