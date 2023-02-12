package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKVCLI(t *testing.T) {
	t.Parallel()
	type args struct {
		fileName string
	}
	tests := []struct {
		name            string
		args            args
		wantLogContents []string
	}{
		{
			name: "tc1: root README testcase",
			args: args{fileName: "tc1"},
			wantLogContents: []string{"hello", "hello", "hello-again", "Key not found: a",
				"Key not found: a", "once-more", "hello", "Exiting..."},
		},
		{
			name: "tc2: case insensitive commands",
			args: args{fileName: "tc2"},
			wantLogContents: []string{"hello", "hello", "hello-again", "Key not found: a",
				"Key not found: a", "once-more", "hello", "Exiting..."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			in := &bytes.Buffer{}
			buf := &bytes.Buffer{}

			f, err := os.Open(fmt.Sprintf("testdata/%s.txt", tt.args.fileName))
			require.NoError(t, err)
			fileContents, err := ioutil.ReadAll(f)
			require.NoError(t, err)
			in.Write(fileContents)

			KVCLI(in, buf)

			logOut := bufio.NewScanner(buf)
			gotOut := make([]string, 0)
			for logOut.Scan() {
				gotOut = append(gotOut, logOut.Text())
			}
			assert.Equal(t, tt.wantLogContents, gotOut)
		})
	}
}
