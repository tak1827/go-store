package pb

import (
	"github.com/lithdew/bytesutil"
)

func (x *Person) StoreKey() []byte {
	id := x.GetId()
	return bytesutil.AppendUint32BE(nil, uint32(id))
}

func PersonIdFromStoreKey(key []byte) int32 {
	return int32(bytesutil.Uint32BE(key))
}
