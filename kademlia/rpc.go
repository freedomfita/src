package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.


// PING
type Ping struct {
    MsgID ID
}

type Pong struct {
    MsgID ID
}

func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
    // This one's a freebie.
    pong.MsgID = CopyID(ping.MsgID)
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
    Nodes []Contact
    Err error
}

func (k *Kademlia) FindNode(req *FindNodeRequest, res *FindNodeResult) error {
    // TODO: Implement.
    //pseudo
    //populate res.Nodes (array of FoundNodes)
    res.Nodes = k.Find_Closest(req.NodeID, 20) //no idea if that should be 20 or not
    //May need to change this, haven't tested it, but we could have to add each entry to
    //res.Nodes in a for loop after returning the array of Contacts, but we'll see
    return nil
}

func (k *Kademlia) Find_Closest(req_id ID, count int) []Contact{
	b_num := req_id.Xor(k.ThisContact.NodeID).PrefixLen() //get bucket number
	b := k.K_Buckets[b_num] //get corresponding bucket
	nodes := make([]Contact, count)  //make node array
	j := 0
	for i:=0;i<len(b) && i<count;i++{ //we copy all contacts from closest bucket
		nodes[i] = *b[i]
		j++
	}
	//then if there is still room, we add neighboring buckets' contacts
	for i:=1; (b_num-i >= 0 || b_num+i < 160) && j<count; i++{
		if b_num-i >= 0{ //copy bucket below
			b = k.K_Buckets[b_num - i]
			for c:=0; j<count && c<len(b);c++{
				nodes[j] = *b[c]
				j++
			}
		}
		if b_num+i < 160{ //copy bucket above
			b = k.K_Buckets[b_num + i]
			for c:=0; j<count && c<len(b);c++{
				nodes[j] = *b[c]
				j++
			}
		}
	}
	//Once full we need to sort. I'm being lazy and saving this for later
	nodes = Sort_Contacts(nodes)
	return nodes
}

//make sort function if above sorting isn't valid
func Sort_Contacts(arr []Contact) []Contact {
	///sorted_arr := make([]Contact, len(arr))
	for i:=1;i<len(arr);i++{
		insert_val := arr[i]
		hole_pos := i
		
		for ;hole_pos>0 && insert_val.NodeID.Less(arr[i-1].NodeID);{
			arr[hole_pos] = arr[hole_pos-1]
			hole_pos--
		}
		arr[hole_pos] = insert_val
	}
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
    if !k.ThisContact.NodeID.Equals(req.Key){//don't know exactly what we're checking
    	
    } //if we don't have value it acts like find node. Don't know if that means that we just
      //can call FindNode with it, because that would be easy, but FindValueResult isn't currently a triple
      
      
      
    return nil
}

