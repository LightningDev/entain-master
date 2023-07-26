package db

const (
	eventsList = "list"
	eventByID  = "byID"
)

func getEventQueries() map[string]string {
	// Status calculated on the fly.
	return map[string]string{
		eventsList: `
			SELECT 
				id, 
				name, 
				location,
				sport,
				advertised_start_time,
				CASE
					WHEN datetime(advertised_start_time) <= datetime('now') THEN 'CLOSED'
					ELSE 'OPEN'
				END as status
			FROM events
		`,
		eventByID: `
			SELECT 
				id, 
				name, 
				location,
				sport,
				advertised_start_time,
				CASE
					WHEN datetime(advertised_start_time) <= datetime('now') THEN 'CLOSED'
					ELSE 'OPEN'
				END as status
			FROM events WHERE id = ?
		`,
	}
}
