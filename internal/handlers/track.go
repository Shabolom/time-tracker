package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"sort"
	"strings"
	"time"
	"timeTracker/internal/model"
)

// StartTrackTask godoc
// @Summary		Track a time for task
// @Tags		Track
// @Produce		json
// @Accept		json
// @Param	id	body		model.TaskTrackRequest	true	"task id"
// @Success		200	{object} model.TaskTrack
// @Router			/api/start-track [post]
func (h *Handlers) StartTrackTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var trackModel model.TaskTrackRequest

	err := json.NewDecoder(r.Body).Decode(&trackModel)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(trackModel)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	trackId := uuid.New()

	trackTimeEntity := model.TaskTrack{
		TaskID: trackModel.TaskID,
		Task: model.Task{
			Base: model.Base{ID: trackModel.TaskID},
		},
		Base: model.Base{ID: trackId},
	}

	trackTimeResponse, err := h.Storage.TrackTime(ctx, trackTimeEntity)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Sender.JSON(w, http.StatusOK, trackTimeResponse)
	if err != nil {
		panic(err)
	}
}

// EndTrackTask godoc
// @Summary		Stop track a time for task
// @Tags		Track
// @Produce		json
// @Accept		json
// @Param	id	body		model.TaskTrackRequest	true	"task id"
// @Success		200	{object} model.User
// @Router			/api/end-track [post]
func (h *Handlers) EndTrackTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var trackModel model.TaskTrackRequest

	err := json.NewDecoder(r.Body).Decode(&trackModel)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(trackModel)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	trackTimeResponse, err := h.Storage.StopTrackTime(ctx, trackModel.TaskID)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Sender.JSON(w, http.StatusOK, trackTimeResponse.Time)
	if err != nil {
		panic(err)
	}
}

// CalcTime godoc
// @Summary		Calculate a time spent on task
// @Tags			Track
// @Produce		json
// @Accept			json
// @Param	id	body		model.CalcTimeRequest	true	"user id"
// @Success		200	"object,object"
// @Router			/api/calc-time [post]
func (h *Handlers) CalcTime(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user model.CalcTimeRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(user)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	storageResult, err := h.Storage.CalcTime(ctx, user.ID)

	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	sort.Slice(storageResult[:], func(i, j int) bool {
		return *storageResult[i].Time > *storageResult[j].Time
	})

	mapDuration := make(map[string]time.Duration)

	for _, el := range storageResult {
		mapDuration[el.TaskID.String()] = *el.Time
	}

	result := make(map[string]string)

	for k, v := range mapDuration {
		out := time.Time{}.Add(v)

		result[k] = fmt.Sprintf(out.Format("15:04"))
	}

	err = h.Sender.JSON(w, http.StatusOK, result)
	if err != nil {
		panic(err)
	}
}
