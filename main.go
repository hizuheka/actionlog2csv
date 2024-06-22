package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
	if len(os.Args) != 4 {
		fmt.Println("使用法: actionlog2csv.exe <ログフォルダパス> <出力ファイルパス> <ワーカー数>")
		return
	}

	folderPath := os.Args[1]
	outputFilePath := os.Args[2]
	numWorkers, err := strconv.Atoi(os.Args[3])
	if err != nil || numWorkers <= 0 {
		fmt.Println("ワーカー数は正の整数で指定してください")
		return
	}

	var wg sync.WaitGroup

	fileChan := make(chan string)
	resultChan := make(chan map[LogEntry]struct{})

	// ワーカーを起動
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				entries := processFile(path)
				resultChan <- entries
			}
		}()
	}

	// フォルダ内の全てのログファイルを解析
	go func() {
		defer close(fileChan)
		err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fileChan <- path
			}
			return nil
		})

		if err != nil {
			fmt.Printf("フォルダ内のファイルを処理中にエラーが発生しました: %v\n", err)
		}
	}()

	// 結果を集約
	var logEntries = make(map[LogEntry]struct{})
	go func() {
		for entries := range resultChan {
			for entry := range entries {
				logEntries[entry] = struct{}{}
			}
		}
	}()

	wg.Wait()
	close(resultChan)

	if err := writeOutputFile(outputFilePath, logEntries); err != nil {
		fmt.Printf("ファイル出力でエラーが発生しました: %v\n", err)
		return
	}
}

func writeOutputFile(outputFilePath string, logEntries map[LogEntry]struct{}) error {
	// 出力ファイルを開く
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// ヘッダーを書き込む
	writer.Write([]string{"action", "src", "dst", "interface", "dir", "rule"})

	// ログエントリを書き込む
	for entry := range logEntries {
		if err := writer.Write([]string{entry.Action, entry.Src, entry.Dst, entry.Interface, entry.Dir, entry.Rule}); err != nil {
			return err
		}
	}

	return nil
}

func processFile(path string) map[LogEntry]struct{} {
	fmt.Printf("対象ファイル: %s\r\n", path)
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("ログファイルを開くことができません: path=%s, error=%v\n", path, err)
		return nil
	}
	defer file.Close()

	logEntries := make(map[LogEntry]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// "action="を含む行のみを対象とする
		if strings.Contains(line, "action=") {
			if entry, err := createEntry(line); err != nil {
				fmt.Printf("ログファイルの解析でエラーが発生しました: path=%s, error=%v\n", path, err)
			} else {
				// 重複を排除するために、entryをキーにする
				logEntries[entry] = struct{}{}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("ログファイルの読み取り中にエラーが発生しました: path=%s, error=%v\n", path, err)
	}

	return logEntries
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
