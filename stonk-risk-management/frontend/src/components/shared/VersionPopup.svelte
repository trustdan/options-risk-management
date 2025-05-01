<script>
  import { onMount } from 'svelte';
  import { VERSION } from '../../version.js';
  
  export let visible = true;
  
  let timer;
  
  onMount(() => {
    // Set visible to true to ensure it shows on startup
    visible = true;
    
    // Auto-close after 4 seconds
    timer = setTimeout(() => {
      visible = false;
    }, 4000);
    
    return () => {
      clearTimeout(timer);
    };
  });
  
  function closePopup() {
    visible = false;
    clearTimeout(timer);
  }
</script>

{#if visible}
  <div class="version-popup">
    <button class="close-button" on:click={closePopup}>×</button>
    <div class="popup-content">
      <h3>Version {VERSION}</h3>
      <p>
        <a href="https://github.com/trustdan/options-risk-management/releases" target="_blank">
          Check for updates
        </a>
      </p>
    </div>
  </div>
{/if}

<style>
  .version-popup {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: var(--card-bg, white);
    border: 1px solid var(--border-color, #e2e8f0);
    border-radius: 6px;
    padding: 15px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 9999;
    min-width: 250px;
    text-align: center;
    animation: fadeIn 0.3s ease-out;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
  
  .close-button {
    position: absolute;
    top: 5px;
    right: 8px;
    background: transparent;
    border: none;
    font-size: 18px;
    cursor: pointer;
    color: var(--text-muted, #718096);
    padding: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
  }
  
  .close-button:hover {
    background: rgba(0, 0, 0, 0.05);
    color: var(--text-color, #333);
  }
  
  .popup-content {
    margin: 0 auto;
  }
  
  h3 {
    margin: 0 0 8px 0;
    font-size: 18px;
    color: var(--heading-color, #2d3748);
  }
  
  p {
    margin: 0;
    font-size: 14px;
  }
  
  a {
    color: var(--primary-button, #38b2ac);
    text-decoration: none;
  }
  
  a:hover {
    text-decoration: underline;
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