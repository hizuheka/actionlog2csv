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
	Dst       string
	Interface string
	Dir       string
	Action    string
	Rule      string
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
					if entry, err := createEntry(line); err != nil {
						return fmt.Errorf("ログファイルの解析でエラーが発生しました: %v", err)
					} else {
						// 重複を排除するために、entryをキーにする
						logEntries[entry] = struct{}{}
					}
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
	writer.Write([]string{"src", "dst", "interface", "dir", "action", "rule"})

	// ログエントリを書き込む
	for entry := range logEntries {
		writer.Write([]string{entry.Src, entry.Dst, entry.Interface})
	}
}

func createEntry(line string) (LogEntry, error) {
	fields := strings.Fields(line)
	var src, dst, iface, dir, action, rule string
	for _, field := range fields {
		switch {
		case strings.HasPrefix(field, "src="):
			src = strings.TrimPrefix(field, "src=")
		case strings.HasPrefix(field, "dst="):
			dst = strings.TrimPrefix(field, "dst=")
		case strings.HasPrefix(field, "interface="):
			iface = strings.TrimPrefix(field, "interface=")
		case strings.HasPrefix(field, "dir="):
			dir = strings.TrimPrefix(field, "dir=")
		case strings.HasPrefix(field, "action="):
			action = strings.TrimPrefix(field, "action=")
		case strings.HasPrefix(field, "rule="):
			rule = strings.TrimPrefix(field, "rule=")
		}
	}

	if src == "" || dst == "" || iface == "" || dir == "" || action == "" || rule == "" {
		return LogEntry{}, fmt.Errorf("ログファイルの形式が不正です: src=%s, dst=%s, interface=%s, dir=%s, action=%s, rule=%s", src, dst, iface, dir, action, rule)
	}
	entry := LogEntry{Src: src, Dst: dst, Interface: iface, Dir: dir, Action: action, Rule: rule}
	return entry, nil
}
