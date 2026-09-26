package test

import (
	"errors"
	"math"
	"slices"
	"testing"
	"time"

	"VSRT-Lang/internal/card"
	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/translators/stub"
	"VSRT-Lang/internal/user"
)

const (
	userAnna  = user.UserId("anna")
	userOther = user.UserId("other")
)

type fixture struct {
	t         *testing.T
	sessions  *session.Service
	cards     *card.Service
	knowledge *memory.KnowledgeRepository
	clock     time.Time
	sessionID int64
}

func setup(t *testing.T) *fixture {
	t.Helper()
	repo := memory.NewSessionRepository()
	f := &fixture{
		t:         t,
		sessions:  session.NewService(repo, stub.Translator{}),
		knowledge: memory.NewKnowledgeRepository(),
		clock:     time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	f.cards = card.NewService(repo, f.knowledge, func() time.Time { return f.clock })
	f.sessionID = f.newSession(userAnna, "Урок 1")
	return f
}

func (f *fixture) newSession(uid user.UserId, name string) int64 {
	f.t.Helper()
	id, err := f.sessions.NewSession(uid, name)
	if err != nil {
		f.t.Fatalf("NewSession(%q): %v", name, err)
	}
	return id
}

func (f *fixture) add(phrases ...string) {
	f.t.Helper()
	for _, p := range phrases {
		cmd := session.AddRecordCommand{UserId: userAnna, SessionId: f.sessionID, Phrase: p}
		if _, err := f.sessions.AddRecord(cmd); err != nil {
			f.t.Fatalf("AddRecord(%q): %v", p, err)
		}
	}
}

func (f *fixture) remove(phrase string) {
	f.t.Helper()
	if err := f.sessions.DeleteRecord(userAnna, f.sessionID, phrase); err != nil {
		f.t.Fatalf("DeleteRecord(%q): %v", phrase, err)
	}
}

func (f *fixture) estimate(phrase string, e card.Estimation) error {
	return f.cards.UpdateKnowledge(card.UpdateKnowledgeCommand{
		UserId: userAnna, SessionId: f.sessionID, Phrase: phrase, Estimation: e,
	})
}

func (f *fixture) mustEstimate(phrase string, e card.Estimation) {
	f.t.Helper()
	if err := f.estimate(phrase, e); err != nil {
		f.t.Fatalf("UpdateKnowledge(%q, %q): %v", phrase, e, err)
	}
}

func (f *fixture) due(limit int) []string {
	f.t.Helper()
	records, err := f.cards.GetCards(userAnna, f.sessionID, limit)
	if err != nil {
		f.t.Fatalf("GetCards: %v", err)
	}
	if records == nil {
		f.t.Fatal("GetCards returned nil slice")
	}
	phrases := make([]string, 0, len(records))
	for _, r := range records {
		phrases = append(phrases, r.Phrase)
	}
	return phrases
}

func (f *fixture) state(phrase string) *card.Knowledge {
	f.t.Helper()
	k, err := f.knowledge.Take(card.KnowledgeId{SessionID: f.sessionID, PhraseKey: session.RecordKey(phrase)})
	if err != nil {
		f.t.Fatalf("Take(%q): %v", phrase, err)
	}
	return k
}

func assertEF(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("ease factor = %v, want %v", got, want)
	}
}

func TestNewRecordIsDueImmediately(t *testing.T) {
	f := setup(t)
	f.add("hello")

	if got := f.due(0); !slices.Equal(got, []string{"hello"}) {
		t.Fatalf("cards = %v, want [hello]", got)
	}
}

func TestNewCardsFirstThenByDue(t *testing.T) {
	f := setup(t)
	f.add("pear", "kiwi", "apple", "fig")
	f.mustEstimate("pear", card.EstimationRepeat)
	f.clock = f.clock.Add(time.Minute)
	f.mustEstimate("apple", card.EstimationRepeat)

	want := []string{"fig", "kiwi", "pear", "apple"}
	if got := f.due(0); !slices.Equal(got, want) {
		t.Fatalf("cards = %v, want %v", got, want)
	}
}

func TestEasyIntervals(t *testing.T) {
	f := setup(t)
	f.add("hello")

	for i, interval := range []int{1, 6, 15} {
		f.mustEstimate("hello", card.EstimationEasy)
		k := f.state("hello")
		if k.IntervalDays != interval || k.Repetitions != i+1 {
			t.Fatalf("step %d: interval %d, repetitions %d; want %d, %d", i, k.IntervalDays, k.Repetitions, interval, i+1)
		}
		assertEF(t, k.EaseFactor, 2.5)
		if !k.DueAt.Equal(f.clock.AddDate(0, 0, interval)) {
			t.Fatalf("step %d: due %v, want now + %d days", i, k.DueAt, interval)
		}

		f.clock = k.DueAt.Add(-time.Second)
		if got := f.due(0); len(got) != 0 {
			t.Fatalf("step %d: card due before its time: %v", i, got)
		}
		f.clock = k.DueAt
		if got := f.due(0); !slices.Equal(got, []string{"hello"}) {
			t.Fatalf("step %d: card missing at due: %v", i, got)
		}
	}
}

