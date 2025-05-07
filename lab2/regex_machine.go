package main

import "errors"

type RegexMachine struct {
	pattern  string
	automata *Automata
}

func (r *RegexMachine) Compile() error {
	r.pattern = "not implemented"
	return errors.New("not implemented")
}

func (r *RegexMachine) Match(matchee string) (bool, error) {
	if r.automata == nil {
		err := r.Compile()
		if err != nil {
			return false, err
		}
	}
	return false, errors.New("not implemented")
}

// Recovers pattern from
func (r *RegexMachine) RecoverPattern() (string, error) {
	if r.automata == nil {
		err := r.Compile()
		if err != nil {
			return "", err
		}
	}
	return "", errors.New("not implemented")
}

// Returns intersection of two given RegexMachines
func (r *RegexMachine) Intersect(other *RegexMachine) (*RegexMachine, error) {
	if r.automata == nil {
		err := r.Compile()
		if err != nil {
			return nil, err
		}
	}
	if other.automata == nil {
		err := other.Compile()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// Subtracts other RegexMachine from this one and returns the result
func (r *RegexMachine) Subtract(other *RegexMachine) (*RegexMachine, error) {
	if r.automata == nil {
		err := r.Compile()
		if err != nil {
			return nil, err
		}
	}
	if other.automata == nil {
		err := other.Compile()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
