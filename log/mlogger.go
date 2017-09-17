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
	`strings`
	`sync`
)

type MLogger struct {
	*log.Logger
	out	*MWriter
	buf	[]byte
	mu	sync.Mutex
}

func NewMLogger(prefix string, flags int, stdout, stderr bool, files ...string) *MLogger {

	mw := NewMWriter(stdout, stderr, files...)
	prefix = strings.TrimSpace(prefix) + ` `
	return &MLogger{Logger: log.New(mw, prefix, flags), out: mw}
}

func (this *MLogger) AddFile(f string) error {
	return this.out.AddFile(f)
}

func (this *MLogger) AddWriter(w io.Writer) {
	this.out.AddWriter(w)
}

func (this *MLogger) SetStdout(b bool) {
	this.out.SetStdout(b)
}

func (this *MLogger) SetStderr(b bool) {
	this.out.SetStderr(b)
}

func (this *MLogger) SetPrefix(p string) {
	this.Logger.SetPrefix(strings.TrimSpace(p) + ` `)
}

func (this *MLogger) Write(b []byte) (n int, err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.buf = this.buf[:0]
	this.buf = append(this.buf, b...)
	if len(b) == 0 || b[len(b)-1] != '\n' {
		this.buf = append(this.buf, '\n')
	}
	return this.out.Write(this.buf)
}

func (this *MLogger) Close() {
	this.out.Close()
}
