import { useState, useEffect } from 'react'
import PromptForm from './components/PromptForm'
import PromptList from './components/PromptList'
import PromptDetail from './components/PromptDetail'
import './App.css'

interface Prompt {
  id: string
  name: string
  content: string
  tags: string[]
  collection: string
  version: number
}

function App() {
  const [prompts, setPrompts] = useState<Prompt[]>([])
  const [selectedPrompt, setSelectedPrompt] = useState<Prompt | null>(null)
  const [view, setView] = useState<'list' | 'detail' | 'create'>('list')

  useEffect(() => {
    // Fetch prompts from API
    fetch('/api/prompts')
      .then(res => res.json())
      .then(data => setPrompts(data))
      .catch(err => console.error('Error fetching prompts:', err))
  }, [])

  const handleCreate = (newPrompt: Prompt) => {
    setPrompts(prev => [...prev, newPrompt])
    setView('list')
  }

  const handleSelect = (prompt: Prompt) => {
    setSelectedPrompt(prompt)
    setView('detail')
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>Prompt Version Control</h1>
        <nav>
          <button onClick={() => setView('list')}>All Prompts</button>
          <button onClick={() => setView('create')}>New Prompt</button>
        </nav>
      </header>

      <main>
        {view === 'list' && (
          <PromptList prompts={prompts} onSelect={handleSelect} />
        )}
        {view === 'create' && (
          <PromptForm onCreate={handleCreate} />
        )}
        {view === 'detail' && selectedPrompt && (
          <PromptDetail prompt={selectedPrompt} onBack={() => setView('list')} />
        )}
      </main>
    </div>
  )
}

export default App