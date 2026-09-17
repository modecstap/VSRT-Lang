package session_repository

func (r *Repository) Delete(sessionId int64) error {
	_, err := r.db.Exec(`
		DELETE FROM sessions
		WHERE id = $1
	`, sessionId)
	return err
}
