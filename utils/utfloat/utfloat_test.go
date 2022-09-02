package utfloat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	vals := map[float64]float64{
		187.319946452476573: 187.3199,
		0.00015678:          0.0002,
		0.00014891:          0.0001,
		0.00019999:          0.0002,
		0.00014447:          0.0001,
		0.00013912:          0.0001,
		0.0001:              0.0001,
		0.001:               0.001,
		0.01:                0.01,
		0.1:                 0.1,
		1:                   1,
	}

	for k, v := range vals {
		vx := Round(k, 4)
		assert.Equal(t, v, vx, fmt.Sprintf("from %v", k))
	}
}

func TestFloor(t *testing.T) {
	vals := map[float64]float64{
		0.00015678: 0.0001,
		0.00014891: 0.0001,
		0.00019999: 0.0001,
		0.00014447: 0.0001,
		0.00013912: 0.0001,
		0.00010599: 0.0001,
		0.00016999: 0.0001,
		0.00017999: 0.0001,
		0.00018999: 0.0001,
		0.0001:     0.0001,
		0.001:      0.001,
		0.01:       0.01,
		0.1:        0.1,
		1:          1,
	}

	for k, v := range vals {
		vx := Floor(k, 4)
		assert.Equal(t, v, vx, fmt.Sprintf("from %v", k))
	}
}

func TestCeil(t *testing.T) {
	vals := map[float64]float64{
		0.00015678: 0.0002,
		0.00014891: 0.0002,
		0.00019999: 0.0002,
		0.00014447: 0.0002,
		0.00013912: 0.0002,
		0.00010999: 0.0002,
		0.0001:     0.0001,
		0.001:      0.001,
		0.01:       0.01,
		0.1:        0.1,
		1:          1,
	}

	for k, v := range vals {
		vx := Ceil(k, 4)
		assert.Equal(t, v, vx, fmt.Sprintf("from %v", k))
	}
}
