package db

const (
	racesList = "list"
)

func getRaceQueries() map[string]string {
	// Status calculated on the fly.
	return map[string]string{
		racesList: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time,
				CASE
					WHEN datetime(advertised_start_time) <= datetime('now') THEN 'CLOSED'
					ELSE 'OPEN'
				END as status
			FROM races
		`,
	}
}
