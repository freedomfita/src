package kademlia
//only populates the kbuckets randomly, not by distance, more just to check things
func (k *Kademlia) Add_Random_Nodes(){
	L := len(k.K_Buckets)
	for j:=0;j<L;j++{
		for i:=0;i<20;i++{
			c := new(Contact)
			c.NodeID = NewRandomID()
			c.IPAddr = "192.168.0.123"
			c.Port = 7890
			//b_num := c.NodeID.Xor(k.ThisContact.NodeID).PrefixLen()
			k.K_Buckets[j][i] = c
		}
	}
}
//use MoritzIP for these nodes
func (k *Kademlia) Random_Nodes() []ID{
	id_list := make([]ID, 50)
	for i:=0;i<50;i++{
		c:= new(Contact)
		c.NodeID = NewRandomID()
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
func (k *Kademlia) Local_Random_Nodes() []ID{
	id_list := make([]ID, 50)
	for i:=0;i<50;i++{
		c:= new(Contact)
		c.NodeID = NewRandomID()
		id_list[i] = c.NodeID
		c.IPAddr = "localhost"
		c.Port = 7890
		b_num := c.NodeID.Xor(k.ThisContact.NodeID).PrefixLen()
		k.Next_Open_Spot(b_num)
		k.K_Buckets[b_num][0] = c
	}
	return id_list

}