package bootstrap

import (
	"testing"
)

// TestSum tests SumAggregator.
func TestSum(t *testing.T) {
	sum := SumAggregator{}
	aggregate := sum.Aggregate([]float64{0, 1, 2, 3, 4})
	if aggregate != 10.0 {
		t.Errorf("expected aggregate %f; got %f", 10.0, aggregate)
	}
}

// TestAverage tests AverageAggregator.
func TestAverage(t *testing.T) {
	avg := AverageAggregator{}
	aggregate := avg.Aggregate([]float64{0, 1, 2, 3, 4})
	if aggregate != 2.0 {
		t.Errorf("expected aggregate %f; got %f", 2.0, aggregate)
	}
}

// TestResampler tests the resampler.
func TestResampler(t *testing.T) {
	resampler := NewResampler(SumAggregator{})
	resampler.r.Seed(0)
	resampler.Resample([]float64{0, 1, 2, 3, 4}, 2000)
	if min := resampler.Quantile(0); min != 0.0 {
		t.Errorf("expected min to be %f; got %f", 0.0, min)
	}
	if median := resampler.Quantile(0.5); median != 10.0 {
		t.Errorf("expected median to be %f; got %f", 10.0, median)
	}
	if max := resampler.Quantile(1); max != 20.0 {
		t.Errorf("expected max to be %f; got %f", 20.0, max)
	}
}
