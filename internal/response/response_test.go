package response

import (
	"reflect"
	"testing"

	in "github.com/blackflagsoftware/go-swagify/internal"
)

func TestBuildResponse(t *testing.T) {
	type args struct {
		comments in.SwagifyComment
	}
	tests := []struct {
		name string
		args args
		want map[string]Response
	}{
		{
			"successful: one response (200) non-ref",
			args{comments: in.SwagifyComment{Comments: map[string][]string{"200": {
				"desc: This is my description",
				"content_name: application/json",
				"content_ref: response_1",
			}}}},
			map[string]Response{"200": {Description: "This is my description", Content: map[string]Content{"application/json": {RefSchema{Ref: "#/components/responses/response_1"}}}}},
		},
		{
			"successful: one response (200) ref",
			args{comments: in.SwagifyComment{Comments: map[string][]string{"200": {
				"ref: response_ref_1",
				"desc: This is my description",
				"content_name: application/json",
				"content_ref: response_1",
			}}}},
			map[string]Response{"200": {Ref: "#/components/responses/response_ref_1", Content: map[string]Content{}}},
		},
		{
			"successful: multiple responses (200), non-ref",
			args{comments: in.SwagifyComment{Comments: map[string][]string{"200": {
				"desc: This is my description",
				"content_name: application/json",
				"content_ref: response_1",
				"content_name: application/text",
				"content_ref: response_2",
			}}}},
			map[string]Response{"200": {Description: "This is my description", Content: map[string]Content{"application/json": {RefSchema{Ref: "#/components/responses/response_1"}}, "application/text": {RefSchema{Ref: "#/components/responses/response_2"}}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildResponse(tt.args.comments); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOperationResponseLines(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want map[string]Response
	}{
		{
			"successful: one",
			args{[]string{
				"resp_name: 400",
				"resp_ref: SomeErrorResponse",
			}},
			map[string]Response{"400": {Ref: "#/components/responses/SomeErrorResponse", Content: map[string]Content{}}},
		},
		{
			"successful: multiple",
			args{[]string{
				"resp_name: 400",
				"resp_ref: SomeErrorResponse",
				"resp_name: 500",
				"resp_ref: SomeServerErrorResponse",
			}},
			map[string]Response{"400": {Ref: "#/components/responses/SomeErrorResponse", Content: map[string]Content{}}, "500": {Ref: "#/components/responses/SomeServerErrorResponse", Content: map[string]Content{}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseOperationResponseLines(tt.args.lines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseOperationResponseLines() = %v, want %v", got, tt.want)
			}
		})
	}
}
