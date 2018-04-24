package service

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"code.ysitd.cloud/component/art/gallery/pkg/modals/artwork"
	"code.ysitd.cloud/component/art/gallery/pkg/modals/exhibition"

	"github.com/sirupsen/logrus"
)

const requestTimeout = 5 * time.Minute

type Handler struct {
	Artwrok    *artwork.Store     `inject:""`
	Exhibition *exhibition.Store  `inject:""`
	Logger     logrus.FieldLogger `inject:"handler logger"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if _, deadline := ctx.Deadline(); !deadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, requestTimeout)
		defer cancel()
	}

	done := make(chan bool)

	go func(done chan bool) {
		h.handleHTTP(ctx, w, r)
		done <- true
	}(done)

	select {
	case <-done:
		close(done)
		break
	case <-ctx.Done():
		if err := ctx.Err(); err == context.DeadlineExceeded {
			http.Error(w, "Timeout", http.StatusGatewayTimeout)
		}
		break
	}
}

func (h *Handler) handleHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	e, err := h.Exhibition.GetExhibition(ctx, r.Host, r.URL.Path)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during routing", http.StatusInternalServerError)
		return
	} else if e == nil {
		http.Error(w, "421 Misdirected Request", 421)
		return
	}

	header := w.Header()

	if match := r.Header.Get("If-None-Match"); match != "" {
		match = strings.Trim(match, "\"'")
		if match == e.Hash {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	header.Set("Content-SHA256", e.Hash)
	header.Set("Etag", "\""+e.Hash+"\"")
	header.Set("Cache-Control", "max-age=36000") // 1 Hours

	origin := r.Header.Get("Origin")

	if origin != "" {
		originUrl, err := url.Parse(origin)
		if err != nil {
			h.Logger.Errorln(err)
			http.Error(w, "Error in parse origin", http.StatusBadRequest)
			return
		}

		if originUrl.Hostname() != r.Host {
			if e.CORS.Valid {
				header.Set("Access-Control-Allow-Origin", e.CORS.String)
			} else {
				return
			}
		}
	}

	blob, err := h.Artwrok.GetWithExhibition(ctx, e)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during loading", http.StatusBadGateway)
		return
	}

	http.ServeContent(w, r, e.Pathname, e.CommitTime, bytes.NewReader(blob))
}
