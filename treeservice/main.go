package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"sync"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"vss/blatt3/blatt-3-stortz/messages"
	"vss/blatt3/blatt-3-stortz/tree"
)

type MyActor struct{}

var nextID int32 = 0
var trees = make(map[int32]*actor.PID)
var tokenToID = make(map[string]int32)

func (state *MyActor) Receive(context actor.Context) {
	switch message := context.Message().(type) {
	case *messages.HelloWorld:
		//fmt.Printf("Service not responding bc syntax")
		context.Respond(&messages.HelloWorld{})
	case *messages.Add:
		fmt.Printf("Service received AddNode-Message\n")
		if tokenToID[message.Token] == message.TreeID { //Valid Request
			//Get Root
			pid := trees[message.TreeID]
			if pid == nil{
				desc := fmt.Sprintf("Service found no tree for the given ID {Token: %s, ID %d}\n", message.Token, message.TreeID)
				fmt.Println(desc)
				context.Respond(&messages.InvalidRequest{Token:message.Token,TreeID: message.TreeID, Description:desc})
			}
			context.Send(pid, &tree.Add{Key:message.Key, Value:message.Value})
			desc := fmt.Sprintf("Service tries to create Node for Tree %d {k: %d, v: %s}\n", message.TreeID, message.Key, message.Value)
			context.Respond(&messages.SuccessfulRequest{Token:message.Token,TreeID: message.TreeID, Description:desc})


		}else { //Invalid Request; Token and ID mismatch
			desc := fmt.Sprintf("Service received mismatched Token and ID {Token: %s, ID %d}\n", message.Token, message.TreeID)
			fmt.Println(desc)
			context.Respond(&messages.InvalidRequest{Token:message.Token,TreeID: message.TreeID, Description:desc})
		}

	case *messages.CreateTree:
		fmt.Printf("Service received CreateTree-Message {Leaf-Size: %d}\n", message.LeafSize)
		token := generateToken()
		context.Respond(&messages.TreeCreated{Token:token,TreeID: nextID})

		//Create Node
		props := actor.PropsFromProducer(func() actor.Actor {
			return &tree.Node{LeafSize:message.LeafSize}
		})
		pid := context.Spawn(props)
		trees[nextID] = pid
		tokenToID[token] = nextID

		nextID ++


	default: // just for linter
	}
}

func NewMyActor() actor.Actor {
	fmt.Printf("Hello-Actor is up and running\n")
	return &MyActor{}
}

// nolint:gochecknoglobals
var flagBind = flag.String("bind", "localhost:8091", "Bind to address")

func main() {
	var wg sync.WaitGroup
	wg.Add(1)


	defer wg.Wait()

	flag.Parse()
	remote.Start(*flagBind)

	remote.Register("hello", actor.PropsFromProducer(NewMyActor))
}

func generateToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
