package auth

import (
	"github.com/diamondburned/cchat"
	"github.com/diamondburned/cchat-tui/tui/app"
	"github.com/diamondburned/cchat-tui/tui/center"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type TextGetter interface {
	GetText() string
}

type Form struct {
	*center.Center
	Pages *tview.Pages

	Form  *tview.Form
	Busy  *center.Text
	Modal *tview.Modal

	Auther cchat.Authenticator
	onAuth func(cchat.Session)
}

// SpawnForm creates a new form with the authenticator. This function is not
// thread-safe, and auth will be called on the UI thread as well.
func SpawnForm(auther cchat.Authenticator, auth func(cchat.Session)) {
	form := newForm(auther, func(ses cchat.Session) {
		app.SwitchToMain()
		auth(ses)
	})
	form.Form.SetCancelFunc(app.SwitchToMain)
	form.spin()
	app.SwitchToView(form)
}

// newForm create a new form with the authenticator. This function is not
// thread-safe, and onAuth will be called on the UI thread as well.
func newForm(auther cchat.Authenticator, auth func(cchat.Session)) *Form {
	form := tview.NewForm()
	form.SetBackgroundColor(-1)
	form.SetFieldBackgroundColor(-1)
	form.SetButtonBackgroundColor(tcell.ColorWhite)
	form.SetButtonTextColor(tcell.ColorBlack)

	busy := center.NewText("Logging in...")

	// used for errors
	modal := tview.NewModal()
	modal.SetBackgroundColor(-1)
	modal.SetButtonBackgroundColor(-1)
	modal.SetTextColor(tcell.ColorRed)
	modal.AddButtons([]string{"Ok"})

	pages := tview.NewPages()
	pages.SetBackgroundColor(-1)
	pages.AddPage("form", form, true, true)
	pages.AddPage("busy", busy, true, false)
	pages.AddPage("modal", modal, true, false)
	pages.SetBorderPadding(0, 0, 0, 10)

	center := center.New(pages)
	center.MaxHeight = 12
	center.MaxWidth = 35 + 10 // 10 from page right padding

	f := &Form{
		center,
		pages,
		form,
		busy,
		modal,
		auther,
		auth,
	}

	return f
}

// Spin populates the form with the given input fields from Authenticator then
// binds the appropriate callbacks. If withErr is not nil, an error will be
// displayed.
func (f *Form) spin() {
	// Wipe the form first.
	f.Form.Clear(true)

	// Add input fields.
	for _, form := range f.Auther.AuthenticateForm() {
		if form.Secret {
			f.Form.AddPasswordField(form.Name, "", 25, '*', nil)
		} else {
			f.Form.AddInputField(form.Name, "", 25, nil, nil)
		}
	}

	// Add the Login button.
	f.Form.AddButton("Login", f.ok)
}

func (f *Form) ok() {
	// Switch to the loading screen.
	f.Pages.SwitchToPage("busy")

	s, err := f.Auther.Authenticate(f.getOutputs())
	if err == nil {
		f.onAuth(s)
		return
	}

	// Show the modal dialog with the error.
	f.Pages.SwitchToPage("modal")
	f.Modal.SetText(err.Error())
	// Essentially this does the same thing, but we're doing this here for
	// clarity.
	f.Modal.SetDoneFunc(func(int, string) {
		f.Pages.SwitchToPage("form")
		f.spin()
		f.Modal.SetText("")
	})
}

func (f *Form) getOutputs() []string {
	var outputs = make([]string, f.Form.GetFormItemCount())
	for i := 0; i < len(outputs); i++ {
		outputs[i] = f.Form.GetFormItem(i).(TextGetter).GetText()
	}
	return outputs
}
