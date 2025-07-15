package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	resp "url-shortener/internals/lib/api/response"
	"url-shortener/internals/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "internals.handlers.redirect.New"

		log := log.With(
			"fn", fn,
			"request_id", middleware.GetReqID(r.Context()),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Error("failed to get alias param")

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		url, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url with such alias not found", "err", err)

			render.JSON(w, r, resp.Error("url with such alias not found: "+alias))

			return
		}

		if err != nil {
			log.Error("failed to get url from storage", "err", err)

			render.JSON(w, r, resp.Error("failed to get url with such alias from storage: "+alias))

			return
		}

		log.Info("redirecting", "alias", alias, "url", url)

		http.Redirect(w, r, url, http.StatusFound)
	}
}
