import { BrowserRouter, Routes, Route } from 'react-router-dom';
import CommentPage from './pages/CommentPage';
import CommentThreadPage from './pages/CommentThreadPage';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<CommentPage />} />
        <Route path="/comment/:id" element={<CommentThreadPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;