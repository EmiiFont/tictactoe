import React, { useEffect, useState } from "react";
import "./App.css";
import { connect, sendMsg } from "./websocket";
import { Link, Outlet } from "react-router-dom";

function BottomWidget(props: any) {
  return (
    <div className="pt-8 text-base font-semibold leading-7">
      <Link
        className="bg-sky-500 hover:bg-sky-700 px-6 py-4
        text-sm leading-5 rounded-md font-semibold text-white"
        to={props.link}>
        {props.name}
      </Link>
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

  return (
    <>
      <div className="col-span-3">
        <BottomWidget link={"/single"} name={"Single game"} />
      </div>
      <div className="col-span-3">
        <BottomWidget link={"/multiplayer"} name={"Multiplayer game"} />
      </div>
    </>
  );
}

export default App;
