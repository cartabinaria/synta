package format

import (
	"testing"

	"github.com/csunibo/synta"
	"github.com/stretchr/testify/assert"
)

func TestFormatBasic(t *testing.T) {
	basicContent := `def = a|b
test = c|d
> def.test
`
	basicSynta, err := synta.ParseSynta(basicContent)
	assert.Nil(t, err)

	formatted := Format(basicSynta)
	formattedContent := `def = a|b

test = c|d

> def.test
`
	assert.Equal(t, formattedContent, formatted)
}

func TestFormatWithOptional(t *testing.T) {
	basicContent := `def = a|b
test = c|d
> def(-test)?.test
`
	basicSynta, err := synta.ParseSynta(basicContent)
	assert.Nil(t, err)

	formatted := Format(basicSynta)
	formattedContent := `def = a|b

test = c|d

> def(-test)?.test
`
	assert.Equal(t, formattedContent, formatted)
}

func TestFormatWithComments(t *testing.T) {
	basicContent := `; a test comment
def = a|b
test = c|d
> def(-test)?.test
`
	basicSynta, err := synta.ParseSynta(basicContent)
	assert.Nil(t, err)

	formatted := Format(basicSynta)
	formattedContent := `; a test comment
def = a|b

test = c|d

> def(-test)?.test
`
	assert.Equal(t, formattedContent, formatted)
}
