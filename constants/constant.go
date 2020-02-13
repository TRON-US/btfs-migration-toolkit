package constants

const (
	OkCode                  = 0
	InsufficientBalanceCode = 20009
	OrderNotExistCode		= 20008

	// file name
	OutputHashFileName    = "output_hash.csv"
	OutputRetryFileName   = "output_retry.csv"
	OutputPendingFileName = "output_P.csv"
	OutputFailFileName    = "output_F.csv"
	OutputSucessFileName  = "output_S.csv"
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
