# CometBFT gRPC Client Example

This program can be used to test the future gRPC service that will be implemented in CometBFT as part of [the ADR-101 PoC](https://github.com/cometbft/cometbft/issues/816)

To test with this clone and check out one of the branches associated with the new gRPC services, e.g. #1094 (BlockService)
or #1095 (BlockResultsService).

### Install cometbft

Open a new terminal and run `make install` from the cometbft repo

### Init the chain

Run `cometbft init`

### Edit the configuration

Open the config file (`$HOME/.cometbft/config/config.toml`) and in the `[grpc]` section ensure that you have a listen
address and the service you're testing is enabled, for example:

```
[grpc]

# TCP or UNIX socket address for the RPC server to listen on. If not specified,
# the gRPC server will be disabled.
laddr = "tcp://0.0.0.0:8080"

#
# Each gRPC service can be turned on/off, and in some cases configured,
# individually. If the gRPC server is not enabled, all individual services'
# configurations are ignored.
#

# The gRPC version service provides version information about the node and the
# protocols it uses.
[grpc.version_service]
enabled = true

# The gRPC block service returns block information
[grpc.block_service]
enabled = true
```

### Start the service

Start cometbft using the kvstore app

`cometbft start --proxy_app=kvstore`


### Run the example

Open another terminal, navigate to the example repo folder and run this program
to test the gRPC client against the gRPC server running on the other terminal

`go run main.go`

> NOTE: You might have to execute `go get` before the run command

If everything works you should be able to get information back from the server. For example testing the `BlockService` and
the `VersionService`

```
VERSION SERVICE => P2P: 8 | Block: 11 | ABCI: 2.0.0 | Node: 0.39.0-dev
BLOCK SERVICE => Block{
Header{
Version:        {11 1}
ChainID:        test-chain-eXx4bz
Height:         2
Time:           2023-07-11 17:36:13.90971809 +0000 UTC
LastBlockID:    E915A74F1C744D0A7CC49831F0332AC2D6AD6A4B8334AFBA2E8CDE95ACA8817E:1:EC9AC2798621
LastCommit:     B96EAA4B20171E4DEB9C8034F42225E07BEE9B4E7D8945D68EA904BAD40A926C
Data:           E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855
Validators:     D53F49A9A2C798DF7C45FF2C3263DE27E913EF777A8D86A73B0CBE1569415E65
NextValidators: D53F49A9A2C798DF7C45FF2C3263DE27E913EF777A8D86A73B0CBE1569415E65
App:            0000000000000000
Consensus:      048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F
Results:        E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855
...
```
