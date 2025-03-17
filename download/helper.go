package download

import "time"

const day = time.Hour * 24

func yesterday() time.Time {
	return time.Now().Add(-day)
}
