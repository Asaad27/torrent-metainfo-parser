package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_parseArgs(t *testing.T) {
	tests := []struct {
		name      string
		arg       string
		Args      []string
		want      Arguments
		wantError bool
	}{
		{name: "test1", arg: "file.go fileName -c comment -a announceUrl -b Asaad -p 100", want: Arguments{Comment: 3, Announce: 5, CreatedBy: 7, SizeOfPiece: 9}, wantError: false},
		{name: "test2", arg: "file.go fileName -c comment -a announceUrl", want: Arguments{Comment: 3, Announce: 5, CreatedBy: -1, SizeOfPiece: -1}, wantError: false},
		{name: "test3", arg: "file.go fileName -c comment -a announceUrl -b Asaad -p 100", want: Arguments{Comment: 3, Announce: 5, CreatedBy: 7, SizeOfPiece: 9}, wantError: false},
		{name: "test4", arg: "file.go fileName -c comment -a announceUrl -b Asaad -p 100", want: Arguments{Comment: 3, Announce: 5, CreatedBy: 7, SizeOfPiece: 9}, wantError: false},
		{name: "test5", arg: "file.go fileName -c comment -a announceUrl -bf Asaad -p 100", want: *NewArguments(), wantError: true},
		{name: "test6", arg: "file.go fileName -c comment -a other other", want: *NewArguments(), wantError: true},
		{name: "test7", arg: "file.go fileName -c comment -a", want: *NewArguments(), wantError: true},
		{name: "test8", arg: "file.go fileName", want: *NewArguments(), wantError: false},
		{name: "test9", arg: "file.go fileName test", want: *NewArguments(), wantError: true},
		{name: "test10", arg: "file.go fileName -t weird", want: *NewArguments(), wantError: true},
		{name: "test4", arg: "file.go fileName -a ann -c comm -p 27 -b as", want: Arguments{Comment: 5, Announce: 3, CreatedBy: 9, SizeOfPiece: 7}, wantError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.Args = strings.Split(tt.arg, " ")
			err, got := ParseArgs(tt.Args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArgs() = %v, want %v", got, tt.want)
			}
			wantError := err != nil
			if wantError != tt.wantError {
				t.Errorf("erreurs ? parseArgs() = %v, want %v", wantError, tt.wantError)
			}
		})
	}
}
