package utils

import (
	"reflect"
	"testing"
)

func TestParseKeyValuePairs(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name           string
		args           args
		expectedMap    map[string]string
		expectedErrors int
	}{
		{"valid single pair", args{"key=value"}, map[string]string{"key": "value"}, 0},
		{"valid single pair with comma", args{"key=value,"}, map[string]string{"key": "value"}, 0},
		{"valid multiple pairs", args{"key1=value1,key2=value2,key3=value3,key4=value4,"},
			map[string]string{"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
			}, 0},
		{"invalid empty key", args{"=value"}, map[string]string{}, 1},
		{"invalid empty value", args{"key="}, map[string]string{}, 1},
		{"valid whitespace", args{" key = value , key2 = value2 ,"}, map[string]string{"key": "value", "key2": "value2"}, 0},
		{"valid empty input", args{" ,"}, map[string]string{}, 0},
		{"invalid empty input", args{" ="}, map[string]string{}, 1},
		{"valid only commas", args{" ,,,,,, "}, map[string]string{}, 0},
		{"invalid only key", args{"key"}, map[string]string{}, 1},
		{"valid value has equal sign", args{"key=val1=val2"}, map[string]string{"key": "val1=val2"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, errors := ParseKeyValuePairs(tt.args.input)
			if !reflect.DeepEqual(got, tt.expectedMap) {
				t.Errorf("ParseKeyValuePairs() got = %v, want %v", got, tt.expectedMap)
			}
			if len(errors) != tt.expectedErrors {
				t.Errorf("ParseKeyValuePairs() got %d errors, want %d, errors: %v",
					len(errors), tt.expectedErrors, errors)
			}

		})
	}
}
