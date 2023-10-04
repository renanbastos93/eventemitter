package eventemitter

import (
	"fmt"
)

// Event is an interface to ensure contract who go to implements it
type Event interface {
	Reset()
	Listen(eventName string, actions ...Actions)
	ListenOnce(eventName string, actions ...Actions)
	Emit(eventName string, param interface{})
	RemoveEvent(evetName string)
}

// Flags are a new type of inheritance from int
type Flags int

// Actions is a new type to define the contract for the function "actions that will be called"
type Actions func(param interface{})

// Topic is a struct to set all actions with your flags according to call when emitting an event from topic
type Topic struct {
	Actions []Actions
	Flags   Flags
	Called  bool
}

// EventActions is a new type using a map with a list as a Topic struct
type EventActions map[string]Topic

// implEvent implements interface Event
type implEvent struct {
	events EventActions
}

// Create all flags
const (
	FlagNone Flags = 0
	FlagOnce Flags = 1 << iota
)

// eventemitter is a global variable to reuse in methods that don't need to call the `New` method
var eventemitter Event = New()

// New is a method to instance `implEvent` using interface `Event`
func New() Event {
	return &implEvent{
		events: EventActions{},
	}
}

// Reset must remove all listeners of the topics
func (e *implEvent) Reset() {
	for eventName := range e.events {
		if topic, ok := e.events[eventName]; ok {
			topic.Actions = []Actions{}
			e.events[eventName] = topic
		}
	}
}

// Listen will add actions for the topic
func (e *implEvent) Listen(eventName string, actions ...Actions) {
	if topic, ok := e.events[eventName]; ok {
		topic.Actions = append(topic.Actions, actions...)
		e.events[eventName] = topic
	} else {
		e.events[eventName] = Topic{
			Actions: actions,
			Flags:   FlagNone,
			Called:  false,
		}
	}
}

// ListenOnce is the same `Listen` method but it will work once time
func (e *implEvent) ListenOnce(eventName string, actions ...Actions) {
	e.events[eventName] = Topic{
		Actions: actions,
		Flags:   FlagOnce,
		Called:  false,
	}
}

// tryRunAction will call to recover case occurs a panic
func tryRunAction(fn Actions, data interface{}) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("Recovered: %+v\n", r)
		}
	}()
	fn(data)
}

// Emit must be called by events triggering your listeners
func (e *implEvent) Emit(eventName string, data interface{}) {
	if topic, ok := e.events[eventName]; !ok {
		return
	} else {
		if topic.Flags == FlagOnce && topic.Called {
			return
		}

		for i := range topic.Actions {
			tryRunAction(topic.Actions[i], data)
			topic.Called = true
			e.events[eventName] = topic
		}
	}
}

// RemoveEvent must removed the listeners from Topic
func (e *implEvent) RemoveEvent(eventName string) {
	if _, ok := e.events[eventName]; ok {
		delete(e.events, eventName)
	}
}

// Reset is a wrapper from Event.Reset()
func Reset() {
	eventemitter.Reset()
}

// Listen is a wrapper from Event.Listen(eventName string, action Actions)
func Listen(eventName string, action Actions) {
	eventemitter.Listen(eventName, action)
}

// Emit is a wrapper from Event.Emit(eventName string, param interface{})
func Emit(eventName string, param interface{}) {
	eventemitter.Emit(eventName, param)
}

// RemoveEvent is a wrapper from Event.RemoveEvent(eventName string)
func RemoveEvent(evetName string) {
	eventemitter.RemoveEvent(evetName)
}

// ListenOnce is a wrapper from Event.ListenOnce(eventName string, action Actions)
func ListenOnce(eventName string, action Actions) {
	eventemitter.ListenOnce(eventName, action)
}
