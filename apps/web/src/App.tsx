import { Routes, Route } from 'react-router-dom'
import HomePage from './pages/HomePage'
import TopicPage from './pages/TopicPage'
import { BTreePage } from './components/btree'
import { MVCCPage } from './components/mvcc'
import { QueryParserPage } from './components/query-parser'

function App() {
  return (
    <div className="app">
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/topic/:id" element={<TopicPage />} />
        <Route path="/btree" element={<BTreePage />} />
        <Route path="/mvcc" element={<MVCCPage />} />
        <Route path="/query-parser" element={<QueryParserPage />} />
      </Routes>
    </div>
  )
}

export default App
