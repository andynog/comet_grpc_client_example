package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/cometbft/cometbft/libs/os"
	blocksvc "github.com/cometbft/cometbft/proto/tendermint/services/block/v1"
	client "github.com/cometbft/cometbft/rpc/grpc/client"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30*time.Second))
	defer cancel()
	grpcURI := "0.0.0.0:8080"
	conn, err := client.New(ctx, grpcURI, client.WithVersionServiceEnabled(true), client.WithInsecure())
	if err != nil {
		fmt.Printf("error new client: %v\n", err)
	}

	// Get Version
	version, err := conn.GetVersion(ctx)
	if err != nil {
		fmt.Printf("error getting block: %v\n", err)
		os.Exit("aborting...")
	}
	fmt.Printf("VERSION SERVICE => P2P: %d | Block: %d | ABCI: %s | Node: %s\n", version.P2P, version.Block, version.ABCI, version.Node)

	//// Get Block
	block, err := conn.GetBlockByHeight(ctx, 10)
	if err != nil {
		fmt.Printf("error getting block: %v\n", err)
		st := status.Convert(err)
		for _, detail := range st.Details() {
			switch t := detail.(type) {
			case *errdetails.BadRequest:
				for _, violation := range t.GetFieldViolations() {
					fmt.Printf("The problem with %q field:\n", violation.GetField())
					fmt.Printf("\t%s\n", violation.GetDescription())
				}
			}
		}
		os.Exit("aborting...")
	}
	fmt.Println("BLOCK SERVICE =>", block.Block.String())

	//// Get Latest Height
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	gRPCConn, err := grpc.Dial("0.0.0.0:8080", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer gRPCConn.Close()

	gRPCClient := blocksvc.NewBlockServiceClient(gRPCConn)
	req := blocksvc.GetLatestHeightRequest{}
	gCtx, cancel := context.WithCancel(context.Background())
	stream, err := gRPCClient.GetLatestHeight(gCtx, &req)
	count := 0
	for {
		height, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving new block: %v", err)
		}
		log.Println("New Block:", height)
		count++
		if count > 15 {
			cancel()
			log.Println("Cancelled stream context")
		}
	}
}
