package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"errors"
	"fmt"
	"sort"
)

// PING
type Ping struct {
    MsgID ID
}

type Pong struct {
    MsgID ID
    Sender Contact
}

func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
    // This one's a freebie.
    pong.MsgID = CopyID(ping.MsgID)
    pong.Sender = *k.ThisContact
    return nil
}


// STORE
type StoreRequest struct {
    MsgID ID
    Key ID
    Value []byte
}

type StoreResult struct {
    MsgID ID
    Err error
}

func (k *Kademlia) Store(req StoreRequest, res *StoreResult) error {
	k.Data[req.Key] = req.Value
	res.MsgID = CopyID(req.MsgID)
    return nil
}


// FIND_NODE
type FindNodeRequest struct {
    MsgID ID
    NodeID ID
}

type FoundNode struct {
    IPAddr string
    Port uint16
    NodeID ID
}

//****NOTE - I changed Nodes from FoundNode to Contact,just so it would be uniform, since FoundNode and Contact are identical
//****If this is bad we can change it back, would just require a little extra code to make everything
//****work smoothly
type FindNodeResult struct {
    MsgID ID
    //Nodes []FoundNode
    Nodes []FoundNode
    Err error
}

func (k *Kademlia) FindNode(req *FindNodeRequest, res *FindNodeResult) error {
    // TODO: Implement.
    //pseudo
    //populate res.Nodes (array of FoundNodes)
    res.MsgID = CopyID(req.MsgID)
    res.Nodes = Bucket_to_FoundNode(k.Find_Closest(req.NodeID, 20)) //no idea if that should be 20 or not
    //May need to change this, haven't tested it, but we could have to add each entry to
    //res.Nodes in a for loop after returning the array of Contacts, but we'll see
    return nil
}

func (k *Kademlia) Find_Closest(req_id ID, count int) []*Contact{
	b_num := req_id.Xor(k.ThisContact.NodeID).PrefixLen() //get bucket number
	b_num--
	fmt.Printf("tried to access bucket %d\n",b_num)
	b := k.K_Buckets[b_num] //get corresponding bucket
	nodes := make([]*Contact, count)  //make node array
	j := 0
	for i:=0;i<len(b) && i<count;i++{ //we copy all contacts from closest bucket
		if b[i] == nil{
			continue
		}
		nodes[i] = b[i]
		j++
	}
	//then if there is still room, we add neighboring buckets' contacts
	for i:=1; (b_num-i >= 0 || b_num+i < 160) && j<count; i++{
		if b_num-i >= 0{ //copy bucket below
			b = k.K_Buckets[b_num - i]
			for c:=0; j<count && c<len(b);c++{
				if b[c] == nil{
					continue
				}
				nodes[j] = b[c]
				j++
			}
		}
		if b_num+i < 160{ //copy bucket above
			b = k.K_Buckets[b_num + i]
			for c:=0; j<count && c<len(b);c++{
				if b[c] == nil{
					continue
				}
				nodes[j] = b[c]
				j++
			}
		}
	}
	//Once full we need to sort. I'm being lazy and saving this for later
	nodes = Sort_Contacts(nodes)
	return nodes
}

// sort function for buckets
func Sort_Contacts(arr Bucket) Bucket {
	sort.Sort(BucketSort_ByNodeID{arr})
	return arr//sorted_arr
}

// FIND_VALUE
type FindValueRequest struct {
    MsgID ID
    Key ID
}

// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
    MsgID ID
    Value []byte
    Nodes []FoundNode
    Err error
}

func (k *Kademlia) FindValue(req FindValueRequest, res *FindValueResult) error {
    // TODO: Implement.
    for key,val := range k.Data {
    	if key.Equals(req.Key) {
    		res.MsgID = CopyID(req.MsgID)
    		// can we just do this here, or need to Copy(val) ?
    		res.Value = val
    		res.Nodes = nil
    		res.Err = nil
    		return nil
    	}
    }
    res.Value = nil
    res.Nodes = Bucket_to_FoundNode(k.Find_Closest(k.ThisContact.NodeID,20))
    res.Err = errors.New("Value not found")
    return nil
}

func Bucket_to_FoundNode(bucket Bucket) []FoundNode {
	b := make([]FoundNode,len(bucket))
	j := 0
	for i := 0; i < len(bucket); i++ {
		if bucket[i] != nil {
			b[j].IPAddr = bucket[i].IPAddr
			b[j].NodeID = bucket[i].NodeID
			b[j].Port = bucket[i].Port
			j++
		}
	}
	return b
}

func FoundNode_to_Bucket(foundNodes []FoundNode) Bucket {
	b := make(Bucket,len(foundNodes))
	for i := 0; i < len(foundNodes); i++ {
			b[i] = new(Contact)
			b[i].IPAddr = foundNodes[i].IPAddr
			b[i].NodeID = foundNodes[i].NodeID
			b[i].Port = foundNodes[i].Port
	}
	return b
}