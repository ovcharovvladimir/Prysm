package tracer

import (
	"errors"

	"github.com/mattn/go-colorable"
	"github.com/ovcharovvladimir/Prysm/shared/p2p"
	"github.com/ovcharovvladimir/essentiaHybrid/log"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

// New creates and initializes a new tracing adapter.
func New(name, endpoint string, sampleFraction float64, enable bool) (p2p.Adapter, error) {
	// Set up the logger
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(3), log.StreamHandler(colorable.NewColorableStdout(), log.TerminalFormat(true))))
	if !enable {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.NeverSample()})
		return adapter, nil
	}

	if name == "" {
		return nil, errors.New("tracing service name cannot be empty")
	}

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(sampleFraction)})

	log.Info("Starting Jaeger exporter endpoint at address = ", endpoint)
	exporter, err := jaeger.NewExporter(jaeger.Options{
		Endpoint: endpoint,
		Process: jaeger.Process{
			ServiceName: name,
		},
	})
	if err != nil {
		return nil, err
	}
	trace.RegisterExporter(exporter)

	return adapter, nil
}

var adapter p2p.Adapter = func(next p2p.Handler) p2p.Handler {
	return func(msg p2p.Message) {
		var messageSpan *trace.Span
		msg.Ctx, messageSpan = trace.StartSpan(msg.Ctx, "handleP2pMessage")
		next(msg)
		messageSpan.End()
	}
}
