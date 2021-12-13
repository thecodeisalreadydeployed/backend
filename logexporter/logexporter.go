package logexporter

import (
	"bufio"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	queue := NewQueue()
	for {
		scanner.Scan()
		text := scanner.Text()
		queue.Enqueue(text)
	}
}
