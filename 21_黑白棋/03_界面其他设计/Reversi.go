package main

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/mattn/go-gtk/gdkpixbuf"

	"github.com/mattn/go-gtk/glib"

	"github.com/mattn/go-gtk/gdk"

	"github.com/mattn/go-gtk/gtk"
)

//控件结构体
type ChessWidget struct {
	window      *gtk.Window
	buttonMin   *gtk.Button //最小化按钮
	buttonClose *gtk.Button //关闭按钮
	labelBlack  *gtk.Label  //记录黑棋个数
	labelWhite  *gtk.Label  //记录白棋个数
	labelTime   *gtk.Label  //记录倒计时
	imageBlack  *gtk.Image  //提示该黑子落子
	imageWhite  *gtk.Image  //提示该白子落子
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

//函数，给按钮设置图片
func ButtonSetImageFromFile(button *gtk.Button, filename string) {
	//获取按钮大小
	w, h := 0, 0
	button.GetSizeRequest(&w, &h)
	//创建pixbuff
	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale(filename, w-10, h-10, false)
	//创建image
	image := gtk.NewImageFromPixbuf(pixbuf)
	//释放pixbuf
	pixbuf.Unref()
	//给按钮设置图片
	button.SetImage(image)
	//去掉按钮的焦距
	button.SetCanFocus(false)
}

//给image 设置图片
func ImageSetPicFromFile(image *gtk.Image, filename string) {
	//获取image大小
	w, h := 0, 0
	image.GetSizeRequest(&w, &h)
	//创建pixbuff
	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale(filename, w-10, h-10, false)
	//给image设置图片
	image.SetFromPixbuf(pixbuf)
	//释放pixbuf
	pixbuf.Unref()
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
	//按钮相关
	obj.buttonMin = gtk.ButtonFromObject(builder.GetObject("buttonMin"))
	obj.buttonClose = gtk.ButtonFromObject(builder.GetObject("buttonClose"))
	//给按钮设置图片
	ButtonSetImageFromFile(obj.buttonMin, "../image/min.png")
	ButtonSetImageFromFile(obj.buttonClose, "../image/close.png")
	//标签相关
	obj.labelBlack = gtk.LabelFromObject(builder.GetObject("labelBlack"))
	obj.labelWhite = gtk.LabelFromObject(builder.GetObject("labelWhite"))
	obj.labelTime = gtk.LabelFromObject(builder.GetObject("labelTime"))
	//设置字体大小
	obj.labelBlack.ModifyFontSize(50)
	obj.labelWhite.ModifyFontSize(50)
	obj.labelTime.ModifyFontSize(30)
	//设置内容
	obj.labelBlack.SetText("2")
	obj.labelWhite.SetText("2")
	obj.labelTime.SetText("20")
	//改变字体颜色
	obj.labelBlack.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))
	obj.labelWhite.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))
	obj.labelTime.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))
	//image相关
	obj.imageBlack = gtk.ImageFromObject(builder.GetObject("imageBlack"))
	obj.imageWhite = gtk.ImageFromObject(builder.GetObject("imageWhite"))
	//设置图片
	ImageSetPicFromFile(obj.imageBlack, "../image/black.png")
	ImageSetPicFromFile(obj.imageWhite, "../image/white.png")
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

//鼠标移动事件函数
func PaintEvent(ctx *glib.CallbackContext) {

	//获取用户传递的参数
	data := ctx.Data()
	obj, ok := data.(*Chessboard) //类型断言
	if ok == false {
		fmt.Println("PaintEvent  Chessboard ERR")
	}
	//获取画家，设置绘图区域
	painter := obj.window.GetWindow().GetDrawable()
	gc := gdk.NewGC(painter)
	//
	//创建pixbuff
	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale("../image/bg.jpg", obj.w, obj.h, false)

	//画图
	painter.DrawPixbuf(gc, pixbuf, 0, 0, 0, 0, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)

	//释放pixbuf
	pixbuf.Unref()
}

//方法：事件、信号处理
func (obj *Chessboard) HandleSignal() {
	//鼠标点击事件
	obj.window.Connect("button-press-event", MousePressEvent, obj)
	//鼠标移动事件
	obj.window.Connect("motion-notify-event", MouseMoveEvent, obj)

	obj.buttonMin.Clicked(func() {
		obj.window.Iconify()
	})
	//按钮的信号处理
	obj.buttonClose.Clicked(func() {
		gtk.MainQuit()
	})
	//绘图相关
	obj.window.Connect("configure-event", func() {
		//重新刷图
		obj.window.QueueDraw()

	})
	//绘图事件
	obj.window.Connect("expose-event", PaintEvent, obj)
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
