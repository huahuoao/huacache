package protocol

import (
	"bytes"
	"fmt"
	"testing"
)

// TestBluebellCodec 测试 Bluebell 的序列化和反序列化
func TestBluebellCodec(t *testing.T) {
	// 定义原始的 Bluebell 结构体
	original := &BluebellRequest{
		Command: "SET",
		Key:     "mykey",
		Value:   []byte("myvalue"),
		Group:   "mygroup",
	}

	// 1. 测试序列化
	serializedData, err := original.Serialize()
	fmt.Print(serializedData)

	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	// 2. 测试反序列化
	deserialized, err := Deserialize(serializedData)
	if err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}

	// 3. 检查反序列化结果是否与原始数据一致
	if deserialized.Command != original.Command {
		t.Errorf("Command 不匹配, 得到: %v, 期望: %v", deserialized.Command, original.Command)
	}

	if deserialized.Key != original.Key {
		t.Errorf("Key 不匹配, 得到: %v, 期望: %v", deserialized.Key, original.Key)
	}

	if !bytes.Equal(deserialized.Value, original.Value) {
		t.Errorf("Value 不匹配, 得到: %v, 期望: %v", deserialized.Value, original.Value)
	}

	if deserialized.Group != original.Group {
		t.Errorf("Group 不匹配, 得到: %v, 期望: %v", deserialized.Group, original.Group)
	}
}
