package groupcache

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/golang/groupcache"
	"github.com/sirupsen/logrus"
)

const dnsTimeout = 5 * time.Second

func UpdatePeer(logger logrus.FieldLogger) {
	ctx, cancel := context.WithTimeout(context.Background(), dnsTimeout)
	defer cancel()
	results, err := net.DefaultResolver.LookupIPAddr(ctx, os.Getenv("SERVICE_NAME"))
	if err != nil {
		logger.Errorln(err)
		return
	}

	peers := make([]string, len(results))
	for _, result := range results {
		peers = append(peers, fmt.Sprintf("http://%s:50005", result.IP.String()))
	}

	pool := picker.(*groupcache.HTTPPool)
	pool.Set(peers...)
}
