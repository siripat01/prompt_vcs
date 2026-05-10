import './PromptList.css'

interface Prompt {
  id: string
  name: string
  content: string
  tags: string[]
  collection: string
  version: number
}

interface PromptListProps {
  prompts: Prompt[]
  onSelect: (prompt: Prompt) => void
}

function PromptList({ prompts, onSelect }: PromptListProps) {
  return (
    <div className="prompt-list">
      <h2>All Prompts</h2>
      {prompts.length === 0 ? (
        <p>No prompts found. Create your first prompt!</p>
      ) : (
        <div className="prompt-grid">
          {prompts.map(prompt => (
            <div 
              key={prompt.id} 
              className="prompt-card"
              onClick={() => onSelect(prompt)}
            >
              <h3>{prompt.name}</h3>
              <p className="prompt-content">{prompt.content.substring(0, 100)}...</p>
              <div className="prompt-meta">
                <span className="collection">{prompt.collection}</span>
                <span className="version">v{prompt.version}</span>
              </div>
              <div className="tags">
                {prompt.tags.map(tag => (
                  <span key={tag} className="tag">{tag}</span>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default PromptList