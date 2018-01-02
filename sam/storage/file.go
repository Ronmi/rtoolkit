// Package storage provides some simple implements of SAM version storage.
package storage

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/Ronmi/rtoolkit/sam"
)

// FakeStorage is only for test purpose. It always return 0.
type FakeStorage struct{}

func (s FakeStorage) Load(app string) (ret int, err error) {
	return
}

func (s FakeStorage) Save(app string, ver int) error {
	return nil
}

type perFile string

func (s perFile) Load(app string) (ret int, err error) {
	fn := fmt.Sprintf(string(s), app)

	// detect file existence
	if _, e := os.Stat(fn); e != nil {
		return
	}

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}

	ret, err = strconv.Atoi(string(data))
	return
}

func (s perFile) Save(app string, ver int) error {
	fn := fmt.Sprintf(string(s), app)
	return ioutil.WriteFile(fn, []byte(strconv.Itoa(ver)), 0600)
}

// PerFile stores version info in plain text, per-app based.
//
// The only argument defines how a file should be named (using fmt.Sprintf).
//
// If specified file does not exists:
//
//     - Load() will return (0, nil)
//     - Save() will try to create the file with permission 0600
func PerFile(format string) sam.Storage {
	if strings.Count(format, "%s") != 1 {
		panic(errors.New(format + " is not a valid PerFile format"))
	}

	return perFile(format)
}

var (
	Default = PerFile("%s-db.ver")
)
