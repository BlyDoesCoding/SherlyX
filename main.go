package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {

	var pathArray []string

	myApp := app.New()
	myWindow := myApp.NewWindow("SherlyX")

	myWindow.Resize(fyne.NewSize(1920/2, 1080/2))

	fileEntry := widget.NewEntry()

	plus := widget.NewIcon(theme.FolderNewIcon())

	vbox := container.NewVBox()

	scrollableVBox := container.NewVScroll(vbox)

	startButton := widget.NewButton("Start", func() {
		start(pathArray, vbox)

	})

	fileEntry.SetPlaceHolder("Drag Folders to check here")

	page_2 := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		startButton,
		layout.NewSpacer(),
	)

	page_1 := container.New(
		layout.NewVBoxLayout(),
		fileEntry,
		layout.NewSpacer(),
		plus,
		layout.NewSpacer(),
	)

	myWindow.SetOnDropped(func(p fyne.Position, u []fyne.URI) {

		pathArray = append(pathArray, uriSliceToFilePaths(u)...)
		fileEntry.MultiLine = true
		fileEntry.SetText(stringArrayToString(pathArray))
	})

	tabs := container.NewAppTabs(
		container.NewTabItem("Folder", page_1),
		container.NewTabItem("Detecting", page_2),
		container.NewTabItem("Output", scrollableVBox),
	)
	tabs.SetTabLocation(container.TabLocationBottom)
	myWindow.SetOnDropped(func(p fyne.Position, u []fyne.URI) {

		pathArray = append(pathArray, uriSliceToFilePaths(u)...)
		fileEntry.MultiLine = true
		fileEntry.SetText(stringArrayToString(pathArray))
		tabs.SelectIndex(1)
	})
	myWindow.SetContent(tabs)

	myWindow.ShowAndRun()
}
