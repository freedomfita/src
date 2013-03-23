package kademlia

import (
	"fmt"
)

const num_test_nodes = 150

//we create lots of nodes, then we can go through and 
//try downloading files from one node, then multiple.
//use localhost for these nodes
func (k *Kademlia) Local_Random_Nodes() []ID{
	id_list := make([]ID, num_test_nodes)
	for i:=0;i<num_test_nodes;i++{
		c:= new(Contact)
		c.NodeID = NewRandomID()
		id_list[i] = c.NodeID
		c.IPAddr = "localhost"
		c.Port = 7890
		b_num := c.NodeID.Xor(k.ThisContact.NodeID).PrefixLen()
		k.next_open_spot(b_num)
		k.K_Buckets[b_num][0] = c
	}
	k.F_Buckets = k.K_Buckets
	return id_list
}



func (k *Kademlia) Download_File_Testing() int {
	//id_list:= k.Local_Random_Nodes()
	//now we have a list of nodes
	for i:= 0; i< 1; i++ {
		fmt.Printf("testing download file\n")
		DownloadFile("Homework1.pdf", "sharing")
		fmt.Printf("made it through test download: %d\n",i)
	}
	return 1
}

func (k *Kademlia) Test_Locking() int {
	fmt.Printf("starting lock testing\n")
	file_info := GetFileInfo("Homework1.pdf")
	fmt.Printf("File ID: %s\n",file_info.FileID)
	k.Request_Lock(file_info.FileID)
	fmt.Printf("made it through lock\n")
	k.Notify_Release_Lock(file_info.FileID)
	fmt.Printf("made it through unlock\n")
	return 1
}


