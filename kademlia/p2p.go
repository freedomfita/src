package kademlia

import (
    "os"
    "log"
)

var gShareDirectory string = "../../sharing"

//! @todo: we shouldn't really be keeping all of the packets in memory as this is too expensive. 
/// Instead, we should load them when a request is made for the file
func InitP2P(localNode *Kademlia) {
    
    localNode.ShareDir = gShareDirectory
    dir,err := os.Open(localNode.ShareDir)
    if err != nil {
        log.Fatal("Open: ", err)
    }
    f_infos,err := dir.Readdir(0)
    if err != nil {
        log.Fatal("Readdir: ", err)
    }
    localNode.FileHeaders = make(map[ID]FileHeader,len(f_infos))
    for i := 0; i < len(f_infos); i++ {
        fi := f_infos[i]
        if fi.IsDir() {
            //! @todo: do some sort of error handling here
            log.Printf("Warning: directory encountered within share directory. Sharing of directories is not currently supported.")
        } else {
            var fh FileHeader
            
            fh.FileName = fi.Name()
            fh.FilePath = localNode.ShareDir + "/" + fi.Name()
            fh.Complete = true
            fh.Info.FileSize = int(fi.Size())
            // split the file represented by fh into packets, which are (currently) stored in process memory only
            fh.LoadPackets()
            
            fh.Info.FileID = sha1hash(fh.FileName)
            fh.Info.PacketIDs = make([]ID,len(fh.Packets))
            
            for p := 0; p < len(fh.Packets); p++ {
                // generate a packet ID for each packet based on its SHA 1 hash
                fh.Info.PacketIDs[p] = fh.Packets[p].sha1hash()
            }
            // add the file header to the list of file headers
            localNode.FileHeaders[fh.Info.FileID] = fh
        }
    }
}