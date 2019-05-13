package tree

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-stortz/messages"
	"sort"
)

type Add struct {
	Key   int
	Value string
}
type Remove struct {
	Key int
}
type Find struct {
	Requester *actor.PID
	Key       int
}

type Found struct {
	Key   int
	Value string
}

type Delete struct{}

type Traverse struct {
	Requester      *actor.PID
	RemainingNodes []*actor.PID
	Data           map[int]string
}

type Node struct {
	MaxLeft   int
	LeftSucc  *actor.PID
	RightSucc *actor.PID
	Data      map[int]string
	LeafSize  int
} //Actor

func (state *Node) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Add:
		fmt.Printf("Node: Add %v\n", msg.Key)

		if state.LeftSucc == nil && state.RightSucc == nil && len(state.Data) < state.LeafSize { //Leaf still has room
			//Create Map if nil
			if state.Data == nil {
				state.Data = make(map[int]string)
			}

			state.Data[msg.Key] = msg.Value
			fmt.Printf("Added pair {k: %d, v: %s} to Node with PID: %s \n", msg.Key, msg.Value, context.Self().Address)
		} else if state.LeftSucc == nil && state.RightSucc == nil && len(state.Data) == state.LeafSize { //Leaf is full -> Split Data
			//Create two new Nodes
			props := actor.PropsFromProducer(func() actor.Actor {
				return &Node{LeafSize: state.LeafSize}
			})
			state.LeftSucc = context.Spawn(props)
			state.RightSucc = context.Spawn(props)

			//Temp assign
			state.Data[msg.Key] = msg.Value

			//Sort keys to split them correctly
			var keys []int
			for k := range state.Data {
				keys = append(keys, int(k))
			}
			sort.Ints(keys)

			//Assign to the correct nodes
			mid := int(len(keys) / 2)
			fmt.Printf("Splitting Data at key: %d\n", keys[mid])
			for i := 0; i <= mid; i++ {
				fmt.Printf("send %d left\n", keys[i])
				context.Send(state.LeftSucc, &Add{Key: keys[i], Value: state.Data[keys[i]]})
			}
			state.MaxLeft = keys[mid]
			for i := mid + 1; i < len(keys); i++ {
				fmt.Printf("send %d right\n", keys[i])
				context.Send(state.RightSucc, &Add{Key: keys[i], Value: state.Data[keys[i]]})
			}

			//Clean Current Data
			state.Data = nil

		} else if state.RightSucc != nil && state.LeftSucc != nil { //Find
			if msg.Key <= state.MaxLeft {
				fmt.Printf("send %d left\n", msg.Key)
				context.Send(state.LeftSucc, &Add{Key: msg.Key, Value: msg.Value})
			} else {
				fmt.Printf("send %d right\n", msg.Key)
				context.Send(state.RightSucc, &Add{Key: msg.Key, Value: msg.Value})
			}
		}

	case *Remove:
		fmt.Printf("Trying to remove pair with key: %d\n", msg.Key)
		if state.Data != nil { //Leaf
			if _, ok := state.Data[msg.Key]; ok {
				delete(state.Data, msg.Key)
				fmt.Printf("Key found in Tree -> Removing: %d\n", msg.Key)
			} else {
				fmt.Printf("Key not found in Tree!: %d\n", msg.Key)
			}
		} else { //Inner Node
			if msg.Key <= state.MaxLeft {
				context.Send(state.LeftSucc, &Remove{Key: msg.Key})
			} else {
				context.Send(state.RightSucc, &Remove{Key: msg.Key})
			}
		}
	case *Find:
		fmt.Printf("Trying to find value for key: %d \n", msg.Key)
		if state.Data != nil { //Leaf
			if val, ok := state.Data[msg.Key]; ok {

				context.Send(msg.Requester, &messages.Found{Key: int32(msg.Key), Value: val})
			} else {
				fmt.Printf("Key not found in Tree!: %d\n", msg.Key)
				context.Send(msg.Requester, &messages.Found{Key: -1, Value: "Not found"})
			}
		} else if state.LeftSucc != nil && state.RightSucc != nil { //Inner Node
			if msg.Key <= state.MaxLeft {
				context.Send(state.LeftSucc, &Find{Key: msg.Key, Requester: msg.Requester})
			} else {
				context.Send(state.RightSucc, &Find{Key: msg.Key, Requester: msg.Requester})
			}
		} else {
			fmt.Println("Something went horribly wrong")
		}
	case *Delete:
		fmt.Printf("Stopping all actors for this tree \n")
		context.Send(state.LeftSucc, &Delete{})
		context.Send(state.RightSucc, &Delete{})
		context.Stop(context.Self())
	case *Traverse:
		fmt.Printf("Traversing tree\n")
		if msg.Data == nil {
			msg.Data = make(map[int]string)
		}

		if state.LeftSucc != nil && state.RightSucc != nil { //Node

			//Add right node to remaining nodes then call left
			msg.RemainingNodes = append(msg.RemainingNodes, state.RightSucc)
			context.Send(state.LeftSucc, &Traverse{Requester: msg.Requester, RemainingNodes: msg.RemainingNodes, Data: msg.Data})
		} else { //Leaf
			//Add Data
			for k, v := range state.Data {
				msg.Data[k] = v
			}

			if len(msg.RemainingNodes) > 0 {
				pid := msg.RemainingNodes[len(msg.RemainingNodes)-1]
				msg.RemainingNodes = msg.RemainingNodes[:len(msg.RemainingNodes)-1]
				context.Send(pid, &Traverse{Requester: msg.Requester, RemainingNodes: msg.RemainingNodes, Data: msg.Data})

			} else { //Visited all Nodes

				//Sort
				var keys []int
				for k := range msg.Data {
					keys = append(keys, int(k))
				}
				sort.Ints(keys)

				pairs := make([]*messages.Pair, 0)

				for i := range keys {
					fmt.Printf("{%d, %s}", keys[i], msg.Data[keys[i]])
					pairs = append(pairs, &messages.Pair{Key: int32(keys[i]), Value: msg.Data[keys[i]]})
				}
				context.Send(msg.Requester, &messages.TraversedAnswer{Pairs: pairs})

			}
		}
	}
}
