package main

import (
	"reflect"
	"testing"
)

func TestParseMessage(t *testing.T) {
	type args struct {
		msg []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Message
		wantErr bool
	}{
		{name: "Good message", args: args{[]byte("{\"id\": 1236, \"code\": \"200\", \"message\": \"shehh\"}")}, want: Message{ID: 1236, Code: "200", Message: "shehh"}, wantErr: false},
		{name: "Bad message", args: args{[]byte("{\"id\": 1236, \"code\": \"200\", \"sdffsfdfdf\": \"shehh\"}")}, want: Message{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMessage(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateHash(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"shplop", args{"shehhh"}, "4acd7e736c8b20955201397800f86b023c478341b670eec2dd8899d941831ec0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateHash(tt.args.text); got != tt.want {
				t.Errorf("GenerateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
