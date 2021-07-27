package twofa

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

const (
	BrightBlack = "\033[48;5;0m  \033[0m"
	BrightWhite = "\033[48;5;7m  \033[0m"
)

func PrintQR(content string) error {
	code, err := qrcode.New(content, qrcode.High)
	if err != nil {
		return err
	}

	bitmap := codeString(code.Bitmap())

	fmt.Println(bitmap)

	return nil
}

func codeString(bitmap [][]bool) (result string) {
	for indexRow, row := range bitmap {
		lineRow := len(row)
		if indexRow == 0 || indexRow == 1 || indexRow == 2 ||
			indexRow == lineRow-1 || indexRow == lineRow-2 || indexRow == lineRow-3 {
			continue
		}

		for indexCol, col := range row {
			lineCount := len(bitmap)
			if indexCol == 0 || indexCol == 1 || indexCol == 2 ||
				indexCol == lineCount-1 || indexCol == lineCount-2 || indexCol == lineCount-3 {
				continue
			}

			if !col {
				result += fmt.Sprint(BrightWhite)
			} else {
				result += fmt.Sprint(BrightBlack)
			}
		}

		result += fmt.Sprintln()
	}

	return
}
