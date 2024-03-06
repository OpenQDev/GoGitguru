// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.bulkInsertCommitsStmt, err = db.PrepareContext(ctx, bulkInsertCommits); err != nil {
		return nil, fmt.Errorf("error preparing query BulkInsertCommits: %w", err)
	}
	if q.checkGithubRepoExistsStmt, err = db.PrepareContext(ctx, checkGithubRepoExists); err != nil {
		return nil, fmt.Errorf("error preparing query CheckGithubRepoExists: %w", err)
	}
	if q.checkGithubUserExistsStmt, err = db.PrepareContext(ctx, checkGithubUserExists); err != nil {
		return nil, fmt.Errorf("error preparing query CheckGithubUserExists: %w", err)
	}
	if q.checkGithubUserRestIdAuthorEmailExistsStmt, err = db.PrepareContext(ctx, checkGithubUserRestIdAuthorEmailExists); err != nil {
		return nil, fmt.Errorf("error preparing query CheckGithubUserRestIdAuthorEmailExists: %w", err)
	}
	if q.deleteRepoURLStmt, err = db.PrepareContext(ctx, deleteRepoURL); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteRepoURL: %w", err)
	}
	if q.getAndUpdateRepoURLStmt, err = db.PrepareContext(ctx, getAndUpdateRepoURL); err != nil {
		return nil, fmt.Errorf("error preparing query GetAndUpdateRepoURL: %w", err)
	}
	if q.getCommitStmt, err = db.PrepareContext(ctx, getCommit); err != nil {
		return nil, fmt.Errorf("error preparing query GetCommit: %w", err)
	}
	if q.getCommitsStmt, err = db.PrepareContext(ctx, getCommits); err != nil {
		return nil, fmt.Errorf("error preparing query GetCommits: %w", err)
	}
	if q.getCommitsWithAuthorInfoStmt, err = db.PrepareContext(ctx, getCommitsWithAuthorInfo); err != nil {
		return nil, fmt.Errorf("error preparing query GetCommitsWithAuthorInfo: %w", err)
	}
	if q.getFirstCommitStmt, err = db.PrepareContext(ctx, getFirstCommit); err != nil {
		return nil, fmt.Errorf("error preparing query GetFirstCommit: %w", err)
	}
	if q.getGithubRepoStmt, err = db.PrepareContext(ctx, getGithubRepo); err != nil {
		return nil, fmt.Errorf("error preparing query GetGithubRepo: %w", err)
	}
	if q.getGithubUserStmt, err = db.PrepareContext(ctx, getGithubUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetGithubUser: %w", err)
	}
	if q.getGroupOfEmailsStmt, err = db.PrepareContext(ctx, getGroupOfEmails); err != nil {
		return nil, fmt.Errorf("error preparing query GetGroupOfEmails: %w", err)
	}
	if q.getLatestCommitterDateStmt, err = db.PrepareContext(ctx, getLatestCommitterDate); err != nil {
		return nil, fmt.Errorf("error preparing query GetLatestCommitterDate: %w", err)
	}
	if q.getLatestUncheckedCommitPerAuthorStmt, err = db.PrepareContext(ctx, getLatestUncheckedCommitPerAuthor); err != nil {
		return nil, fmt.Errorf("error preparing query GetLatestUncheckedCommitPerAuthor: %w", err)
	}
	if q.getRepoURLStmt, err = db.PrepareContext(ctx, getRepoURL); err != nil {
		return nil, fmt.Errorf("error preparing query GetRepoURL: %w", err)
	}
	if q.getRepoURLsStmt, err = db.PrepareContext(ctx, getRepoURLs); err != nil {
		return nil, fmt.Errorf("error preparing query GetRepoURLs: %w", err)
	}
	if q.getReposStatusStmt, err = db.PrepareContext(ctx, getReposStatus); err != nil {
		return nil, fmt.Errorf("error preparing query GetReposStatus: %w", err)
	}
	if q.getUserCommitsForReposStmt, err = db.PrepareContext(ctx, getUserCommitsForRepos); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserCommitsForRepos: %w", err)
	}
	if q.insertCommitStmt, err = db.PrepareContext(ctx, insertCommit); err != nil {
		return nil, fmt.Errorf("error preparing query InsertCommit: %w", err)
	}
	if q.insertGithubRepoStmt, err = db.PrepareContext(ctx, insertGithubRepo); err != nil {
		return nil, fmt.Errorf("error preparing query InsertGithubRepo: %w", err)
	}
	if q.insertRestIdToEmailStmt, err = db.PrepareContext(ctx, insertRestIdToEmail); err != nil {
		return nil, fmt.Errorf("error preparing query InsertRestIdToEmail: %w", err)
	}
	if q.insertUserStmt, err = db.PrepareContext(ctx, insertUser); err != nil {
		return nil, fmt.Errorf("error preparing query InsertUser: %w", err)
	}
	if q.multiRowInsertCommitsStmt, err = db.PrepareContext(ctx, multiRowInsertCommits); err != nil {
		return nil, fmt.Errorf("error preparing query MultiRowInsertCommits: %w", err)
	}
	if q.updateStatusAndUpdatedAtStmt, err = db.PrepareContext(ctx, updateStatusAndUpdatedAt); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateStatusAndUpdatedAt: %w", err)
	}
	if q.upsertRepoURLStmt, err = db.PrepareContext(ctx, upsertRepoURL); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertRepoURL: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.bulkInsertCommitsStmt != nil {
		if cerr := q.bulkInsertCommitsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing bulkInsertCommitsStmt: %w", cerr)
		}
	}
	if q.checkGithubRepoExistsStmt != nil {
		if cerr := q.checkGithubRepoExistsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing checkGithubRepoExistsStmt: %w", cerr)
		}
	}
	if q.checkGithubUserExistsStmt != nil {
		if cerr := q.checkGithubUserExistsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing checkGithubUserExistsStmt: %w", cerr)
		}
	}
	if q.checkGithubUserRestIdAuthorEmailExistsStmt != nil {
		if cerr := q.checkGithubUserRestIdAuthorEmailExistsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing checkGithubUserRestIdAuthorEmailExistsStmt: %w", cerr)
		}
	}
	if q.deleteRepoURLStmt != nil {
		if cerr := q.deleteRepoURLStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteRepoURLStmt: %w", cerr)
		}
	}
	if q.getAndUpdateRepoURLStmt != nil {
		if cerr := q.getAndUpdateRepoURLStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAndUpdateRepoURLStmt: %w", cerr)
		}
	}
	if q.getCommitStmt != nil {
		if cerr := q.getCommitStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCommitStmt: %w", cerr)
		}
	}
	if q.getCommitsStmt != nil {
		if cerr := q.getCommitsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCommitsStmt: %w", cerr)
		}
	}
	if q.getCommitsWithAuthorInfoStmt != nil {
		if cerr := q.getCommitsWithAuthorInfoStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCommitsWithAuthorInfoStmt: %w", cerr)
		}
	}
	if q.getFirstCommitStmt != nil {
		if cerr := q.getFirstCommitStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getFirstCommitStmt: %w", cerr)
		}
	}
	if q.getGithubRepoStmt != nil {
		if cerr := q.getGithubRepoStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getGithubRepoStmt: %w", cerr)
		}
	}
	if q.getGithubUserStmt != nil {
		if cerr := q.getGithubUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getGithubUserStmt: %w", cerr)
		}
	}
	if q.getGroupOfEmailsStmt != nil {
		if cerr := q.getGroupOfEmailsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getGroupOfEmailsStmt: %w", cerr)
		}
	}
	if q.getLatestCommitterDateStmt != nil {
		if cerr := q.getLatestCommitterDateStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLatestCommitterDateStmt: %w", cerr)
		}
	}
	if q.getLatestUncheckedCommitPerAuthorStmt != nil {
		if cerr := q.getLatestUncheckedCommitPerAuthorStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getLatestUncheckedCommitPerAuthorStmt: %w", cerr)
		}
	}
	if q.getRepoURLStmt != nil {
		if cerr := q.getRepoURLStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRepoURLStmt: %w", cerr)
		}
	}
	if q.getRepoURLsStmt != nil {
		if cerr := q.getRepoURLsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRepoURLsStmt: %w", cerr)
		}
	}
	if q.getReposStatusStmt != nil {
		if cerr := q.getReposStatusStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getReposStatusStmt: %w", cerr)
		}
	}
	if q.getUserCommitsForReposStmt != nil {
		if cerr := q.getUserCommitsForReposStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserCommitsForReposStmt: %w", cerr)
		}
	}
	if q.insertCommitStmt != nil {
		if cerr := q.insertCommitStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertCommitStmt: %w", cerr)
		}
	}
	if q.insertGithubRepoStmt != nil {
		if cerr := q.insertGithubRepoStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertGithubRepoStmt: %w", cerr)
		}
	}
	if q.insertRestIdToEmailStmt != nil {
		if cerr := q.insertRestIdToEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertRestIdToEmailStmt: %w", cerr)
		}
	}
	if q.insertUserStmt != nil {
		if cerr := q.insertUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertUserStmt: %w", cerr)
		}
	}
	if q.multiRowInsertCommitsStmt != nil {
		if cerr := q.multiRowInsertCommitsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing multiRowInsertCommitsStmt: %w", cerr)
		}
	}
	if q.updateStatusAndUpdatedAtStmt != nil {
		if cerr := q.updateStatusAndUpdatedAtStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateStatusAndUpdatedAtStmt: %w", cerr)
		}
	}
	if q.upsertRepoURLStmt != nil {
		if cerr := q.upsertRepoURLStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertRepoURLStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                         DBTX
	tx                                         *sql.Tx
	bulkInsertCommitsStmt                      *sql.Stmt
	checkGithubRepoExistsStmt                  *sql.Stmt
	checkGithubUserExistsStmt                  *sql.Stmt
	checkGithubUserRestIdAuthorEmailExistsStmt *sql.Stmt
	deleteRepoURLStmt                          *sql.Stmt
	getAndUpdateRepoURLStmt                    *sql.Stmt
	getCommitStmt                              *sql.Stmt
	getCommitsStmt                             *sql.Stmt
	getCommitsWithAuthorInfoStmt               *sql.Stmt
	getFirstCommitStmt                         *sql.Stmt
	getGithubRepoStmt                          *sql.Stmt
	getGithubUserStmt                          *sql.Stmt
	getGroupOfEmailsStmt                       *sql.Stmt
	getLatestCommitterDateStmt                 *sql.Stmt
	getLatestUncheckedCommitPerAuthorStmt      *sql.Stmt
	getRepoURLStmt                             *sql.Stmt
	getRepoURLsStmt                            *sql.Stmt
	getReposStatusStmt                         *sql.Stmt
	getUserCommitsForReposStmt                 *sql.Stmt
	insertCommitStmt                           *sql.Stmt
	insertGithubRepoStmt                       *sql.Stmt
	insertRestIdToEmailStmt                    *sql.Stmt
	insertUserStmt                             *sql.Stmt
	multiRowInsertCommitsStmt                  *sql.Stmt
	updateStatusAndUpdatedAtStmt               *sql.Stmt
	upsertRepoURLStmt                          *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                        tx,
		tx:                        tx,
		bulkInsertCommitsStmt:     q.bulkInsertCommitsStmt,
		checkGithubRepoExistsStmt: q.checkGithubRepoExistsStmt,
		checkGithubUserExistsStmt: q.checkGithubUserExistsStmt,
		checkGithubUserRestIdAuthorEmailExistsStmt: q.checkGithubUserRestIdAuthorEmailExistsStmt,
		deleteRepoURLStmt:                          q.deleteRepoURLStmt,
		getAndUpdateRepoURLStmt:                    q.getAndUpdateRepoURLStmt,
		getCommitStmt:                              q.getCommitStmt,
		getCommitsStmt:                             q.getCommitsStmt,
		getCommitsWithAuthorInfoStmt:               q.getCommitsWithAuthorInfoStmt,
		getFirstCommitStmt:                         q.getFirstCommitStmt,
		getGithubRepoStmt:                          q.getGithubRepoStmt,
		getGithubUserStmt:                          q.getGithubUserStmt,
		getGroupOfEmailsStmt:                       q.getGroupOfEmailsStmt,
		getLatestCommitterDateStmt:                 q.getLatestCommitterDateStmt,
		getLatestUncheckedCommitPerAuthorStmt:      q.getLatestUncheckedCommitPerAuthorStmt,
		getRepoURLStmt:                             q.getRepoURLStmt,
		getRepoURLsStmt:                            q.getRepoURLsStmt,
		getReposStatusStmt:                         q.getReposStatusStmt,
		getUserCommitsForReposStmt:                 q.getUserCommitsForReposStmt,
		insertCommitStmt:                           q.insertCommitStmt,
		insertGithubRepoStmt:                       q.insertGithubRepoStmt,
		insertRestIdToEmailStmt:                    q.insertRestIdToEmailStmt,
		insertUserStmt:                             q.insertUserStmt,
		multiRowInsertCommitsStmt:                  q.multiRowInsertCommitsStmt,
		updateStatusAndUpdatedAtStmt:               q.updateStatusAndUpdatedAtStmt,
		upsertRepoURLStmt:                          q.upsertRepoURLStmt,
	}
}
