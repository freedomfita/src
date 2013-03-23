package kademlia

//only populates the kbuckets randomly, not by distance, more just to check things
import (
	"fmt"
	"log"
)

func (k *Kademlia) Random_Nodes() []ID{
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
	return id_list

}

//takes in id_list from Random_Nodes and runs test
func TestDHTFunctionality(k *Kademlia){
	k.ThisContact.IPAddr= "localhost"
	k.ThisContact.Port = 7890
	id_list := k.Local_Random_Nodes()
	fmt.Printf("*****************\n*****************\n*****************\n")
	Test_Iterative_Find_Node(id_list)
	fmt.Printf("*****************\n*****************\n*****************\n")
	k.Test_Find_Nodes(id_list)
	fmt.Printf("*****************\n*****************\n*****************\n")
	//k.Print_KBuckets()
	fmt.Printf("*****************\n*****************\n*****************\n")
	id_key_list := Test_Iterative_Store()
	fmt.Printf("*****************\n*****************\n*****************\n")
	Test_Iterative_Find_Value(id_key_list)
	fmt.Printf("*****************\n*****************\n*****************\n")
	k.Print_KBuckets_bare()
  fmt.Printf("*****************\n*****************\n*****************\n")
}

//Tests find Node
func (k *Kademlia) Test_Find_Nodes(id_list []ID){

	for i:=0;i<len(id_list);i++{
		req := new(FindNodeRequest)
		req.NodeID = id_list[i]
		req.MsgID = NewRandomID()
		var k_res FindNodeResult
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

func Test_Iterative_Find_Node(id_list []ID){
	num_to_test := len(id_list)
	for i:=0;i<num_to_test;i++{
		b := IterativeFindNode(id_list[i])
		fmt.Printf("Bucket#%v :\n%v\n",i,b)
	}
}

func Test_Iterative_Store() []ID{
	//pass in random value
	val := make([]byte,20)
	val[0] = 1
	id_list := make([]ID,num_test_nodes)
	for i:=0;i<num_test_nodes;i++{
		id_list[i] = NewRandomID()
		fmt.Printf("Created New ID : %v\n",id_list[i])
		err := IterativeStore(id_list[i], val)
		if err != 1{
			fmt.Printf("ERROR at IterativeStore\n")
		}
	}
	return id_list
	
}

func Test_Iterative_Find_Value(id_list []ID) {
	for i:=0;i<len(id_list);i++{
		val,_ := IterativeFindValue(id_list[i])
		if val == nil {
			fmt.Printf("ERROR at IterativeFindValue for %v\n",id_list[i])
		}
	}
}

func (k *Kademlia) Print_KBuckets(){
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

func (k *Kademlia) Print_KBuckets_bare(){
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

