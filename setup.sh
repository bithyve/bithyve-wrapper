# Get electrs and bitcoind on a fresh Ubuntu 18 server

# Install build-essential
sudo apt-get update
sudo apt-get -y upgrade
sudo apt-get install build-essential

# Install Go
wget https://dl.google.com/go/go1.13.8.linux-amd64.tar.gz
sudo tar -xvf go1.13.8.linux-amd64.tar.gz
sudo mv go /usr/local

# set GO variables
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
go version

# Install Rust
curl https://sh.rustup.rs -sSf | sh
source $HOME/.cargo/env
cargo version

# Get and Run bitcoind
wget https://bitcoincore.org/bin/bitcoin-core-0.19.0.1/bitcoin-0.19.0.1-x86_64-linux-gnu.tar.gz
tar -xvf bitcoin-0.19.0.1-x86_64-linux-gnu.tar.gz
cd bitcoin-0.19.0.1/bin
mv bitcoind /usr/local/ # or wherever you would like to have bitcoind
bitcoind -daemon -testnet # start bitcoind on testnet or mainnet

# Clone, build and run electrs
git clone -b new-index https://github.com/Blockstream/electrs.git
cd electrs
sudo apt-get install libclang-dev
sudo apt-get install clang 
# change below parameters based on your bitcoind parameters
screen -SL electrs cargo run --release --bin electrs -- -vvvv --daemon-dir ~/.bitcoin --cookie=username:password --daemon-rpc-addr 127.0.0.1:18332 --network testnet --cors 0.0.0.0/0
# Go get and build the Bithyve Wrapper
go get github.com/bithyve/bithyve-wrapper
cd ~/go/src/github.com/bithyve/bithyve-wrapper

# Get an SSL certificate
sudo certbot certonly --standalone --preferred-challenges http-01 -d api.bithyve.com
sudo cd /etc/letsencrypt/live/api.bithyve.com
cp fullchain.pem server.crt ; cp privkey.pem server.key ; mv server.* ~/go/src/github.com/bithyve/bithyve-wrapper/ssl/

# Run the bithyve wrapper
cd ~/go/src/github.com/bithyve/bithyve-wrapper
go get ./...
go build
sudo screen -SL wrapper ./bithyve-wrapper

sudo screen -SL socat443 socat tcp-listen:443,reuseaddr,fork tcp:localhost:445 # 445-443 for the wrapper
sudo screen -SL socat80 socat tcp-listen:80,reuseaddr,fork tcp:localhost:3001 # 3001-80 for electrsn