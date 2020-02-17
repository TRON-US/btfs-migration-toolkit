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
```bash
./btfs-migration-toolkit --config=config.yaml --method=batch_upload --input=input.csv
```

#### Single upload
```bash
./btfs-migration-toolkit --config=config.yaml --method=single_upload --hash=Qmxxxxxx
```

#### Batch verify
```bash
./btfs-migration-toolkit --config=config.yaml --method=batch_verify --input=output_hash.csv
```

#### Single verify
```bash
./btfs-migration-toolkit --config=config.yaml --method=single_verify --request-id=xxxxx
```