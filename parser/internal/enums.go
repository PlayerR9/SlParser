package internal

// ActionType is the type of an action.
type ActionType uint

const (
	// ActAccept is the action type for accepting the parse.
	ActAccept ActionType = iota

	// ActReduce is the action type for reducing the parse.
	ActReduce

	// ActShift is the action type for shifting the parse.
	ActShift
)

// PhaseType is the type of a phase.
type PhaseType int

const (
	// PhaseReduction is the phase type for reduction.
	PhaseReduction PhaseType = iota

	// PhasePrediction is the phase type for prediction.
	PhasePrediction

	// PhaseShifting is the phase type for shifting.
	PhaseShifting

	// PhaseCheckBranch is the phase type for checking branch.
	PhaseCheckBranch
)
