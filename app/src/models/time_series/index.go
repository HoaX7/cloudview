package models

import (
	"cloudview/app/dbschema/cloudview/public/table"
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/api/middleware/logger"
	"cloudview/app/src/database"
	models "cloudview/app/src/models/structs"
	"cloudview/app/src/utility"
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type GetByMetricIdInput struct {
	MetricPanelID     uuid.UUID
	ProviderAccountID uuid.UUID
	MetricType        string
	StartTime         *int64 `json:"startTime"` // ms
	EndTime           *int64 `json:"endTime"`   // ms
	Period            *int   `json:"period"`    // sec
	/*
		`avg` or `sum`.
		default: will return 5 minutes worth of data betwee start and end time.
	*/
	Aggregation *string `json:"aggregation"`
}

type GetExistingTsInput struct {
	MetricPanelID      uuid.UUID
	TruncatedTimestamp int64
	MetricType         string
}
type GetExistingTsOutput struct {
	ID int32 `sql:"primary_key" json:"id"`
}

func _appendTsData(db *database.DB, data models.TimeSeries) error {
	query := "update time_series set series = series || $1 where id = $2;"
	series := pq.Array(*data.Series)
	_, err := db.Postgres.Query(query, series, data.ID)
	if err != nil {
		logger.Logger.Error("time_series._appendTsData: error", err)
		return custom_errors.DBErrors(err)
	}
	return nil
}

func _create(db *database.DB, data models.TimeSeries) error {
	columnList := postgres.ColumnList{
		table.TimeSeries.Series,
		table.TimeSeries.Type,
		table.TimeSeries.MetricPanelID,
		table.TimeSeries.TruncatedTimestamp,
	}

	insertData := struct {
		Series             any
		Type               models.MetricType
		MetricPanelID      uuid.UUID
		TruncatedTimestamp int64
	}{
		Series:             pq.Array(*data.Series),
		Type:               data.Type,
		MetricPanelID:      data.MetricPanelID,
		TruncatedTimestamp: data.TruncatedTimestamp,
	}

	stmt := table.TimeSeries.INSERT(columnList).MODEL(insertData)

	queryString, args := stmt.Sql()
	logger.Logger.Log("Inserting into Time Series table with data: ", queryString, args)

	_, err := db.Postgres.Query(queryString, args...)
	if err != nil {
		logger.Logger.Error("models.time_series.Create: ERROR", err)
		return custom_errors.DBErrors(err)
	}
	return nil
}

func queryMetricFromBuilder(db *database.DB, opts GetByMetricIdInput) (*sql.Rows, error) {
	stmt := table.TimeSeries.SELECT(table.TimeSeries.AllColumns).
		FROM(table.TimeSeries.LEFT_JOIN(table.MetricPanels, table.MetricPanels.ID.EQ(table.TimeSeries.MetricPanelID))).
		WHERE(postgres.AND(
			table.TimeSeries.MetricPanelID.EQ(postgres.UUID(opts.MetricPanelID)),
			table.TimeSeries.Type.EQ(postgres.NewEnumValue(opts.MetricType)),
			table.MetricPanels.ProviderAccountID.EQ(postgres.UUID(opts.ProviderAccountID)),
		))

	debug := stmt.DebugSql()
	logger.Logger.Log("time_series._getByMetricId: debug query", debug)
	query, args := stmt.Sql()
	/*
		When using `db.Postgres.Query` it is necessary
		to check for `row.Next()` since it does not
		return an error if there are no rows.
	*/
	return db.Postgres.Query(query, args...)
}

type queryStruct struct {
	StartTime *int64
	EndTime   *int64
}

/*
TODO - Need to add aggregation query. i.e ('avg', 'sum')
This is useful to fetch historical data that have large intervals (30 days+)
We can then aggregate data to return average metric points.

This function is also used by default to fetch data in real-time.
*/
func runRawMetricQuery(db *database.DB, opts GetByMetricIdInput) (*sql.Rows, error) {
	query := `SELECT id, truncated_timestamp, 
	ARRAY_AGG(ARRAY[ts, val] ORDER BY ts) AS series, type, metric_panel_id,
	created_at, updated_at FROM (
		SELECT id, truncated_timestamp, type, metric_panel_id, created_at,
		updated_at, unnest(series[:][1:1]) AS ts, unnest(series[:][2:2]) AS val
		FROM time_series
	) AS x
	WHERE ts BETWEEN $1 AND $2
	GROUP BY id, truncated_timestamp, type, metric_panel_id, created_at, updated_at;`

	return db.Postgres.Query(query, *opts.StartTime, *opts.EndTime)
}

/*
This method is used to fetch usage metric data with aggregation or
real-time.
*/
func _getByMetricId(db *database.DB, opts GetByMetricIdInput) ([]models.TimeSeries, error) {
	var result []models.TimeSeries
	/*
		If StartTime and EndTime is present need to run raw query
		to fetch the correct data.
	*/
	var (
		rows *sql.Rows
		err  error
	)
	if *opts.StartTime > 0 {
		rows, err = runRawMetricQuery(db, opts)
	} else {
		rows, err = queryMetricFromBuilder(db, opts)
	}
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return result, custom_errors.NoDataFound
		}
		return result, err
	}
	for rows.Next() {
		var item models.TimeSeries
		var series []uint8
		if err := rows.Scan(&item.ID, &item.TruncatedTimestamp, &series,
			&item.Type, &item.MetricPanelID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			logger.Logger.Error("error scanning row values:", err)
			continue
		}
		item.Series = utility.ParseDoublePrecision2DToFloat2D(series)
		result = append(result, item)
	}
	logger.Logger.Log(result)
	return result, nil
}

/*
This method allows us to identify if we should append the time-series to existing row
or create new row.
*/
func _getExistingTs(db *database.DB, opts GetExistingTsInput) (*GetExistingTsOutput, error) {
	stmt := table.TimeSeries.SELECT(table.TimeSeries.ID.AS("id")).
		WHERE(postgres.AND(
			table.TimeSeries.MetricPanelID.EQ(postgres.UUID(opts.MetricPanelID)),
			table.TimeSeries.Type.EQ(postgres.NewEnumValue(opts.MetricType)),
			table.TimeSeries.TruncatedTimestamp.EQ(postgres.Int64(opts.TruncatedTimestamp)),
		))

	debug := stmt.DebugSql()
	logger.Logger.Log("time_series._getExistingTs: debug query", debug)
	var result GetExistingTsOutput
	query, args := stmt.Sql()
	row, err := db.Postgres.Query(query, args...)
	if err != nil {
		return nil, custom_errors.DBErrors(err)
	}
	if !row.Next() {
		return nil, custom_errors.NoDataFound
	}
	if err := row.Scan(&result.ID); err != nil {
		logger.Logger.Error("time_series._getExistingTs: error scanning row", err)
	}
	logger.Logger.Log("time_series._getExistingTs: row found, id:", result.ID)
	return &result, nil
}
