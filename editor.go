package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"strconv"
)

type CodeEditor struct {
	*widgets.QPlainTextEdit
	line_area           *widgets.QWidget
	line_hilight        *gui.QBrush
	line_area_bg        *gui.QBrush
	line_area_fontcolor *gui.QColor
}

func New_CodeEditor(window *widgets.QMainWindow) *CodeEditor {
	editor := widgets.NewQPlainTextEdit(window)
	doc := gui.NewQTextDocument2("", window)
	doc.SetDefaultFont(codefont)
	doc.SetDocumentMargin(10.0)
	editor.SetDocument(doc)
	_ = New_GolangHighlighter(editor.Document())
	editor.SetLineWrapMode(widgets.QPlainTextEdit__NoWrap)
	editor.SetTabStopDistance(editor.TabStopDistance() / 2)
	editor.SetFont(codefont)
	editor.SetAutoFillBackground(true)

	line_area := widgets.NewQWidget(editor, core.Qt__Widget)
	codeedit := &CodeEditor{
		QPlainTextEdit: editor,
		line_area:      line_area,
	}
	line_area.ConnectPaintEvent(codeedit.line_paintevent)
	codeedit.set_colors()
	codeedit.ConnectPaintEvent(codeedit.paintEvent)
	codeedit.ConnectWheelEvent(codeedit.wheel_event)
	codeedit.ConnectCursorPositionChanged(codeedit.cursor_position_changed)
	codeedit.ConnectBlockCountChanged(codeedit.block_count_changed)
	codeedit.ConnectUpdateRequest(codeedit.update_line_area)
	codeedit.ConnectResizeEvent(codeedit.resize_event)
	codeedit.block_count_changed(0)
	return codeedit
}

func (editor *CodeEditor) wheel_event(e *gui.QWheelEvent) {
	if e.Modifiers() == core.Qt__ControlModifier {
		if e.AngleDelta().Y() > 0 {
			editor.ZoomIn(1)
		} else {
			editor.ZoomOut(1)
		}
		editor.update_viewport()
	} else {
		editor.WheelEventDefault(e)
	}
}

func (editor *CodeEditor) set_colors() {
	color := gui.NewQColor6("e5e5e5")
	editor.line_area_fontcolor = color
	color.SetAlpha(40)
	editor.line_area_bg = gui.NewQBrush3(color, core.Qt__SolidPattern)
	color.SetAlpha(20)
	editor.line_hilight = gui.NewQBrush3(color, core.Qt__SolidPattern)
}

func (editor *CodeEditor) block_count_changed(newBlockCount int) {
	editor.update_viewport()
}

func (editor *CodeEditor) update_viewport() {
	editor.SetViewportMargins(editor.line_area_width(), 0, 0, 0)
}

func (editor *CodeEditor) line_area_width() int {
	digits, max := 1, 1
	if editor.BlockCount() > max {
		max = editor.BlockCount()
	}
	for ; max >= 10; digits++ {
		max /= 10
	}
	space := 35 + editor.FontMetrics().HorizontalAdvance("9", 1)*digits
	return space
}

func (editor *CodeEditor) update_line_area(rect *core.QRect, dy int) {
	line_area := editor.line_area
	if dy > 0 {
		line_area.Scroll(0, dy)
	} else {
		line_area.Update2(0, rect.Y(), line_area.Width(), rect.Height())
	}
	if rect.Contains2(editor.Viewport().Rect(), true) {
		editor.block_count_changed(0)
	}
}

func (editor *CodeEditor) resize_event(event *gui.QResizeEvent) {
	editor.ResizeEventDefault(event)
	cr := editor.ContentsRect()
	newrec := core.NewQRect4(cr.Left(), cr.Top(), editor.line_area_width(), cr.Height())
	editor.line_area.SetGeometry(newrec)
}

func (editor *CodeEditor) cursor_position_changed() {
	editor.Viewport().Update()
}

func (editor *CodeEditor) paintEvent(event *gui.QPaintEvent) { /// line hilight
	painter := gui.NewQPainter2(editor.Viewport())
	painter.SetRenderHint(gui.QPainter__Antialiasing, true)

	rect := editor.CursorRect2()
	rect.SetX(0)
	rect.SetWidth(editor.Viewport().Width())
	painter.FillRect3(rect, editor.line_hilight)
	editor.PaintEventDefault(event)
	painter.DestroyQPainter()
}

func (editor *CodeEditor) line_paintevent(event *gui.QPaintEvent) { /// line numbers
	painter := gui.NewQPainter2(editor.line_area)
	painter.FillRect3(event.Rect(), editor.line_area_bg)

	block := editor.FirstVisibleBlock()
	blockNumber := block.BlockNumber()
	top := int(editor.BlockBoundingGeometry(block).Translated2(editor.ContentOffset()).Top())
	bottom := top + int(editor.BlockBoundingRect(block).Height())

	var pos *core.QRect
	painter.SetPen2(editor.line_area_fontcolor)
	for ; block.IsValid() && top <= event.Rect().Bottom(); blockNumber++ {
		if block.IsVisible() && bottom >= event.Rect().Top() {
			number := strconv.Itoa(blockNumber + 1)
			pos = core.NewQRect4(0, top, editor.line_area.Width(), editor.FontMetrics().Height())
			boundingrect := painter.BoundingRect2(pos, int(core.Qt__AlignHCenter), number)
			painter.DrawText6(pos, int(core.Qt__AlignHCenter), number, boundingrect)
		}
		block = block.Next()
		top = bottom
		bottom = top + int(editor.BlockBoundingRect(block).Height())
	}
	editor.line_area.PaintEventDefault(event)
	painter.DestroyQPainter()
}
