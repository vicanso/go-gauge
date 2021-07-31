// Copyright 2021 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gauge

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Gauge struct {
	// unix nano time
	createdAt int64
	// sum of value
	sum int64
	// count of value
	count int64
	// reset count
	resetCount int64
	// reset cum
	resetSum int64
	// reset on fail
	resetOnFail bool
	// period of gauge
	period time.Duration
}

type Option func(m *Gauge)

// Returns a new gauge
// g := New(ResetCountOption(5))
func New(opts ...Option) *Gauge {
	g := &Gauge{}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Set the reset count of gauge
// g := New(ResetCountOption(5))
func ResetCountOption(resetCount int64) Option {
	return func(g *Gauge) {
		g.resetCount = resetCount
	}
}

// Set the reset sum of gauge
// m := New(ResetSumOption(10))
func ResetSumOption(sum int64) Option {
	return func(g *Gauge) {
		g.resetSum = sum
	}
}

// Set the value of reset on fail
// g := New(ResetOnFailOption())
func ResetOnFailOption() Option {
	return func(g *Gauge) {
		g.resetOnFail = true
	}
}

// Set period of guage
// g := New(PeriodOption(time.Minute))
func PeriodOption(period time.Duration) Option {
	return func(g *Gauge) {
		g.period = period
	}
}

// Resets the values of gauge
func (g *Gauge) Reset() {
	atomic.StoreInt64(&g.createdAt, 0)
	atomic.StoreInt64(&g.sum, 0)
	atomic.StoreInt64(&g.count, 0)
}

func (g *Gauge) before() {
	// 如果并行触发，只导致重复写入createAt
	// 并不会导致程序异常
	createdAt := atomic.LoadInt64(&g.createdAt)
	// 如果设置了周期参数，而已过一次周期，则重置
	if g.period != 0 && time.Now().UnixNano()-createdAt > int64(g.period) {
		g.Reset()
	}
	if atomic.LoadInt64(&g.createdAt) == 0 {
		atomic.StoreInt64(&g.createdAt, time.Now().UnixNano())
	}
}

// Adds value to gauge
func (g *Gauge) Add(value int64) (sum, count int64) {
	g.before()
	sum = atomic.AddInt64(&g.sum, value)
	count = atomic.AddInt64(&g.count, 1)
	if g.resetCount != 0 && count >= g.resetCount {
		g.Reset()
	} else if g.resetSum != 0 && sum > g.resetSum {
		g.Reset()
	}
	return
}

// Set max value to the gauge
func (g *Gauge) SetMax(value int64) (max, count int64) {
	g.before()
	max = atomic.LoadInt64(&g.sum)
	if value > max {
		atomic.StoreInt64(&g.sum, value)
		max = value
	}
	count = atomic.AddInt64(&g.count, 1)
	if g.resetCount != 0 && count >= g.resetCount {
		g.Reset()
	}
	return
}

// Adds value to gauge and check the mean value,
// if the mean value is greater than, it will return error
func (g *Gauge) AddCheckMean(value, max int64) (mean int64, err error) {
	g.Add(value)
	mean = g.Mean()
	if mean > max {
		if g.resetOnFail {
			g.Reset()
		}
		err = fmt.Errorf("mean is %d gt %d", mean, max)
		return
	}
	return
}

// Adds value to gauge and check the sum value,
// if the sum value is greater than, it will return error
func (g *Gauge) AddCheckSum(value, max int64) (sum int64, err error) {
	g.Add(value)
	sum = g.Sum()
	if sum > max {
		if g.resetOnFail {
			g.Reset()
		}
		err = fmt.Errorf("sum is %d gt %d", sum, max)
		return
	}
	return
}

// Mean returns the mean value of gauge
func (g *Gauge) Mean() int64 {
	// 两次操作存在时间差，有可能导致非同一次的数据
	// 由于仅是用于计算值，因此影响可接受
	count := atomic.LoadInt64(&g.count)
	if count == 0 {
		return 0
	}
	return atomic.LoadInt64(&g.sum) / count
}

// Count returns the count of gauge
func (g *Gauge) Count() int64 {
	return atomic.LoadInt64(&g.count)
}

// Sum returns the sum value of gauge
func (g *Gauge) Sum() int64 {
	return atomic.LoadInt64(&g.sum)
}
