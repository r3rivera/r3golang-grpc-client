package main

import (
	"context"
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
	client := pb.NewHealthCheckStatusServiceClient(conn)

	doUnaryCall(client)

}

//Unary API Call
func doUnaryCall(client pb.HealthCheckStatusServiceClient) {
	log.Println("Performing a UNARY API Call...")
	//Creating a payload request
	rqst := pb.HealthCheckStatusRequest{
		AppName: "R3APPClient",
	}

	response, err := client.HealthCheckStatus(context.Background(), &rqst)
	if err != nil {
		log.Fatalf("Error getting the health response ! %v", err)
		panic(err)
	}

	log.Printf("Client is connected :: %v", response)
}
