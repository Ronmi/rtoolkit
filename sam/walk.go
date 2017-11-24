package sam

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

// Walker defines an interface so sam can support both physics and virtual fs
type Walker func(root string, f filepath.WalkFunc) error

var reFileName *regexp.Regexp

func init() {
	reFileName = regexp.MustCompile(`^([0-9]+)-.+\.sql$`)
}

type errMismatch string

func (e errMismatch) Error() string {
	return string(e) + " does not match sql name pattern"
}

func walkfunc(path string, info os.FileInfo, err error) (int, error) {
	if err != nil {
		return 0, err
	}

	if info.IsDir() {
		return 0, filepath.SkipDir
	}

	arr := reFileName.FindStringSubmatch(filepath.Base(path))
	if len(arr) == 0 {
		return 0, errMismatch(path)
	}

	i, _ := strconv.ParseInt(arr[1], 10, 64)
	return int(i), nil
}

type sqlFile struct {
	id int
	fn string
}

type byID []sqlFile

func (b byID) Len() int           { return len(b) }
func (b byID) Less(i, j int) bool { return b[i].id < b[j].id }
func (b byID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

func sqlFiles(w Walker, root string, cur int) []sqlFile {
	ret := make([]sqlFile, 0)
	base := filepath.Base(root)

	w(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && filepath.Base(path) == base {
			return nil
		}

		i, e := walkfunc(path, info, err)
		if e != nil {
			if _, ok := e.(errMismatch); ok {
				return nil
			}

			return e
		}

		if i > cur {
			ret = append(ret, sqlFile{id: i, fn: path})
		}

		return nil
	})

	sort.Sort(byID(ret))

	return ret
}
