package main

import (
	"log"

	pb "github.com/r3rivera/r3app-protobuffer-repo/basic-test"

	"google.golang.org/grpc"
)

func main() {

	log.Println("GRPC Client application is running!")

	//Creates a connection to the server (r3golang-grpc) with Insecure connection.
	//By default, grpc uses secure connection but for this exercise, we explicitly specified unsecure connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalln("Error establishing a connection!", err)
		panic(err)
	}

	defer conn.Close()

	//Creating a client from the proto buffer
	client := pb.NewGreeterServiceClient(conn)

	log.Printf("Client is connected :: %f", client)

}
