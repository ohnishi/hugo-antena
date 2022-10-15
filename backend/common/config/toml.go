package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ReadToml はTOMLのファイルを読んでdataに格納する。
func ReadToml(path string, data interface{}) error {
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read TOML file: %w", err)
	}
	err = v.Unmarshal(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal TOML data: %w", err)
	}
	return nil
}

// TomlWriter はTOMLのファイルを書くための機能を提供する
type TomlWriter struct {
	v *viper.Viper
}

// NewTomlWriter は新しいTomlWriterを作って返す。
func NewTomlWriter() *TomlWriter {
	return &TomlWriter{viper.New()}
}

// Set は設定ファイルのnameというフィールドにdataの値がセットされるようにする。
func (w *TomlWriter) Set(name string, data interface{}) {
	w.v.Set(name, data)
}

// SetWithPrefix は設定ファイルのprefix.nameというフィールドにdataの値がセットされるようにする。
func (w *TomlWriter) SetWithPrefix(name string, prefix string, data interface{}) {
	w.v.Set(prefix+"."+name, data)
}

// WriteFile はセットされた各項目を指定されたパスのファイルに書き込む。
func (w *TomlWriter) WriteFile(path string) error {
	w.v.SetConfigFile(path)
	w.v.SetConfigType("toml")
	err := w.v.WriteConfig()
	if err != nil {
		return fmt.Errorf("failed to write config flie: %w", err)
	}
	return nil
}
