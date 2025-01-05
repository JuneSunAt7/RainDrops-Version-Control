package init_dir

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "time"

    "github.com/pterm/pterm"
)

type Keep struct {
    Message string    `json:"message"`
    Author  string    `json:"author"`
    Date    time.Time `json:"date"`
    Files   []File    `json:"files"`
    Branch  string    `json:"branch"` // Добавляем информацию о ветке
}

type File struct {
    Name    string `json:"name"`
    Content string `json:"content"`
}

type VCS struct {
    RepoPath      string // Путь к репозиторию
    CurrentBranch string // Текущая ветка
}

func NewVCS(repoPath string) *VCS {
    return &VCS{
        RepoPath:      repoPath,
        CurrentBranch: "main", // "main" - основная ветка по умолчанию
    }
}

// Создание новой ветки
func (v *VCS) CreateBranch(branchName string) error {
    branchPath := filepath.Join(v.RepoPath, ".rdvc", "branches", branchName)
    err := os.MkdirAll(branchPath, os.ModePerm)
    if err != nil {
        return err
    }
    pterm.Success.Printfln("Created line: %s", branchName)
    return nil
}

// Переключение на ветку
func (v *VCS) CheckoutBranch(branchName string) error {
    // Проверьте, существует ли ветка
    branchPath := filepath.Join(v.RepoPath, ".rdvc", "branches", branchName)
    if _, err := os.Stat(branchPath); os.IsNotExist(err) {
        return fmt.Errorf("line %s not exists", branchName)
    }
    v.CurrentBranch = branchName
    pterm.Success.Printfln("Checkout to line: %s", branchName)
    return nil
}

// Модификация функции MakeKeep для поддержки веток
func (v *VCS) MakeKeep(message string, author string) error {
    keep := Keep{
        Message: message,
        Author:  author,
        Date:    time.Now(),
        Branch:  v.CurrentBranch, // Сохраняем текущую ветку
    }

    err := filepath.Walk(v.RepoPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }

        content, err := ioutil.ReadFile(path)
        if err != nil {
            return err
        }

        // добавляем в keep
        keep.Files = append(keep.Files, File{
            Name:    path,
            Content: string(content),
        })
        return nil
    })

    if err != nil {
        return err
    }

    // Сохранить keep в файл
    keepFileName := fmt.Sprintf("%s_%d.keep.json", v.CurrentBranch, time.Now().Unix()) // Добавляем ветку в имя файла
    keepFilePath := filepath.Join(v.RepoPath, ".rdvc", "keeps", keepFileName)

    os.MkdirAll(filepath.Dir(keepFilePath), os.ModePerm)

    jsonData, err := json.Marshal(keep)
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(keepFilePath, jsonData, 0644)
    if err != nil {
        return err
    }

    pterm.Success.Printfln("Successfuly keeped in: %s\n", keepFilePath)
    return nil
}