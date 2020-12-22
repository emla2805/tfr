package utils

import (
	protobuf "github.com/emla2805/tfr/protobuf"
	"google.golang.org/protobuf/proto"
	"testing"
)

var (
	age = &protobuf.Feature{
		Kind: &protobuf.Feature_Int64List{
			Int64List: &protobuf.Int64List{Value: []int64{29}},
		},
	}
	movie = &protobuf.Feature{
		Kind: &protobuf.Feature_BytesList{
			BytesList: &protobuf.BytesList{Value: [][]byte{
				[]byte("The Shawshank Redemption"), []byte("Fight Club")}},
		},
	}
	movieRating = &protobuf.Feature{
		Kind: &protobuf.Feature_FloatList{
			FloatList: &protobuf.FloatList{Value: []float32{9.0, 9.7}},
		},
	}

	example = &protobuf.Example{
		Features: &protobuf.Features{
			Feature: map[string]*protobuf.Feature{
				"age":           age,
				"movie":         movie,
				"movie_ratings": movieRating,
			},
		},
	}

	movieNames   = &protobuf.FeatureList{Feature: []*protobuf.Feature{movie}}
	movieRatings = &protobuf.FeatureList{Feature: []*protobuf.Feature{movieRating}}
	movie1Actors = &protobuf.Feature{
		Kind: &protobuf.Feature_BytesList{
			BytesList: &protobuf.BytesList{Value: [][]byte{
				[]byte("Tim Robbins"), []byte("Morgan Freeman")}},
		},
	}
	movie2Actors = &protobuf.Feature{
		Kind: &protobuf.Feature_BytesList{
			BytesList: &protobuf.BytesList{Value: [][]byte{
				[]byte("Brad Pitt"), []byte("Edward Norton"), []byte("Helena Bonham Carter"),
			}},
		},
	}
	actors = &protobuf.FeatureList{Feature: []*protobuf.Feature{movie1Actors, movie2Actors}}

	sequenceExample = &protobuf.SequenceExample{
		Context: &protobuf.Features{
			Feature: map[string]*protobuf.Feature{"age": age},
		},
		FeatureLists: &protobuf.FeatureLists{
			FeatureList: map[string]*protobuf.FeatureList{
				"movie_names":   movieNames,
				"movie_ratings": movieRatings,
				"actors":        actors,
			},
		},
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

	sequenceExampleJSON = `{` +
		`"context":{` +
		`"feature":{` +
		`"age":{"int64List":{"value":[29]}}` +
		`}` +
		`},` +
		`"featureLists":{` +
		`"featureList":{` +
		`"actors":{` +
		`"feature":[` +
		`{"bytesList":{"value":["Tim Robbins","Morgan Freeman"]}},` +
		`{"bytesList":{"value":["Brad Pitt","Edward Norton","Helena Bonham Carter"]}}` +
		`]` +
		`},` +
		`"movie_names":{` +
		`"feature":[` +
		`{"bytesList":{"value":["The Shawshank Redemption","Fight Club"]}}` +
		`]` +
		`},` +
		`"movie_ratings":{` +
		`"feature":[` +
		`{"floatList":{"value":[9,9.7]}}` +
		`]` +
		`}` +
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
	{"sequenceExample object", sequenceExample, sequenceExampleJSON},
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
