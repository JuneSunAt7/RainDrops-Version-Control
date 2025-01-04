package init_dir

import (
    "os"
    "github.com/pterm/pterm"
    "syscall"
    "time"
)

func InitInvisible(directory string) {
    if directory == "" {
        pterm.Error.Println("Error: No directory specified. Use -p <path_to_directory>.")
        return
    }
    // Создайте директорию ".rdvc"
    err := os.Mkdir(directory + "\\.rdvc", os.ModePerm) // Лучше использовать "\\" для Windows
    if err != nil {
        pterm.Error.Printfln("Failed to create directory: %v\n", err)
        return
    }

    // Измените атрибуты папки на скрытые
    err = syscall.SetFileAttributes(syscall.StringToUTF16Ptr(directory+"\\.rdvc"), syscall.FILE_ATTRIBUTE_HIDDEN)
    if err != nil {
        pterm.Error.Println("Error setting directory attributes:", err)
        return
    }
    p, _ := pterm.DefaultProgressbar.WithTotal(10).WithTitle("Initializing directory").Start()

    for i := 0; i < p.Total; i++ {
        p.UpdateTitle("Initializing directory...")
        p.Increment()
        time.Sleep(time.Millisecond * 50)
    }
    pterm.Success.Printfln("Successful init: %s\n", directory)
}