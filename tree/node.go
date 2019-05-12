package tree

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type Add struct {
	Key   int32
	Value string
}
type Remove struct {
	Key int32
}
type Find struct {
	Key int32
}

type Delete struct{}

type Traverse struct{}

type Node struct {
	MaxLeft   int32
	LeftSucc  *actor.PID
	RightSucc *actor.PID
	Data     map[int]string
	LeafSize int32
} //Actor

func (state *Node) Receive(context actor.Context) {
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

/*func main() {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &node{Amount:2,}
	})
	pid := context.Spawn(props)
	context.Send(pid, &Add{Key: 12, Value: 13})
	context.Send(pid, &Remove{Key: 14})
	console.ReadLine() // nolint:errcheck
}
*/
