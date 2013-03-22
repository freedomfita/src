package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
    //"net"
    "fmt"
    "log"
    "net"
    "net/http"
    "net/rpc"
    "strings"
    "strconv"
)

const NumBuckets = 160

var ThisNode *Kademlia

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
    F_Buckets []Bucket
    Data map[ID]([]byte)
    // stuff for p2p
    ShareDir string
    Lock_Acquired ID //this will be the FileID, if no lock, null
    FileHeaders map[ID](FileHeader)
}

func NewKademlia() *Kademlia {
    // Assign this node a semi-random ID and prepare other state here.
    kadem := new(Kademlia)
    kadem.ThisContact = new(Contact)
    kadem.ThisContact.NodeID = NewRandomID()
    kadem.K_Buckets = make([]Bucket,160)
    kadem.F_Buckets = make([]Bucket,160)
    for i := 0; i < 160; i++ {
	kadem.K_Buckets[i] = make(Bucket,20)
	kadem.F_Buckets[i] = make(Bucket,20)
    }
    kadem.Data = make(map[ID]([]byte))
    return kadem
}

func LookupContact(node *Kademlia, lookupID ID) *Contact {
    dist := lookupID.Xor(node.ThisContact.NodeID)
    bucket,_ := node.getBucket(dist)
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

// execute the ping RPC given either a nodeID or a host:port pair
// returns the contact info of the ping'ed node
func Ping2(nodeToPing string) *Contact {

	// if the argument is not of the form host:port, assume that it is a nodeID and 
	// look up the corresponding host/port pair
	if len(strings.Split(nodeToPing, ":")) != 2 {
		id, err := FromString(nodeToPing)
		if err != nil {
			fmt.Printf("Error: Could not convert from string to nodeID (%e)\n",err)
			return nil
		}
		contact := LookupContact(ThisNode,id)
		if contact == nil {
			fmt.Printf("Error: Could not find node with NodeID %s\n",nodeToPing)
			return nil
		}
		nodeToPing = get_host_port(contact)
		//fmt.Printf("Host/Port: %s\n",nodeToPing)
	}
    client, err := rpc.DialHTTP("tcp", nodeToPing)
    if err != nil {
		log.Fatal("DialHTTP: ", err)
    }
    ping := new(Ping)
    ping.MsgID = NewRandomID()
    var pong Pong
    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
	log.Fatal("Call: ", err)
    }
    client.Close()
    fmt.Printf("ping msgID: %s\n", ping.MsgID.AsString())
    fmt.Printf("pong msgID: %s\n", pong.MsgID.AsString())
  sender := new(Contact)
  sender.NodeID = pong.Sender.NodeID
  sender.IPAddr = pong.Sender.IPAddr
  sender.Port = pong.Sender.Port
	return sender
}

func authenticate(hostAndPort string) int {
	client, err := rpc.DialHTTP("tcp", hostAndPort)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
    	}
    
    	req := new(AuthRequest)
    	req.MsgID = ThisNode.ThisContact.NodeID
    	
    	var res AuthResult
        err = client.Call("Kademlia.Authenticate", req, &res)
        client.Close()
    	if err != nil {
    		//Maybe change to not fatal, rather a fail
        	log.Fatal("Call: ", err)
    	}
	return res.isFriend
}

func acquire_lock(hostAndPort string, FileID ID) int {
	client, err := rpc.DialHTTP("tcp", hostAndPort)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
    	}
    
    	req := new(LockRequest)
    	req.MsgID = ThisNode.ThisContact.NodeID
    	req.FileID = FileID
    	
    	var res LockResult
        err = client.Call("Kademlia.Acquire_Lock", req, &res)
        client.Close()
    	if err != nil {
    	//Maybe change to not fatal, rather a fail
        log.Fatal("Call: ", err)
    	}
	return res.is_locked
}

func release_lock(hostAndPort string, FileID ID) int {
	client, err := rpc.DialHTTP("tcp", hostAndPort)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
	    }
	    
	req := new(LockRequest)
	req.MsgID = ThisNode.ThisContact.NodeID
	req.FileID = FileID
	    	
	var res LockResult
	err = client.Call("Kademlia.Acquire_Lock", req, &res)
	client.Close()
	if err != nil {
	    	//Maybe change to not fatal, rather a fail
	        log.Fatal("Call: ", err)
    	}
	return 0
}

