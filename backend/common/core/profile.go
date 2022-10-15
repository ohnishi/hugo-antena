package core

import (
	"log"
	"os"
	"runtime/pprof"
)

// SetupProfile はCPUプロファイラーを有効化してfilenameに結果を保存する。
//
// 使い方:
// 下記のようにmain関数開始直後にSetupProfileを呼び、返された関数をdeferで実行する。
//
// func main() {
//     defer SetupProfile("profile.prof")()  // ***後ろの括弧を忘れないようにする***
//     ...
// }
//
// 保存されたファイルを下記コマンドで解析する。
//
// go tool pprof myprog profile.prof
func SetupProfile(filename string) func() {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	return func() { pprof.StopCPUProfile() }
}
