// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

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
	if q.batchInsertRepoDependenciesStmt, err = db.PrepareContext(ctx, batchInsertRepoDependencies); err != nil {
		return nil, fmt.Errorf("error preparing query BatchInsertRepoDependencies: %w", err)
	}
	if q.bulkInsertCommitsStmt, err = db.PrepareContext(ctx, bulkInsertCommits); err != nil {
		return nil, fmt.Errorf("error preparing query BulkInsertCommits: %w", err)
	}
	if q.bulkInsertDependenciesStmt, err = db.PrepareContext(ctx, bulkInsertDependencies); err != nil {
		return nil, fmt.Errorf("error preparing query BulkInsertDependencies: %w", err)
	}
	if q.bulkInsertUserDependenciesStmt, err = db.PrepareContext(ctx, bulkInsertUserDependencies); err != nil {
		return nil, fmt.Errorf("error preparing query BulkInsertUserDependencies: %w", err)
	}
	if q.bulkUpsertFilePatternsStmt, err = db.PrepareContext(ctx, bulkUpsertFilePatterns); err != nil {
		return nil, fmt.Errorf("error preparing query BulkUpsertFilePatterns: %w", err)
	}
	if q.checkGithubRepoExistsStmt, err = db.PrepareContext(ctx, checkGithubRepoExists); err != nil {
		return nil, fmt.Errorf("error preparing query CheckGithubRepoExists: %w", err)
	}
	if q.checkGithubUserExistsStmt, err = db.PrepareContext(ctx, checkGithubUserExists); err != nil {
		return nil, fmt.Errorf("error preparing query CheckGithubUserExists: %w", err)
	}
	if q.checkGithubUserIdExistsStmt, err = db.PrepareContext(ctx, checkGithubUserIdExists); err != nil {
		return nil, fmt.Errorf("error preparing query CheckGithubUserIdExists: %w", err)
	}
	if q.checkGithubUserRestIdAuthorEmailExistsStmt, err = db.PrepareContext(ctx, checkGithubUserRestIdAuthorEmailExists); err != nil {
		return nil, fmt.Errorf("error preparing query CheckGithubUserRestIdAuthorEmailExists: %w", err)
	}
	if q.deleteRepoURLStmt, err = db.PrepareContext(ctx, deleteRepoURL); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteRepoURL: %w", err)
	}
	if q.deleteUnusedDependenciesStmt, err = db.PrepareContext(ctx, deleteUnusedDependencies); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUnusedDependencies: %w", err)
	}
	if q.getAllFilePatternsStmt, err = db.PrepareContext(ctx, getAllFilePatterns); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllFilePatterns: %w", err)
	}
	if q.getAllUserDependenciesByUserStmt, err = db.PrepareContext(ctx, getAllUserDependenciesByUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllUserDependenciesByUser: %w", err)
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
	if q.getDependenciesStmt, err = db.PrepareContext(ctx, getDependencies); err != nil {
		return nil, fmt.Errorf("error preparing query GetDependencies: %w", err)
	}
	if q.getDependenciesByFilesStmt, err = db.PrepareContext(ctx, getDependenciesByFiles); err != nil {
		return nil, fmt.Errorf("error preparing query GetDependenciesByFiles: %w", err)
	}
	if q.getDependenciesByNamesStmt, err = db.PrepareContext(ctx, getDependenciesByNames); err != nil {
		return nil, fmt.Errorf("error preparing query GetDependenciesByNames: %w", err)
	}
	if q.getDependencyStmt, err = db.PrepareContext(ctx, getDependency); err != nil {
		return nil, fmt.Errorf("error preparing query GetDependency: %w", err)
	}
	if q.getFirstAndLastCommitStmt, err = db.PrepareContext(ctx, getFirstAndLastCommit); err != nil {
		return nil, fmt.Errorf("error preparing query GetFirstAndLastCommit: %w", err)
	}
	if q.getFirstCommitStmt, err = db.PrepareContext(ctx, getFirstCommit); err != nil {
		return nil, fmt.Errorf("error preparing query GetFirstCommit: %w", err)
	}
	if q.getGithubRepoStmt, err = db.PrepareContext(ctx, getGithubRepo); err != nil {
		return nil, fmt.Errorf("error preparing query GetGithubRepo: %w", err)
	}
	if q.getGithubRepoByUrlStmt, err = db.PrepareContext(ctx, getGithubRepoByUrl); err != nil {
		return nil, fmt.Errorf("error preparing query GetGithubRepoByUrl: %w", err)
	}
	if q.getGithubUserStmt, err = db.PrepareContext(ctx, getGithubUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetGithubUser: %w", err)
	}
	if q.getGithubUserByCommitEmailStmt, err = db.PrepareContext(ctx, getGithubUserByCommitEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetGithubUserByCommitEmail: %w", err)
	}
	if q.getGithubUserByRestIdStmt, err = db.PrepareContext(ctx, getGithubUserByRestId); err != nil {
		return nil, fmt.Errorf("error preparing query GetGithubUserByRestId: %w", err)
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
	if q.getRepoAuthorsInfoStmt, err = db.PrepareContext(ctx, getRepoAuthorsInfo); err != nil {
		return nil, fmt.Errorf("error preparing query GetRepoAuthorsInfo: %w", err)
	}
	if q.getRepoDependenciesStmt, err = db.PrepareContext(ctx, getRepoDependencies); err != nil {
		return nil, fmt.Errorf("error preparing query GetRepoDependencies: %w", err)
	}
	if q.getRepoDependenciesByURLStmt, err = db.PrepareContext(ctx, getRepoDependenciesByURL); err != nil {
		return nil, fmt.Errorf("error preparing query GetRepoDependenciesByURL: %w", err)
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
	if q.getUserDependenciesByUpdatedAtStmt, err = db.PrepareContext(ctx, getUserDependenciesByUpdatedAt); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserDependenciesByUpdatedAt: %w", err)
	}
	if q.getUserDependenciesByUserStmt, err = db.PrepareContext(ctx, getUserDependenciesByUser); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserDependenciesByUser: %w", err)
	}
	if q.initializeRepoDependenciesStmt, err = db.PrepareContext(ctx, initializeRepoDependencies); err != nil {
		return nil, fmt.Errorf("error preparing query InitializeRepoDependencies: %w", err)
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
	if q.setAllCommitsToCheckedStmt, err = db.PrepareContext(ctx, setAllCommitsToChecked); err != nil {
		return nil, fmt.Errorf("error preparing query SetAllCommitsToChecked: %w", err)
	}
	if q.switchReposRelationToSimpleStmt, err = db.PrepareContext(ctx, switchReposRelationToSimple); err != nil {
		return nil, fmt.Errorf("error preparing query SwitchReposRelationToSimple: %w", err)
	}
	if q.switchUsersRelationToSimpleStmt, err = db.PrepareContext(ctx, switchUsersRelationToSimple); err != nil {
		return nil, fmt.Errorf("error preparing query SwitchUsersRelationToSimple: %w", err)
	}
	if q.updateStatusAndUpdatedAtStmt, err = db.PrepareContext(ctx, updateStatusAndUpdatedAt); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateStatusAndUpdatedAt: %w", err)
	}
	if q.updateStatusAndUpdatedAtepoUrlV2Stmt, err = db.PrepareContext(ctx, updateStatusAndUpdatedAtepoUrlV2); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateStatusAndUpdatedAtepoUrlV2: %w", err)
	}
	if q.upsertMissingDependenciesStmt, err = db.PrepareContext(ctx, upsertMissingDependencies); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertMissingDependencies: %w", err)
	}
	if q.upsertRepoToUserByIdStmt, err = db.PrepareContext(ctx, upsertRepoToUserById); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertRepoToUserById: %w", err)
	}
	if q.upsertRepoURLStmt, err = db.PrepareContext(ctx, upsertRepoURL); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertRepoURL: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.batchInsertRepoDependenciesStmt != nil {
		if cerr := q.batchInsertRepoDependenciesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing batchInsertRepoDependenciesStmt: %w", cerr)
		}
	}
	if q.bulkInsertCommitsStmt != nil {
		if cerr := q.bulkInsertCommitsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing bulkInsertCommitsStmt: %w", cerr)
		}
	}
	if q.bulkInsertDependenciesStmt != nil {
		if cerr := q.bulkInsertDependenciesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing bulkInsertDependenciesStmt: %w", cerr)
		}
	}
	if q.bulkInsertUserDependenciesStmt != nil {
		if cerr := q.bulkInsertUserDependenciesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing bulkInsertUserDependenciesStmt: %w", cerr)
		}
	}
	if q.bulkUpsertFilePatternsStmt != nil {
		if cerr := q.bulkUpsertFilePatternsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing bulkUpsertFilePatternsStmt: %w", cerr)
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
	if q.checkGithubUserIdExistsStmt != nil {
		if cerr := q.checkGithubUserIdExistsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing checkGithubUserIdExistsStmt: %w", cerr)
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
	if q.deleteUnusedDependenciesStmt != nil {
		if cerr := q.deleteUnusedDependenciesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUnusedDependenciesStmt: %w", cerr)
		}
	}
	if q.getAllFilePatternsStmt != nil {
		if cerr := q.getAllFilePatternsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllFilePatternsStmt: %w", cerr)
		}
	}
	if q.getAllUserDependenciesByUserStmt != nil {
		if cerr := q.getAllUserDependenciesByUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllUserDependenciesByUserStmt: %w", cerr)
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
	if q.getDependenciesStmt != nil {
		if cerr := q.getDependenciesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDependenciesStmt: %w", cerr)
		}
	}
	if q.getDependenciesByFilesStmt != nil {
		if cerr := q.getDependenciesByFilesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDependenciesByFilesStmt: %w", cerr)
		}
	}
	if q.getDependenciesByNamesStmt != nil {
		if cerr := q.getDependenciesByNamesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDependenciesByNamesStmt: %w", cerr)
		}
	}
	if q.getDependencyStmt != nil {
		if cerr := q.getDependencyStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getDependencyStmt: %w", cerr)
		}
	}
	if q.getFirstAndLastCommitStmt != nil {
		if cerr := q.getFirstAndLastCommitStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getFirstAndLastCommitStmt: %w", cerr)
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
	if q.getGithubRepoByUrlStmt != nil {
		if cerr := q.getGithubRepoByUrlStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getGithubRepoByUrlStmt: %w", cerr)
		}
	}
	if q.getGithubUserStmt != nil {
		if cerr := q.getGithubUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getGithubUserStmt: %w", cerr)
		}
	}
	if q.getGithubUserByCommitEmailStmt != nil {
		if cerr := q.getGithubUserByCommitEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getGithubUserByCommitEmailStmt: %w", cerr)
		}
	}
	if q.getGithubUserByRestIdStmt != nil {
		if cerr := q.getGithubUserByRestIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getGithubUserByRestIdStmt: %w", cerr)
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
	if q.getRepoAuthorsInfoStmt != nil {
		if cerr := q.getRepoAuthorsInfoStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRepoAuthorsInfoStmt: %w", cerr)
		}
	}
	if q.getRepoDependenciesStmt != nil {
		if cerr := q.getRepoDependenciesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRepoDependenciesStmt: %w", cerr)
		}
	}
	if q.getRepoDependenciesByURLStmt != nil {
		if cerr := q.getRepoDependenciesByURLStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRepoDependenciesByURLStmt: %w", cerr)
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
	if q.getUserDependenciesByUpdatedAtStmt != nil {
		if cerr := q.getUserDependenciesByUpdatedAtStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserDependenciesByUpdatedAtStmt: %w", cerr)
		}
	}
	if q.getUserDependenciesByUserStmt != nil {
		if cerr := q.getUserDependenciesByUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserDependenciesByUserStmt: %w", cerr)
		}
	}
	if q.initializeRepoDependenciesStmt != nil {
		if cerr := q.initializeRepoDependenciesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing initializeRepoDependenciesStmt: %w", cerr)
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
	if q.setAllCommitsToCheckedStmt != nil {
		if cerr := q.setAllCommitsToCheckedStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setAllCommitsToCheckedStmt: %w", cerr)
		}
	}
	if q.switchReposRelationToSimpleStmt != nil {
		if cerr := q.switchReposRelationToSimpleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing switchReposRelationToSimpleStmt: %w", cerr)
		}
	}
	if q.switchUsersRelationToSimpleStmt != nil {
		if cerr := q.switchUsersRelationToSimpleStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing switchUsersRelationToSimpleStmt: %w", cerr)
		}
	}
	if q.updateStatusAndUpdatedAtStmt != nil {
		if cerr := q.updateStatusAndUpdatedAtStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateStatusAndUpdatedAtStmt: %w", cerr)
		}
	}
	if q.updateStatusAndUpdatedAtepoUrlV2Stmt != nil {
		if cerr := q.updateStatusAndUpdatedAtepoUrlV2Stmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateStatusAndUpdatedAtepoUrlV2Stmt: %w", cerr)
		}
	}
	if q.upsertMissingDependenciesStmt != nil {
		if cerr := q.upsertMissingDependenciesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertMissingDependenciesStmt: %w", cerr)
		}
	}
	if q.upsertRepoToUserByIdStmt != nil {
		if cerr := q.upsertRepoToUserByIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertRepoToUserByIdStmt: %w", cerr)
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
	batchInsertRepoDependenciesStmt            *sql.Stmt
	bulkInsertCommitsStmt                      *sql.Stmt
	bulkInsertDependenciesStmt                 *sql.Stmt
	bulkInsertUserDependenciesStmt             *sql.Stmt
	bulkUpsertFilePatternsStmt                 *sql.Stmt
	checkGithubRepoExistsStmt                  *sql.Stmt
	checkGithubUserExistsStmt                  *sql.Stmt
	checkGithubUserIdExistsStmt                *sql.Stmt
	checkGithubUserRestIdAuthorEmailExistsStmt *sql.Stmt
	deleteRepoURLStmt                          *sql.Stmt
	deleteUnusedDependenciesStmt               *sql.Stmt
	getAllFilePatternsStmt                     *sql.Stmt
	getAllUserDependenciesByUserStmt           *sql.Stmt
	getAndUpdateRepoURLStmt                    *sql.Stmt
	getCommitStmt                              *sql.Stmt
	getCommitsStmt                             *sql.Stmt
	getCommitsWithAuthorInfoStmt               *sql.Stmt
	getDependenciesStmt                        *sql.Stmt
	getDependenciesByFilesStmt                 *sql.Stmt
	getDependenciesByNamesStmt                 *sql.Stmt
	getDependencyStmt                          *sql.Stmt
	getFirstAndLastCommitStmt                  *sql.Stmt
	getFirstCommitStmt                         *sql.Stmt
	getGithubRepoStmt                          *sql.Stmt
	getGithubRepoByUrlStmt                     *sql.Stmt
	getGithubUserStmt                          *sql.Stmt
	getGithubUserByCommitEmailStmt             *sql.Stmt
	getGithubUserByRestIdStmt                  *sql.Stmt
	getGroupOfEmailsStmt                       *sql.Stmt
	getLatestCommitterDateStmt                 *sql.Stmt
	getLatestUncheckedCommitPerAuthorStmt      *sql.Stmt
	getRepoAuthorsInfoStmt                     *sql.Stmt
	getRepoDependenciesStmt                    *sql.Stmt
	getRepoDependenciesByURLStmt               *sql.Stmt
	getRepoURLStmt                             *sql.Stmt
	getRepoURLsStmt                            *sql.Stmt
	getReposStatusStmt                         *sql.Stmt
	getUserCommitsForReposStmt                 *sql.Stmt
	getUserDependenciesByUpdatedAtStmt         *sql.Stmt
	getUserDependenciesByUserStmt              *sql.Stmt
	initializeRepoDependenciesStmt             *sql.Stmt
	insertGithubRepoStmt                       *sql.Stmt
	insertRestIdToEmailStmt                    *sql.Stmt
	insertUserStmt                             *sql.Stmt
	setAllCommitsToCheckedStmt                 *sql.Stmt
	switchReposRelationToSimpleStmt            *sql.Stmt
	switchUsersRelationToSimpleStmt            *sql.Stmt
	updateStatusAndUpdatedAtStmt               *sql.Stmt
	updateStatusAndUpdatedAtepoUrlV2Stmt       *sql.Stmt
	upsertMissingDependenciesStmt              *sql.Stmt
	upsertRepoToUserByIdStmt                   *sql.Stmt
	upsertRepoURLStmt                          *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                              tx,
		tx:                              tx,
		batchInsertRepoDependenciesStmt: q.batchInsertRepoDependenciesStmt,
		bulkInsertCommitsStmt:           q.bulkInsertCommitsStmt,
		bulkInsertDependenciesStmt:      q.bulkInsertDependenciesStmt,
		bulkInsertUserDependenciesStmt:  q.bulkInsertUserDependenciesStmt,
		bulkUpsertFilePatternsStmt:      q.bulkUpsertFilePatternsStmt,
		checkGithubRepoExistsStmt:       q.checkGithubRepoExistsStmt,
		checkGithubUserExistsStmt:       q.checkGithubUserExistsStmt,
		checkGithubUserIdExistsStmt:     q.checkGithubUserIdExistsStmt,
		checkGithubUserRestIdAuthorEmailExistsStmt: q.checkGithubUserRestIdAuthorEmailExistsStmt,
		deleteRepoURLStmt:                          q.deleteRepoURLStmt,
		deleteUnusedDependenciesStmt:               q.deleteUnusedDependenciesStmt,
		getAllFilePatternsStmt:                     q.getAllFilePatternsStmt,
		getAllUserDependenciesByUserStmt:           q.getAllUserDependenciesByUserStmt,
		getAndUpdateRepoURLStmt:                    q.getAndUpdateRepoURLStmt,
		getCommitStmt:                              q.getCommitStmt,
		getCommitsStmt:                             q.getCommitsStmt,
		getCommitsWithAuthorInfoStmt:               q.getCommitsWithAuthorInfoStmt,
		getDependenciesStmt:                        q.getDependenciesStmt,
		getDependenciesByFilesStmt:                 q.getDependenciesByFilesStmt,
		getDependenciesByNamesStmt:                 q.getDependenciesByNamesStmt,
		getDependencyStmt:                          q.getDependencyStmt,
		getFirstAndLastCommitStmt:                  q.getFirstAndLastCommitStmt,
		getFirstCommitStmt:                         q.getFirstCommitStmt,
		getGithubRepoStmt:                          q.getGithubRepoStmt,
		getGithubRepoByUrlStmt:                     q.getGithubRepoByUrlStmt,
		getGithubUserStmt:                          q.getGithubUserStmt,
		getGithubUserByCommitEmailStmt:             q.getGithubUserByCommitEmailStmt,
		getGithubUserByRestIdStmt:                  q.getGithubUserByRestIdStmt,
		getGroupOfEmailsStmt:                       q.getGroupOfEmailsStmt,
		getLatestCommitterDateStmt:                 q.getLatestCommitterDateStmt,
		getLatestUncheckedCommitPerAuthorStmt:      q.getLatestUncheckedCommitPerAuthorStmt,
		getRepoAuthorsInfoStmt:                     q.getRepoAuthorsInfoStmt,
		getRepoDependenciesStmt:                    q.getRepoDependenciesStmt,
		getRepoDependenciesByURLStmt:               q.getRepoDependenciesByURLStmt,
		getRepoURLStmt:                             q.getRepoURLStmt,
		getRepoURLsStmt:                            q.getRepoURLsStmt,
		getReposStatusStmt:                         q.getReposStatusStmt,
		getUserCommitsForReposStmt:                 q.getUserCommitsForReposStmt,
		getUserDependenciesByUpdatedAtStmt:         q.getUserDependenciesByUpdatedAtStmt,
		getUserDependenciesByUserStmt:              q.getUserDependenciesByUserStmt,
		initializeRepoDependenciesStmt:             q.initializeRepoDependenciesStmt,
		insertGithubRepoStmt:                       q.insertGithubRepoStmt,
		insertRestIdToEmailStmt:                    q.insertRestIdToEmailStmt,
		insertUserStmt:                             q.insertUserStmt,
		setAllCommitsToCheckedStmt:                 q.setAllCommitsToCheckedStmt,
		switchReposRelationToSimpleStmt:            q.switchReposRelationToSimpleStmt,
		switchUsersRelationToSimpleStmt:            q.switchUsersRelationToSimpleStmt,
		updateStatusAndUpdatedAtStmt:               q.updateStatusAndUpdatedAtStmt,
		updateStatusAndUpdatedAtepoUrlV2Stmt:       q.updateStatusAndUpdatedAtepoUrlV2Stmt,
		upsertMissingDependenciesStmt:              q.upsertMissingDependenciesStmt,
		upsertRepoToUserByIdStmt:                   q.upsertRepoToUserByIdStmt,
		upsertRepoURLStmt:                          q.upsertRepoURLStmt,
	}
}
