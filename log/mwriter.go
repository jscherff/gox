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
	`log`
	`os`
	`sync`
	`github.com/jscherff/gox`
)

type MWriter struct {
	out	io.Writer
	buf	[]byte
	writers	[]io.Writer
	files	[]*os.File
	stdout	bool
	stderr	bool
	mu	sync.Mutex
}

func NewMWriter(stdout, stderr bool, files ...string) *MWriter {

	var f []*os.File

	for _, fn := range files {
		if fh, err := gox.MkdirOpen(fn); err != nil {
			log.Println(err)
		} else {
			f = append(f, fh)
		}
	}

	this := new(MWriter)

	this.mu.Lock()
	defer this.mu.Unlock()

	this.files = f
	this.stdout = stdout
	this.stderr = stderr
	this.refresh()

	return this
}

func (this *MWriter) AddFile(f string) (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if fh, err := gox.MkdirOpen(f); err == nil {
		this.files = append(this.files, fh)
		this.refresh()
	}
	return err
}

func (this *MWriter) AddWriter(w io.Writer) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.writers = append(this.writers, w)
	this.refresh()
}

func (this *MWriter) SetStdout(b bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.stdout = b
	this.refresh()
}

func (this *MWriter) SetStderr(b bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.stderr = b
	this.refresh()
}

func (this *MWriter) Write(b []byte) (int, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.out.Write(b)
}

func (this *MWriter) Close() {
	this.mu.Lock()
	defer this.mu.Unlock()
	for _, f := range this.files {
		f.Sync()
		f.Close()
	}
}

func (this *MWriter) refresh() {
	var writers []io.Writer
	if this.stdout {
		writers = append(writers, os.Stdout)
	}
	if this.stderr {
		writers = append(writers, os.Stderr)
	}
	for _, w := range this.writers {
		writers = append(writers, w)
	}
	for _, f := range this.files {
		f.Sync()
		writers = append(writers, f)
	}
	this.out = io.MultiWriter(writers...)
}
