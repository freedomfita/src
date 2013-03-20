package kademlia

import (
    "os"
    "strconv"
    "io"
    "io/ioutil"
    "log"
    "crypto/sha1"
    "encoding/hex"
    "strings"
)

var gPacketSize int = 16384 // == 16kb 

type FileInfo struct {
    FileSize int64
    FileName string
    FileID ID
    Complete bool
    PacketIDs []ID
    
}

type FileHeader struct {
    Info FileInfo
    FilePath string
    PacketsLoaded bool
    PacketDir string
    Packets map[ID]Packet
    UpdateNodes []ID
}

type Packet struct {
    PacketID ID
    PacketSize int
    Data []byte
}

type DirInfo struct {
    DirName string
    DirID ID
    Checksum ID
}

type Directory struct {
    Info DirInfo
    //Parent DirInfo
    ChildDirs []DirInfo
    Files []FileInfo
}

// get info for a file from the network
func GetFileInfo(fname string) (fi FileInfo) {
    return fi
}

// check if a file has already been split into packets
func (f *FileHeader) PacketsExist() bool {
    return false
    _,err := os.Stat(f.FilePath + ".1")
    if err == nil {
        log.Printf("Packets exist.")
        return true
    }
    return !os.IsNotExist(err)
}
// split a single file into packets
func (f *FileHeader) MakePackets() {
    // calculate number of packets
    numPackets := int(f.Info.FileSize)/gPacketSize
    
    // read the file into a byte array
    fbytes, err := ioutil.ReadFile(gShareDirectory+f.FilePath)
    if err != nil {
        log.Fatal("ReadFile: ",err)
    }
    
    // save the current directory and switch to sharing dir. for this file
    curDir,_ := os.Getwd()
    nameparts := strings.Split(f.Info.FileName, ".")
    f.PacketDir = gDataDirectory + "/" + strings.Join(nameparts[:len(nameparts)-1],"")
    os.Mkdir(f.PacketDir,os.ModeDir | 0x1ff)
    os.Chdir(f.PacketDir)
    
    // fill each packet with the appropriate data
    i := 0
    for ; i < numPackets; i++ {
        data := fbytes[i*gPacketSize:(i+1)*gPacketSize]
        fname := f.Info.FileName + "." + strconv.Itoa(i+1)
        file, err := os.Create(fname)
        if err != nil {
            log.Fatal("Create: ",err)
        }
        file.Write(data)
        file.Close()
    }
    if int(f.Info.FileSize) % gPacketSize != 0 {
        data := fbytes[i*gPacketSize:]
        fname := f.Info.FileName + "." + strconv.Itoa(i+1)
        file, err := os.Create(fname)
        if err != nil {
            log.Fatal("Create: ",err)
        }
        file.Write(data)
        file.Close()
    }
    // restore cwd
    os.Chdir(curDir)
    return
}

// load the packets for a file
func (f *FileHeader) LoadPackets() {
    curDir,_ := os.Getwd()
    // calculate number of packets
    numPackets := int(f.Info.FileSize)/gPacketSize
    
    f.Packets = make(map[ID]Packet,numPackets+1)
    // open the directory containing the packet files
    dir,err := os.Open(gShareDirectory+"/"+f.PacketDir)
    if err != nil {
        log.Fatal("Open: ", err)
    }
    packet_files,err := dir.Readdir(0)
    if err != nil {
        log.Fatal("Readdir: ", err)
    }
    os.Chdir(f.PacketDir)
    // fill each packet with the appropriate data
    for i := 0; i < len(packet_files); i++ {
        fi := packet_files[i]

        var p Packet
        // read the file into a byte array
        fbytes, err := ioutil.ReadFile(fi.Name())
        if err != nil {
            log.Fatal("ReadFile: ",err)
        }
        
        copy(p.Data,fbytes)
        p.PacketID = p.sha1hash()
        p.PacketSize = len(fbytes)
        f.Packets[p.PacketID] = p
    }
    os.Chdir(curDir)
    return
}

// join all packets into a single file
func (f *FileHeader) JoinPackets(loc string) {
    //curDir,_ := os.Getwd()
    
    file,err := os.Create(loc+f.Info.FileName)
    if err != nil {
        log.Fatal("os.Create: ",err)
    }
    
    dir,err := os.Open(f.PacketDir)
    if err != nil {
        log.Fatal("Open: ", err)
    }
    packet_files,err := dir.Readdir(0)
    if err != nil {
        log.Fatal("Readdir: ", err)
    }
    
    //os.Chdir(f.PacketDir)
    
    var pos int64 = 0
    for i := 0; i < len(packet_files); i++ {
        fi := packet_files[i]
        if !fi.IsDir() {
            fbytes, err := ioutil.ReadFile(f.PacketDir+"/"+fi.Name())
            if err != nil {
                log.Fatal("ReadFile in JoinPackets: ",err)
            }
            increment,_ := file.WriteAt(fbytes,pos)
            pos += int64(increment)
        }
    }
    file.Close()
    //os.Chdir(curDir)
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
    io.WriteString(shaGen,hex.EncodeToString(p.Data))
    return FromByteArray(shaGen.Sum(nil))
}

func (p *Packet) Deserialize(data []byte) {
    p.PacketSize = len(data)
    copy(p.Data,data)
    p.PacketID = p.sha1hash()
}

func (p *Packet) Write(dir string, pnum int) {
    file,err := os.Create(dir+"/"+strconv.Itoa(pnum))
    if err != nil {
        log.Fatal("os.Create: ",err)
    }
    _,err = file.Write(p.Data)
    if err != nil {
        log.Fatal("file.Write: ",err)
    }
    file.Close()
}