package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/jiura/tview"
)

var req_form *tview.Form
var instructions *tview.TextView

var headers_form *tview.Form

var resp_status *tview.TextView
var resp_headers *tview.TextView
var resp_body *tview.TextView

var left_flex *tview.Flex
var right_flex *tview.Flex
var main_flex *tview.Flex

func sendRequest() (string, string, string) {
	_, method := req_form.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
	url := req_form.GetFormItem(1).(*tview.InputField).GetText()
	body := req_form.GetFormItem(2).(*tview.TextArea).GetText()

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return "", err.Error(), ""
	}

	for i := 0; i < 10; i += 2 {
		key := headers_form.GetFormItem(i).(*tview.InputField).GetText()

		if key == "" {
			continue
		}

		val := headers_form.GetFormItem(i + 1).(*tview.InputField).GetText()

		req.Header.Set(key, val)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err.Error(), ""
	}
	defer resp.Body.Close()

	var builder strings.Builder
	for key, values := range resp.Header {
		builder.WriteString(key)
		builder.WriteString(": ")
		builder.WriteString(strings.Join(values, ", "))
		builder.WriteRune('\n')
	}
	resp_headers_string := builder.String()

	resp_body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.Status, resp_headers_string, err.Error()
	}

	return resp.Status, resp_headers_string, string(resp_body_bytes)
}

func main() {
	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlQ:
			fallthrough
		case tcell.KeyEsc:
			app.Stop()
		case tcell.KeyCtrlC:
			return tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
		case tcell.KeyCtrlS:
			resp_status.SetText("Loading...")
			resp_headers.SetText("Loading...")
			resp_body.SetText("Loading...")

			s, h, b := sendRequest()

			resp_status.SetText(s)
			resp_headers.SetText(h)
			resp_body.SetText(b)
		}

		return event
	})

	/* LEFT SIDE START */

	req_form = tview.NewForm().
		AddDropDown("Method", []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, 0, nil).
		AddInputField("URL", "", 0, nil, nil).
		AddTextArea("Body", "", 0, 15, 0, nil)
		//		AddButton("Send", nil)

	req_form.SetBorder(true).SetTitle("Request Form").SetTitleAlign(tview.AlignLeft)

	req_form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModAlt {
			switch event.Key() {
			case tcell.KeyUp:
				return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
			case tcell.KeyDown:
				return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
			case tcell.KeyRight:
				app.SetFocus(headers_form)
				return nil
			}
		}

		switch event.Key() {
		case tcell.KeyTab:
			return tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone)
		case tcell.KeyBacktab:
			return nil
		}

		return event
	})

	req_form_body := req_form.GetFormItem(2).(*tview.TextArea)
	req_form_body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			if req_form_body.HasSelection() {
				text, _, _ := req_form_body.GetSelection()
				app.GetScreen().SetClipboard([]byte(text))
			}
			return nil
		}

		return event
	})

	instructions = tview.NewTextView().
		SetText("Alt + Up/Down - Next/prev field\n" +
			"Alt + Right/Left - Next/prev menu\n" +
			"Ctrl + s - Send request\n" +
			"Esc | Ctrl + q/c - Quit")

	instructions.SetBorder(true).SetTitle("Commands").SetTitleAlign(tview.AlignLeft)

	left_flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(req_form, 0, 5, true).
		AddItem(instructions, 0, 1, false)

	/* LEFT SIDE END */

	/* CENTER START */

	headers_form = tview.NewForm().
		//		SetHorizontal(true).
		AddInputField("Key", "", 0, nil, nil).
		AddInputField("Value", "", 0, nil, nil).
		AddInputField("Key", "", 0, nil, nil).
		AddInputField("Value", "", 0, nil, nil).
		AddInputField("Key", "", 0, nil, nil).
		AddInputField("Value", "", 0, nil, nil).
		AddInputField("Key", "", 0, nil, nil).
		AddInputField("Value", "", 0, nil, nil).
		AddInputField("Key", "", 0, nil, nil).
		AddInputField("Value", "", 0, nil, nil)

	headers_form.SetBorder(true).SetTitle("Headers").SetTitleAlign(tview.AlignCenter)

	headers_form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModAlt {
			switch event.Key() {
			case tcell.KeyUp:
				return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
			case tcell.KeyDown:
				return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
			case tcell.KeyRight:
				app.SetFocus(resp_headers)
				return nil
			case tcell.KeyLeft:
				app.SetFocus(req_form)
				return nil
			}
		}

		switch event.Key() {
		case tcell.KeyTab:
			return tcell.NewEventKey(tcell.KeyRune, ' ', tcell.ModNone)
		case tcell.KeyBacktab:
			return nil
		}

		return event
	})

	/* CENTER END */

	/* RIGHT SIDE START */

	resp_status = tview.NewTextView()
	resp_status.SetBorder(true).SetTitle("Response Status").SetTitleAlign(tview.AlignRight)
	resp_status.SetChangedFunc(func() {
		app.ForceDraw()
	})

	resp_headers = tview.NewTextView().SetScrollable(true)
	resp_headers.SetBorder(true).SetTitle("Response Headers").SetTitleAlign(tview.AlignRight)
	resp_headers.SetChangedFunc(func() {
		app.ForceDraw()
	})

	resp_body = tview.NewTextView().SetScrollable(true)
	resp_body.SetBorder(true).SetTitle("Response Body").SetTitleAlign(tview.AlignRight)
	resp_body.SetChangedFunc(func() {
		app.ForceDraw()
	})

	right_flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(resp_status, 0, 1, false).
		AddItem(resp_headers, 0, 5, false).
		AddItem(resp_body, 0, 5, false)

	right_flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Modifiers() == tcell.ModAlt {
			switch event.Key() {
			case tcell.KeyUp:
				app.SetFocus(resp_headers)
				return nil
			case tcell.KeyDown:
				app.SetFocus(resp_body)
				return nil
			case tcell.KeyLeft:
				app.SetFocus(headers_form)
				return nil
			}
		}

		switch event.Key() {
		case tcell.KeyTab:
			return nil
		case tcell.KeyBacktab:
			return nil
		}

		return event
	})

	/* RIGHT SIDE END */

	main_flex = tview.NewFlex().
		AddItem(left_flex, 0, 1, true).
		AddItem(headers_form, 0, 1, true).
		AddItem(right_flex, 0, 1, false)

	if err := app.SetTitle("vcurl").SetRoot(main_flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
