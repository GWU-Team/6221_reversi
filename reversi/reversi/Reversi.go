package main

import (
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"github.com/mattn/go-gtk/gdkpixbuf"

	"github.com/mattn/go-gtk/glib"

	"github.com/mattn/go-gtk/gdk"

	"github.com/mattn/go-gtk/gtk"
)

type ChessWidget struct {
	window      *gtk.Window
	buttonMin   *gtk.Button 
	buttonClose *gtk.Button 
	labelBlack  *gtk.Label  
	labelWhite  *gtk.Label  
	labelTime   *gtk.Label  
	imageBlack  *gtk.Image  
	imageWhite  *gtk.Image  
}

type ChessInfo struct {
	w, h           int 
	x, y           int
	startX, startY int 
	gridW, gridH   int 

}


const (
	Empty = iota 
	Black        
	White       
)


type Chessboard struct {
	ChessWidget 
	ChessInfo
	currentRole    int 
	tipTimeId      int 
	machineTimerId int 
	leftTimerId    int
	rightTimerId   int
	timeNum        int 

	chess [8][8]int 
}


func ButtonSetImageFromFile(button *gtk.Button, filename string) {
	w, h := 0, 0
	button.GetSizeRequest(&w, &h)
	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale(filename, w-10, h-10, false)
	image := gtk.NewImageFromPixbuf(pixbuf)
	pixbuf.Unref()
	button.SetImage(image)
	button.SetCanFocus(false)
}


func ImageSetPicFromFile(image *gtk.Image, filename string) {
	w, h := 0, 0
	image.GetSizeRequest(&w, &h)
	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale(filename, w-10, h-10, false)
	image.SetFromPixbuf(pixbuf)
	pixbuf.Unref()
}

func (obj *Chessboard) CreatWindow() {
	builder := gtk.NewBuilder()
	builder.AddFromFile("ui.glade")
	obj.window = gtk.WindowFromObject(builder.GetObject("window1"))
	obj.window.SetAppPaintable(true)       
	obj.window.SetPosition(gtk.WIN_POS_CENTER) 
	obj.w, obj.h = 800, 480                    
	obj.window.SetSizeRequest(800, 480)
	obj.window.SetDecorated(false)

	obj.window.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

	obj.buttonMin = gtk.ButtonFromObject(builder.GetObject("buttonMin"))
	obj.buttonClose = gtk.ButtonFromObject(builder.GetObject("buttonClose"))

	ButtonSetImageFromFile(obj.buttonMin, "../image/min.png")
	ButtonSetImageFromFile(obj.buttonClose, "../image/close.png")

	obj.labelBlack = gtk.LabelFromObject(builder.GetObject("labelBlack"))
	obj.labelWhite = gtk.LabelFromObject(builder.GetObject("labelWhite"))
	obj.labelTime = gtk.LabelFromObject(builder.GetObject("labelTime"))

	obj.labelBlack.ModifyFontSize(50)
	obj.labelWhite.ModifyFontSize(50)
	obj.labelTime.ModifyFontSize(30)

	obj.labelBlack.SetText("2")
	obj.labelWhite.SetText("2")
	obj.labelTime.SetText("20")

	obj.labelBlack.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))
	obj.labelWhite.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))
	obj.labelTime.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	obj.imageBlack = gtk.ImageFromObject(builder.GetObject("imageBlack"))
	obj.imageWhite = gtk.ImageFromObject(builder.GetObject("imageWhite"))

	ImageSetPicFromFile(obj.imageBlack, "../image/black.png")
	ImageSetPicFromFile(obj.imageWhite, "../image/white.png")


	obj.startX, obj.startY = 200, 60
	obj.gridW, obj.gridH = 50, 40
}


