package kelvin

// Cipher provides encryption for your Kelvin databases.
type Cipher interface {
	Encrypt([]byte) []byte
	Decrpyt([]byte) []byte
}
