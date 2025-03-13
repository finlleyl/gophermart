package luhn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckLuhn(t *testing.T) {
	// "79927398713" — валидный номер по алгоритму Луна
	valid := "79927398713"
	// Изменённый номер — не проходит проверку
	invalid := "79927398710"

	assert.True(t, CheckLuhn(valid), "Номер %s должен проходить проверку Луна", valid)
	assert.False(t, CheckLuhn(invalid), "Номер %s не должен проходить проверку Луна", invalid)
}
