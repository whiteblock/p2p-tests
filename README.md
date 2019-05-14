# ETH2.0 P2P Tests 
## Test Phase v1.0 

## Overview

The following tests are designed to observe and measure the performance of various protocols 
responsible for the dissemination of messages within the network. Within these tests, we will 
analyzing the following protocols: 
    - Gossipsub (libp2p)
    - Floodsub (libp2p)
    - Plumtree (Artemis)

This work began in early 2018 in collaboration with between Whiteblock and the [ETH Research team](https://github.com/ethresearch/sharding-p2p-poc).

Please reference [this document](https://notes.ethereum.org/s/ByYhlJBs7) for further details pertaining to this initial test plan.

This document is a work in progress and will be updated accordingly as we progress with these initiatives. 

## Test Methodology

Tests will be conducted using Whiteblock's [Genesis(www.github.com/whiteblock/genesis)] testing framework in accordance with the proposed scope of work outlined within this document. 
Libp2p will be tested using our own custom client located within this repo. The Plumtree implementation is natively supported within the Genesis framework.

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

## Performance Metrics

| Value | Description | 
| -------- | -------- | 
| Subscription Time | The length of time it takes for a node to subscribe to a topic, or in otherwords, join a shard, and begin receiving and broadcasting messages pertaining to that topic. |
| Discovery Time | The length of time it takes for a node to become aware of its peers within their subscribed shard. | 
| Message Propagation Time | (Broadcast time) The length of time it takes for a message, once broadcast, to be received by a majority (99%) of peers within the shard.|

## Performance Tests

The following tables define each test series within this test phase. A test series focuses on observing and documenting the effects of certain conditions on performance. Each test series is comprised of three separate test cases which define the variable to be tested. 

It is important to note that each test series may yield unexpected results which may influence the configuration of subsequential test cases or series. Accounting for this notion, this test plan should be considered a living document subject to change. Based on the results of this test phase, a consecutive test phase may be developed. 

### Series 1: Control

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 50          | 50          | 50          |
| Rx Nodes         | 50          | 50          | 50          |
| Static Nodes     | 0           | 0           | 0           |
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
| Static Nodes     | 0           | 0           | 0           |
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
| Static Nodes     | 0           | 0           | 0           |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 1Gb         | 1Gb         | 1Gb         |
| Network Latency  | 0ms         | 0ms         | 0ms         |
| Packet Loss      | 0%          | 0%          | 0%          |

### Series 4: Static Nodes

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 10          | 30          | 40          |
| Rx Nodes         | 10          | 30          | 40          |
| Static Nodes     | 80          | 40          | 20          |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 1Gb         | 1Gb         | 1Gb         |
| Network Latency  | 0ms         | 0ms         | 0ms         |
| Packet Loss      | 0%          | 0%          | 0%          |


### Series 5: Total Nodes

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 20          | 40          | 80          |
| Tx Nodes         | 100         | 100         | 100         |
| Rx Nodes         | 100         | 100         | 100         |
| Static Nodes     | 0           | 0           | 0           |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 1Gb         | 1Gb         | 1Gb         |
| Network Latency  | 0ms         | 0ms         | 0ms         |
| Packet Loss      | 0%          | 0%          | 0%          |


### Series 6: Bandwidth

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 50          | 50          | 50          |
| Rx Nodes         | 50          | 50          | 50          |
| Static Nodes     | 0           | 0           | 0           |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 50Mb        | 250Mb       | 750Mb       |
| Network Latency  | 0ms         | 0ms         | 0ms         |
| Packet Loss      | 0%          | 0%          | 0%          |


### Series 7: Network Latency

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 50          | 50          | 50          |
| Rx Nodes         | 50          | 50          | 50          |
| Static Nodes     | 0           | 0           | 0           |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 1Gb         | 1Gb         | 1Gb         |
| Network Latency  | 10ms        | 100ms       | 500ms       |
| Packet Loss      | 0%          | 0%          | 0%          |

### Series 8: Packet Loss

| Variable         | Test Case A | Test Case B | Test Case C |
|------------------|------------:|------------:|------------:|
| Total Nodes      | 100         | 100         | 100         |
| Tx Nodes         | 50          | 50          | 50          |
| Rx Nodes         | 50          | 50          | 50          |
| Static Nodes     | 0           | 0           | 0           |
| Peers/Node       | 10          | 10          | 10          |
| Message Size     | 200B        | 200B        | 200B        |
| Bandwidth        | 1Gb         | 1Gb         | 1Gb         |
| Network Latency  | 0ms         | 0ms         | 0ms         |
| Packet Loss      | 0.01%       | 0.1%        | 1%          |

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
