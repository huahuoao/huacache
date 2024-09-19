package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/bytedance/sonic"
	huacache "github.com/huahuoao/huacache/core"
	"github.com/panjf2000/gnet/v2"
)

// Bluebell 消息结构
type BluebellRequest struct {
	Command string
	Key     string // 键，通常是用于标识数据的字符串
	Value   []byte // 值，存储数据的字节数组
	Group   string // 组，表示消息所属的组或类别
}
type BluebellResponse struct {
	Code   string
	Result []byte // 响应数据
}

func (b *BluebellResponse) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := writeString(buf, b.Code); err != nil {
		return nil, err
	}

	if err := writeBytes(buf, b.Result); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
func (b *BluebellResponse) Encode() ([]byte, error) {
	// Serialize the message body
	body, err := b.Serialize()
	if err != nil {
		return nil, err
	}

	// Create a buffer to hold the final encoded message
	buf := new(bytes.Buffer)

	// Write the length of the body as the first 4 bytes
	bodyLength := uint32(len(body))
	if err := binary.Write(buf, binary.BigEndian, bodyLength); err != nil {
		return nil, err
	}

	// Write the serialized body
	buf.Write(body)

	return buf.Bytes(), nil
}
func DeserializeResponse(data []byte) (*BluebellResponse, error) {
	buf := bytes.NewBuffer(data)

	code, err := readString(buf)
	if err != nil {
		return nil, err
	}

	result, err := readBytes(buf)
	if err != nil {
		return nil, err
	}

	return &BluebellResponse{
		Code:   code,
		Result: result,
	}, nil
}
func (b *BluebellRequest) String() string {
	return fmt.Sprintf("Bluebell{\n  Command: %s,\n  Key: %s,\n  Value: %s,\n  Group: %s\n}",
		b.Command,
		b.Key,
		string(b.Value), // 将 []byte 转换为 string
		b.Group,
	)
}
func (b *BluebellRequest) Encode() ([]byte, error) {
	// 1. 序列化 Bluebell 结构体
	serializedData, err := b.Serialize()
	if err != nil {
		return nil, err
	}
	// 2. 计算总长度：头部长度（4 字节） + 实际序列化数据长度
	totalLength := len(serializedData) + 4
	// 3. 创建最终数据字节数组
	finalData := make([]byte, totalLength)
	// 4. 将总长度写入前4个字节
	binary.BigEndian.PutUint32(finalData, uint32(len(serializedData)))
	// 5. 将序列化数据复制到总长度后的部分
	copy(finalData[4:], serializedData)
	return finalData, nil
}

// 序列化：将 Bluebell 结构体序列化为二进制
func (b *BluebellRequest) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Command 字段
	if err := writeString(buf, b.Command); err != nil {
		return nil, err
	}

	// Key 字段
	if err := writeString(buf, b.Key); err != nil {
		return nil, err
	}

	// Value 字段
	if err := writeBytes(buf, b.Value); err != nil {
		return nil, err
	}

	// Group 字段
	if err := writeString(buf, b.Group); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 反序列化：将二进制数据反序列化为 Bluebell 结构体
func Deserialize(data []byte) (*BluebellRequest, error) {
	buf := bytes.NewReader(data)
	b := &BluebellRequest{}

	// Command 字段
	command, err := readString(buf)
	if err != nil {
		return nil, err
	}
	b.Command = command

	// Key 字段
	key, err := readString(buf)
	if err != nil {
		return nil, err
	}
	b.Key = key

	// Value 字段
	value, err := readBytes(buf)
	if err != nil {
		return nil, err
	}
	b.Value = value

	// Group 字段
	group, err := readString(buf)
	if err != nil {
		return nil, err
	}
	b.Group = group

	return b, nil
}

// writeString 将字符串以长度+内容的形式写入到缓冲区
func writeString(buf *bytes.Buffer, s string) error {
	length := uint32(len(s))
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := buf.Write([]byte(s))
	return err
}

// writeBytes 将 []byte 以长度+内容的形式写入到缓冲区
func writeBytes(buf *bytes.Buffer, b []byte) error {
	length := uint32(len(b))
	if err := binary.Write(buf, binary.BigEndian, length); err != nil {
		return err
	}
	_, err := buf.Write(b)
	return err
}

// readString 从缓冲区中读取字符串（先读取长度，再读取内容）
func readString(buf io.Reader) (string, error) {
	var length uint32
	if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
		return "", err
	}
	strBuf := make([]byte, length)
	if _, err := io.ReadFull(buf, strBuf); err != nil {
		return "", err
	}
	return string(strBuf), nil
}

// readBytes 从缓冲区中读取 []byte（先读取长度，再读取内容）
func readBytes(buf io.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	byteBuf := make([]byte, length)
	if _, err := io.ReadFull(buf, byteBuf); err != nil {
		return nil, err
	}
	return byteBuf, nil
}

// BluebellServer 实现 gnet 的 Server
type BluebellServer struct {
	*gnet.BuiltinEventEngine
	buffer       map[gnet.Conn]*bytes.Buffer // 用于存储半包数据
	eng          gnet.Engine
	Network      string
	Addr         string
	Multicore    bool
	connected    int32
	disconnected int32
	inBufferPool *sync.Pool
}

// 创建新服务
func NewBluebellServer(network, addr string, multicore bool) *BluebellServer {
	return &BluebellServer{
		buffer:    make(map[gnet.Conn]*bytes.Buffer),
		Network:   network,
		Addr:      addr,
		Multicore: multicore,
		inBufferPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, huacache.LIMIT_SIZE) // 预先创建缓冲区
			},
		}}
}

func SonicSerialize(b interface{}) []byte {
	jsonBytes, err := sonic.Marshal(b)
	if err != nil {
		return nil
	}
	return jsonBytes
}
