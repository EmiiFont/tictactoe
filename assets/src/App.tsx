import React, {useEffect, useState} from 'react';
import logo from './logo.svg';
import './App.css';

function App() {
  const [data, setData] = useState({ winner: { winner: 45, positions: [] }, board: [], movesLeft: true, roomId: "", turn: ""});
  const [trColor, setTrColor] = useState({ color: ''})

  useEffect( () => {
    const fetchData = async () => {
      try {
        const response = await fetch("http://localhost:8080",{mode:'cors'});
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
  function mapSymbol(ascii: number): string {
      return String.fromCharCode(ascii);
  }
  async function handleClick(col: number, row: number) {
      console.log(`user clicked (${col}, ${row}, ${data.roomId})`)
      try {
          const requestOptions: any = {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({
                  "player": mapSymbol(+data.turn),
                  "column": col,
                  "row": row,
                  "roomId": data.roomId
              })
          };
        const response = await fetch("http://localhost:8080/move", requestOptions);
        const json = await response.json();
        console.log(json);
          setData(json);
      } catch (error) {
        console.log("error", error);
      }
  }

  function colorWinner(winnerObj: any, col: number, row: number) {
      if (mapSymbol(winnerObj.winner) != '_' ) {
          for (let i=0; i < winnerObj.positions.length; i++) {
              if (winnerObj.positions[i][0] === col && winnerObj.positions[i][1] === row) {
                  return {color: 'green'}
              }
          }
      }
  }

  return (
    <div className="App">
      <header className="App-header">
          <div className="flex flex-col">
              <div className="overflow-x-auto sm:-mx-6 lg:-mx-8">
                  <div className="py-2 inline-block min-w-full sm:px-6 lg:px-8">
                      <div className="overflow-hidden">
                          <table className="hover:table-auto border-separate border-spacing-2 border border-slate-500 min-w-full">
                              <tbody>
                              {data.board.map((row: [], i) => (
                                  <tr key={i} style={trColor}>
                                      { row.map((col, j) => (
                                          <td key={j} className="border border-slate-600" style={colorWinner(data.winner, i,j)}
                                              onClick={(e) => handleClick(i, j)}>
                                              {
                                                  data.board[i][j] === 95 ? "_" :
                                                      mapSymbol(data.board[i][j])
                                              }
                                          </td>
                                      ))}
                                  </tr>
                              ))}
                              </tbody>
                          </table>
                      </div>
                  </div>
              </div>
          </div>
      </header>

    </div>
  );
}

export default App;
