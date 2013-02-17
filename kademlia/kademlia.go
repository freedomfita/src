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
  ThisContact *Contact
  K_Buckets []Bucket
  Data map[ID]([]byte)
  bucket_sizes []int
}

func NewKademlia() *Kademlia {
    // Assign this node a semi-random ID and prepare other state here.
    kadem := new(Kademlia)
    kadem.ThisContact = new(Contact)
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
    	if bucket[i] != nil {
        	if bucket[i].NodeID == id {
            	return bucket[i]
        	}
        }
    }
    return nil
}

func (kadem *Kademlia) GetBucket(dist ID) (Bucket,int) {
    // if PrefixLen == x, then 2^(160-(x+1)) <= ID < 2^(160-x), so the bucket # is (159-x)
    bucketNum := NumBuckets - (dist.PrefixLen() + 1)
    return kadem.K_Buckets[bucketNum], bucketNum
}

func (k *Kademlia) Next_Open_Spot(b_num int) {
	b := k.K_Buckets[b_num]
	open_spot := -1
	b_len := len(b)
	if b[0] ==nil{
		return
	}
	for i:=1;i<b_len;i++{
		if b[i]==nil{
			open_spot=i
			break
		}
	}
	//if open_spot==-1, list is full
	//so pop last entry(which is really the first) and shift list one spot to the right
	if open_spot==-1{
		b[b_len-1] = nil //make last entry nil
		//shift list
		for i:=b_len-2;i>0;i--{
			b[i+1] = b[i]
		}
		b[0] = nil
		
	}
	//else, shift list over one, with last entry at open_spot-1
	//shift 0 to openspot -1 to 1 to openspot
	for i:=open_spot-1;i>0;i--{
		b[i+1] = b[i]
	}
	b[0]=nil
	return
}
/*
[a][ ][ ]
[b][a][ ]
[c][b][a]
[c][b][ ]
[ ][c][b]
 ^
[d][c][b]

*/
func (kadem *Kademlia) AddContactToBuckets(node *Contact) int {

    _, idx := kadem.GetBucket(node.NodeID)
    //frees up first
    kadem.Next_Open_Spot(idx)
    kadem.K_Buckets[idx][0] = node
    
    return 0
}