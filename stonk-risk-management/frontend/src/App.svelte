<script>
  import RiskDashboard from './components/risk/RiskDashboard.svelte';
  import StockDashboard from './components/stock/StockDashboard.svelte';
  import TradeCalendar from './components/trade/TradeCalendar.svelte';
  import { WindowSetDarkTheme, WindowSetLightTheme } from '../wailsjs/runtime/runtime';
  
  let activeView = 'stock'; // Default to 'stock' view to match screenshots
  let isDarkMode = true; // Default to dark mode to match screenshots
  const version = "v0.1.0";
  
  function switchView(view) {
    activeView = view;
  }
  
  function toggleTheme() {
    isDarkMode = !isDarkMode;
    
    if (isDarkMode) {
      WindowSetDarkTheme();
      document.body.classList.add('dark-mode');
    } else {
      WindowSetLightTheme();
      document.body.classList.remove('dark-mode');
    }
  }
  
  // On mount, set the initial theme
  import { onMount } from 'svelte';
  
  onMount(() => {
    if (isDarkMode) {
      WindowSetDarkTheme();
      document.body.classList.add('dark-mode');
    } else {
      WindowSetLightTheme();
    }
  });
</script>

<main class:dark-mode={isDarkMode}>
  <header>
    <div class="header-content">
      <h1>Trading Dashboard</h1>
      <div class="header-right">
        <span class="version">{version}</span>
        <button class="theme-toggle" on:click={toggleTheme}>
          {isDarkMode ? '☀️' : '🌙'}
        </button>
      </div>
    </div>
    <nav>
      <button 
        class:active={activeView === 'risk'} 
        on:click={() => switchView('risk')}
      >
        Risk Management
      </button>
      <button 
        class:active={activeView === 'stock'} 
        on:click={() => switchView('stock')}
      >
        Stock Rating
      </button>
      <button 
        class:active={activeView === 'trade'} 
        on:click={() => switchView('trade')}
      >
        Trade Calendar
      </button>
    </nav>
  </header>
  
  <div class="content">
    {#if activeView === 'risk'}
      <RiskDashboard />
    {:else if activeView === 'stock'}
      <StockDashboard />
    {:else if activeView === 'trade'}
      <TradeCalendar />
    {/if}
  </div>
  
  <footer>
    <p>© {new Date().getFullYear()} Stonk Risk Management (SRM) - v1.0.0</p>
  </footer>
</main>

<style>
  :global(:root) {
    /* Light theme variables */
    --bg-color: #f7f9fc;
    --text-color: #333;
    --header-bg: #2d3748;
    --header-text: white;
    --footer-bg: #e2e8f0;
    --footer-text: #4a5568;
    --active-button: #38b2ac;
    --button-hover: rgba(255, 255, 255, 0.1);
    --content-bg: white;
    --card-bg: white;
    --border-color: #e2e8f0;
    --input-bg: white;
    --input-border: #cbd5e0;
    --primary-button: #38b2ac;
    --secondary-button: #e2e8f0;
  }
  
  :global(body.dark-mode) {
    /* Dark theme variables */
    --bg-color: #1a202c;
    --text-color: #e2e8f0;
    --header-bg: #1e2533;
    --header-text: white;
    --footer-bg: #1e2533;
    --footer-text: #e2e8f0;
    --active-button: #f05252;
    --button-hover: rgba(255, 255, 255, 0.2);
    --content-bg: #1e2533;
    --card-bg: #2d3748;
    --border-color: #4a5568;
    --input-bg: #2d3748;
    --input-border: #4a5568;
    --primary-button: #f05252;
    --secondary-button: #4a5568;
  }

  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    margin: 0;
    padding: 0;
    color: var(--text-color);
    background-color: var(--bg-color);
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  main {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }
  
  header {
    background-color: var(--header-bg);
    color: var(--header-text);
    padding: 1rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  .header-right {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  .version {
    font-size: 0.8rem;
    opacity: 0.8;
    color: var(--header-text);
  }
  
  h1 {
    margin: 0;
    font-size: 1.5rem;
  }
  
  .theme-toggle {
    background: transparent;
    border: none;
    color: var(--header-text);
    font-size: 1.25rem;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 50%;
    transition: background-color 0.2s;
  }
  
  .theme-toggle:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
  
  nav {
    display: flex;
    gap: 0;
  }
  
  nav button {
    background-color: transparent;
    border: none;
    color: var(--header-text);
    padding: 0.75rem 1.5rem;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s;
    border-radius: 0;
  }
  
  nav button:hover {
    background-color: var(--button-hover);
  }
  
  nav button.active {
    background-color: var(--active-button);
    color: white;
  }
  
  /* Special styling for active tabs */
  body:not(.dark-mode) nav button.active {
    background-color: #38b2ac;
  }
  
  body.dark-mode nav button.active {
    background-color: #f05252;
  }
  
  .content {
    flex: 1;
    padding: 0;
    width: 100%;
    max-width: 100%;
  }
  
  footer {
    text-align: center;
    padding: 0.5rem;
    background-color: var(--footer-bg);
    font-size: 0.8rem;
    color: var(--footer-text);
    transition: background-color 0.3s ease, color 0.3s ease;
  }
</style>
