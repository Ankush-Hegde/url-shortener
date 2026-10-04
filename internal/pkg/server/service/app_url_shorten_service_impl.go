package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	createshorturl "url-shortener/internal/app/createshorturl"
	resolveshorturl "url-shortener/internal/app/resolveshorturl"
	"url-shortener/internal/pkg/database/mongodb"
	"url-shortener/internal/pkg/redis"
	"url-shortener/internal/pkg/server/gen"
)

type URLStore interface {
	createshorturl.MappingStore
	resolveshorturl.MappingStore
}

type URLCache interface {
	createshorturl.URLCache
	resolveshorturl.URLCache
}

type URLShortenAPIServiceImpl struct {
	store   URLStore
	cache   URLCache
	baseURL string
}

func NewUrlShortenAPIServiceImpl(store URLStore, cache URLCache, baseURL string) gen.DefaultAPIServicer {
	return &URLShortenAPIServiceImpl{
		store:   store,
		cache:   cache,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (s *URLShortenAPIServiceImpl) ShortenPost(ctx context.Context, request gen.ShortenPostRequest) (gen.ImplResponse, error) {
	result, err := createshorturl.CreateShortURL(ctx, request.LongUrl, s.baseURL, s.store, s.cache)
	if err != nil {
		return gen.Response(http.StatusInternalServerError, nil), err
	}
	return gen.Response(http.StatusCreated, gen.ShortenPost201Response{
		ShortCode: result.ShortCode,
		ShortUrl:  result.ShortURL,
	}), nil
}

func (s *URLShortenAPIServiceImpl) ShortCodeGet(ctx context.Context, shortCode string) (gen.ImplResponse, error) {
	response, err := resolveshorturl.ResolveShortURL(ctx, shortCode, s.store, s.cache)
	if err != nil {
		if errors.Is(err, resolveshorturl.ErrShortURLNotFound) {
			return gen.Response(http.StatusNotFound, nil), err
		}
		return gen.Response(http.StatusInternalServerError, nil), err
	}
	return gen.Response(http.StatusFound, response), nil
}

var _ URLStore = (*mongodb.Client)(nil)
var _ URLCache = (*redis.Client)(nil)
