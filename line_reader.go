/*
 * Copyleft 2017, Simone Margaritelli <evilsocket at protonmail dot com>
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *   * Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *   * Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *   * Neither the name of ARM Inject nor the names of its contributors may be used
 *     to endorse or promote products derived from this software without
 *     specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */
package xray

import (
	"bufio"
	"os"
)

// LineReader will accept the name of a file and offset as argument
// and will return a channel from which lines can be read
// one at a time.
func LineReader(filename string, noff int64) (chan string, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// if offset defined then start from there
	if noff > 0 {
		// and go to the start of the line
		b := make([]byte, 1)
		for b[0] != '\n' {
			noff--
			fp.Seek(noff, os.SEEK_SET)
			fp.Read(b)
		}
		noff++
	}

	out := make(chan string)
	go func() {
		defer fp.Close()
		// we need to close the out channel in order
		// to signal the end-of-data condition
		defer close(out)
		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			noff, _ = fp.Seek(0, os.SEEK_CUR)
			out <- scanner.Text()
		}
	}()

	return out, nil
}
