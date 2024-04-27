package trace_test

import (
	"testing"

	"github.com/ronny/clog/trace"
	"github.com/stretchr/testify/assert"
)

func TestCloudTraceContext(t *testing.T) {
	validExample := "70e0091f6f5d4643bb4eca9d81320c76/97123319527522;o=1"

	tr, err := trace.ParseCloudTraceContext(validExample)
	if !assert.Nil(t, err) {
		return
	}

	assert.Equal(t, "70e0091f6f5d4643bb4eca9d81320c76", tr.GetTraceID())
	assert.Equal(t, "97123319527522", tr.GetSpanID())
	assert.Equal(t, true, tr.Sampled())
}
