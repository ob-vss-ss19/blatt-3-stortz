package main

import (
	"flag"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-stortz/messages"
	"github.com/ob-vss-ss19/blatt-3-stortz/tree"
	"sync"
)

type MyActor struct{}

var nextID int32 = 0
var trees = make(map[int32]*actor.PID)
var tokenToID = make(map[string]int32)

func (state *MyActor) Receive(context actor.Context) {
	switch message := context.Message().(type) {
	case *messages.Add:
		fmt.Printf("Service received AddNode-Message\n")
		if validateTokenAndID(message.Token, message.TreeID, context) {
			pid := trees[message.TreeID]
			context.Send(pid, &tree.Add{Key: int(message.Key), Value: message.Value})
			desc := fmt.Sprintf("Service tries to create Node for Tree %d {k: %d, v: %s}\n", message.TreeID, message.Key, message.Value)
			context.Respond(&messages.SuccessfulRequest{Token: message.Token, TreeID: message.TreeID, Description: desc})
		}

	case *messages.CreateTree:
		fmt.Printf("Service received CreateTree-Message {Leaf-Size: %d}\n", message.LeafSize)
		token := generateToken()
		//Create Node
		props := actor.PropsFromProducer(func() actor.Actor {
			return &tree.Node{LeafSize: int(message.LeafSize)}
		})
		pid := context.Spawn(props)
		trees[nextID] = pid
		tokenToID[token] = nextID

		context.Respond(&messages.TreeCreated{Token: token, TreeID: nextID})
		nextID++
	case *messages.Delete:
		fmt.Printf("Service received Delete-Message for Tree %d with Token %s \n", message.TreeID, message.Token)
		if validateTokenAndID(message.Token, message.TreeID, context) {
			pid := trees[message.TreeID]
			context.Send(pid, &tree.Delete{})
			desc := fmt.Sprintf("Service tries to delete tree %d {k: %d, v: %s}\n", message.TreeID)
			context.Respond(&messages.SuccessfulRequest{Token: message.Token, TreeID: message.TreeID, Description: desc})
		}
	case *messages.Find:
		fmt.Printf("Service received Find-Message for Tree %d with Token %s {Key: %d}\n", message.TreeID, message.Token, message.Key)
		if validateTokenAndID(message.Token, message.TreeID, context) {
			pid := trees[message.TreeID]
			context.Send(pid, &tree.Find{Key: int(message.Key), Requester: context.Sender()})
			//desc := fmt.Sprintf("Service tries to find pair for Tree %d {k: %d, }\n", message.TreeID, message.Key, message.Value)
			//context.Respond(&messages.SuccessfulRequest{Token:message.Token,TreeID: message.TreeID, Description:desc})
		}
	case *messages.Remove:
		fmt.Printf("Service received Remove-Message for Tree %d with Token %s {Key: %d}\n", message.TreeID, message.Token, message.Key)
		if validateTokenAndID(message.Token, message.TreeID, context) {
			pid := trees[message.TreeID]
			context.Send(pid, &tree.Remove{Key: int(message.Key)})
			desc := fmt.Sprintf("Service tries to remove pair for Tree %d {k: %d}\n", message.TreeID, message.Key)
			context.Respond(&messages.SuccessfulRequest{Token: message.Token, TreeID: message.TreeID, Description: desc})
		}
	case *messages.Traverse:
		fmt.Printf("Service received Traverse-Message for Tree %d with Token %s\n", message.TreeID, message.Token)
		if validateTokenAndID(message.Token, message.TreeID, context) {
			pid := trees[message.TreeID]
			context.Send(pid, &tree.Traverse{Requester: context.Sender()})
			desc := fmt.Sprintf("Service tries to traverse Tree %d \n", message.TreeID)
			context.Respond(&messages.SuccessfulRequest{Token: message.Token, TreeID: message.TreeID, Description: desc})
		}

	default: // just for linter
	}
}

func validateTokenAndID(token string, id int32, context actor.Context) bool {
	if tokenToID[token] == id { //Valid Request
		//Get Root
		pid := trees[id]
		if pid == nil {
			desc := fmt.Sprintf("Service found no tree for the given ID {Token: %s, ID %d}\n", token, id)
			fmt.Println(desc)
			context.Respond(&messages.InvalidRequest{Token: token, TreeID: id, Description: desc})
			return false
		} else {
			return true
		}
	} else { //Invalid Request; Token and ID mismatch
		desc := fmt.Sprintf("Service received mismatched Token and ID {Token: %s, ID %d}\n", token, id)
		fmt.Println(desc)
		context.Respond(&messages.InvalidRequest{Token: token, TreeID: id, Description: desc})
		return false
	}
}

func ServiceActor() actor.Actor {
	fmt.Printf("Service-Actor is up and running\n")
	return &MyActor{}
}

// nolint:gochecknoglobals
var flagBind = flag.String("bind", "localhost:8091", "Bind to address")

func generateToken() string {
	return "a"
	/*	b := make([]byte, 4)
		rand.Read(b)
		return fmt.Sprintf("%x", b)*/
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	flag.Parse()
	remote.Start(*flagBind)
	remote.Register("ServiceActor", actor.PropsFromProducer(ServiceActor))
}
