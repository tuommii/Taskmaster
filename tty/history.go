package tty

import "strings"

const historyLimit = 3

// HistoryAdd ...
func (s *State) HistoryAdd(item string) {
	if item == "" || (s.HistoryCount > 0 && item == s.History[s.HistoryCount-1]) {
		return
	}
	s.History = append(s.History, item)
	s.HistoryCount++
	if s.HistoryCount > historyLimit {
		s.History = s.History[1:]
		s.HistoryCount--
	}
	if s.HistoryCount > 0 {
		s.HistoryPos = s.HistoryCount - 1
	}
}

// HistoryPop ...
func (s *State) HistoryPop() string {
	if s.HistoryCount == 0 {
		return ""
	}
}

// HistorySearch ...
func (s *State) HistorySearch(prefix string) []string {
	var result []string
	for _, item := range s.History {
		if strings.HasPrefix(item, prefix) {
			result = append(result, item)
		}
	}
	return result
}
