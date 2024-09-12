package traceparent

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const (
	// TraceparentVersion is the version of the traceparent header.
	TraceparentVersion = "00"
	// TraceparentFlag indicates that the trace is sampled.
	TraceparentFlag = "01" // sampled

	// TraceparentInvalidParentID is the value used for an invalid parentID.
	TraceparentInvalidParentID = "0000000000000000" // invalid parentID

	// TraceparentHTTPHeaderTag is the HTTP header tag for traceparent.
	TraceparentHTTPHeaderTag = "traceparent"

	// traceparentParts is the number of parts in a traceparent header.
	traceparentParts = 4
)

var (
	// reSpanID is a regular expression to validate span IDs.
	reSpanID = regexp.MustCompile(`^[0-9a-f]{16}$`)

	// errTraceparentInvalidFormat is returned when the traceparent format is invalid.
	errTraceparentInvalidFormat = errors.New("invalid traceparent format")
	// errSpanIDInvalidFormat is returned when the span ID format is invalid.
	errSpanIDInvalidFormat = errors.New("invalid spanID format")
)

// Traceparent represents the traceparent value.
type Traceparent struct {
	version  string
	traceID  string
	parentID string
	flags    string
}

// Serialize converts the Traceparent struct to a string.
func (tp *Traceparent) Serialize() string {
	return fmt.Sprintf("%s-%s-%s-%s", tp.version, tp.traceID, tp.parentID, tp.flags)
}

// WithNewParentID returns a new Traceparent with the provided parentID.
func (tp *Traceparent) WithNewParentID(parentID string) (Traceparent, error) {
	if !reSpanID.MatchString(parentID) {
		return Traceparent{}, fmt.Errorf("%w: %s", errSpanIDInvalidFormat, parentID)
	}

	return Traceparent{
		version:  tp.version,
		traceID:  tp.traceID,
		parentID: parentID,
		flags:    tp.flags,
	}, nil
}

// New creates a new Traceparent with a generated traceID.
func New() (Traceparent, error) {
	traceID, err := newTraceID()
	if err != nil {
		return Traceparent{}, fmt.Errorf("failed to generate traceID: %w", err)
	}

	return Traceparent{
		version:  TraceparentVersion,
		traceID:  traceID,
		parentID: TraceparentInvalidParentID,
		flags:    TraceparentFlag,
	}, nil
}

// newTraceID generates a new trace ID using UUID v7.
func newTraceID() (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	return strings.ReplaceAll(uuid.String(), "-", ""), nil
}

// Deserialize parses a traceparent string and returns a Traceparent struct.
func Deserialize(str string) (Traceparent, error) {
	var tpt Traceparent

	n, err := fmt.Sscanf(str, "%2s-%32s-%16s-%2s", &tpt.version, &tpt.traceID, &tpt.parentID, &tpt.flags)
	if err != nil {
		return Traceparent{}, fmt.Errorf("failed to parse traceparent: %w", err)
	}

	if n != traceparentParts {
		return Traceparent{}, fmt.Errorf("%w: %s", errTraceparentInvalidFormat, str)
	}

	return tpt, nil
}
