package transporthttp

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jithin-kg/webpa-common/logging"
	"github.com/jithin-kg/webpa-common/xhttp"
)

// GetBody returns a go-kit RequestFunc that sets the request's GetBody function.  This ensures
// that redirects are properly followed automatically.
func GetBody(logger log.Logger) func(context.Context, *http.Request) context.Context {
	return func(ctx context.Context, request *http.Request) context.Context {
		err := xhttp.EnsureRewindable(request)
		if err != nil {
			logger.Log(level.Key(), level.ErrorValue(), logging.MessageKey(), "Unable to setup request for rewind", logging.ErrorKey(), err)
		}

		return ctx
	}
}
