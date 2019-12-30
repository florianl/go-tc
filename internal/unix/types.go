/*
Package internal/unix contains constants, that are needed to use github.com/florianl/go-tc.
Some of these constants are copied from golang.org/x/sys/unix to make them available also to
other OS than linux. In the end, the only source of truth will be the Linux kernel itself.
*/
package unix

const (
	LINKLAYER_UNSPEC = iota
	LINKLAYER_ETHERNET
	LINKLAYER_ATM
)

const (
	ATM_CELL_PAYLOAD = 48
	ATM_CELL_SIZE    = 53
)