func store(hostAndPort string, key ID, data []byte) int {
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
    	req := new(StoreRequest)
    	req.MsgID = NewRandomID()
    	req.Key = key
    	req.Value = data
    	
    	var res StoreResult
      err = client.Call("Kademlia.Store", req, &res)
      client.Close()
    	if err != nil {
        log.Fatal("Call: ", err)
    	}
    /* } else {
    	fmt.Printf("Node already had value.\n")
    } */
	return 1
}
func Find_node(key ID) int {
	bucket,_ := ThisNode.getBucket(key)
	//nodes := bucket.FindNode(key)
	bucket_Find_Node(bucket, key)

	return 0
}

func bucket_Find_Node(b Bucket, key ID){
	for i:=0;i<len(b);i++{
		if key.Equals(b[i].NodeID){
		fmt.Printf("%v %v\n", b[i].IPAddr, b[i].Port)
		return
		}
	}
	fmt.Printf("ERR")
}

func Get_local_value(key ID) int {
    if ThisNode.Data[key] != nil {
		fmt.Printf("OK: %v\n", ThisNode.Data[key])
    } else {
    	fmt.Printf("ERR\n")
    }
    return 0

}
func Whoami() int {
	fmt.Printf("Node ID of this node: %s\n",ThisNode.ThisContact.NodeID.AsString())
	//fmt.Printf("IP/Port: %v %v\n",ThisNode.ThisContact.IPAddr,ThisNode.ThisContact.Port)
	return 0
}

/*
Perform the iterativeStore operation and then print the ID of the node that
received the final STORE operation.
*/
func IterativeStore(key ID, value []byte) int {

	prevDistance := key.Xor(ThisNode.ThisContact.NodeID)
	
	//var closestNode kademlia.FoundNode
	closestNode := ThisNode.ThisContact
	kClosestNodes := make(Bucket,20)
  
	hostPortStr := get_host_port(ThisNode.ThisContact)
	
	//closestnode may want to be its own function that we call from FindNode, or at least
	//that code should be in FindNode, since we need to populate res.Nodes with more than one bucket
	for true {
		//fmt.Printf("%s\n",hostPortStr)
		client, err := rpc.DialHTTP("tcp", hostPortStr)
		if err != nil {
			//fmt.Printf("1\n")
			log.Fatal("DialHTTP: ", err)
		}
		req := new(FindNodeRequest)
		req.MsgID = NewRandomID()
		req.NodeID = key
	
		var res FindNodeResult
		//if FindNode works, all of the closest nodes should be in res.
		err = client.Call("Kademlia.FindNode", req, &res)
    		if err != nil {
			log.Fatal("Call: ", err)
    		}
    client.Close()
    		// obviously we need to do something with the array here, not just take the first element
    		//fmt.Printf("Node 0: %v\n",res.Nodes[0])
    		nextClosestNode, dist := res.Nodes[0], key.Xor(res.Nodes[0].NodeID)
    		for i:= 0; i < len(res.Nodes); i++ {
    			if res.Nodes[i].Port != 0 {
            if res.Nodes[i].NodeID.Xor(key).Less(dist) {
              dist = res.Nodes[i].NodeID.Xor(key)
              nextClosestNode = res.Nodes[i]
            }
          // update kClosestNodes
            replace_idx := -1
            for j := 0; j < len(kClosestNodes); j++ {
              if kClosestNodes[j] == nil {
                kClosestNodes[j] = res.Nodes[i].ToContactPtr()
                replace_idx = -1
                break
              } else if res.Nodes[i].ToContactPtr().NodeID.Xor(key).Less(kClosestNodes[j].NodeID.Xor(key)) {
                if replace_idx != -1 {
                  if kClosestNodes[replace_idx].NodeID.Xor(key).Less(kClosestNodes[j].NodeID.Xor(key)) {
                    replace_idx = j
                  }
                }
              }
            }
            if replace_idx != -1 {
              kClosestNodes[replace_idx] = res.Nodes[i].ToContactPtr()
            }
          }
    		}
    		curDistance := key.Xor(nextClosestNode.NodeID)
    		if !curDistance.Less(prevDistance) {
    			break
    		} else {
    			closestNode = nextClosestNode.ToContactPtr()
          prevDistance = curDistance
    		}
    		hostPortStr = get_host_port(closestNode)
		}
	hostPortStr = get_host_port(closestNode)
	store(hostPortStr, key, value)
  // replicate data across k closest nodes
  for i:=0; i < len(kClosestNodes); i++ {
    if kClosestNodes[i] != nil {
      if !kClosestNodes[i].NodeID.Equals(closestNode.NodeID) {
        hostPortStr = get_host_port(kClosestNodes[i])
        store(hostPortStr, key, value)
      }
    }
  }
	//fmt.Printf("%v\n",closestNode.NodeID)
	return 1
}

