package main

import (
  "log"
  "os"
  "reflect"
)

type FileStateOS struct {
  Vol    uint32  `json:"vol,omitempty"`
  IdxHi  uint32  `json:"idxhi,omitempty"`
  IdxLo  uint32  `json:"idxlo,omitempty"`
}

func (fs *FileStateOS) PopulateFileIds(info os.FileInfo) {
  // For information on the following, see Go source: src/pkg/os/types_windows.go
  // This is the only way we can get at the idxhi and idxlo
  // Unix it is much easier as syscall.Stat_t is exposed and os.FileInfo interface has a Sys() method to get a syscall.Stat_t
  // Unfortunately, the relevant Windows information is in a private struct so we have to dig inside

  // NOTE: This WILL be prone to break if Go source changes, but I'd rather just fix it if it does or make it fail gracefully

  // info is os.FileInfo which is an interface to a
  // - *os.fileStat (holding methods) which is a pointer to a
  // - os.fileStat (holding data)

  // Ensure that the numbers are loaded by calling os.SameFile
  // os.SameFile will call sameFile (types_windows.go) which will call *os.fileStat's loadFileId
  // Reflection panics if we try to call loadFileId directly as its a hidden method; regardless this is much safer and more reliable
  os.SameFile(info, info)

  // If any of the following fails, report the library has changed and recover and return 0s
  defer func() {
    if r := recover(); r != nil {
      log.Printf("WARNING: File rotations that occur while LogStash Forwarder is not running will NOT be detected due to an incompatible change to the Go library used for compiling. This is a bug, please report it.\n")
      fs.Vol = 0
      fs.IdxHi = 0
      fs.IdxLo = 0
    }
  }()

  // Following makes fstat hold os.fileStat
  fstat := reflect.ValueOf(info).Elem().Elem()

  // To get the data, we need the os.fileStat that fstat points to, so one more Elem()
  fs.Vol = uint32(fstat.FieldByName("vol").Uint())
  fs.IdxHi = uint32(fstat.FieldByName("idxhi").Uint())
  fs.IdxLo = uint32(fstat.FieldByName("idxlo").Uint())
}

func (fs *FileStateOS) SameAs(info os.FileInfo) bool {
  state := &FileStateOS{}
  state.PopulateFileIds(info)
  return (fs.Vol == state.Vol && fs.IdxHi == state.IdxHi && fs.IdxLo == state.IdxLo)
}
