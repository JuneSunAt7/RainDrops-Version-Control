package main

import (
    "rdvc/init_dir"
	"os"

	"github.com/pterm/pterm"
)

func printHelp() {
    pterm.FgGreen.Println("rdvc - RainDrops Version Control")
    pterm.FgBlue.Println("Usage:")
    pterm.FgCyan.Println("  rdvc init -p <path_to_directory>     Initialize a controlled directory")
    pterm.FgCyan.Println("  rdvc help                             Show this help")
}
func check_args(){
    if len(os.Args) < 2 {
        printHelp()
        return
    }

    if os.Args[1] == "help" {
        printHelp()
        return
    } else if os.Args[1] == "init" {
		if os.Args[2] != "-p" {
			pterm.Error.Printfln("Error: Invalid argument. Use -p <path_to_directory>.")
			return
		}
        pathFlag := os.Args[3]
        init_dir.InitInvisible(pathFlag)
    } else {
        pterm.Error.Printfln("Unknown command: %s\n", os.Args[1])
        printHelp()
    }
}
func main() {
    check_args()
}