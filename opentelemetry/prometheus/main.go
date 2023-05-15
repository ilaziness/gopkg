// Copyright The OpenTelemetry Authors
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

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	prometheusExp "go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

var meter api.Meter
var provider *metric.MeterProvider

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx := context.Background()

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheusExp.New()
	if err != nil {
		log.Fatal(err)
	}
	provider = metric.NewMeterProvider(metric.WithReader(exporter))
	meter = provider.Meter("prometheus/example")

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics()

	opt := api.WithAttributes(
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	)

	// This is the equivalent of prometheus.NewCounterVec
	counter, err := meter.Float64Counter("foo", api.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	counter.Add(ctx, 5, opt)

	start := 1.
	counter2, err := meter.Float64ObservableCounter("foo2", api.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	meter.RegisterCallback(func(ctx context.Context, o api.Observer) error {
		start = start + 1.0
		o.ObserveFloat64(counter2, start)
		return nil
	}, counter2)

	gauge, err := meter.Float64ObservableGauge("bar", api.WithDescription("a fun little gauge"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
		n := -10. + rng.Float64()*(90.) // [-10, 100)
		o.ObserveFloat64(gauge, n, opt)
		return nil
	}, gauge)
	if err != nil {
		log.Fatal(err)
	}

	// This is the equivalent of prometheus.NewHistogramVec
	histogram, err := meter.Float64Histogram("baz", api.WithDescription("a very nice histogram"))
	if err != nil {
		log.Fatal(err)
	}
	histogram.Record(ctx, 23, opt)
	histogram.Record(ctx, 7, opt)
	histogram.Record(ctx, 101, opt)
	histogram.Record(ctx, 105, opt)

	gent()

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

func gent() {
	reg := prometheus.NewRegistry()

	cpuTemp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temp",
		Help: "current temp of the cpu",
	})
	reg.MustRegister(cpuTemp)

	cpuTemp.Set(65.3)
	//cpuTemp.

	http.Handle("/metrics2", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
}

func serveMetrics() {
	log.Printf("serving metrics at localhost:2223/metrics /metrics2")

	// 默认指标,默认包括go的版本，gc，内存，线程的指标信息
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":2223", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
