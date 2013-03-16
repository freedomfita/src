package authentication

import (
    "kademlia"
)

type User struct {
    username string
    userID kademlia.ID
}

type UserGroup struct {
    name string
    access map[User]int
}


func (k *Kademlia) Authenticate(req AuthRequest, *res AuthResult) error {
	res.MsgID = CopyID(req.MsgID)
	//Now run function to test whether or not we have node in F_Buckets
	res.isFriend = k.find_friend(req.MsgID)
}



func (k *Kademlia) find_friend(req_id ID) int{
	//fmt.Printf("Prepare to Xor:\n|%v|\n|%v|\n", req_id, k.ThisContact.NodeID)
	b_num := req_id.Xor(k.ThisContact.NodeID).PrefixLen() //get bucket number
	
	// if req_id == k.NodeID, b_num will be 160. In this case use b_num = 159
	if b_num == 160{ 
    		b_num--
	}
	//MODIFIED * changing the function just to find the one node in corresponding bucket
	//fmt.Printf("tried to access bucket %d\n",b_num)
	b := k.F_Buckets[b_num] //get corresponding bucket
	for i:=0;i<len(b);i++{ //we copy all contacts from closest bucket
		if b[i] == nil{
			continue
		}
		else if b[i].NodeID.Equals(req_id){
			return 1
		}
	}
	return 0
}