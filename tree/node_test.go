package tree

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-stortz/messages"
	"sync"
	"testing"
	"time"
)

type TestActor struct {
	t     *testing.T
	wg    *sync.WaitGroup
	pairs []KVPair
}

type KVPair struct {
	key   int
	value string
}

func (state *TestActor) Receive(context actor.Context) {
	switch message := context.Message().(type) {
	case *messages.TreeCreated:
		fmt.Printf("CLI received TreeCreated {Token: %s, ID: %d}\n", message.Token, message.TreeID)
		state.wg.Done()
	case *messages.InvalidRequest:
		fmt.Printf("CLI received InvalidRequest {Token: %s, ID: %d}: %s\n", message.Token, message.TreeID, message.Description)
		state.wg.Done()
	case *messages.SuccessfulRequest:
		fmt.Printf("CLI received ValidRequest {Token: %s, ID: %d}: %s\n", message.Token, message.TreeID, message.Description)
		state.wg.Done()
	case *messages.Found:
		fmt.Printf("Expected{%d, %s}; Got{%d, %s}\n", state.pairs[0].key, state.pairs[0].value, message.Key, message.Value)
		if state.pairs[0].key != int(message.Key) || state.pairs[0].value != message.Value {
			fmt.Printf("Mismatch!\n")
			state.t.Error()
		}
		state.wg.Done()
	case *messages.TraversedAnswer:
		for i, k := range message.GetPairs() {
			fmt.Printf("Expected{%d, %s}; Got{%d, %s}\n", state.pairs[i].key, state.pairs[i].value, k.Key, k.Value)
			if state.pairs[i].key != int(k.Key) || state.pairs[i].value != k.Value {
				fmt.Printf("Mismatch!\n")
				state.t.Error()
			}
		}
		state.wg.Done()
	}

}

//Create Tree and add 2 -> 3 -> 1
func TestAdd(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Node{LeafSize: 2}
	})
	root := context.Spawn(props)
	var wg sync.WaitGroup

	context.Send(root, &Add{Key: 2, Value: "zwei"})
	context.Send(root, &Add{Key: 3, Value: "drei"})
	context.Send(root, &Add{Key: 1, Value: "eins"})

	time.Sleep(1 * time.Second)
	pairs := make([]KVPair, 0)
	pairs = append(pairs, KVPair{1, "eins"})
	pairs = append(pairs, KVPair{2, "zwei"})
	pairs = append(pairs, KVPair{3, "drei"})

	props = actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1)
		return &TestActor{t, &wg, pairs}
	})
	evalActor := context.Spawn(props)
	context.Send(root, &Traverse{Requester: evalActor})
	time.Sleep(1 * time.Second)
	wg.Wait()
}

//Create Tree
func TestTraverse(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Node{LeafSize: 2}
	})
	root := context.Spawn(props)
	var wg sync.WaitGroup

	context.Send(root, &Add{Key: 2, Value: "zwei"})
	context.Send(root, &Add{Key: 3, Value: "drei"})
	context.Send(root, &Add{Key: 1, Value: "eins"})
	context.Send(root, &Add{Key: 4, Value: "4"})
	context.Send(root, &Add{Key: 6, Value: "6"})
	context.Send(root, &Add{Key: 9, Value: "9"})
	context.Send(root, &Add{Key: 13, Value: "13"})
	context.Send(root, &Add{Key: 5, Value: "5"})
	context.Send(root, &Add{Key: 7, Value: "7"})

	time.Sleep(1 * time.Second)
	pairs := make([]KVPair, 0)
	pairs = append(pairs, KVPair{1, "eins"})
	pairs = append(pairs, KVPair{2, "zwei"})
	pairs = append(pairs, KVPair{3, "drei"})
	pairs = append(pairs, KVPair{4, "4"})
	pairs = append(pairs, KVPair{5, "5"})
	pairs = append(pairs, KVPair{6, "6"})
	pairs = append(pairs, KVPair{7, "7"})
	pairs = append(pairs, KVPair{9, "9"})
	pairs = append(pairs, KVPair{13, "13"})

	props = actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1)
		return &TestActor{t, &wg, pairs}
	})
	evalActor := context.Spawn(props)
	context.Send(root, &Traverse{Requester: evalActor})
	time.Sleep(1 * time.Second)
	wg.Wait()
}

func TestFind(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Node{LeafSize: 2}
	})
	root := context.Spawn(props)
	var wg sync.WaitGroup

	context.Send(root, &Add{Key: 2, Value: "zwei"})
	context.Send(root, &Add{Key: 3, Value: "drei"})
	context.Send(root, &Add{Key: 1, Value: "eins"})

	time.Sleep(1 * time.Second)
	pairs := make([]KVPair, 0)
	pairs = append(pairs, KVPair{2, "zwei"})

	props = actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1)
		return &TestActor{t, &wg, pairs}
	})
	evalActor := context.Spawn(props)
	context.Send(root, &Find{Requester: evalActor, Key: 2})
	time.Sleep(1 * time.Second)
	wg.Wait()
}

func TestRemove(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &Node{LeafSize: 2}
	})
	root := context.Spawn(props)
	var wg sync.WaitGroup

	context.Send(root, &Add{Key: 2, Value: "zwei"})
	context.Send(root, &Add{Key: 3, Value: "drei"})
	context.Send(root, &Add{Key: 1, Value: "eins"})
	context.Send(root, &Remove{Key: 3})

	time.Sleep(1 * time.Second)
	pairs := make([]KVPair, 0)
	pairs = append(pairs, KVPair{1, "eins"})
	pairs = append(pairs, KVPair{2, "zwei"})

	props = actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1)
		return &TestActor{t, &wg, pairs}
	})
	evalActor := context.Spawn(props)
	context.Send(root, &Traverse{Requester: evalActor})
	time.Sleep(1 * time.Second)
	wg.Wait()
}
