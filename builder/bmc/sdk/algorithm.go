package oraclesdk
//Original inspiration :
//https://github.com/99designs/httpsignatures-go
//We should consider contributing this back.
import (
	"crypto/sha256"
	"crypto/sha1"
	"errors"
	"hash"
)

var (
	AlgorithmRsaSha256 = &Algorithm{"rsa-sha256", sha256.New}
	AlgorithmRsaSha1 = &Algorithm{"rsa-sha1", sha1.New}
	ErrorUnknownAlgorithm = errors.New("Unknown Algorithm")
)

type Algorithm struct {
	name string
	hash func() hash.Hash
}

func algorithmFromString(name string) (*Algorithm, error) {
	switch name {
	case AlgorithmRsaSha256.name:
		return AlgorithmRsaSha256, nil
	case AlgorithmRsaSha1.name:
		return AlgorithmRsaSha1, nil
	}

	return nil, ErrorUnknownAlgorithm
}
