package groupcache

import (
	"os"

	"github.com/golang/groupcache"
)

var picker groupcache.PeerPicker
var SingleNode = os.Getenv("SINGLE_NODE") != ""

func init() {
	if SingleNode {
		picker = &groupcache.NoPeers{}
	} else {
		picker = groupcache.NewHTTPPoolOpts(os.Getenv("POD_ID"), &groupcache.HTTPPoolOptions{
			Replicas: 2,
		})
	}
}
