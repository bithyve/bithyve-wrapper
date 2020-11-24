#!/bin/bash
# Get electrs and bitcoind on a fresh Ubuntu 18 server

sudo apt-get update
sudo apt-get -y upgrade
sudo apt-get install build-essential

wget https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
sudo tar -xvf go1.14.1.linux-amd64.tar.gz
sudo mv go /usr/bin/

export GOROOT=/usr/bin/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
go version

curl https://sh.rustup.rs -sSf | sh
cargo version
source ~/.profile

wget https://bitcoincore.org/bin/bitcoin-core-0.20.1/bitcoin-0.20.1-x86_64-linux-gnu.tar.gz
tar -xvf bitcoin-0.20.1-x86_64-linux-gnu.tar.gz
cd bitcoin-0.20.1/bin
sudo mv bitcoind /usr/local/ # or wherever you would like to have bitcoind
sudo mv bitcoin-cli /usr/locals
# bitcoind -daemon -testnet # start bitcoind on testnet or mainnet
sudo screen -SL btc bitcoind -server=1 -txindex=0 -prune=0

# Clone, build and run electrs
git clone -b new-index https://github.com/Blockstream/electrs.git
cd electrs
sudo apt-get install gcc
sudo apt-get install libclang-dev
sudo apt-get install clang
sudo apt-get install librocksdb-sys # or librocksdb-dev
cargo build
# screen -SL electrs cargo run --release --bin electrs -- -vvvv --daemon-dir /home/ubuntu/.bitcoin --daemon-rpc-addr 127.0.0.1:18332 --network testnet --cors 0.0.0.0/0
sudo screen -SL indexer ./target/release/electrs -vvvv --daemon-dir /home/ubuntu/.bitcoin --daemon-rpc-addr 127.0.0.1:8332 --cors 0.0.0.0/0 --network mainnet --bulk-index-threads 4 --index-batch-size 1000 --tx-cache-size 1000000

go get github.com/bithyve/bithyve-wrapper
cd ~/go/src/github.com/bithyve/bithyve-wrapper

# Get an SSL certificate
sudo certbot certonly --standalone --preferred-challenges http-01 -d api.bithyve.com
sudo cd /etc/letsencrypt/live/api.bithyve.com
cp fullchain.pem server.crt ; cp privkey.pem server.key ; mv server.* /home/ubuntu/go/src/github.com/bithyve/bithyve-wrapper/ssl/

# Run the bithyve wrapper
cd ~/go/src/github.com/bithyve/bithyve-wrapper
go get ./...
go build
sudo screen -SL wrapper ./bithyve-wrapper -m

sudo screen -SL socat80 socat tcp-listen:80,reuseaddr,fork tcp:localhost:3001
sudo screen -SL socat443 socat tcp-listen:443,reuseaddr,fork tcp:localhost:445