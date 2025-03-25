package tools

import (
	context "context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func NewClient() (*grpc.ClientConn, error) {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func GetWebResourceContent(ctx context.Context, url string) (*ExtractResponse, error) {

	conn, err := NewClient()
	if err != nil {
		log.Printf("error connecting to grpc server: %v", err)
		return nil, err
	}
	defer conn.Close()

	c := NewExtractServiceClient(conn)

	response, err := c.Extract(ctx, &ExtractRequest{Url: url})
	if err != nil {
		log.Printf("error sending request grpc: %v", err)
		return nil, err
	}

	return response, nil
}
