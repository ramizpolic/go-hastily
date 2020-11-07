package global

import (
	"fmt"

	. "github.com/jedib0t/go-pretty/v6/text"
	"github.com/schollz/progressbar/v3"
)

// CLI exports message controller.
var CLI = NewController("Message controller")

type msgControl struct {
	Controller string
}

// Clier wraps functions of message controller.
type Clier interface {
	Header(string, string)
	Title(string, ...interface{})
	Subtitle(string, ...interface{})
	Desc(string, ...interface{})
	Info(string, ...interface{})
	Success(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Progress(string, int) *progressbar.ProgressBar
	NewLine()
}

// NewController creates new Clier controller.
func NewController(controller string) Clier {
	return &msgControl{
		Controller: controller,
	}
}

func (cli *msgControl) Progress(title string, size int) *progressbar.ProgressBar {
	return progressbar.NewOptions(size,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(fmt.Sprintf("  [yellow][+][reset] %s...", title)),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionShowCount(),
		//progressbar.OptionOnCompletion(func() { fmt.Printf("\n\n") }),
	)
}

func (cli *msgControl) NewLine() {
	fmt.Printf("\n")
}

func (cli *msgControl) Header(title string, desc string) {
	cli.Title(title)
	if desc == "" {
		cli.NewLine()
	} else {
		cli.Desc(desc)
	}
}

func (cli *msgControl) Title(format string, a ...interface{}) {
	fmt.Printf("🚀 %s\n",
		Colors{Bold, Underline, BgHiBlack}.Sprintf(format, a...),
	)
}

func (cli *msgControl) Subtitle(format string, a ...interface{}) {
	fmt.Printf("  %s %s\n",
		Colors{FgYellow}.Sprintf(">"),
		Colors{Bold}.Sprintf(format, a...),
	)
}

func (cli *msgControl) Desc(format string, a ...interface{}) {
	fmt.Printf("%s\n\n",
		Colors{Italic, Faint}.Sprintf(format, a...),
	)
}

func (cli *msgControl) Info(format string, a ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprintf(format, a...))
}

func (cli *msgControl) Success(format string, a ...interface{}) {
	fmt.Printf("✔️ %s\n",
		Colors{Bold}.Sprintf(format, a...),
	)
}

func (cli *msgControl) Error(format string, a ...interface{}) {
	fmt.Printf("❌ %s\n",
		Colors{Bold}.Sprintf(format, a...),
	)
}

func (cli *msgControl) Warn(format string, a ...interface{}) {
	fmt.Printf("⚠️ %s\n",
		Colors{Faint}.Sprintf(format, a...),
	)
}
