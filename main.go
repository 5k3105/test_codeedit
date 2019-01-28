package main

import (
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
)

var (
	ap             *Application
	codefont       *gui.QFont
	fontfile       = ":/fonts/saxmono.ttf"
	codefontfamily = "saxMono"
)

type Application struct {
	*widgets.QApplication
	Window *widgets.QMainWindow
}

func main() {
	ap = &Application{}
	ap.QApplication = widgets.NewQApplication(len(os.Args), os.Args)
	window := widgets.NewQMainWindow(nil, 0)
	ap.Window = window
	ap.Window.SetWindowTitle("Code Editor")
	_ = gui.QFontDatabase_AddApplicationFont(fontfile)
	codefont = gui.NewQFont2(codefontfamily, 12, 50, false)
	codeedit := New_CodeEditor(window)
	window.SetCentralWidget(codeedit)
	widgets.QApplication_SetStyle2("fusion")
	window.ShowMaximized()
	widgets.QApplication_Exec()
}
