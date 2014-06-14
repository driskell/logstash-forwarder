package main

import "os"

type Event map[string]interface{};

type FileEvent struct {
  ProspectorInfo *ProspectorInfo
  Offset         int64
  Event          *Event
}

type RegistrarEvent interface {
  Process(state map[*ProspectorInfo]*FileState)
}

type NewFileEvent struct {
  ProspectorInfo *ProspectorInfo
  Source         string
  Offset         int64
  fileinfo       os.FileInfo
}

type DeletedEvent struct {
  ProspectorInfo *ProspectorInfo
}

type RenamedEvent struct {
  ProspectorInfo *ProspectorInfo
  Source         string
}

type EventsEvent struct {
  Events []*FileEvent
}

func NewEvent(fields map[string]*string, file *string, offset int64, line uint64, message *string) *Event {
  event := Event{
    "file":   file,
    "offset": offset,
    "line":   message, // The lumberjack receiver expects "line" and not "message"
  }
  for k, v := range fields {
    event[k] = v
  }
  return &event
}
