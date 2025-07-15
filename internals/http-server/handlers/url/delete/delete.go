package delete

import (
	"log/slog"
	"net/http"
	"url-shortener/internals/lib/api/response"
	resp "url-shortener/internals/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) (int64, error)
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "internals.handlers.url.delete.New"

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

		_, err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Error("failed to delete: ", "alias", alias, "err", err)

			render.JSON(w, r, response.Error("failed to delete: "+alias))

			return
		}

		log.Info("alias was deleted or abscent", "alias", alias)

		render.JSON(w, r, resp.OK())
	}
}
