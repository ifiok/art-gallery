package service

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"code.ysitd.cloud/art/gallery/pkg/modals/artwork"
	"code.ysitd.cloud/art/gallery/pkg/modals/exhibition"

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
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		h.handleGet(ctx, w, r)
		return
	} else if r.Method == http.MethodOptions {
		h.handlePreflight(ctx, w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Set("Allow", "OPTIONS, GET, HEAD")
}

func (h *Handler) handlePreflight(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	e, err := h.Exhibition.GetExhibitionWithHost(ctx, r.Host)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during routing", http.StatusInternalServerError)
		return
	} else if e == nil {
		http.Error(w, "421 Misdirected Request", 421)
		return
	}
	if !h.handleOrigin(e, w, r) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	e, err := h.Exhibition.GetExhibitionWithPath(ctx, r.Host, r.URL.Path)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during routing", http.StatusInternalServerError)
		return
	} else if e == nil {
		http.Error(w, "421 Misdirected Request", 421)
		return
	}

	if h.handleOrigin(e, w, r) {
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
	header.Set("Cache-Control", "max-age=14400") // 4 Hours

	blob, err := h.Artwrok.GetWithExhibition(ctx, e)
	if err != nil {
		h.Logger.Errorln(err)
		http.Error(w, "Error during loading", http.StatusBadGateway)
		return
	}

	http.ServeContent(w, r, e.Pathname, e.CommitTime, bytes.NewReader(blob))
}

func (h *Handler) handleOrigin(e *exhibition.Exhibition, w http.ResponseWriter, r *http.Request) (finish bool) {
	if origin := r.Header.Get("Origin"); origin != "" {
		originUrl, err := url.Parse(origin)
		if err != nil {
			h.Logger.Errorln(err)
			http.Error(w, "Error in parse origin", http.StatusBadRequest)
			return true
		}

		header := w.Header()

		if e.CORS.Valid && validateCorsOrigin(originUrl.Hostname(), e.CORS.String) {
			header.Set("Access-Control-Allow-Origin", e.CORS.String)
			if e.CORS.String == "*" {
				header.Add("Vary", "Origin")
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
		return true
	}

	return false
}
