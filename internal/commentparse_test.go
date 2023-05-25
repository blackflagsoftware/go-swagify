package internal

import (
	"reflect"
	"testing"
)

func TestParseSwagifyComment(t *testing.T) {
	tests := []struct {
		name     string
		comments []string
		want     Component
	}{
		{
			"successful",
			[]string{"/* go-swagify\n@@test: name1\n@@prop: prop_name\n@@\n@@again: name2\n@@another_prop: doh\n*/"},
			Component{Types: map[string]SwagifyComment{"test": {map[string][][]string{"name1": {{"prop: prop_name"}}}}, "again": {map[string][][]string{"name2": {{"another_prop: doh"}}}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSwagifyComment(tt.comments); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSwagifyComment() = %v, want %v", got, tt.want)
			}
		})
	}
}
