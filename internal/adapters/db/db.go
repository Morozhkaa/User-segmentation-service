// The db package provides methods for working directly with the database.
package db

import (
	"context"
	"fmt"
	"log"
	"segmentation-service/internal/domain/models"
	"segmentation-service/internal/ports"
	"segmentation-service/pkg/infra/logger"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	Pool *pgxpool.Pool
}

var _ ports.SegmentStorage = (*DBStorage)(nil)

// New establishes one connection and returns a new instance of DBStorage.
func New(ctx context.Context, conn string) (*DBStorage, error) {
	time.Sleep(time.Second)
	pool, err := pgxpool.Connect(ctx, conn)
	if err != nil {
		return nil, err
	}
	return &DBStorage{
		Pool: pool,
	}, nil
}

func (db *DBStorage) FindSegment(ctx context.Context, slug string) (count int, err error) {
	const query = `
	SELECT COUNT(*) FROM segments WHERE name = $1;
	`
	err = db.Pool.QueryRow(ctx, query, slug).Scan(&count)
	return count, err
}

func (db *DBStorage) SaveSegment(ctx context.Context, slug string) (err error) {
	const query = `
	INSERT INTO segments (name) VALUES ($1);
	`
	_, err = db.Pool.Exec(ctx, query, slug)
	return err
}

// DeleteSegment removes a segment and all users from it.
func (db *DBStorage) DeleteSegment(ctx context.Context, slug string) (err error) {
	// start a transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	// remove all users from a segment
	const queryDeleteUsers = `	
	DELETE FROM segments_users WHERE segments_id = (SELECT id FROM segments WHERE segments.name = $1);
	`
	_, err = tx.Exec(ctx, queryDeleteUsers, slug)
	if err != nil {
		return err
	}

	// delete the segment
	const queryDeleteSegment = `
	DELETE FROM segments WHERE name = $1;
	`
	_, err = tx.Exec(ctx, queryDeleteSegment, slug)
	if err != nil {
		return err
	}

	if tx.Commit(ctx) != nil {
		if rb := tx.Rollback(ctx); rb != nil {
			log.Fatalf("query failed: %v, unable to abort: %v", err, rb)
		}
	}
	return nil
}

// UpdateUserSegments adds and removes segments from a user. If one of the segments is not in the database, an error will be returned.
func (db *DBStorage) UpdateUserSegments(ctx context.Context, data models.UpdateRequest, userID uuid.UUID) error {
	// start a transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	logger := logger.Get()

	logger.Debug("start processing the list of segments for deletion")
	for _, slug := range data.SegmentsToRemove {
		// check that the segment with the given slug exists and get segment_id
		const querySegmId = `
		SELECT id FROM segments WHERE name = $1;
		`
		var segment_id int
		if err = tx.QueryRow(ctx, querySegmId, slug).Scan(&segment_id); err != nil {
			return models.ErrSegmentNotFound
		}

		// remove row from segments_users table
		const queryDelete = `
		DELETE FROM segments_users WHERE segments_id = $1 AND user_id = $2;
		`
		_, err = tx.Exec(ctx, queryDelete, segment_id, userID)
		if err != nil {
			logger.Debug("failed to delete entry from segments_users table")
			return err
		}

		// add delete entry to report table
		const queryReport = `
			INSERT INTO report (user_id, segments_id, action)
			VALUES ($1, $2, $3)
			`
		_, err = tx.Exec(ctx, queryReport, userID, segment_id, models.ActRemove)
		if err != nil {
			logger.Debug("failed to add delete record to report table")
			return err
		}
	}

	logger.Debug("start processing the list of segments to be added")
	for _, slug := range data.SegmentsToAdd {
		// check that the segment with the given slug exists
		const querySegmId = `
		SELECT id FROM segments WHERE name = $1;
		`
		var segment_id int
		if err = tx.QueryRow(ctx, querySegmId, slug).Scan(&segment_id); err != nil {
			return models.ErrSegmentNotFound
		}

		// check if the user is added to the segment
		const queryCheck = `
		SELECT COUNT(*) FROM segments_users WHERE segments_id = $1 AND user_id = $2;
		`
		var cnt int
		if err = tx.QueryRow(ctx, queryCheck, segment_id, userID).Scan(&cnt); err != nil {
			return err
		}

		// if there is no user in the segment, add it
		if cnt == 0 {
			const queryInsert = `
			INSERT INTO segments_users (segments_id, user_id)
			VALUES ($1, $2)
			`
			_, err = tx.Exec(ctx, queryInsert, segment_id, userID)
			if err != nil {
				logger.Debug("failed to add record to segments_users tablee")
				return err
			}

			// write add record to report table
			const queryReport = `
			INSERT INTO report (user_id, segments_id, action)
			VALUES ($1, $2, $3)
			`
			_, err = tx.Exec(ctx, queryReport, userID, segment_id, models.ActAdd)
			if err != nil {
				logger.Debug("failed to write add entry to report table")
				return err
			}
		}
	}
	if tx.Commit(ctx) != nil {
		if rb := tx.Rollback(ctx); rb != nil {
			log.Fatalf("query failed: %v, unable to abort: %v", err, rb)
		}
	}
	return nil
}

