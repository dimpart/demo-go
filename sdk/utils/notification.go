/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package utils

import . "github.com/dimchat/mkm-go/types"

// Notification defines the interface for a generic notification object
//
// Carries a named event with associated sender and extra info
type Notification interface {

	// Name returns the unique identifier/type of the notification
	Name() string

	// Sender returns the originator of the notification
	Sender() interface{}

	// Info returns extra info associated with the notification
	Info() StringKeyMap
}

// NotificationObserver defines the interface for objects that listen to notifications
//
// Implements the observer pattern for reactive handling of named events
type NotificationObserver interface {

	// OnNotificationReceived is the callback invoked when a notification is dispatched
	//
	// Parameter notify: Notification instance containing event details (name, sender, info)
	OnNotificationReceived(notify Notification)
}

// BaseNotification provides a concrete implementation of the Notification interface
//
// Serves as the base struct for all concrete notification types
type BaseNotification struct {

	// name is the unique identifier/type of the notification
	name string

	// sender is the originator of the notification
	sender interface{}

	// info contains extra info for the notification
	info StringKeyMap
}

func NewNotification(name string, sender interface{}, info StringKeyMap) Notification {
	if info == nil {
		info = NewMap()
	}
	return &BaseNotification{
		name:   name,
		sender: sender,
		info:   info,
	}
}

func (notify *BaseNotification) Name() string {
	return notify.name
}

func (notify *BaseNotification) Sender() interface{} {
	return notify.sender
}

func (notify *BaseNotification) Info() StringKeyMap {
	return notify.info
}

// NotificationCenter implements the publish-subscribe pattern for notification dispatching
//
// Manages a registry of observers and dispatches notifications to relevant subscribers
type NotificationCenter struct {

	// observers maps notification names to lists of subscribed observers
	//
	// When a notification is dispatched, all observers for its name receive it
	observers map[string][]NotificationObserver
}

func NewNotificationCenter() *NotificationCenter {
	return &NotificationCenter{
		observers: make(map[string][]NotificationObserver, 128),
	}
}

func (center *NotificationCenter) getObservers(name string) []NotificationObserver {
	return center.observers[name]
}

// Add observer with notification name
func (center *NotificationCenter) Add(observer NotificationObserver, name string) {
	array := center.observers[name]
	if array == nil {
		array = make([]NotificationObserver, 0, 8)
	} else {
		for _, item := range array {
			if item == observer {
				// already exists
				return
			}
		}
	}
	center.observers[name] = append(array, observer)
}

// Remove observer from notification center with name
func (center *NotificationCenter) Remove(observer NotificationObserver, name string) {
	array := center.observers[name]
	if array != nil {
		array = remove(array, observer)
		if len(array) == 0 {
			delete(center.observers, name)
		} else {
			center.observers[name] = array
		}
	}
}

// Remove observer from notification center, no matter what names
func (center *NotificationCenter) RemoveAll(observer NotificationObserver) {
	count := len(center.observers)
	names := make([]string, 0, count)
	for key := range center.observers {
		names = append(names, key)
	}
	for _, name := range names {
		center.Remove(observer, name)
	}
}

func find(observer NotificationObserver, list []NotificationObserver) int {
	for index, item := range list {
		if item == observer {
			return index
		}
	}
	return -1
}

func remove(list []NotificationObserver, item NotificationObserver) []NotificationObserver {
	pos := find(item, list)
	if pos < 0 {
		return list
	} else if pos == 0 {
		return list[1:]
	}
	length := len(list) - 1
	if pos == length {
		return list[:length]
	}
	out := make([]NotificationObserver, length)
	index := 0
	for ; index < pos; index++ {
		out[index] = list[index]
	}
	for ; index < length; index++ {
		out[index] = list[index+1]
	}
	return out
}

// Default notification center
var defaultCenter = NewNotificationCenter()

// Add observer with notification name
func NotificationAddObserver(observer NotificationObserver, name string) {
	defaultCenter.Add(observer, name)
}

// Remove observer from default center
func NotificationRemoveObserver(observer NotificationObserver, name string) {
	if name == "" {
		defaultCenter.RemoveAll(observer)
	} else {
		defaultCenter.Remove(observer, name)
	}
}

// Post a notification (with name, sender and extra info)
func NotificationPost(name string, sender interface{}, info StringKeyMap) Notification {
	observers := defaultCenter.getObservers(name)
	if observers == nil {
		return nil
	}
	notify := NewNotification(name, sender, info)
	for _, item := range observers {
		item.OnNotificationReceived(notify)
	}
	return notify
}
