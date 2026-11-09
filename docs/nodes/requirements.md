---
title: Compute and network requirements
description: Review the approximate hardware and network requirements for an Upside L1 node.
slug: /nodes/requirements
---

# Compute and network requirements

Start with these minimums when planning a mainnet validator:


| Resource           | Approximate minimum |
| ------------------ | ------------------- |
| CPU                | 8 cores             |
| RAM                | 16 GB               |
| Storage            | 1 TB                |
| Network connection | 1 Gbps              |


## Network access

Nodes need access to the published P2P peers. Configure the advertised P2P address and firewall to match the network topology. A private validator behind sentries may have different inbound rules from a publicly reachable full node.

Keep management interfaces and consensus key material private. Bind local RPC/API services to loopback unless intentionally operating a public service. Public services need their own access, resource, and rate controls.

Check the release's P2P and RPC settings before configuring public-node ports. The localnet guide supplies separate loopback ports as a guide.
