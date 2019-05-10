tc  [![GoDoc](https://godoc.org/github.com/florianl/go-tc?status.svg)](https://godoc.org/github.com/florianl/go-tc)
==
This is a work in progress version of `tc`.  It provides a [C](https://en.wikipedia.org/wiki/C_(programming_language))-binding free API to the netlink based traffic control system of [rtnetlink](http://man7.org/linux/man-pages/man7/rtnetlink.7.html).

Example
-------

```golang
func main() {
	// open a rtnetlink socket
	rtnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open rtnetlink socket: %v\n", err)
		return
	}
	defer func() {
		if err := rtnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

    // get all the qdiscs from all interfaces
	qdiscs, err := rtnl.Qdisc().Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get qdiscs: %v\n", err)
		return
	}

	for _, qdisc := range qdiscs {
		iface, err := net.InterfaceByIndex(int(qdisc.Ifindex))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get interface from id %d: %v", qdisc.Ifindex, err)
			return
		}
		fmt.Printf("%20s\t%s\n", iface.Name, qdisc.Kind)
	}
}
```