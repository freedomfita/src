package main

import (
    "flag"
    "fmt"
    "log"
    "math/rand"
    "net"
    "net/http"
    "net/rpc"
    "time"
    "os"
    "bufio"
    "strings"
    "strconv"
)

import (
    "kademlia"
)

var running bool
var thisNode *kademlia.Kademlia
//myIP 192.168.0.145:7890
//MoritzIP 192.168.0.123:7890
func main() {

	running = false
    // By default, Go seeds its RNG with 1. This would cause every program to
    // generate the same sequence of IDs.
    rand.Seed(time.Now().UnixNano())
	    // Get the bind and connect connection strings from command-line arguments.
    flag.Parse()
    //args := flag.Args()
    log.Printf("\nType q to quit\nFormat of Run\nrun str str")
    
    
    for true {
	    reader := bufio.NewReader(os.Stdin)
	    input,_:= reader.ReadString('\n')
	    //input includes both a carriage return and newline, trim whitespace
	    input = strings.TrimSpace(input)
	    /*l := len(input)
	    for i:=0; i<l; i++ {
	    	log.Printf("%d", input[i])}
	    	*/
	    
	    arg_s := strings.Split(input, " ")
	    length := len(arg_s)
	    log.Printf("there were %d\n", length)
	    for i := 0; i < length; i++ {
	    	log.Printf("number %d: %s", i, arg_s[i])
	    }
	    
	    if len(arg_s) == 0 { //do nothing if command line is empty
	    	continue
	    } else if arg_s[0] == "q" {
	    	return
	    } else if arg_s[0] == "run" && is_cmd_valid(arg_s,2,false) {
	    	run(arg_s[1],arg_s[2])
	    } else if arg_s[0] == "ping" && is_cmd_valid(arg_s,1,true) {
	    	ping(arg_s[1])
	    } else if arg_s[0] == "store" && is_cmd_valid(arg_s,3,true) {
	    	// k and b are just placeholders for now
	    	var k kademlia.ID
	    	var b []byte
	    	iterativeStore(k,b)
	    } else if arg_s[0] == "find_node" && is_cmd_valid(arg_s,1,true) {
	    	    	id, err := kademlia.FromString(arg_s[1])
		    	if err != nil {
		    		log.Fatal("Find Node: ",err)
		    	}
		    	find_node(id)

	    } else if arg_s[0] == "find_value" && is_cmd_valid(arg_s,1,true) {
	    	find_value(arg_s[1])
	    } else if arg_s[0] == "get_local_value" && is_cmd_valid(arg_s,1,true) {
	    	    	id, err := kademlia.FromString(arg_s[1])
		    	if err != nil {
		    		log.Fatal("Get Local Value: ",err)
		    	}
		    	get_local_value(id)

	    } else if arg_s[0] == "get_node_id" && is_cmd_valid(arg_s,0,true) {
	    	get_node_id()
	    } else {
	    	log.Printf("Command/s unknown.")
	    }
    }
    
}
//enter in shell as "run host:port host:port"
func run(listenStr string, firstPeerStr string) int {

    fmt.Printf("Kademlia starting up!\n")
    thisNode = kademlia.NewKademlia()
    
    //Register on server
    rpc.Register(thisNode)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", listenStr)
    if err != nil {
		log.Fatal("Listen: ", err)
    }

    // Serve forever.
    go http.Serve(l, nil)
	running = true
	
	
	/*
	Add the first contact. For now, just create a new contact with host, port and a random nodeID
	*/
	// ping the first peer
	firstPeerContact := ping(firstPeerStr)
	fmt.Printf("Made it to before iterative\n")
	// find and add the closest contacts to this node
	closestContacts := iterativeFindNode(firstPeerContact.NodeID)
	fmt.Printf("Made it through iterativeFindNode\n")
	for i := 0; i < len(closestContacts); i++ {
		fmt.Printf("finding closest contact %d\n",i)
		thisNode.AddContactToBuckets(closestContacts[i])
	}
	id_list := thisNode.Local_Random_Nodes()
	fmt.Printf("Made it through, have %d random nodes now in our buckets\n", len(id_list))
	return 1

}

// check if the number of parameters is correct and global var running is equal to "status"
func is_cmd_valid(cmd []string, argc int, status bool) bool {
	if status != running {
		if running {
			log.Printf("Error: Kademlia is already running.")
		} else {
			log.Printf("Error: Kademlia is not yet running.")
		}
		return false
	} else if argc != len(cmd)-1 {
		log.Printf("Error: Command '%s' must be invoked with %d arguments!\n",cmd[0],argc)
		return false
	}
	return true
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
		hostPort := make([]string, 2)
		hostPort[0] = contact.IPAddr
		hostPort[1] = strconv.FormatUint(uint64(contact.Port),10)
		nodeToPing = strings.Join(hostPort, ":")
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
	return &pong.Sender
}

func store(hostAndPort string, key kademlia.ID, data []byte) int {
	client, err := rpc.DialHTTP("tcp", hostAndPort)
	if err != nil {
		log.Fatal("DialHTTP: ", err)
    }
    
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
	return 1
}
func find_node(key kademlia.ID) int {
	bucket,_ := thisNode.GetBucket(key)
	nodes := bucket.FindNode(key)
	fmt.Println(nodes)

	return 0
}
func find_value(key string) int {
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
	return 0
}

