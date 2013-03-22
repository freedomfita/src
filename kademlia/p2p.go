package kademlia

import (
    "os"
    "log"
    "fmt"
    "strings"
    "net/rpc"
)

//! TODO: allow user to set these, and set dynamically on run
var gShareDirectory string = "/Users/MoritzGellner/Desktop/Projects/Kademlia_DHT/Kademlia/sharing"
var gDataDirectory string = "/Users/MoritzGellner/Desktop/Projects/Kademlia_DHT/Kademlia/packets"

var DIRFLAGS = os.ModeDir | 0x1ff

func InitP2P(localNode *Kademlia) {
    
    ThisNode.FileHeaders = make(map[ID]FileHeader,1000)
    LoadDir(gShareDirectory)
}

func LoadDir(directory string) (ret Directory) {
    truncatedName := strings.SplitAfter(directory,gShareDirectory)[1]+"/"
    ret.Info.DirName = truncatedName
    ret.Info.DirID = sha1hash(truncatedName)
    
    dir,err := os.Open(directory)
    if err != nil {
        log.Fatal("Open: ", err)
    }
    f_infos,err := dir.Readdir(0)
    if err != nil {
        log.Fatal("Readdir: ", err)
    }
    for i := 0; i < len(f_infos); i++ {
        fi := f_infos[i]
        if fi.IsDir() {
            childdir := LoadDir(directory+"/"+fi.Name())
            ret.ChildDirs = append(ret.ChildDirs,childdir.Info)
        } else if fi.Name()[0] == 0x2e {
            log.Printf("Ignoring special file %v\n",fi.Name())
        } else {
            LoadFile(fi,&ret)
        }
    }
    ThisNode.Data[ret.Info.DirID] = ret.Serialize()
    return ret
}

func LoadFile(fi os.FileInfo, dir *Directory) {
    
    var fh FileHeader
    
    fh.Info.FileName = fi.Name()
    fh.FilePath = dir.Info.DirName + fi.Name()
    fh.Info.Complete = true
    fh.Info.FileSize = fi.Size()
    // split the file represented by fh into persistent packets, if this has not already happened
    if !fh.PacketsExist() {
        fh.MakePackets()
    }
    
    fh.PacketsLoaded = false
    fh.Info.FileID = sha1hash(fh.Info.FileName)
    
    // add the file header to the list of file headers
    ThisNode.FileHeaders[fh.Info.FileID] = fh
    ThisNode.Data[fh.Info.FileID] = fh.Info.Serialize()
    dir.Files = append(dir.Files,fh.Info)
}

func DownloadFile(fname string, dest string, wantUpdate bool) {
    // calculate file ID
    fid := sha1hash(fname)
    // get the file header
    result,node := IterativeFindValue(fid)
    // if we want to be notified when this file changes, let the node with the data know
    if wantUpdate {
        hostPortStr := get_host_port(&node)
        client, err := rpc.DialHTTP("tcp", hostPortStr)
        if err != nil {
            log.Fatal("rpc.DialHTTP:",err)
        }
        req := new(UpdateListenerRequest)
        req.MsgID = NewRandomID()
        req.FileID = fid
        req.ListenerID = ThisNode.ThisContact.NodeID
        
        var res FindValueResult
        err = client.Call("Kademlia.AddUpdateListener", req, &res)
        if err != nil {
            log.Fatal("Call: ", err)
        }
        client.Close()
    }
    var fi FileInfo
    if result == nil {
        fmt.Printf("Error: could not find file associated with key %v\n",fid)
        return
    }
    fi.Deserialize(result)
    
    fi.Complete = false
    var fh FileHeader
    fh.Info.FileName = fname
    fh.FilePath = gShareDirectory + "/" + fname
    nameparts := strings.Split(fname, ".")
    fh.PacketDir = gDataDirectory + "/" + strings.Join(nameparts[:len(nameparts)-1],"")
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
    fh.JoinPackets(dest)
    fh.Info.Complete = true
    fmt.Printf("File with key %v downloaded successfully.\n",fid)
    
    return
}

func DownloadDirectory(dirname string, dest string, updateDir string, wantUpdate bool) {
    if wantUpdate == false {
        if updateDir == dirname {
            wantUpdate = true
        }
    }
    dirid := sha1hash(dirname)
    result,_ := IterativeFindValue(dirid)
    var dir Directory
    if result == nil {
        fmt.Printf("Error: could not find directory associated with key %v\n",dirid)
        return
    }
    dir.Deserialize(result)
    
    os.Mkdir(dest+dirname,DIRFLAGS)
    
    for fi := 0; fi < len(dir.Files); fi++ {
        curFile := dir.Files[fi]
        DownloadFile(curFile.FileName,dest+dirname,wantUpdate)
    }
    
    for di := 0; di < len(dir.ChildDirs); di++ {
        curChild := dir.ChildDirs[di].DirName
        DownloadDirectory(curChild,dest,updateDir,wantUpdate)
    }
    fmt.Printf("Directory with key %v downloaded successfully.\n",dirid)
    
}

func GetPacket(fh FileHeader, packetID ID, doneChannel chan int,pnum int) {
    // first, we check if we already have the packet
    if _,ok := fh.Packets[packetID]; ok {
        doneChannel <- 1
        return
    }
    // if we don't have the packet, get it from the network
    packetData,_ := IterativeFindValue(packetID)
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