package main

import (
	"github.com/TRON-US/btfs-migration-toolkit/core"
	"github.com/TRON-US/btfs-migration-toolkit/uploader"
)

func main() {

	if err := core.InitConfig(""); err != nil {
		panic(err)
	}

	uploader.BatchUpload("input.csv")
	//err := uploader.DownloadFromIPFS("QmV12MeejQUCyU5QVbihZCRwajzSnGLcTpBC7FhK3NzHby")
	//if err != nil {
	//	panic(err)
	//}
}
