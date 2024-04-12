package ui

import (
	"pokapoka-viewer/pkg/blog"
	"pokapoka-viewer/pkg/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	app      *tview.Application
	frontEnd *FrontEnd
}

type FrontEnd struct {
	mainFlex  *tview.Flex
	urlInput  *tview.InputField
	topPane   *tview.List
	MediaURLPane *tview.List
	commentPane *tview.List
	modal *tview.Modal
}

func NewApp() *App {
	app := &App{
		app: tview.NewApplication(),
	}
	app.frontEnd = app.setupUI()
	return app
}

func (app *App) Run() error {
	return app.app.SetRoot(app.frontEnd.mainFlex, true).Run()
}

func (app *App) setupUI() *FrontEnd {
	frontEnd := &FrontEnd{
		mainFlex: tview.NewFlex().SetDirection(tview.FlexRow),
		urlInput:  tview.NewInputField(),
		topPane:  tview.NewList(),
		MediaURLPane: tview.NewList(),
		commentPane: tview.NewList(),
		modal: tview.NewModal(),
	}

	frontEnd.mainFlex.SetBorder(true).SetTitle("ぽかぽかビューア").SetTitleAlign(tview.AlignCenter)
	frontEnd.mainFlex.AddItem(frontEnd.urlInput, 1, 1, true)
	frontEnd.mainFlex.AddItem(frontEnd.topPane, 0, 1, false)
	frontEnd.mainFlex.AddItem(frontEnd.MediaURLPane, 0, 1, true)
	frontEnd.mainFlex.AddItem(frontEnd.commentPane, 0, 3, true)

	frontEnd.topPane.SetBorder(true).SetTitle("概要").SetTitleAlign(tview.AlignCenter)
	frontEnd.MediaURLPane.SetBorder(true).SetTitle("画像URL一覧").SetTitleAlign(tview.AlignCenter)
	frontEnd.commentPane.SetBorder(true).SetTitle("コメント一覧").SetTitleAlign(tview.AlignCenter)

	// reset focus to mainFlex
	frontEnd.modal.AddButtons([]string{"OK"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.app.SetRoot(frontEnd.mainFlex, true)
	})

	app.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.app.Stop()
			return nil
		}
		return event
	})
	

	frontEnd.urlInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.app.SetFocus(frontEnd.MediaURLPane)
			return nil
		} else if event.Key() == tcell.KeyBacktab {
			app.app.SetFocus(frontEnd.commentPane)
			return nil
		}
		return event
	})
	frontEnd.MediaURLPane.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.app.SetFocus(frontEnd.commentPane)
			return nil
		} else if event.Key() == tcell.KeyBacktab {
			app.app.SetFocus(frontEnd.urlInput)
			return nil
		}
		return event
	})
	frontEnd.commentPane.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.app.SetFocus(frontEnd.urlInput)
			return nil
		} else if event.Key() == tcell.KeyBacktab {
			app.app.SetFocus(frontEnd.MediaURLPane)
			return nil
		}
		return event
	})

	frontEnd.urlInput.SetDoneFunc(func(key tcell.Key) {
		blogURL := frontEnd.urlInput.GetText()
		b, err := blog.GetBlog(blogURL)
		if err != nil {
			frontEnd.modal.SetText(err.Error())
			app.app.SetRoot(frontEnd.modal, true)
			app.app.SetFocus(frontEnd.modal)
			return
		}

		frontEnd.topPane.Clear()
		frontEnd.topPane.AddItem("タイトル: " + b.Title, "", ' ', nil)
		frontEnd.topPane.AddItem("詳細: " + b.Detail, "", ' ', nil)
		frontEnd.MediaURLPane.Clear()
		for _, url := range b.MediaURLs {
			frontEnd.MediaURLPane.AddItem(url, "", ' ', nil)
		}
		frontEnd.commentPane.Clear()
		for _, comment := range b.Comments {
			frontEnd.commentPane.AddItem(comment, "", ' ', nil)
		}
	})

	// mediaURLPane enter key event
	frontEnd.MediaURLPane.SetSelectedFunc(func(i int, mainText string, secondaryText string, shortcut rune) {
		url, _ := frontEnd.MediaURLPane.GetItemText(i)
		utils.OpenBrowser(url)
	})

	return frontEnd
}
