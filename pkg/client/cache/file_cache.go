/*
Copyright 2021 CodeNotary, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cache

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rogpeppe/go-internal/lockedfile"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/golang/protobuf/proto"
)

// STATE_FN ...
const STATE_FN = ".state-"

type fileCache struct {
	Dir string
}

// NewFileCache returns a new file cache
func NewFileCache(dir string) Cache {
	return &fileCache{Dir: dir}
}

func (w *fileCache) Get(serverUUID, db string) (*schema.ImmutableState, error) {
	fn := filepath.Join(w.Dir, string(getRootFileName([]byte(STATE_FN), []byte(serverUUID))))

	raw, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(raw), "\n")
	for _, line := range lines {
		if strings.Contains(line, db+":") {
			r := strings.Split(line, ":")
			if r[1] == "" {
				return nil, ErrPrevStateNotFound
			}
			oldState, err := base64.StdEncoding.DecodeString(r[1])
			if err != nil {
				return nil, ErrLocalStateCorrupted
			}
			state := &schema.ImmutableState{}
			if err = proto.Unmarshal(oldState, state); err != nil {
				return nil, ErrLocalStateCorrupted
			}
			return state, nil
		}
	}
	return nil, ErrPrevStateNotFound
}

func (w *fileCache) Set(serverUUID, db string, state *schema.ImmutableState) error {
	raw, err := proto.Marshal(state)
	if err != nil {
		return err
	}
	fn := filepath.Join(w.Dir, string(getRootFileName([]byte(STATE_FN), []byte(serverUUID))))

	input, _ := ioutil.ReadFile(fn)
	lines := strings.Split(string(input), "\n")

	newState := db + ":" + base64.StdEncoding.EncodeToString(raw) + "\n"
	var exists bool
	for i, line := range lines {
		if strings.Contains(line, db+":") {
			exists = true
			lines[i] = newState
		}
	}
	if !exists {
		lines = append(lines, newState)
	}
	output := strings.Join(lines, "\n")

	if err = ioutil.WriteFile(fn, []byte(output), 0644); err != nil {
		return err
	}
	return nil
}

func getRootFileName(prefix []byte, serverUUID []byte) []byte {
	l1 := len(prefix)
	l2 := len(serverUUID)
	var fn = make([]byte, l1+l2)
	copy(fn[:], STATE_FN)
	copy(fn[l1:], serverUUID)
	return fn
}

func (w *fileCache) GetLocker(serverUUID string) Locker {
	fn := filepath.Join(w.Dir, string(getRootFileName([]byte(STATE_FN), []byte(serverUUID))))
	fm := lockedfile.MutexAt(fn)
	return &FileLocker{lm: fm}
}

type FileLocker struct {
	lm         *lockedfile.Mutex
	unlockFunc func()
}

func (fl *FileLocker) Lock() (err error) {
	fl.unlockFunc, err = fl.lm.Lock()
	return err
}

func (fl *FileLocker) Unlock() (err error) {
	if fl.unlockFunc == nil {
		return fmt.Errorf("try to lock a not locked file")
	}
	fl.unlockFunc()
	return nil
}