/*
Print a list of ≤ k closest nodes and print their IDs. You should collect
the IDs in a slice and print that.
*/
func IterativeFindNode(id ID) Bucket { 
	//Get 20 closest nodes from us.
	req := new(FindNodeRequest)
	req.NodeID = id
	req.MsgID = NewRandomID()
	var k_res FindNodeResult
	//fmt.Printf("In Iterative Find Node, before finding initial nodes closest to NodeID %v\n",id)
	err := ThisNode.FindNode(req,&k_res)
	
	k_closest := foundNodeArr_to_Bucket(k_res.Nodes)
	//fmt.Printf("In Iterative Find Node, after finding initial nodes closest to NodeID\n")
	if err != nil {
		log.Fatal("Call: ", err)
    	}
	//initialize array to hold all 20+20^2 contacts, which we'll sort later
	big_arr := make(Bucket, 420)
  for j := 0; j<len(k_closest); j++{
    big_arr[400+j] = k_closest[j]
  }
	for i :=0;i<len(k_closest);i++{
		//find 20 closest for each node.
		if k_closest[i] == nil {
			// do nothing
		} else if k_closest[i].Port != 0 {
			hostPortStr := get_host_port(k_closest[i])
			//fmt.Printf("Host/Port: %s\n",hostPortStr)
			client, err := rpc.DialHTTP("tcp", hostPortStr)
			if err != nil {
				log.Fatal("DialHTTP: ", err)
			}
			req := new(FindNodeRequest)
			req.MsgID = NewRandomID()
			req.NodeID = id
			
			var res FindNodeResult
			//if FindNode works, all of the closest nodes should be in res.
			err = client.Call("Kademlia.FindNode", req, &res)
			//fmt.Printf("NODES: %v\n",res.Nodes)
			if err != nil {
				log.Fatal("Call: ", err)
    			}
      client.Close()
    			offset:= 20 * i
    			// convert to Bucket type so we can call funcs on it
    			resBucket := foundNodeArr_to_Bucket(res.Nodes)
			for j := 0; j<len(resBucket); j++{
				big_arr[j+offset] = resBucket[j]
			}
		}
		
	}
	// print slice of <= k closest NodeIDs
	//fmt.Printf("%v\n",bucket_to_FoundNodeArr(sort_contacts(big_arr)[:20]))
	return (sort_contacts(big_arr)[:20])
}

