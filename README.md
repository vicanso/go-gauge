# go-gauge

Simple gauge for int, it support mean and sum.

```go
g := gauge.New(
    gauge.ResetCountOption(5),
    gauge.ResetSumOption(5),
    gauge.ResetOnFailOption(),
)
sum, count := g.Add(1)
mean, err := g.AddCheckMean(1, 5)
sum, err = g.AddCheckSum(2, 10)
```
