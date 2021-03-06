package userio

import (
    "os"
    "log"
    "kademlia"
)
func do_main_testing(argv []string){
	kademlia.Main_Testing()
}

func notify_release_lock(argv []string){
	id, err := kademlia.FromString(argv[1])
	if err != nil {
		log.Fatal("Find Node: ",err)
	}
	kademlia.ThisNode.Notify_Release_Lock(id)
}
func request_lock(argv []string){
	id, err := kademlia.FromString(argv[1])
	if err != nil {
		log.Fatal("Find Node: ",err)
	}
	kademlia.ThisNode.Request_Lock(id)
}
func do_authenticate(argv []string){
	kademlia.Authenticate(argv[1])
}

func do_quit(argv []string) {
    os.Exit(0)
}

func do_ping(argv []string) {
    kademlia.Ping2(argv[1])
}

func do_get_contact(argv []string) {
    id, err := kademlia.FromString(argv[1])
    if err != nil {
        log.Fatal("Find Node: ",err)
    }
    kademlia.Find_node(id)
}

func do_iterative_store(argv []string) {
    k,err := kademlia.FromString(argv[1])
    if err != nil {
        log.Fatal("Iterative Store: ",err)
    }
    b := []byte(argv[2])
    kademlia.IterativeStore(k,b)
}

func do_iterative_findnode(argv []string) {
    id, err := kademlia.FromString(argv[1])
    if err != nil {
        log.Fatal("Find Node: ",err)
    }
    kademlia.IterativeFindNode(id)
}

func do_iterative_findvalue(argv []string) {
    k,_ := kademlia.FromString(argv[1])
    kademlia.IterativeFindValue(k)
}

func do_local_findvalue(argv []string) {
    id, err := kademlia.FromString(argv[1])
    if err != nil {
        log.Fatal("Get Local Value: ",err)
    }
    kademlia.Get_local_value(id)
}

func do_whoami(argv []string) {
    kademlia.Whoami()
}

func do_download(argv []string) {
    kademlia.DownloadFile(argv[1],argv[2],false)
}

func do_download_dir(argv []string) {
    kademlia.DownloadDirectory(argv[1],argv[2],argv[3],false)
}