// Provides methods for interfacing with the accounts database table
package database

import (
	"net/http"

	"github.com/cfrank/auth.fun/api/apierror"
)

func InsertAccount(accountId, emailLocal, emailDomain, passwordHash string, emailVerified bool) *apierror.ApiError {
	stmt, stmtError := MyDb.Db.Prepare(`INSERT INTO account(userid,email_local,email_domain,password_hash,email_verified) VALUES (?,?,?,?,?)`)

	if stmtError != nil {
		return apierror.New("Error creating user", http.StatusInternalServerError)
	}

	defer stmt.Close()

	_, resultError := stmt.Exec(accountId, emailLocal, emailDomain, passwordHash, emailVerified)

	if resultError != nil {
		return apierror.New("Error creating user", http.StatusInternalServerError)
	}

	return nil
}

func UniqueEmail(email string) bool {
	return true
}

func UniqueAccountId(id string) bool {
	return true
}
