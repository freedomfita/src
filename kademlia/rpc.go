package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"errors"
  "fmt"
)

// PING
type Ping struct {
    MsgID ID
}

type Pong struct {
    MsgID ID
    Sender FoundNode
}

func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
    // This one's a freebie.
    pong.MsgID = CopyID(ping.MsgID)
    pong.Sender.NodeID = CopyID(k.ThisContact.NodeID)
    pong.Sender.IPAddr = k.ThisContact.IPAddr
    pong.Sender.Port = k.ThisContact.Port
    return nil
}

//LOCKING/RELEASING
type LockRequest struct {
	MsgID ID
	FileID ID
}

type LockResult struct {
	MsgID ID
	is_locked int
}

type UnlockRequest struct {
	MsgID ID
	FileID ID
}

type UnlockResult struct {
	MsgID ID
}

//AUTHENTICATE
type AuthRequest struct {
	MsgID ID
}

type AuthResult struct {
	MsgID ID
	isFriend int //0 on rejection, 1 on success
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
  fmt.Printf("\n")
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
    // print list of nodes
    fmt.Printf("%v\n",res.Nodes)
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
        fmt.Printf("%v\n",res.Value)
    		return nil
    	}
    }
    res.Value = nil
    res.Nodes = bucket_to_FoundNodeArr(k.find_closest(k.ThisContact.NodeID,20))
    fmt.Printf("%v\n",res.Nodes)
    res.Err = errors.New("Value not found")
    return nil
}

type UpdateListenerRequest struct {
    MsgID ID
    FileID ID
    ListenerID ID
}

type UpdateListenerResult struct {
    MsgID ID
    Err error
}

func (k *Kademlia) AddUpdateListener(req UpdateListenerRequest, res *UpdateListenerResult) error {
    res.MsgID = CopyID(req.MsgID)
    fh := ThisNode.FileHeaders[req.FileID]
    fh.UpdateNodes = append(fh.UpdateNodes,req.ListenerID)
    return nil
}