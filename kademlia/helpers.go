package kademlia

import (
  "strings"
  "strconv"
  "sort"
)

// sort function for buckets
func sort_contacts(arr Bucket) Bucket {
	sort.Sort(BucketSort_ByNodeID{arr})
	return arr//sorted_arr
}

func bucket_to_FoundNodeArr(bucket Bucket) []FoundNode {
	b := make([]FoundNode,len(bucket))
	j := 0
	for i := 0; i < len(bucket); i++ {
		if bucket[i] != nil {
			b[j].IPAddr = bucket[i].IPAddr
			b[j].NodeID = bucket[i].NodeID
			b[j].Port = bucket[i].Port
			j++
		}
	}
	return b
}

func foundNodeArr_to_Bucket(foundNodes []FoundNode) Bucket {
	b := make(Bucket,len(foundNodes))
	for i := 0; i < len(foundNodes); i++ {
			b[i] = new(Contact)
			b[i].IPAddr = foundNodes[i].IPAddr
			b[i].NodeID = foundNodes[i].NodeID
			b[i].Port = foundNodes[i].Port
	}
	return b
}

func copyData(data []byte) (ret []byte) {
  ret = make([]byte,len(data))
    for i := 0; i < len(data); i++ {
        ret[i] = data[i]
    }
    return
}

func get_host_port(c *Contact) string {
	if c == nil {
		return ""
	}
	hostPort := make([]string,2)
	hostPort[0] = c.IPAddr
	hostPort[1] = strconv.FormatUint(uint64(c.Port),10)
	hostPortStr := strings.Join(hostPort, ":")
	return hostPortStr
}

func (kadem *Kademlia) addContactToBuckets(node *Contact) int {

    _, idx := kadem.getBucket(node.NodeID)
    if node.NodeID == kadem.ThisContact.NodeID {
      // we don't want to add our own contact info to buckets, so return
      return 0
    }
    //frees up first
    kadem.next_open_spot(idx)
    kadem.K_Buckets[idx][0] = node
    
    return 0
}

func (k *Kademlia) next_open_spot(b_num int) {
	b := k.K_Buckets[b_num]
	open_spot := -1
	b_len := len(b)
	//fmt.Printf("Looking for next open spot in bucket %v\n", b_num)
	if b[0] ==nil{
		return
	}
	for i:=1;i<b_len;i++{
		if b[i]==nil{
			open_spot=i
			//fmt.Printf("Open spot at %v\n",i)
			break
		}
	}
	//if open_spot==-1, list is full
	//so pop last entry(which is really the first) and shift list one spot to the right
	if open_spot==-1{
		//fmt.Printf("Popping %v\n", b[b_len-1])
		b[b_len-1] = nil //make last entry nil
		//shift list
		for i:=b_len-2;i>0;i--{
			b[i+1] = b[i]
		}
		b[0] = nil
		
	} else{
	//else, shift list over one, with last entry at open_spot-1
	//shift 0 to openspot -1 to 1 to openspot
		for i:=open_spot;i>0;i--{
			//fmt.Printf("moving %v to %v\n",i-1,i)
			//fmt.Printf("Values: %v\n %v\n",b[i],b[i-1])
			b[i] = b[i-1]
		}
		b[0]=nil
		return
	}
}

  func (kadem *Kademlia) getBucket(dist ID) (Bucket,int) {
    bucketNum := dist.PrefixLen()-1
    if bucketNum == -1 {
      return kadem.K_Buckets[0],0
    }
    return kadem.K_Buckets[bucketNum], bucketNum
  }
  
  func next_open_spot(b Bucket) {
    open_spot := -1
    b_len := len(b)
    if b[0] ==nil{
      return
    }
    for i:=1;i<b_len;i++{
      if b[i]==nil{
        open_spot=i
        break
      }
    }
    //if open_spot==-1, list is full
    //so pop last entry(which is really the first) and shift list one spot to the right
    if open_spot==-1{
      b[b_len-1] = nil //make last entry nil
      //shift list
      for i:=b_len-2;i>0;i--{
        b[i+1] = b[i]
      }
      b[0] = nil
      
    }
    //else, shift list over one, with last entry at open_spot-1
    //shift 0 to openspot -1 to 1 to openspot
    for i:=open_spot-1;i>0;i--{
      b[i+1] = b[i]
    }
    b[0]=nil
    return
  }
  /*
   [a][ ][ ]
   [b][a][ ]
   [c][b][a]
   [c][b][ ]
   [ ][c][b]
   ^
   [d][c][b]
   
   */

  // interface to allow for sorting within buckets
  func (s Bucket) Len() int      { return len(s) }
  func (s Bucket) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
  
  // BucketSort_ByNodeID implements sort.Interface by providing Less and using the Len and
  // Swap methods of the embedded Organs value.
  type BucketSort_ByNodeID struct{ Bucket }
  
  func (s BucketSort_ByNodeID) Less(i, j int) bool {
    if s.Bucket[i] == nil {
      return false //nil's go at the end
    } else if s.Bucket[j] == nil {
      return true
    }
    return s.Bucket[i].NodeID.Less(s.Bucket[j].NodeID)
  }