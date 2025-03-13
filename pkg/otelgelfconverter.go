package pkg

import (
	"github.com/tomsobpl/otel-gelf-converter/pkg/factory"
	"github.com/tomsobpl/otel-gelf-converter/pkg/message"
	"go.uber.org/zap"
)

// CreateFactory creates a new factory
func CreateFactory(logger *zap.Logger) *factory.Factory {
	return factory.NewFactory(logger)
}

// CreateMessage creates a new message
func CreateMessage() *message.Message {
	return message.NewMessage()
}
