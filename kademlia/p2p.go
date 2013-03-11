package kademlia

import (
    "os"
    "log"
    "fmt"
    "strings"
)

var gShareDirectory string = "../../sharing"

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
        } else if fi.Name()[0] == 0x2e {
            log.Printf("Ignoring special file %v\n",fi.Name())
        } else {
            var fh FileHeader
            
            fh.FileName = fi.Name()
            fh.FilePath = localNode.ShareDir + "/" + fi.Name()
            fh.Info.Complete = true
            fh.Info.FileSize = fi.Size()
            // split the file represented by fh into persistent packets, if this has not already happened
            if !fh.PacketsExist() {
                fh.MakePackets()
            }
            //! @todo: we need to load the packets when a request for the associated File Header is made
            fh.PacketsLoaded = false
            fh.Info.FileID = sha1hash(fh.FileName)
            
            // add the file header to the list of file headers
            localNode.FileHeaders[fh.Info.FileID] = fh
            localNode.Data[fh.Info.FileID] = fh.Info.Serialize()
        }
    }
}

func DownloadFile(fname string) {
    // calculate file ID
    fid := sha1hash(fname)
    // get the file header
    result := IterativeFindValue(fid)
    var fi FileInfo
    fi.Deserialize(result)
    if result == nil {
        fmt.Printf("Error: could not find file associated with key %v\n",fid)
        return
    }
    fi.Complete = false
    var fh FileHeader
    fh.FileName = fname
    fh.FilePath = gShareDirectory + "/" + fname
    nameparts := strings.Split(fname, ".")
    fh.PacketDir = gShareDirectory + "/" + strings.Join(nameparts[:len(nameparts)-1],"")
    fh.Info = fi
    // load existing packets
    
    doneChannel := make(chan int, len(fi.PacketIDs))
    for pIdx := 0; pIdx < len(fi.PacketIDs); pIdx++ {
        pid := fi.PacketIDs[pIdx]
        go GetPacket(fh,pid,doneChannel,pIdx)
    }
    // clear the channel
    for i := 0; i < len(fi.PacketIDs); i++ {
        <- doneChannel
    }
    // at this point, all the GetPacket goroutines will have finished downloading
    fh.JoinPackets()
    fh.Info.Complete = true
    fmt.Printf("File with key %v downloaded successfully.\n",fid)
    return
}

func GetPacket(fh FileHeader, packetID ID, doneChannel chan int,pnum int) {
    // first, we check if we already have the packet
    if _,ok := fh.Packets[packetID]; ok {
        doneChannel <- 1
        return
    }
    // if we don't have the packet, get it from the network
    packetData := IterativeFindValue(packetID)
    // deserialize the raw byte stream
    var packet Packet
    packet.Deserialize(packetData)
    if packet.PacketID != packetID {
        log.Fatal("ERROR: Packet IDs do not match.\n")
    }
    // add packet to the map in the file header
    fh.Packets[packetID] = packet
    // write the packet to disk
    packet.Write(fh.PacketDir,pnum)
    // signal the channel when we are done downloading
    doneChannel <- 1
}