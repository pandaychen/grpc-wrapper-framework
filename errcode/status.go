package errcode

import (
	"grpc-wrapper-framework/errcode/types"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

//incase of import error
var _ Codes = &Status{}

// Status statusError is an alias of a status proto
// implement ecode.Codes

//Codes接口的另外一个实例化实现
type Status struct {
	s *types.Status //封装类似grpc.Status pb结构体
}

// Proto return origin protobuf message
func (s *Status) Proto() *types.Status {
	return s.s
}

// Error implement error
func (s *Status) Error() string {
	return s.Message()
}

// Code return error code
func (s *Status) Code() int {
	return int(s.s.Code)
}

// Message return error message for developer
func (s *Status) Message() string {
	if s.s.Message == "" {
		return strconv.Itoa(int(s.s.Code))
	}

	//
	return s.s.Message
}

// Details return error details
func (s *Status) Details() []interface{} {
	if s == nil || s.s == nil {
		return nil
	}
	details := make([]interface{}, 0, len(s.s.Details))
	for _, any := range s.s.Details {
		detail := &ptypes.DynamicAny{}
		if err := ptypes.UnmarshalAny(any, detail); err != nil {
			details = append(details, err)
			continue
		}
		details = append(details, detail.Message)
	}
	return details
}

// WithDetails WithDetails
func (s *Status) WithDetails(pbs ...proto.Message) (*Status, error) {
	for _, pb := range pbs {
		anyMsg, err := ptypes.MarshalAny(pb)
		if err != nil {
			return s, err
		}
		s.s.Details = append(s.s.Details, anyMsg)
	}
	return s, nil
}
