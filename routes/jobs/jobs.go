package jobs

import (
	"net/http"

	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/db"
	"github.com/gernest/zedlist/modules/flash"
	"github.com/gernest/zedlist/modules/forms"
	"github.com/gernest/zedlist/modules/query"
	"github.com/gernest/zedlist/modules/tmpl"
	"github.com/gernest/zedlist/modules/utils"
	"github.com/labstack/echo"
)

var msgInvalidorm = "some fish happende"

func New(ctx echo.Context) error {
	utils.SetData(ctx, "PageTitle", "new job")
	return ctx.Render(http.StatusOK, tmpl.JobsNewTpl, utils.GetData(ctx))
}

func NewPost(ctx echo.Context) error {
	var flashMessages = flash.New()
	f := forms.New(utils.GetLang(ctx))
	jf, err := f.DecodeJobForm(ctx.Request())
	if err != nil {
		// TODO: improve flash message ?
		flashMessages.Err(msgInvalidorm)
		flashMessages.Save(ctx)
		ctx.Redirect(http.StatusFound, "/jobs/new")
		return nil
	}
	if !jf.Valid() {
		// TODO: improve flash message ?
		flashMessages.Err(msgInvalidorm)
		flashMessages.Save(ctx)
		ctx.Redirect(http.StatusFound, "/jobs/new")
		return nil
	}

	if isLoged := ctx.Get("IsLoged"); isLoged != nil {
		person := ctx.Get("User").(*models.Person)
		if jerr := query.PersonCreateJob(db.Conn, person, *jf); jerr != nil {
			// TODO: improve flash message ?
			flashMessages.Err("some really bad fish happened")
			flashMessages.Save(ctx)
			ctx.Redirect(http.StatusFound, "/dash/jobs/new")
			return nil
		}
		// add flash message
		flashMessages.Success("new job was created successful")
		flashMessages.Save(ctx)

		ctx.Redirect(http.StatusFound, "/dash/")
		return nil
	}
	he := echo.NewHTTPError(http.StatusUnauthorized)
	ctx.Error(he)
	return nil
}
