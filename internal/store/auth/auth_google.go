package store_auth

import (
	"git/ymoldabe/forum/models"
)

func (a *Auth) GoogleAuthUser(form *models.GoogleLoginUserData) (int, error) {
	ok, err := a.checkUserInTable(form.Email)
	if err != nil {
		return 0, err
	}

	if ok {
		return a.InsertGoogleUser(form)
	} else {

		if err := a.UpdatePassGoogleUser(form); err != nil {
			return 0, err
		}

		id, err := a.GetIdGoogleUser(form.Email)

		if err != nil {
			return 0, err
		}
		return id, nil
	}
}

func (a *Auth) GetIdGoogleUser(email string) (int, error) {
	var id int
	stmt := `
	SELECT 
		id
	FROM
		users
	WHERE
		email = ?`
	if err := a.db.QueryRow(stmt, email).Scan(&id); err != nil {
		return id, err
	}
	return id, nil

}

func (a *Auth) UpdatePassGoogleUser(form *models.GoogleLoginUserData) error {
	stmt := `
	UPDATE 
		users
	SET 
		name = ?,
		hash_password = ?
	WHERE 
		email = ?`
	_, err := a.db.Exec(stmt, form.Name, form.Password, form.Email)
	return err
}

func (a *Auth) InsertGoogleUser(form *models.GoogleLoginUserData) (int, error) {
	stmt := `
	INSERT INTO 
		users(name, email, hash_password, provider)
	VALUES
		(?, ?, ?, ?)`
	res, err := a.db.Exec(stmt, form.Name, form.Email, form.Password, models.ProviderGoogle)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err

}

func (a *Auth) checkUserInTable(email string) (bool, error) {
	var checkUser int
	stmt := `
	SELECT COUNT(*) FROM users WHERE email = ?;`

	if err := a.db.QueryRow(stmt, email).Scan(&checkUser); err != nil {
		return false, err
	}
	if checkUser > 0 {
		return false, nil
	}
	return true, nil
}
