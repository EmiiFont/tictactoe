import React, { useState } from "react";
import Board from "../components/board";

function singleGame() {
  const [data, setData] = useState({
    winner: { winner: 45, positions: [] },
    board: [],
    movesLeft: true,
    roomId: "",
    turn: "",
  });
  return (
    <main style={{ padding: "1rem 0" }}>
      <h1>Single game</h1>
      <Board></Board>
    </main>
  );
}

export default singleGame;
