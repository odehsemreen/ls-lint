package main

import (
	"sync"
	"unicode"
)

type RuleSnakeCase struct {
	Name string
	*sync.RWMutex
}

func (rule *RuleSnakeCase) Init() Rule {
	rule.Name = "snakecase"
	rule.RWMutex = new(sync.RWMutex)

	return rule
}

func (rule *RuleSnakeCase) GetName() string {
	rule.Lock()
	defer rule.Unlock()

	return rule.Name
}

func (rule *RuleSnakeCase) SetParameters(params []string) error {
	return nil
}

// Validate checks if string is sneak case
// false if rune is no lowercase letter, digit or _
func (rule *RuleSnakeCase) Validate(value string) (bool, error) {
	for _, c := range value {
		if c == 95 || unicode.IsDigit(c) { // 95 => _
			continue
		}

		if !unicode.IsLetter(c) {
			return false, nil
		}

		if !unicode.IsLower(c) {
			return false, nil
		}
	}

	return true, nil
}

func (rule *RuleSnakeCase) GetErrorMessage() string {
	return rule.GetName()
}
