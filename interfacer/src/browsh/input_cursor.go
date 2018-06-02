package browsh

import "unicode/utf8"

func (i *inputBox) renderCursor() {
	var xCursor int
	xFrameOffset := CurrentTab.frame.xScroll
	yFrameOffset := CurrentTab.frame.yScroll - uiHeight
	if urlInputBox.isActive {
		xFrameOffset = 0
		yFrameOffset = 0
	}
	if i.isMultiLine() {
		xCursor = i.xCursor
	} else {
		xCursor = i.textCursor
	}
	x := (i.X + xCursor) - i.xScroll - xFrameOffset
	y := (i.Y + i.yCursor) - i.yScroll - yFrameOffset
	mainRune, combiningRunes, style, _ := screen.GetContent(x, y)
	style = style.Reverse(true)
	screen.SetContent(x, y, mainRune, combiningRunes, style)
}

func (i *inputBox) cursorLeft() {
	i.xCursor--
	i.textCursor--
	i.updateAllCursors()
}

func (i *inputBox) cursorRight() {
	i.xCursor++
	i.textCursor++
	i.updateAllCursors()
}

func (i *inputBox) cursorUp() {
	i.multiLiner.moveYCursorBy(-1)
	i.updateAllCursors()
}

func (i *inputBox) cursorDown() {
	i.multiLiner.moveYCursorBy(1)
	i.updateAllCursors()
}

func (i *inputBox) cursorBackspace() {
	if (utf8.RuneCountInString(i.text) == 0) { return }
	if (i.textCursor == 0) { return }
	start := i.text[:i.textCursor - 1]
	end := i.text[i.textCursor:]
	i.text = start + end
	i.cursorLeft()
	i.sendInputBoxToBrowser()
}

func (i *inputBox) cursorInsertRune(theRune rune) {
	character := string(theRune)
	start := i.text[:i.textCursor]
	end := i.text[i.textCursor:]
	i.text = start + character + end
	i.cursorRight()
	i.sendInputBoxToBrowser()
}

func (i *inputBox) isCursorOverRightEdge() bool {
	return i.textCursor - i.xScroll >= i.Width
}

func (i *inputBox) isCursorOverLeftEdge() bool {
	return i.textCursor - i.xScroll <= -1
}

func (i *inputBox) isCursorOverTopEdge() bool {
	return i.yCursor - i.yScroll <= -1
}

func (i *inputBox) isCursorOverBottomEdge() bool {
	return i.yCursor - i.yScroll > i.Height
}

func (i *inputBox) getCharacterAt() string {
	var index int
	var c rune
	for index, c = range i.text {
		if index == i.textCursor {
			return string(c)
		}
	}
	return ""
}

func (i *inputBox) putCursorAtEnd() {
	i.textCursor = utf8.RuneCountInString(urlInputBox.text)
	// TODO: Do for multiline
}

func (i *inputBox) updateAllCursors() {
	i.updateXYCursors()
	if (i.isCursorOverLeftEdge() || !i.isBestFit()) { i.xScrollBy(-1) }
	if (i.isCursorOverTopEdge()) { i.yScrollBy(-1) }
	if (i.isCursorOverRightEdge()) { i.xScrollBy(1) }
	if (i.isCursorOverBottomEdge()) { i.yScrollBy(1) }
	i.limitTextCursor()
	i.updateXYCursors()
}

func (i *inputBox) limitTextCursor() {
	if (i.textCursor < 0) {
		i.textCursor = 0
	}
	if (i.textCursor > utf8.RuneCountInString(i.text)) {
		i.textCursor = utf8.RuneCountInString(i.text)
	}
}

func (i *inputBox) updateXYCursors() {
	if !i.isMultiLine() { return }
	i.multiLiner.updateCursor()
	i.renderCursor()
}
