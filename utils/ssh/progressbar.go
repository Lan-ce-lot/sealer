package ssh

import "github.com/schollz/progressbar/v3"

var (
	Width                  = 50
	optionEnableColorCodes = progressbar.OptionEnableColorCodes(true)
	optionSetWidth         = progressbar.OptionSetWidth(Width)
	optionSetTheme         = progressbar.OptionSetTheme(progressbar.Theme{
		Saucer:        "=",
		SaucerHead:    ">",
		SaucerPadding: " ",
		BarStart:      "[",
		BarEnd:        "]",
	})
)