// GetUserSegments returns all segments the user is a member of.
func (db *DBStorage) GetUserSegments(ctx context.Context, userID uuid.UUID) (models.SegmentsList, error) {
	segments := models.SegmentsList{}
	const query = `
	SELECT segments.name FROM segments
	INNER JOIN segments_users ON segments.id=segments_users.segments_id WHERE user_id = $1;
	`
	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return segments, fmt.Errorf("can't get segments by user: %v", err)
	}

	for rows.Next() {
		var slug string
		if err = rows.Scan(&slug); err != nil {
			return segments, err
		}
		segments.S = append(segments.S, slug)
	}
	return segments, err
}

// GetMonthlyReport returns all entries about adding / removing users from segments for the specified month (in the format: yyyy-mm).
func (db *DBStorage) GetReport(ctx context.Context, period string) ([][]string, error) {
	var result [][]string

	const queryCheck = `
	SELECT user_id, segments.name, action, created_at FROM segments
	INNER JOIN report ON segments.id=report.segments_id
	WHERE created_at between date($1) and date($1) + interval '1 month';
	`
	rows, err := db.Pool.Query(ctx, queryCheck, period)
	if err != nil {
		return result, fmt.Errorf("getting report for the month since '%s' failed: %v", period, err)
	}

	for rows.Next() {
		var line models.ReportRow
		err = rows.Scan(
			&line.UserID,
			&line.SegmentName,
			&line.Action,
			&line.Time,
		)
		if err != nil {
			return result, err
		}
		arr := []string{fmt.Sprintf("%v", line.UserID), string(line.SegmentName), string(line.Action), line.Time.Add(3 * time.Hour).Format("2006-01-02 15:04:05")}
		result = append(result, arr)
	}
	return result, nil
}

// GetMonthlyReport returns all entries about adding / removing users from segments for the specified month (in the format: yyyy-mm).
func (db *DBStorage) GetUserReport(ctx context.Context, period string, userID uuid.UUID) ([][]string, error) {
	var result [][]string

	const queryCheck = `
	SELECT user_id, segments.name, action, created_at FROM segments
	INNER JOIN report ON segments.id=report.segments_id
	WHERE (created_at between date($1) and date($1) + interval '1 month') AND user_id = $2;
	`
	rows, err := db.Pool.Query(ctx, queryCheck, period, userID)
	if err != nil {
		return result, fmt.Errorf("getting report for the month since '%s' failed: %v", period, err)
	}

	for rows.Next() {
		var line models.ReportRow
		err = rows.Scan(
			&line.UserID,
			&line.SegmentName,
			&line.Action,
			&line.Time,
		)
		if err != nil {
			return result, err
		}
		arr := []string{fmt.Sprintf("%v", line.UserID), string(line.SegmentName), string(line.Action), line.Time.Add(3 * time.Hour).Format("2006-01-02 15:04:05")}
		result = append(result, arr)
	}
	return result, nil
}
