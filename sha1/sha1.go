package sha1

// #if defined( __APPLE__ )
//     #define COMMON_DIGEST_FOR_OPENSSL
//     #include <CommonCrypto/CommonDigest.h>
//     #define SHA1                CC_SHA1
// #else
//    #cgo LDFLAGS: -lcrypto
//    #include <openssl/sha.h>
// #endif
import "C"

import (
	"crypto"
	"errors"
	"hash"
	"unsafe"
)

func init() {
	crypto.RegisterHash(crypto.SHA1, New)
}

// The size of an SHA1 checksum in bytes.
const Size = 20

// The blocksize of SHA1 in bytes.
const BlockSize = 64

type digest struct {
	context *_Ctype_SHA_CTX
}

// New returns a new hash.Hash computing the SHA1 checksum.
func New() hash.Hash {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Reset() {
	d.context = &_Ctype_SHA_CTX{}
	C.SHA1_Init(d.context)
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Write(p []byte) (nn int, err error) {
	if len(p) == 0 || C.SHA1_Update(d.context, unsafe.Pointer(&p[0]),
		C.size_t(len(p))) == 1 {
		return len(p), nil
	}

	return 0, errors.New("SHA1_Update failed")
}

func (d *digest) Sum(in []byte) []byte {
	context := *d.context
	defer func() { *d.context = context }()

	md := make([]byte, Size)
	if C.SHA1_Final((*_Ctype_unsignedchar)(&md[0]), d.context) == 1 {
		return append(in, md...)
	}

	return nil
}
