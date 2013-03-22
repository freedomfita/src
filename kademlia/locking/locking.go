package locking

import (
	"kademlia"
	"authentication"
	"strings"
)

/**********************************************************************************
* When we call locking, do locking rpc to check other nodes. If any return with a 
* "using_lock" response, we abort. Otherwise, change permissions on file from read-only
* to read/write. When we unlock, we do a "file_modified" rpc to all nodes with that file,
* then they will all attempt to download the modified file. -Not scalable
**********************************************************************************/

//file_header.Update_Nodes[] for friends who want to be notified of a file.

func (k *Kademlia) Request_Lock(f_id ID) {
	// here we loop through and do acquire_lock to each friend with the file.
	f_header := IterativeFindValue(f_id)
	u_n_len := len(f_header.Update_Nodes)
	k.Votes_Needed = u_n_len
	//now we have header, look in Update_Nodes for list of nodes to send request to
	for i:=0; i<u_n_len;i++ {
		if Update_Nodes[i] { //we have NodeID
			node := k.find_friend(Update_Nodes[i])
			s := []string{node.IPAddr, node.Port}
			go acquire_lock(strings.Join(s,""), f_id)
		}
		else{
			break
		}
	}
	//wait until we have all votes * in event a node is down, 
	//we count that as unlocked vote
	for ; k.Vote_Total != u_n_len; {
		//do nothing, we wait until all votes received. Not practical in
		//a real system by any means, but I couldn't think of a simple
		//better way to do it without worrying about timers and shit
		//which aren't hard to implement but I don't think are necessary
		//for a project of this scale. Besides, even when a node is down
		//we will still get a correct vote total. Hopefully
	}
	//Now we have all votes
	// ***LOOK HERE FOR CHMOD INFO *** 
	if k.Vote_Total == k.Votes_Acquired {
		//we have lock, change permissions on file
		//files should initially be
		//os.Chmod(f_header.FilePath, 555) 
		//read = 4
		//write = 2
		//execute = 1
		//first number = owner permission
		//second = group permissions
		//third = world permissions
		//so we set it to 755 when they have the lock.
		k.Lock_Acquired = f_id
		os.Chmod(f_header.FilePath, 755)
	}
	
}


func (k *Kademlia) Acquire_Lock(req LockRequest, *res LockResult) error {
	res.MsgID = CopyID(req.MsgID)
	res.is_locked = k.Lock_Acquired.Equal(req.FileID)
}

func (k *Kademlia) Notify_Release_Lock(f_id ID){
	k.Lock_Acquired = nil
	f_header := IterativeFindValue(f_id)
	os.Chmod(f_header.FilePath, 555)
	u_n_len := len(f_header.Update_Nodes)
	//now we have header, look in Update_Nodes for list of nodes to send request to
	for i:=0; i<u_n_len;i++ {
		if Update_Nodes[i] { //we have NodeID
			node := k.find_friend(Update_Nodes[i])
			s := []string{node.IPAddr, node.Port}
			go release_lock(strings.Join(s,""), f_id, f_header.FilePath)
		} 
		else {
			break
		}
	}

}

func (k *Kademlia) Release_Lock(req UnlockRequest, *res UnlockResult) error {
	res.MsgID = CopyID(req.MsgID)
	//call function to re-download file
	//func DownloadFile(fname string, dest string
	DownloadFile(req.FileName, req.FilePath)
	
}
