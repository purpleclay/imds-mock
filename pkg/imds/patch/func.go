/*
Copyright (c) 2022 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package patch

import (
	"fmt"
	"strings"
)

// JSONPairs transforms a map of key value pairs into a comma separated list
// of JSON field pairs.
//
//	pairs := map[string]string{
//		"FirstName": "joe",
//		"LastName": "bloggs",
//	}
//	out := JSONPairs(pairs)
//	fmt.Println(out)
//
//	"\"FirstName\": \"joe\", \"LastName\": \"bloggs\""
func JSONPairs(in map[string]string) string {
	pairs := make([]string, 0, len(in))
	for k, v := range in {
		pairs = append(pairs, fmt.Sprintf(`"%s": "%s"`, k, v))
	}

	return strings.Join(pairs, ", ")
}
