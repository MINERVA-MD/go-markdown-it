package pkg

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

type Cache map[string][]RuleFunction
type RuleFunction func(*StateCore, *StateBlock, *StateInline, int, int, bool) bool
type Rule struct {
	Name    string
	Enabled bool
	Fn      RuleFunction
	Alt     []string
	// TODO: convert ^this tuple to its own type
}

type Ruler struct {
	Rules []Rule
	Cache Cache
}

// Helper methods

// Find finds rule index by name
// Returns the index if found, otherwise -1.
func (ruler *Ruler) Find(name string) int {
	for idx, rule := range ruler.Rules {
		if rule.Name == name {
			return idx
		}
	}
	return -1
}

func (ruler *Ruler) Compile() {
	var chains = []string{""}

	// Collect unique names
	for _, rule := range ruler.Rules {
		if !rule.Enabled {
			continue
		}

		for _, alt := range rule.Alt {
			if IndexOf(alt, chains) < 0 {
				chains = append(chains, alt)
			}
		}
	}

	ruler.Cache = Cache{}

	for _, chain := range chains {
		ruler.Cache[chain] = []RuleFunction{}

		for _, rule := range ruler.Rules {
			if !rule.Enabled {
				continue
			}
			if utf8.RuneCountInString(chain) > 0 && IndexOf(chain, rule.Alt) < 0 {
				continue
			}

			ruler.Cache[chain] = append(ruler.Cache[chain], rule.Fn)
		}
	}
}

func (ruler *Ruler) At(name string, fn RuleFunction, options Rule) error {
	var idx = ruler.Find(name)

	if idx == -1 {
		return errors.New(fmt.Sprintf("Parser rule not found: %s", name))
	}

	ruler.Rules[idx].Fn = fn

	if options.Alt != nil {
		ruler.Rules[idx].Alt = options.Alt
	} else {
		ruler.Rules[idx].Alt = []string{}
	}

	ruler.Cache = nil
	return nil
}

// Before adds new rule to chain before one with given name.
func (ruler *Ruler) Before(
	beforeName string,
	ruleName string,
	fn RuleFunction,
	options Rule) error {

	var idx = ruler.Find(beforeName)
	return ruler.Insert(idx, options, ruleName, fn)
}

// After adds new rule to chain after one with given name.
func (ruler *Ruler) After(
	afterName string,
	ruleName string,
	fn RuleFunction,
	options Rule) error {

	var idx = ruler.Find(afterName)

	if idx == -1 {
		return errors.New(fmt.Sprintf("Parser rule not found: %s", ruleName))
	}

	return ruler.Insert(idx+1, options, ruleName, fn)
}

func (ruler *Ruler) Insert(idx int, options Rule, ruleName string, fn RuleFunction) error {
	if idx == -1 {
		return errors.New(fmt.Sprintf("Parser rule not found: %s", ruleName))
	}

	var alt []string
	if options.Alt != nil {
		alt = options.Alt
	} else {
		alt = []string{}
	}

	ruler.Rules = ruler.InsertAt(ruler.Rules, idx, Rule{
		Name:    ruleName,
		Enabled: true,
		Fn:      fn,
		Alt:     alt,
	})

	ruler.Cache = nil
	return nil
}

func (ruler *Ruler) Push(
	ruleName string,
	fn RuleFunction,
	options Rule,
) {

	var alt []string
	if options.Alt != nil {
		alt = options.Alt
	} else {
		alt = []string{}
	}

	ruler.Rules = append(ruler.Rules, Rule{
		Name:    ruleName,
		Enabled: true,
		Fn:      fn,
		Alt:     alt,
	})
}

func (ruler *Ruler) Enable(list []string, ignoreInvalid bool) ([]string, error) {
	var result []string

	// Search by name and enable
	for _, name := range list {
		var idx = ruler.Find(name)

		if idx < 0 {
			if ignoreInvalid {
				continue
			}
			return nil, errors.New(fmt.Sprintf("Rules manager: invalid rule name '%s'", name))
		}

		ruler.Rules[idx].Enabled = true
		result = append(result, name)
	}

	ruler.Cache = nil
	return result, nil
}

func (ruler *Ruler) EnableOnly(list []string, ignoreInvalid bool) error {

	for i := 0; i < len(ruler.Rules); i++ {
		ruler.Rules[i].Enabled = false
	}

	_, err := ruler.Enable(list, ignoreInvalid)

	return err
}

func (ruler *Ruler) Disable(list []string, ignoreInvalid bool) ([]string, error) {
	var result []string

	// Search by name and enable
	for _, name := range list {
		var idx = ruler.Find(name)

		if idx < 0 {
			if ignoreInvalid {
				continue
			}
			return nil, errors.New(fmt.Sprintf("Rules manager: invalid rule name '%s'", name))
		}

		ruler.Rules[idx].Enabled = false
		result = append(result, name)
	}

	ruler.Cache = nil
	return result, nil
}

func (ruler *Ruler) GetRules(chainName string) []RuleFunction {
	if ruler.Cache == nil {
		ruler.Compile()
	}

	// Chain can be empty, if rules disabled. But we still have to return Array.
	if ruler.Cache[chainName] != nil {
		return ruler.Cache[chainName]
	}
	return []RuleFunction{}
}

func (ruler *Ruler) InsertAt(rules []Rule, idx int, rule Rule) []Rule {
	if idx >= len(rules) { // nil or empty slice or after last element
		return append(rules, rule)
	}
	rules = append(rules[:idx+1], rules[idx:]...) // index < len(a)
	rules[idx] = rule
	return rules
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
