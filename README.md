# ETH2.0 P2P Tests 
[![Gitter](https://badges.gitter.im/whiteblock-io/community.svg)](https://gitter.im/whiteblock-io/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Overview

The following tests are designed to observe and measure the performance of various protocols responsible for the dissemination of messages within the network. Within these tests,we will analyzing the following protocols: 

    * Gossipsub (libp2p)
    * Floodsub (libp2p)
    * Plumtree (Apache Tuweni)

This work began in early 2018 in collaboration with between Whiteblock and the [ETH Research team](https://github.com/ethresearch/sharding-p2p-poc).

Please reference [this document](https://notes.ethereum.org/s/ByYhlJBs7) for further details pertaining to this initial test plan.

This document is a work in progress and will be updated accordingly as we progress with these initiatives. 

## Test Methodology

Tests will be conducted using Whiteblock's [Genesis](www.github.com/whiteblock/genesis) testing framework in accordance with the proposed scope of work outlined within this document. 
Libp2p will be tested using our own custom client located within this repo. The Plumtree implementation is natively supported within the Genesis framework.

## Network Topology 
![Network Topology](/topology.png)

Within most topologies, peering with every other node within the network is ineffective and likely impossible. Within a live, global network, we can assume that
nodes will be organized according to the topology illustrated within the above diagram. 

For example, a (cluster specific) node within Cluster 1 may be peered with N number of nodes within its own cluster, however, based on proximity, certain nodes on the edge of this cluster may also be peered with nodes within Cluster 2 (inter cluster nodes). If Node X within Cluster 1 would like to transmit a message to Node Y within Cluster 4, these messages must propogate through each consecutive cluster in order to reach its destination. 

While this topology may present an oversimplification, within most cases, we can expect the results to be reflective of real-world performance. As we establish an appropriate dataset that is indicative of baseline performance, we can develop additional test series' and cases for future test phases. 

Since peer discovery is outside the scope of work for this test phase, peering within the client implementation presented within this repository is handled statically. 

## Client Behavior
Nodes within the network will be running the client application included within this repo. This client application is responsible for constructing or relaying messages, interpreting these messages, outputting this data to a log in accordance with the defined message struct, and then relaying those messages according to the rules defined by the pusub router (floodsub, gossipsub)

## Test Procedure

Per test case:
1. Build network
2. Provision nodes
3. Configure network conditions between nodes according to specified test case
4. Configure actions and behavior between nodes according to specified test case
5. Log output from each node is aggregated as raw data
6. Raw data is parsed 
8. Parsed data is pushed to appropriate repo
9. Reset environment

## Message Struct

The message struct defines the data which is written to the node's log. These logs are aggregated at runtime to be parsed after each test series is complete. The included data is as follows:

* Timestamp of message received
* Message type 
* Message origin (sender)
* Message destination (receiver)
* Last relaying node (node that sent to you)
* Message value 
* Message nonce (chronology of the sent message)
* Message size 
* MessageID - unique string associated with that message

## Performance Tests

The following tables define each test series within this test phase. A test series focuses on observing and documenting the effects of certain conditions on performance. Each test series is comprised of three separate test cases which define the variable to be tested. 

It is important to note that each test series may yield unexpected results which may influence the configuration of subsequential test cases or series. Accounting for this notion, this test plan should be considered a living document subject to change. Based on the results of this test phase, a consecutive test phase may be developed. 

### Series 1: Control

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 50          | 50          | 50          |
| Rx Nodes         | 50          | 50          | 50          |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 1Gb         | 1Gb         | 1Gb         |
| Network Latency  | 0ms         | 0ms         | 0ms         |
| Packet Loss      | 0%          | 0%          | 0%          |

### Series 2: Message Size

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 50          | 50          | 50          |
| Rx Nodes         | 50          | 50          | 50          |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 500B        | 500KB       | 1MB         |
| Bandwidth        | 1Gb         | 1Gb         | 1Gb         |
| Network Latency  | 0ms         | 0ms         | 0ms         |
| Packet Loss      | 0%          | 0%          | 0%          |


### Series 3: Tx/Rx Nodes

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 10          | 40          | 90          |
| Rx Nodes         | 90          | 60          | 10          |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 1Gb         | 1Gb         | 1Gb         |
| Network Latency  | 0ms         | 0ms         | 0ms         |
| Packet Loss      | 0%          | 0%          | 0%          |


### Series 4: Bandwidth

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 50          | 50          | 50          |
| Rx Nodes         | 50          | 50          | 50          |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 50Mb        | 250Mb       | 750Mb       |
| Network Latency  | 0ms         | 0ms         | 0ms         |
| Packet Loss      | 0%          | 0%          | 0%          |


### Series 5: Network Latency

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 50          | 50          | 50          |
| Rx Nodes         | 50          | 50          | 50          |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 1Gb         | 1Gb         | 1Gb         |
| Network Latency  | 10ms        | 100ms       | 500ms       |
| Packet Loss      | 0%          | 0%          | 0%          |

### Series 6: Packet Loss

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 50          | 50          | 50          |
| Rx Nodes         | 50          | 50          | 50          |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 1Gb         | 1Gb         | 1Gb         |
| Network Latency  | 0ms         | 0ms         | 0ms         |
| Packet Loss      | 0.01%       | 0.1%        | 1%          |

### Series 7: Stress Test

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 50          | 50          | 50          |
| Rx Nodes         | 50          | 50          | 50          |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 10MB        | 10MB        | 10MB        |
| Network Latency  | 150ms       | 150ms       | 150ms       |
| Packet Loss      | 0.1%        | 0.1%        | 0.1%        |



## Future Test Scenarios

* Observe and measure performance under the presence of various network conditions.
  * Latency between nodes:
    * What is the maximum amount of network latency each individual node can tolerate before performance begins to degrade?
    * What are the security implications of high degrees of latency?
    * Are there any other unforeseen issues which may arise from network conditions for which we can’t immediately accommodate?
  * Intermittent blackout conditions
  * High degrees of packet loss
  * Bandwidth constraints (various bandwidth sizes)
* Introduce new nodes to network:
  * Add/remove nodes at random.
  * Add/remove nodes at set intervals.
  * Introduce a high volume of nodes simultaneously.
* Partition tolerance
  * Prevent segments of nodes from communicating with one another.
* Measure the performance of sending/receiving messages within set time periods and repeat for N epochs.
* Observe the process of introducing and removing nodes from the network.
