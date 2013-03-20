package kademlia

/*
 serialization.go
 Contains convenience functions to serialize and deserialize our data structures
 */

import (
    "encoding/gob"
    "bytes"
    "log"
)

func (fi *FileInfo) Serialize() []byte {
    var buf bytes.Buffer
    
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(fi)
    if err != nil {
        log.Fatal("encode error:", err)
    }
    return buf.Bytes()
}

func (fi *FileInfo) Deserialize(b []byte) {
    buf := bytes.NewBuffer(b)
    
    dec := gob.NewDecoder(buf)
    err := dec.Decode(fi)
    if err != nil {
        log.Fatal("decode error:", err)
    }
}

func (dir *Directory) Serialize() []byte {
    var buf bytes.Buffer
        
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(dir)
    if err != nil {
        log.Fatal("encode error:", err)
    }
    return buf.Bytes()
}
    
func (dir *Directory) Deserialize(b []byte) {
    buf := bytes.NewBuffer(b)
    
    dec := gob.NewDecoder(buf)
    err := dec.Decode(dir)
    if err != nil {
        log.Fatal("decode error:", err)
    }
}

func (p *Packet) Serialize() []byte {
    var buf bytes.Buffer
    
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(p)
    if err != nil {
        log.Fatal("encode error:", err)
    }
    return buf.Bytes()
}