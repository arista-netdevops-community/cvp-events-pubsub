package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"time"

	kafkastream "github.com/arista-netdevops-community/cvp-events-pubsub/pkg"
	"github.com/aristanetworks/cloudvision-go/api/arista/event.v1"
	cvgrpc "github.com/aristanetworks/cloudvision-go/grpc"
	"github.com/aristanetworks/glog"
)

const (
	timeout = 25 * time.Second
)

func main() {
	// Flags
	flag.StringVar(&server, "server", "", "CloudVision IP Address")
	yamlfile := flag.String("yamlfile", "", "CloudVision server to connect to")
	flag.BoolVar(&ssl, "ssl", true, "Download certificate for server")
	flag.BoolVar(&verify, "verify", false, "Verify certificate for server")
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")
	authConfig := cvgrpc.AuthFlag()
	flag.Parse()
	// Initialize the token to create the cvp.crt and token.txt to put into the config file.
	ReturnToken()
	// Get to the YAML files from the config/data.yaml if necessary
	var k kafkastream.Config
	k.GetConf(*yamlfile)
	// Create a blank ctx
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Start the grpc connection to cvp
	conn, err := cvgrpc.DialWithAuth(ctx, server, authConfig)
	//conn, err := cvgrpc.DialWithToken(ctx, server, InitToken)
	if err != nil {
		glog.Fatalf("failed to dial server: %s", err)
	}
	eventClient := event.NewEventServiceClient(conn)
	esr := event.EventStreamRequest{}
	stream, err := eventClient.Subscribe(context.Background(), &esr)
	if err != nil {
		glog.Fatalf("This broke because %s", err)
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			glog.Fatalf("failed to receive on stream: %s", err)
		}
		//This prints out any updated new events after this has ran.
		// Types events UPDATED,INITIAL,DELETED
		if resp.Type.String() == "UPDATED" {
			//Print to output of the response value + description of this.
			fmt.Println(resp.Value.Title.Value, "\t", resp.Value.Description.Value)
			//Pass into the StreamToKafka function.
			kafkastream.StreamToKafka(k.Kafka_broker, k.Kafka_topic, resp.Value.Title.Value, resp.Value.Description.Value)
		}
	}
}
