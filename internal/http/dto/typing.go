package dto

type TypingOperation struct {
	Position *int   `json:"position,omitempty" binding:"omitempty,gte=0,lte=4096"                   example:"5"`
	Length   *int   `json:"length,omitempty"   binding:"omitempty,gte=1,lte=4096"                   example:"3"`
	Type     string `json:"type"               binding:"required,oneof=insert delete replace clear" example:"insert"`
	Text     string `json:"text,omitempty"     binding:"omitempty,min=1,max=4096"                   example:"привет"`
}

func (op *TypingOperation) Validate() bool {
	switch op.Type {
	case "insert":
		return op.Position != nil && op.Text != "" && op.Length == nil
	case "delete":
		return op.Position != nil && op.Text == "" && op.Length != nil
	case "replace":
		return op.Position != nil && op.Text != "" && op.Length == nil
	case "clear":
		return op.Position == nil && op.Text == "" && op.Length == nil
	default:
		return false
	}
}

type TypingUpdateRequest struct {
	Operations []TypingOperation `json:"operations" binding:"required,min=1,max=50"`
}

func (dto *TypingUpdateRequest) Validate() bool {
	for _, op := range dto.Operations {
		if !op.Validate() {
			return false
		}
	}

	return true
}

type TypingEventResponse struct {
	Operations []TypingOperation `json:"operations"`
	SenderID   int64             `json:"sender_id"  example:"123"`
	ChatID     int64             `json:"chat_id"    example:"456"`
}

func (dto *TypingEventResponse) Fill(senderID, chatID int64, operations []TypingOperation) {
	dto.SenderID = senderID
	dto.ChatID = chatID
	dto.Operations = operations
}
