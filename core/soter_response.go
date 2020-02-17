package core

type SoterAddFileResponse struct {
	Cid       string `json:"cid"`
	RequestId string `json:"request_id"`
}

type SoterOrderDetailsResponse struct {
	FileHash    string `json:"file_hash"`
	FileSize    int64  `json:"file_size"`
	FileName    string `json:"file_name"`
	Fee         int64  `json:"fee"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Status      string `json:"status"`
}