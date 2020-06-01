package tty

const historyLimit = 3

// Proposer (autocompleter) takes current input and returns all possible completions
type Proposer func(input string, jobNames []string) []string

type autocomplete struct {
	// autocomplete func
	proposer Proposer
	// index in suggestions arr
	proposerPos int
	// Currently available suggestions
	suggestions []string
	// All job names, getting these from server
	jobNames []string
}

type hist struct {
	history      []string
	historyCount int
	historyPos   int
}

// SetJobNames ...
func (s *State) SetJobNames(names []string) {
	s.jobNames = names
	// fmt.Println(s.jobNames)
}

// SetProposer sets autocomplete function
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

// func (s *State) historySearch(prefix string) []string {
// 	var result []string
// 	for _, item := range s.history {
// 		if strings.HasPrefix(item, prefix) {
// 			result = append(result, item)
// 		}
// 	}
// 	return result
// }
