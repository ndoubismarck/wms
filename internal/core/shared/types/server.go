package types

type ApiResponseBody struct {
	Code       ServiceResultCode `json:"code"`
	Data       any               `json:"data,omitempty"`
	Validation ValidationResult  `json:"validation,omitempty"`
	Pagination *PaginationResult `json:"pagination,omitempty"`
}

type ServerEvent struct {
	Name    ServerEventName `json:"name"`
	UserID  string          `json:"user_id"`
	Payload any             `json:"payload"`
}

type ServerEventName string

type IServer interface {
	SendEvent(event *ServerEvent)
}

const (
	ServerEventNameDisconnectClient ServerEventName = "client.disconnect"
)

func NewServerEvent(name ServerEventName, userId string, payload any) *ServerEvent {
	return &ServerEvent{
		Name:    name,
		UserID:  userId,
		Payload: payload,
	}
}
