package tests

import (
	"fmt"
	"groupsCrawl/watcher/services/models"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := models.NewQueue()
	queue.Insert(1, "prior 1")
	queue.Insert(2, "prior 2")
	queue.Insert(3, "prior 3")
	fmt.Println(queue)
}
