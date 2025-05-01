<script>
  import { onMount } from 'svelte';
  
  export let visible = false;
  
  function closeHelp() {
    visible = false;
  }
  
  function handleKeydown(event) {
    if (event.key === 'Escape' && visible) {
      closeHelp();
    }
  }
  
  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    
    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  });
</script>

{#if visible}
  <div class="keybindings-help">
    <button class="close-button" on:click={closeHelp}>×</button>
    <div class="help-content">
      <h2>Keyboard Shortcuts</h2>
      
      <h3>Navigation</h3>
      <div class="shortcuts-grid">
        <div class="shortcut">
          <span class="key">h</span>
          <span class="description">Scroll left</span>
        </div>
        <div class="shortcut">
          <span class="key">j</span>
          <span class="description">Scroll down</span>
        </div>
        <div class="shortcut">
          <span class="key">k</span>
          <span class="description">Scroll up</span>
        </div>
        <div class="shortcut">
          <span class="key">l</span>
          <span class="description">Scroll right</span>
        </div>
        <div class="shortcut">
          <span class="key">u</span>
          <span class="description">Page up</span>
        </div>
        <div class="shortcut">
          <span class="key">d</span>
          <span class="description">Page down</span>
        </div>
      </div>
      
      <h3>Link Hints</h3>
      <div class="shortcuts-grid">
        <div class="shortcut">
          <span class="key">f</span>
          <span class="description">Show link hints</span>
        </div>
        <div class="shortcut">
          <span class="key">Esc</span>
          <span class="description">Exit link selection mode</span>
        </div>
      </div>
      
      <h3>Other</h3>
      <div class="shortcuts-grid">
        <div class="shortcut">
          <span class="key">?</span>
          <span class="description">Show this help dialog</span>
        </div>
      </div>
      
      <p class="help-note">
        Note: Keyboard shortcuts don't work when typing in form fields.
      </p>
    </div>
  </div>
{/if}

<style>
  .keybindings-help {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: var(--card-bg, white);
    border: 1px solid var(--border-color, #e2e8f0);
    border-radius: 6px;
    padding: 20px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 9999;
    max-width: 500px;
    width: 90%;
    max-height: 90vh;
    overflow-y: auto;
    animation: fadeIn 0.3s ease-out;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
  
  .close-button {
    position: absolute;
    top: 10px;
    right: 10px;
    background: transparent;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: var(--text-muted, #718096);
    padding: 0;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
  }
  
  .close-button:hover {
    background: rgba(0, 0, 0, 0.05);
    color: var(--text-color, #333);
  }
  
  h2 {
    margin: 0 0 20px 0;
    font-size: 24px;
    color: var(--heading-color, #2d3748);
    text-align: center;
  }
  
  h3 {
    margin: 20px 0 10px 0;
    font-size: 18px;
    color: var(--heading-color, #2d3748);
    border-bottom: 1px solid var(--border-color, #e2e8f0);
    padding-bottom: 5px;
  }
  
  .shortcuts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 10px;
    margin-bottom: 15px;
  }
  
  .shortcut {
    display: flex;
    align-items: center;
    padding: 5px 0;
  }
  
  .key {
    display: inline-block;
    background: var(--secondary-button, #4a5568);
    color: white;
    border-radius: 4px;
    padding: 2px 8px;
    margin-right: 10px;
    font-family: monospace;
    font-size: 14px;
    min-width: 20px;
    text-align: center;
  }
  
  .description {
    color: var(--text-color, #333);
    font-size: 14px;
  }
  
  .help-note {
    font-size: 13px;
    color: var(--text-secondary, #666);
    margin-top: 20px;
    border-top: 1px solid var(--border-color, #e2e8f0);
    padding-top: 10px;
    font-style: italic;
  }
  
  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
</style> 