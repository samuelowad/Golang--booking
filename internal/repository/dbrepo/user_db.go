package dbrepo

import "github.com/samuelowad/bookings/internal/models"

//GetUserByID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx := getCtx()

	query := `SELECT id,first_name,last_name,email,access_token,password,created_at,updated_at FROM users WHERE id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)
	var u models.User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.AccessLevel, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, err
	}
	return u, nil
}
