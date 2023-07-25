package db

const (
	racesList = "list"
	raceByID  = "byID" // const for query a single race by id.
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
					WHEN advertised_start_time <= datetime('now') THEN 'CLOSED'
					ELSE 'OPEN'
				END as status
			FROM races
		`,
		raceByID: `
			SELECT 
					id, 
					meeting_id, 
					name, 
					number, 
					visible, 
					advertised_start_time,
					CASE
						WHEN advertised_start_time <= datetime('now') THEN 'CLOSED'
						ELSE 'OPEN'
					END as status
				FROM races
				WHERE id = ? 
		`,
	}
}
