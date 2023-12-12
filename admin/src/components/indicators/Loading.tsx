import React from "react";
//import Spinner from "react-bootstrap/Spinner";
import { PuffLoader } from "react-spinners";

interface LoadingProps {
  height?: string;
}

const Loading: React.FC<LoadingProps> = ({ height }) => {
  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: height !== undefined ? height : "100vh",
      }}
    >
      {/*<Spinner animation="grow" />*/}
      <PuffLoader color="#8a8a8a" />
    </div>
  );
};

export default Loading;
