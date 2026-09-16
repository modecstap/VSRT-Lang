package session_repository


func (r *Repository) Delete(sessionId int) error {
	_, err := r.db.Exec(`
		DELETE FROM sessions
		WHERE id = $1
	`, sessionId)
	return err
}
