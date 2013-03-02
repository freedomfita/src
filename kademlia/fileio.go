package kademlia

import (
    "io"
    "log"
    "crypto/sha1"
    "encoding/hex"
)

var gPacketSize int = 16384 // == 16kb 

type FileInfo struct {
    FileSize int
    FileID ID
    PacketIDs []ID
}

type FileHeader struct {
    Info FileInfo
    FileName string
    FilePath string
    Complete bool
    Packets []Packet
}

type Packet struct {
    PacketID ID
    PacketSize int
    Data []byte
}

// get info for a file from the network
func GetFileInfo(fname string) (fi FileInfo) {
    return fi
}

// split a single file into packets
func (f *FileHeader) LoadPackets() {
    // calculate number of packets
    numPackets := (f.FileInfo.FileSize)/gPacketSize
    if f.FileInfo.FileSize % gPacketSize != 0 {
        numPackets += 1
    }
    // fill each packet with the appropriate data
    for i := 0; i < numPackets; i++ {
    
    }
    
    return
}

// join all packets into a single file
func (f *FileHeader) JoinPackets() {
    return
}

// hash a string or Packet struct to a 160-bit ID
func sha1hash(s string) ID {
    // used to generate file and packet IDs later
    shaGen := sha1.New()
    io.WriteString(shaGen,s)
    return FromByteArray(shaGen.Sum(nil))
}

func (p *Packet) sha1hash() ID {
    // used to generate file and packet IDs later
    shaGen := sha1.New()
    log.Printf("%v\n",p.Data)
    io.WriteString(shaGen,hex.EncodeToString(p.Data))
    return FromByteArray(shaGen.Sum(nil))
}