// Copyright 2014 Dirk Jablonowski. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io4

import (
	"fmt"
	"github.com/dirkjabl/bricker"
	"github.com/dirkjabl/bricker/device"
	"github.com/dirkjabl/bricker/net/packet"
	"github.com/dirkjabl/bricker/subscription"
	"github.com/dirkjabl/bricker/util/hash"
)

// SetInterrupt creates the subscriber to set the interrupt bitmask.
func SetInterrupt(id string, uid uint32, i *Interrupt, handler func(device.Resulter, error)) *device.Device {
	fid := function_set_interrupt
	si := device.New(device.FallbackId(id, "SetInterrupt"))
	p := packet.NewSimpleHeaderPayload(uid, fid, true, i)
	sub := subscription.New(hash.ChoosenFunctionIDUid, uid, fid, p, false)
	si.SetSubscription(sub)
	si.SetResult(&device.EmptyResult{})
	si.SetHandler(handler)
	return si
}

// SetInterruptFuture is a future pattern version for a synchronized call of the subscriber.
// If an error occur, the result is false.
func SetInterruptFuture(brick *bricker.Bricker, connectorname string, uid uint32, i *Interrupt) bool {
	future := make(chan bool)
	defer close(future)
	sub := SetInterrupt("setinterruptfuture"+device.GenId(), uid, i,
		func(r device.Resulter, err error) {
			future <- device.IsEmptyResultOk(r, err)
		})
	err := brick.Subscribe(sub, connectorname)
	if err != nil {
		return false
	}
	return <-future
}

// GetInterrupt creates the subscriber to get the interrupt bitmask.
func GetInterrupt(id string, uid uint32, handler func(device.Resulter, error)) *device.Device {
	fid := function_get_interrupt
	gi := device.New(device.FallbackId(id, "GetInterrupt"))
	p := packet.NewSimpleHeaderOnly(uid, fid, true)
	sub := subscription.New(hash.ChoosenFunctionIDUid, uid, fid, p, false)
	gi.SetSubscription(sub)
	gi.SetResult(&Interrupt{})
	gi.SetHandler(handler)
	return gi
}

// GetInterruptFuture is a future pattern version for a synchronized all of the subscriber.
// If an error occur, the result is nil.
func GetInterruptFuture(brick *bricker.Bricker, connectorname string, uid uint32) *Interrupt {
	future := make(chan *Interrupt)
	defer close(future)
	sub := GetInterrupt("getinterruptfuture"+device.GenId(), uid,
		func(r device.Resulter, err error) {
			var v *Interrupt = nil
			if err == nil {
				if value, ok := r.(*Interrupt); ok {
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

// InterruptTrigger creates a subscriber for the interrupt callback.
// This callback is triggered whenever a change of the voltage level is detected
// on pins where the interrupt was activated with SetInterrupt.
func InterruptTrigger(id string, uid uint32, handler func(device.Resulter, error)) *device.Device {
	fid := callback_interrupt
	it := device.New(device.FallbackId(id, "InterruptTrigger"))
	sub := subscription.New(hash.ChoosenFunctionIDUid, uid, fid, nil, true)
	it.SetSubscription(sub)
	it.SetResult(&Interrupts{})
	it.SetHandler(handler)
	return it
}

// Interrupt bitmask type.
// Interrupts are triggered on changes of the voltage level of the pin,
// i.e. changes from high to low and low to high.
type Interrupt struct {
	Mask uint8 // bitmask 4bit
}

// FromPacket creates a Interrupt from a packet.
func (i *Interrupt) FromPacket(p *packet.Packet) error {
	if err := device.CheckForFromPacket(i, p); err != nil {
		return err
	}
	return p.Payload.Decode(i)
}

// String fullfill the stringer interface.
func (i *Interrupt) String() string {
	txt := "Interrupt "
	if i == nil {
		txt += "[nil]"
	} else {
		txt += fmt.Sprintf("[Mask: %d (%s)]",
			i.Mask, MaskToString(i.Mask))
	}
	return txt
}

// Interrupts is the result type of the interrupt callback.
type Interrupts struct {
	InterruptMask uint8 // bitmap 4bit
	ValueMask     uint8 // bitmap 4nit
}

// FromPacket creates a Interrupts object from a packet.
func (i *Interrupts) FromPacket(p *packet.Packet) error {
	if err := device.CheckForFromPacket(i, p); err != nil {
		return err
	}
	return p.Payload.Decode(i)
}

// String fullfill the stringer interface.
func (i *Interrupts) String() string {
	txt := "Interrupts "
	if i == nil {
		txt += "[nil]"
	} else {
		txt += fmt.Sprintf("[Interrupt Mask: %d (%s), Value Mask: %d (%s)]",
			i.InterruptMask, MaskToString(i.InterruptMask),
			i.ValueMask, MaskToString(i.ValueMask))
	}
	return txt
}