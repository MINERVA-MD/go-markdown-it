package pkg

type Attribute struct {
	Name  string
	Value string
}

type PluginMetaData struct {
	Delimiters []string
}

type Token struct {
	// Type of the token (string, e.g. "paragraph_open")
	Type string

	// HTML Tag name, e.g. "p"
	Tag string

	// Html attributes. Format: `[ [ name1, value1 ], [ name2, value2 ]
	Attrs []Attribute

	// Source map info. Format: `[ line_begin, line_end ]`
	Map []int

	/**
	 * Level change (number in {-1, 0, 1} set), where:
	 *
	 * -  `1` means the tag is opening
	 * -  `0` means the tag is self-closing
	 * - `-1` means the tag is closing
	 **/
	Nesting int

	// nesting level, the same as `state.level
	Level int

	// An array of child nodes (inline and img tokens)
	Children []*Token

	// In a case of self-closing tag (code, html, fence, etc.),
	// it has contents of this tag.
	Content string

	// '*' or '_' for emphasis, fence string for fence, etc.
	Markup string

	/**
	 * Additional information:
	 *
	 * - Info string for "fence" tokens
	 * - The value "auto" for autolink "link_open" and "link_close" tokens
	 * - The string value of the item marker for ordered-list "list_item_open" tokens
	 **/
	Info string

	// A place for plugins to store an arbitrary data
	Meta PluginMetaData

	// True for block-level tokens, false for inline tokens.
	//  Used in renderer to calculate line breaks
	Block bool

	// If it's true, ignore this element when rendering. Used for tight lists
	// to hide paragraphs.
	Hidden bool
}

// AttrIndex Searches for attribute index by name.
// It returns the index if found, otherwise -1;
func (t *Token) AttrIndex(name string) int {
	if t.Attrs == nil {
		return -1
	}

	for i, s := range t.Attrs {
		if s.Name == name {
			return i
		}
	}

	return -1
}

// AttrPush adds attribute to list. Initializes if necessary.
func (t *Token) AttrPush(attribute Attribute) {
	if t.Attrs != nil {
		t.Attrs = append(t.Attrs, attribute)
	} else {
		t.Attrs = []Attribute{attribute}
	}
}

// AttrSet sets attribute.
// Overrides old value, if exists.
func (t *Token) AttrSet(attribute Attribute) {
	var idx = t.AttrIndex(attribute.Name)

	if idx < 0 {
		t.Attrs = append(t.Attrs, attribute)
	} else {
		t.Attrs[idx] = attribute
	}
}

// AttrGet gets the value of attribute `name`.
// Returns nil if it does not exist, otherwise the value of Attribute
func (t *Token) AttrGet(name string) (string, bool) {
	var value string
	var idx = t.AttrIndex(name)

	if idx >= 0 {
		value = t.Attrs[idx].Value
	} else {
		return "", false
	}

	return value, true
}

// AttrJoin joins attribute.Value to existing attribute via space.
// Or create new attribute if not exists.
func (t *Token) AttrJoin(attribute Attribute) {
	var idx = t.AttrIndex(attribute.Name)

	if idx < 0 {
		t.Attrs = append(t.Attrs, attribute)
	} else {
		t.Attrs[idx].Value += " " + attribute.Value
	}
}

func GenerateToken(Type string, Tag string, Nesting int) Token {
	return Token{
		Type:     Type,
		Tag:      Tag,
		Attrs:    nil,
		Map:      []int{},
		Nesting:  Nesting,
		Level:    0,
		Children: nil,
		Content:  "",
		Markup:   "",
		Info:     "",
		Meta:     PluginMetaData{},
		Block:    false,
		Hidden:   false,
	}
}