/*
Perform the iterativeStore operation and then print the ID of the node that
received the final STORE operation.
*/
func iterativeStore(key kademlia.ID, value []byte) int {

	prevDistance := key.Xor(thisNode.ThisContact.NodeID)
	
	//var closestNode kademlia.FoundNode
	var closestNode *kademlia.Contact
	
	hostPort := make([]string, 2)
	hostPort[0] = thisNode.ThisContact.IPAddr
	hostPort[1] = strconv.FormatUint(uint64(thisNode.ThisContact.Port),10)
	hostPortStr := strings.Join(hostPort, ":")
	
	//closestnode may want to be its own function that we call from FindNode, or at least
	//that code should be in FindNode, since we need to populate res.Nodes with more than one bucket
	for true {
		
		client, err := rpc.DialHTTP("tcp", hostPortStr)
		if err != nil {
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
    		curDistance := key.Xor(res.Nodes[0].NodeID)
    	
    		if !curDistance.Less(prevDistance) {
    			closestNode = res.Nodes[0]
    			break
    		}
    		hostPort[0] = closestNode.IPAddr
			hostPort[1] = strconv.FormatUint(uint64(closestNode.Port),10)
			hostPortStr = strings.Join(hostPort, ":")
		}
	
	hostPort[0] = closestNode.IPAddr
	hostPort[1] = strconv.FormatUint(uint64(closestNode.Port),10)
	hostPortStr = strings.Join(hostPort, ":")
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
	k_closest := k_res.Nodes
	fmt.Printf("In Iterative Find Node, after finding initial nodes closest to NodeID\n")
	if err != nil {
		log.Fatal("Call: ", err)
    	}
	//initialize array to hold all 20^2 contacts, which we'll sort later
	big_arr := make(kademlia.Bucket, 400)
	for i :=0;i<len(k_closest);i++{
		//find 20 closest for each node.
		hostPort := make([]string, 2)
		hostPort[0] = k_closest[i].IPAddr
		hostPort[1] = strconv.FormatUint(uint64(k_closest[i].Port),10)
		hostPortStr := strings.Join(hostPort, ":")
		
		client, err := rpc.DialHTTP("tcp", hostPortStr)
		if err != nil {
			log.Fatal("DialHTTP: ", err)
		}
		req := new(kademlia.FindNodeRequest)
		req.MsgID = kademlia.NewRandomID()
		req.NodeID = id
			
		var res kademlia.FindNodeResult
		//if FindNode works, all of the closest nodes should be in res.
		err = client.Call("Kademlia.FindNode", req, &res)
		if err != nil {
			log.Fatal("Call: ", err)
    		}
    		offset:= 20 * i
		for j := 0; j<len(res.Nodes); j++{
			big_arr[j+offset] = res.Nodes[j]
		}
		
	}
	fmt.Printf("Finished IterativeFindNode and returning array of contacts\n")
	return (kademlia.Sort_Contacts(big_arr))[:20]
}

/*
printf("%v %v\n", ID, value), where ID refers to the node that finally
returned the value. If you do not find a value, print "ERR".
*/
func iterativeFindValue(key /*kademlia.ID*/ string) *kademlia.Contact {
	/*const alpha := 3
	shortlist := make(kademlia.Bucket,3)
	shortlist_size := 0
	// The search begins by selecting alpha contacts from the non-empty k-bucket closest to the 
	// bucket appropriate to the key being searched on.
	_, bucket_num := thisNode.GetBucket(key.Xor(thisNode.ThisContact.NodeID))
	for i := 0; i < 20; i++ {
		if thisNode.K_Buckets[i] != nil {
			shortlist[shortlist_size] = thisNode.K_Buckets[i]
		}
		if shortlist_size >= 
	}
	// If there are fewer than alpha contacts in that bucket, contacts are selected from other buckets.
	for b_idx := 0; b < kademlia.NumBuckets; b++ {
		if b_idx != bucket_num {
			
		}
	}
	// The contact closest to the target key, closestNode, is noted.
	closestNode := kademlia.Sort_Contacts(shortlist)[0]
// The node then sends parallel, asynchronous FIND_* RPCs to the alpha contacts in the 
// shortlist. Each contact, if it is live, should normally return k triples. If any of the 
// alpha contacts fails to reply, it is removed from the shortlist, at least temporarily.

// The node then fills the shortlist with contacts from the replies received. These are those closest to the target. From the shortlist it selects another alpha contacts. The only condition for this selection is that they have not already been contacted. Once again a FIND_* RPC is sent to each in parallel.

// Each such parallel search updates closestNode, the closest node seen so far.

// The sequence of parallel searches is continued until either no node in the sets returned is closer than the closest node already seen or the initiating node has accumulated k probed and known to be active contacts.

// If a cycle doesn't find a closer node, if closestNode is unchanged, then the 
// initiating node sends a FIND_* RPC to each of the k closest nodes that it has not already queried.

// At the end of this process, the node will have accumulated a set of k active contacts 
// or (if the RPC was FIND_VALUE) may have found a data value. Either a set of triples or 
// the value is returned to the caller.

Kademlia uses a value of 3 for alpha, the degree of parallelism used. It appears that (see stutz06) this value is optimal.

There are at least three approaches to managing parallelism. The first is to launch alpha probes and wait until all have succeeded or timed out before iterating. This is termed strict parallelism. The second is to limit the number of probes in flight to alpha; whenever a probe returns a new one is launched. We might call this bounded parallelism. A third is to iterate after what seems to be a reasonable delay (duration unspecified), so that the number of probes in flight is some low multiple of alpha. This is loose parallelism and the approach used by Kademlia.
*/
/*
	req := new(kademlia.FindValueRequest)
	req.MsgID = kademlia.NewRandomID()
	req.Key = key
	var res kademlia.FindValueResult
	
	for true {
	
	}
	*/
	return nil
}