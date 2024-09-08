package keys

import (
	"fmt"
	"syscall"
	"time"
)

var user32 = syscall.NewLazyDLL("user32.dll")
var procKeyBd = user32.NewProc("keybd_event")

type KeyBd struct {
	keys []int
}

func Run(str string) {
	RunKeys(GetKeys(str)...)
}

func RunKeys(keys ...int) {
	fmt.Println("RunKeys", keys)
	keyBd := NewKeyBd(keys...)
	keyBd.Execute()
}

func NewKeyBd(keys ...int) KeyBd {
	keyBd := KeyBd{}
	keyBd.SetKeys(keys...)
	return keyBd
}

func (k *KeyBd) Clear() *KeyBd {
	k.keys = []int{}
	return k
}

func (k *KeyBd) SetKeys(keys ...int) {
	k.keys = keys
}

// AddKey add one key pressed
func (k *KeyBd) AddKey(key int) {
	k.keys = append(k.keys, key)
}

// Execute key bounding
func (k *KeyBd) Execute() error {
	err := k.Press()
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	err = k.Release()
	return err
}

// Press key(s)
func (k *KeyBd) Press() error {
	for _, key := range k.keys {
		downKey(key)
	}
	return nil
}

// Release key(s)
func (k *KeyBd) Release() error {
	for _, key := range k.keys {
		upKey(key)
	}
	return nil
}

func downKey(key int) {
	flag := 0
	if key < 0xFFF { // Detect if the key code is virtual or no
		flag |= _KEYEVENTF_SCANCODE
	} else {
		key -= 0xFFF
	}
	vkey := key + 0x80
	procKeyBd.Call(uintptr(key), uintptr(vkey), uintptr(flag), 0)
}

func upKey(key int) {
	flag := _KEYEVENTF_KEYUP
	if key < 0xFFF {
		flag |= _KEYEVENTF_SCANCODE
	} else {
		key -= 0xFFF
	}
	vkey := key + 0x80
	procKeyBd.Call(uintptr(key), uintptr(vkey), uintptr(flag), 0)
}
