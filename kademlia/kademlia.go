package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
  "net"
  "fmt"
)

type Contact struct {
  NodeID ID
  IPAddr string
  Port uint16
}

// Core Kademlia type. You can put whatever state you want in this.
type Kademlia struct {
  ThisNode Contact
  K_Buckets[][] Contact
}

func NewKademlia() *Kademlia {
    // Assign this node a semi-random ID and prepare other state here.

}

