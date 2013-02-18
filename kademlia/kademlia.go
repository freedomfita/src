package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
  //"net"
  "fmt"
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
}

func NewKademlia() *Kademlia {
    // Assign this node a semi-random ID and prepare other state here.
    kadem := new(Kademlia)
    kadem.ThisContact = new(Contact)
    kadem.ThisContact.NodeID = NewRandomID()
    kadem.K_Buckets = make([]Bucket,160)
    for i := 0; i < 160; i++ {
	kadem.K_Buckets[i] = make(Bucket,20)
    }
    kadem.Data = make(map[ID]([]byte))
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

func Next_Open_Spot(b Bucket) {
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

func (k *Kademlia) Next_Open_Spot(b_num int) {
	b := k.K_Buckets[b_num]
	open_spot := -1
	b_len := len(b)
	fmt.Printf("Looking for next open spot in bucket %v\n", b_num)
	if b[0] ==nil{
		return
	}
	for i:=1;i<b_len;i++{
		if b[i]==nil{
			open_spot=i
			fmt.Printf("Open spot at %v\n",i)
			break
		}
	}
	//if open_spot==-1, list is full
	//so pop last entry(which is really the first) and shift list one spot to the right
	if open_spot==-1{
		fmt.Printf("Popping %v\n", b[b_len-1])
		b[b_len-1] = nil //make last entry nil
		//shift list
		for i:=b_len-2;i>0;i--{
			b[i+1] = b[i]
		}
		b[0] = nil
		
	} else{
	//else, shift list over one, with last entry at open_spot-1
	//shift 0 to openspot -1 to 1 to openspot
		for i:=open_spot;i>0;i--{
			fmt.Printf("moving %v to %v\n",i-1,i)
			fmt.Printf("Values: %v\n %v\n",b[i],b[i-1])
			b[i] = b[i-1]
		}
		b[0]=nil
		return
	}
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

// interface to allow for sorting within buckets
func (s Bucket) Len() int      { return len(s) }
func (s Bucket) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// BucketSort_ByNodeID implements sort.Interface by providing Less and using the Len and
// Swap methods of the embedded Organs value.
type BucketSort_ByNodeID struct{ Bucket }

func (s BucketSort_ByNodeID) Less(i, j int) bool { 
if s.Bucket[i] == nil {
	return false //nil's go at the end
} else if s.Bucket[j] == nil {
	return true
}
	return s.Bucket[i].NodeID.Less(s.Bucket[j].NodeID) 
}

// execute the ping RPC given either a nodeID or a host:port pair
// returns the contact info of the ping'ed node
func ping(nodeToPing string) *kademlia.Contact {

	// if the argument is not of the form host:port, assume that it is a nodeID and 
	// look up the corresponding host/port pair
	if len(strings.Split(nodeToPing, ":")) != 2 {
		id, err := kademlia.FromString(nodeToPing)
		if err != nil {
			log.Printf("Error: Could not convert from string to nodeID (%e)\n",err)
			return nil
		}
		contact := kademlia.LookupContact(thisNode,id)
		if contact == nil {
			log.Printf("Error: Could not find node with NodeID %s\n",nodeToPing)
			return nil
		}
		nodeToPing = get_host_port(contact)
		//log.Printf("Host/Port: %s\n",nodeToPing)
	}
    client, err := rpc.DialHTTP("tcp", nodeToPing)
    if err != nil {
		log.Fatal("DialHTTP: ", err)
    }
    ping := new(kademlia.Ping)
    ping.MsgID = kademlia.NewRandomID()
    var pong kademlia.Pong
    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
	log.Fatal("Call: ", err)
    }

    log.Printf("ping msgID: %s\n", ping.MsgID.AsString())
    log.Printf("pong msgID: %s\n", pong.MsgID.AsString())
    //fmt.Printf("%s\n",pong.Sender.NodeID.AsString())
	return &pong.Sender
}

func store(hostAndPort string, key kademlia.ID, data []byte) int {
	client, err := rpc.DialHTTP("tcp", hostAndPort)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
    }
    /*
    // check if the receiver has the value already
    req := new(kademlia.FindValueRequest)
    req.MsgID = kademlia.NewRandomID()
    req.Key = key
    
    var res kademlia.FindValueResult
	err = client.Call("Kademlia.Store", req, &res)
    if err != nil {
    	log.Fatal("Call: ", err)
    }
    // if res.Err == nil, then the node has the value for the key already
    if res.Err != nil {
    */
    /* initialize the Request and Result structs */
    	req := new(kademlia.StoreRequest)
    	req.MsgID = kademlia.NewRandomID()
    	req.Key = key
    	req.Value = data
    	
    	var res kademlia.StoreResult
		err = client.Call("Kademlia.Store", req, &res)
    	if err != nil {
			log.Fatal("Call: ", err)
    	}
    /* } else {
    	fmt.Printf("Node already had value.\n")
    } */
	return 1
}
func find_node(key kademlia.ID) int {
	bucket,_ := thisNode.GetBucket(key)
	nodes := bucket.FindNode(key)
	fmt.Println(nodes)

	return 0
}
func get_local_value(key kademlia.ID) int {
    if thisNode.Data[key] != nil {
		log.Printf("OK: %v\n", thisNode.Data[key])
    } else {
    	log.Printf("ERR\n")
    }
    return 0

}
func get_node_id() int {
	log.Printf("Node ID of this node: %s\n",thisNode.ThisContact.NodeID.AsString())
	log.Printf("IP/Port: %v %v\n",thisNode.ThisContact.IPAddr,thisNode.ThisContact.Port)
	return 0
}

/*
Perform the iterativeStore operation and then print the ID of the node that
received the final STORE operation.
*/
func iterativeStore(key kademlia.ID, value []byte) int {

	prevDistance := key.Xor(thisNode.ThisContact.NodeID)
	
	//var closestNode kademlia.FoundNode
	closestNode := thisNode.ThisContact
	
	hostPortStr := get_host_port(thisNode.ThisContact)
	
	//closestnode may want to be its own function that we call from FindNode, or at least
	//that code should be in FindNode, since we need to populate res.Nodes with more than one bucket
	for true {
		log.Printf("%s\n",hostPortStr)
		client, err := rpc.DialHTTP("tcp", hostPortStr)
		if err != nil {
			log.Printf("1\n")
			log.Fatal("DialHTTP: ", err)
		}
		req := new(kademlia.FindNodeRequest)
		req.MsgID = kademlia.NewRandomID()
		req.NodeID = key
	
		var res kademlia.FindNodeResult
		//if FindNode works, all of the closest nodes should be in res.
		err = client.Call("Kademlia.FindNode", req, &res)
    		if err != nil {
			log.Fatal("Call: ", err)
    		}
    		// obviously we need to do something with the array here, not just take the first element
    		log.Printf("Node 0: %v\n",res.Nodes[0])
    		nextClosestNode, dist := res.Nodes[0], key.Xor(res.Nodes[0].NodeID)
    		for i:= 0; i < len(res.Nodes); i++ {
    			if res.Nodes[i].NodeID.Xor(key).Less(dist) {
    				dist = res.Nodes[i].NodeID.Xor(key)
    				nextClosestNode = res.Nodes[i]
    			}
    		}
    		curDistance := key.Xor(nextClosestNode.NodeID)
    	
    		if !curDistance.Less(prevDistance) {
    			break
    		} else {
    			closestNode = nextClosestNode.ToContactPtr()
    		}
    		hostPortStr = get_host_port(closestNode)
		}
	hostPortStr = get_host_port(closestNode)
	store(hostPortStr, key, value)
	log.Printf("NodeID receiving STORE operation: %d\n",closestNode.NodeID)
	return 1
}

/*
Print a list of ≤ k closest nodes and print their IDs. You should collect
the IDs in a slice and print that.
*/
func iterativeFindNode(id kademlia.ID) kademlia.Bucket { 
	//Get 20 closest nodes from us.
	req := new(kademlia.FindNodeRequest)
	req.NodeID = id
	req.MsgID = kademlia.NewRandomID()
	var k_res kademlia.FindNodeResult
	fmt.Printf("In Iterative Find Node, before finding initial nodes closest to NodeID %v\n",id)
	err := thisNode.FindNode(req,&k_res)
	
	k_closest := kademlia.FoundNode_to_Bucket(k_res.Nodes)
	fmt.Printf("In Iterative Find Node, after finding initial nodes closest to NodeID\n")
	if err != nil {
		log.Fatal("Call: ", err)
    	}
	//initialize array to hold all 20^2 contacts, which we'll sort later
	big_arr := make(kademlia.Bucket, 400)

	for i :=0;i<len(k_closest);i++{
		//find 20 closest for each node.
		if k_closest[i] == nil {
			
		} else if k_closest[i].Port != 0 {
			hostPortStr := get_host_port(k_closest[i])
			log.Printf("Host/Port: %s\n",hostPortStr)
			client, err := rpc.DialHTTP("tcp", hostPortStr)
			if err != nil {
				log.Printf("2\n")
				log.Fatal("DialHTTP: ", err)
			}
			req := new(kademlia.FindNodeRequest)
			req.MsgID = kademlia.NewRandomID()
			req.NodeID = id
			
			var res kademlia.FindNodeResult
			//if FindNode works, all of the closest nodes should be in res.
			err = client.Call("Kademlia.FindNode", req, &res)
			log.Printf("NODES: %v\n",res.Nodes)
			if err != nil {
				log.Fatal("Call: ", err)
    			}
    			offset:= 20 * i
    			resBucket := kademlia.FoundNode_to_Bucket(res.Nodes)
			for j := 0; j<len(resBucket); j++{
				big_arr[j+offset] = resBucket[j]
			}
		}
		
	}
	fmt.Printf("Finished IterativeFindNode and returning array of contacts\n")
	// print slice of <= k closest NodeIDs
	//fmt.Printf("%v\n",kademlia.Sort_Contacts(big_arr)[:20])
	return (kademlia.Sort_Contacts(big_arr))
}

/*
printf("%v %v\n", ID, value), where ID refers to the node that finally
returned the value. If you do not find a value, print "ERR".
*/
func iterativeFindValue(key kademlia.ID) int {

	// check if this node has the value
	if thisNode.Data[key] {
		log.Printf("%v %v\n", thisNode.ThisContact.IPAddr, thisNode.Data[key])
		return 0
	}

	const alpha = 3
	
	contacted_nodes := make(kademlia.Bucket,1600)
	shortlist := make(kademlia.Bucket,20)
	shortlist_size := 0
	// The search begins by selecting alpha contacts from the non-empty k-bucket closest to the 
	// bucket appropriate to the key being searched on.
	_, bucket_num := thisNode.GetBucket(key.Xor(thisNode.ThisContact.NodeID))
	for i := 0; i < 20 && shortlist_size < alpha; i++ {
		if thisNode.K_Buckets[bucket_num][i] != nil {
			shortlist[shortlist_size] = thisNode.K_Buckets[bucket_num][i]
			shortlist_size++
		}
	}
	// If there are fewer than alpha contacts in that bucket, contacts are selected from other buckets.
	for b_idx := 0; b_idx < kademlia.NumBuckets && shortlist_size < alpha; b_idx++ {
		if b_idx != bucket_num {
			for i := 0; i < 20 && shortlist_size < alpha; i++ {
				if thisNode.K_Buckets[bucket_num][i] != nil {
					shortlist[shortlist_size] = thisNode.K_Buckets[bucket_num][i]
					shortlist_size++
					if shortlist_size >= alpha {
						break
					}
				}
			}
		}
	}
	
	for true {
	// The node then sends parallel, asynchronous FIND_* RPCs to the alpha contacts in the 
	// shortlist. Each contact, if it is live, should normally return k triples. If any of the 
	// alpha contacts fails to reply, it is removed from the shortlist, at least temporarily.
	
		// TODO: this isn't parallel yet.
		new_shortlist := make(kademlia.Bucket,400)
		for i := 0; i < len(shortlist); i++ {
			kademlia.Next_Open_Spot(contacted_nodes)
			contacted_nodes[0] = shortlist[i]
			if shortlist[i] != nil {
				hostPortStr := get_host_port(contacted_nodes[i])
		
				client, err := rpc.DialHTTP("tcp", hostPortStr)
				if err != nil {
					// TODO: this definitely shouldn't be a fatal error, it should just go on to the next node
					log.Printf("3\n")
					log.Printf("DialHTTP: %e\n", err)
				} else {
					req := new(kademlia.FindValueRequest)
					req.MsgID = kademlia.NewRandomID()
					req.Key = key
			
					var res kademlia.FindValueResult
					err = client.Call("Kademlia.FindValue", req, &res)
					if err != nil {
						log.Fatal("Call: ", err)
    				}	
    				// if res.Err is nil, the node contains the value
    				if res.Err == nil {
    					log.Printf("%v %v\n", shortlist[i].IPAddr, res.Value)
    					return 0
    				} else {
    					offset:= 20 * i
    					resBucket := kademlia.FoundNode_to_Bucket(res.Nodes)
						for j := 0; j<len(resBucket); j++{
							new_shortlist[j+offset] = resBucket[j]
						}
    				}
    			}
    		}
    		shortlist_size = 0
    	// assign contacts from new_shortlist to shortlist IF they haven't been contacted already
    		for i := 0; i < len(new_shortlist) && shortlist_size < alpha; i++ {
    		if new_shortlist[i] != nil {
    			if shortlist.FindNode(new_shortlist[i].NodeID) == nil && contacted_nodes.FindNode(new_shortlist[i].NodeID) == nil {
    				shortlist[shortlist_size] = new_shortlist[i]
    				shortlist_size++
    			}
    		}
    	}
    		if shortlist_size == 0 {
    			log.Printf("ERR\n")
    			return 1
    		}
    	}
    }
    return 0
}

func get_host_port(c *kademlia.Contact) string {
	if c == nil {
		return ""
	}
	hostPort := make([]string,2)
	hostPort[0] = c.IPAddr
	hostPort[1] = strconv.FormatUint(uint64(c.Port),10)
	hostPortStr := strings.Join(hostPort, ":")
	return hostPortStr
}