func (obj *Chessboard) JudgeResult() {
	isOver := true 
	blackNum, whiteNum := 0, 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if obj.chess[i][j] == Black {
				blackNum++
			} else if obj.chess[i][j] == White {
				whiteNum++
			}

			if obj.JudgeRule(i, j, Black, false) > 0 || obj.JudgeRule(i, j, White, false) > 0 {
				isOver = false
			}
		}
	}

	obj.labelBlack.SetText(strconv.Itoa(blackNum))
	obj.labelWhite.SetText(strconv.Itoa(whiteNum))
	if isOver == false {
		return
	}


	glib.TimeoutRemove(obj.tipTimeId)
	glib.TimeoutRemove(obj.leftTimerId)

	var result string
	if blackNum > whiteNum {
		result = "You win \n Continue?"

	} else if whiteNum > blackNum {
		result = "You lose \n Continue?"
	} else {
		result = "dogfall\n  Continue?"
	}


	dialog := gtk.NewMessageDialog(
		obj.window,           
		gtk.DIALOG_MODAL,     
		gtk.MESSAGE_QUESTION, 
		gtk.BUTTONS_YES_NO,  
		result)               

	ret := dialog.Run()
	if ret == gtk.RESPONSE_YES {
		obj.InitChess() 
	}
	dialog.Destroy()
}

func (obj *Chessboard) MachinePlay() {

	glib.TimeoutRemove(obj.machineTimerId)
	max, px, py := 0, -1, -1

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			num := obj.JudgeRule(i, j, obj.currentRole, false)
			if num > 0 {

				if (i == 0 && j == 0) || (i == 7 && j == 0) || (i == 0 && j == 7) || (i == 7 && j == 7) {
					px, py = i, j
					goto End
				}
				if num > max {
					max, px, py = num, i, j
				}
			}
		}
	}
End:
	if px == -1 {
		obj.ChangeRole()
		return
	}

	obj.JudgeRule(px, py, obj.currentRole, true)
	obj.window.QueueDraw()
	obj.ChangeRole()
}


func (obj *Chessboard) ChangeRole() {

	obj.timeNum = 20
	obj.labelTime.SetText(strconv.Itoa(obj.timeNum))

	obj.imageBlack.Hide()
	obj.imageWhite.Hide()

	if obj.currentRole == Black {
		obj.currentRole = White
	} else {
		obj.currentRole = Black
	}

	if obj.currentRole == White {
		obj.machineTimerId = glib.TimeoutAdd(1000, func() bool {
			obj.MachinePlay() 
			return true
		})
	}
}
func (obj *Chessboard) JudgeRule(x, y int, role int, eatChess bool) (eatNum int) {

	dir := [8][2]int{{1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}}
	tempX, tempY := x, y                 
	if obj.chess[tempX][tempY] != Empty {
		return 0
	}

	for i := 0; i < 8; i++ {
		tempX += dir[i][0]
		tempY += dir[i][1] 

		if (tempX < 8 && tempX >= 0 && tempY < 8 && tempY >= 0) && (obj.chess[tempX][tempY] != role) && (obj.chess[tempX][tempY] != Empty) {
			tempX += dir[i][0]
			tempY += dir[i][1] 
			for tempX < 8 && tempX >= 0 && tempY < 8 && tempY >= 0 {
				if obj.chess[tempX][tempY] == Empty { 
					break
				}
				if obj.chess[tempX][tempY] == role {
					if eatChess == true { 
						obj.chess[x][y] = role 
						tempX -= dir[i][0]
						tempY -= dir[i][1] 
						for (tempX != x) || (tempY != y) {
	
							obj.chess[tempX][tempY] = role 
							tempX -= dir[i][0]
							tempY -= dir[i][1] 
							eatNum++          
						}

					} else {
						tempX -= dir[i][0]
						tempY -= dir[i][1]               
						for (tempX != x) || (tempY != y) { 
							tempX -= dir[i][0]
							tempY -= dir[i][1] 
							eatNum++
						}

					}
					break 
				}

				tempX += dir[i][0]
				tempY += dir[i][1] 

			}

		}
		tempX, tempY = x, y
	}
	a := 10
	a = a + 1

	return
}


func MousePressEvent(ctx *glib.CallbackContext) {


	data := ctx.Data()
	obj, ok := data.(*Chessboard)
	if ok == false {
		fmt.Println("MousePressEvent  Chessboard ERR")
	}

	arg := ctx.Args(0)
	event := *(**gdk.EventButton)(unsafe.Pointer(&arg))

	obj.x, obj.y = int(event.X), int(event.Y)
	//fmt.Println("x = ", obj.x, ",y = ", obj.y)
	i := (obj.x - obj.startX) / obj.gridW
	j := (obj.y - obj.startY) / obj.gridH
	if obj.currentRole == White { 
		return
	}

	if i >= 0 && i <= 7 && j >= 0 && j <= 7 {
		fmt.Printf("(%d,%d)\n", i, j)
		if obj.JudgeRule(i, j, obj.currentRole, true) > 0 {
			obj.window.QueueDraw()
			obj.ChangeRole()
		}

	}

}


