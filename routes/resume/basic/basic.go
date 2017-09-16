package basic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/zedio/zedlist/models"
	"github.com/zedio/zedlist/modules/db"
	"github.com/zedio/zedlist/modules/log"
	"github.com/zedio/zedlist/modules/query"
	"github.com/zedio/zedlist/modules/utils"
)

func Get(ctx echo.Context) error {
	id, err := utils.GetInt64(ctx.Param("id"))
	if err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	b, err := query.GetBasicResumeByID(db.Conn, id)
	if err != nil {
		log.Error(ctx, err)
		if query.NotFound(err) {
			return ctx.JSON(http.StatusNotFound, models.NewJSONErr(http.StatusText(
				http.StatusNotFound,
			)))
		}
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(http.StatusText(
			http.StatusInternalServerError,
		)))
	}
	return ctx.JSONPretty(http.StatusOK, b, "\t")
}

func Put(ctx echo.Context) error {
	r := ctx.Request()
	c := r.Header.Get(echo.HeaderContentType)
	if c != echo.MIMEApplicationJSON {
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	b := &models.Basic{}
	o, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	err = json.Unmarshal(o, b)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, models.NewJSONErr(http.StatusText(
			http.StatusUnprocessableEntity,
		)))
	}
	if err = query.Update(db.Conn, b); err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(http.StatusText(
			http.StatusInternalServerError,
		)))
	}
	return ctx.JSONPretty(http.StatusOK, b, "\t")
}

func Post(ctx echo.Context) error {
	r := ctx.Request()
	c := r.Header.Get(echo.HeaderContentType)
	if c != echo.MIMEApplicationJSON {
		ctx.Echo().DefaultHTTPErrorHandler(echo.ErrUnsupportedMediaType, ctx)
		return nil
	}
	b := &models.Basic{}
	o, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	err = json.Unmarshal(o, b)
	if err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusBadRequest, models.NewJSONErr(http.StatusText(
			http.StatusBadRequest,
		)))
	}
	if err = query.Create(db.Conn, b); err != nil {
		log.Error(ctx, err)
		return ctx.JSON(http.StatusInternalServerError, models.NewJSONErr(http.StatusText(
			http.StatusInternalServerError,
		)))
	}
	return ctx.JSONPretty(http.StatusCreated, b, "\t")
}
