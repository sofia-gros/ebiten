// package main provides a clean docserver runner that launches a local pkgsite godoc server
// for all modules in this workspace using the active mise Go toolchain.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("エラー: カレントディレクトリを取得できません: %v\n", err)
		os.Exit(1)
	}

	rootDir := dir
	if filepath.Base(dir) == "tools" {
		rootDir = filepath.Dir(dir)
	}

	port := "6060"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	url := fmt.Sprintf("http://localhost:%s", port)
	fmt.Println("==========================================================")
	fmt.Printf(" 📚 Ebitengine ライブラリ群 Godoc サーバーを起動します\n")
	fmt.Printf(" 🌐 URL: %s\n", url)
	fmt.Println(" ※ 終了するには Ctrl+C を押してください")
	fmt.Println("==========================================================")

	// システム標準/miseの go コマンドで pkgsite を一括起動
	cmd := exec.Command("go", "run", "golang.org/x/pkgsite/cmd/pkgsite@latest", "-http", "localhost:"+port, rootDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = rootDir

	if err := cmd.Run(); err != nil {
		fmt.Printf("Godoc サーバー終了: %v\n", err)
	}
}
