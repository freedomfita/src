package userio

import (
    "flag"
    "fmt"
    "log"
    "math/rand"
    "time"
    "os"
    "bufio"
    "strings"
    "kademlia"
)

//** Globals **************
var gLoggedIn bool = false
var gPromptStr string = "> "
var gMaxFailedLogins int = 0
// maps command names to corresponding command structs
var gCmdMap map[string]Command = make(map[string]Command)
//*************************

type Command struct {
    name string
    numArgs int
    cmdFunc func([]string)
}

func MainLoop() {
    
    // fill command map with the appropriate structs
    gCmdMap["ping"] = Command{"ping",1,do_ping }
    gCmdMap["q"] = Command{ "q",0,do_quit }
    gCmdMap["get_contact"] = Command{ "get_contact",1,do_get_contact }
    gCmdMap["iterativeStore"] = Command{ "iterativeStore",2,do_iterative_store }
    gCmdMap["iterativeFindNode"] = Command{ "iterativeFindNode",1,do_iterative_findnode }
    gCmdMap["iterativeFindValue"] = Command{ "iterativeFindValue",1,do_iterative_findvalue }
    gCmdMap["local_find_value"] = Command{ "local_find_value",1,do_local_findvalue }
    gCmdMap["whoami"] = Command{ "whoami",0,do_whoami }
    
    var listenStr,firstPeerStr string
    // By default, Go seeds its RNG with 1. This would cause every program to
    // generate the same sequence of IDs.
    rand.Seed(time.Now().UnixNano())
    
    // Get the bind and connect connection strings from command-line arguments.
    flag.Parse()
    args := flag.Args()
    
    // if not initialized with two command line args, prompt for them
    if len(args) != 2 {
        log.Println("Invoke Kademlia with 'run <host:port> <host:port>, or type 'q' to exit.")
        for true {
            args = get_input(gPromptStr)
            if len(args) == 1 { 
                if args[0] == "q" || args[0] == "Q" { return } 
            } else if len(args) == 3 {
                if args[0] == "run" {
                    listenStr = args[1]
                    firstPeerStr = args[2]
                    break
                }
            }
        } 
    } else {
        listenStr = args[0]
        firstPeerStr = args[1]
    }
    kademlia.Run(listenStr,firstPeerStr)
    
    // loop until the user exits
    for true {
        args = get_input(gPromptStr)
	    /*log.Printf("there were %d\n", length)
         for i := 0; i < length; i++ {
         log.Printf("number %d: %s", i, arg_s[i])
         }*/
        if !gLoggedIn {
            // loop until the user provides valid credential or fails too often
            numTries := 0
            for !process_login() {
                if numTries >= gMaxFailedLogins {
                    fmt.Printf("Too many failed login attempts.\n")
                    os.Exit(1)
                }
                numTries++
            }
        }
        interpret_cmdline(args)
    }
}


// print a prompt, get user input and return array of args (split input at whitespace)
func get_input(prompt string) (ret []string) {
    fmt.Printf(prompt)
    reader := bufio.NewReader(os.Stdin)
    input,_:= reader.ReadString('\n')
    //input includes both a carriage return and newline, trim whitespace
    input = strings.TrimSpace(input)
    ret = strings.Split(input, " ")
    return
}

// check if the number of parameters is correct for a command
func is_cmd_valid(cmd []string, argc int) bool {
    if argc != len(cmd)-1 {
        fmt.Printf("Error: Command '%s' must be invoked with %d arguments!\n",cmd[0],argc)
        return false
    }
    return true
}

// parse the command line, given an array of arguments, and execute the appropriate command or print an error msg
func interpret_cmdline(argv []string) {
    if len(argv) == 0 { //do nothing if command line is empty
        return
    } else if cmd,ok := gCmdMap[argv[0]]; ok {
        if (is_cmd_valid(argv,cmd.numArgs)) {
            cmd.cmdFunc(argv)
        }
        return
    }
    fmt.Println("Command/s unknown.")
    return
}

func process_login() bool {
    return true
}