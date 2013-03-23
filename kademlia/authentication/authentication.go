package kademlia

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



