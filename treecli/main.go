package main

import (
	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"vss/blatt3/blatt-3-stortz/messages"
)

type MyActor struct {
}

func (state *MyActor) Receive(context actor.Context) {
	switch message := context.Message().(type) {
	case *messages.TreeCreated:
		fmt.Printf("CLI received TreeCreated {Token: %s, ID: %d}\n", message.Token, message.TreeID)
		//context.Respond(&messages.Add{Value:"23", Key:7, Token:message.Token, TreeID:message.TreeID })
		wg.Done()
	case *messages.InvalidRequest:
		fmt.Printf("CLI received InvalidRequest {Token: %s, ID: %d}: %s\n", message.Token, message.TreeID, message.Description)
		wg.Done()
	case *messages.SuccessfulRequest:
		fmt.Printf("CLI received ValidRequest {Token: %s, ID: %d}: %s\n", message.Token, message.TreeID, message.Description)
		wg.Done()
	}
}

var (
	// nolint:gochecknoglobals
	flagBind = flag.String("bind", "localhost:8092", "Bind to address")
	// nolint:gochecknoglobals
	flagRemote = flag.String("remote", "127.0.0.1:8091", "remote host:port")

	flagToken = flag.String("token", "", "token")
	flagID = flag.Int("id", -1, "Tree-ID")

	//flagCommand = flag.String("cmd", "", "specify command")
	//flagKey = flag.Int("key", -1, "Key")
	//flagValue = flag.Int("value", -1, "Value")
	//flagLeafSize = flag.Int("leafsize", 2, "Leaf-Size")
)

var wg sync.WaitGroup

func main() {

	flag.Parse()

	remote.Start(*flagBind)
	props := actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1)
		return &MyActor{}
	})
	rootContext := actor.EmptyRootContext
	pid := rootContext.Spawn(props)

	fmt.Println("Sleeping 5 seconds...")
	time.Sleep(5 * time.Second)
	fmt.Println("Awake...")

	//this is the remote actor we want to communicate with
	fmt.Printf("Trying to connect to %s\n", *flagRemote)

	pidResp, err := remote.SpawnNamed(*flagRemote, "remote", "hello", 5*time.Second)
	if err != nil {
		panic(err)
	}
	remotePid := pidResp.Pid

	/*for i := 0; i < 10; i++ {
		rootContext.RequestWithCustomSender(remotePid, message, pid)
	}*/

	switch flag.Args()[0] {
	case "create":
		if len(flag.Args()) != 2 {
			println("invalid amount of args")
			return
		}
		leafsize := parseToInt32(1)
		rootContext.RequestWithCustomSender(remotePid, &messages.CreateTree{leafsize}, pid)
		wg.Wait()
	case "delete":
	case "add":
		println("Trying to add node")
		if len(flag.Args()) != 3 {
			println("invalid amount of args")
			return
		}
		key := parseToInt32(1)
		val := flag.Args()[2]
		rootContext.RequestWithCustomSender(remotePid, &messages.Add{TreeID:int32(*flagID), Token:*flagToken, Key:key, Value:val}, pid)
		wg.Wait()
	case "find":
	case "remove":
	case "traverse":
	case "":
		fmt.Println("No command specified!")
		wg.Done()
	}
}

func checkArgAmount(amount int){

}

func parseToInt32(pos int) int32 {
	i, err := strconv.ParseInt(flag.Args()[pos], 10, 32)
	if err != nil {
		panic(err)
	}
	result := int32(i)
	return result
}
