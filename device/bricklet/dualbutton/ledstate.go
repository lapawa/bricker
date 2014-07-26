// Copyright 2014 Dirk Jablonowski. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dualbutton

import (
	"fmt"
	"github.com/dirkjabl/bricker"
	"github.com/dirkjabl/bricker/device"
	"github.com/dirkjabl/bricker/net/packet"
	"github.com/dirkjabl/bricker/subscription"
	"github.com/dirkjabl/bricker/util/hash"
)

/*
SetLedState creates a subscriber to set the led states.
With auto toggle the led is switched on and off (toggle) by pressing the button.
For setting only one led state use SetSelectedLedState.
*/
func SetLedState(id string, uid uint32, ls *LedState, handler func(device.Resulter, error)) *device.Device {
	fid := function_set_led_state
	sls := device.New(device.FallbackId(id, "SetLedState"))
	p := packet.NewSimpleHeaderPayload(uid, fid, true, ls)
	sub := subscription.New(hash.ChoosenFunctionIDUid, uid, fid, p, false)
	sls.SetSubscription(sub)
	sls.SetResult(&device.EmptyResult{})
	sls.SetHandler(handler)
	return sls
}

// SetLedStateFuture is a future pattern version for a synchronized call of the subscriber.
// If an error occur, the result is false.
func SetLedStateFuture(brick *bricker.Bricker, connectorname string, uid uint32, ls *LedState) bool {
	future := make(chan bool)
	defer close(future)
	sub := SetLedState("setledstatefuture"+device.GenId(), uid, ls,
		func(r device.Resulter, err error) {
			future <- device.IsEmptyResultOk(r, err)
		})
	err := brick.Subscribe(sub, connectorname)
	if err != nil {
		return false
	}
	return <-future
}

// GetLedState creates the subscriber to get the led states.
func GetLedState(id string, uid uint32, handler func(device.Resulter, error)) *device.Device {
	fid := function_get_led_state
	gls := device.New(device.FallbackId(id, "GetLedState"))
	p := packet.NewSimpleHeaderOnly(uid, fid, true)
	sub := subscription.New(hash.ChoosenFunctionIDUid, uid, fid, p, false)
	gls.SetSubscription(sub)
	gls.SetResult(&LedState{})
	gls.SetHandler(handler)
	return gls
}

// GetLedStateFuture is a future pattern version for a synchronized all of the subscriber.
// If an error occur, the result is nil.
func GetLedStateFuture(brick *bricker.Bricker, connectorname string, uid uint32) *LedState {
	future := make(chan *LedState)
	defer close(future)
	sub := GetLedState("getledstatefuture"+device.GenId(), uid,
		func(r device.Resulter, err error) {
			var v *LedState = nil
			if err == nil {
				if value, ok := r.(*LedState); ok {
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

// SetSelectedLedState creates a subscriber for setting a selected led state.
// The other led stays untouched.
func SetSelectedLedState(id string, uid uint32, sls *SelectedLedState, handler func(device.Resulter, error)) *device.Device {
	fid := function_set_selected_led_state
	ssls := device.New(device.FallbackId(id, "SetSelectedLedState"))
	p := packet.NewSimpleHeaderPayload(uid, fid, true, sls)
	sub := subscription.New(hash.ChoosenFunctionIDUid, uid, fid, p, false)
	ssls.SetSubscription(sub)
	ssls.SetResult(&device.EmptyResult{})
	ssls.SetHandler(handler)
	return ssls
}

// SetSelectedLedStateFuture is a future pattern version for a synchronized call of the subscriber.
// If an error occur, the result is false.
func SetSelectedLedStateFuture(brick *bricker.Bricker, connectorname string, uid uint32, sls *SelectedLedState) bool {
	future := make(chan bool)
	defer close(future)
	sub := SetSelectedLedState("setselectedledstatefuture"+device.GenId(), uid, sls,
		func(r device.Resulter, err error) {
			future <- device.IsEmptyResultOk(r, err)
		})
	err := brick.Subscribe(sub, connectorname)
	if err != nil {
		return false
	}
	return <-future
}

// StateChanged creates a subscriber for the state changed callback.
// This callback is triggered whenever a button is pressed.
func StateChanged(id string, uid uint32, handler func(device.Resulter, error)) *device.Device {
	fid := callback_state_changed
	sc := device.New(device.FallbackId(id, "InterruptTrigger"))
	sub := subscription.New(hash.ChoosenFunctionIDUid, uid, fid, nil, true)
	sc.SetSubscription(sub)
	sc.SetResult(&States{})
	sc.SetHandler(handler)
	return sc
}

/*
Led states for the left and the right led.

    0 - activate auto toggle
    1 - deactivate auto toggle
    2 - led on (auto toggle disabled)
    3 - led off (auto toggle disabled)
*/
type LedState struct {
	LedLeft  uint8 // led state for the left led
	LedRight uint8 // led state for the right led
}

// Converts a packet to a LedState type.
func (ls *LedState) FromPacket(p *packet.Packet) error {
	if err := device.CheckForFromPacket(ls, p); err != nil {
		return err
	}
	return p.Payload.Decode(ls)
}

// String fullfill the stringer interface.
func (ls *LedState) String() string {
	txt := "Led state "
	if ls == nil {
		txt += "[nil]"
	} else {
		txt += fmt.Sprintf("[Led left: %s (%d), Led right: %s (%d)]",
			LedStateName(ls.LedLeft), ls.LedLeft,
			LedStateName(ls.LedRight), ls.LedRight)
	}
	return txt
}

// SelectedLedState is a type for select a led and set the state.
type SelectedLedState struct {
	Led   uint8 // led 0 (left) or 1 (right)
	State uint8 // led state
}

// States is the type for the state changed callback.
type States struct {
	ButtonLeft  uint8 // state from button left
	ButtonRight uint8 // state from button right
	LedLeft     uint8 // state from led left
	LedRight    uint8 // state from led right
}

// Converts a packet to a States type.
func (s *States) FromPacket(p *packet.Packet) error {
	if err := device.CheckForFromPacket(s, p); err != nil {
		return err
	}
	return p.Payload.Decode(s)
}

// String fullfill the stringer interface.
func (s *States) String() string {
	txt := "States "
	if s == nil {
		txt += "[nil]"
	} else {
		txt += fmt.Sprintf("[Button left: %s (%d), Button right: %s (%d), Led left: %s (%d), Led right: %s (%d)]",
			ButtonStateName(s.ButtonLeft), s.ButtonLeft,
			ButtonStateName(s.ButtonRight), s.ButtonRight,
			LedStateName(s.LedLeft), s.LedLeft,
			LedStateName(s.LedRight), s.LedRight)
	}
	return txt
}