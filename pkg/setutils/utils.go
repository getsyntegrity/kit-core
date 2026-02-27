package setutils

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MsgType returns the type name of a protobuf message.
func MsgType(m protoreflect.ProtoMessage) string {
	return string(m.ProtoReflect().Descriptor().FullName())
}
