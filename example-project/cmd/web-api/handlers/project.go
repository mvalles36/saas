package handlers

import (
	"context"
	"net/http"
	"strconv"

	"geeks-accelerator/oss/saas-starter-kit/example-project/internal/platform/auth"
	"geeks-accelerator/oss/saas-starter-kit/example-project/internal/platform/web"
	"geeks-accelerator/oss/saas-starter-kit/example-project/internal/project"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Project represents the Project API method handler set.
type Project struct {
	MasterDB *sqlx.DB

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED.
}

// List returns all the existing projects in the system.
func (p *Project) Find(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var req project.ProjectFindRequest
	if err := web.Decode(r, &req); err != nil {
		return errors.Wrap(err, "")
	}

	res, err := project.Find(ctx, claims, p.MasterDB, req)
	if err != nil {
		return err
	}

	return web.RespondJson(ctx, w, res, http.StatusOK)
}

// Read returns the specified project from the system.
func (p *Project) Read(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var includeArchived bool
	if qv := r.URL.Query().Get("include-archived"); qv != "" {
		var err error
		includeArchived, err = strconv.ParseBool(qv)
		if err != nil {
			return errors.Wrapf(err, "Invalid value for include-archived : %s", qv)
		}
	}

	res, err := project.Read(ctx, claims, p.MasterDB, params["id"], includeArchived)
	if err != nil {
		switch err {
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.RespondJson(ctx, w, res, http.StatusOK)
}

// Create inserts a new project into the system.
func (p *Project) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var req project.ProjectCreateRequest
	if err := web.Decode(r, &req); err != nil {
		return errors.Wrap(err, "")
	}

	res, err := project.Create(ctx, claims, p.MasterDB, req, v.Now)
	if err != nil {
		switch err {
		case project.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "Project: %+v", &req)
		}
	}

	return web.RespondJson(ctx, w, res, http.StatusCreated)
}

// Update updates the specified project in the system.
func (p *Project) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	var req project.ProjectUpdateRequest
	if err := web.Decode(r, &req); err != nil {
		return errors.Wrap(err, "")
	}
	req.ID = params["id"]

	err := project.Update(ctx, claims, p.MasterDB, req, v.Now)
	if err != nil {
		switch err {
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s Update: %+v", params["id"], req)
		}
	}

	return web.RespondJson(ctx, w, nil, http.StatusNoContent)
}

// Archive soft-deletes the specified project from the system.
func (p *Project) Archive(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	err := project.Archive(ctx, claims, p.MasterDB, params["id"], v.Now)
	if err != nil {
		switch err {
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "Id: %s", params["id"])
		}
	}

	return web.RespondJson(ctx, w, nil, http.StatusNoContent)
}

// Delete removes the specified project from the system.
func (p *Project) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	err := project.Delete(ctx, claims, p.MasterDB, params["id"])
	if err != nil {
		switch err {
		case project.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case project.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case project.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "Id: %s", params["id"])
		}
	}

	return web.RespondJson(ctx, w, nil, http.StatusNoContent)
}
