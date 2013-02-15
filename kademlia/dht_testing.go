package kademlia
//***NOT NEARLY FUNCTIONAL
func (k *Kademlia) Add_Random_Nodes(){
	for i:=0;i<50;i++{
			c := new(Contact)
			c.NodeID = NewRandomID()
			c.IPAddr = "192.168.0.123"
			c.Port = 7890
			//b_num := c.NodeID.Xor(k.ThisContact.NodeID).PrefixLen()
			k.K_Buckets[i][0] = c
	}
}