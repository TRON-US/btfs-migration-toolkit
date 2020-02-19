# btfs-migration-tools
> This is a command line tool to migrate ipfs file to btfs network

## Get the source code
```bash
go get -u github.com/TRON-US/btfs-migration-toolkit
```

## Build the binary
```bash
go build
```

## Usage
The file migration from IPFS network to BTFS network has to be done in two phases. 

Phase one: download a file from IPFS to local file system and upload this file to BTFS network 
through Soter using soter-go-sdk. 

Phase two: verify a file has indeed been stored in BTFS network by BTT Integration Architecture.
This is achieved by using the request id to query soter status record.

Before using this command line tool, users should configure the config.yaml file first and provide a input.csv file.
### Example
#### Batch upload
The content in input.csv is just a list of lines of IPFS file hashes. Batch upload will produce two files, output_hash.csv
and output_retry.csv. 

output_hash.csv stores all the (ipfs_hash, request_id, btfs_hash) tuples in one file,
indicating these files identified by corresponding IPFS hash have been migrated to BTFS network,and waiting to be verified.

And output_retry stores all the IPFS hashes that need to re-migrate again due to errors such as internal errors returned
from Soter, insufficient balance, invalid IPFS hash provided, and files being not able to be downloaded from IPFS network.

```bash
./btfs-migration-toolkit --config=config.yaml --method=batch_upload --input=input.csv
```

#### Single upload
Single upload requires user to provide IPFS hash from the command line. After the command being executed, some information
will be printed on the console as stdout if it's been migrated successfully. These information include IPFS hash, 
BTFS hash, and request_id in Soter system.

```bash
./btfs-migration-toolkit --config=config.yaml --method=single_upload --hash=Qmxxxxxx
```

#### Batch verify
Batch verify takes the hash_out.csv file from batch upload as its input, and it will produce a output_timestamp director
as its output. There are 4 files under the output_timestamp folder: output_P.csv, output_F.csv, output_S.csv, and output_E.csv.

output_P.csv stores records that are in pending status, waiting for BTFS BTT Integration to be completed; output_F.csv
means this file failed to store on BTFS network during BTFS BTT Integration phase; output_S.csv means this file has been
migrated completely and successfully, and pass through BTFS BTT Integration phase; output_E.csv file contains records that 
invalid request id is provided and records that stores unkonwn status in Soter database, and ideally this file should be 
an empty output file.

Note: Please wait a while to verify the files after you migrate them through batch upload or single upload, because BTFS 
BTT Integration phase takes some time to complete its internal service and the storage on BTFS network. Otherwise, there
will be plenty of records in pending status rather in failed or success status, and you will need to verify them again.

```bash
./btfs-migration-toolkit --config=config.yaml --method=batch_verify --input=output_hash.csv
```

#### Single verify
Single verify takes soter request-id from batch upload or single upload as the input. And it prints out the status in
BTFS BTT Integration phase. It should be in one of the pending, failed, or success status.

```bash
./btfs-migration-toolkit --config=config.yaml --method=single_verify --request-id=xxxxx
```