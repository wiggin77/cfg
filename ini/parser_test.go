package ini

import (
	"reflect"
	"testing"
)

func Test_buildLineArray(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"1", args{" line 1 \n line 2 \nline 3"}, []string{"line 1", "line 2", "line 3"}},
		{"2", args{" line 1 \r\n line 2 \r\nline 3"}, []string{"line 1", "line 2", "line 3"}},
		{"3", args{" line 1 \n ;line 2 \n#line 3"}, []string{"line 1"}},
		{"4", args{""}, []string{}},
		{"5", args{"\n\n\n  \n"}, []string{}},
		{"6", args{"#hello"}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildLineArray(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildLineArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSection(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1", args{"[sec1]"}, "sec1", false},
		{"2", args{"  []  "}, "", false},
		{"3", args{" [   ]  "}, "", false},
		{"4", args{"[  sec1  ]"}, "sec1", false},
		{"5", args{"[sec1"}, "", true},
		{"6", args{"sec1]"}, "", true},
		{"7", args{"blap"}, "", true},
		{"8", args{""}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSection(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseProp(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name        string
		args        args
		wantKey     string
		wantVal     string
		wantComment bool
		wantErr     bool
	}{
		{"1", args{" num = 77 "}, "num", "77", false, false},
		{"2", args{"#num = 77 "}, "", "", true, false},
		{"3", args{" num =   "}, "num", "", false, false},
		{"4", args{" blap! "}, "", "", false, true},
		{"5", args{"  = 77 "}, "", "", false, true},
		{"6", args{" num = 77 77"}, "num", "77 77", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotVal, gotComment, err := parseProp(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseProp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotComment != tt.wantComment {
				t.Errorf("parseProp() gotComment = %v, wantComment %v", gotComment, tt.wantComment)
			}
			if gotKey != tt.wantKey {
				t.Errorf("parseProp() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotVal != tt.wantVal {
				t.Errorf("parseProp() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}
