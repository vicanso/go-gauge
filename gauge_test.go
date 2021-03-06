package gauge

import (
	"testing"
	"time"

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

	// sum超过后重置再处理
	sum, count = g.Add(2)
	assert.Equal(int64(2), sum)
	assert.Equal(int64(1), count)
	g.Reset()

	for i := 0; i < 5; i++ {
		_, count := g.Add(1)
		assert.Equal(int64(i+1), count)
	}
	// 次数超过后，重置处理
	_, _ = g.Add(2)
	assert.Equal(int64(2), g.Sum())
	assert.Equal(int64(1), g.Count())
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

func TestGaugeSetMax(t *testing.T) {
	assert := assert.New(t)

	g := New(PeriodOption(time.Millisecond))
	max, count := g.SetMax(10)
	assert.Equal(int64(10), max)
	assert.Equal(int64(1), count)

	max, count = g.SetMax(1)
	assert.Equal(int64(10), max)
	assert.Equal(int64(2), count)

	max, count = g.SetMax(11)
	assert.Equal(int64(11), max)
	assert.Equal(int64(3), count)

	// 别一个区间，重置
	time.Sleep(2 * time.Millisecond)
	max, count = g.SetMax(10)
	assert.Equal(int64(10), max)
	assert.Equal(int64(1), count)
}
