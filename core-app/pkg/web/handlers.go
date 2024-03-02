package web

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"be/pkg/backend"
	"be/pkg/config"
	"be/pkg/storage"
	"be/pkg/web/types"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id int `form:"id"`
}

type routerGroup interface {
	Handle(httpMethod string, path string, handlers ...gin.HandlerFunc) gin.IRoutes
}

type webHandlers struct {
	ctx     context.Context
	log     *slog.Logger
	conf    *config.Config
	storage storage.Storer
	back    *backend.Backend

	frontend func(routerGroup) (string, string, []gin.HandlerFunc)
	backend  func(routerGroup) (string, string, []gin.HandlerFunc)
}

func (wh *webHandlers) frontendHandlers(rg routerGroup) func() {
	return func() {
		rg.Handle(wh.healthzHandler())
		rg.Handle(wh.indexHandler())
	}
}

func (wh *webHandlers) backendHandlers(rg routerGroup) func() {
	return func() {
		rg.Handle(wh.healthzHandler())
		rg.Handle(wh.getAboutMeHandler())
		rg.Handle(wh.getTitleHandler())
	}
}

func (wh *webHandlers) healthzHandler() (string, string, gin.HandlerFunc) {
	return http.MethodGet, "/healthz", func(c *gin.Context) {
		_, err := wh.storage.GetTitle(1)
		if err != nil {
			wh.log.Error("Database unresponsive", slog.String("err", err.Error()))
			c.AbortWithStatusJSON(types.InternalErrorResponse(err))
			return
		}

		c.JSON(types.SuccessResponse("all systems up and running"))
	}
}

func (wh *webHandlers) indexHandler() (string, string, gin.HandlerFunc) {
	return http.MethodGet, "/", func(c *gin.Context) {
		var u User
		if err := c.ShouldBindQuery(&u); err != nil {
			// Mocking for now
			u.Id = 1
		}

		if u.Id == 0 {
			// Mocking for now
			u.Id = 1
		}

		uc, err := wh.back.GetUserContent(u.Id)
		if err != nil {
			wh.log.Error("Error from backend", slog.String("err", err.Error()))
			c.AbortWithStatusJSON(types.InternalErrorResponse(err))
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   uc.Title,
			"aboutme": uc.AboutMe,
		})
	}
}

func (wh *webHandlers) getTitleHandler() (string, string, gin.HandlerFunc) {
	return http.MethodGet, "/title", func(c *gin.Context) {
		var u User
		if err := c.ShouldBindQuery(&u); err != nil {
			wh.log.Error("Could not bind query", slog.String("err", err.Error()))
			c.AbortWithStatusJSON(types.InternalErrorResponse(err))
			return
		}

		title, err := wh.storage.GetTitle(int32(u.Id))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(types.SuccessResponse("user id doesn't exist"))
				return
			}

			wh.log.Error("Could not get title from DB", slog.String("err", err.Error()))
			c.AbortWithStatusJSON(types.InternalErrorResponse(err))
			return
		}

		c.JSON(types.SuccessResponse(title))
	}
}

func (wh *webHandlers) getAboutMeHandler() (string, string, gin.HandlerFunc) {
	return http.MethodGet, "/aboutme", func(c *gin.Context) {
		var u User
		if err := c.ShouldBindQuery(&u); err != nil {
			wh.log.Error("Could not bind query", slog.String("err", err.Error()))
			c.AbortWithStatusJSON(types.InternalErrorResponse(err))
			return
		}

		aboutme, err := wh.storage.GetAboutMe(int32(u.Id))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(types.SuccessResponse("user id doesn't exist"))
				return
			}

			wh.log.Error("Could not get aboutme from DB", slog.String("err", err.Error()))
			c.AbortWithStatusJSON(types.InternalErrorResponse(err))
			return
		}

		c.JSON(types.SuccessResponse(aboutme))
	}
}
