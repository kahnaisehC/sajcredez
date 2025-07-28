package board

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	type testCase struct {
		s                   Sajcredez
		inputMoves          []Move
		expectedBoardString string
	}

	runCases := []testCase{
		{
			s:                   CreateSajcredez(),
			expectedBoardString: "whiteEnhances: 0 blackEnhances: 0\ncastlingRights: KQkq\nturn: WHITE\n\t1\t2\t3\t4\t5\t6\t7\t\n1\t♖\t♘\t♗\t♔\t♗\t♘\t♖\t\n2\t♙\t♙\t♙\t♙\t♙\t♙\t♙\t\n3\t_\t_\t_\t_\t_\t_\t_\t\n4\t_\t_\t_\t_\t_\t_\t_\t\n5\t_\t_\t_\t_\t_\t_\t_\t\n6\t♟\t♟\t♟\t♟\t♟\t♟\t♟\t\n7\t♜\t♞\t♝\t♚\t♝\t♞\t♜\t\n",
		},
	}

	for _, test := range runCases {
		output := test.s.GetBoardString()
		if output != test.expectedBoardString {
			t.Errorf(`---------------------
Input moves: (%v)
Expecting: %s 
Actual: %s 
FAIL
				`, "", test.expectedBoardString, output)
		} else {
			fmt.Printf(`---------------------
Input moves: (%v)
Expecting: %s 
Actual: %s 
PASS
				`, "", test.expectedBoardString, output)
		}
	}
}

/*
white
	_
	♙
	♘
	♗
	♖
	♔
black
	♟
	♞
	♝
	♜
	♚

\t1\t2\t3\t4\t5\t6\t7\n1\t♖\t♘\t♗\t♔\t♗\t♘\t♖\n2\t♙\t♙\t♙\t♙\t♙\t♙\t♙\n3\t_\t_\t_\t_\t_\t_\t_\n4\t_\t_\t_\t_\t_\t_\t_\n5\t_\t_\t_\t_\t_\t_\t_\n6\t♜\t♞\t♝\t♚\t♝\t♞\t♜\n7\t♟\t♟\t♟\t♟\t♟\t♟\t♟\n
*/
