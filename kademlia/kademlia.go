package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
  //"net"
  //"fmt"
)

const NumBuckets = 160

type Bucket []*Contact

type Contact struct {
  NodeID ID
  IPAddr string
  Port uint16
}

// Core Kademlia type. You can put whatever state you want in this.
type Kademlia struct {
  ThisContact Contact
  K_Buckets []Bucket
  Data map[ID]([]byte)
  bucket_sizes []int
}

func NewKademlia() *Kademlia {
    // Assign this node a semi-random ID and prepare other state here.
kadem := new(Kademlia)
kadem.ThisContact.NodeID = NewRandomID()
kadem.K_Buckets = make([]Bucket,160)
kadem.bucket_sizes = make([]int,160)
for i := 0; i < 160; i++ {
kadem.K_Buckets[i] = make(Bucket,20)
kadem.bucket_sizes[i] = 0
}
return kadem
}

func LookupContact(thisNode *Kademlia, lookupID ID) *Contact {
dist := lookupID.Xor(thisNode.ThisContact.NodeID)
bucket,_ := thisNode.GetBucket(dist)
return bucket.FindNode(lookupID)
}

func (bucket Bucket) FindNode(id ID) *Contact {
for i := 0; i < len(bucket); i++ {
if bucket[i].NodeID == id {
return bucket[i]
}
}
return nil
}

func (kadem *Kademlia) GetBucket(dist ID) (Bucket,int) {
// if PrefixLen == x, then 2^(160-(x+1)) <= ID < 2^(160-x), so the bucket # is (159-x)
bucketNum := NumBuckets - (dist.PrefixLen() + 1)
return kadem.K_Buckets[bucketNum], bucketNum
}

func (kadem *Kademlia) AddContactToBuckets(node *Contact) int {

bucket, idx := kadem.GetBucket(node.NodeID)

bucket[kadem.bucket_sizes[idx]] = node
kadem.bucket_sizes[idx] += 1

/* check if bucket is full; if it is, we need to remove a node first. */
return 0
}