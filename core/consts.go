package huacache

const (
	MB                         = 1 << 20
	GB                         = 1 << 30
	HTTP_BODY_DEFAULT_MAX_SIZE = 32 * MB
	LIMIT_SIZE                 = 15 * MB
)

// command
const (
	GET_KEY    = "get"
	SET_KEY    = "set"
	DEL_KEY    = "del"
	NEW_GROUP  = "new_group"
	LIST_GROUP = "list_group"
	DEL_GROUP  = "del_group"
	GET_KEYS   = "keys"
)

const (
	CONSISTENTHASH_VIRTUAL_NODE_NUM = 160
)
