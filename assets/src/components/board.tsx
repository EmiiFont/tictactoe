import React, { useState } from "react";

export function board(props: any) {
  const [trColor, setTrColor] = useState({ color: "" });

  function mapSymbol(ascii: number): string {
    return String.fromCharCode(ascii);
  }

  function colorWinner(winnerObj: any, col: number, row: number) {
    if (mapSymbol(winnerObj.winner) != "_") {
      for (let i = 0; i < winnerObj.positions.length; i++) {
        if (
          winnerObj.positions[i][0] === col &&
          winnerObj.positions[i][1] === row
        ) {
          return { color: "green" };
        }
      }
    }
  }

  return (
    <table className="table-fixed">
      <tbody>
        {props.data.board.map((row: [], i: number) => (
          <tr key={i} style={trColor}>
            {row.map((col, j) => (
              <td
                key={j}
                className="border w-40 h-40 text-black-700 text-center text-4xl"
                style={colorWinner(props.data.winner, i, j)}
                onClick={(e) => props.handleClick(i, j)}>
                {props.data.board[i][j] === 95
                  ? " "
                  : mapSymbol(props.data.board[i][j])}
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
