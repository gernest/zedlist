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
		flashMessages.Err(settings.ErrInvalidForm)
		flashMessages.Save(ctx)
		return ctx.Redirect(http.StatusFound, "/jobs/new")
	}
	if !jf.Valid() {
		flashMessages.Err(settings.ErrInvalidForm)
		flashMessages.Save(ctx)
		return ctx.Redirect(http.StatusFound, "/jobs/new")
	}
	person := ctx.Get("User").(*models.Person)
	jb, jerr := query.PersonCreateJob(db.Conn, person, *jf)
	if jerr != nil {
		flashMessages.Err(settings.FlashFailedNewJob)
		flashMessages.Save(ctx)
		return ctx.Redirect(http.StatusFound, "/jobs/new")
	}
	flashMessages.Success(settings.FlashCreateJobSuccess)
	flashMessages.Save(ctx)
	return ctx.Redirect(http.StatusFound, fmt.Sprintf("/jobs/view/%d", jb.ID))
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

func Update(ctx echo.Context) error {
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
		return ctx.Render(http.StatusInternalServerError, tmpl.ErrServerTpl, utils.GetData(ctx))
	}
	utils.SetData(ctx, "PageTitle", "new job")
	utils.SetData(ctx, "Job", job)
	return ctx.Render(http.StatusOK, tmpl.JobsUpdateTpl, utils.GetData(ctx))
}

func UpdatePost(ctx echo.Context) error {
	id, err := utils.GetInt64(ctx.Param("id"))
	if err != nil {
		utils.SetData(ctx, settings.Message, settings.ErrBadRequest)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}
	flashMessages := flash.New()
	f := forms.New(utils.GetLang(ctx))
	jf, err := f.DecodeJobForm(ctx.Request())
	if err != nil {
		flashMessages.Err(settings.ErrInvalidForm)
		flashMessages.Save(ctx)
		return ctx.Redirect(http.StatusFound, fmt.Sprintf("/jobs/update/%d", id))
	}
	if !jf.Valid() {
		flashMessages.Err(settings.ErrInvalidForm)
		flashMessages.Save(ctx)
		return ctx.Redirect(http.StatusFound, fmt.Sprintf("/jobs/update/%d", id))
	}
	job, err := query.GetJobByID(db.Conn, id)
	if err != nil {
		if query.NotFound(err) {
			utils.SetData(ctx, settings.Message, settings.ErrorResourceNotFound)
			return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
		}
		utils.SetData(ctx, settings.Message, settings.ErrorInternalServer)
		return ctx.Render(http.StatusInternalServerError, tmpl.ErrServerTpl, utils.GetData(ctx))
	}
	job.Title = jf.Title
	job.Description = jf.Description
	if err = query.Update(db.Conn, job); err != nil {
		log.Error(ctx, err)
		utils.SetData(ctx, settings.Message, settings.ErrorInternalServer)
		return ctx.Render(http.StatusInternalServerError, tmpl.ErrServerTpl, utils.GetData(ctx))
	}
	flashMessages.Success(settings.FlashSuccessUpdate)
	flashMessages.Save(ctx)
	return ctx.Redirect(http.StatusFound, fmt.Sprintf("/jobs/view/%d", id))
}

func Delete(ctx echo.Context) error {
	id, err := utils.GetInt64(ctx.Param("id"))
	if err != nil {
		utils.SetData(ctx, settings.Message, settings.ErrBadRequest)
		return ctx.Render(http.StatusBadRequest, tmpl.ErrBadRequest, utils.GetData(ctx))
	}
	flashMessages := flash.New()
	job, err := query.GetJobByID(db.Conn, id)
	if err != nil {
		if query.NotFound(err) {
			utils.SetData(ctx, settings.Message, settings.ErrorResourceNotFound)
			return ctx.Render(http.StatusNotFound, tmpl.ErrNotFoundTpl, utils.GetData(ctx))
		}
		utils.SetData(ctx, settings.Message, settings.ErrorInternalServer)
		return ctx.Render(http.StatusInternalServerError, tmpl.ErrServerTpl, utils.GetData(ctx))
	}
	person := ctx.Get("User").(*models.Person)
	if person.ID != job.PersonID {
		//TODO ; render unauthorized template
	}
	if err = query.Delete(db.Conn, job); err != nil {
		log.Error(ctx, err)
		utils.SetData(ctx, settings.Message, settings.ErrorInternalServer)
		return ctx.Render(http.StatusInternalServerError, tmpl.ErrServerTpl, utils.GetData(ctx))
	}
	flashMessages.Success(settings.FlashSuccessDelete)
	flashMessages.Save(ctx)
	return ctx.Redirect(http.StatusFound, "/jobs/")
}
