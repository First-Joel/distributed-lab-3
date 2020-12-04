package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"pairbroker/stubs"
)


type Factory struct {}

var mulch = make(chan int,10)

func getOutboundIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.String()
	return localAddr
}

func makedivision(ch chan int, broker *rpc.Client){
	for {
		x := <-ch
		y := <-ch
		newpair := stubs.Pair{X: x, Y: y}
		towork := stubs.PublishRequest{Topic: "divide", Pair: newpair}
		status := new(stubs.StatusReport)
		err := broker.Call(stubs.Publish, towork, status)
		if err != nil {
			fmt.Println("RPC client returned error:")
			fmt.Println(err)
			fmt.Println("Dropping division.")
		}
	}
}

//TODO: Define a Multiply function to be accessed via RPC. 
//Check the previous weeks' examples to figure out how to do this.
func (f *Factory) Multiply (intoken stubs.Pair, outtoken *stubs.JobReport)(err error){
	x := intoken.X
	fmt.Println("hi")
	y := intoken.Y
	outtoken.Result = x*y
	mulch <- outtoken.Result
	fmt.Println(*outtoken)
	return
}

func (f *Factory) Divide (req stubs.Pair, resp *stubs.JobReport) (err error){
	resp.Result = req.X /req.Y
	fmt.Println(resp.Result)
	return
}


func main(){
	pAddr := flag.String("port","8050","Port to listen on")
	brokerAddr := flag.String("broker","127.0.0.1:8030", "Address of broker instance")
	flag.Parse()

	broker, _ := rpc.Dial("tcp", *brokerAddr)
	listener, _ := net.Listen("tcp",":"+*pAddr)
	rpc.Register(&Factory{})
	fmt.Println(*pAddr)

	status := new(stubs.StatusReport)


	subscription1 := stubs.Subscription{Topic:"multiply",FactoryAddress:"127.0.0.1:"+*pAddr,Callback:"Factory.Multiply"}
	subscription2 := stubs.Subscription{Topic:"divide",FactoryAddress:"127.0.0.1:"+*pAddr,Callback:"Factory.Divide"}
	broker.Call(stubs.Subscribe,subscription1,status)
	broker.Call(stubs.Subscribe,subscription2,status)
	defer listener.Close()
	go makedivision(mulch, broker)
	rpc.Accept(listener)




	//TODO: You'll need to set up the RPC server, and subscribe to the running broker instance.
}
