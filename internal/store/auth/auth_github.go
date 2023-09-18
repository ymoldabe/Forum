package store_auth

import "git/ymoldabe/forum/models"

func (a *Auth) GithubAuthUser(form *models.GithubLoginUserData) (int, error) {
	ok, err := a.checkUserInTable(form.UserName)
	if err != nil {
		return 0, err
	}

	if ok {
		return a.InsertGithubUser(form)
	} else {

		if err := a.UpdatePassGithubUser(form); err != nil {
			return 0, err
		}

		id, err := a.GetIdGithubUser(form.UserName)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
}

func (a *Auth) GetIdGithubUser(userName string) (int, error) {
	var id int
	stmt := `
	SELECT 
		id
	FROM
		users
	WHERE
		email = ?`
	if err := a.db.QueryRow(stmt, userName).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (a *Auth) UpdatePassGithubUser(form *models.GithubLoginUserData) error {
	stmt := `
	UPDATE 
		users
	SET 
		name = ?,
		hash_password = ?
	WHERE 
		email = ?`
	_, err := a.db.Exec(stmt, form.Name, form.Password, form.UserName)
	return err
}

func (a *Auth) InsertGithubUser(form *models.GithubLoginUserData) (int, error) {
	stmt := `
	INSERT INTO 
		users(name, email, hash_password, provider)
	VALUES
		(?, ?, ?, ?)`
	res, err := a.db.Exec(stmt, form.Name, form.UserName, form.Password, models.ProviderGit)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}
