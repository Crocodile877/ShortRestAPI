package redirect

import (
	resp "ShortRestAPI/internal/lib/api/response"
	"ShortRestAPI/internal/lib/logger/sl"
	"ShortRestAPI/internal/storage"
	"ShortRestAPI/internal/storage/sqlite"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	slog2 "golang.org/x/exp/slog"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog2.Logger, urlGetter *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op))
		slog.String("request_id", middleware.GetReqID(r.Context()))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resURL, err := urlGetter.GetUrl(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		//redirect to found url (without error)
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
