package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

// IsNotExist はerrがファイルが存在しないことを表すerrorであるときtrueを返す。
// os.IsNotExistとの違いはerrがerrors.Wrapなどでラップされている場合にも動作する点。
func IsNotExist(err error) bool {
	return os.IsNotExist(errors.Cause(err))
}

// ReadFileJSON はJSONを読み込んでoutputに入れる
func ReadFileJSON(file string, output interface{}) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(json.Unmarshal(bytes, output))
}

// WriteFileJSON はJSONをpathに書き込む
func WriteFileJSON(file string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	err = os.MkdirAll(path.Dir(file), os.ModePerm)
	if err != nil {
		return errors.WithStack(nil)
	}
	return errors.WithStack(ioutil.WriteFile(file, bytes, os.ModePerm))
}

// WriteFilePrettyJSON は整形したJSONをpathに書き込む
func WriteFilePrettyJSON(file string, data interface{}) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}
	err = os.MkdirAll(path.Dir(file), os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(ioutil.WriteFile(file, bytes, os.ModePerm))
}
