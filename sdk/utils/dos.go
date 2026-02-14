/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package utils

import (
	"os"
	"path"

	. "github.com/dimchat/mkm-go/format"
)

func MakeDirs(dir string) bool {
	err := os.MkdirAll(dir, os.ModePerm)
	if err == nil {
		return true
	}
	panic(err)
}

func PathJoin(elem ...string) string {
	return path.Join(elem...)
}
func PathDir(filepath string) string {
	return path.Dir(filepath)
}

func PathIsExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err) == false
}
func PathRemove(path string) bool {
	err := os.Remove(path)
	return err == nil
}

//
//  Binary File
//

func ReadBinaryFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err == nil {
		return data
	}
	return nil
}
func WriteBinaryFile(path string, data []byte) bool {
	err := os.WriteFile(path, data, 0644)
	return err == nil
}
func AppendBinaryFile(path string, data []byte) bool {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err == nil {
		cnt, err := fd.Write(data)
		return fd.Close() == nil && err == nil && cnt == len(data)
	}
	return false
}

//
//  Text File
//

func ReadTextFile(path string) string {
	data := ReadBinaryFile(path)
	if data == nil {
		return ""
	}
	return UTF8Decode(data)
}
func WriteTextFile(path string, text string) bool {
	return WriteBinaryFile(path, UTF8Encode(text))
}
func AppendTextFile(path string, text string) bool {
	return AppendBinaryFile(path, UTF8Encode(text))
}

//
//  JSON File
//

func ReadJSONFile(path string) interface{} {
	data := ReadBinaryFile(path)
	if data == nil {
		return nil
	}
	json := UTF8Decode(data)
	return JSONDecode(json)
}
func WriteJSONFile(path string, object interface{}) bool {
	json := JSONEncode(object)
	data := UTF8Encode(json)
	return WriteBinaryFile(path, data)
}
