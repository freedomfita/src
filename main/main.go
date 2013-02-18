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
	    } else if arg_s[0] == "t" {
	    	fmt.Printf("Kademlia starting up!\n")
	    	thisNode = kademlia.NewKademlia()
	    	    
	    	//Register on server
	    	rpc.Register(thisNode)
	    	rpc.HandleHTTP()
	    	// why is localhost:7890 hardcoded here??
	    	l, err := net.Listen("tcp", "localhost:7890")
	    	if err != nil {
	    		log.Fatal("Listen: ", err)
	    	}
	    	
	    	// Serve forever.
	    	go http.Serve(l, nil)
		running = true
	    	thisNode.Main_Testing()
	    } else if arg_s[0] == "run" && is_cmd_valid(arg_s,2,false) {
	    	run(arg_s[1],arg_s[2])
	    } else if arg_s[0] == "ping" && is_cmd_valid(arg_s,1,true) {
	    	ping(arg_s[1])
	    } else if arg_s[0] == "store" && is_cmd_valid(arg_s,2,true) {
	    	// k and b are just placeholders for now
	    	k,_ := kademlia.FromString(arg_s[1])
	    	b := []byte(arg_s[2])
	    	iterativeStore(k,b)
	    } else if arg_s[0] == "find_node" && is_cmd_valid(arg_s,1,true) {
	    	    	id, err := kademlia.FromString(arg_s[1])
		    	if err != nil {
		    		log.Fatal("Find Node: ",err)
		    	}
		    	find_node(id)

	    } else if arg_s[0] == "find_value" && is_cmd_valid(arg_s,1,true) {
	    	k,_ := kademlia.FromString(arg_s[1])
	    	iterativeFindValue(k)
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
    thisNode.ThisContact.IPAddr = strings.Split(listenStr, ":")[0]
    port,_ := strconv.Atoi(strings.Split(listenStr, ":")[1])
    thisNode.ThisContact.Port = uint16(port)
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
	thisNode.AddContactToBuckets(firstPeerContact)
	fmt.Printf("Made it to before iterative\n")
	// find and add the closest contacts to this node
	closestContacts := iterativeFindNode(thisNode.ThisContact.NodeID)
	fmt.Printf("Made it through iterativeFindNode\n")
	for i := 0; i < len(closestContacts); i++ {
		if closestContacts[i] != nil {
			fmt.Printf("contact: %v\n",closestContacts[i])
			thisNode.AddContactToBuckets(closestContacts[i])
		}
	}
	//id_list := thisNode.Local_Random_Nodes()
	//fmt.Printf("Made it through, have %d random nodes now in our buckets\n", len(id_list))
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

