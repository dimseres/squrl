package models

type Node struct {
	Value string
	Next  *Node
}

type Container struct {
	Priority uint8
	Root     *Node
	Next     *Container
	Prev     *Container
}

type Queue struct {
	Head *Container
	Tail *Container
}

func NewQueue() Queue {
	return Queue{
		Head: &Container{},
		Tail: &Container{},
	}
}

func (queue *Queue) makeContainer(prior uint8, value string) Container {
	return Container{
		Priority: prior,
		Root: &Node{
			Value: value,
		},
	}
}

func (queue *Queue) Insert(prior uint8, value string) {
	container := queue.makeContainer(prior, value)

	if queue.Head.Priority < prior {
		container.Next = queue.Head
		queue.Head = &container
		queue.Tail = container.Next.Next
		return
	}

	if queue.Tail.Priority > prior {
		container.Prev = queue.Tail
		queue.Tail = &container
	}
}
