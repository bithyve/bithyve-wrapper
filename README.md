# BitHyve Wrapper

BitHyve wrapper is a server instance that wraps around [electrs](https://github.com/Blockstream/electrs) to provide some additional functionality on top of electrs as required by [hexa](https://github.com/bithyve/hexa)

## Prerequisites for running BitHyve Wrappe

BitHyve Wrapper requires **electrs** (https://github.com/Blockstream/electrs) and electrs requires a **Bitcoin Core** (v0.16+)

**Bitcoin Core daemon**

Bitcoin Core can be downloaded from https://bitcoincore.org/en/download/ 

Detailed instructions on installing, configuring and running Bitcoin Core as daemon are available here  https://bitcoin.org/en/full-node

**Electrs**

Electrs can be installed from https://github.com/Blockstream/electrs

Please follow the instructions here https://github.com/Blockstream/electrs#installing--indexing for help in installing and setting up electrs.

## Installing and running BitHyve Wrapper

1. Install golang (replace 1.15.4 with your favorite version)

```
wget https://dl.google.com/go/go1.15.4.linux-amd64.tar.gz
sudo tar -xvf go1.15.4.linux-amd64.tar.gz
sudo mv go /usr/bin/
```

2. Update path

```
export GOROOT=/usr/bin/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

3. Download and build

```
go get github.com/bithyve/bithyve-wrapper
cd ~/go/src/github.com/bithyve/bithyve-wrapper
go get ./...
go build
```

4. Run the wrapper
   - Before running the wrapper please ensure that Bitcoin Core daemon and electrs have been setup and and are running as these are required for BitHyve Wrapper.
   - You would need to know the ip number of the machine/server runnig the BitHye Wrapper and port number

```
sudo screen -SL wrapper ./bithyve-wrapper -m
```

5. Ensure your server accepts http traffic

```
sudo screen -SL socat80 socat tcp-listen:80,reuseaddr,fork tcp:localhost:3001
```
