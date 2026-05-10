import { useState } from 'react'

interface Prompt {
  id: string
  name: string
  content: string
  tags: string[]
  collection: string
  version: number
}

interface PromptDetailProps {
  prompt: Prompt
  onBack: () => void
}

function PromptDetail({ prompt, onBack }: PromptDetailProps) {
  const [commits, setCommits] = useState<string[]>([])
  const [showCommits, setShowCommits] = useState(false)

  const fetchCommits = async () => {
    try {
      const response = await fetch(`/api/prompts/${prompt.id}/commits`)
      if (response.ok) {
        const data = await response.json()
        setCommits(data)
      }
    } catch (error) {
      console.error('Error fetching commits:', error)
    }
  }

  const handleVersionControl = () => {
    setShowCommits(!showCommits)
    if (!showCommits) {
      fetchCommits()
    }
  }

  return (
    <div className="prompt-detail">
      <button onClick={onBack} className="back-button">← Back</button>
      
      <div className="detail-header">
        <h2>{prompt.name}</h2>
        <div className="detail-meta">
          <span>Collection: {prompt.collection}</span>
          <span>Version: {prompt.version}</span>
        </div>
      </div>

      <div className="detail-content">
        <h3>Content</h3>
        <pre>{prompt.content}</pre>
      </div>

      <div className="detail-tags">
        <h3>Tags</h3>
        {prompt.tags.map(tag => (
          <span key={tag} className="tag">{tag}</span>
        ))}
      </div>

      <div className="version-control">
        <button onClick={handleVersionControl}>
          {showCommits ? 'Hide' : 'Show'} Version History
        </button>
        
        {showCommits && (
          <div className="commits-list">
            <h3>Version History</h3>
            {commits.length === 0 ? (
              <p>No commits yet</p>
            ) : (
              <ul>
                {commits.map((commit, index) => (
                  <li key={index}>{commit}</li>
                ))}
              </ul>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

export default PromptDetail