package message

import (
	"github.com/tomsobpl/otel-gelf-converter/internal/helpers"
	"go.opentelemetry.io/collector/pdata/plog"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

type Message struct {
	rawmsg *gelf.Message
}

// NewMessage creates a new message
func NewMessage() *Message {
	return &Message{rawmsg: &gelf.Message{
		Version:  "1.1",
		Host:     "UNKNOWN",
		Extra:    map[string]interface{}{},
		RawExtra: nil,
	}}
}

// GetRawMessage returns the raw message
func (m *Message) GetRawMessage() *gelf.Message {
	return m.rawmsg
}

// SetBody sets the body of the message
func (m *Message) SetBody(body string) {
	m.rawmsg.Short = body
}

// SetHost sets the host of the message
func (m *Message) SetHost(host string) {
	m.rawmsg.Host = host
}

// SetSeverity sets the severity of the message
func (m *Message) SetSeverity(severity int32) {
	m.rawmsg.Level = severity
}

// SetTimestamp sets the timestamp of the message
func (m *Message) SetTimestamp(timestamp float64) {
	m.rawmsg.TimeUnix = timestamp
}

// UpdateExtraFields updates the extra fields of the message
func (m *Message) UpdateExtraFields(fields map[string]interface{}) {
	for k, v := range fields {
		m.rawmsg.Extra[k] = v
	}
}

// UpdateFromLogRecord updates the message from an otel plog.LogRecord data
func (m *Message) UpdateFromLogRecord(lr plog.LogRecord) {
	m.SetBody(lr.Body().AsString())
	m.SetSeverity(helpers.ConvertOtelSeverityToSyslogLevel(int32(lr.SeverityNumber())))
	m.SetTimestamp(helpers.ConvertOtelTimestampToGelfTimeUnix(lr.Timestamp(), lr.ObservedTimestamp()))

	m.UpdateExtraFields(map[string]interface{}{
		"otel_log_dropped_attributes_count": lr.DroppedAttributesCount(),
		"otel_log_event_name":               lr.EventName(),
		"otel_log_severity_number":          lr.SeverityNumber(),
		"otel_log_severity_text":            lr.SeverityText(),
		"otel_log_span_id":                  lr.SpanID().String(),
		"otel_log_trace_id":                 lr.TraceID().String(),
	})

	m.UpdateExtraFields(helpers.MapOtelAttributes(lr.Attributes()))
}
