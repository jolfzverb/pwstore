package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"

	"github.com/jolfzverb/pwstore/internal/api"
	"github.com/jolfzverb/pwstore/internal/contextkey"
	"github.com/jolfzverb/pwstore/internal/dependencies"
	sessioninfoget "github.com/jolfzverb/pwstore/internal/views/session/info/get"
	sessionnewpost "github.com/jolfzverb/pwstore/internal/views/session/new/post"
	sessionsubmitpost "github.com/jolfzverb/pwstore/internal/views/session/submit/post"
)

type Handlers struct{}

func (h Handlers) GetSessionInfo(
	ctx context.Context,
	request api.GetSessionInfoRequestObject,
) (api.GetSessionInfoResponseObject, error) {
	return sessioninfoget.GetSessionInfo(ctx, request)
}

func (h Handlers) PostSessionNew(
	ctx context.Context,
	request api.PostSessionNewRequestObject,
) (api.PostSessionNewResponseObject, error) {
	return sessionnewpost.PostSessionNew(ctx, request)
}

func (h Handlers) PostSessionSubmit(
	ctx context.Context,
	request api.PostSessionSubmitRequestObject,
) (api.PostSessionSubmitResponseObject, error) {
	return sessionsubmitpost.PostSessionSubmit(ctx, request)
}

func logRequestAndResponse(
	f strictnethttp.StrictHTTPHandlerFunc,
	operationID string,
) strictnethttp.StrictHTTPHandlerFunc {
	ff := func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		request interface{},
	) (response interface{}, err error) {
		requestJSON, err := json.Marshal(request)
		if err == nil {
			slog.Info(
				fmt.Sprintf("Start handling %s %s", r.Method, r.URL),
				slog.Any("body", requestJSON),
				slog.String("operation_id", operationID),
			)
		} else {
			slog.Info(fmt.Sprintf("Start handling %s %s", r.Method, r.URL),
				slog.String("operation_id", operationID))
		}

		result, err := f(ctx, w, r, request)

		if err != nil {
			slog.Error(
				fmt.Sprintf("Error processing request %s %s", r.Method, r.URL),
				slog.Any("err", err),
				slog.String("operation_id", operationID),
			)
		} else {
			var s strings.Builder
			err2 := json.NewEncoder(&s).Encode(result)
			if err2 == nil {
				slog.Info(fmt.Sprintf("Finish handling %s %s", r.Method, r.URL),
					slog.Any("body", s.String()),
					slog.String("operation_id", operationID))
			} else {
				slog.Info(fmt.Sprintf("Finish handling %s %s", r.Method, r.URL),
					slog.String("operation_id", operationID))
			}
		}

		return result, err
	}
	return ff
}

func GetHandler() http.Handler {
	handlers := Handlers{}
	return api.Handler(
		api.NewStrictHandler(handlers, []api.StrictMiddlewareFunc{logRequestAndResponse}),
	)
}

func InitializeServer(deps *dependencies.Collection) (*http.Server, error) {
	h := GetHandler()

	server := http.Server{
		Addr:              ":8080",
		Handler:           h,
		ReadHeaderTimeout: 2 * time.Second,
	}
	server.ConnContext = func(ctx context.Context, _ net.Conn) context.Context {
		return context.WithValue(ctx, contextkey.Deps, deps)
	}
	return &server, nil
}
