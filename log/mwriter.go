// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	`io`
	`io/ioutil`
	`log`
	`os`
	`path/filepath`
	`sync`
)

const (
	FileFlagsAppend = os.O_APPEND|os.O_CREATE|os.O_WRONLY
	FileNameStdout = `/dev/stdout`
	FileNameStderr = `/dev/stderr`
	FileModeDefault = 0640
	DirModeDefault = 0750
)

type MWriter interface {
	io.Writer
	AddFile(fn string) (error)
	AddWriter(io.Writer)
	Close()
}

type mWriter struct {
	out	io.Writer
	writers	[]io.Writer
	mu	sync.Mutex
}

func NewMWriter(stdout, stderr bool, files ...string) MWriter {

	var writers []io.Writer

	for _, fn := range files {
		if fh, err := createOrAppendFile(fn); err != nil {
			log.Println(err)
		} else {
			writers = append(writers, fh)
		}
	}
	if stdout {
		writers = append(writers, os.Stdout)
	}
	if stderr {
		writers = append(writers, os.Stderr)
	}

	this := &mWriter{writers: writers}
	this.reset()

	return this
}

func (this *mWriter) AddFile(fn string) (error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if fh, err := createOrAppendFile(fn); err != nil {
		return err
	} else {
		this.writers = append(this.writers, fh)
	}
	this.reset()
	return nil
}

func (this *mWriter) AddWriter(w io.Writer) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.writers = append(this.writers, w)
	this.reset()
}

func (this *mWriter) Write(b []byte) (int, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.out.Write(b)
}

func (this *mWriter) Close() {

	this.mu.Lock()
	defer this.mu.Unlock()

	for _, writer := range this.writers {

		if w, ok := writer.(*os.File); ok {
			w.Sync()
			if w.Name() == `/dev/stdout` || w.Name() == `/dev/stderr` {
				continue
			}
			w.Close()
		}
	}
}

func (this *mWriter) reset() {

	var writers []io.Writer

	for _, writer := range this.writers {

		if w, ok := writer.(*os.File); ok {
			w.Sync()
		}
		writers = append(writers, writer)
	}

	if len(writers) == 0 {
		writers = append(writers, ioutil.Discard)
	}

	this.out = io.MultiWriter(writers...)
}

func createOrAppendFile(fn string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(fn), DirModeDefault); err != nil {
		return nil, err
	}
	return os.OpenFile(fn, FileFlagsAppend, FileModeDefault)
}
