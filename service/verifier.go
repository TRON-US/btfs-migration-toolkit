package service

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

func BatchVerify(filename string) {
	inputFile, err := os.Open(filename)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to open file %s, reason=[%v]", filename, err))
		os.Exit(1)
	}
	defer func() {
		if err := inputFile.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to close file %s, reason=[%v]", filename, err))
		}
	}()

	outputDir := fmt.Sprintf("output-%d", time.Now().Unix())
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to create dir %s, reason=[%v]", outputDir, err))
			os.Exit(1)
		}
	}

	outputP, err := os.Create(fmt.Sprintf("./%s/%s", outputDir, constants.OutputPendingFileName))
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to open file %s, reason=[%v]",
			constants.OutputPendingFileName, err))
		os.Exit(1)
	}
	defer func() {
		if err := outputP.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to close file %s, reason=[%v]",
				constants.OutputPendingFileName, err))
		}
	}()
	outputF, err := os.Create(fmt.Sprintf("./%s/%s", outputDir, constants.OutputFailFileName))
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to open file %s, reason=[%v]",
			constants.OutputFailFileName, err))
		os.Exit(1)
	}
	defer func() {
		if err := outputF.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to close file %s, reason=[%v]",
				constants.OutputFailFileName, err))
		}
	}()
	outputS, err := os.Create(fmt.Sprintf("./%s/%s", outputDir, constants.OutputSucessFileName))
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to open file %s, reason=[%v]",
			constants.OutputSucessFileName, err))
		os.Exit(1)
	}
	defer func() {
		if err := outputS.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to close file %s, reason=[%v]",
				constants.OutputSucessFileName, err))
		}
	}()
	outputE, err := os.Create(constants.OutputErrorFileName)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Failed to open file %s, reason=[%v]",
			constants.OutputErrorFileName, err))
		os.Exit(1)
	}
	defer func() {
		if err := outputE.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("Failed to close file %s, reason=[%v]",
				constants.OutputErrorFileName, err))
		}
	}()

	wg := sync.WaitGroup{}
	scanner := bufio.NewScanner(inputFile)
	counter := 0
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		counter++
		log.Logger().Debug(fmt.Sprintf("Verifying %s", line))

		go func(str string, pFile, fFile, sFile, eFile *os.File) {
			defer wg.Done()

			words := strings.Split(str, constants.Delimiter)
			s, err := verify(words[1])
			if err != nil {
				log.Logger().Error(fmt.Sprintf("[request_id=%s] Failed to verify, reason=[%v]", words[1], err))
				_, err := fmt.Fprintln(eFile, str)
				if err != nil {
					errMsg := fmt.Sprintf("[request_id=%s] Failed to write to file %s, reason=[%v]",
						words[1], constants.OutputErrorFileName, err)
					log.Logger().Error(errMsg)
				}
				return
			}
			if s == constants.StatusP {
				_, err := fmt.Fprintln(pFile, str)
				if err != nil {
					errMsg := fmt.Sprintf("[request_id=%s] Failed to write to file %s, reason=[%v]",
						words[1], constants.OutputPendingFileName, err)
					log.Logger().Error(errMsg)
				}
			} else if s == constants.StatusS {
				_, err := fmt.Fprintln(sFile, str)
				if err != nil {
					errMsg := fmt.Sprintf("[request_id=%s] Failed to write to file %s, reason=[%v]",
						words[1], constants.OutputSucessFileName, err)
					log.Logger().Error(errMsg)
				}
			} else if s == constants.StatusF {
				_, err := fmt.Fprintln(fFile, words[0])
				if err != nil {
					errMsg := fmt.Sprintf("[request_id=%s] Failed to write to file %s, reason=[%v]",
						words[1], constants.OutputFailFileName, err)
					log.Logger().Error(errMsg)
				}
			} else {
				_, err := fmt.Fprintln(eFile, str)
				if err != nil {
					errMsg := fmt.Sprintf("[request_id=%s] Failed to write to file %s, reason=[%v]",
						words[1], constants.OutputErrorFileName, err)
					log.Logger().Error(errMsg)
				}
			}
		}(line, outputP, outputF, outputS, outputE)

		if counter % core.Conf.BatchSize == 0 {
			wg.Wait()
			counter = 0
		}
	}
	wg.Wait()
	if err := scanner.Err(); err != nil {
		errMsg := fmt.Sprintf("Failed to scan input file, reason=[%v]", err)
		log.Logger().Error(errMsg)
	}

	fmt.Printf("\nVerification complete.\n" +
		"Please checkout %s, %s, %s, and %s for batch verification.\n",
		constants.OutputPendingFileName, constants.OutputFailFileName,
		constants.OutputSucessFileName, constants.OutputErrorFileName)
}

func SingleVerify(requestId string) {
	s, err := verify(requestId)
	if err != nil {
		fmt.Printf("Failed to verify %s, reason=[%v]", requestId, err)
		os.Exit(1)
	}
	status := "unknown"
	switch s {
	case constants.StatusP:
		status = "pending"
	case constants.StatusF:
		status = "failed"
	case constants.StatusS:
		status = "success"
	}

	fmt.Printf(
		"The request id is %s\n" +
		"The upload status is %s\n",
		requestId, status)
}

func verify(requestId string) (string, error) {
	sh := soter.NewShell(core.Conf.PrivateKey, core.Conf.UserAddress, core.Conf.SoterUrl)
	resp, err := sh.QueryOrderDetails(requestId)
	if err != nil {
		return "", fmt.Errorf("failed to query order details, %v", err)
	}
	if resp.Code != constants.OkCode {
		errMsg := fmt.Sprintf("response code error: code=%d, message=%s", resp.Code, resp.Message)
		return "", fmt.Errorf(errMsg)
	}
	s, err := json.Marshal(resp.Data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response data, %v", err)
	}
	var soterResponse core.SoterOrderDetailsResponse
	err = json.Unmarshal(s, &soterResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal soter response data, %v", err)
	}

	return soterResponse.Status, nil
}
