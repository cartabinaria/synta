package json

import (
	"testing"

	"github.com/cartabinaria/synta"
	"github.com/stretchr/testify/assert"
)

type StringDefintions = map[string]StringDefinition
type StringDefinition struct {
	Regexp   string
	Comments []string
}

func checkSegments(t *testing.T, synSeg []synta.Segment, expSeg []Segment) {
	for i, e := range synSeg {
		assert.Equal(t, uint(e.Kind), uint(expSeg[i].Kind))
		switch e.Kind {
		case synta.SegmentTypeIdentifier:
			assert.Equal(t, string(*e.Value), string(expSeg[i].Value))
		case synta.SegmentTypeOptional:
			if e.Value == nil {
				assert.True(t, expSeg[i].Value == "")
			} else {
				assert.True(t, string(expSeg[i].Value) == string(*e.Value))
			}
			checkSegments(t, e.Subsegments, expSeg[i].Subsegments)
		}
	}
}

func checkSynta(t *testing.T, syn synta.Synta, expected Synta) {
	assert.Equal(t, len(syn.Definitions), len(expected.Definitions))

	for id, def := range syn.Definitions {
		if len(def.Comments) == 0 {
			assert.Equal(t, len(def.Comments), len(expected.Definitions[string(id)].Comments))
		} else {
			assert.Equal(t, def.Comments, expected.Definitions[string(id)].Comments)
		}
	}
	checkSegments(t, syn.Filename.Segments, expected.Filename.Segments)
	assert.Equal(t, string(syn.Filename.Extension), string(expected.Filename.Extension))
}

func TestConvertBasic(t *testing.T) {
	input := `test = a|b
> test-test.test`
	syn, err := synta.ParseSynta(input)
	assert.Nil(t, err)

	expectedConvert := Convert(syn)
	assert.Nil(t, err)
	checkSynta(t, syn, expectedConvert)
}

func TestConvertWithOptional(t *testing.T) {
	input := `test = a|b
> test(-test)?.test`
	syn, err := synta.ParseSynta(input)
	assert.Nil(t, err)

	expectedConvert := Convert(syn)
	assert.Nil(t, err)
	checkSynta(t, syn, expectedConvert)
}

func TestConvertWithNestedOptional(t *testing.T) {
	input := `atest = a|b
btest = a|b
ctest = a|b
etest = a|b
> atest(-btest(-ctest)?(-ctest)?)?.etest`
	syn, err := synta.ParseSynta(input)
	assert.Nil(t, err)

	expectedConvert := Convert(syn)
	assert.Nil(t, err)
	checkSynta(t, syn, expectedConvert)
}
