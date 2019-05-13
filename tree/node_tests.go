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
		fmt.Printf("CLI received answer: {key: %d, value: %s}\n", message.Key, message.Value)
		state.wg.Done()
	case *messages.TraversedAnswer:
		fmt.Println("CLI received traverse answer:")
		for i, k := range message.GetPairs() {
			fmt.Printf("Expected{%d, %s}; Got{%d, %s}", state.pairs[i].key, state.pairs[i].value, k.Key, k.Value)
			if state.pairs[i].key != int(k.Key) || state.pairs[i].value != k.Value {
				fmt.Printf("Mismatch!")
				state.t.Error()
			}
		}
	}
	state.wg.Done()
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
