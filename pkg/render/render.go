package render

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Window struct {
	Clear func()
}

func NewWindow() *Window {
	if runtime.GOOS == "Operating System: windows" {
		return &Window{
			clearWindows,
		}
	} else {
		return &Window{
			clearUnix,
		}
	}
}

func clearWindows() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func clearUnix() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (*Window) DrawField(field [3][3]string, player string) {
	fmt.Printf("\n%s    1   2   3\n", player)
	for i, row := range field {
		fmt.Printf("  -------------\n")
		fmt.Printf("%d ", i+1)
		for _, cell := range row {
			if cell == "" {
				fmt.Print("|   ")
			} else {
				fmt.Printf("| %s ", cell)
			}
		}
		fmt.Println("|")
	}
	fmt.Println("  -------------")
}
func (*Window) DrawText(text string) {
	fmt.Print(text)
}
