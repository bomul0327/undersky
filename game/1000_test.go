package game_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hellodhlyn/undersky/game"
)

var _ = Describe("TicTacToeBoard", func() {
	Describe("#GetInputText()", func() {
		It("", func() {
			board := game.NewTicTacToeBoard()
			board.Set("A2", int8(1))
			board.Set("B3", int8(2))

			Expect(board.GetInputText()).To(Equal("000\n100\n020"))
		})
	})

	Describe("#FindWinner()", func() {
		var board *game.TicTacToeBoard

		BeforeEach(func() {
			board = game.NewTicTacToeBoard()
		})

		Context("No winner", func() {
			It("before actions", func() {
				Expect(board.FindWinner()).To(Equal(int8(0)))
			})

			It("after actions", func() {
				// 112
				// 221
				// 122
				board.Set("A1", int8(1))
				board.Set("A2", int8(2))
				board.Set("A3", int8(1))
				board.Set("B1", int8(1))
				board.Set("B2", int8(2))
				board.Set("B3", int8(2))
				board.Set("C1", int8(2))
				board.Set("C2", int8(1))
				board.Set("C3", int8(2))

				Expect(board.FindWinner()).To(Equal(int8(0)))
			})
		})

		Context("Winner exists", func() {
			It("horizontal", func() {
				// 110
				// 222
				// 120
				board.Set("A1", int8(1))
				board.Set("A2", int8(2))
				board.Set("A3", int8(1))
				board.Set("B1", int8(1))
				board.Set("B2", int8(2))
				board.Set("B3", int8(2))
				board.Set("C1", int8(0))
				board.Set("C2", int8(2))
				board.Set("C3", int8(0))

				Expect(board.FindWinner()).To(Equal(int8(2)))
			})

			It("vertical", func() {
				// 111
				// 122
				// 220
				board.Set("A1", int8(1))
				board.Set("A2", int8(1))
				board.Set("A3", int8(2))
				board.Set("B1", int8(1))
				board.Set("B2", int8(2))
				board.Set("B3", int8(2))
				board.Set("C1", int8(1))
				board.Set("C2", int8(2))
				board.Set("C3", int8(0))

				Expect(board.FindWinner()).To(Equal(int8(1)))
			})

			It("diagonal", func() {
				// 102
				// 210
				// 121
				board.Set("A1", int8(1))
				board.Set("A2", int8(2))
				board.Set("A3", int8(1))
				board.Set("B1", int8(0))
				board.Set("B2", int8(1))
				board.Set("B3", int8(2))
				board.Set("C1", int8(2))
				board.Set("C2", int8(0))
				board.Set("C3", int8(1))

				Expect(board.FindWinner()).To(Equal(int8(1)))
			})
		})
	})
})
