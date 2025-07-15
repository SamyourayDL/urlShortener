package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internals/lib/api/response"
	"url-shortener/internals/lib/random"
	"url-shortener/internals/storage/sqlite"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
	GetURL(alias string) (string, error)
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "internals.handlers.url.save.New"

		log := log.With(
			"fn", fn,
			"request_id", middleware.GetReqID(r.Context()),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", "err", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request decoded successfully", "request", req)

		if err := validator.New().Struct(req); err != nil {
			validatorErrors := err.(validator.ValidationErrors)

			log.Error("req validation failed", "err", err)

			render.JSON(w, r, resp.ValidationError(validatorErrors))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
			for url, _ := urlSaver.GetURL(alias); url != ""; {
				alias = random.NewRandomString(aliasLength)
			}
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, sqlite.StorageErrAlreadyExists{}) {
			log.Info("url already exists", "url", req.URL)

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", "url", "err", err)

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", "id", id)

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
