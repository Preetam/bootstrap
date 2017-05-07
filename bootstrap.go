package bootstrap

import (
	"math/rand"
	"sort"
)

// An aggregator aggregates a slice of floats.
type Aggregator interface {
	Aggregate(values []float64) float64
}

// SumAggregator generates a sum.
type SumAggregator struct{}

// Aggregate returns the sum of values.
func (a SumAggregator) Aggregate(values []float64) float64 {
	result := 0.0
	for _, v := range values {
		result += v
	}
	return result
}

// AverageAggregator generates an average.
type AverageAggregator struct{}

// Aggregate returns the average of values.
func (a AverageAggregator) Aggregate(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}

	result := 0.0
	for _, v := range values {
		result += v
	}
	return result / float64(len(values))
}

// Resampler is a bootstrap resampler.
type Resampler struct {
	aggregator       Aggregator
	sampleAggregates []float64
	r                rand.Source
}

// NewResampler returns a Resampler that aggregates values using aggregator.
func NewResampler(aggregator Aggregator) *Resampler {
	return &Resampler{
		aggregator:       aggregator,
		sampleAggregates: make([]float64, 0, 100),
		r:                rand.NewSource(0),
	}
}

// Resample resamples from values for the given number of iterations and
// saves the aggregate values.
func (r *Resampler) Resample(values []float64, iterations int) {
	length := len(values)
	scratch := make([]float64, length)
	for i := 0; i < iterations; i++ {
		for i := range values {
			scratch[i] = values[int(r.r.Int63())%length]
		}
		r.sampleAggregates = append(r.sampleAggregates, r.aggregator.Aggregate(scratch))
	}
	sort.Float64s(r.sampleAggregates)
}

// Quantile returns the q quantile of resampled aggregate values.
// Resample must be called before this method.
func (r *Resampler) Quantile(q float64) float64 {
	length := len(r.sampleAggregates)
	index := int(q * float64(length-1))
	return r.sampleAggregates[index]
}
