# Verifying Digital Files Via Hyperledger Fabric Blockchain
The codes are mainly to provide users with intuitive interfaces to save the hashes of valuable digital files or archichives into Hyperledger Blockchain as proofs as well as fetching proofs from Blockchain. The codes also include font-end ones and back-end ones in which the back-end ones include APIs dealing with requests and responses from font end, accessing the Blockchain through FabricClient framework, and operating the Mssql database.

The following content is for how to configurate Blockchain Nodes including Peer nodes, Orderer nodes and CA nodes as well as how to deploy Smart Contracts to Blockchain and operate(save and fetch) data in Blockchain using Smart Contracts.

Here, for purpose of test, we utilize virtual machines installed Centos7 and Dockers as the infrastructure to operate Blockchain nodes. We use one virtual machine for each Organization where several Peer nodes could be settled down.

## Two Peer Blockchain Nodes and One Orderer Node
In this section, we show the codes for configuring and using two Peer nodes and one Orderer node.

Sychronize the system time using this line of command in Centos.
```sh
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```

Set the production path in Centos.
```sh
$GOPATH/src/github.com/hyperledger/fabric/helloworld
```
We assume Docker has been used installed in Centos.
```sh
systemctl restart docker;
docker stop $(docker ps -a -q);
docker rm $(docker ps -a -q); # Remove all unused local volumes
docker volume prune
docker network rm config_test # delete the specific network config_test which is used for communication among blockchain nodes
```
We assume that Hyperledger Fabric 2.3 has been installed in Centos. This command below is used to generate the public and private certificates for accessing blockchain.
```sh
./bin/cryptogen generate --config=./config/crypto-config.yaml
```
Here we have two Organizations one of which will be deployed one Orderer blockchain node and the other will be deployed two Peer blockchain nodes. The command below is for generating ccp files for the Organization one and Organization two.
```sh
./organizations/ccp-generate.sh
```
We use codes below to generate the First Block here we named channel1.block using configtx.yaml.
```sh
export PATH=${PWD}/bin:$PATH; 
export FABRIC_CFG_PATH=${PWD}/config; 
./bin/configtxgen -profile TwoOrgsApplicationGenesis -outputBlock ./channel-artifacts/channel1.block -channelID channel1
```
Start orderer、peer nodes using the command below.
```sh
docker-compose -f ./config/orderer-peer-cli.yaml up -d  
or docker-compose -f ${PWD}/config/orderer-peer-cli.yaml up -d
```
Generate one channel named channel1 and let Orderer1 join the channel.
```sh
./bin/osnadmin channel join --channel-id channel1 --config-block ./channel-artifacts/channel1.block -o localhost:7053 \
--ca-file ./organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
--client-cert ./organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt \
--client-key ./organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key
```
Make Organization one join the channel1 using the core.yaml located in the installation path of Hypberledger Fabric. 
```sh
export FABRIC_CFG_PATH=${PWD}/config;  
export CORE_PEER_TLS_ENABLED=true;
export CORE_PEER_LOCALMSPID="Org1MSP";
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt;
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp;
export CORE_PEER_ADDRESS=localhost:7051;
./bin/peer channel join -b ./channel-artifacts/channel1.block
```
Make Organization two join the channel1.
```sh
export FABRIC_CFG_PATH=${PWD}/config;
export CORE_PEER_TLS_ENABLED=true;
export CORE_PEER_LOCALMSPID="Org2MSP";
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt;
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp;
export CORE_PEER_ADDRESS=localhost:9051;
./bin/peer channel join -b ./channel-artifacts/channel1.block
```
Anchor peer set for org 'Org1MSP' on channel 'channel1'.
```sh
export CORE_PEER_TLS_ENABLED=true;
export CORE_PEER_LOCALMSPID="Org1MSP";
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt;
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp;
export CORE_PEER_ADDRESS=localhost:7051;

./bin/peer channel fetch config config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c channel1 --tls --cafile \
${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

./bin/configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config >Org1MSPconfig.json;
jq '.channel_group.groups.Application.groups.Org1MSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.org1.example.com","port": 7051}]},"version": "0"}}' Org1MSPconfig.json >Org1MSPmodified_config.json; 
./bin/configtxlator proto_encode --input Org1MSPconfig.json --type common.Config >original_config.pb;
./bin/configtxlator proto_encode --input Org1MSPmodified_config.json --type common.Config >modified_config.pb;
./bin/configtxlator compute_update --channel_id channel1 --original original_config.pb --updated modified_config.pb >config_update.pb;
./bin/configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate >config_update.json;
echo '{"payload":{"header":{"channel_header":{"channel_id":"channel1", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . >config_update_in_envelope.json; 
./bin/configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope >Org1MSPanchors.tx

./bin/peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c channel1 -f Org1MSPanchors.tx --tls --cafile \
${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```
