import React from "react";
import { Outlet } from "react-router-dom";

function layout(props: any) {
  return (
    <div className="App">
      <div className="relative flex min-h-screen flex-col justify-center overflow-hidden bg-gray-50 py-6 bg:py-12">
        <img
          src=""
          alt=""
          className="absolute top-1/2 left-1/2 max-w-none -translate-x-1/2 -translate-y-1/2"
          width="1308"
        />
        <p className="h-6"> Tic Tac Toe </p>
        <div className="absolute inset-0 bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))]"></div>
        <div className="relative bg-white px-6 pt-10 pb-8 shadow-xl ring-1 ring-gray-900/5 sm:mx-auto sm:max-w-lg sm:rounded-lg sm:px-10">
          <div className="max-w-md">
            <div className="divide-y divide-gray-300/50">
              <div className="space-x-16 px-16 py-12 text-gray-600">
                <div className="grid grid-rows-3 grid-flow-col">
                  <Outlet></Outlet>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default layout;
