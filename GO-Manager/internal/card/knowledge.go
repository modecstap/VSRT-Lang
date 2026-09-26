package card

import (
	"errors"
	"math"
	"time"
)

const (
	initialEaseFactor = 2.5
	minEaseFactor     = 1.3
)

var (
	ErrUnknownEstimation = errors.New("unknown estimation")
	ErrCardNotFound      = errors.New("card not found")
)

type Estimation string

const (
	EstimationRepeat    Estimation = "REPEAT"
	EstimationDifficult Estimation = "DIFFICULT"
	EstimationEasy      Estimation = "EASY"
	EstimationMomental  Estimation = "MOMENTAL"
)

func (e Estimation) quality() (int, error) {
	switch e {
	case EstimationRepeat:
		return 0, nil
	case EstimationDifficult:
		return 3, nil
	case EstimationEasy:
		return 4, nil
	case EstimationMomental:
		return 5, nil
	default:
		return 0, ErrUnknownEstimation
	}
}

type KnowledgeId struct {
	SessionID int64
	PhraseKey string
}

type Knowledge struct {
	ID           KnowledgeId
	EaseFactor   float64
	Repetitions  int
	IntervalDays int
	DueAt        time.Time
}

func NewKnowledge(id KnowledgeId) *Knowledge {
	return &Knowledge{
		ID:         id,
		EaseFactor: initialEaseFactor,
	}
}

func (k *Knowledge) Estimate(e Estimation, now time.Time) error {
	q, err := e.quality()
	if err != nil {
		return err
	}

	miss := float64(5 - q)
	k.EaseFactor = math.Max(k.EaseFactor+(0.1-miss*(0.08+miss*0.02)), minEaseFactor)

	if q < 3 {
		k.Repetitions = 0
		k.IntervalDays = 0
		k.DueAt = now
		return nil
	}

	switch k.Repetitions {
	case 0:
		k.IntervalDays = 1
	case 1:
		k.IntervalDays = 6
	default:
		k.IntervalDays = int(math.Round(float64(k.IntervalDays) * k.EaseFactor))
	}
	k.Repetitions++
	k.DueAt = now.AddDate(0, 0, k.IntervalDays)
	return nil
}

func (k Knowledge) DueBy(now time.Time) bool {
	return !k.DueAt.After(now)
}
