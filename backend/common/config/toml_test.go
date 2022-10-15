package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ohnishi/antena/backend/common/config"
)

func testTempFile(t *testing.T) (string, func()) {
	f, err := ioutil.TempFile("", "test_temp_file")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	f.Close()
	return f.Name(), func() { os.Remove(f.Name()) }
}

func testTempDir(t *testing.T) (string, func()) {
	d, err := ioutil.TempDir("", "test_temp_dir")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return d, func() { os.RemoveAll(d) }
}

type testToml struct {
	Hostname string
	Port     int
	Logging  bool
}

var testTomlValue = testToml{
	Hostname: "localhost",
	Port:     12345,
	Logging:  true,
}

var testTomlStr = `hostname = "localhost"
port = 12345
logging = true
`

func TestReadToml(t *testing.T) {
	testFile, remover := testTempFile(t)
	defer remover()
	ioutil.WriteFile(testFile, []byte(testTomlStr), os.ModePerm)
	c := testToml{}
	config.ReadToml(testFile, &c)

	if testTomlValue.Hostname != c.Hostname {
		t.Errorf("Hostname differs: expected: %s, actual: %s", testTomlValue.Hostname, c.Hostname)
	}
	if testTomlValue.Port != c.Port {
		t.Errorf("Port differs: expected %d, actual: %d", testTomlValue.Port, c.Port)
	}
	if testTomlValue.Logging != c.Logging {
		t.Errorf("Logging differs: expected %v, actual: %v", testTomlValue.Logging, c.Logging)
	}
}

func TestTomlWriter_WriteFile(t *testing.T) {
	// testTempFileでは拡張子を指定できないためtestTempDirを使う
	testDir, remover := testTempDir(t)
	testFile := filepath.Join(testDir, "test_toml_writer_write_file.toml")
	defer remover()

	w := config.NewTomlWriter()
	w.Set("hostname", testTomlValue.Hostname)
	w.Set("port", testTomlValue.Port)
	w.Set("logging", testTomlValue.Logging)
	err := w.WriteFile(testFile)
	if err != nil {
		t.Errorf("err: %s", err)
	}
	configBytes, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Errorf("err: %s", err)
	}
	configStr := strings.Trim(string(configBytes), "\n")
	configLines := strings.Split(configStr, "\n")
	if len(configLines) != 3 {
		t.Errorf("Unexpected config line length: expected: 3, actual: %d", len(configLines))
	}
	// 書かれる順番は変わり得るので各行を比較
	for _, line := range configLines {
		if !strings.Contains(testTomlStr, line) {
			t.Errorf("Unexpected line written: %s", line)
		}
	}
}
