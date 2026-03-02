package setutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestMsgType(t *testing.T) {
	msg := &emptypb.Empty{}
	typ := MsgType(msg)
	assert.Contains(t, typ, "Empty")
}