func MouseMoveEvent(ctx *glib.CallbackContext) {
	data := ctx.Data()
	obj, ok := data.(*Chessboard) 
	if ok == false {
		fmt.Println("MouseMoveEvent  Chessboard ERR")
	}
	arg := ctx.Args(0)
	event := *(**gdk.EventButton)(unsafe.Pointer(&arg))
	x, y := int(event.XRoot)-obj.x, int(event.YRoot)-obj.y
	obj.window.Move(x, y)

}


func PaintEvent(ctx *glib.CallbackContext) {

	data := ctx.Data()
	obj, ok := data.(*Chessboard) 
	if ok == false {
		fmt.Println("PaintEvent  Chessboard ERR")
	}
	painter := obj.window.GetWindow().GetDrawable()
	gc := gdk.NewGC(painter)

	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale("../image/bg.jpg", obj.w, obj.h, false)
	blackPixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale("../image/black.png", obj.gridW, obj.gridH, false)
	whitePixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale("../image/white.png", obj.gridW, obj.gridH, false)
	painter.DrawPixbuf(gc, pixbuf, 0, 0, 0, 0, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if obj.chess[i][j] == Black {
				painter.DrawPixbuf(gc, blackPixbuf, 0, 0, obj.startX+i*obj.gridW, obj.startY+j*obj.gridH,
					-1, -1, gdk.RGB_DITHER_NONE, 0, 0)
			} else if obj.chess[i][j] == White {
				painter.DrawPixbuf(gc, whitePixbuf, 0, 0, obj.startX+i*obj.gridW, obj.startY+j*obj.gridH,
					-1, -1, gdk.RGB_DITHER_NONE, 0, 0)
			}
		}
	}
	pixbuf.Unref()
	blackPixbuf.Unref()
	whitePixbuf.Unref()
}

func (obj *Chessboard) HandleSignal() {
	//Mouse press event
	obj.window.Connect("button-press-event", MousePressEvent, obj)
	//Mouse move event
	obj.window.Connect("motion-notify-event", MouseMoveEvent, obj)

	obj.buttonMin.Clicked(func() {
		obj.window.Iconify()
	})

	obj.buttonClose.Clicked(func() {
		glib.TimeoutRemove(obj.tipTimeId)
		glib.TimeoutRemove(obj.leftTimerId)
		gtk.MainQuit()
	})
	obj.window.Connect("configure-event", func() {
		obj.window.QueueDraw()

	})
	obj.window.Connect("expose-event", PaintEvent, obj)
}

func ShowTip(obj *Chessboard) {
	if obj.currentRole == Black { 
		obj.imageWhite.Hide()
		if obj.imageBlack.GetVisible() == true {
			obj.imageBlack.Hide()
		} else {
			obj.imageBlack.Show()
		}
	} else {
		obj.imageBlack.Hide()
		if obj.imageWhite.GetVisible() == true {
			obj.imageWhite.Hide()
		} else {
			obj.imageWhite.Show()
		}
	}
}

func (obj *Chessboard) InitChess() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			obj.chess[i][j] = Empty
		}
	}
	obj.chess[3][3] = Black
	obj.chess[4][4] = Black
	obj.chess[4][3] = White
	obj.chess[3][4] = White
	obj.window.QueueDraw()
	obj.labelBlack.SetText("2")
	obj.labelWhite.SetText("2")
	obj.imageBlack.Hide()
	obj.imageWhite.Hide()
	obj.currentRole = Black

	obj.tipTimeId = glib.TimeoutAdd(500, func() bool {
		ShowTip(obj)
		return true
	})
	obj.timeNum = 20
	obj.labelTime.SetText(strconv.Itoa(obj.timeNum))

	obj.leftTimerId = glib.TimeoutAdd(1000, func() bool {

		obj.timeNum--
		obj.labelTime.SetText(strconv.Itoa(obj.timeNum))
		if obj.timeNum == 0 {
			obj.ChangeRole()
		}
		return true
	})
}
func main() {
	gtk.Init(&os.Args)
	var obj Chessboard
	obj.CreatWindow()  
	obj.HandleSignal() 
	obj.InitChess()    
	obj.window.Show()
	gtk.Main()
}
