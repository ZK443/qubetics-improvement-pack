package store

var (
	KeyParams = []byte("bridge/params") // сериализованные types.Params
	KeyACL    = []byte("bridge/acl/")   // per-address: bridge/acl/<addr> -> {0|1}
)
