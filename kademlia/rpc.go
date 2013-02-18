package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"errors"
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
	res.MsgID = CopyID(req.MsgID)
	k.Data[req.Key] = req.Value
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

func (f *FoundNode) ToContactPtr() *Contact {
	c := new(Contact)
	c.NodeID = CopyID(f.NodeID)
	c.Port = f.Port
	c.IPAddr = f.IPAddr
	return c
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
    res.Nodes = bucket_to_FoundNodeArr(k.find_closest(req.NodeID, 20)) //no idea if that should be 20 or not
    //May need to change this, haven't tested it, but we could have to add each entry to
    //res.Nodes in a for loop after returning the array of Contacts, but we'll see
    return nil
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
    		res.Value = copyData(val)
    		res.Nodes = nil
    		res.Err = nil
    		return nil
    	}
    }
    res.Value = nil
    res.Nodes = bucket_to_FoundNodeArr(k.find_closest(k.ThisContact.NodeID,20))
    res.Err = errors.New("Value not found")
    return nil
}