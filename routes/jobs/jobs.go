package jobs

import (
	"fmt"
	"net/http"

	"github.com/gernest/zedlist/modules/settings"

	"github.com/gernest/zedlist/models"
	"github.com/gernest/zedlist/modules/db"
	"github.com/gernest/zedlist/modules/flash"
	"github.com/gernest/zedlist/modules/forms"
	"github.com/gernest/zedlist/modules/log"
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
		jb, jerr := query.PersonCreateJob(db.Conn, person, *jf)
		if jerr != nil {
			// TODO: improve flash message ?
			flashMessages.Err("some really bad fish happened")
			flashMessages.Save(ctx)
			ctx.Redirect(http.StatusFound, "/jobs/new")
			return nil
		}
		// add flash message
		flashMessages.Success("new job was created successful")
		flashMessages.Save(ctx)
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/jobs/view/%d", jb.ID))
		return nil
	}
	he := echo.NewHTTPError(http.StatusUnauthorized)
	ctx.Error(he)
	return nil
}

func View(ctx echo.Context) error {
	id, err := utils.GetInt64(ctx.Param("id"))
	if err != nil {
		utils.SetData(ctx, settings.Message, settings.ErrBadRequest)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}
	job, err := query.GetJobByID(db.Conn, id)
	if err != nil {
		if query.NotFound(err) {
			utils.SetData(ctx, settings.Message, settings.ErrorResourceNotFound)
			return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
		}
		utils.SetData(ctx, settings.Message, settings.ErrorInternalServer)
		return ctx.Render(http.StatusNotFound, tmpl.ErrServerTpl, utils.GetData(ctx))
	}
	if job != nil {
		utils.SetData(ctx, "Job", job)
		utils.SetData(ctx, "PageTitle", job.Title)
		owner, err := query.GetPersonByID(db.Conn, job.PersonID)
		if err != nil {
			log.Error(ctx, err)
			utils.SetData(ctx, settings.Message, settings.ErrorInternalServer)
			return ctx.Render(http.StatusNotFound, tmpl.ErrServerTpl, utils.GetData(ctx))
		}
		utils.SetData(ctx, "Owner", owner)
	}
	return ctx.Render(http.StatusOK, tmpl.JobsViewTpl, utils.GetData(ctx))
}

func List(ctx echo.Context) error {
	jobs, err := query.GetLatestJobs(db.Conn)
	if err != nil {
		utils.SetData(ctx, settings.Message, settings.ErrorResourceNotFound)
		return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
	}
	utils.SetData(ctx, "Jobs", jobs)
	utils.SetData(ctx, settings.PageTitle, "jobs")
	return ctx.Render(http.StatusOK, tmpl.JobsListTpl, utils.GetData(ctx))
}
