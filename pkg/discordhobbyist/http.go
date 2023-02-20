package discordhobbyist

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type HTTPServer struct {
	log                     logrus.FieldLogger
	router                  *httprouter.Router
	channelRequestCallbacks []func(ctx context.Context, params httprouter.Params, body []byte) error
}

func NewHTTPServer(log logrus.FieldLogger, channelRequestCallbacks []func(ctx context.Context, params httprouter.Params, body []byte) error) *HTTPServer {
	return &HTTPServer{
		log:                     log,
		router:                  httprouter.New(),
		channelRequestCallbacks: channelRequestCallbacks,
	}
}

func (h *HTTPServer) Start(addr string) error {
	h.log.WithField("addr", addr).Info("starting http server")

	if err := h.RegisterRoutes(); err != nil {
		return err
	}

	//nolint:gosec // don't really care about this
	return http.ListenAndServe(addr, h.router)
}

func (h *HTTPServer) RegisterRoutes() error {
	h.router.POST("/channels/:group/:channel", h.wrappedHandler(h.handleChannelRequest))

	return nil
}

func (h *HTTPServer) wrappedHandler(handler func(ctx context.Context, r *http.Request, p httprouter.Params) (*http.Response, error)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		logCtx := h.log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		})

		defer r.Body.Close()

		resp, err := handler(ctx, r, ps)
		if err != nil {
			logCtx.WithError(err).Error("error handling request")

			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
	}
}

func (h *HTTPServer) handleChannelRequest(ctx context.Context, r *http.Request, ps httprouter.Params) (*http.Response, error) {
	if err := r.ParseMultipartForm(0); err != nil {
		return nil, errors.New("error parsing multipart form")
	}

	body := []byte(r.FormValue("payload_json"))

	allErrors := make([]error, 0)

	for _, callback := range h.channelRequestCallbacks {
		err := callback(ctx, ps, body)
		if err != nil {
			allErrors = append(allErrors, err)

			h.log.WithError(err).Error("error handling channel request")
		}
	}

	if len(allErrors) > 0 {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(nil),
		}, allErrors[0]
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(nil),
	}, nil
}
