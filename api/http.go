package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/the-code-innovator/squeezer/serializer/json"
	"github.com/the-code-innovator/squeezer/serializer/msgpack"
	"github.com/the-code-innovator/squeezer/shortener"
)

type ShortLinkHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	shortLinkService shortener.ShortLinkService
}

func NewHandler(shortLinkService shortener.ShortLinkService) ShortLinkHandler {
	return &handler{shortLinkService: shortLinkService}
}

func SetupResponse(writer http.ResponseWriter, contentType string, body []byte, statusCode int) {
	writer.Header().Set("Content-Type", contentType)
	writer.WriteHeader(statusCode)
	_, err := writer.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) Serializer(contentType string) shortener.ShortLinkSerializer {
	if contentType == "application/x-msgpack" {
		return &msgpack.ShortLink{}
	} else if contentType == "application/json" {
		return &json.ShortLink{}
	} else {
		return nil
	}
}

func (h *handler) Get(writer http.ResponseWriter, request *http.Request) {
	code := chi.URLParam(request, "code")
	shortLink, err := h.shortLinkService.Find(code)
	if err != nil {
		if errors.Cause(err) == shortener.ErrorShortLinkNotFound {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, shortLink.URL, http.StatusMovedPermanently)
}

func (h *handler) Post(writer http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	shortLink, err := h.Serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = h.shortLinkService.Store(shortLink)
	if err != nil {
		if errors.Cause(err) == shortener.ErrorShortLinkInvalid {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseBody, err := h.Serializer(contentType).Encode(shortLink)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	SetupResponse(writer, contentType, responseBody, http.StatusCreated)
}
