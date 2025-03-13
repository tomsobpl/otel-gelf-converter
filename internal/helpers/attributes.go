package helpers

import "go.opentelemetry.io/collector/pdata/pcommon"

// MapOtelAttributes converts OpenTelemetry attributes to GELF extra fields.
func MapOtelAttributes(attributes pcommon.Map) map[string]interface{} {
	return MapOtelAttributesWithPrefix(attributes, "")
}

// MapOtelAttributesWithPrefix converts OpenTelemetry attributes to GELF extra fields with a prefix.
func MapOtelAttributesWithPrefix(attributes pcommon.Map, prefix string) map[string]interface{} {
	fields := make(map[string]interface{})

	attributes.Range(func(k string, v pcommon.Value) bool {
		if prefix != "" {
			k = prefix + "." + k
		}

		switch v.Type() {
		case pcommon.ValueTypeMap:
			//@TODO Handle nested maps if needed
		case pcommon.ValueTypeSlice:
			//@TODO Handle slices if needed
		default:
			fields[k] = v.AsString()
		}

		return true
	})

	return fields
}
