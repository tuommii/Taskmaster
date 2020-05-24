package tty

import "strings"

const historyLimit = 3

// Proposer takes current input and returns all possible completions
type Proposer func(input string) []string

// SetProposer ...
func (s *State) SetProposer(f Proposer) {
	s.proposer = f
}

func (s *State) historyAdd(item string) {
	if item == "" || (s.historyCount > 0 && item == s.history[s.historyCount-1]) {
		return
	}
	s.history = append(s.history, item)
	s.historyCount++
	if s.historyCount > historyLimit {
		s.history = s.history[1:]
		s.historyCount--
	}
	if s.historyCount > 0 {
		s.historyPos = s.historyCount - 1
	}
}

func (s *State) historySearch(prefix string) []string {
	var result []string
	for _, item := range s.history {
		if strings.HasPrefix(item, prefix) {
			result = append(result, item)
		}
	}
	return result
}
