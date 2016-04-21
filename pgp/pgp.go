package pgp

import (
	"bytes"

	"github.com/optdyn/gorjun/db"

	"github.com/subutai-io/base/agent/log"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/clearsign"
)

func Verify(name, message string) string {
	entity, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(db.UserKey(name)))
	log.Check(log.WarnLevel, "Reading user public key", err)

	block, _ := clearsign.Decode([]byte(message))

	_, err = openpgp.CheckDetachedSignature(entity, bytes.NewBuffer(block.Bytes), block.ArmoredSignature.Body)
	if log.Check(log.WarnLevel, "Checking signature", err) {
		return ""
	}
	return string(block.Bytes)
}