package uploader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TRON-US/btfs-migration-toolkit/constants"
	"github.com/TRON-US/btfs-migration-toolkit/core"
	"github.com/TRON-US/btfs-migration-toolkit/log"
	"github.com/TRON-US/soter-sdk-go/soter"
)

func BatchUpload(inputFilename string) {
	batchSize := core.Conf.BatchSize

	inputHashFile, err := os.Open(inputFilename)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to open file %s, reason=[%v]", inputFilename, err))
		os.Exit(1)
	}
	defer func() {
		if err := inputHashFile.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to close file %s, reason=[%v]", inputFilename, err))
		}
	}()

	outputHashFile, err := os.Create(fmt.Sprintf("./%s", constants.OutputHashFileName))
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to open file %s, reason=[%v]", constants.OutputHashFileName, err))
		os.Exit(1)
	}
	defer func() {
		if err := outputHashFile.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to close file %s, reason=[%v]",
				constants.OutputHashFileName, err))
		}
	}()
	outputRetryFile, err := os.Create(fmt.Sprintf("./%s", constants.OutputRetryFileName))
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to open %s, reason=[%v]", constants.OutputRetryFileName, err))
	}

	wg := sync.WaitGroup{}
	scanner := bufio.NewScanner(inputHashFile)
	counter := 0
	for scanner.Scan() {
		hash := scanner.Text()
		wg.Add(1)
		counter++
		go func(h string, outFile *os.File, retryFile *os.File) {
			defer wg.Done()
			time.Sleep(time.Second * 1)
			res, err := migrate(h)
			if err != nil {
				log.Logger().Error(fmt.Sprintf("ipfs_hash=%s, reason=[%v]", h, err))
				// definitely failed to upload through soter; write to output_retry.csv
				_, err = fmt.Fprintln(retryFile, h)
				if err != nil {
					errMsg := fmt.Sprintf("Failed to write to file %s, hash=%s, reason=[%v]",
						constants.OutputRetryFileName, h, err)
					log.Logger().Error(errMsg)
				}
				return
			}
			log.Logger().Debug(fmt.Sprintf("[%s,%s,%s]", h, res[0], res[1]))
			// write <ipfs_hash, request_id, btfs_hash> to output_hash.csv
			line := fmt.Sprintf("%s,%s,%s", h, res[0], res[1])
			_, err = fmt.Fprintln(outFile, line)
			if err != nil {
				log.Logger().Error(err.Error())
			}
		}(hash, outputHashFile, outputRetryFile)
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

func SingleUpload(ipfsHash string)  {
	res, err := migrate(ipfsHash)
	if err != nil {
		log.Logger().Error(err.Error())
		os.Exit(1)
	}
	fmt.Printf("IPFS hash: %s\n", ipfsHash)
	fmt.Printf("BTFS hash: %s\n", res[1])
	fmt.Printf("Request ID: %s\n", res[0])
}

func migrate(ipfsHash string) ([]string, error) {
	if !strings.HasPrefix(ipfsHash, "Qm") {
		errMsg := fmt.Sprintf("input with invalid IPFS hash [%s]", ipfsHash)
		log.Logger().Debug(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	// download file from IPFS network to local file system
	log.Logger().Debug(fmt.Sprintf("downloading file from IPFS network, %s", ipfsHash))
	if err := downloadFromIPFS(ipfsHash); err != nil {
		return nil, err
	}

	// upload the file to BTFS through soter
	res, err := uploadToBTFS(ipfsHash)
	if err != nil {
		return nil, err
	}
	// delete local files
	if err := os.Remove(fmt.Sprintf("./%s", ipfsHash)); err != nil {
		errMsg := fmt.Sprintf("Failed to delete file %s", ipfsHash)
		log.Logger().Error(errMsg)
	}
	return res, nil
}

func downloadFromIPFS(hash string) error {
	// go-ipfs-api, sdk: get
	if err := core.Sh.Get(hash, hash); err != nil {
		return err
	}

	return nil
}

func uploadToBTFS(filename string) ([]string, error) {
	sh := soter.NewShell(core.Conf.PrivateKey, core.Conf.UserAddress, core.Conf.SoterUrl)
	filePath := fmt.Sprintf("./%s", filename)
	resp, err := sh.AddFile(core.Conf.UserAddress, filePath)
	if err != nil {
		log.Logger().Error(err.Error())
		return nil, err
	}
	if resp.Code != constants.OkCode {
		if resp.Code == constants.InsufficientBalanceCode {
			os.Exit(1)
		}
		errMsg := fmt.Sprintf("Error: code=%d, message=%s", resp.Code, resp.Message)
		log.Logger().Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	s, err := json.Marshal(resp.Data)
	if err != nil {
		log.Logger().Error(err.Error())
		return nil, err
	}
	var soterResponse core.SoterResponse
	err = json.Unmarshal(s, &soterResponse)
	if err != nil {
		log.Logger().Error(err.Error())
		return nil, err
	}
	res := [...]string{soterResponse.RequestId, soterResponse.Cid}
	return res[:], nil
}