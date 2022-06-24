package wasmhelper

import (
	"reflect"
	"unsafe"

	"github.com/averagemarcus/wasmhook/pkg/wasmhelper/types"

	"github.com/pquerna/ffjson/ffjson"
)

var Fn func(types.HookRequest) types.HookRequest

//export Run
func Run(ptr, size uint32) uint64 {
	input := ptrToString(ptr, size)

	var out types.HookRequest
	if err := ffjson.Unmarshal([]byte(input), &out); err != nil {
		panic(err)
	}

	out = Fn(out)

	result, _ := ffjson.Marshal(out)
	return stringToPtr(string(result))
}

func ptrToString(ptr uint32, size uint32) (ret string) {
	strHdr := (*reflect.StringHeader)(unsafe.Pointer(&ret))
	strHdr.Data = uintptr(ptr)
	strHdr.Len = uintptr(size)
	return
}

func stringToPtr(s string) uint64 {
	buf := []byte(s)
	ptr := &buf[0]
	unsafePtr := uintptr(unsafe.Pointer(ptr))
	return (uint64(unsafePtr) << uint64(32)) | uint64(len(buf))
}
