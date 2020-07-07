package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/mattn/go-gtk/glib"

	"github.com/mattn/go-gtk/gdk"

	"github.com/mattn/go-gtk/gtk"
)

//控件结构体
type ChessWidget struct {
	window *gtk.Window
}

//控件属性结构体
type ChessInfo struct {
	w, h int //窗口的宽度和高度
	x, y int //鼠标点击，相当于窗口的坐标
}

//黑白棋结构体
type Chessboard struct {
	ChessWidget //匿名字段
	ChessInfo
}

//方法：创建控件，设置从简属性
func (obj *Chessboard) CreatWindow() {
	//加载glade文件
	builder := gtk.NewBuilder()
	builder.AddFromFile("ui.glade")
	//窗口相关
	obj.window = gtk.WindowFromObject(builder.GetObject("window1"))
	obj.window.SetAppPaintable(true)           //允许绘图
	obj.window.SetPosition(gtk.WIN_POS_CENTER) //剧中显示
	obj.w, obj.h = 800, 480                    //窗口宽度和高度
	obj.window.SetSizeRequest(800, 480)
	obj.window.SetDecorated(false)
	//设置事件，让窗口可以捕捉鼠标点击和移动
	obj.window.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

}

//鼠标点击事件函数
func MousePressEvent(ctx *glib.CallbackContext) {

	//获取用户传递的参数
	data := ctx.Data()
	obj, ok := data.(*Chessboard) //类型断言
	if ok == false {
		fmt.Println("MousePressEvent  Chessboard ERR")
	}
	//获取鼠标按下结构体变量，系统内部的变量，不是用户传参变量
	arg := ctx.Args(0)
	event := *(**gdk.EventButton)(unsafe.Pointer(&arg))
	//保存点击的X，Y坐标
	obj.x, obj.y = int(event.X), int(event.Y)
	fmt.Println("x = ", obj.x, ",y = ", obj.y)

}

//鼠标移动事件函数
func MouseMoveEvent(ctx *glib.CallbackContext) {

	//获取用户传递的参数
	data := ctx.Data()
	obj, ok := data.(*Chessboard) //类型断言
	if ok == false {
		fmt.Println("MouseMoveEvent  Chessboard ERR")
	}
	//获取鼠标按下结构体变量，系统内部的变量，不是用户传参变量
	arg := ctx.Args(0)
	event := *(**gdk.EventButton)(unsafe.Pointer(&arg))
	//保存点击的X，Y坐标
	x, y := int(event.XRoot)-obj.x, int(event.YRoot)-obj.y
	obj.window.Move(x, y)

}

//方法：事件、信号处理
func (obj *Chessboard) HandleSignal() {
	//鼠标点击事件
	obj.window.Connect("button-press-event", MousePressEvent, obj)
	//鼠标移动事件
	obj.window.Connect("motion-notify-event", MouseMoveEvent, obj)
}
func main() {
	gtk.Init(&os.Args)
	//创建结构体变量
	var obj Chessboard
	obj.CreatWindow()  //创建控件，设置控件属性
	obj.HandleSignal() //事件信号处理
	obj.window.ShowAll()
	gtk.Main()
}
