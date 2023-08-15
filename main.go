package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cometbft/cometbft/libs/os"
	client "github.com/cometbft/cometbft/rpc/grpc/client"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30*time.Second))
	defer cancel()
	grpcURI := "0.0.0.0:8080"
	conn, err := client.New(ctx, grpcURI, client.WithVersionServiceEnabled(true), client.WithInsecure())
	defer conn.Close()
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

	//// Get Block
	blockResults, err := conn.GetBlockResults(ctx, 9)
	if err != nil {
		fmt.Printf("error getting block results: %v\n", err)
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
	fmt.Println("BLOCK RESULTS SERVICE =>", blockResults.Height, blockResults.ValidatorUpdates, blockResults.FinalizeBlockEvents, blockResults.AppHash)

	gCtx, cancel := context.WithCancel(context.Background())
	ch, err := conn.GetLatestHeight(gCtx)
	if err != nil {
		fmt.Println("Error getting latest height:", err)
	}
	count := 0

	for {
		result := <-ch
		if result.Error != nil {
			fmt.Println("Error on new block:", result.Error)
		} else {
			fmt.Println("New block at height:", result.Height)
		}
		count++
		if count > 5 {
			break
		}
	}

	os.Exit("finished receiving new blocks")
}
