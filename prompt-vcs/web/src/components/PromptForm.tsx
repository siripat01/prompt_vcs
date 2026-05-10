import { useState } from 'react'

interface Prompt {
  id: string
  name: string
  content: string
  tags: string[]
  collection: string
  version: number
}

interface PromptFormProps {
  onCreate: (prompt: Prompt) => void
}

function PromptForm({ onCreate }: PromptFormProps) {
  const [name, setName] = useState('')
  const [content, setContent] = useState('')
  const [tags, setTags] = useState('')
  const [collection, setCollection] = useState('default')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    const newPrompt = {
      id: Date.now().toString(),
      name,
      content,
      tags: tags.split(',').map(t => t.trim()).filter(t => t),
      collection,
      version: 1
    }

    // Send to API
    try {
      const response = await fetch('/api/prompts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newPrompt)
      })
      
      if (response.ok) {
        onCreate(newPrompt)
      }
    } catch (error) {
      console.error('Error creating prompt:', error)
    }
  }

  return (
    <div className="prompt-form">
      <h2>Create New Prompt</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Name:</label>
          <input 
            type="text" 
            value={name} 
            onChange={e => setName(e.target.value)} 
            required 
          />
        </div>
        
        <div className="form-group">
          <label>Content:</label>
          <textarea 
            value={content} 
            onChange={e => setContent(e.target.value)} 
            rows={10}
            required 
          />
        </div>

        <div className="form-group">
          <label>Tags (comma-separated):</label>
          <input 
            type="text" 
            value={tags} 
            onChange={e => setTags(e.target.value)} 
          />
        </div>

        <div className="form-group">
          <label>Collection:</label>
          <input 
            type="text" 
            value={collection} 
            onChange={e => setCollection(e.target.value)} 
          />
        </div>

        <button type="submit">Create Prompt</button>
      </form>
    </div>
  )
}

export default PromptForm