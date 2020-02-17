package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/TRON-US/btfs-migration-toolkit/constants"
	"github.com/TRON-US/btfs-migration-toolkit/core"
	"github.com/TRON-US/btfs-migration-toolkit/log"
	"github.com/TRON-US/soter-sdk-go/soter"

	"github.com/ipfs/go-ipfs-api"
)

func BatchUpload(inputFilename string) {
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

	outputHashFile, err := os.Create(constants.OutputHashFileName)
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

	outputRetryFile, err := os.Create(constants.OutputRetryFileName)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to open file %s, reason=[%v]", constants.OutputRetryFileName, err))
		os.Exit(1)
	}
	defer func() {
		if err := outputRetryFile.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to close file %s, reason=[%v]",
				constants.OutputRetryFileName, err))
		}
	}()

	wg := sync.WaitGroup{}
	scanner := bufio.NewScanner(inputHashFile)
	counter := 0
	for scanner.Scan() {
		hash := scanner.Text()
		wg.Add(1)
		counter++
		go func(h string, outFile, retryFile *os.File) {
			defer wg.Done()
			res, err := migrate(h)
			if err != nil {
				log.Logger().Error(fmt.Sprintf("[ipfs_hash=%s] Failed to migrate, reason=[%v]", h, err))
				// definitely failed to upload through soter; write to output_retry.csv
				_, err = fmt.Fprintln(retryFile, h)
				if err != nil {
					errMsg := fmt.Sprintf("[ipfs_hash=%s] Failed to write to file %s, reason=[%v]",
						h, constants.OutputRetryFileName, err)
					log.Logger().Error(errMsg)
				}
				return
			}
			// write (ipfs_hash, request_id, btfs_hash) to output_hash.csv
			log.Logger().Debug(fmt.Sprintf("[ipfs_hash=%s] Write to file %s, (%s,%s,%s)",
				h, constants.OutputHashFileName, h, res[0], res[1]))
			line := fmt.Sprintf("%s,%s,%s", h, res[0], res[1])
			_, err = fmt.Fprintln(outFile, line)
			if err != nil {
				log.Logger().Error(err.Error())
			}
		}(hash, outputHashFile, outputRetryFile)
		if counter % core.Conf.BatchSize == 0 {
			wg.Wait()
			counter = 0
		}
	}
	// wait here because counter < batchSize and no more lines to read
	wg.Wait()
	if err := scanner.Err(); err != nil {
		errMsg := fmt.Sprintf("Failed to scan input file, reason=[%v]", err)
		log.Logger().Error(errMsg)
	}

	fmt.Printf("\nMigration complete.\n" +
		"Please checkout %s and %s for batch migration.\n",
		constants.OutputHashFileName, constants.OutputRetryFileName)
}

func SingleUpload(ipfsHash string)  {
	res, err := migrate(ipfsHash)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("[ipfs_hash=%s] Failed to migrate, reason=[%v]", ipfsHash, err))
		os.Exit(1)
	}
	fmt.Printf("IPFS hash: %s\n", ipfsHash)
	fmt.Printf("BTFS hash: %s\n", res[1])
	fmt.Printf("Request ID: %s\n", res[0])
}

func migrate(ipfsHash string) ([]string, error) {
	if !strings.HasPrefix(ipfsHash, "Qm") {
		errMsg := fmt.Sprintf("[ipfs_hash=%s] Input with invalid ipfs hash", ipfsHash)
		return nil, fmt.Errorf(errMsg)
	}
	// download file from IPFS network to local file system
	log.Logger().Debug(fmt.Sprintf("[ipfs_hash=%s] Downloading the file from IPFS network", ipfsHash))
	if err := downloadFromIPFS(ipfsHash); err != nil {
		return nil, err
	}

	// defer os.Remove(the_downloaded_file)
	defer func(h string) {
		// delete local files
		if err := os.Remove(fmt.Sprintf("./%s", h)); err != nil {
			errMsg := fmt.Sprintf("[ipfs_hash=%s] Failed to delete local file", h)
			log.Logger().Error(errMsg)
		}
	}(ipfsHash)

	// upload the file to BTFS through soter
	log.Logger().Debug(fmt.Sprintf("[ipfs_hash=%s] Uploading the file to BTFS network", ipfsHash))
	res, err := uploadToBTFS(ipfsHash)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func downloadFromIPFS(hash string) error {
	// go-ipfs-api, sdk: get
	sh := shell.NewShell(core.Conf.IpfsUrl)
	if err := sh.Get(hash, hash); err != nil {
		return fmt.Errorf("downloading from IPFS errors out, %v", err)
	}

	return nil
}

func uploadToBTFS(filename string) ([]string, error) {
	sh := soter.NewShell(core.Conf.PrivateKey, core.Conf.UserAddress, core.Conf.SoterUrl)
	filePath := fmt.Sprintf("./%s", filename)
	resp, err := sh.AddFile(core.Conf.UserAddress, filePath)
	if err != nil {
		errMsg := fmt.Sprintf("failed to add file, reason=[%v]", err)
		return nil, fmt.Errorf(errMsg)
	}
	if resp.Code != constants.OkCode {
		errMsg := fmt.Sprintf("response code error: code=%d, message=%s", resp.Code, resp.Message)
		if resp.Code == constants.InsufficientBalanceCode {
			os.Exit(1)
		}
		return nil, fmt.Errorf(errMsg)
	}
	s, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response data, %v", err)
	}
	var soterResponse core.SoterAddFileResponse
	err = json.Unmarshal(s, &soterResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal soter response data, %v", err)
	}
	res := [...]string{soterResponse.RequestId, soterResponse.Cid}
	return res[:], nil
}