package mochatex

import (
	"encoding/json"
	"fmt"
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
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

var fyneApp fyne.App
var window fyne.Window
var errLog, infoLog *log.Logger

var detailsBtn *widget.Button
var mainContent fyne.Container
var fields map[string]*widget.Entry

var texPath string
var tmpl *template.Template
var jsonPath string
var dtls map[string]interface{}

func Gui(err, info *log.Logger) {
	errLog = err
	infoLog = info
	fyneApp = app.NewWithID("mochatex")
	window = fyneApp.NewWindow("MochaTeX")
	window.SetIcon(theme.FileIcon())
	window.Resize(fyne.NewSize(900, 1000))
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
		debugInfoContainer(),
		nil, nil,
		&mainContent,
	)
}

func toolbar() *fyne.Container {
	detailsBtn = widget.NewButtonWithIcon("Load Details", theme.FileTextIcon(), openDetails)
	if tmpl == nil {
		detailsBtn.Disable()
	}
	return container.NewBorder(
		nil,
		container.NewVBox(
			widget.NewSeparator(),
			layout.NewSpacer(),
		),
		container.NewHBox(
			widget.NewButtonWithIcon("Open Template", theme.DocumentIcon(), openTemplate),
			detailsBtn,
		),
		widget.NewToolbar(
			widget.NewToolbarAction(theme.HelpIcon(), displayHelp),
		),
	)
}

func debugInfoContainer() *fyne.Container {
	if _, ok := os.LookupEnv("MOCHATEX_DEBUG"); !ok {
		return container.NewVBox()
	}
	return container.NewVBox(widget.NewLabel("Debug mode enabled"))
}

func displayTmplFields() {
	templateFields := parsing.MapTemplateFields(tmpl)
	infoLog.Printf("Template Fields: %v\n", templateFields)

	keys := make([]string, 0, len(templateFields))
	for k := range templateFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	form := widget.NewForm()
	fields = make(map[string]*widget.Entry)
	for _, field := range keys {
		entry := widget.NewEntry()
		entry.OnChanged = func(s string) {
			onDtlEntryChanged(field, s)
		}
		fields[field] = entry
		form.Append(field, entry)
	}

	mainContent = *container.NewVBox(
		widget.NewLabel("Template loaded:   "+tmpl.Name()),
		form,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Export details", theme.UploadIcon(), exportDetails),
			widget.NewButtonWithIcon("Apply details & save Pdf", theme.DocumentSaveIcon(), applyTemplate),
			layout.NewSpacer(),
		),
	)
	mainContent.Refresh()
}

func fillDetails() {
	fltn := parsing.FlattenJson(dtls)
	for key, value := range fltn {
		_, ok := fields[key]
		if ok {
			onChangeFn := fields[key].OnChanged
			fields[key].OnChanged = nil
			fields[key].SetText(value)
			fields[key].OnChanged = onChangeFn
		}
	}
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
		texPath = closer.URI().Path()
		infoLog.Printf("Opened file: %v", texPath)

		tmpl = parsing.ParseTexFile(&texPath, errLog, infoLog)
		if tmpl != nil {
			detailsBtn.Enable()
			displayTmplFields()
		} else {
			texPath = ""
			detailsBtn.Disable()
		}
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
		jsonPath = closer.URI().Path()
		infoLog.Printf("Opened file: %v", jsonPath)

		dtls = parsing.ParseJsonFile(&jsonPath, errLog, infoLog)
		if dtls != nil {
			fillDetails()
			infoLog.Println("Loaded Details.")
		} else {
			jsonPath = ""
		}
	}, window)
	dlg.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	dlg.Resize(fyne.NewSize(800, 600))
	dlg.Show()
}

func applyTemplate() {
	infoLog.Println("Apply details to template.")
	dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			errLog.Printf("error saving pdf: %v", err)
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		infoLog.Println("Saved as Pdf.")
		p := filepath.Dir(texPath)
		if writer.URI() == nil {
			return
		} else {
			p = writer.URI().Path()
		}
		pdfPath, err := StartJob(tmpl, dtls, p)
		if err != nil {
			errLog.Fatalf("error while compiling pdf: %v", err)
		}
		infoLog.Printf("Successfully created PDF at location: %s", filepath.Join(p, pdfPath))

	}, window)
	u, err := storage.ListerForURI(storage.NewFileURI(filepath.Dir(texPath)))
	if err != nil {
		errLog.Printf("error setting up save file path: %v\n", err)
	} else {
		dlg.SetLocation(u)
	}
	dlg.Resize(fyne.NewSize(800, 600))
	dlg.Show()
}

func exportDetails() {
	infoLog.Println("Export details as JSON.")
	dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			errLog.Printf("error saving json: %v", err)
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		p := jsonPath
		if writer.URI() == nil {
			return
		} else {
			p = writer.URI().Path()
		}
		err = saveToJSONFile(p, dtls)
		if err != nil {
			errLog.Printf("error saving json: %v", err)
			return
		}
		infoLog.Printf("Successfully created PDF at location: %s", jsonPath)

	}, window)
	u, err := storage.ListerForURI(storage.NewFileURI(filepath.Dir(jsonPath)))
	if err != nil {
		errLog.Printf("error setting up save file path: %v\n", err)
	} else {
		dlg.SetLocation(u)
	}
	dlg.Resize(fyne.NewSize(800, 600))
	dlg.Show()
}

func displayHelp() {
	infoLog.Println("Display help dialog.")
}

func onDtlEntryChanged(field, value string) {
	if dtls == nil {
		dtls = make(map[string]interface{})
	}
	splits := strings.Split(strings.TrimPrefix(field, "."), ".")
	if len(splits) == 0 {
		return
	}

	// Start from the top-level map
	subDtlRef := dtls
	// Iterate through the keys except the last one
	for i := 0; i < len(splits)-1; i++ {
		key := splits[i]
		// Check if the subDtlRef key exists and is a map
		if next, ok := subDtlRef[key].(map[string]interface{}); ok {
			// Move down to the next level
			subDtlRef = next
		} else {
			// If the key doesn't exist or isn't a map, create a new map
			newMap := make(map[string]interface{})
			subDtlRef[key] = newMap
			subDtlRef = newMap
		}
	}
	// Set the value at the final key
	finalKey := splits[len(splits)-1]
	subDtlRef[finalKey] = value
}

func saveToJSONFile(filename string, data interface{}) error {
	// Marshal the interface{} to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ") // Use MarshalIndent for pretty printing
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Create or open the JSON file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
