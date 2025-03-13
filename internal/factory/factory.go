package factory

import (
	"github.com/tomsobpl/otel-gelf-converter/internal/helpers"
	"github.com/tomsobpl/otel-gelf-converter/internal/message"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type Factory struct {
	logger *zap.Logger
}

// NewFactory creates a new factory
func NewFactory(logger *zap.Logger) *Factory {
	return &Factory{
		logger: logger,
	}
}

func (f *Factory) FromOtelLogsData(logs plog.Logs) []*message.Message {
	messages := make([]*message.Message, 0)

	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		messages = append(messages, f.parseResourceLogs(logs.ResourceLogs().At(i))...)
	}

	return messages
}

func (f *Factory) parseResourceLogs(rl plog.ResourceLogs) []*message.Message {
	messages := make([]*message.Message, 0)

	for i := 0; i < rl.ScopeLogs().Len(); i++ {
		messages = append(messages, f.parseScopeLogs(rl.ScopeLogs().At(i))...)
	}

	for _, m := range messages {
		m.UpdateExtraFields(map[string]interface{}{
			"otel_resource_dropped_attributes_count": rl.Resource().DroppedAttributesCount(),
		})

		m.UpdateExtraFields(helpers.MapOtelAttributesWithPrefix(rl.Resource().Attributes(), "resource"))

		host, hostExist := rl.Resource().Attributes().Get("host.name")

		if hostExist {
			m.SetHost(host.AsString())
		}
	}

	return messages
}

func (f *Factory) parseScopeLogs(sl plog.ScopeLogs) []*message.Message {
	messages := make([]*message.Message, 0)

	for i := 0; i < sl.LogRecords().Len(); i++ {
		messages = append(messages, f.parseLogRecord(sl.LogRecords().At(i)))
	}

	for _, m := range messages {
		m.UpdateExtraFields(map[string]interface{}{
			"otel_scope_dropped_attributes_count": sl.Scope().DroppedAttributesCount(),
			"otel_scope_name":                     sl.Scope().Name(),
			"otel_scope_version":                  sl.Scope().Version(),
		})

		m.UpdateExtraFields(helpers.MapOtelAttributesWithPrefix(sl.Scope().Attributes(), "scope"))
	}

	return messages
}

func (f *Factory) parseLogRecord(lr plog.LogRecord) *message.Message {
	m := message.NewMessage()
	m.UpdateFromLogRecord(lr)
	return m
}
