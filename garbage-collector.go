/*
   Garbage Collector - clear outdated files
   Version: 1.0.0

   Copyright (c) 2024 https://github.com/utilmind/

   This software is licensed under the MIT License.
   You can find the full license text at https://opensource.org/licenses/MIT

    TODO
        1. Complete regular script.
        2. Add support for multiple directories.
 */
package main

import (
    "flag"
    "fmt"
    "strconv"
    "os"
    "path/filepath"
    "strings"
    "github.com/fatih/color"
//    "time"
)

// @CONFIG
const DefExpireDays = 90;

// @private functions
func showError() {
    errorColor := color.New(color.FgRed, color.Bold)
    errorColor.Print("ERROR: ")
}

func die(str string) {
    showError()
    fmt.Println(str)
    os.Exit(0)
}

func main() {
    cliArgs := map[string]*string{
        "dir":      flag.String("dir", "", "directory name"),
        "sub":      flag.String("sub", "", "include subdirectories"),
        "ext":      flag.String("ext", "", "file extension(s). Comma-separated if multiple."),
        "expire":   flag.String("expire", strconv.Itoa(DefExpireDays), "expire after N days. 0 = don't check date, delete all."), // AK: actually it's integer, but I'd prefer to parse it myself
        "confirm":  flag.String("confirm", "", "'y' or 'yes' auto-confirms file deletions. Otherwise you'll need to confirm file deletions one by one."),
    }

    flag.Usage = func() {
        thisExeName := filepath.Base(os.Args[0])

		fmt.Fprintf(os.Stderr, "Usage of %s:\n", thisExeName)

        for arg, str := range cliArgs {
            fmt.Printf("Key: %s, Value: %d\n", arg, str)
        }

		flag.PrintDefaults()
		// Add custom description
		fmt.Fprintf(os.Stderr, "\nExample: %s -dir=/var/www/project-name/data/cache -ext=jpg,jpeg,png,gif,webp -expire=60\n", thisExeName)
	}

    flag.Parse()

    if "" == *cliArgs["dir"] {
        showError()
        fmt.Println("-dir argument is required.\n")
        flag.Usage()
        os.Exit(0)
    }

    if '/' == []rune(*cliArgs["dir"])[0] && 2 > strings.Count(*cliArgs["dir"], "/") {
        die(fmt.Sprintf("No, it doesn’t works with root or any top-level directory. It will not process \"%s\" or any other \"/directory/\" under root.", *cliArgs["dir"]))
    }

    expire, err := strconv.Atoi(*cliArgs["expire"])
    if nil != err {
        die(fmt.Sprintf("Invalid integer value in argument -expire=%s. Please use integer value to specify the number of days, or skip it to use default %d days.", *cliArgs["expire"], DefExpireDays))
    }

    fmt.Sprintf("test %s", expire);

    // Смотри, если не указан параметр -, то пусть будет подтверждение. Действительно ли удалить все файлы старше N дней из такой-то директории?


/* TODO:
    dir := "путь/к/директории"
    cutoff := time.Now().AddDate(0, 0, -30) // Время 30 дней назад

    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Проверяем, является ли объект файлом и старше ли он 30 дней
        if !info.IsDir() && info.ModTime().Before(cutoff) {
            fmt.Println("Удаляю:", path)
            if err := os.Remove(path); err != nil {
                return err
            }
        }

        return nil
    })

    if err != nil {
        fmt.Println("Ошибка:", err)
    } else {
        fmt.Println("Удаление завершено.")
    }
*/
}
