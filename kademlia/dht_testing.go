package kademlia
//only populates the kbuckets randomly, not by distance, more just to check things
import (
	"fmt"
	"log"
)
import (
	"main"
)
const num_test_nodes = 50

//use MoritzIP for these nodes
func (k *kademlia.Kademlia) Random_Nodes() []kademlia.ID{
	id_list := make([]kademlia.ID, num_test_nodes)
	for i:=0;i<num_test_nodes;i++{
		c:= new(kademlia.Contact)
		c.NodeID = kademlia.NewRandomID()
		id_list[i] = c.NodeID
		c.IPAddr = "192.168.0.123"
		c.Port = 7890
		b_num := c.NodeID.Xor(k.ThisContact.NodeID).PrefixLen()
		k.Next_Open_Spot(b_num)
		k.K_Buckets[b_num][0] = c
	}
	return id_list

}
//use localhost for these nodes
func (k *kademlia.Kademlia) Local_Random_Nodes() []kademlia.ID{
	id_list := make([]kademlia.ID, num_test_nodes)
	for i:=0;i<num_test_nodes;i++{
		c:= new(kademlia.Contact)
		c.NodeID = kademlia.NewRandomID()
		id_list[i] = c.NodeID
		c.IPAddr = "localhost"
		c.Port = 7890
		b_num := c.NodeID.Xor(k.ThisContact.NodeID).PrefixLen()
		k.Next_Open_Spot(b_num)
		k.K_Buckets[b_num][0] = c
	}
	return id_list

}

//takes in id_list from Random_Nodes and runs test
func (k *kademlia.Kademlia) Main_Testing(){

	//id_list := k.Local_Random_Nodes()
	fmt.Printf("*****************\n*****************\n*****************\n")
	//Test_Iterative_Find_Node(id_list)
	fmt.Printf("*****************\n*****************\n*****************\n")
	//k.Test_Find_Nodes(id_list)
	fmt.Printf("*****************\n*****************\n*****************\n")
	//k.Print_KBuckets()
	fmt.Printf("*****************\n*****************\n*****************\n")
	//k.Print_KBuckets_bare()
}

//Tests find Node
func (k *kademlia.Kademlia) Test_Find_Nodes(id_list []kademlia.ID){

	for i:=0;i<len(id_list);i++{
		req := new(kademlia.FindNodeRequest)
		req.NodeID = id_list[i]
		req.MsgID = kademlia.NewRandomID()
		var k_res kademlia.FindNodeResult
		err := k.FindNode(req,&k_res)
		if err != nil {
			log.Fatal("Call: ", err)
		}
		b := k_res.Nodes
		//b := main.iterativeFindNode(id_list[i])
		fmt.Printf("results for b%v found, begin printing\n",i)
		for j:=0;j<len(b);j++{
			fmt.Printf("#%v: %v\n",j,b[j])
			if b[j].Port == 0{
				fmt.Printf("B%v has %v elements\n", i,j)
				break
			}
		}
		fmt.Printf("finished printing b results\n")

	}

}
//Test Iterative find node

func Test_Iterative_Find_Node(id_list []kademlia.ID){
	num_to_test := 2
	for i:=0;i<num_to_test;i++{
		b := iterativeFindNode(id_list[i])
		fmt.Printf("Bucket#%v :\n%v\n",i,b)
	}
}

func (k *kademlia.Kademlia) Print_KBuckets(){
	for i:=0;i<15;i++{//len(k.K_Buckets);i++{
		fmt.Printf("Printing Bucket #%v\n",i)
		kb := k.K_Buckets[i]
		for j:=0;j<len(kb);j++{
			/*if kb[j] == nil{
				fmt.Printf("Bucket #%v printed with %v elements\n",i,j)
				break
			} else {*/
			fmt.Printf("B%v E%v :%v\n", i,j,kb[j])
			//}
			
		}
	}
}

func (k *kademlia.Kademlia) Print_KBuckets_bare(){
	for i:=0;i<160;i++{
		count:= -1
		kb:= k.K_Buckets[i]
		for j:=0;j<len(kb);j++{
			if kb[j] == nil{
			count = j
				//fmt.Printf("Bucket %v has %v elements\n",i,j)
				break
			}
		}
		if count==-1{ //bucket is full
			count = len(kb)
		}
		fmt.Printf("Bucket %v has %v elements\n",i,count)
	
		
	}
}