/*
printf("%v %v\n", ID, value), where ID refers to the node that finally
returned the value. If you do not find a value, print "ERR".
*/
func IterativeFindValue(key ID) ([]byte, Contact) {
    // dummy to return on error
    var dummy Contact
	// check if this node has the value
	if ThisNode.Data[key] != nil {
		//fmt.Printf("%v %v\n", ThisNode.ThisContact.NodeID, ThisNode.Data[key])
		return ThisNode.Data[key], *ThisNode.ThisContact
	}
  // if this node doesn't have the value, search among the known nodes
	const alpha = 3
	
	contacted_nodes := make(Bucket,1600)
	shortlist := make(Bucket,20)
	shortlist_size := 0
	// The search begins by selecting alpha contacts from the non-empty k-bucket closest to the 
	// bucket appropriate to the key being searched on.
	_, bucket_num := ThisNode.getBucket(key.Xor(ThisNode.ThisContact.NodeID))
	for i := 0; i < 20 && shortlist_size < alpha; i++ {
		if ThisNode.K_Buckets[bucket_num][i] != nil {
			shortlist[shortlist_size] = ThisNode.K_Buckets[bucket_num][i]
			shortlist_size++
		}
	}
	// If there are fewer than alpha contacts in that bucket, contacts are selected from other buckets.
	for b_idx := 0; b_idx < NumBuckets && shortlist_size < alpha; b_idx++ {
		if b_idx != bucket_num {
			for i := 0; i < 20 && shortlist_size < alpha; i++ {
				if ThisNode.K_Buckets[bucket_num][i] != nil {
					shortlist[shortlist_size] = ThisNode.K_Buckets[bucket_num][i]
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
	
		new_shortlist := make(Bucket,400)
		for i := 0; i < len(shortlist); i++ {
			next_open_spot(contacted_nodes)
			contacted_nodes[0] = shortlist[i]
			if shortlist[i] != nil {
				hostPortStr := get_host_port(contacted_nodes[i])
		
				client, err := rpc.DialHTTP("tcp", hostPortStr)
				if err != nil {
                    // if we can't contact the node, remove it from shortlist
                    shortlist[i] = nil
				} else {
					req := new(FindValueRequest)
					req.MsgID = NewRandomID()
					req.Key = key
			
					var res FindValueResult
					err = client.Call("Kademlia.FindValue", req, &res)
					if err != nil {
						log.Fatal("Call: ", err)
    				}
          client.Close()
    				// if res.Err is nil, the node contains the value
    				if res.Err == nil {
              				//fmt.Printf("OK\n")
    					fmt.Printf("%v %v\n", shortlist[i].NodeID, res.Value)
    					return res.Value, *contacted_nodes[i]
    				} else {
    					offset:= 20 * i
    					resBucket := foundNodeArr_to_Bucket(res.Nodes)
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
    			//fmt.Printf("ERR\n")
    			return nil,dummy
    		}
    	}
    }
    return nil,dummy
}

func (k *Kademlia) find_closest(req_id ID, count int) []*Contact{
	//fmt.Printf("Prepare to Xor:\n|%v|\n|%v|\n", req_id, k.ThisContact.NodeID)
	b_num := req_id.Xor(k.ThisContact.NodeID).PrefixLen() //get bucket number
	if b_num == 160{ // if req_id == k.NodeID, b_num will be 160. In this case use b_num = 159
    b_num--
	}
	//fmt.Printf("tried to access bucket %d\n",b_num)
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
	nodes = sort_contacts(nodes)
	return nodes
}

//enter in shell as "run host:port host:port"
func Run(listenStr string, firstPeerStr string) int {
  
  fmt.Printf("Kademlia starting up!\n")
  ThisNode = NewKademlia()
  ThisNode.ThisContact.IPAddr = strings.Split(listenStr, ":")[0]
  port,_ := strconv.Atoi(strings.Split(listenStr, ":")[1])
  ThisNode.ThisContact.Port = uint16(port)
  //Register on server
  rpc.Register(ThisNode)
  rpc.HandleHTTP()
  l, err := net.Listen("tcp", listenStr)
  if err != nil {
		log.Fatal("Listen: ", err)
  }
  
  // Serve forever.
  go http.Serve(l, nil)
	
	
	/*
   Add the first contact. For now, just create a new contact with host, port and a random nodeID
   */
	// ping the first peer
	firstPeerContact := Ping2(firstPeerStr)
	ThisNode.addContactToBuckets(firstPeerContact)
	//fmt.Printf("Made it to before iterative\n")
	// find and add the closest contacts to this node
	closestContacts := IterativeFindNode(ThisNode.ThisContact.NodeID)
	//fmt.Printf("Made it through iterativeFindNode\n")
	for i := 0; i < len(closestContacts); i++ {
		if closestContacts[i] != nil {
			ThisNode.addContactToBuckets(closestContacts[i])
		}
	}
	// initialize the p2p functionality
    InitP2P(ThisNode)
	return 1
  
}