package main

import (
    "flag"
    //"fmt"
    "log"
    "math/rand"
    //"net"
    //"net/http"
    //"net/rpc"
    "time"
    "os"
    "bufio"
    "strings"
)

import (
    "kademlia"
)

//myIP 192.168.0.145:7890
//MoritzIP 192.168.0.123:7890
func main() {

    // By default, Go seeds its RNG with 1. This would cause every program to
    // generate the same sequence of IDs.
    rand.Seed(time.Now().UnixNano())
	    // Get the bind and connect connection strings from command-line arguments.
    
    flag.Parse()
    args := flag.Args()
    if len(args) != 2 {
        log.Fatal("Must be invoked with exactly two arguments!\n")
    }
    listenStr := args[0]
    firstPeerStr := args[1]
    kademlia.Run(listenStr,firstPeerStr)
    /*
    log.Printf("\nType q to quit\nFormat of Run\nrun str str")
    */
    
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
	    /*log.Printf("there were %d\n", length)
	    for i := 0; i < length; i++ {
	    	log.Printf("number %d: %s", i, arg_s[i])
	    }*/
	    
	    if len(arg_s) == 0 { //do nothing if command line is empty
	    	continue
	    } else if arg_s[0] == "q" {
	    	return
	    } else if arg_s[0] == "ping" && is_cmd_valid(arg_s,1) {
	    	kademlia.Ping2(arg_s[1])
	    } else if arg_s[0] == "get_contact" && is_cmd_valid(arg_s,1){
	    		id, err := kademlia.FromString(arg_s[1])
			if err != nil {
				log.Fatal("Find Node: ",err)
			}
		    	kademlia.Find_node(id)
	    } else if arg_s[0] == "iterativeStore" && is_cmd_valid(arg_s,2) {
	    	k,_ := kademlia.FromString(arg_s[1])
	    	b := []byte(arg_s[2])
	    	kademlia.IterativeStore(k,b)
	    } else if arg_s[0] == "iterativeFindNode" && is_cmd_valid(arg_s,1) {
	    	    	id, err := kademlia.FromString(arg_s[1])
		    	if err != nil {
		    		log.Fatal("Find Node: ",err)
		    	}
		    	kademlia.IterativeFindNode(id)

	    } else if arg_s[0] == "iterativeFindValue" && is_cmd_valid(arg_s,1) {
	    	k,_ := kademlia.FromString(arg_s[1])
	    	kademlia.IterativeFindValue(k)
	    } else if arg_s[0] == "local_find_value" && is_cmd_valid(arg_s,1) {
	    	    	id, err := kademlia.FromString(arg_s[1])
		    	if err != nil {
		    		log.Fatal("Get Local Value: ",err)
		    	}
		    	kademlia.Get_local_value(id)

	    } else if arg_s[0] == "whoami" && is_cmd_valid(arg_s,0) {
	    	kademlia.Whoami()
	    } else {
	    	log.Printf("Command/s unknown.")
	    }
    }
    
}

// check if the number of parameters is correct and global var running is equal to "status"
func is_cmd_valid(cmd []string, argc int) bool {
  if argc != len(cmd)-1 {
		log.Printf("Error: Command '%s' must be invoked with %d arguments!\n",cmd[0],argc)
		return false
	}
	return true
}

