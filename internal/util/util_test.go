package util

import "testing"

func TestBuildAlternateFieldName(t *testing.T) {
	type args struct {
		name string
		mode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"SnakeCase",
			args{
				"HelloWorld",
				"snakeCase",
			},
			"hello_world",
		},
		{
			"KebabCase",
			args{
				"HelloWorld",
				"kebabCase",
			},
			"hello-world",
		},
		{
			"PascalCase",
			args{
				"hello_world",
				"pascalCase",
			},
			"HelloWorld",
		},
		{
			"CamelCase",
			args{
				"Hello_World",
				"camelCase",
			},
			"helloWorld",
		},
		{
			"LowerCase",
			args{
				"Hello_World",
				"lowerCase",
			},
			"hello_world",
		},
		{
			"KebabCase 2",
			args{
				"hello_world",
				"kebabCase",
			},
			"hello-world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildAlternateFieldName(tt.args.name, tt.args.mode); got != tt.want {
				t.Errorf("BuildAlternateFieldName() = %v, want %v", got, tt.want)
			}
		})
	}
}
