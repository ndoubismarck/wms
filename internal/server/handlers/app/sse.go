package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/data/entities"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *handler) sse(context *gin.Context) {
	context.Writer.Header().Set("Content-Type", "text/event-stream")
	context.Writer.Header().Set("Cache-Control", "no-cache")
	context.Writer.Header().Set("Connection", "keep-alive")
	context.Writer.Header().Set("X-Accel-Buffering", "no")
	flusher, ok := context.Writer.(http.Flusher)
	if !ok {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	user, ok := h.helpers.Auth().GetUser(context)
	if !ok {
		_, _ = fmt.Fprint(context.Writer, ": unauthorized\n\n")
		flusher.Flush()
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	_, _ = fmt.Fprint(context.Writer, ": connected\n\n")
	flusher.Flush()
	heartbeat := time.NewTicker(20 * time.Second)
	userEventsChan := h.helpers.SSE().Subscribe(user.ID)
	defer func(user *entities.User, heartbeat *time.Ticker) {
		heartbeat.Stop()
		h.helpers.SSE().Unsubscribe(user.ID)
	}(user, heartbeat)
	for {
		select {
		case event := <-userEventsChan:
			if event != nil && user.ID == event.UserID {
				if event.Name == types.ServerEventNameDisconnectClient {
					_, _ = fmt.Fprint(context.Writer, ": bye\n\n")
					flusher.Flush()
					h.helpers.SSE().Unsubscribe(event.UserID)
					return
				}
				data, err := json.Marshal(event.Payload)
				if err != nil {
					_, _ = fmt.Fprint(context.Writer, ": error\n\n")
					flusher.Flush()
					h.ctx.Logger().Error(err)
					return
				}
				if _, err = fmt.Fprintf(context.Writer, "id: %s\nevent: %s\ndata: %s\n\n",
					cryptoutil.UID(), event.Name, string(data)); err != nil {
					h.ctx.Logger().Error(err)
					return
				}
				flusher.Flush()
			}
		case <-heartbeat.C:
			if !h.helpers.SSE().Subscriber(user.ID) {
				_, _ = fmt.Fprint(context.Writer, ": bye\n\n")
				flusher.Flush()
				return
			}
			data, err := json.Marshal(gin.H{
				"time": time.Now().UTC().Format(time.RFC3339),
			})
			if err != nil {
				_, _ = fmt.Fprint(context.Writer, ": error\n\n")
				flusher.Flush()
				h.ctx.Logger().Error(err)
				return
			}
			if _, err = fmt.Fprintf(context.Writer, "id: %s\nevent: %s\ndata: %s\n\n",
				cryptoutil.UID(), "ping", string(data)); err != nil {
				h.ctx.Logger().Error(err)
				return
			}
			flusher.Flush()
		case <-context.Request.Context().Done():
			return
		}
	}
}
