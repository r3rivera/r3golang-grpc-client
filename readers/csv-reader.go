package readers

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/r3rivera/r3app-protobuffer-repo/basicpb"
)

//StoreMockCSV Reads the mock data and send it to a channel
func StoreMockCSV(client pb.DataUploadMessageServiceClient) {

	stream, err := client.DataUploadMessage(context.Background())

	if err != nil {
		log.Fatalf("Error connectin to stream! %v", err)

	} else {

		csvFile, _ := os.Open("./mock_data/MOCK_DATA.csv")
		reader := csv.NewReader(bufio.NewReader(csvFile))

		var counter int = 0
		for {
			line, err := reader.Read()
			if err == io.EOF {

				log.Println("Done reading the CSV file...")
				break
			}

			if err != nil {
				log.Fatalln("Error reading the file...", err)
				panic(err)
			}

			if counter != 0 {
				val := fmt.Sprintf("{%s|%s|%s|%s|%s|%s", line[0], line[1], line[2], line[3], line[4], line[5])
				log.Println(val)

				payload := pb.DataUploadMessage{
					Message: val,
				}

				stream.Send(&pb.DataUploadMessageRequest{
					Payload: &payload,
				})

				time.Sleep(2 * time.Millisecond)
			}
			counter++

		}
		resp, err := stream.CloseAndRecv() // Need to have this otherwise throws an exception.
		if err != nil {
			panic(err)
		}

		if resp.StatusCode == "200" {
			log.Println("Done sending the file...")
		}
	}

}

//SpitMockCSV extract the value from channel and return it
func SpitMockCSV(c <-chan string) string {
	value := <-c
	log.Println("Value from channel ::[", value, "]")
	return value
}
