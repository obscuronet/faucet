# Obscuro Faucet

This repository contains a Faucet to allow allocation of OBX tokens within an Obscuro network. For more information 
on Obscuro see the [Obscuro repo](https://github.com/obscuronet/go-obscuro) and [documentation](https://docs.obscu.ro/).

## Repository Structure
The top level structure of the project is as below;

```
├── .github                    # GitHub workflow actions for deployment 
├── Dockerfile                 # Docker file to build container
├── README.md                  # Readme
├── cmd                        # Source code for the CLI application
├── container_build.sh         # Build a local container
├── container_run.sh           # Run a local container
├── faucet                     # Source code for faucet implementation
├── go.mod                     # Golang dependency management 
└── go.sum                     # Goland dependency checksums
```

## Building a local container
To build a local container and run the Faucet use the below;

```bash
$ ./container_build.sh 
$ ./container_run.sh 
```

By default, when running locally the Faucet will connect to a local testnet started as described in the go-obscuro 
project repo [readme](https://github.com/obscuronet/go-obscuro#building-and-running-a-local-testnet). It will connect 
to the Obcuro node running within the local testnet on host `testnet-host-1` and port `13000`. The Faucet opens 
on port `80` within the container, but maps port `8080` on the host machine to this.


## Allocating OBX to and EOA
Allocating OBX to an externally owned account is done through a POST command to the `/fund/obx` endpoint, where the 
data in the POST command specifies the address e.g. for the account `0x0d2166b7b3A1522186E809e83d925d7b0B6db084`

```bash
curl --location --request POST 'http://127.0.0.1:8080/fund/obx' \
--header 'Content-Type: application/json' \
--data-raw '{ "address":"0x0d2166b7b3A1522186E809e83d925d7b0B6db084" }'
```

