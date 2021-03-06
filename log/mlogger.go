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

const (
        LUTC =		log.LUTC
        Ldate =		log.Ldate
        Ltime =		log.Ltime
        Llongfile =	log.Llongfile
        Lshortfile =	log.Lshortfile
        LstdFlags =	log.LstdFlags
)

var loggerFlags = map[string]int {
	`utc`:		log.LUTC,
	`date`:		log.Ldate,
	`time`:		log.Ltime,
	`longfile`:	log.Llongfile,
	`shortfile`:	log.Lshortfile,
	`standard`:	log.LstdFlags,
}

type MLogger interface {

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
	Flags() int
	Output(int, string) error
	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
	Prefix() string
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
	SetFlags(int)
	SetOutput(io.Writer)
	SetPrefix(string)

	AddFile(string) (error)
	AddWriter(io.Writer)
	Write([]byte) (int, error)
	Close()
}

type mLogger struct {
	*log.Logger
	out	MWriter
	buf	[]byte
	mu	sync.Mutex
}

func NewMLogger(prefix string, flags int, stdout, stderr bool, files ...string) MLogger {

	mw := NewMWriter(stdout, stderr, files...)
	prefix = strings.TrimSpace(prefix) + ` `
	return &mLogger{Logger: log.New(mw, prefix, flags), out: mw}
}

func LoggerFlags(fs ...string) (lf int) {

	for _, f := range fs {
		lf |= loggerFlags[f]
	}

	return lf
}

func (this *mLogger) AddFile(f string) error {
	return this.out.AddFile(f)
}

func (this *mLogger) AddWriter(w io.Writer) {
	this.out.AddWriter(w)
}

func (this *mLogger) SetPrefix(p string) {
	this.Logger.SetPrefix(strings.TrimSpace(p) + ` `)
}

func (this *mLogger) Write(b []byte) (n int, err error) {

	this.mu.Lock()
	defer this.mu.Unlock()

	this.buf = this.buf[:0]
	this.buf = append(this.buf, b...)

	if len(b) == 0 || b[len(b)-1] != '\n' {
		this.buf = append(this.buf, '\n')
	}

	return this.out.Write(this.buf)
}

func (this *mLogger) Close() {
	this.out.Close()
}
