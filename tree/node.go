package main

import (
	"fmt"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type Add struct {
	Key   int
	Value int
}
type Remove struct {
	Key int
}
type Find struct {
	Key int
}

type Delete struct{}

type Traverse struct{}

type node struct {
	MaxLeft   int
	LeftSucc  *actor.PID
	RightSucc *actor.PID
} //Actor

func (state *node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Add:
		fmt.Printf("Add %v\n", msg.Key)
	case *Remove:
		fmt.Printf("Remove %v\n", msg.Key)
	case *Find:
		fmt.Printf("Find %v\n", 14)
	case *Delete:
		fmt.Printf("Delete %v\n", 13)
	case *Traverse:
		fmt.Printf("Traverse %v\n", 23)

	}
}

func main() {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &node{}
	})
	pid := context.Spawn(props)
	context.Send(pid, &Add{Key: 12, Value: 13})
	context.Send(pid, &Remove{Key: 14})
	console.ReadLine() // nolint:errcheck
}
