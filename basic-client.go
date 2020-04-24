package main

import (
	"context"
	"io"
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
	//healthClient := pb.NewHealthCheckStatusServiceClient(conn)
	//doUnaryCall(healthClient)

	//calcClient := pb.NewCalculatorServiceClient(conn)
	//doUnaryCalc(calcClient)

	//Server-streaming client
	notifyClient := pb.NewNotificationMessageServiceClient(conn)
	doServerStream(notifyClient)
}

//Unary API Call
func doUnaryCall(client pb.HealthCheckStatusServiceClient) {
	log.Println("Performing a UNARY API Call...")
	//Creating a payload request
	rqst := pb.HealthCheckStatusRequest{
		AppName: "R3APPClient",
	}

	resp, err := client.HealthCheckStatus(context.Background(), &rqst)
	if err != nil {
		log.Fatalf("Error getting the health response ! %v", err)
		panic(err)
	}

	log.Printf("Client is connected :: %v", resp)
}

func doUnaryCalc(client pb.CalculatorServiceClient) {
	log.Println("Performing a UNARY API Calculator Call...")

	calc := pb.Calculator{
		NumOne: 20,
		NumTwo: 30,
	}
	rqst := pb.CalculatorRequest{
		Payload: &calc,
	}

	resp, err := client.Calculator(context.Background(), &rqst)
	if err != nil {
		log.Fatalf("Error getting the calculator response ! %v", err)
		panic(err)
	}
	log.Printf("Client is connected :: %v", resp)
}

func doServerStream(client pb.NotificationMessageServiceClient) {
	log.Println("Performing a Server Stream API Call...")

	rqst := pb.NotificationMessageRequest{
		Requester: "Tester",
	}
	respStream, err := client.NotificationMessage(context.Background(), &rqst)
	if err != nil {
		log.Fatalf("Error sending the notification message stream ! %v", err)
		panic(err)
	}

	for {
		msg, err := respStream.Recv()

		//Check if we still have message from stream
		if err == io.EOF {
			log.Println("Streaming message is complete. Exiting...")
			break
		}

		if err != nil {
			log.Fatalf("Error getting the notification message stream ! %v", err)
			panic(err)
		}
		log.Println("Response is ", msg)
	}

}
