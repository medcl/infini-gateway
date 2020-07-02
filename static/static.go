package public

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	log "github.com/cihub/seelog"
	"infini.sh/framework/core/errors"
	"infini.sh/framework/core/util"
	"infini.sh/framework/core/vfs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
    "strings"
)

func (vfs StaticFS) prepare(name string) (*vfs.VFile, error) {
	
	log.Trace("check virtual file, ", name)
	
	name = path.Clean(name)

	if strings.HasSuffix(name,"/"){
		name=name+"index.html"
	}

	log.Trace("clean virtual file, ", name)

	f, present := data[name]
	if !present {
		log.Trace("virtual file not found, ", name)
		return nil, os.ErrNotExist
	}
	var err error
	vfs.once.Do(func() {
		f.FileName = path.Base(name)
		if f.FileSize == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.Compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			log.Error(err)
			return
		}
		f.Data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return f, nil
}

func (vfs StaticFS) Open(name string) (http.File, error) {

	name = path.Clean(name)

	if vfs.CheckLocalFirst {

		if vfs.TrimLeftPath !=""{
			name= util.TrimLeftStr(name, vfs.TrimLeftPath)	
		}

		localFile := path.Join(vfs.StaticFolder, name)

		log.Trace("check local file, ", localFile)

		if util.FileExists(localFile) {

			f2, err := os.Open(localFile)
			if err == nil {
				return f2, err
			}
		}

		log.Debug("local file not found,", localFile)
	}

	if vfs.SkipVFS{
		log.Trace("file was not found on vfs, ",name)
		return nil,errors.New("file not found")
	}

	f, err := vfs.prepare(name)
	if err != nil {
		return nil, err
	}
	log.Trace(f.FileName,",",f.ModifyTime,",",f.FileSize,",",f.Mode(),",",f.Name())
	return f.File()
}

type StaticFS struct {
	once            sync.Once
	StaticFolder    string
	TrimLeftPath    string
	CheckLocalFirst bool
	SkipVFS         bool

}

var data = map[string]*vfs.VFile{
}
