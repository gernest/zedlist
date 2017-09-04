package jobs

import (
	"net/http"

	"github.com/gernest/zedlist/modules/forms"
	"github.com/gernest/zedlist/modules/tmpl"
	"github.com/gernest/zedlist/modules/utils"
	"github.com/labstack/echo"
)

func New(ctx echo.Context) error {
	f := forms.New(utils.GetLang(ctx))
	utils.SetData(ctx, "PageTitle", "new job")
	utils.SetData(ctx, "form", f)
	return ctx.Render(http.StatusOK, tmpl.JobsNewTpl, utils.GetData(ctx))
}
