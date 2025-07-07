package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/rugi123/go-shortener/internal/domain/model"
	"github.com/rugi123/go-shortener/internal/domain/service"
	"github.com/rugi123/go-shortener/internal/storage/postgres"
)

func Shortener(original_url string, r *gin.Engine, conn *pgx.Conn, ctx context.Context) (string, error) {
	short_url, err := service.GenerateShortURL()
	if err != nil {
		return "", fmt.Errorf("ошибка гинерации алиаса: %s", err)
	}
	link := model.Link{
		OriginalUrl: original_url,
		ShortUrl:    short_url,
	}
	err = postgres.SaveLink(link, ctx, *conn)
	if err != nil {
		return "", fmt.Errorf("ошибка сохранения link: %s", err)
	}
	r.GET(short_url, func(c *gin.Context) {
		path := c.FullPath()
		link, err := postgres.SearchLink(path, ctx, *conn)
		if err != nil {
			fmt.Println("ошибка поиска в db: ", err)
		}
		if link != nil {
			c.Redirect(http.StatusPermanentRedirect, link.OriginalUrl)
		}
	})
	return short_url, err
}
