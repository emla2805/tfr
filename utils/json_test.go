package utils

import (
	protobuf "github.com/emla2805/tfr/protobuf"
	"google.golang.org/protobuf/proto"
	"testing"
)

var (
	fs = map[string]*protobuf.Feature{
		"age": {
			Kind: &protobuf.Feature_Int64List{
				Int64List: &protobuf.Int64List{Value: []int64{29}},
			},
		},
		"movie": {
			Kind: &protobuf.Feature_BytesList{
				BytesList: &protobuf.BytesList{Value: [][]byte{[]byte("The Shawshank Redemption"), []byte("Fight Club")}},
			},
		},
		"movie_ratings": {
			Kind: &protobuf.Feature_FloatList{
				FloatList: &protobuf.FloatList{Value: []float32{9.0, 9.7}},
			},
		},
	}
	example = &protobuf.Example{
		Features: &protobuf.Features{Feature: fs},
	}

	exampleJSON = `{` +
		`"features":{` +
		`"feature":{` +
		`"age":{"int64List":{"value":[29]}},` +
		`"movie":{"bytesList":{"value":["The Shawshank Redemption","Fight Club"]}},` +
		`"movie_ratings":{"floatList":{"value":[9,9.7]}}` +
		`}` +
		`}` +
		`}`
)
var marshalingTests = []struct {
	desc string
	pb   proto.Message
	json string
}{
	{"example object", example, exampleJSON},
}

func TestMarshaling(t *testing.T) {
	for _, tt := range marshalingTests {
		json, err := Marshal(tt.pb)
		jsonString := string(json)
		if err != nil {
			t.Errorf("%s: marshaling error: %v", tt.desc, err)
		} else if tt.json != jsonString {
			t.Errorf("%s:\ngot:  %v\nwant: %v", tt.desc, jsonString, tt.json)
		}
	}
}
