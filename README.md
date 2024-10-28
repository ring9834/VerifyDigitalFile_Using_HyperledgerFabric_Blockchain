# Verifying Digital Files Via Hyperledger Fabric Blockchain
The codes are mainly to provide users with intuitive interfaces to save the hashes of valuable digital files or archichives into Hyperledger Blockchain as proofs as well as fetching proofs from Blockchain. The codes also include font-end ones and back-end ones in which the back-end ones include APIs dealing with requests and responses from font end, accessing the Blockchain through FabricClient framework, and operating the Mssql database.

The following content is for how to configurate Blockchain Nodes including Peer nodes, Orderer nodes and CA nodes as well as how to deploy Smart Contracts to Blockchain and operate(save and fetch) data in Blockchain using Smart Contracts.

Here, for purpose of test, we utilize virtual machines installed Centos7 and Dockers as the infrastructure to operate Blockchain nodes. We use one virtual machine for each Organization where several Peer nodes could be settled down.

## Two Peer Blockchain Nodes and One Orderer Node
### Configure Blockchain Nodes
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
Anchor peer set for org 'Org2MSP' on channel 'channel1'.
```sh
export CORE_PEER_TLS_ENABLED=true;
export CORE_PEER_LOCALMSPID="Org2MSP";
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt;
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp;
export CORE_PEER_ADDRESS=localhost:9051;

./bin/peer channel fetch config config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c channel1 --tls --cafile \
${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

./bin/configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config >Org2MSPconfig.json;
jq '.channel_group.groups.Application.groups.Org2MSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.org2.example.com","port": 9051}]},"version": "0"}}' Org2MSPconfig.json >Org2MSPmodified_config.json; 
./bin/configtxlator proto_encode --input Org2MSPconfig.json --type common.Config >original_config.pb;
./bin/configtxlator proto_encode --input Org2MSPmodified_config.json --type common.Config >modified_config.pb;
./bin/configtxlator compute_update --channel_id channel1 --original original_config.pb --updated modified_config.pb >config_update.pb;
./bin/configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate >config_update.json;
echo '{"payload":{"header":{"channel_header":{"channel_id":"channel1", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . >config_update_in_envelope.json;
./bin/configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope >Org2MSPanchors.tx

./bin/peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c channel1 -f Org2MSPanchors.tx --tls --cafile \
${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```
Pack the Chaincode (Smart contract).
```sh
./bin/peer lifecycle chaincode package hellocc.tar.gz -p ./chaincode/go/helloworld --label hello_1  
```
### Install Smart Contract on Blockchain Nodes
Install the chaincode on the Organization one peer node.
```sh
export CORE_PEER_TLS_ENABLED=true;
export CORE_PEER_LOCALMSPID="Org1MSP";
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt;
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp;
export CORE_PEER_ADDRESS=localhost:7051;
./bin/peer lifecycle chaincode install hellocc.tar.gz
```

