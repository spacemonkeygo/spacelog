// Copyright (C) 2016 Space Monkey, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spacelog

import (
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

func TestFileWriterOutput(t *testing.T) {
	stack := func() string {
		var buf [4096]byte
		return string(buf[:runtime.Stack(buf[:], false)])
	}
	assertNoError := func(err error) {
		if err != nil {
			t.Fatalf("error: %v\n%s", err, stack())
		}
	}
	assertContents := func(path, contents string) {
		data, err := ioutil.ReadFile(path)
		assertNoError(err)
		if string(data) != contents {
			t.Fatalf("%q != %q\n%s", data, contents, stack())
		}
	}

	fh, err := ioutil.TempFile("", "spacelog-")
	assertNoError(err)
	assertNoError(fh.Close())

	name := fh.Name()
	rotated_name := name + ".1"
	defer os.Remove(name)
	defer os.Remove(rotated_name)

	fwo, err := NewFileWriterOutput(fh.Name())
	assertNoError(err)

	fwo.Output(Critical, []byte("hello world"))
	assertContents(name, "hello world\n")

	fwo.OnHup()

	assertNoError(os.Rename(name, rotated_name))
	assertContents(rotated_name, "hello world\n")

	fwo.Output(Critical, []byte("hello universe"))
	assertContents(name, "hello universe\n")
	assertContents(rotated_name, "hello world\n")
}
