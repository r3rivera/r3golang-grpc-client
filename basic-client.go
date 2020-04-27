package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	mock "r3golang-grpc-client/readers"

	pb "github.com/r3rivera/r3app-protobuffer-repo/basicpb"
	pb2 "github.com/r3rivera/r3app-protobuffer-repo/messagepb"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("GRPC Client application is starting!")

	port := fmt.Sprintf("%s:%s", os.Getenv("CLIENT_SERVER_HOST"), os.Getenv("CLIENT_SERVER_PORT"))
	isTLS, _ := strconv.ParseBool(os.Getenv("SECURE_MODE"))

	opts := []grpc.DialOption{}
	if isTLS {
		log.Println("Connecting to server in secure mode!")
		credFile := os.Getenv("CLIENT_CERT_PATH")
		creds, err := credentials.NewClientTLSFromFile(credFile, "")

		if err != nil {
			log.Fatal("Unable to load the certificate!", err)
			panic(err)
		}
		//Append the credentials
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	conn, err := grpc.Dial(port, opts...)

	//Creates a connection to the server (r3golang-grpc) with Insecure connection.
	//By default, grpc uses secure connection but for this exercise, we explicitly specified unsecure connection
	//conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

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
	//notifyClient := pb.NewNotificationMessageServiceClient(conn)
	//doServerStream(notifyClient)

	//Client-Stream client
	//dataUpload := pb.NewDataUploadMessageServiceClient(conn)
	//doClientStream(dataUpload)

	//Bi-Directional Stream
	//chatClient := pb.NewChatSupportMessageServiceClient(conn)
	//doBiDirectionalStream(chatClient)

	//Unary Message
	msgClient := pb2.NewMessagePostServiceClient(conn)
	doContactUs(msgClient)
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

func doClientStream(client pb.DataUploadMessageServiceClient) {
	log.Println("Performing a Client Stream API Call...")
	mock.StoreMockCSV(client)
}

func doBiDirectionalStream(client pb.ChatSupportMessageServiceClient) {
	log.Println("Performing a Bi-Directional Stream API Call...")

	bidiStream, err := client.ChatSupportMessage(context.Background())
	if err != nil {
		log.Fatalf("Error with the bi-direction streaming. Error is %v", err)
		panic(err)
	}

	waitc := make(chan struct{})
	//Send a bunch of messsage
	go func() {

		for _, chatRqst := range mockChatSender() {
			log.Printf("Sending chat message of %v \n", chatRqst)
			bidiStream.Send(chatRqst)
			time.Sleep(3 * time.Second)
		}

		err := bidiStream.CloseSend()
		if err != nil {
			log.Println("Error with CloseSend() of Bi-directional stream. ", err)
		}
	}()

	//Receive a response
	go func() {

		for {
			res, err := bidiStream.Recv()
			if err == io.EOF {
				log.Println("Bi-direction stream is EOF!", err)
				break
			}
			if err != nil {
				log.Println("Error receiving bi-directional response. ", err)
				break
			}
			fmt.Println("Got response of ", res)
		}
		close(waitc)

	}()

	//Block until everything is done
	<-waitc

}

func mockChatSender() []*pb.ChatSupportMessageRequest {

	chats := []*pb.ChatSupportMessageRequest{
		{
			ChatId: "123",
			ChatMessage: &pb.ChatSupportMessage{
				Sender:     "R2D2",
				Receipient: "Admin",
				Message:    "We need help",
			},
		},
		{
			ChatId: "456",
			ChatMessage: &pb.ChatSupportMessage{
				Sender:     "James",
				Receipient: "Admin",
				Message:    "Need some ticket",
			},
		},
		{
			ChatId: "789",
			ChatMessage: &pb.ChatSupportMessage{
				Sender:     "John",
				Receipient: "Admin",
				Message:    "Need IT help",
			},
		},
		{
			ChatId: "ABC",
			ChatMessage: &pb.ChatSupportMessage{
				Sender:     "John",
				Receipient: "Admin",
				Message:    "Need food",
			},
		},
	}

	return chats
}

func doContactUs(client pb2.MessagePostServiceClient) {
	log.Println("Invoking contact us GRPC...")

	payload := pb2.MessagePost{
		FirstName:   "Customer",
		LastName:    "Interest",
		Email:       "abc@test.com",
		PhoneNumber: "+18007775555",
		Message:     "What is going on here?",
	}
	rqst := pb2.MessagePostRequest{
		Post: &payload,
	}

	resp, err := client.MessagePost(context.Background(), &rqst)

	if err != nil {
		log.Println(err)
	}

	log.Println(resp)
}
