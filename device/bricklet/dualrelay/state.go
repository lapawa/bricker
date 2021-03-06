// Copyright 2014 Dirk Jablonowski. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dualrelay

import (
	"fmt"
	"github.com/dirkjabl/bricker"
	"github.com/dirkjabl/bricker/device"
	"github.com/dirkjabl/bricker/net/packet"
	misc "github.com/dirkjabl/bricker/util/miscellaneous"
)

// SetState creates a subscriber to set the dual relays.
// This subscriber set all two relays at once.
//
// If you do not know the state of one of the relays, you can read the states with GetState or
// use for setting SetSelectedState.
//
// Default state is &State{false, false}.
func SetState(id string, uid uint32, s *State, handler func(device.Resulter, error)) *device.Device {
	return device.Generator{
		Id:         device.FallbackId(id, "SetState"),
		Fid:        function_set_state,
		Uid:        uid,
		Data:       NewStateRaw(s),
		Handler:    handler,
		WithPacket: true}.CreateDevice()
}

// SetStateFuture is a future pattern version for a synchronized call of the subscriber.
// If an error occur, the result is false.
func SetStateFuture(brick bricker.Bricker, connectorname string, uid uint32, s *State) bool {
	future := make(chan bool)
	defer close(future)
	sub := SetState("setstatefuture"+device.GenId(), uid, s,
		func(r device.Resulter, err error) {
			future <- device.IsEmptyResultOk(r, err)
		})
	err := brick.Subscribe(sub, connectorname)
	if err != nil {
		return false
	}
	return <-future
}

// GetState creates a subscriber to get the relay states.
func GetState(id string, uid uint32, handler func(device.Resulter, error)) *device.Device {
	return device.Generator{
		Id:         device.FallbackId(id, "GetState"),
		Fid:        function_get_state,
		Uid:        uid,
		Result:     &State{},
		Handler:    handler,
		WithPacket: true}.CreateDevice()
}

// GetStateFuture is a future pattern version for a synchronized call of the subscriber.
// If an error occur, the result is nil.
func GetStateFuture(brick bricker.Bricker, connectorname string, uid uint32) *State {
	future := make(chan *State)
	defer close(future)
	sub := GetState("getstatefuture"+device.GenId(), uid,
		func(r device.Resulter, err error) {
			var v *State = nil
			if err == nil {
				if value, ok := r.(*State); ok {
					v = value
				}
			}
			future <- v
		})
	err := brick.Subscribe(sub, connectorname)
	if err != nil {
		return nil
	}
	return <-future
}

// SetSelectedState creates a subscriber to set only one relay.
// The not seleced relay remains untouched.
func SetSelectedState(id string, uid uint32, s *SelectedState, handler func(device.Resulter, error)) *device.Device {
	return device.Generator{
		Id:         device.FallbackId(id, "SetSelectedState"),
		Fid:        function_set_selected_state,
		Uid:        uid,
		Data:       NewSelectedStateRaw(s),
		Handler:    handler,
		WithPacket: true}.CreateDevice()
}

// SetSelectedStateFuture is a future pattern version for a synchronized call of the subscriber.
// If an error occur, the result is false.
func SetSelectedStateFuture(brick bricker.Bricker, connectorname string, uid uint32, s *SelectedState) bool {
	future := make(chan bool)
	defer close(future)
	sub := SetSelectedState("setselectedstatefuture"+device.GenId(), uid, s,
		func(r device.Resulter, err error) {
			future <- device.IsEmptyResultOk(r, err)
		})
	err := brick.Subscribe(sub, connectorname)
	if err != nil {
		return false
	}
	return <-future
}

// State holds the state of the relays.
// If a relay is on than the state is true (off/false).
type State struct {
	Relay1 bool
	Relay2 bool
}

// FromPacket converts the payload data from a packet into a State type object.
func (s *State) FromPacket(p *packet.Packet) error {
	if err := device.CheckForFromPacket(s, p); err != nil {
		return err
	}
	sr := new(StateRaw)
	err := p.Payload.Decode(sr)
	if err == nil && sr != nil {
		s.FromStateRaw(sr)
	}
	return err
}

// String fullfill the stringer interface
func (s *State) String() string {
	txt := "State "
	if s == nil {
		txt += "[nil]"
	} else {
		txt += fmt.Sprintf("[Relay1: %t, Relay2: %t]", s.Relay1, s.Relay2)
	}
	return txt
}

// Copy creates a copy of the content.
func (s *State) Copy() device.Resulter {
	if s == nil {
		return nil
	}
	return &State{
		Relay1: s.Relay1,
		Relay2: s.Relay2}
}

// FromStateRaw converts a StateRaw to a State
func (s *State) FromStateRaw(sr *StateRaw) {
	if s == nil || sr == nil {
		return
	}
	s.Relay1 = misc.Uint8ToBool(sr.Relay1)
	s.Relay2 = misc.Uint8ToBool(sr.Relay2)
}

// StateRaw is a de/encoding type for State.
type StateRaw struct {
	Relay1 uint8
	Relay2 uint8
}

// Creates a new StateRaw from a State.
func NewStateRaw(s *State) *StateRaw {
	if s == nil {
		return nil
	}
	sr := new(StateRaw)
	sr.Relay1 = misc.BoolToUint8(s.Relay1)
	sr.Relay2 = misc.BoolToUint8(s.Relay2)
	return sr
}

// SelectedState is a type to address one specific relay (1 or 2).
//
// Relay could be 1 or 2.
// If State is true this means on (false/off).
type SelectedState struct {
	Relay uint8
	State bool
}

// FromPacket converts the payload data from a packet into a SelectedState type object.
func (s *SelectedState) FromPacket(p *packet.Packet) error {
	if err := device.CheckForFromPacket(s, p); err != nil {
		return err
	}
	sr := new(SelectedStateRaw)
	err := p.Payload.Decode(sr)
	if err == nil && sr != nil {
		s.FromSelectedStateRaw(sr)
	}
	return err
}

// String fullfill the stringer interface.
func (s *SelectedState) String() string {
	txt := "Selected State "
	if s == nil {
		txt += "[nil]"
	} else {
		txt += fmt.Sprintf("[Relay: %d, State: %t]", s.Relay, s.State)
	}
	return txt
}

// Copy creates a copy of the content.
func (s *SelectedState) Copy() device.Resulter {
	if s == nil {
		return nil
	}
	return &SelectedState{
		Relay: s.Relay,
		State: s.State}
}

// FromSelectedStateRaw converts the contet of a SelectedStateRaw into the object.
func (s *SelectedState) FromSelectedStateRaw(sr *SelectedStateRaw) {
	if s == nil || sr == nil {
		return
	}
	s.Relay = sr.Relay
	s.State = (sr.State & 0x01) == 0x01
}

// SelectedStateRaw is a en/decoding type for SelectedState.
type SelectedStateRaw struct {
	Relay uint8
	State uint8
}

// Creates a new SelectedStateRaw from a SelectedState.
func NewSelectedStateRaw(s *SelectedState) *SelectedStateRaw {
	if s == nil {
		return nil
	}
	sr := new(SelectedStateRaw)
	sr.FromSelectedState(s)
	return sr
}

// FromState converts a State into a StateRaw.
func (sr *SelectedStateRaw) FromSelectedState(s *SelectedState) {
	if sr == nil || s == nil {
		return
	}
	sr.Relay = s.Relay
	if s.State {
		sr.State = 0x01
	} else {
		sr.State = 0x00
	}
}
