package keys

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var dll = syscall.NewLazyDLL("user32.dll")
var procKeyBd = dll.NewProc("keybd_event")
var procSetCursorPos = dll.NewProc("SetCursorPos")
var procGetCursorPos = dll.NewProc("GetCursorPos")
var procMouseEvent = dll.NewProc("mouse_event")

// 定义 POINT 结构体
type POINT struct {
	X int32
	Y int32
}

type KeyBd struct {
	keys []int
}

// 定义鼠标事件常量
const (
	MOUSEEVENTF_LEFTDOWN   = 0x0002 // 鼠标左键按下
	MOUSEEVENTF_LEFTUP     = 0x0004 // 鼠标左键抬起
	MOUSEEVENTF_RIGHTDOWN  = 0x0008 // 鼠标右键按下
	MOUSEEVENTF_RIGHTUP    = 0x0010 // 鼠标右键抬起
	MOUSEEVENTF_MIDDLEDOWN = 0x0020 // 鼠标中键按下
	MOUSEEVENTF_MIDDLEUP   = 0x0040 // 鼠标中键抬起
)

var CursorPos POINT

func SetMouse(x int, y int, isDiff bool) {
	if isDiff {
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&CursorPos)))
		fmt.Println("cursorPos: ", CursorPos.X, CursorPos.Y)
		procSetCursorPos.Call(uintptr(CursorPos.X+int32(x)), uintptr(CursorPos.Y+int32(y)))
	} else {
		procSetCursorPos.Call(uintptr(int32(x)), uintptr(int32(y)))
	}
}

func ClickMouse(str string) {
	if str == "L" {
		procMouseEvent.Call(MOUSEEVENTF_LEFTDOWN, 0, 0, 0, 0)
		time.Sleep(50 * time.Millisecond) // 短暂延迟
		procMouseEvent.Call(MOUSEEVENTF_LEFTUP, 0, 0, 0, 0)
	} else if str == "R" {
		procMouseEvent.Call(MOUSEEVENTF_RIGHTDOWN, 0, 0, 0, 0)
		time.Sleep(50 * time.Millisecond) // 短暂延迟
		procMouseEvent.Call(MOUSEEVENTF_RIGHTUP, 0, 0, 0, 0)
	} else if str == "M" {
		procMouseEvent.Call(MOUSEEVENTF_MIDDLEDOWN, 0, 0, 0, 0)
		time.Sleep(50 * time.Millisecond) // 短暂延迟
		procMouseEvent.Call(MOUSEEVENTF_MIDDLEUP, 0, 0, 0, 0)
	}
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
