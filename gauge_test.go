package gauge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGauge(t *testing.T) {
	assert := assert.New(t)

	g := New(
		ResetCountOption(5),
		ResetSumOption(10),
		ResetOnFailOption(),
	)

	_, _ = g.Add(1)
	_, _ = g.Add(2)

	assert.Equal(int64(3), g.Sum())
	assert.Equal(int64(2), g.Count())
	assert.Equal(int64(1), g.Mean())

	sum, count := g.Add(10)
	assert.Equal(int64(13), sum)
	assert.Equal(int64(3), count)
	// 重置
	assert.Equal(int64(0), g.Sum())
	assert.Equal(int64(0), g.Count())

	for i := 0; i < 5; i++ {
		_, count := g.Add(1)
		assert.Equal(int64(i+1), count)
	}
	// 重置
	assert.Equal(int64(0), g.Sum())
	assert.Equal(int64(0), g.Count())
}

func TestGaugeAddAndCheck(t *testing.T) {
	assert := assert.New(t)

	g := New(ResetOnFailOption())
	assert.Equal(int64(0), g.Mean())

	_, _ = g.Add(10)
	mean, err := g.AddCheckMean(2, 5)
	assert.Equal(int64(6), mean)
	assert.Equal("mean is 6 gt 5", err.Error())

	assert.Equal(int64(0), g.Count())
	assert.Equal(int64(0), g.Sum())

	_, _ = g.Add(10)
	sum, err := g.AddCheckSum(2, 10)
	assert.Equal(int64(12), sum)
	assert.Equal("sum is 12 gt 10", err.Error())
}
