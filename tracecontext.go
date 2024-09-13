package tracecontext

import (
	"encoding/hex"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

const (
	// TraceparentVersion is the version of the traceparent header.
	traceparentVersion = "00"

	// TraceparentHTTPHeaderTag is the HTTP header tag for traceparent.
	TraceparentHTTPHeaderTag = "traceparent"

	// TraceparentHTTPHeaderTag is the HTTP header tag for traceparent.
	TracestateHTTPHeaderTag = "tracestate"

	// traceparentParts is the number of parts in a traceparent header.
	traceparentParts = 4
)

var (
	// errTraceparentInvalidFormat is returned when the traceparent format is invalid.
	errTraceparentInvalidFormat = errors.New("invalid traceparent format")
	// errTraceparentInvalidVersion is returned when the traceparent version is invalid.
	errTraceparentInvalidVersion = errors.New("invalid traceparent version")
)

func MarshalSpanContext(sc trace.SpanContext) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		traceparentVersion, sc.TraceID().String(), sc.SpanID().String(), sc.TraceFlags().String())
}

func UnmarshalSpanContext(traceparent, tracestate string) (trace.SpanContextConfig, error) {
	var version, traceID, parentID, flags string

	if n, err := fmt.Sscanf(traceparent, "%2s-%32s-%16s-%2s", &version, &traceID, &parentID, &flags); err != nil {
		return trace.SpanContextConfig{}, fmt.Errorf("failed to parse traceparent: %w", err)
	} else if n != traceparentParts {
		return trace.SpanContextConfig{}, fmt.Errorf("%w: %s", errTraceparentInvalidFormat, traceparent)
	}

	if version != traceparentVersion {
		return trace.SpanContextConfig{}, fmt.Errorf("%w: %s", errTraceparentInvalidVersion, version)
	}

	var cfgTraceID, cfgSpanID, cgfTraceFlags []byte

	var cfgTraceState trace.TraceState

	var err error

	if cfgTraceID, err = hex.DecodeString(traceID); err != nil {
		return trace.SpanContextConfig{}, fmt.Errorf("failed to decode trace ID: %w", err)
	}

	if cfgSpanID, err = hex.DecodeString(parentID); err != nil {
		return trace.SpanContextConfig{}, fmt.Errorf("failed to decode parent ID: %w", err)
	}

	if cgfTraceFlags, err = hex.DecodeString(flags); err != nil {
		return trace.SpanContextConfig{}, fmt.Errorf("failed to decode flags: %w", err)
	}

	if cfgTraceState, err = trace.ParseTraceState(tracestate); err != nil {
		return trace.SpanContextConfig{}, fmt.Errorf("failed to parse tracestate: %w", err)
	}

	return trace.SpanContextConfig{
		TraceID:    trace.TraceID(cfgTraceID),
		SpanID:     trace.SpanID(cfgSpanID),
		TraceFlags: trace.TraceFlags(cgfTraceFlags[0]),
		TraceState: cfgTraceState,
		Remote:     true,
	}, nil
}
