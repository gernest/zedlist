package jobs

import (
	"net/http"

	"github.com/gernest/zedlist/modules/tmpl"
	"github.com/gernest/zedlist/modules/utils"
	"github.com/labstack/echo"
)

func New(ctx echo.Context) error {
	utils.SetData(ctx, "PageTitle", "new job")
	return ctx.Render(http.StatusOK, tmpl.JobsNewTpl, utils.GetData(ctx))
}
