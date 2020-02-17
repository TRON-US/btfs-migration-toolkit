package constants

const (
	OkCode                  = 0
	InsufficientBalanceCode = 20009

	// files has been downloaded from IPFS network and uploaded to BTFS network through Soter
	// BTFS hash and request id have been return from Soter and saved to this file
	OutputHashFileName    = "output_hash.csv"
	// the migration has been failed, and IPFS hash is saved in this file for retry
	OutputRetryFileName   = "output_retry.csv"
	// record identified by this request id is in Pending status
	OutputPendingFileName = "output_P.csv"
	// record identified by this request id is in Failed status
	OutputFailFileName    = "output_F.csv"
	// record identified by this request id is in Success status
	OutputSucessFileName  = "output_S.csv"
	// this request id is invalid or verification is abort due to internal error
	OutputErrorFileName   = "output_E.csv"

	// method
	BatchUpload  = "batch_upload"
	SingleUpload = "single_upload"
	BatchVerify  = "batch_verify"
	SingleVerify = "single_verify"

	//
	Delimiter = ","
	StatusP   = "P"
	StatusF   = "F"
	StatusS   = "S"
)
