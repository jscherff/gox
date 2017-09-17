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

type MultiLogger struct {
	*log.Logger
	writers	[]io.Writer
	files	[]*os.File
	stdout	bool
	stderr	bool
	buf		[]byte
	out		io.Writer
	mu		sync.Mutex
}

func NewMultiLogger(prefix string, flags int, stdout, stderr bool, files ...string) *MultiLogger {

	var f []*os.File
	
	for _, fn := range files {
		if fh, err := gox.MkdirOpen(fn); err != nil {
			log.Println(err)
		} else {
			f = append(f, fh)
		}
	}
	
	this := new(MultiLogger)
	
	this.mu.Lock()
	defer this.mu.Unlock()
	
	this.files = f
	this.stdout = stdout
	this.stderr = stderr
	this.SetPrefix(prefix)
	this.SetFlags(flags)

	return this
}

func (this *MultiLogger) AddFile(fn string) (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if fh, err := gox.MkdirOpen(fn); err == nil {
		this.files = append(this.files, fh)
		this.refreshWriters()
	}
	return err
}

func (this *MultiLogger) AddWriter(writer io.Writer) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.writers = append(this.writers, writer)
	this.refreshWriters()
}

func (this *MultiLogger) SetStdout(opt bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.stdout = opt
	this.refreshWriters()
}

func (this *MultiLogger) SetStderr(opt bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.stderr = opt
	this.refreshWriters()
}

func (this *MultiLogger) Write(b []byte) (n int, err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.buf = this.buf[:0]
	this.buf = append(this.buf, b...)
	if len(b) == 0 || b[len(b)-1] != '\n' {
		this.buf = append(this.buf, '\n')
	}
	return this.out.Write(this.buf)
}

func (this *MultiLogger) refreshWriters() {
	this.mu.Lock()
	defer this.mu.Unlock()
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
		writers = append(writers, f)
	}
	this.out = io.MultiWriter(writers...)
	this.SetOutput(this.out)
}