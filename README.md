# Verifying Digital Files Via Hyperledger Fabric Blockchain
The codes are mainly to provide users with intuitive interfaces to save the hashes of valuable digital files or archichives into Hyperledger Blockchain as proofs as well as fetching proofs from Blockchain. The codes also include font-end ones and back-end ones in which the back-end ones include APIs dealing with requests and responses from font end, accessing the Blockchain through FabricClient framework, and operating the Mssql database.

The following content is for how to configurate Blockchain Nodes including Peer nodes, Orderer nodes and CA nodes as well as how to deploy Smart Contracts to Blockchain and operate(save and fetch) data in Blockchain using Smart Contracts.

Here, for purpose of test, we utilize virtual machines installed Centos7 and Dockers as the infrastructure to operate Blockchain nodes. We use one virtual machine for each Organization where several Peer nodes could be settled down.

## Two Peer nodes with one Orderer node
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
