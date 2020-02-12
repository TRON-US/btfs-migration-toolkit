package main

import (
	"fmt"
	"github.com/TRON-US/btfs-migration-toolkit/constants"
	"github.com/TRON-US/btfs-migration-toolkit/core"
	"github.com/TRON-US/btfs-migration-toolkit/uploader"
	"os"

	"github.com/spf13/pflag"
)

var (
	config = pflag.StringP("config", "c", "", "soter controller server config file path")
	method = pflag.StringP("method", "m", "", "choose a method to run")
	inputFile = pflag.StringP("input", "i", "", "input .csv file with a list of IPFS QmHash")
	ipfsHash = pflag.StringP("hash", "h", "", "IPFS hash to migrate")
)

func main() {
	pflag.Parse()

	if err := core.InitConfig(*config); err != nil {
		panic(err)
	}

	switch *method {
	case constants.BatchUpload:
		uploader.BatchUpload(*inputFile)
	case constants.SingleUpload:
		uploader.SingleUpload(*ipfsHash)
	default:
		fmt.Println("unknown method")
		os.Exit(0)
	}

	fmt.Printf("Migration complete.\n" +
		"Please checkout %s and %s for batch migration.\n" +
		"Or review screen output for single migration.\n",
		constants.OutputHashFileName, constants.OutputRetryFileName)
}
