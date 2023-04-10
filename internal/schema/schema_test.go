package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseTag(t *testing.T) {
	type args struct {
		fieldName string
		fieldType string
		tagValue  string
		schemas   map[string]Schema
	}
	tests := []struct {
		name          string
		args          args
		wantKey       string
		isRequiredKey string
	}{
		{
			"successful",
			args{
				fieldName: "Address",
				fieldType: "string",
				tagValue:  "sw:\"AddressRequest*;AddressResponse\" sw_ex:\"some example\" sw_desc:\"some desc\"",
				schemas:   make(map[string]Schema),
			},
			"AddressRequest",
			"AddressRequest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseTag(tt.args.fieldName, tt.args.fieldType, tt.args.tagValue, tt.args.schemas)
			_, ok := tt.args.schemas[tt.wantKey]
			assert.Equal(t, true, ok, "No key")
			for v, k := range tt.args.schemas {
				if len(k.Required) > 0 {
					assert.Equal(t, tt.isRequiredKey, v, "Required")
				}
			}
		})
	}
}

func Test_determineRequired(t *testing.T) {
	type args struct {
		schemaName string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			"required",
			args{schemaName: "TestHere*"},
			"TestHere",
			true,
		},
		{
			"not required",
			args{schemaName: "TestHere"},
			"TestHere",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := determineRequired(tt.args.schemaName)
			assert.Equal(t, tt.want, got, "name not equal")
			assert.Equal(t, tt.want1, got1, "bool value not equal")
		})
	}
}
