package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LogEntry はログファイルの各エントリを表します。
type LogEntry struct {
	Src       string
	Dest      string
	Interface string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("使用法: actionlog2csv.exe <ログフォルダパス> <出力ファイルパス>")
		return
	}

	folderPath := os.Args[1]
	outputFilePath := os.Args[2]

	var logEntries = make(map[LogEntry]struct{})

	// フォルダ内の全てのログファイルを解析
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("ログファイルを開くことができません: %v", err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				// "action="を含む行のみを対象とする
				if strings.Contains(line, "action=") {
					fields := strings.Fields(line)
					var src, dest, iface string
					for _, field := range fields {
						switch {
						case strings.HasPrefix(field, "src="):
							src = strings.TrimPrefix(field, "src=")
						case strings.HasPrefix(field, "dest="):
							dest = strings.TrimPrefix(field, "dest=")
						case strings.HasPrefix(field, "interface="):
							iface = strings.TrimPrefix(field, "interface=")
						}
					}
					entry := LogEntry{Src: src, Dest: dest, Interface: iface}
					// 重複を排除するために、entryをキーにする
					logEntries[entry] = struct{}{}
				}
			}

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("ログファイルの読み取り中にエラーが発生しました: %v", err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("フォルダ内のファイルを処理中にエラーが発生しました: %v\n", err)
		return
	}

	// 出力ファイルを開く
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("出力ファイルを作成できません: %v\n", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// ヘッダーを書き込む
	writer.Write([]string{"src", "dest", "interface"})

	// ログエントリを書き込む
	for entry := range logEntries {
		writer.Write([]string{entry.Src, entry.Dest, entry.Interface})
	}
}
