import ConfigForm from "./components/ConfigForm";
import { toast, ToastContainer } from "react-toastify";

const TOAST_CLOSE_DELAY_MS = 5000;

function App() {
  return (
    <>
      <div className="container-fluid">
        <ConfigForm />
      </div>
      <ToastContainer
        position={toast.POSITION.TOP_RIGHT}
        autoClose={TOAST_CLOSE_DELAY_MS}
        hideProgressBar
        newestOnTop
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover={false}
      />
    </>
  );
}

export default App;
