package sessions

import (
	"crypto/rand"
	"encoding/hex"
)

// uuid representation compliant with specification
// described in RFC 4122.
type uuid [16]byte

// Used in string method conversion
const dash byte = '-'

func safeRandom(dest []byte) {
	if _, err := rand.Read(dest); err != nil {
		panic(err)
	}
}

// Bytes returns bytes slice representation of UUID.
func (u uuid) Bytes() []byte {
	return u[:]
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u uuid) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = dash
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = dash
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = dash
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = dash
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

// SetVersion sets version bits.
func (u *uuid) SetVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

// SetVariant sets variant bits as described in RFC 4122.
func (u *uuid) SetVariant() {
	u[8] = (u[8] & 0xbf) | 0x80
}

// newToken returns random generated UUID.
func newToken() string {
	u := uuid{}
	safeRandom(u[:])
	u.SetVersion(4)
	u.SetVariant()
	return u.String()
}
