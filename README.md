# ETH2.0 P2P Tests

## Overview
The focal point of the following tests is to analyze the performance of various messaging patterns, including 
the Libp2p implementation of gossipsub and floodsub, as well as Artemis' current implementation of the Plumtree algorithm.

## Test Methodology

This work started in early 2018 in collaboration with the [Eth Research team](https://github.com/ethresearch/sharding-p2p-poc).

Please reference [this document](https://notes.ethereum.org/s/ByYhlJBs7) for further details pertaining to this initial test plan.

This document is a work in progress and will be updated accordingly as we progress with these initiatives. 

## Test Utility

* Number of nodes: <100 (more if necessary)
* Bandwidth: 1G (standard, up to 10G if necessary) can be configured and assigned to each individual node.
* Latency: >1 second of network latency can be applied to each node’s individual link.
* Data aggregation and visualization

## Test Scenarios

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

## Need to Define

* Configuration specifications for relevant test cases to create templates which allow for quicker test automation.
* Code which should be tested.
* Preliminary testing methodology should be established based on everyone's input.
  * We can make adjustments to this methodology based on the results of each test case.
  * It's generally best (in my experience) to create a high-level overview which provides a more granular definition of the first three test cases and then make adjustments to each subsequent test series based on the results of those three.

## Other Notes

* This document acts as a high-level overview to communicate potential test scenarios based on our initial assumptions of the existing codebase. It is meant to act as a starting point to guide the development of a preliminary test series.
* Although tests will be run locally within our lab, access to the test network can be granted to appropriate third-parties for the sake of due diligence, validation, or other purposes as deemed necessary.
* Network statistics and performance data dashboard can be assigned a public IP to allow for public access.
* Raw and formatted data will be shared within appropriate repos.
* Please voice any suggestions, comments, or concerns in this thread and feel free to contact me on Gitter.
