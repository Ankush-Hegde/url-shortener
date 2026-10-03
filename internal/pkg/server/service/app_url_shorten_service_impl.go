package service

import (
	"context"
	"net/http"

	createshorturl "url-shortener/internal/app/createshorturl"
	resolveshorturl "url-shortener/internal/app/resolveshorturl"
	"url-shortener/internal/pkg/server/gen"
)

type URLShortenAPIServiceImpl struct{}

func NewUrlShortenAPIServiceImpl() gen.DefaultAPIServicer {
	return &URLShortenAPIServiceImpl{}
}

func (s *URLShortenAPIServiceImpl) ShortenPost(ctx context.Context, request gen.ShortenPostRequest) (gen.ImplResponse, error) {
	response, err := createshorturl.CreateShortURL()
	if err != nil {
		return gen.Response(http.StatusInternalServerError, nil), err
	}
	return gen.Response(http.StatusCreated, response), nil
}

func (s *URLShortenAPIServiceImpl) ShortCodeGet(ctx context.Context, shortCode string) (gen.ImplResponse, error) {
	response, err := resolveshorturl.ResolveShortURL(shortCode)
	if err != nil {
		return gen.Response(http.StatusNotFound, nil), err
	}
	return gen.Response(http.StatusFound, response), nil
}
