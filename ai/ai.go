package ai

import (
	"fmt"

	"github.com/48thFlame/Checkers/checkers"
)

type moveEval struct {
	move checkers.Move
	eval int
}

func (me moveEval) String() string {
	return fmt.Sprintf("(%d,%d|%d)",
		me.move.StartI, me.move.EndI, me.eval)
}

func calculateAllMoves(g *checkers.Game, depth int) []moveEval {
	moveEvalsChannel := make(chan moveEval)

	legalMoves := g.GetLegalMoves()

	depth-- // because playing a move and then min-max

	for _, move := range legalMoves {
		futureGame := *g
		(&futureGame).PlayMove(move)

		go func(m checkers.Move) {
			eval := minMax(futureGame, depth, depth, lowestE, highestE)
			moveEvalsChannel <- moveEval{move: m, eval: eval}
			// moveEvals = append(moveEvals, moveEval{move: move, eval: eval})
		}(move)
	}

	moveEvals := make([]moveEval, 0)
	for i := 0; i < len(legalMoves); i++ {
		me := <-moveEvalsChannel

		moveEvals = append(moveEvals, me)
	}

	return moveEvals
}

func SmartAi(g checkers.Game) checkers.Move {
	moveEvals := calculateAllMoves(&g, 8)
	sortMoveEvalsHighToLow(moveEvals)

	var bestMoveEval moveEval

	if g.PlrTurn == checkers.BluePlayer {
		// take first move
		bestMoveEval = moveEvals[0]
	} else {
		// red wants lowest so take last
		bestMoveEval = moveEvals[len(moveEvals)-1]
	}

	fmt.Println(bestMoveEval)

	return bestMoveEval.move
}
