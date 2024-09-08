package keys

import (
	"dcontrol/server/setting"
	
	"fmt"
	"time"
	"unsafe"
	"syscall"
)

var (
	procSetCursorPos = user32.NewProc("SetCursorPos")
	procGetCursorPos = user32.NewProc("GetCursorPos")
	procMouseEvent   = user32.NewProc("mouse_event")
	// 监听鼠标滚动
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessage     = user32.NewProc("DispatchMessageW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	// 获取屏幕的高度
	procGetSystemMetrics    = user32.NewProc("GetSystemMetrics")
)

// 监听鼠标滚动
const (
	WH_MOUSE_LL = 14
	WM_MOUSEWHEEL = 0x020A
	SCREEN_BOTTOM_BUFFER = 5 // 距离屏幕底部5像素监听滚动
	SM_CYSCREEN = 1
)

// 定义鼠标事件常量
const (
	MOUSEEVENTF_LEFTDOWN   = 0x0002 // 鼠标左键按下
	MOUSEEVENTF_LEFTUP     = 0x0004 // 鼠标左键抬起
	MOUSEEVENTF_RIGHTDOWN  = 0x0008 // 鼠标右键按下
	MOUSEEVENTF_RIGHTUP    = 0x0010 // 鼠标右键抬起
	MOUSEEVENTF_MIDDLEDOWN = 0x0020 // 鼠标中键按下
	MOUSEEVENTF_MIDDLEUP   = 0x0040 // 鼠标中键抬起
)


type MSLLHOOKSTRUCT struct {
	Pt  POINT
	MouseData uint32
	Flags     uint32
	Time      uint32
	DwExtra   uintptr
}

type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// 定义 POINT 结构体
type POINT struct {
	X int32
	Y int32
}

var hHook syscall.Handle

var CursorPos POINT

func SetMouse(x int, y int, isDiff bool) {
	fmt.Println("cursorPos: ", CursorPos.X, CursorPos.Y, ", diff: ", x, y)
	if isDiff {
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&CursorPos)))
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

func ScrollMouse(scroll int, weight int) {
	var direct = -1 // 1: 向上滚动, -1: 向下滚动
	if scroll > 0 {
		direct = 1
	}
	if weight > 4 {
		weight = 4
	} else if weight < 1 {
		weight = 1
	}
	var delta = 120 * weight * direct
	procMouseEvent.Call(uintptr(0x0800), 0, 0, uintptr(delta), 0)
}


// Callback function for the mouse hook
func mouseHookCallback(nCode int32, wParam, lParam uintptr) uintptr {
	if nCode >= 0 && wParam == WM_MOUSEWHEEL && setting.Conf.Volume {
		mouseHookStruct := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		ret, _, _ := procGetSystemMetrics.Call(SM_CYSCREEN)
		screenBottom := int32(ret)
		mouseY := mouseHookStruct.Pt.Y
		if mouseY >= screenBottom - SCREEN_BOTTOM_BUFFER {
			var dy = int32(mouseHookStruct.MouseData>>16)
			fmt.Printf("Mouse wheel delta: %d, Mouse position: %d (within %d pixels of screen bottom)\n",
				dy, mouseY, SCREEN_BOTTOM_BUFFER)
			if dy > 1000 { // down
				RunKeys(KeyMap["VOLUME_DOWN"])
			} else {
				RunKeys(KeyMap["VOLUME_UP"])
			}
		}
	}
	ret, _, _ := procCallNextHookEx.Call(
			uintptr(hHook),
			uintptr(nCode),
			wParam,
			lParam,
	)
	return ret
}

func setHook() syscall.Handle {
	hook, _, _ := procSetWindowsHookEx.Call(
			WH_MOUSE_LL,
			syscall.NewCallback(mouseHookCallback),
			0,
			0,
	)
	return syscall.Handle(hook)
}

func removeHook(hook syscall.Handle) {
	procUnhookWindowsHookEx.Call(uintptr(hook))
}

func ListenScroll() {
	hHook = setHook()
	if hHook == 0 {
			fmt.Println("Failed to set mouse hook.")
			return
	}
	defer removeHook(hHook)
	fmt.Println("ListenScroll start...")
	var msg MSG
	for {
			procGetMessage.Call(
					uintptr(unsafe.Pointer(&msg)),
					0,
					0,
					0,
			)
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}