# Gossipsub tester

This program takes on the libp2p daemon with gossipsub to create a simple application to send messages to all other peers.

## packages
The packages will require golang version 1.11+. The application will not be able to compile with a older version.

# Introduction

This is a simple tester that runs a gossip implementation, is able to publish messages passed on it to it via a RPC call, and can log received messages to a log file.

## Startup

The program takes a few arguments to start:
* the destination log file
* The network interface to bind to
* The port on which it will listen for incoming gossip messages
* The port on which it will listen for RPC messages

## Execution

The program runs until it is interrupted by a signal from terminal, such as interrupt.

## Publishing messages

The program will publish messages using a POST RPC call under the path /publish.
The program will absorb the payload of the body and post it as is to all other peers.

## Logging messages

Messages are logged using JSON for easy ingestion.
The format of the JSON object is:
{"timestamp":<ISO 8601 timestamp>, "value":<payload received>}

This format may evolve over time.

# License

Copyright 2019 Antoine Toulme

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
