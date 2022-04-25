package tests

import (
	"fmt"
	"groupsCrawl/watcher/services/models"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := models.PriorityQueue{}
	queue.Push(&models.Item{
		Value:    "hello",
		Priority: 0,
	})
	queue.Push(&models.Item{
		Value:    "hello2",
		Priority: 0,
	})
	queue.Push(&models.Item{
		Value:    "hello3",
		Priority: 10,
	})
	fmt.Println(queue)
	value := queue.Pop()
	fmt.Println(value)
}