func TestEaseFactorByEstimation(t *testing.T) {
	tests := []struct {
		estimation card.Estimation
		ef         float64
	}{
		{card.EstimationRepeat, 1.7},
		{card.EstimationDifficult, 2.36},
		{card.EstimationEasy, 2.5},
		{card.EstimationMomental, 2.6},
	}
	for _, tt := range tests {
		t.Run(string(tt.estimation), func(t *testing.T) {
			f := setup(t)
			f.add("hello")
			f.mustEstimate("hello", tt.estimation)
			assertEF(t, f.state("hello").EaseFactor, tt.ef)
		})
	}
}

func TestRepeatResetsAndKeepsCardDue(t *testing.T) {
	f := setup(t)
	f.add("hello")
	f.mustEstimate("hello", card.EstimationEasy)
	f.clock = f.clock.AddDate(0, 0, 1)
	f.mustEstimate("hello", card.EstimationRepeat)

	k := f.state("hello")
	if k.Repetitions != 0 || k.IntervalDays != 0 {
		t.Fatalf("repetitions %d, interval %d; want 0, 0", k.Repetitions, k.IntervalDays)
	}
	if got := f.due(0); !slices.Equal(got, []string{"hello"}) {
		t.Fatalf("cards = %v, want [hello]", got)
	}

	f.mustEstimate("hello", card.EstimationRepeat)
	f.mustEstimate("hello", card.EstimationRepeat)
	assertEF(t, f.state("hello").EaseFactor, 1.3)
}

func TestEstimateNormalizesPhrase(t *testing.T) {
	f := setup(t)
	f.add("hello")
	f.mustEstimate("  HELLO ", card.EstimationEasy)

	if k := f.state("hello"); k == nil || k.Repetitions != 1 {
		t.Fatalf("knowledge = %+v, want one repetition", k)
	}
}

func TestEstimateErrorsLeaveKnowledgeUntouched(t *testing.T) {
	tests := []struct {
		name       string
		phrase     string
		estimation card.Estimation
		want       error
	}{
		{"missing phrase", "world", card.EstimationEasy, card.ErrCardNotFound},
		{"lowercase estimation", "hello", "easy", card.ErrUnknownEstimation},
		{"empty estimation", "hello", "", card.ErrUnknownEstimation},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := setup(t)
			f.add("hello")

			if err := f.estimate(tt.phrase, tt.estimation); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
			if k := f.state(tt.phrase); k != nil {
				t.Fatalf("knowledge created: %+v", k)
			}

			f.mustEstimate("hello", card.EstimationEasy)
			before := *f.state("hello")
			if err := f.estimate(tt.phrase, tt.estimation); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
			if after := *f.state("hello"); after != before {
				t.Fatalf("knowledge changed: %+v, want %+v", after, before)
			}
		})
	}
}

func TestForeignSessionNotFound(t *testing.T) {
	f := setup(t)
	foreign := f.newSession(userOther, "Чужой урок")

	if _, err := f.cards.GetCards(userAnna, foreign, 0); !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf("GetCards error = %v, want %v", err, session.ErrSessionNotFound)
	}
	err := f.cards.UpdateKnowledge(card.UpdateKnowledgeCommand{
		UserId: userAnna, SessionId: foreign, Phrase: "hello", Estimation: card.EstimationEasy,
	})
	if !errors.Is(err, session.ErrSessionNotFound) {
		t.Fatalf("UpdateKnowledge error = %v, want %v", err, session.ErrSessionNotFound)
	}
}

func TestLimitAndEmptyResult(t *testing.T) {
	f := setup(t)
	f.add("c", "a", "b")

	if got := f.due(2); !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("cards = %v, want [a b]", got)
	}

	for _, p := range []string{"a", "b", "c"} {
		f.mustEstimate(p, card.EstimationEasy)
	}
	if got := f.due(0); len(got) != 0 {
		t.Fatalf("cards = %v, want empty", got)
	}
}

func TestKnowledgeSurvivesRecordChanges(t *testing.T) {
	f := setup(t)
	f.add("hello")
	f.mustEstimate("hello", card.EstimationEasy)
	before := *f.state("hello")

	f.add("world")
	f.remove("world")
	f.add("hello")
	if after := *f.state("hello"); after != before {
		t.Fatalf("knowledge changed: %+v, want %+v", after, before)
	}

	f.remove("hello")
	f.add("hello")
	f.clock = before.DueAt
	f.mustEstimate("hello", card.EstimationEasy)
	if k := f.state("hello"); k.Repetitions != 2 || k.IntervalDays != 6 {
		t.Fatalf("repetitions %d, interval %d; want 2, 6", k.Repetitions, k.IntervalDays)
	}
}
