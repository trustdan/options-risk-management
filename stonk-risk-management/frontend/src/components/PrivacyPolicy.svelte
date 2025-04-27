<script>
  // Simple privacy policy component
  import { onMount, createEventDispatcher } from 'svelte';
  
  const dispatch = createEventDispatcher();
  
  function closeModal() {
    dispatch('close');
  }
  
  // Handle escape key to close the modal
  function handleKeydown(event) {
    if (event.key === 'Escape') {
      closeModal();
    }
  }
  
  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  });
</script>

<div class="privacy-container">
  <div class="privacy-content">
    <h1>Privacy Policy</h1>
    
    <div class="policy-text">
      <p>
        <strong>Options Risk Management does not collect any user data.</strong>
      </p>
      
      <p>
        All your data is stored locally on your device, and we do not track, collect, or share 
        any personal information or usage statistics.
      </p>
      
      <p>
        Your privacy is important to us. This application operates entirely offline and does not 
        transmit any data over the internet except when explicitly requested by you (such as when using the 
        "Buy Me a Coffee" button or accessing GitHub).
      </p>
    </div>
    
    <div class="action-buttons">
      <button class="close-button" on:click={closeModal}>
        Close
      </button>
    </div>
  </div>
</div>

<style>
  .privacy-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }
  
  .privacy-content {
    background-color: var(--card-bg);
    border-radius: 8px;
    width: 90%;
    max-width: 700px;
    max-height: 90vh;
    padding: 2rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    overflow-y: auto;
  }
  
  h1 {
    color: var(--heading-color);
    margin-top: 0;
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 1rem;
  }
  
  .policy-text {
    line-height: 1.6;
    color: var(--text-color);
  }
  
  .action-buttons {
    margin-top: 2rem;
    display: flex;
    justify-content: flex-end;
  }
  
  .close-button {
    background-color: var(--primary-button);
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    transition: background-color 0.2s;
  }
  
  .close-button:hover {
    background-color: var(--primary-button-hover);
  }
</style> 