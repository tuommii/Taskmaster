package main

func main() {
	client := create()
	// _ = logger.Get()
	// log.Println("CLIENT TEST")
	go client.listenSignals()
	go client.readInput()
	client.quit()
}
