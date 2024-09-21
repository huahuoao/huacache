package protocol

import (
	"fmt"
	"strconv"

	huacache "github.com/huahuoao/huacache/core"
)

func HandleSetKey(request *BluebellRequest) *BluebellResponse {
	group, _ := huacache.GetGroup(request.Group)
	err := group.AddOrUpdate(request.Key, huacache.ByteView{B: request.Value})
	if err != nil {
		return &BluebellResponse{
			Code:   "500",
			Result: []byte("failed to set key"),
		}
	}
	return &BluebellResponse{
		Code:   "200",
		Result: []byte("OK"),
	}
}

func HandleGetKey(request *BluebellRequest) *BluebellResponse {
	group, _ := huacache.GetGroup(request.Group)
	value, err := group.Get(request.Key)
	if err != nil {
		return &BluebellResponse{
			Code:   "500",
			Result: []byte("failed to get key"),
		}
	}
	return &BluebellResponse{
		Code:   "200",
		Result: value.B,
	}
}

func HandleDeleteKey(request *BluebellRequest) *BluebellResponse {
	group, _ := huacache.GetGroup(request.Group)
	err := group.Delete(request.Key)
	if err != nil {
		return &BluebellResponse{
			Code:   "500",
			Result: []byte("failed to delete key"),
		}
	}
	return &BluebellResponse{
		Code:   "200",
		Result: []byte("OK"),
	}
}

func HandleNewGroup(request *BluebellRequest) *BluebellResponse {
	size, err := strconv.ParseInt(request.Key, 10, 64)
	if err != nil {
		return &BluebellResponse{
			Code:   "400",
			Result: []byte("invalid size"),
		}
	}
	_, err = huacache.NewGroup(request.Group, size)
	if err != nil {
		return &BluebellResponse{
			Code:   "500",
			Result: []byte(err.Error()),
		}
	}
	return &BluebellResponse{
		Code:   "200",
		Result: []byte("OK"),
	}
}

func HandleDeleteGroup(request *BluebellRequest) *BluebellResponse {
	err := huacache.DelGroup(request.Group)
	if err != nil {
		return &BluebellResponse{
			Code:   "500",
			Result: []byte("failed to delete group"),
		}
	}
	return &BluebellResponse{
		Code:   "200",
		Result: []byte("OK"),
	}
}

func HandleListGroup(request *BluebellRequest) *BluebellResponse {
	groups, err := huacache.ListGroups()
	if err != nil {
		return &BluebellResponse{
			Code:   "500",
			Result: []byte("failed to list groups"),
		}
	}
	return &BluebellResponse{
		Code:   "200",
		Result: []byte(fmt.Sprintf("%v", groups)),
	}
}
