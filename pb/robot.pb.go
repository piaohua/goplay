// Code generated by protoc-gen-gogo.
// source: robot.proto
// DO NOT EDIT!

package pb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// 机器人消息
type RobotMsg struct {
	Code  string `protobuf:"bytes,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Num   uint32 `protobuf:"varint,2,opt,name=Num,proto3" json:"Num,omitempty"`
	Rtype uint32 `protobuf:"varint,3,opt,name=Rtype,proto3" json:"Rtype,omitempty"`
}

func (m *RobotMsg) Reset()                    { *m = RobotMsg{} }
func (*RobotMsg) ProtoMessage()               {}
func (*RobotMsg) Descriptor() ([]byte, []int) { return fileDescriptorRobot, []int{0} }

func (m *RobotMsg) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *RobotMsg) GetNum() uint32 {
	if m != nil {
		return m.Num
	}
	return 0
}

func (m *RobotMsg) GetRtype() uint32 {
	if m != nil {
		return m.Rtype
	}
	return 0
}

type RobotLogin struct {
	Phone string `protobuf:"bytes,1,opt,name=Phone,proto3" json:"Phone,omitempty"`
}

func (m *RobotLogin) Reset()                    { *m = RobotLogin{} }
func (*RobotLogin) ProtoMessage()               {}
func (*RobotLogin) Descriptor() ([]byte, []int) { return fileDescriptorRobot, []int{1} }

func (m *RobotLogin) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

type RobotReLogin struct {
	Phone string `protobuf:"bytes,1,opt,name=Phone,proto3" json:"Phone,omitempty"`
	Code  string `protobuf:"bytes,2,opt,name=Code,proto3" json:"Code,omitempty"`
	Rtype uint32 `protobuf:"varint,3,opt,name=Rtype,proto3" json:"Rtype,omitempty"`
}

func (m *RobotReLogin) Reset()                    { *m = RobotReLogin{} }
func (*RobotReLogin) ProtoMessage()               {}
func (*RobotReLogin) Descriptor() ([]byte, []int) { return fileDescriptorRobot, []int{2} }

func (m *RobotReLogin) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *RobotReLogin) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *RobotReLogin) GetRtype() uint32 {
	if m != nil {
		return m.Rtype
	}
	return 0
}

type RobotLogout struct {
	Phone string `protobuf:"bytes,1,opt,name=Phone,proto3" json:"Phone,omitempty"`
	Code  string `protobuf:"bytes,2,opt,name=Code,proto3" json:"Code,omitempty"`
}

func (m *RobotLogout) Reset()                    { *m = RobotLogout{} }
func (*RobotLogout) ProtoMessage()               {}
func (*RobotLogout) Descriptor() ([]byte, []int) { return fileDescriptorRobot, []int{3} }

func (m *RobotLogout) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *RobotLogout) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func init() {
	proto.RegisterType((*RobotMsg)(nil), "pb.RobotMsg")
	proto.RegisterType((*RobotLogin)(nil), "pb.RobotLogin")
	proto.RegisterType((*RobotReLogin)(nil), "pb.RobotReLogin")
	proto.RegisterType((*RobotLogout)(nil), "pb.RobotLogout")
}
func (this *RobotMsg) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*RobotMsg)
	if !ok {
		that2, ok := that.(RobotMsg)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Code != that1.Code {
		return false
	}
	if this.Num != that1.Num {
		return false
	}
	if this.Rtype != that1.Rtype {
		return false
	}
	return true
}
func (this *RobotLogin) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*RobotLogin)
	if !ok {
		that2, ok := that.(RobotLogin)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Phone != that1.Phone {
		return false
	}
	return true
}
func (this *RobotReLogin) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*RobotReLogin)
	if !ok {
		that2, ok := that.(RobotReLogin)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Phone != that1.Phone {
		return false
	}
	if this.Code != that1.Code {
		return false
	}
	if this.Rtype != that1.Rtype {
		return false
	}
	return true
}
func (this *RobotLogout) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*RobotLogout)
	if !ok {
		that2, ok := that.(RobotLogout)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Phone != that1.Phone {
		return false
	}
	if this.Code != that1.Code {
		return false
	}
	return true
}
func (this *RobotMsg) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&pb.RobotMsg{")
	s = append(s, "Code: "+fmt.Sprintf("%#v", this.Code)+",\n")
	s = append(s, "Num: "+fmt.Sprintf("%#v", this.Num)+",\n")
	s = append(s, "Rtype: "+fmt.Sprintf("%#v", this.Rtype)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *RobotLogin) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&pb.RobotLogin{")
	s = append(s, "Phone: "+fmt.Sprintf("%#v", this.Phone)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *RobotReLogin) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&pb.RobotReLogin{")
	s = append(s, "Phone: "+fmt.Sprintf("%#v", this.Phone)+",\n")
	s = append(s, "Code: "+fmt.Sprintf("%#v", this.Code)+",\n")
	s = append(s, "Rtype: "+fmt.Sprintf("%#v", this.Rtype)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *RobotLogout) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&pb.RobotLogout{")
	s = append(s, "Phone: "+fmt.Sprintf("%#v", this.Phone)+",\n")
	s = append(s, "Code: "+fmt.Sprintf("%#v", this.Code)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringRobot(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *RobotMsg) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RobotMsg) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Code) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintRobot(dAtA, i, uint64(len(m.Code)))
		i += copy(dAtA[i:], m.Code)
	}
	if m.Num != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintRobot(dAtA, i, uint64(m.Num))
	}
	if m.Rtype != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintRobot(dAtA, i, uint64(m.Rtype))
	}
	return i, nil
}

func (m *RobotLogin) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RobotLogin) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Phone) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintRobot(dAtA, i, uint64(len(m.Phone)))
		i += copy(dAtA[i:], m.Phone)
	}
	return i, nil
}

func (m *RobotReLogin) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RobotReLogin) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Phone) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintRobot(dAtA, i, uint64(len(m.Phone)))
		i += copy(dAtA[i:], m.Phone)
	}
	if len(m.Code) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintRobot(dAtA, i, uint64(len(m.Code)))
		i += copy(dAtA[i:], m.Code)
	}
	if m.Rtype != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintRobot(dAtA, i, uint64(m.Rtype))
	}
	return i, nil
}

func (m *RobotLogout) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RobotLogout) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Phone) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintRobot(dAtA, i, uint64(len(m.Phone)))
		i += copy(dAtA[i:], m.Phone)
	}
	if len(m.Code) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintRobot(dAtA, i, uint64(len(m.Code)))
		i += copy(dAtA[i:], m.Code)
	}
	return i, nil
}

func encodeFixed64Robot(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Robot(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintRobot(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *RobotMsg) Size() (n int) {
	var l int
	_ = l
	l = len(m.Code)
	if l > 0 {
		n += 1 + l + sovRobot(uint64(l))
	}
	if m.Num != 0 {
		n += 1 + sovRobot(uint64(m.Num))
	}
	if m.Rtype != 0 {
		n += 1 + sovRobot(uint64(m.Rtype))
	}
	return n
}

func (m *RobotLogin) Size() (n int) {
	var l int
	_ = l
	l = len(m.Phone)
	if l > 0 {
		n += 1 + l + sovRobot(uint64(l))
	}
	return n
}

func (m *RobotReLogin) Size() (n int) {
	var l int
	_ = l
	l = len(m.Phone)
	if l > 0 {
		n += 1 + l + sovRobot(uint64(l))
	}
	l = len(m.Code)
	if l > 0 {
		n += 1 + l + sovRobot(uint64(l))
	}
	if m.Rtype != 0 {
		n += 1 + sovRobot(uint64(m.Rtype))
	}
	return n
}

func (m *RobotLogout) Size() (n int) {
	var l int
	_ = l
	l = len(m.Phone)
	if l > 0 {
		n += 1 + l + sovRobot(uint64(l))
	}
	l = len(m.Code)
	if l > 0 {
		n += 1 + l + sovRobot(uint64(l))
	}
	return n
}

func sovRobot(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozRobot(x uint64) (n int) {
	return sovRobot(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *RobotMsg) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RobotMsg{`,
		`Code:` + fmt.Sprintf("%v", this.Code) + `,`,
		`Num:` + fmt.Sprintf("%v", this.Num) + `,`,
		`Rtype:` + fmt.Sprintf("%v", this.Rtype) + `,`,
		`}`,
	}, "")
	return s
}
func (this *RobotLogin) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RobotLogin{`,
		`Phone:` + fmt.Sprintf("%v", this.Phone) + `,`,
		`}`,
	}, "")
	return s
}
func (this *RobotReLogin) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RobotReLogin{`,
		`Phone:` + fmt.Sprintf("%v", this.Phone) + `,`,
		`Code:` + fmt.Sprintf("%v", this.Code) + `,`,
		`Rtype:` + fmt.Sprintf("%v", this.Rtype) + `,`,
		`}`,
	}, "")
	return s
}
func (this *RobotLogout) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&RobotLogout{`,
		`Phone:` + fmt.Sprintf("%v", this.Phone) + `,`,
		`Code:` + fmt.Sprintf("%v", this.Code) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringRobot(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *RobotMsg) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRobot
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RobotMsg: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RobotMsg: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRobot
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Code = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Num", wireType)
			}
			m.Num = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Num |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rtype", wireType)
			}
			m.Rtype = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Rtype |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRobot(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRobot
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RobotLogin) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRobot
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RobotLogin: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RobotLogin: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Phone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRobot
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Phone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRobot(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRobot
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RobotReLogin) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRobot
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RobotReLogin: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RobotReLogin: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Phone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRobot
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Phone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRobot
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Code = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rtype", wireType)
			}
			m.Rtype = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Rtype |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRobot(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRobot
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *RobotLogout) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRobot
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RobotLogout: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RobotLogout: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Phone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRobot
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Phone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRobot
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Code = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRobot(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRobot
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipRobot(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRobot
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRobot
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthRobot
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRobot
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipRobot(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthRobot = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRobot   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("robot.proto", fileDescriptorRobot) }

var fileDescriptorRobot = []byte{
	// 202 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0xca, 0x4f, 0xca,
	0x2f, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x72, 0xe3, 0xe2, 0x08,
	0x02, 0x09, 0xf9, 0x16, 0xa7, 0x0b, 0x09, 0x71, 0xb1, 0x38, 0xe7, 0xa7, 0xa4, 0x4a, 0x30, 0x2a,
	0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x42, 0x02, 0x5c, 0xcc, 0x7e, 0xa5, 0xb9, 0x12, 0x4c, 0x0a,
	0x8c, 0x1a, 0xbc, 0x41, 0x20, 0xa6, 0x90, 0x08, 0x17, 0x6b, 0x50, 0x49, 0x65, 0x41, 0xaa, 0x04,
	0x33, 0x58, 0x0c, 0xc2, 0x51, 0x52, 0xe2, 0xe2, 0x02, 0x9b, 0xe3, 0x93, 0x9f, 0x9e, 0x99, 0x07,
	0x52, 0x13, 0x90, 0x91, 0x9f, 0x07, 0x33, 0x0a, 0xc2, 0x51, 0xf2, 0xe3, 0xe2, 0x01, 0xab, 0x09,
	0x4a, 0xc5, 0xa3, 0x0a, 0xee, 0x0a, 0x26, 0x24, 0x57, 0x60, 0xb7, 0xd3, 0x9c, 0x8b, 0x1b, 0x66,
	0x67, 0x7e, 0x69, 0x09, 0xf1, 0xc6, 0x39, 0xe9, 0x5c, 0x78, 0x28, 0xc7, 0x70, 0xe3, 0xa1, 0x1c,
	0xc3, 0x87, 0x87, 0x72, 0x8c, 0x0d, 0x8f, 0xe4, 0x18, 0x57, 0x3c, 0x92, 0x63, 0x3c, 0xf1, 0x48,
	0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x5f, 0x3c, 0x92, 0x63, 0xf8, 0xf0,
	0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x24, 0x36, 0x70, 0x68, 0x19, 0x03, 0x02, 0x00, 0x00,
	0xff, 0xff, 0x4a, 0xac, 0x99, 0x16, 0x3c, 0x01, 0x00, 0x00,
}