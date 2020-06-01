package cli

var commands = map[string]string{
	"help":   "",
	"h":      "",
	"status": "",
	"st":     "",
	"reload": "",
	"start":  "",
	"run":    "",
	"stop":   "",
	"exit":   "",
	"quit":   "",
	"fg":     "",
	"bg":     "",
}

// CommandNames ...
func CommandNames() []string {
	var names []string
	for key := range commands {
		names = append(names, key)
	}
	return names
}
