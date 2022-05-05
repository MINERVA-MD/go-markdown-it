package test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go-markdown-it/pkg"
	"testing"
)

func TestShouldReplaceRuleAt(t *testing.T) {
	res := 0
	ruler := pkg.Ruler{
		Rules: nil,
		Cache: nil,
	}

	ruler.Push("test", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		res = 1
		return true
	}, pkg.Rule{})

	ruler.At("test", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		res = 2
		return true
	}, pkg.Rule{})

	rules := ruler.GetRules("")

	assert.Equal(t, 1, len(rules))
	rules[0](nil, nil, nil, 0, 0, false)
	assert.Equal(t, 2, res)
}

func TestShouldInjectBeforeOrAfterRule(t *testing.T) {
	res := 0
	ruler := pkg.Ruler{
		Rules: nil,
		Cache: nil,
	}

	ruler.Push("test", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		res = 1
		return true
	}, pkg.Rule{})

	ruler.Before("test", "before_test", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		res = -10
		return true
	}, pkg.Rule{})

	ruler.After("test", "after_test", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		res = 10
		return true
	}, pkg.Rule{})

	rules := ruler.GetRules("")

	assert.Equal(t, 3, len(rules))
	rules[0](nil, nil, nil, 0, 0, false)
	assert.Equal(t, -10, res)
	rules[1](nil, nil, nil, 0, 0, false)
	assert.Equal(t, 1, res)
	rules[2](nil, nil, nil, 0, 0, false)
	assert.Equal(t, 10, res)
}

func TestShouldEnableDisableRule(t *testing.T) {
	ruler := pkg.Ruler{
		Rules: nil,
		Cache: nil,
	}

	ruler.Push("test", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	ruler.Push("test2", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	rules := ruler.GetRules("")
	assert.Equal(t, 2, len(rules))

	ruler.Disable([]string{"test"}, false)
	rules = ruler.GetRules("")
	assert.Equal(t, 1, len(rules))

	ruler.Disable([]string{"test2"}, false)
	rules = ruler.GetRules("")
	assert.Equal(t, 0, len(rules))

	ruler.Enable([]string{"test"}, false)
	rules = ruler.GetRules("")
	assert.Equal(t, 1, len(rules))

	ruler.Enable([]string{"test2"}, false)
	rules = ruler.GetRules("")
	assert.Equal(t, 2, len(rules))
}

func TestShouldEnableRulesByWhitelist(t *testing.T) {
	ruler := pkg.Ruler{
		Rules: nil,
		Cache: nil,
	}

	ruler.Push("test", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	ruler.Push("test2", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	rules := ruler.GetRules("")
	assert.Equal(t, 2, len(rules))

	ruler.EnableOnly([]string{"test"}, false)
	rules = ruler.GetRules("")
	assert.Equal(t, 1, len(rules))
}

func TestShouldSupportMultipleChains(t *testing.T) {
	ruler := pkg.Ruler{
		Rules: nil,
		Cache: nil,
	}

	ruler.Push("test", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	ruler.Push("test2", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{
		Name:    "",
		Enabled: false,
		Fn:      nil,
		Alt:     []string{"alt1"},
	})

	ruler.Push("test2", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{
		Name:    "",
		Enabled: false,
		Fn:      nil,
		Alt:     []string{"alt1", "alt2"},
	})

	rules := ruler.GetRules("")
	assert.Equal(t, 3, len(rules))

	rules = ruler.GetRules("alt1")
	assert.Equal(t, 2, len(rules))

	rules = ruler.GetRules("alt2")
	assert.Equal(t, 1, len(rules))
}

func TestShouldFailOnInvalidRuleName(t *testing.T) {
	ruler := pkg.Ruler{
		Rules: nil,
		Cache: nil,
	}

	ruler.Push("test", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	baseErrorMsg := "Parser rule not found: "
	err := ruler.At("invalid name", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	assert.Equal(t, errors.New(baseErrorMsg+"invalid name"), err)

	err = ruler.Before("invalid name", "invalid name", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	assert.Equal(t, errors.New(baseErrorMsg+"invalid name"), err)

	err = ruler.Before("invalid name", "invalid name", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	assert.Equal(t, errors.New(baseErrorMsg+"invalid name"), err)

	err = ruler.After("invalid name", "invalid name", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	assert.Equal(t, errors.New(baseErrorMsg+"invalid name"), err)

	_, err = ruler.Enable([]string{"invalid name"}, false)
	assert.Equal(t, errors.New("Rules manager: invalid rule name "+"'invalid name'"), err)

	_, err = ruler.Disable([]string{"invalid name"}, false)
	assert.Equal(t, errors.New("Rules manager: invalid rule name "+"'invalid name'"), err)
}

func TestInvalidShouldNotFailInSilentMode(t *testing.T) {
	ruler := pkg.Ruler{
		Rules: nil,
		Cache: nil,
	}

	ruler.Push("test", func(
		_ *pkg.StateCore,
		_ *pkg.StateBlock,
		_ *pkg.StateInline,
		_ int,
		_ int,
		_ bool,
	) bool {
		return true
	}, pkg.Rule{})

	_, err := ruler.Enable([]string{"invalid name"}, true)
	assert.Equal(t, nil, err)

	err = ruler.EnableOnly([]string{"invalid name"}, true)
	assert.Equal(t, nil, err)

	_, err = ruler.Disable([]string{"invalid name"}, true)
	assert.Equal(t, nil, err)
}
