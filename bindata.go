package ry

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _lisp_runtime_syp = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x6c\x90\x3d\x6e\xc3\x30\x0c\x85\x77\x9d\xe2\x39\x1d\x4a\x0d\x06\xdc\xd4\x99\xd2\xb1\x47\xe8\x56\x74\x50\x22\xc6\x15\xac\x48\xae\x2c\x23\xc8\xed\xab\x1f\xa3\x30\xd0\x6c\xfc\x79\xfc\x1e\x49\xd2\x7c\xc1\xcf\x62\x62\x7b\xfe\x56\x0e\x74\x55\x23\x97\x50\x4a\x21\x72\xd3\x95\x2e\x3e\xbf\x04\x40\x73\xf4\x53\xcb\xda\x44\x1f\x64\xc9\xd9\xe9\x66\x33\xfe\x4c\x32\xcf\xad\x83\xa9\xa2\x2d\xb7\x23\xdf\xa7\xc0\xf3\xcc\xf3\x4a\xb1\x9c\x78\x23\xe8\xad\x01\x55\xd8\x46\xb3\x9a\x67\x61\x31\x48\x64\xb6\x16\x2f\x38\x60\x44\x87\x4e\xd6\xc6\x3f\x76\x31\x3e\x1e\xf1\xce\xa7\x65\x18\x8c\x1b\x40\xce\xe3\xa6\xee\xf0\xf5\xc0\x98\x6b\x17\x1f\xe0\xfc\x4d\x0a\x1a\x3c\xe8\xc4\x49\x58\xee\xb0\xcc\x13\x0e\x5d\x57\xf0\x94\xe5\x25\x98\x82\x71\xd1\x3a\xec\x3e\xcc\x95\xfd\x12\x9b\x5d\xb1\x49\x7f\x50\x21\xfe\x3d\x82\xce\x96\x55\x58\x53\xec\xfb\x1e\xfb\xd7\x3e\x95\x37\xcb\xf7\x78\xd2\x75\xfb\xea\xfc\x60\x7d\xf1\x1b\x00\x00\xff\xff\x0b\xf7\xe5\x1e\x8b\x01\x00\x00")

func lisp_runtime_syp_bytes() ([]byte, error) {
	return bindata_read(
		_lisp_runtime_syp,
		"lisp/runtime.syp",
	)
}

func lisp_runtime_syp() (*asset, error) {
	bytes, err := lisp_runtime_syp_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "lisp/runtime.syp", size: 395, mode: os.FileMode(420), modTime: time.Unix(1429375746, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"lisp/runtime.syp": lisp_runtime_syp,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"lisp": &_bintree_t{nil, map[string]*_bintree_t{
		"runtime.syp": &_bintree_t{lisp_runtime_syp, map[string]*_bintree_t{
		}},
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

