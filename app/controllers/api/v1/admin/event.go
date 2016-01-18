package admin

import (
	"github.com/gitDashboard/client/v1/admin/request"
	"github.com/gitDashboard/client/v1/admin/response"
	"github.com/gitDashboard/client/v1/misc"
	basicResponse "github.com/gitDashboard/client/v1/response"
	"github.com/gitDashboard/gitDashboard/app/controllers"
	"github.com/gitDashboard/gitDashboard/app/models"
	"github.com/revel/revel"
	"time"
)

type AdminEvent struct {
	controllers.AdminController
}

func (ctrl *AdminEvent) Search() revel.Result {
	var req request.FindEventRequest
	var resp response.FindEventResponse
	var dbEvents []models.Event

	err := ctrl.GetJSONBody(&req)
	if err != nil {
		revel.ERROR.Println(err.Error())
		return ctrl.RenderError(err)
	}

	query := ctrl.Tx.Table("events")
	if req.RepoID != 0 {
		query = query.Where("repo_id = ? ", req.RepoID)
	}
	if req.User != "" {
		query = query.Where("user = ? ", req.User)
	}
	if req.Type != "" {
		query = query.Where("type = ? ", req.Type)
	}
	if len(req.Levels) != 0 {
		query = query.Where("level IN (?) ", req.Levels)
	}
	if req.Description != "" {
		query = query.Where("description like ? ", "%"+req.Description+"%")
	}
	if req.Reference != "" {
		query = query.Where("reference like ? ", "%"+req.Reference+"%")
	}
	if req.Since != 0 {
		dtSince := time.Unix(req.Since, 0)
		query = query.Where("started >= ? ", dtSince)
	}
	if req.To != 0 {
		dtTo := time.Unix(req.To, 0)
		query = query.Where("started <= ? ", dtTo)
	}
	query = query.Order("started desc").Limit(req.Count).Offset(req.First).Find(&dbEvents)
	if query.Error != nil {
		controllers.ErrorResp(&resp, basicResponse.DbError, query.Error)
		return ctrl.RenderJson(resp)
	}

	resp.Events = make([]response.Event, len(dbEvents), len(dbEvents))
	for i, dbEvent := range dbEvents {
		var uxFinished int64
		if dbEvent.Finished != nil {
			uxFinished = dbEvent.Finished.UnixNano() / 1000000
		}
		resp.Events[i] = response.Event{ID: dbEvent.ID,
			Type:        dbEvent.Type,
			User:        dbEvent.User,
			Level:       misc.EventLevel(dbEvent.Level),
			Description: dbEvent.Description,
			Reference:   dbEvent.Reference,
			Started:     dbEvent.Started.UnixNano() / 1000000,
			Finished:    uxFinished}
	}

	resp.Success = true
	return ctrl.RenderJson(resp)
}