Install the chaincode on the Organization two peer node.
```sh
export CORE_PEER_TLS_ENABLED=true;
export CORE_PEER_LOCALMSPID="Org2MSP";
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt;
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp;
export CORE_PEER_ADDRESS=localhost:9051;
./bin/peer lifecycle chaincode install hellocc.tar.gz
```
Appove the chaincode defination on Organization one peer node.
```sh
export CORE_PEER_TLS_ENABLED=true;
export CORE_PEER_LOCALMSPID="Org1MSP";
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt;
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp;
export CORE_PEER_ADDRESS=localhost:7051;
./bin/peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
--cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
--channelID channel1 --name hello_1 --version 1.0 --package-id hello_1:5cd76591329d8c8fd9d23516484735adf574e88f13b81c0f09ff0330e71dc719 --sequence 1
```
Appove the chaincode defination on Organization two peer node.
```sh
export CORE_PEER_TLS_ENABLED=true;
export CORE_PEER_LOCALMSPID="Org2MSP";
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt;
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp;
export CORE_PEER_ADDRESS=localhost:9051;
./bin/peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
--cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
--channelID channel1 --name hello_1 --version 1.0 --package-id hello_1:5cd76591329d8c8fd9d23516484735adf574e88f13b81c0f09ff0330e71dc719 --sequence 1
```
Check the defination of the chaincode.
```sh
./bin/peer lifecycle chaincode checkcommitreadiness --channelID channel1 --name hello_1 --version 1.0 --sequence 1 --output json
```
Submit the chaincode defination to the channel1.
```sh
./bin/peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
--cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
--channelID channel1 --name hello_1 \
--peerAddresses localhost:7051 \
--tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 \
--tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
--version 1.0 --sequence 1
```
### Accessing Data on Blockchain
Invode the InitLedger function in the chaincode.
```sh
export PATH=${PWD}/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/config/ #设置FABRIC_CFG_PATH 指向core.yaml文件
export CORE_PEER_TLS_ENABLED=true #以下为设置环境变量为org1的，即在org1上调用链码
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
./bin/peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C channel1 -n hello_1 \
--peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
-c '{"function":"InitLedger","Args":[]}'
```
Query the data saved on Blockchain.
```sh
./bin/peer chaincode query -C channel1 -n hello_1 -c '{"Args":["GetAllAssets"]}'
```
Transfer assets from one account to another.
```sh
./bin/peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C channel1 -n hello_1 --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"TransferAsset","Args":["asset6","Christopher"]}'
```
Query the data saved on Blockchain again to check if the transferring is a success.
```sh
./bin/peer chaincode query -C channel1 -n hello_1 -c '{"Args":["GetAllAssets"]}'
```
### Attachment: configtx.yaml
Codes below are the content of file configtx.yaml for Fabric CA.
```sh
################################################################################
#
#   Section: Organizations
#
#   - This section defines the different organizational identities which will
#   be referenced later in the configuration.
#
################################################################################
Organizations:

    # SampleOrg defines an MSP using the sampleconfig.  It should never be used
    # in production but may be used as a template for other definitions
    - &OrdererOrg
        # DefaultOrg defines the organization which is used in the sampleconfig
        # of the fabric.git development environment
        Name: OrdererOrg

        # ID to load the MSP definition as
        ID: OrdererMSP

        # MSPDir is the filesystem path which contains the MSP configuration
        MSPDir: ../organizations/ordererOrganizations/example.com/msp

        # Policies defines the set of policies at this level of the config tree
        # For organization policies, their canonical path is usually
        #   /Channel/<Application|Orderer>/<OrgName>/<PolicyName>
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('OrdererMSP.member')"
            Writers:
                Type: Signature
                Rule: "OR('OrdererMSP.member')"
            Admins:
                Type: Signature
                Rule: "OR('OrdererMSP.admin')"

        OrdererEndpoints:
            #- orderer.example.com:7050
            - 192.168.181.131:7050
            - 192.168.181.137:7050
            - 192.168.181.138:7050

    - &Org1
        # DefaultOrg defines the organization which is used in the sampleconfig
        # of the fabric.git development environment
        Name: Org1MSP

        # ID to load the MSP definition as
        ID: Org1MSP

        MSPDir: ../organizations/peerOrganizations/org1.example.com/msp

        # Policies defines the set of policies at this level of the config tree
        # For organization policies, their canonical path is usually
        #   /Channel/<Application|Orderer>/<OrgName>/<PolicyName>
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('Org1MSP.member')"
                #Rule: "OR('Org1MSP.admin', 'Org1MSP.peer', 'Org1MSP.client')"
            Writers:
                Type: Signature
                Rule: "OR('Org1MSP.member')"
                #Rule: "OR('Org1MSP.admin', 'Org1MSP.client')"
            Admins:
                Type: Signature
                Rule: "OR('Org1MSP.admin')"
            Endorsement:
                Type: Signature
                Rule: "OR('Org1MSP.member')"
                #Rule: "OR('Org1MSP.peer')"
                
    - &Org2
        # DefaultOrg defines the organization which is used in the sampleconfig
        # of the fabric.git development environment
        Name: Org2MSP

        # ID to load the MSP definition as
        ID: Org2MSP

        MSPDir: ../organizations/peerOrganizations/org2.example.com/msp

        # Policies defines the set of policies at this level of the config tree
        # For organization policies, their canonical path is usually
        #   /Channel/<Application|Orderer>/<OrgName>/<PolicyName>
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('Org2MSP.member')"
                #Rule: "OR('Org2MSP.admin', 'Org2MSP.peer', 'Org2MSP.client')"
            Writers:
                Type: Signature
                Rule: "OR('Org2MSP.member')"
                #Rule: "OR('Org2MSP.admin', 'Org2MSP.client')"
            Admins:
                Type: Signature
                Rule: "OR('Org2MSP.admin')"
            Endorsement:
                Type: Signature
                Rule: "OR('Org2MSP.member')"
                #Rule: "OR('Org2MSP.peer')"


Capabilities:
    # Channel capabilities apply to both the orderers and the peers and must be
    # supported by both.
    # Set the value of the capability to true to require it.
    Channel: &ChannelCapabilities
        # V2_0 capability ensures that orderers and peers behave according
        # to v2.0 channel capabilities. Orderers and peers from
        # prior releases would behave in an incompatible way, and are therefore
        # not able to participate in channels at v2.0 capability.
        # Prior to enabling V2.0 channel capabilities, ensure that all
        # orderers and peers on a channel are at v2.0.0 or later.
        V2_0: true

    # Orderer capabilities apply only to the orderers, and may be safely
    # used with prior release peers.
    # Set the value of the capability to true to require it.
    Orderer: &OrdererCapabilities
        # V2_0 orderer capability ensures that orderers behave according
        # to v2.0 orderer capabilities. Orderers from
        # prior releases would behave in an incompatible way, and are therefore
        # not able to participate in channels at v2.0 orderer capability.
        # Prior to enabling V2.0 orderer capabilities, ensure that all
        # orderers on channel are at v2.0.0 or later.
        V2_0: true

    # Application capabilities apply only to the peer network, and may be safely
    # used with prior release orderers.
    # Set the value of the capability to true to require it.
    Application: &ApplicationCapabilities
        # V2_0 application capability ensures that peers behave according
        # to v2.0 application capabilities. Peers from
        # prior releases would behave in an incompatible way, and are therefore
        # not able to participate in channels at v2.0 application capability.
        # Prior to enabling V2.0 application capabilities, ensure that all
        # peers on channel are at v2.0.0 or later.
        V2_0: true


Application: &ApplicationDefaults

    # Organizations is the list of orgs which are defined as participants on
    # the application side of the network
    Organizations:

    # Policies defines the set of policies at this level of the config tree
    # For Application policies, their canonical path is
    #   /Channel/Application/<PolicyName>
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
        LifecycleEndorsement:
            Type: ImplicitMeta
            Rule: "MAJORITY Endorsement"
        Endorsement:
            Type: ImplicitMeta
            Rule: "MAJORITY Endorsement"

    Capabilities:
        <<: *ApplicationCapabilities

Orderer: &OrdererDefaults

    # Orderer Type: The orderer implementation to start
    OrdererType: etcdraft
    
    # Addresses used to be the list of orderer addresses that clients and peers
    # could connect to.  However, this does not allow clients to associate orderer
    # addresses and orderer organizations which can be useful for things such
    # as TLS validation.  The preferred way to specify orderer addresses is now
    # to include the OrdererEndpoints item in your org definition
    Addresses:
        - 192.168.181.131:7050
        - 192.168.181.137:7050
        - 192.168.181.138:7050

    EtcdRaft:
        Consenters:
        - Host: 192.168.181.131
          Port: 7050
          ClientTLSCert: ../organizations/ordererOrganizations/example.com/orderers/orderer1.example.com/tls/server.crt
          ServerTLSCert: ../organizations/ordererOrganizations/example.com/orderers/orderer1.example.com/tls/server.crt
        - Host: 192.168.181.137
          Port: 7050
          ClientTLSCert: ../organizations/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt
          ServerTLSCert: ../organizations/ordererOrganizations/example.com/orderers/orderer2.example.com/tls/server.crt
        - Host: 192.168.181.138
          Port: 7050
          ClientTLSCert: ../organizations/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt
          ServerTLSCert: ../organizations/ordererOrganizations/example.com/orderers/orderer3.example.com/tls/server.crt

    # Batch Timeout: The amount of time to wait before creating a batch
    BatchTimeout: 2s

    # Batch Size: Controls the number of messages batched into a block
    BatchSize:

        # Max Message Count: The maximum number of messages to permit in a batch
        MaxMessageCount: 10

        # Absolute Max Bytes: The absolute maximum number of bytes allowed for
        # the serialized messages in a batch.
        AbsoluteMaxBytes: 99 MB

        # Preferred Max Bytes: The preferred maximum number of bytes allowed for
        # the serialized messages in a batch. A message larger than the preferred
        # max bytes will result in a batch larger than preferred max bytes.
        PreferredMaxBytes: 512 KB

    # Organizations is the list of orgs which are defined as participants on
    # the orderer side of the network
    Organizations:

    # Policies defines the set of policies at this level of the config tree
    # For Orderer policies, their canonical path is
    #   /Channel/Orderer/<PolicyName>
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
        # BlockValidation specifies what signatures must be included in the block
        # from the orderer for the peer to validate it.
        BlockValidation:
            Type: ImplicitMeta
            Rule: "ANY Writers"

Channel: &ChannelDefaults
    # Policies defines the set of policies at this level of the config tree
    # For Channel policies, their canonical path is
    #   /Channel/<PolicyName>
    Policies:
        # Who may invoke the 'Deliver' API
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        # Who may invoke the 'Broadcast' API
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        # By default, who may modify elements at this config level
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"

    # Capabilities describes the channel level capabilities, see the
    # dedicated Capabilities section elsewhere in this file for a full
    # description
    Capabilities:
        <<: *ChannelCapabilities

Profiles:

    ClusterOrgsApplicationGenesis:
        <<: *ChannelDefaults
        Orderer:
            <<: *OrdererDefaults
            Organizations:
                - *OrdererOrg
            Capabilities:
                <<: *OrdererCapabilities
        Application:
            <<: *ApplicationDefaults
            Organizations:
                - *Org1
                - *Org2
            Capabilities:
                <<: *ApplicationCapabilities
```
