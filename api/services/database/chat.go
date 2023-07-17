package database

import (
	"time"
)

type ChatSession struct {
	ID        int
	UserID    int
	PublicKey string
	Name      string
	CreateAt  time.Time
	UpdateAt  time.Time
	Deleted   bool
}

type ChatSessionRepository struct{}

func (csr *ChatSessionRepository) Query(sessionID int) ([]ChatSession, error) {
	// query
	query := `SELECT id, user_id, public_key, name, create_at, update_at, deleted FROM chat_sessions WHERE id = ?`
	rows, err := DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []ChatSession
	for rows.Next() {
		cs := ChatSession{}
		err := rows.Scan(&cs.ID, &cs.UserID, &cs.PublicKey, &cs.Name, &cs.CreateAt, &cs.UpdateAt, &cs.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, cs)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (csr *ChatSessionRepository) QueryByUserID(userID int) ([]ChatSession, error) {
	// query
	query := `SELECT id, user_id, public_key, name, create_at, update_at, deleted FROM chat_sessions WHERE user_id = ?`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []ChatSession
	for rows.Next() {
		cs := ChatSession{}
		err := rows.Scan(&cs.ID, &cs.UserID, &cs.PublicKey, &cs.Name, &cs.CreateAt, &cs.UpdateAt, &cs.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, cs)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (csr *ChatSessionRepository) UpdateName(id int, name string) error {
	// query
	query := `UPDATE chat_sessions SET name = ?, update_at = NOW() WHERE id = ?`
	_, err := DB.Exec(query, name, id)
	if err != nil {
		return err
	}
	return nil
}

func (csr *ChatSessionRepository) Create(userID int, publicKey string, name string) error {
	// query
	query := `INSERT INTO chat_sessions (user_id, public_key, name) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, userID, publicKey, name)
	if err != nil {
		return err
	}
	return nil
}

func (csr *ChatSessionRepository) Delete(sessionID int) error {
	// query
	query := `DELETE FROM chat_sessions WHERE id = ?`
	_, err := DB.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

/////////////////////////////
//chat session participants//
/////////////////////////////

type ChatSessionParticipant struct {
	ID        int
	SessionID int
	UserID    int
	Status    ParticipantStatus
	CreateAt  time.Time
	UpdateAt  time.Time
	Deleted   bool
}

type ParticipantStatus string

const (
	Invited  ParticipantStatus = "invited"
	Joined   ParticipantStatus = "joined"
	Rejected ParticipantStatus = "rejected"
)

type ChatSessionParticipantRepository struct{}

func (cspr *ChatSessionParticipantRepository) QueryBySessionID(sessionID int) ([]ChatSessionParticipant, error) {
	// query
	query := `SELECT id, chat_session_id, user_id, status, create_at, update_at, deleted FROM chat_session_participants WHERE chat_session_id = ?`
	rows, err := DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []ChatSessionParticipant
	for rows.Next() {
		csp := ChatSessionParticipant{}
		err := rows.Scan(&csp.ID, &csp.SessionID, &csp.UserID, &csp.Status, &csp.CreateAt, &csp.UpdateAt, &csp.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, csp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (cspr *ChatSessionParticipantRepository) QueryBySessionAndUser(sessionID int, userID int) (*ChatSessionParticipant, error) {
	// query
	query := `SELECT id, chat_session_id, user_id, status, create_at, update_at, deleted FROM chat_session_participants WHERE chat_session_id = ? AND user_id = ?`
	row := DB.QueryRow(query, sessionID, userID)

	// result
	csp := ChatSessionParticipant{}
	err := row.Scan(&csp.ID, &csp.SessionID, &csp.UserID, &csp.Status, &csp.CreateAt, &csp.UpdateAt, &csp.Deleted)
	if err != nil {
		return nil, err
	}

	return &csp, nil
}

func (cspr *ChatSessionParticipantRepository) Create(sessionID int, userID int, status ParticipantStatus) error {
	// query
	query := `INSERT INTO chat_session_participants (chat_session_id, user_id, status) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, sessionID, userID, status)
	if err != nil {
		return err
	}
	return nil
}

func (cspr *ChatSessionParticipantRepository) Delete(sessionID int) error {
	// query
	query := `DELETE FROM chat_session_participants WHERE chat_session_id = ?`
	_, err := DB.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (cspr *ChatSessionParticipantRepository) UpdateStatus(id int, status ParticipantStatus) error {
	// query
	query := `UPDATE chat_session_participants SET status = ? WHERE id = ?`
	_, err := DB.Exec(query, status, id)
	if err != nil {
		return err
	}
	return nil
}

/////////////////////////////
//         chats           //
/////////////////////////////

type Chat struct {
	ID        int
	SessionID int
	UserID    int
	Content   string
	CreateAt  time.Time
	UpdateAt  time.Time
	Deleted   bool
}

type ChatRepository struct{}

func (cr *ChatRepository) QueryBySessionID(sessionID int) ([]Chat, error) {
	// query
	query := `SELECT id, chat_session_id, user_id, content, create_at, update_at, deleted FROM chats WHERE chat_session_id = ?`
	rows, err := DB.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []Chat
	for rows.Next() {
		c := Chat{}
		err := rows.Scan(&c.ID, &c.SessionID, &c.UserID, &c.Content, &c.CreateAt, &c.UpdateAt, &c.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// QueryByUserIDAndTimeRange retrieves chat messages created within a specified time range for a given user ID.
// It takes the userID as the ID of the user to query, endTime as the upper limit of the time range.
// The function returns a slice of Chat structs representing the retrieved messages, or an error if the query fails.
func (cr *ChatRepository) QueryByUserIDAndTimeRange(userID int, endTime time.Time) ([]Chat, error) {
	startTime := time.Now().Add(-24 * time.Hour)

	// query
	query := `
	SELECT id, chat_session_id, user_id, content, create_at, update_at, deleted 
	FROM chats 
	WHERE user_id = ?
		AND create_at >= ? 
		AND create_at <= ?`
	rows, err := DB.Query(query, userID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// results
	var results []Chat
	for rows.Next() {
		c := Chat{}
		err := rows.Scan(&c.ID, &c.SessionID, &c.UserID, &c.Content, &c.CreateAt, &c.UpdateAt, &c.Deleted)
		if err != nil {
			return nil, err
		}
		results = append(results, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (cr *ChatRepository) Create(sessionID int, userID int, content string) error {
	// query
	query := `INSERT INTO chats (chat_session_id, user_id, content) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, sessionID, userID, content)
	if err != nil {
		return err
	}
	return nil
}

func (cr *ChatRepository) Delete(sessionID int) error {
	// query
	query := `DELETE FROM chats WHERE chat_session_id = ?`
	_, err := DB.Exec(query, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (cr *ChatRepository) DeleteBySessionAndUser(sessionID int, userID int) error {
	// query
	query := `DELETE FROM chats WHERE chat_session_id = ? AND user_id = ?`
	_, err := DB.Exec(query, sessionID, userID)
	if err != nil {
		return err
	}
	return nil
}
