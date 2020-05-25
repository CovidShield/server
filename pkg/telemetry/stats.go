package telemetry

import (
	"github.com/shirou/gopsutil/mem"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/unit"
)

// func initSystemStatsObserver() {
// 	meter := global.Meter("covidshield")

// 	var memTotal metric.Int64ValueObserver
// 	var memUsedPercent metric.Float64ValueObserver
// 	var memUsed metric.Int64ValueObserver
// 	var memAvailable metric.Int64ValueObserver

// 	cb := metric.Must(meter).NewBatchObserver(func(_ context.Context, result metric.BatchObserverResult) {
// 		v, _ := mem.VirtualMemory()
// 		result.Observe(nil,
// 			memTotal.Observation(int64(v.Total)),
// 			memUsedPercent.Observation(v.UsedPercent),
// 			memUsed.Observation(int64(v.Used)),
// 			memAvailable.Observation(int64(v.Available)),
// 		)
// 	})

// 	memTotal = cb.NewInt64ValueObserver("covidshield.system.memory.total",
// 		metric.WithDescription("Total amount of RAM on this system"),
// 		metric.WithUnit(unit.Bytes),
// 	)
// 	memUsedPercent = cb.NewFloat64ValueObserver("covidshield.system.memory.usedpercent",
// 		metric.WithDescription("RAM available for programs to allocate"),
// 	)
// 	memUsed = cb.NewInt64ValueObserver("covidshield.system.memory.used",
// 		metric.WithDescription("RAM used by programs"),
// 		metric.WithUnit(unit.Bytes),
// 	)
// 	memAvailable = cb.NewInt64ValueObserver("covidshield.system.memory.free",
// 		metric.WithDescription("Percentage of RAM used by programs"),
// 		metric.WithUnit(unit.Bytes),
// 	)
// }

func initSystemStatsRecorder() {
	meter := global.Meter("covidshield")
	v, _ := mem.VirtualMemory()

	memTotal := metric.Must(meter).NewInt64ValueRecorder("covidshield.system.memory.total",
		metric.WithDescription("Total amount of RAM on this system"),
		metric.WithUnit(unit.Bytes),
	)
	memUsedPercent := metric.Must(meter).NewFloat64ValueRecorder("covidshield.system.memory.usedpercent",
		metric.WithDescription("RAM available for programs to allocate"),
	)
	memUsed := metric.Must(meter).NewInt64ValueRecorder("covidshield.system.memory.used",
		metric.WithDescription("RAM used by programs"),
		metric.WithUnit(unit.Bytes),
	)
	memAvailable := metric.Must(meter).NewInt64ValueRecorder("covidshield.system.memory.free",
		metric.WithDescription("Percentage of RAM used by programs"),
		metric.WithUnit(unit.Bytes),
	)

	meter.RecordBatch(nil, nil,
		memTotal.Measurement(int64(v.Total)),
		memUsedPercent.Measurement(v.UsedPercent),
		memUsed.Measurement(int64(v.Used)),
		memAvailable.Measurement(int64(v.Available)))
}
