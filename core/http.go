package huacache

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const defaultBasePath = "/huacache/"

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 添加 CORS 头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<action> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 1)
	if len(parts) < 1 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	action := ""
	if len(parts) == 1 {
		action = parts[0]
	}
	switch action {
	case GET_KEY:
		p.handleGetAction(w, r)
	case SET_KEY:
		p.handleSetAction(w, r)
	case DEL_KEY:
		p.handleDelAction(w, r)
	case LIST_GROUP:
		p.handleListGroupsAction(w)
	case NEW_GROUP:
		p.handleNewGroupAction(w, r)
	default:
		http.Error(w, "not supported action: "+action, http.StatusBadRequest)
	}
}

func (p *HTTPPool) handleGetAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(HTTP_BODY_DEFAULT_MAX_SIZE); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("handleGetAction: %s", r.Form)
	// 从表单中获取 "key" 的值
	key := r.FormValue("key")
	groupName := r.FormValue("group")
	group, err := GetGroup(groupName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

func (p *HTTPPool) handleSetAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(HTTP_BODY_DEFAULT_MAX_SIZE); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// 从表单中获取 "key" 的值
	key := r.FormValue("key")
	value := r.FormValue("value")
	groupName := r.FormValue("group")
	group, err := GetGroup(groupName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := group.AddOrUpdate(key, ByteView{B: []byte(value)}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write([]byte("set key success"))
}

func (p *HTTPPool) handleDelAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(HTTP_BODY_DEFAULT_MAX_SIZE); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// 从表单中获取 "key" 的值
	key := r.FormValue("key")
	groupName := r.FormValue("group")
	group, err := GetGroup(groupName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := group.Delete(key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write([]byte("del key success"))
}

func (p *HTTPPool) handleListGroupsAction(w http.ResponseWriter) {
	// 获取所有组的名称
	groups, _ := ListGroups()
	// 以 JSON 格式返回组的名称列表
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *HTTPPool) handleNewGroupAction(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(HTTP_BODY_DEFAULT_MAX_SIZE); err != nil {
		http.Error(w, err.Error(), http.StatusOK)
	}
	name := r.FormValue("name")
	capcity := r.FormValue("capacity")
	//capacity转int64
	capacity, err := strconv.ParseInt(capcity, 10, 64)
	if err != nil {
		http.Error(w, "capacity must be a number", http.StatusBadRequest)
		return
	}
	capacity *= MB
	_, err = NewGroup(name, capacity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}
	response, _ := json.Marshal("success create group:" + name)
	w.Write(response)
}
