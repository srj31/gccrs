// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zip

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

type ZipTest struct {
	Name    string
	Comment string
	File    []ZipTestFile
	Error   os.Error // the error that Opening this file should return
}

type ZipTestFile struct {
	Name    string
	Content []byte // if blank, will attempt to compare against File
	File    string // name of file to compare to (relative to testdata/)
	Mtime   string // modified time in format "mm-dd-yy hh:mm:ss"
}

// Caution: The Mtime values found for the test files should correspond to
//          the values listed with unzip -l <zipfile>. However, the values
//          listed by unzip appear to be off by some hours. When creating
//          fresh test files and testing them, this issue is not present.
//          The test files were created in Sydney, so there might be a time
//          zone issue. The time zone information does have to be encoded
//          somewhere, because otherwise unzip -l could not provide a different
//          time from what the archive/zip package provides, but there appears
//          to be no documentation about this.

var tests = []ZipTest{
	{
		Name:    "test.zip",
		Comment: "This is a zipfile comment.",
		File: []ZipTestFile{
			{
				Name:    "test.txt",
				Content: []byte("This is a test text file.\n"),
				Mtime:   "09-05-10 12:12:02",
			},
			{
				Name:  "gophercolor16x16.png",
				File:  "gophercolor16x16.png",
				Mtime: "09-05-10 15:52:58",
			},
		},
	},
	{
		Name: "r.zip",
		File: []ZipTestFile{
			{
				Name:  "r/r.zip",
				File:  "r.zip",
				Mtime: "03-04-10 00:24:16",
			},
		},
	},
	{Name: "readme.zip"},
	{Name: "readme.notzip", Error: FormatError},
	{
		Name: "dd.zip",
		File: []ZipTestFile{
			{
				Name:    "filename",
				Content: []byte("This is a test textfile.\n"),
				Mtime:   "02-02-11 13:06:20",
			},
		},
	},
}

func TestReader(t *testing.T) {
	for _, zt := range tests {
		readTestZip(t, zt)
	}
}

func readTestZip(t *testing.T, zt ZipTest) {
	z, err := OpenReader("testdata/" + zt.Name)
	if err != zt.Error {
		t.Errorf("error=%v, want %v", err, zt.Error)
		return
	}

	// bail if file is not zip
	if err == FormatError {
		return
	}
	defer z.Close()

	// bail here if no Files expected to be tested
	// (there may actually be files in the zip, but we don't care)
	if zt.File == nil {
		return
	}

	if z.Comment != zt.Comment {
		t.Errorf("%s: comment=%q, want %q", zt.Name, z.Comment, zt.Comment)
	}
	if len(z.File) != len(zt.File) {
		t.Errorf("%s: file count=%d, want %d", zt.Name, len(z.File), len(zt.File))
	}

	// test read of each file
	for i, ft := range zt.File {
		readTestFile(t, ft, z.File[i])
	}

	// test simultaneous reads
	n := 0
	done := make(chan bool)
	for i := 0; i < 5; i++ {
		for j, ft := range zt.File {
			go func() {
				readTestFile(t, ft, z.File[j])
				done <- true
			}()
			n++
		}
	}
	for ; n > 0; n-- {
		<-done
	}

	// test invalid checksum
	if !z.File[0].hasDataDescriptor() { // skip test when crc32 in dd
		z.File[0].CRC32++ // invalidate
		r, err := z.File[0].Open()
		if err != nil {
			t.Error(err)
			return
		}
		var b bytes.Buffer
		_, err = io.Copy(&b, r)
		if err != ChecksumError {
			t.Errorf("%s: copy error=%v, want %v", z.File[0].Name, err, ChecksumError)
		}
	}
}

func readTestFile(t *testing.T, ft ZipTestFile, f *File) {
	if f.Name != ft.Name {
		t.Errorf("name=%q, want %q", f.Name, ft.Name)
	}

	mtime, err := time.Parse("01-02-06 15:04:05", ft.Mtime)
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := f.Mtime_ns()/1e9, mtime.Seconds(); got != want {
		t.Errorf("%s: mtime=%s (%d); want %s (%d)", f.Name, time.SecondsToUTC(got), got, mtime, want)
	}

	size0 := f.UncompressedSize

	var b bytes.Buffer
	r, err := f.Open()
	if err != nil {
		t.Error(err)
		return
	}

	if size1 := f.UncompressedSize; size0 != size1 {
		t.Errorf("file %q changed f.UncompressedSize from %d to %d", f.Name, size0, size1)
	}

	_, err = io.Copy(&b, r)
	if err != nil {
		t.Error(err)
		return
	}
	r.Close()

	var c []byte
	if len(ft.Content) != 0 {
		c = ft.Content
	} else if c, err = ioutil.ReadFile("testdata/" + ft.File); err != nil {
		t.Error(err)
		return
	}

	if b.Len() != len(c) {
		t.Errorf("%s: len=%d, want %d", f.Name, b.Len(), len(c))
		return
	}

	for i, b := range b.Bytes() {
		if b != c[i] {
			t.Errorf("%s: content[%d]=%q want %q", f.Name, i, b, c[i])
			return
		}
	}
}

func TestInvalidFiles(t *testing.T) {
	const size = 1024 * 70 // 70kb
	b := make([]byte, size)

	// zeroes
	_, err := NewReader(sliceReaderAt(b), size)
	if err != FormatError {
		t.Errorf("zeroes: error=%v, want %v", err, FormatError)
	}

	// repeated directoryEndSignatures
	sig := make([]byte, 4)
	binary.LittleEndian.PutUint32(sig, directoryEndSignature)
	for i := 0; i < size-4; i += 4 {
		copy(b[i:i+4], sig)
	}
	_, err = NewReader(sliceReaderAt(b), size)
	if err != FormatError {
		t.Errorf("sigs: error=%v, want %v", err, FormatError)
	}
}

type sliceReaderAt []byte

func (r sliceReaderAt) ReadAt(b []byte, off int64) (int, os.Error) {
	copy(b, r[int(off):int(off)+len(b)])
	return len(b), nil
}
