package bootstrap

import (
	"math"
	"math/rand"
	"sort"
)

// Aggregator aggregates a slice of floats.
type Aggregator interface {
	Aggregate(values []float64) float64
}

// Resampler resamples floats with an Aggregator.
type Resampler interface {
	Resample(values []float64)
	Quantile(float64) float64
	Reset()
}

// SumAggregator generates a sum.
type SumAggregator struct{}

// NewSumAggregator returns a new SumAggregator.
func NewSumAggregator() SumAggregator {
	return SumAggregator{}
}

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

// NewAverageAggregator returns a new AverageAggregator.
func NewAverageAggregator() AverageAggregator {
	return AverageAggregator{}
}

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

// QuantileAggregator generates a quantile.
type QuantileAggregator struct {
	quantile float64
}

// NewQuantileAggregator returns a new QuantileAggregator with the given quantile.
func NewQuantileAggregator(quantile float64) QuantileAggregator {
	return QuantileAggregator{
		quantile: quantile,
	}
}

// Aggregate returns the a.quantile quantile of values.
func (a QuantileAggregator) Aggregate(values []float64) float64 {
	sort.Float64s(values)
	return quantile(values, a.quantile)
}

// BasicResampler is a basic bootstrap resampler.
type BasicResampler struct {
	aggregator       Aggregator
	iterations       int
	sampleAggregates []float64
	r                rand.Source
}

// NewBasicResampler returns a BasicResampler that aggregates values using aggregator.
func NewBasicResampler(aggregator Aggregator, iterations int) *BasicResampler {
	return &BasicResampler{
		aggregator:       aggregator,
		iterations:       iterations,
		sampleAggregates: make([]float64, 0, 100),
		r:                rand.NewSource(0),
	}
}

// Resample resamples from values for the given number of iterations and
// saves the aggregate values.
func (r *BasicResampler) Resample(values []float64) {
	length := len(values)
	scratch := make([]float64, length)
	for i := 0; i < r.iterations; i++ {
		for i := range values {
			scratch[i] = values[int(r.r.Int63())%length]
		}
		r.sampleAggregates = append(r.sampleAggregates, r.aggregator.Aggregate(scratch))
	}
	sort.Float64s(r.sampleAggregates)
}

// Quantile returns the q quantile of resampled aggregate values.
// Resample must be called before this method or NaN is returned.
func (r *BasicResampler) Quantile(q float64) float64 {
	return quantile(r.sampleAggregates, q)
}

// Reset resets any sampled state.
func (r *BasicResampler) Reset() {
	r.sampleAggregates = r.sampleAggregates[:0]
}

// PresampledResampler is a bootstrap resampler that precomputes sample indexes
// on creation.
type PresampledResampler struct {
	aggregator       Aggregator
	iterations       int
	sampleAggregates []float64
	samples          [][]int
}

// NewPresampledResampler returns a PresampledResampler that aggregates values using aggregator.
func NewPresampledResampler(aggregator Aggregator, iterations int, numValues int) *PresampledResampler {
	r := rand.NewSource(0)
	samples := make([][]int, iterations)
	for i := range samples {
		sampledInts := make([]int, numValues)
		for _ = range sampledInts {
			sampledInts[int(r.Int63())%numValues]++
		}
		samples[i] = sampledInts
	}
	return &PresampledResampler{
		aggregator:       aggregator,
		iterations:       iterations,
		sampleAggregates: make([]float64, 0, 100),
		samples:          samples,
	}
}

// Resample resamples from values for the given number of iterations and
// saves the aggregate values.
func (r *PresampledResampler) Resample(values []float64) {
	length := len(values)
	scratch := make([]float64, length)
	for i := 0; i < r.iterations; i++ {
		for j := range values {
			scratch[j] = values[j] * float64(r.samples[i][j])
		}
		r.sampleAggregates = append(r.sampleAggregates, r.aggregator.Aggregate(scratch))
	}
	sort.Float64s(r.sampleAggregates)
}

// Quantile returns the q quantile of resampled aggregate values.
// Resample must be called before this method or NaN is returned.
func (r *PresampledResampler) Quantile(q float64) float64 {
	return quantile(r.sampleAggregates, q)
}

// Reset resets any sampled state.
func (r *PresampledResampler) Reset() {
	r.sampleAggregates = r.sampleAggregates[:0]
}

// quantile returns the q quantile in vals.
func quantile(vals []float64, q float64) float64 {
	length := len(vals)
	if length == 0 {
		return math.NaN()
	}
	index := int(q * float64(length-1))
	return vals[index]
}
