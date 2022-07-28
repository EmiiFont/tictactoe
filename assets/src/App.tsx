import React, { useEffect, useState } from "react";
import "./App.css";
import { connect, sendMsg } from "./websocket";
import { Link, Outlet } from "react-router-dom";

function BottomWidget(props: any) {
  return (
    <div className="pt-8 text-base font-semibold leading-7">
      <button
        onClick={props.newGame}
        className="bg-sky-500 hover:bg-sky-700 px-5 py-2.5
        text-sm leading-5 rounded-md font-semibold text-white">
        New game
      </button>
    </div>
  );
}

function App() {
  const [data, setData] = useState({
    winner: { winner: 45, positions: [] },
    board: [],
    movesLeft: true,
    roomId: "",
    turn: "",
  });
  const [trColor, setTrColor] = useState({ color: "" });

  useEffect(() => {
    connect();
    const fetchData = async () => {
      try {
        const response = await fetch("http://localhost:8080", { mode: "cors" });
        console.log(response);
        const json = await response.json();
        setData(json);
        console.log(json);
      } catch (error) {
        console.log("error", error);
      }
    };
    fetchData();
  }, []);
  function send() {
    sendMsg("Hello");
  }
  function mapSymbol(ascii: number): string {
    return String.fromCharCode(ascii);
  }

  async function newGame() {
    try {
      const requestOptions: any = {
        method: "GET",
        headers: { "Content-Type": "application/json" },
      };
      const response = await fetch(
        "http://localhost:8080/newGame",
        requestOptions
      );
      const json = await response.json();
      window.location.href = `http://localhost:8080/ws/${json}`;
    } catch (error) {
      console.log("error", error);
    }
  }
  async function handleClick(col: number, row: number) {
    console.log(`user clicked (${col}, ${row}, ${data.roomId})`);
    sendMsg(`user clicked (${col}, ${row}, ${data.roomId})`);
    try {
      const requestOptions: any = {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          player: mapSymbol(+data.turn),
          column: col,
          row: row,
          roomId: data.roomId,
        }),
      };
      const response = await fetch(
        "http://localhost:8080/move",
        requestOptions
      );
      const json = await response.json();
      console.log(json);
      setData(json);
    } catch (error) {
      console.log("error", error);
    }
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
    <div className="App">
      <div className="relative flex min-h-screen flex-col justify-center overflow-hidden bg-gray-50 py-6 bg:py-12">
        <img
          src=""
          alt=""
          className="absolute top-1/2 left-1/2 max-w-none -translate-x-1/2 -translate-y-1/2"
          width="1308"
        />
        <div className="absolute inset-0 bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]"></div>
        <div className="relative bg-white px-6 pt-10 pb-8 shadow-xl ring-1 ring-gray-900/5 sm:mx-auto sm:max-w-lg sm:rounded-lg sm:px-10">
          <div className="max-w-md">
            <p className="h-6"> Tic Tac Toe </p>
            <div className="divide-y divide-gray-300/50">
              <div className="space-y-14 py-14 text-gray-600">
                <table className="table-fixed">
                  <tbody>
                    {data.board.map((row: [], i) => (
                      <tr key={i} style={trColor}>
                        {row.map((col, j) => (
                          <td
                            key={j}
                            className="border w-40 h-40 text-black-700 text-center text-4xl"
                            style={colorWinner(data.winner, i, j)}
                            onClick={(e) => handleClick(i, j)}>
                            {data.board[i][j] === 95
                              ? " "
                              : mapSymbol(data.board[i][j])}
                          </td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              <Link to="/single">Single game</Link> | <BottomWidget />
            </div>
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
