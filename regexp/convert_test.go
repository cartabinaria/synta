package regexp

import (
	"testing"

	"github.com/cartabinaria/synta"
	"github.com/stretchr/testify/assert"
)

func TestConvertBasic(t *testing.T) {
	content := `test = a|b
> test-test.test`
	basicSynta, err := synta.ParseSynta(content)
	assert.Nil(t, err)

	expr, err := Convert(basicSynta)
	assert.Nil(t, err)
	assert.Equal(t, "^(a|b)-(a|b)\\.(a|b)$", expr.String())
}

func TestConvertBasicOptional(t *testing.T) {
	content := `test = a|b
> test(-test)?.test`
	basicSynta, err := synta.ParseSynta(content)
	assert.Nil(t, err)

	expr, err := Convert(basicSynta)
	assert.Nil(t, err)
	assert.Equal(t, "^(a|b)(-(a|b))?\\.(a|b)$", expr.String())
}

func TestConvertMutiple(t *testing.T) {
	content := `test = a|b
castoro = roditore|anfibio
> test-castoro(-test)?.castoro`
	basicSynta, err := synta.ParseSynta(content)
	assert.Nil(t, err)

	expr, err := Convert(basicSynta)
	assert.Nil(t, err)
	assert.Equal(t, "^(a|b)-(roditore|anfibio)(-(a|b))?\\.(roditore|anfibio)$", expr.String())
}

func TestConvertExapleOnReadme(t *testing.T) {
	content := `; La tipologia della prova
tipo = scritto|orale
; Una data del tipo yyyy-mm-dd
data = \d{4}-\d{2}-\d{2}
; La fila e' un numero
fila = \d
; Una qualunque parola alfanumerica
extra = (\w|\d)+
; Estensione del file. Possibili valori:
; - txt, tex, md, pdf, doc, docx
ext = txt|tex|md|pdf|doc|docx

> tipo-data(-fila)?-extra.ext`
	basicSynta, err := synta.ParseSynta(content)
	assert.Nil(t, err)

	expr, err := Convert(basicSynta)
	assert.Nil(t, err)
	assert.Equal(t, "^(scritto|orale)-(\\d{4}-\\d{2}-\\d{2})(-(\\d))?-((\\w|\\d)+)\\.(txt|tex|md|pdf|doc|docx)$", expr.String())
}

func TestConvertMutipleNestedOptional(t *testing.T) {
	content := `test = a|b
castoro = roditore|anfibio
> test-castoro(-test(-castoro)?(-castoro)?)?.castoro`
	basicSynta, err := synta.ParseSynta(content)
	assert.Nil(t, err)

	expr, err := Convert(basicSynta)
	str := "^(a|b)-(roditore|anfibio)(-(a|b)(-(roditore|anfibio))?(-(roditore|anfibio))?)?\\.(roditore|anfibio)$"
	assert.Nil(t, err)
	assert.Equal(t, str, expr.String())

}
