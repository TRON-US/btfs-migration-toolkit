package uploader

import (
	"github.com/TRON-US/btfs-migration-toolkit/core"
)

func downloadFromIPFS(hash string) error {
	// method 1: curl: http://127.0.0.1:8080/ipfs/QmHash

	// method 2: go-ipfs-sdk: cat/get
	if err := core.Sh.Get(hash, hash); err != nil {
		return err
	}

	return nil
}
