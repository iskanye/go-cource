package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func dirTree(file io.Writer, rootDir string, printFiles bool) error {
	var recursion func(string, string)

	recursion = func(root string, prefix string) {
		dir, err := os.ReadDir(root)
		if err != nil {
			return
		}

		// Фильтрация: оставляем только директории, если printFiles false
		var entries []os.DirEntry
		for _, entry := range dir {
			if entry.IsDir() || printFiles {
				entries = append(entries, entry)
			}
		}

		sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

		for i, entry := range entries {
			// Определяем, последний ли это элемент
			isLast := i == len(entries)-1

			// Рисуем соединение для текущего элемента
			if isLast {
				fmt.Fprint(file, prefix+"└───"+entry.Name())
			} else {
				fmt.Fprint(file, prefix+"├───"+entry.Name())
			}

			// Если это директория, рекурсивно обрабатываем её содержимое
			if entry.IsDir() {
				fmt.Fprintln(file)
				newPrefix := prefix

				if isLast {
					newPrefix += "\t"
				} else {
					newPrefix += "│\t"
				}

				recursion(root+string(os.PathSeparator)+entry.Name(), newPrefix)
			} else {
				// Если это файл, анализируем его размер
				info, err := entry.Info()
				if err != nil {
					return
				}
				size := info.Size()

				if size == 0 {
					fmt.Fprintln(file, " (empty)")
				} else {
					fmt.Fprintf(file, " (%db)\n", size)
				}
			}
		}
	}

	recursion(rootDir, "")
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
