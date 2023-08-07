package main

func main() {
	server := NewAPIServer(":4444")

	server.Run()
}
