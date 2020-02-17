package main

import (
	"fmt"
	"os"

	"github.com/TRON-US/btfs-migration-toolkit/constants"
	"github.com/TRON-US/btfs-migration-toolkit/core"
	"github.com/TRON-US/btfs-migration-toolkit/service"

	"github.com/spf13/pflag"
)

var (
	config = pflag.StringP("config", "c", "", "soter controller server config file path")
	method = pflag.StringP("method", "m", "", "choose a method to run")
	inputFile = pflag.StringP("input", "i", "", "input .csv file with a list of IPFS QmHash")
	ipfsHash = pflag.StringP("hash", "h", "", "IPFS hash to migrate")
	requestId = pflag.StringP("request-id", "r", "", "request id to query for")
)

func main() {
	pflag.Parse()

	if err := core.InitConfig(*config); err != nil {
		panic(err)
	}

	switch *method {
	case constants.BatchUpload:
		service.BatchUpload(*inputFile)
	case constants.SingleUpload:
		service.SingleUpload(*ipfsHash)
	case constants.BatchVerify:
		service.BatchVerify(*inputFile)
	case constants.SingleVerify:
		service.SingleVerify(*requestId)
	default:
		fmt.Println("unknown method")
		os.Exit(0)
	}
}
