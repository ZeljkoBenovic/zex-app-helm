package main

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type handlerFn func(*slog.Logger, *db, *Config) (string, string, gin.HandlerFunc)

var handlersBE = []handlerFn{
	healthzHandler,
	getTitleHandler,
	getAboutMeHandler,
}

var handlersFE = []handlerFn{
	indexHandler,
	healthzHandler,
}

type User struct {
	Id int `form:"id"`
}

type backendData struct {
	title   string
	aboutme string
}

type Response[T any] struct {
	Success bool   `json:"status"`
	Error   string `json:"error"`
	Data    T      `json:"data,omitempty"`
}

func SuccessResponse(data any) (int, Response[any]) {
	return http.StatusOK, Response[any]{
		Success: true,
		Data:    data,
	}
}

func InternalErrorResponse(err error) (int, Response[string]) {
	return http.StatusInternalServerError, Response[string]{
		Success: false,
		Error:   err.Error(),
	}
}

func indexHandler(log *slog.Logger, _ *db, conf *Config) (string, string, gin.HandlerFunc) {
	return http.MethodGet, "/", func(c *gin.Context) {
		var u User
		if err := c.ShouldBindQuery(&u); err != nil {
			u.Id = 1
		}

		if u.Id == 0 {
			u.Id = 1
		}

		beData, err := fetchBEData(os.Getenv("BE_URL"), u.Id, c, log)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   beData.title,
			"aboutme": beData.aboutme,
		})
	}
}

func fetchBEData(beUrl string, id int, c *gin.Context, log *slog.Logger) (backendData, error) {
	var (
		title Response[string]
		about Response[string]
	)

	// fetch title
	rawTitle, err := http.Get(fmt.Sprintf("%s/title?id=%d", beUrl, id))
	if err != nil {
		log.Error("Error fetching title from backend", "err", err)
		//TODO: create error UI error page
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return backendData{}, err
	}

	// fetch aboutme
	rawAbout, err := http.Get(fmt.Sprintf("%s/aboutme?id=%d", beUrl, id))
	if err != nil {
		log.Error("Error fetching title from backend", "err", err)
		//TODO: create error UI error page
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return backendData{}, err
	}

	bTitle, _ := io.ReadAll(rawTitle.Body)
	json.Unmarshal(bTitle, &title)

	bAbout, _ := io.ReadAll(rawAbout.Body)
	json.Unmarshal(bAbout, &about)

	return backendData{
		title:   title.Data,
		aboutme: about.Data,
	}, nil

}

func healthzHandler(log *slog.Logger, d *db, _ *Config) (string, string, gin.HandlerFunc) {
	return http.MethodGet, "/healthz", func(c *gin.Context) {
		if d != nil {
			_, err := d.db.GetTitle(context.Background(), 1)
			if err != nil {
				log.Error("Database unresponsive", slog.String("err", err.Error()))
				c.AbortWithStatusJSON(InternalErrorResponse(err))
				return
			}
		}

		c.JSON(SuccessResponse("all systems up and running"))
	}
}

func getTitleHandler(log *slog.Logger, d *db, _ *Config) (string, string, gin.HandlerFunc) {
	return http.MethodGet, "/title", func(c *gin.Context) {
		var u User
		if err := c.ShouldBindQuery(&u); err != nil {
			log.Error("Could not bind query", slog.String("err", err.Error()))
			c.AbortWithStatusJSON(InternalErrorResponse(err))
			return
		}

		title, err := d.db.GetTitle(context.Background(), int32(u.Id))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(SuccessResponse("user id doesn't exist"))
				return
			}

			log.Error("Could not get title from DB", slog.String("err", err.Error()))
			c.AbortWithStatusJSON(InternalErrorResponse(err))
			return
		}

		if title.Valid {
			c.JSON(SuccessResponse(title.String))
		} else {
			c.JSON(SuccessResponse("There Is No Title Yet!"))
		}
	}
}

func getAboutMeHandler(log *slog.Logger, d *db, _ *Config) (string, string, gin.HandlerFunc) {
	return http.MethodGet, "/aboutme", func(c *gin.Context) {
		var u User
		if err := c.ShouldBindQuery(&u); err != nil {
			log.Error("Could not bind query", slog.String("err", err.Error()))
			c.AbortWithStatusJSON(InternalErrorResponse(err))
			return
		}

		abm, err := d.db.GetAboutMe(context.Background(), int32(u.Id))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(SuccessResponse("user id doesn't exist"))
				return
			}

			log.Error("Could not get about me from DB", slog.String("err", err.Error()))
			c.AbortWithStatusJSON(InternalErrorResponse(err))
			return
		}

		if abm.Valid {
			c.JSON(SuccessResponse(abm.String))
		} else {
			c.JSON(SuccessResponse("There is nothing about me yet!"))
		}
	}
}
