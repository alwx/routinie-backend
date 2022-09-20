package db

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"habiko-go/pkg/models"
)

func getSimpleTrackerMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "20211230-1",
			Migrate: func(tx *gorm.DB) error {
				trackers := []models.SampleTracker{
					{
						ID:          "cook-healthy-food",
						Title:       "Cook healthy food",
						Description: "Establish healthy cooking & eating habits.",
						Emoji:       "🥗",
						Tags:        []string{"food", "healthy", "eat", "cook", "salads"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     1,
							"default_change": 1,
						}),
						Priority: 10,
					},
					{
						ID:          "read-10-pages",
						Title:       "Read at least 10 pages",
						Description: "Challenge yourself by reading more than you usually do.",
						Emoji:       "📚",
						Tags:        []string{"read", "pages", "books", "magazines", "articles"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     10,
							"default_change": 1,
							"is_infinite":    true,
						}),
						Priority: 10,
					},
					{
						ID:          "50-pushups",
						Title:       "Do at least 50 push ups",
						Description: "Try adding some sports to your life.",
						Emoji:       "💪",
						Tags:        []string{"sport", "workout", "pushups", "gym", "push ups", "fitness"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     50,
							"default_change": 10,
							"is_infinite":    true,
						}),
						Priority: 10,
					},
					{
						ID:          "meditate-for-5-minutes",
						Title:       "Meditate for 5 minutes",
						Description: "Try to establish a meditation routine.",
						Emoji:       "🧘",
						Tags:        []string{"health", "meditation", "yoga", "sport", "gym", "fitness"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     10,
							"default_change": 5,
							"is_infinite":    true,
						}),
						Priority: 5,
					},
					{
						ID:          "do-yoga",
						Title:       "Do yoga",
						Description: "Start practicing yoga regularly.",
						Emoji:       "🧘‍️",
						Tags:        []string{"health", "yoga", "sport", "gym"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     1,
							"default_change": 1,
						}),
						Priority: 5,
					},
					{
						ID:          "run-3km",
						Title:       "Run 3km",
						Description: "Do short distance running daily.",
						Emoji:       "🏃‍️️",
						Tags:        []string{"sport", "workout", "run", "gym", "fitness", "cardio", "distance"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     3,
							"default_change": 1,
							"is_infinite":    true,
						}),
						Priority: 5,
					},
					{
						ID:          "eat-an-apple",
						Title:       "Eat an apple",
						Description: "Because an apple a day keeps doctor away!",
						Emoji:       "🍏",
						Tags:        []string{"eat", "food", "apple", "health"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     1,
							"default_change": 1,
						}),
						Priority: 3,
					},
					{
						ID:          "wake-up-at-6am",
						Title:       "Wake up at 6am",
						Description: "Start waking up earlier than you normally do.",
						Emoji:       "⏰",
						Tags:        []string{"wake", "health", "sleep", "time", "bed"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     1,
							"default_change": 1,
						}),
						Priority: 10,
					},
					{
						ID:          "sleep-at-11pm",
						Title:       "Sleep at 11pm",
						Description: "Start going to bed earlier than you normally do.",
						Emoji:       "⏰",
						Tags:        []string{"wake", "health", "sleep", "time", "bed"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     1,
							"default_change": 1,
						}),
						Priority: 5,
					},
					{
						ID:          "play-music",
						Title:       "Play guitar",
						Description: "Practice guitar regularly.",
						Emoji:       "🎸",
						Tags:        []string{"play", "music", "guitar", "education", "learn", "practice"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     1,
							"default_change": 1,
						}),
						Priority: 5,
					},
					{
						ID:          "learn-german",
						Title:       "Learn German",
						Description: "Because the best way to learn a language is by practicing daily.",
						Emoji:       "🗣",
						Tags:        []string{"learn", "language", "education", "practice", "speaking", "german"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     1,
							"default_change": 1,
						}),
						Priority: 10,
					},
					{
						ID:          "learn-german-30min",
						Title:       "Spend 30 min learning German",
						Description: "Challenge yourself by practicing a language more than you normally do.",
						Emoji:       "🗣",
						Tags:        []string{"learn", "language", "education", "practice", "speaking", "german"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     30,
							"default_change": 10,
							"is_infinite":    true,
						}),
						Priority: 5,
					},
					{
						ID:          "no-smoking",
						Title:       "No smoking",
						Description: "Because zero cigarettes a day keeps doctor away!",
						Emoji:       "🚭",
						Tags:        []string{"smoke", "smoking", "health", "cigarette", "cigar", "weed", "joint"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  1,
							"goal_value":     0,
							"default_change": 1,
						}),
						Priority: 1,
					},
					{
						ID:          "eat-breakfast",
						Title:       "Eat breakfast",
						Description: "Don't skip breakfast to stay healthy.",
						Emoji:       "🍳",
						Tags:        []string{"eat", "breakfast", "food", "cook", "health"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     1,
							"default_change": 1,
						}),
						Priority: 10,
					},
					{
						ID:          "brush-teeth-twice-per-day",
						Title:       "Brush teeth twice per day",
						Description: "Because it's important for your teeth!",
						Emoji:       "🦷",
						Tags:        []string{"teeth", "tooth", "brush", "hygiene", "health"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  0,
							"goal_value":     2,
							"default_change": 1,
						}),
						Priority: 5,
					},
					{
						ID:          "dont-order-food",
						Title:       "Don't order food",
						Description: "No Uber Eats and Deliveroo, only cooking.",
						Emoji:       "🍕",
						Tags:        []string{"order", "delivery", "food", "groceries", "health", "eat"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  1,
							"goal_value":     0,
							"default_change": 1,
						}),
						Priority: 3,
					},
					{
						ID:          "no-alcohol",
						Title:       "No alcohol",
						Description: "Because life is brighter without alcohol.",
						Emoji:       "🍻",
						Tags:        []string{"drinks", "alcohol", "health"},
						Data: jsonMarshal(map[string]interface{}{
							"default_value":  1,
							"goal_value":     0,
							"default_change": 1,
						}),
						Priority: 1,
					},
				}

				return tx.Create(&trackers).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Delete(models.SampleTracker{}, []string{
					"cook-healthy-food",
					"read-10-pages",
					"50-pushups",
					"meditate-for-5-minutes",
					"do-yoga",
					"run-3km",
					"eat-an-apple",
					"wake-up-at-6am",
					"sleep-at-11pm",
					"play-music",
					"learn-german",
					"learn-german-30min",
					"no-smoking",
					"eat-breakfast",
					"brush-teeth-twice-per-day",
					"dont-order-food",
					"no-alcohol",
				}).Error
			},
		},
		{
			ID: "20210121-1",
			Migrate: func(tx *gorm.DB) error {
				err := tx.Delete(models.SampleTracker{}, []string{
					"meditate-for-5-minutes",
					"learn-german-30min",
				}).Error

				if err != nil {
					return err
				}

				trackers := []models.SampleTracker{
					{
						ID:          "meditate-for-5-minutes",
						Title:       "Meditate for 5 minutes",
						Description: "Try to establish a meditation routine.",
						Emoji:       "🧘",
						Tags:        []string{"health", "meditation", "yoga", "sport", "gym", "fitness", "timer"},
						Data: jsonMarshal(map[string]interface{}{
							"type":           "timer",
							"default_value":  0,
							"goal_value":     5 * 60,
							"default_change": 1,
							"is_infinite":    true,
						}),
						Priority: 5,
					},
					{
						ID:          "learn-german-30min",
						Title:       "Spend 30 min learning German",
						Description: "Challenge yourself by practicing a language more than you normally do.",
						Emoji:       "🗣",
						Tags:        []string{"learn", "language", "education", "practice", "speaking", "german", "timer"},
						Data: jsonMarshal(map[string]interface{}{
							"type":           "timer",
							"default_value":  0,
							"goal_value":     30 * 60,
							"default_change": 1,
							"is_infinite":    true,
						}),
						Priority: 5,
					},
					{
						ID:          "focused-work-60min",
						Title:       "60 min of focused work",
						Description: "Focus on work for at least 60 minutes.",
						Emoji:       "⌛️",
						Tags:        []string{"work", "build", "make", "focus", "timer"},
						Data: jsonMarshal(map[string]interface{}{
							"type":           "timer",
							"default_value":  0,
							"goal_value":     60 * 60,
							"default_change": 1,
							"is_infinite":    true,
						}),
						Priority: 10,
					},
					{
						ID:          "focused-studies-60min",
						Title:       "60 min of focused studies",
						Description: "Focus on studies for at least 60 minutes.",
						Emoji:       "⌛️",
						Tags:        []string{"learn", "study", "education", "focus", "timer"},
						Data: jsonMarshal(map[string]interface{}{
							"type":           "timer",
							"default_value":  0,
							"goal_value":     60 * 60,
							"default_change": 1,
							"is_infinite":    true,
						}),
						Priority: 10,
					},
					{
						ID:          "read-10min",
						Title:       "Read for 10 minutes",
						Description: "Carve out some time to read in order to establish the routine.",
						Emoji:       "️📚",
						Tags:        []string{"read", "books", "focus", "timer"},
						Data: jsonMarshal(map[string]interface{}{
							"type":           "timer",
							"default_value":  0,
							"goal_value":     10 * 60,
							"default_change": 1,
							"is_infinite":    true,
						}),
						Priority: 5,
					},
				}

				return tx.Create(&trackers).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Delete(models.SampleTracker{}, []string{
					"meditate-for-5-minutes",
					"learn-german-30min",
					"focused-work-60min",
					"focused-studies-60min",
					"read-10min",
				}).Error
			},
		},
	}
}

type SampleTrackerRepository struct {
	*TransactionProvider
}

func (r *SampleTrackerRepository) FindAllBySimilarTag(tag string, limit int) ([]models.SampleTracker, error) {
	var sampleTrackers []models.SampleTracker

	operation := r.DBConn.Raw(`
		SELECT * FROM sample_trackers st INNER JOIN (
			SELECT DISTINCT(id), priority FROM sample_trackers, unnest(tags) tag 
			WHERE lower(tag) LIKE ? 
			ORDER BY priority DESC
			LIMIT ?
		) st2 ON st.id = st2.id;
	`, "%"+tag+"%", limit).Find(&sampleTrackers)

	return sampleTrackers, operation.Error
}
