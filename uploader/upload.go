package uploader

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/TRON-US/btfs-migration-toolkit/core"
	"github.com/TRON-US/btfs-migration-toolkit/log"

	"github.com/satori/go.uuid"
)

func BatchUpload(inputHashFile string) {
	batchSize := core.Conf.BatchSize

	file, err := os.Open(inputHashFile)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to open file %s, reason=[%v]", inputHashFile, err))
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to close file %s, reason=[%v]", inputHashFile, err))
		}
	}()

	outputHash, err := os.OpenFile("output_hash.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// TODO
	}
	bw := bufio.NewWriter(outputHash)
	wg := sync.WaitGroup{}
	scanner := bufio.NewScanner(file)
	counter := 0
	for scanner.Scan() {
		hash := scanner.Text()
		wg.Add(1)
		counter++
		go func(h string, writer *bufio.Writer) {
			defer wg.Done()
			fmt.Println(h)
			time.Sleep(time.Second * 1)
			res, err := SingleUpload(h)
			if err != nil {
				log.Logger().Error(err.Error())
				// TODO: definitely failed to upload through soter; write to output_retry.csv
			}
			// write <ipfs_hash, request_id, btfs_hash> to output_hash.csv
			line := fmt.Sprintf("%s,%s,%s\n", h, res[0], res[1])
			_, err = writer.WriteString(line)
			if err != nil {
				log.Logger().Error(err.Error())
			}
		}(hash, bw)
		if counter % batchSize == 0 {
			wg.Wait()
			counter = 0
		}
	}
	// wait here because counter < batchSize and no more lines to read
	wg.Wait()
	if err := scanner.Err(); err != nil {
		log.Logger().Error(err.Error())
	}
}

func SingleUpload(ipfsHash string) ([]string, error) {
	// download file from IPFS network to local file system
	if err := downloadFromIPFS(ipfsHash); err != nil {
		return nil, err
	}

	// upload the file to BTFS through soter
	requestId := uuid.NewV4().String()
	btfsHash, err := uploadToBTFS(ipfsHash, requestId)
	if err != nil {
		// TODO:
	}
	res := [...]string{requestId, btfsHash}
	return res[:], nil
}

func uploadToBTFS(filename, requestId string) (string, error) {
	// TODO: make request to Soter

	return "", nil
}