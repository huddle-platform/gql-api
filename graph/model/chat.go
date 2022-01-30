package model

type Chat struct {
	// first participant in the chat
	P1 *ChatParticipant `json:"p1"`
	// second participant in the chat
	P2 *ChatParticipant `json:"p2"`
	// me if logged in user is part of the chat (as user or project owner), null otherwise
	Me MessageAuthor `json:"me"`
}
