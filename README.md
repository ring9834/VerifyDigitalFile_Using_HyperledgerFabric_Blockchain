# Verifying Digital Files Via Hyperledger Fabric Blockchain
The codes are mainly to provide users with intuitive interfaces to save the hashes of valuable digital files or archichives into Hyperledger Blockchain as proofs as well as fetching proofs from Blockchain. The codes also include font-end ones and back-end ones in which the back-end ones include APIs dealing with requests and responses from font end, accessing the Blockchain through FabricClient framework, and operating the Mssql database.

The following content is for how to configurate Blockchain Nodes including Peer nodes, Orderer nodes and CA nodes as well as how to deploy Smart Contracts to Blockchain and operate(save and fetch) data in Blockchain using Smart Contracts.

Here, for purpose of test, we utilize virtual machines installed Centos7 and Dockers as the infrastructure to operate Blockchain nodes. We use one virtual machine for each Organization where several Peer nodes could be settled down.

## Two Peer nodes with one Orderer node
In this section, we show the codes for configuring and using two Peer nodes and one Orderer node.

