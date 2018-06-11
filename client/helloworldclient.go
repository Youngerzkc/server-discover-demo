package main

import (
	"flag"
	"fmt"
	"time"
	grpclb "service-discover/etcdv3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "service-discover/protobuf"
	"strconv"
)

var (
	serv   = flag.String("service", "hello_service", "service name")
	reg    = flag.String("reg", "http://127.0.0.1:2379", "register etcd address")
	// target = flag.String("target", "127.0.0.1:50001", "usage....")
)

func main() {
	flag.Parse()
	r := grpclb.NewResolver(*serv)
	b := grpc.RoundRobin(r)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(ctx, *reg, grpc.WithInsecure(), grpc.WithBalancer(b))
	// conn, err := grpc.DialContext(context.Background(), *target, grpc.WithInsecure())
	// conn, err := grpc.Dial(*target, grpc.WithInsecure())
	if err != nil {
		fmt.Println("建立连接失败！")
		panic(err)
	}
	// client := pb.NewGreeterClient(conn)
	// request := new(pb.HelloRequest)
	// request.Name = "world"
	// resp, err := client.SayHello(context.Background(), request)
	// if err == nil {
	// 	fmt.Println(resp.Message)
	// }

	ticker := time.NewTicker(1 * time.Second)
	for t := range ticker.C {
	    // fmt.Println("helloe ")
	    client := pb.NewGreeterClient(conn)
	    resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "world " + strconv.Itoa(t.Second())})
	    if err == nil {
	        fmt.Printf("%v: Reply is %s\n", t, resp.Message)
	    }
	}
}
