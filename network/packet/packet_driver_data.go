package packet

import (
	"bytes"
	"fmt"
	"strings"
)

type mood string

const (
	MoodTired         mood = "TIRED"
	MoodInexperienced mood = "INEXPERIENCED"
	MoodDistracted    mood = "DISTRACTED"
	MoodImpetuous     mood = "IMPETUOUS"
)

func getMoodFromString(moodStr string) mood {
	allowedMoods := []mood{MoodTired, MoodInexperienced, MoodDistracted, MoodImpetuous}
	for _, mood := range allowedMoods {
		if moodStr == string(mood) {
			return mood
		}
	}

	return ""
}

type DriverData struct {
	Moods []mood `json:"moods"`
}

func (dt DriverData) Encode() []byte {
	var bfr bytes.Buffer
	writeSubseparatedValue(&bfr, dt.getMoodsBytes())

	return bfr.Bytes()
}

func (dt *DriverData) Decode(bts []byte) error {
	parts := bytes.Split(bts, []byte{innerSeparator})
	const mustPartsCount = 1
	if len(parts) != mustPartsCount {
		return fmt.Errorf("Driver Data: Unexpected number of parts, expected %d, received %d.", mustPartsCount, len(parts))
	}

	err := dt.decodeMoods(parts[0])
	return err
}

func (dt *DriverData) decodeMoods(moodBts []byte) error {
	parts := strings.Split(string(moodBts), ",")
	dt.Moods = nil
	for _, moodStr := range parts {
		if m := getMoodFromString(moodStr); m != "" {
			dt.Moods = append(dt.Moods, mood(m))
		}
	}

	return nil
}

func (dt DriverData) getMoodsBytes() []byte {
	var bfr bytes.Buffer
	once := false
	for _, mood := range dt.Moods {
		if once {
			bfr.WriteByte(',')
		} else {
			once = true
		}

		bfr.WriteString(string(mood))
	}

	return bfr.Bytes()
}

func Moods(mds ...mood) []mood {
	return mds
}
