package mochatex

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/floholz/mochatex/internal/parsing"
	"log"
	"sort"
	"text/template"
)

var fyneApp fyne.App
var window fyne.Window
var errLog, infoLog *log.Logger
var mainContent fyne.Container

func Gui(err, info *log.Logger) {
	errLog = err
	infoLog = info
	fyneApp = app.NewWithID("mochatex")
	window = fyneApp.NewWindow("MochaTeX")
	window.SetIcon(theme.FileIcon())
	window.Resize(fyne.NewSize(1600, 1000))
	window.SetContent(contentLayout())
	infoLog.Println("Show and Run fyne window.")
	window.ShowAndRun()
}

func contentLayout() *fyne.Container {
	mainContent = *container.NewCenter(
		container.NewVBox(
			widget.NewLabel("No template loaded"),
			widget.NewButtonWithIcon("Open Template", theme.DocumentIcon(), openTemplate),
		),
	)
	return container.NewBorder(
		toolbar(),
		nil, nil, nil,
		&mainContent,
	)
}

func toolbar() *fyne.Container {
	return container.NewBorder(
		nil,
		container.NewVBox(
			widget.NewSeparator(),
			layout.NewSpacer(),
		),
		container.NewHBox(
			widget.NewButtonWithIcon("Open Template", theme.DocumentIcon(), openTemplate),
			widget.NewButtonWithIcon("Load Details", theme.FileTextIcon(), openDetails),
		),
		widget.NewToolbar(
			widget.NewToolbarAction(theme.HelpIcon(), displayHelp),
		),
	)
}

func displayTmplFields(tmpl *template.Template) {
	templateFields := parsing.MapTemplateFields(tmpl)
	infoLog.Printf("Template Fields: %v\n", templateFields)

	keys := make([]string, 0, len(templateFields))
	for k := range templateFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	form := widget.NewForm()
	for _, field := range keys {
		form.Append(field, widget.NewEntry())
	}
	mainContent = *container.NewVBox(
		widget.NewLabel("Template loaded:   "+tmpl.Name()),
		form,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Export details", theme.UploadIcon(), func() {
				infoLog.Println("Export details as JSON.")
			}),
			widget.NewButtonWithIcon("Apply details & save Pdf", theme.DocumentSaveIcon(), func() {
				infoLog.Println("Saved as Pdf.")
			}),
			layout.NewSpacer(),
		),
	)
	mainContent.Refresh()
}

func openTemplate() {
	infoLog.Println("Open template LaTeX file.")
	dlg := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
		if err != nil {
			errLog.Printf("error opening template file: %v", err)
			return
		}
		if closer == nil {
			return
		}
		texPath := closer.URI().Path()
		infoLog.Printf("Opened file: %v", texPath)

		tmpl := parsing.ParseTexFile(&texPath, errLog, infoLog)
		displayTmplFields(tmpl)
	}, window)
	dlg.SetFilter(storage.NewExtensionFileFilter([]string{".tex"}))
	dlg.Resize(fyne.NewSize(800, 600))
	dlg.Show()
}

func openDetails() {
	infoLog.Println("Open details JSON file.")
	dlg := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
		if err != nil {
			errLog.Printf("error opening template file: %v", err)
			return
		}
		if closer == nil {
			return
		}
		infoLog.Printf("Opened file: %v", closer.URI())
	}, window)
	dlg.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	dlg.Resize(fyne.NewSize(800, 600))
	dlg.Show()
}

func applyTemplate() {
	infoLog.Println("Apply details to template.")
	infoLog.Println("Saved as Pdf.")
}

func displayHelp() {
	infoLog.Println("Display help dialog.")
}
