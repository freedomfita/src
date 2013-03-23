package kademlia

import (
	"strings"
)

/**********************************************************************************
* When we call locking, do locking rpc to check other nodes. If any return with a 
* "using_lock" response, we abort. Otherwise, change permissions on file from read-only
* to read/write. When we unlock, we do a "file_modified" rpc to all nodes with that file,
* then they will all attempt to download the modified file. -Not scalable
**********************************************************************************/

//file_header.Update_Nodes[] for friends who want to be notified of a file.


func (k *Kademlia) Acquire_Lock(req LockRequest, *res LockResult) error {
	res.MsgID = CopyID(req.MsgID)
	res.is_locked = k.Lock_Acquired.Equal(req.FileID)
}



func (k *Kademlia) Release_Lock(req UnlockRequest, *res UnlockResult) error {
	res.MsgID = CopyID(req.MsgID)
	//call function to re-download file
	//func DownloadFile(fname string, dest string
	DownloadFile(req.FileName, req.FilePath)
	
}
