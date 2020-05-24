package main

import (
	"strings"
	"time"
)

// TraceEvent is the event that git_trace2 logs use
type TraceEvent string

// EVENT Format
// https://git-scm.com/docs/api-trace2#_event_format
const (
	TraceEventVersion TraceEvent = "version"
	TraceEventDiscard TraceEvent = "discard"

	TraceEventStart  TraceEvent = "start"
	TraceEventExit   TraceEvent = "exit"
	TraceEventAtExit TraceEvent = "atexit"

	TraceEventSignal TraceEvent = "signal"
	TraceEventError  TraceEvent = "error"

	TraceEventCommandPath TraceEvent = "cmd_path"
	TraceEventCommandName TraceEvent = "cmd_name"
	TraceEventCommandMode TraceEvent = "cmd_mode"

	TraceEventAlias TraceEvent = "alias"

	TraceEventChildStart TraceEvent = "child_start"
	TraceEventChildExit  TraceEvent = "child_exit"

	TraceEventExec       TraceEvent = "exec"
	TraceEventExecResult TraceEvent = "exec_result"

	TraceEventThreadStart TraceEvent = "thread_start"
	TraceEventThreadExit  TraceEvent = "thread_exit"

	TraceEventDefParam TraceEvent = "def_param"
	TraceEventDefRepo  TraceEvent = "def_repo"

	TraceEventRegionEnter TraceEvent = "region_enter"
	TraceEventRegionLeave TraceEvent = "region_leave"

	TraceEventData     TraceEvent = "data"
	TraceEventDataJSON TraceEvent = "data_JSON"
)

/*
GitTrace2Event typical format for git_trace2
Extracted json logs using
  `tail -200 git_trace.log | jq --slurp '.' | pbcopy`
then pasted into  https://mholt.github.io/json-to-go/
to have the struct generated dynamically.

Note that this should be run accross multiple times to get all structs fields
and struct fields may differ depends on git version.

Meaning that if you have multiple git version, you should sample the struct separately

Git trace2 is modified almost every 2 minor version of git
*/
type GitTrace2Event struct {
	Event     TraceEvent `json:"event"`
	Evt       *string    `json:"evt,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Hierarchy *string    `json:"hierarchy,omitempty"`

	Sid SessionID `json:"sid"`

	Thread string     `json:"thread"`
	Time   *time.Time `json:"time"`
	File   string     `json:"file"`
	Line   int        `json:"line"`

	TimeRelative *float64 `json:"t_rel,omitempty"`
	TimeAbsolute *float64 `json:"t_abs,omitempty"`

	ArgumentVector []string `json:"argv,omitempty"`
	Code           *int     `json:"code,omitempty"`

	Repo       *int    `json:"repo,omitempty"`
	Nesting    *int    `json:"nesting,omitempty"`
	Category   *string `json:"category,omitempty"`
	Label      *string `json:"label,omitempty"`
	Key        *string `json:"key,omitempty"`
	Value      *string `json:"value,omitempty"`
	ChildID    *int    `json:"child_id,omitempty"`
	ProcessID  *int    `json:"pid,omitempty"`
	ChildClass *string `json:"child_class,omitempty"`
	UseShell   *bool   `json:"use_shell,omitempty"`

	GitVersion *string `json:"exe,omitempty"`
}

// IsParentEvent check if the event is on parent Session
func (e *GitTrace2Event) IsParentEvent() bool {
	return e.Sid.getChildSessionID() == ""
}

// IsChildEvent check if the event is on child Session
func (e *GitTrace2Event) IsChildEvent() bool {
	return e.Sid.getChildSessionID() != ""
}

// IsCompletedSession checked if the slice of Events represent a completed Session
func IsCompletedSession(events []*GitTrace2Event) bool {
	count := len(events)
	if count == 0 {
		return false
	}

	hasFirstEvent := false
	hasLastEvent := false

	for _, event := range events {
		if hasFirstEvent && hasLastEvent {
			return true
		}

		if event.IsParentEvent() && event.Event == TraceEventVersion {
			hasFirstEvent = true
			continue
		}

		if event.IsParentEvent() && event.Event == TraceEventAtExit {
			hasLastEvent = true
			continue
		}
	}

	return hasFirstEvent && hasLastEvent
}

// SessionID is one or two session IDs,
// parent session and child session,
// concated and separated by '/'.
// Each session id is composed of a timestamp and some traces followed
type SessionID string

func (s *SessionID) getParentSessionID() string {
	a := strings.Split(string(*s), "/")
	return a[0]
}

func (s *SessionID) getChildSessionID() string {
	a := strings.Split(string(*s), "/")
	if len(a) > 1 {
		return a[len(a)-1]
	}
	return ""
}
