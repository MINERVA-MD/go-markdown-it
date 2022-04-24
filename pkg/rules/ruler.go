package rules

import (
	"errors"
	"fmt"
	. "go-markdown-it/pkg/types"
)

// Helper methods

// RulerFind finds rule index by name
// Returns the index if found, otherwise -1.
func RulerFind(ruler *Ruler, name string) int {
	for idx, rule := range ruler.Rules {
		if rule.Name == name {
			return idx
		}
	}
	return -1
}

func RulerCompile(ruler *Ruler) {
	var chains = []string{""}

	// Collect unique names
	for _, rule := range ruler.Rules {
		if !rule.Enabled {
			return
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
				return
			}
			if len(chain) > 0 && IndexOf(chain, rule.Alt) < 0 {
				return
			}

			ruler.Cache[chain] = append(ruler.Cache[chain], rule.Fn)
		}
	}
}

func RulerAt(ruler *Ruler, name string, fn RuleFunction, options Rule) error {
	var idx = RulerFind(ruler, name)

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

// RulerBefore adds new rule to chain before one with given name.
func RulerBefore(
	ruler *Ruler,
	beforeName string,
	ruleName string,
	fn RuleFunction,
	options Rule) error {

	var idx = RulerFind(ruler, beforeName)
	return RulerInsert(ruler, idx, options, ruleName, fn)
}

// RulerAfter adds new rule to chain after one with given name.
func RulerAfter(
	ruler *Ruler,
	afterName string,
	ruleName string,
	fn RuleFunction,
	options Rule) error {

	var idx = RulerFind(ruler, afterName)
	return RulerInsert(ruler, idx+1, options, ruleName, fn)
}

func RulerInsert(ruler *Ruler, idx int, options Rule, ruleName string, fn RuleFunction) error {
	if idx == -1 {
		return errors.New(fmt.Sprintf("Parser rule not found: %s", ruleName))
	}

	var alt []string
	if options.Alt != nil {
		alt = options.Alt
	} else {
		alt = []string{}
	}

	ruler.Rules = InsertAt(ruler.Rules, idx, Rule{
		Name:    ruleName,
		Enabled: true,
		Fn:      fn,
		Alt:     alt,
	})

	ruler.Cache = nil
	return nil
}

func RulerPush(
	ruler *Ruler,
	ruleName string,
	fn RuleFunction,
	options Rule) {

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

func RulerEnable(ruler *Ruler, list []string, ignoreInvalid bool) ([]string, error) {
	var result []string

	// Search by name and enable
	for _, name := range list {
		var idx = RulerFind(ruler, name)

		if idx < 0 {
			if ignoreInvalid {
				return nil, nil
			}
			return nil, errors.New(fmt.Sprintf("Rules manager: invalid rule name %s", name))
		}

		ruler.Rules[idx].Enabled = true
		result = append(result, name)
	}

	ruler.Cache = nil
	return result, nil
}

func RulerEnableOnly(ruler *Ruler, list []string, ignoreInvalid bool) {

	for _, rule := range ruler.Rules {
		rule.Enabled = false
	}

	_, _ = RulerEnable(ruler, list, ignoreInvalid)
}

func RulerDisable(ruler *Ruler, list []string, ignoreInvalid bool) ([]string, error) {
	var result []string

	// Search by name and enable
	for _, name := range list {
		var idx = RulerFind(ruler, name)

		if idx < 0 {
			if ignoreInvalid {
				return nil, nil
			}
			return nil, errors.New(fmt.Sprintf("Rules manager: invalid rule name %s", name))
		}

		ruler.Rules[idx].Enabled = false
		result = append(result, name)
	}

	ruler.Cache = nil
	return result, nil
}

func RulerGetRules(ruler *Ruler, chainName string) []RuleFunction {
	if ruler.Cache == nil {
		RulerCompile(ruler)
	}

	// Chain can be empty, if rules disabled. But we still have to return Array.
	if ruler.Cache[chainName] != nil {
		return ruler.Cache[chainName]
	}
	return []RuleFunction{}
}

func InsertAt(rules []Rule, idx int, rule Rule) []Rule {
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
