package job

// Tasks ...
type Tasks map[string]*Process

// Status ...
func (tasks Tasks) Status() string {
	var res string
	for _, task := range tasks {
		res += task.Name + ",  " + task.Status + "\n"
	}
	return res
}

// Start ...
func (tasks Tasks) Start(name string) error {
	return tasks[name].Launch()
}

// Stop ...
func (tasks Tasks) Stop(name string) error {
	return tasks[name].Kill()
}
