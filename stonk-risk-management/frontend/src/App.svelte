<script>
  import RiskDashboard from './components/risk/RiskDashboard.svelte';
  import StockDashboard from './components/stock/StockDashboard.svelte';
  import TradeCalendar from './components/trade/TradeCalendar.svelte';
  import TradingKoans from './components/koans/TradingKoans.svelte';
  import PrivacyPolicy from './components/PrivacyPolicy.svelte';
  import { WindowSetDarkTheme, WindowSetLightTheme } from '../wailsjs/runtime/runtime';
  import { VERSION } from './version.js';
  import { RunDatabaseMaintenance } from '../wailsjs/go/main/App';
  
  let activeView = 'stock'; // Default to 'stock' view to match screenshots
  let isDarkMode = true; // Default to dark mode to match screenshots
  const version = "v" + VERSION;
  let showPrivacyPolicy = false;
  
  // Reference to component instances
  let riskDashboardComponent;
  let stockDashboardComponent;
  let tradeCalendarComponent;
  
  // Determine data directory path based on operating system
  function getDataPath() {
    const isWindows = navigator.platform.indexOf('Win') !== -1;
    if (isWindows) {
      return `C:\\Users\\<YourUsername>\\.options-risk-management 
      
(Replace <YourUsername> with your Windows username)`;
    } else {
      return `/home/<YourUsername>/.options-risk-management
      
(Replace <YourUsername> with your username)`;
    }
  }
  
  const dataPath = getDataPath();
  
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
  
  // Export data function
  function exportData() {
    const message = `EXPORT DATA INSTRUCTIONS:

Your data is stored in:
${dataPath}

To export your data:
1. Navigate to this directory in your file explorer
2. Copy the entire folder or its contents to your backup location
3. All your risk assessments, stock ratings, and trades are stored in this directory

Note: The app must be closed before copying to ensure all data is properly saved.`;
    alert(message);
  }
  
  // Import data function
  function importData() {
    const message = `IMPORT DATA INSTRUCTIONS:

Your data should be placed in:
${dataPath}

To import your data:
1. Ensure the app is closed
2. Navigate to this directory in your file explorer
3. Copy your backed-up files into this directory, replacing any existing files
4. Start the application again to see your imported data

Warning: This will overwrite any existing data in the application. Make a backup first if needed.`;
    alert(message);
  }
  
  // Database maintenance function
  async function runDatabaseMaintenance() {
    try {
      const result = await RunDatabaseMaintenance();
      alert(result);
    } catch (error) {
      alert(`Error: ${error}`);
    }
  }
  
  // Function to refresh the entire application
  function refreshApp() {
    console.log('Global refresh triggered');
    
    // Check which component is active and refresh it
    if (activeView === 'risk' && riskDashboardComponent && riskDashboardComponent.loadRiskAssessments) {
      riskDashboardComponent.loadRiskAssessments();
    } else if (activeView === 'stock' && stockDashboardComponent && stockDashboardComponent.loadStockRatings) {
      stockDashboardComponent.loadStockRatings();
    } else if (activeView === 'trade' && tradeCalendarComponent && tradeCalendarComponent.forceRefresh) {
      tradeCalendarComponent.forceRefresh();
    }
    
    // Reload the window as a fallback if component methods don't exist
    setTimeout(() => {
      // Force a window refresh in case component methods didn't work
      window.location.reload();
    }, 500);
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
    
    // Load the Buy Me A Coffee script
    const script = document.createElement('script');
    script.src = "https://cdnjs.buymeacoffee.com/1.0.0/button.prod.min.js";
    script.setAttribute('data-name', 'bmc-button');
    script.setAttribute('data-slug', 'admitchu');
    script.setAttribute('data-color', '#40DCA5');
    script.setAttribute('data-emoji', '☕');
    script.setAttribute('data-font', 'Cookie');
    script.setAttribute('data-text', 'Buy me a coffee');
    script.setAttribute('data-outline-color', '#000000');
    script.setAttribute('data-font-color', '#ffffff');
    script.setAttribute('data-coffee-color', '#FFDD00');
    document.body.appendChild(script);
    
    // Load Cookie font
    const fontLink = document.createElement('link');
    fontLink.rel = 'stylesheet';
    fontLink.href = 'https://fonts.googleapis.com/css2?family=Cookie&display=swap';
    document.head.appendChild(fontLink);
  });
</script>

<main class:dark-mode={isDarkMode}>
  <a href="https://github.com/trustdan/options-risk-management" target="_blank" class="github-corner" aria-label="View source on GitHub">
    <svg width="80" height="80" viewBox="0 0 250 250" style="fill:{isDarkMode ? '#f05252' : '#38b2ac'}; color:#fff; position: absolute; top: 0; border: 0; right: 0;" aria-hidden="true">
      <path d="M0,0 L115,115 L130,115 L142,142 L250,250 L250,0 Z"></path>
      <path d="M128.3,109.0 C113.8,99.7 119.0,89.6 119.0,89.6 C122.0,82.7 120.5,78.6 120.5,78.6 C119.2,72.0 123.4,76.3 123.4,76.3 C127.3,80.9 125.5,87.3 125.5,87.3 C122.9,97.6 130.6,101.9 134.4,103.2" fill="currentColor" style="transform-origin: 130px 106px;" class="octo-arm"></path>
      <path d="M115.0,115.0 C114.9,115.1 118.7,116.5 119.8,115.4 L133.7,101.6 C136.9,99.2 139.9,98.4 142.2,98.6 C133.8,88.0 127.5,74.4 143.8,58.0 C148.5,53.4 154.0,51.2 159.7,51.0 C160.3,49.4 163.2,43.6 171.4,40.1 C171.4,40.1 176.1,42.5 178.8,56.2 C183.1,58.6 187.2,61.8 190.9,65.4 C194.5,69.0 197.7,73.2 200.1,77.6 C213.8,80.2 216.3,84.9 216.3,84.9 C212.7,93.1 206.9,96.0 205.4,96.6 C205.1,102.4 203.0,107.8 198.3,112.5 C181.9,128.9 168.3,122.5 157.7,114.1 C157.9,116.9 156.7,120.9 152.7,124.9 L141.0,136.5 C139.8,137.7 141.6,141.9 141.8,141.8 Z" fill="currentColor" class="octo-body"></path>
    </svg>
  </a>
  
  <header>
    <div class="header-content">
      <h1>Options Trading Dashboard</h1>
      <div class="header-right">
        <button class="data-btn" on:click={exportData} title="Export your data">
          💾 Export Data
        </button>
        <button class="data-btn" on:click={importData} title="Import your data">
          📥 Import Data
        </button>
        <button class="data-btn" on:click={runDatabaseMaintenance} title="Clean up database and free space">
          🧹 Database Cleanup
        </button>
        <a href="https://docs.google.com/forms/d/e/1FAIpQLSdy5CEPbAFFli589aR1DPEVUJli5MxRNYsl8PmS4E-srgt7IA/viewform" target="_blank" class="feedback-btn" title="Provide Feedback">
          📝 Feedback
        </a>
        <button class="refresh-btn" on:click={refreshApp} title="Refresh application data">
          ↻ Refresh
        </button>
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
        Market, Sector, Stock Rating
      </button>
      <button 
        class:active={activeView === 'trade'} 
        on:click={() => switchView('trade')}
      >
        Trades + Options Calendar
      </button>
      <button 
        class:active={activeView === 'koans'} 
        on:click={() => switchView('koans')}
      >
        Trading Koans
      </button>
      <div class="bmc-container">
        <a class="bmc-button" target="_blank" href="https://www.buymeacoffee.com/admitchu">
          <span class="coffee-emoji">☕</span>
          Buy me a coffee
        </a>
      </div>
    </nav>
  </header>
  
  <div class="content">
    {#if activeView === 'risk'}
      <RiskDashboard bind:this={riskDashboardComponent} />
    {:else if activeView === 'stock'}
      <StockDashboard bind:this={stockDashboardComponent} />
    {:else if activeView === 'trade'}
      <TradeCalendar bind:this={tradeCalendarComponent} />
    {:else if activeView === 'koans'}
      <TradingKoans />
    {/if}
  </div>
  
  <footer>
    <p>
      © {new Date().getFullYear()} Options Risk Management - v{VERSION} | 
      <a href="#privacy" on:click|preventDefault={() => showPrivacyPolicy = true}>Privacy Policy</a>
    </p>
  </footer>
  
  {#if showPrivacyPolicy}
    <PrivacyPolicy on:close={() => showPrivacyPolicy = false} />
  {/if}
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
    --card-bg-light: #f9f9f9;
    --border-color: #e2e8f0;
    --input-bg: white;
    --input-border: #cbd5e0;
    --primary-button: #38b2ac;
    --primary-button-hover: #319795;
    --secondary-button: #4a5568;
    --secondary-button-hover: #3c4655;
    --text-secondary: #666;
    --success-color: #10b981;
    --danger-color: #ef4444;
    --warning-color: #f59e0b;
    --table-header-color: #4a5568;
    --table-row-alternate: #f9fafb;
    --table-row-hover: #edf2f7;
    --text-muted: #718096;
    --heading-color: #2d3748;
    --label-color: #4a5568;
    --warning-bg: rgba(255, 0, 0, 0.1);
    --success-bg: rgba(0, 200, 0, 0.1);
    --neutral-bg: rgba(200, 200, 200, 0.1);
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
    --card-bg-light: #2a3141;
    --border-color: #4a5568;
    --input-bg: #2d3748;
    --input-border: #4a5568;
    --primary-button: #f05252;
    --primary-button-hover: #e02424;
    --secondary-button: #4a5568;
    --secondary-button-hover: #3c4655;
    --text-secondary: #a0aec0;
    --success-color: #0d9488;
    --danger-color: #dc2626;
    --warning-color: #d97706;
    --table-header-color: #a0aec0;
    --table-row-alternate: #283141;
    --table-row-hover: #2d3748;
    --text-muted: #718096;
    --heading-color: #e2e8f0;
    --label-color: #cbd5e0;
    --warning-bg: rgba(239, 68, 68, 0.2);
    --success-bg: rgba(16, 185, 129, 0.2);
    --neutral-bg: rgba(160, 174, 192, 0.2);
  }

  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    margin: 0;
    padding: 0;
    color: var(--text-color);
    background-color: var(--bg-color);
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  /* GitHub corner styles */
  .github-corner:hover .octo-arm {
    animation: octocat-wave 560ms ease-in-out;
    transform-origin: 130px 106px; 
  }
  
  @keyframes octocat-wave {
    0%, 100% { transform: rotate(0); }
    20%, 60% { transform: rotate(-25deg); }
    40%, 80% { transform: rotate(10deg); }
  }
  
  @media (max-width: 500px) {
    .github-corner:hover .octo-arm {
      animation: none;
    }
    .github-corner .octo-arm {
      animation: octocat-wave 560ms ease-in-out;
    }
  }
  
  main {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    position: relative;
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
    margin-right: 100px; /* Add space to avoid overlap with GitHub corner */
  }
  
  .data-btn {
    background-color: transparent;
    color: var(--header-text);
    border: none;
    padding: 0.5rem 0.75rem;
    border-radius: 0;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 0.9rem;
    font-weight: bold;
    text-decoration: none;
    transition: all 0.2s ease;
    margin-right: 0.5rem;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
  }
  
  .data-btn:hover {
    background-color: var(--button-hover);
  }
  
  .feedback-btn {
    background-color: transparent;
    color: var(--header-text);
    border: none;
    padding: 0.5rem 0.75rem;
    border-radius: 0;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 0.9rem;
    font-weight: bold;
    text-decoration: none;
    transition: all 0.2s ease;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
  }
  
  .feedback-btn:hover {
    background-color: var(--button-hover);
  }
  
  .refresh-btn {
    background-color: transparent;
    color: var(--header-text);
    border: none;
    padding: 0.5rem;
    border-radius: 0;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.9rem;
    font-weight: bold;
    text-decoration: none;
    transition: all 0.2s ease;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
  }
  
  .refresh-btn:hover {
    background-color: var(--button-hover);
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
    align-items: center;
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
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
  }
  
  nav button:hover {
    background-color: var(--button-hover);
    box-shadow: none;
  }
  
  nav button.active {
    background-color: var(--active-button);
    color: white;
    text-shadow: none;
  }
  
  body:not(.dark-mode) nav button.active {
    background-color: #38b2ac;
  }
  
  body.dark-mode nav button.active {
    background-color: #f05252;
  }
  
  .bmc-container {
    margin-left: 1rem;
    display: flex;
    align-items: center;
  }
  
  .bmc-button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background-color: #40DCA5;
    color: white;
    padding: 0.5rem 1rem;
    border-radius: 24px;
    font-family: 'Cookie', cursive;
    font-size: 1.5rem;
    text-decoration: none;
    transition: all 0.2s;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
    border: 1px solid rgba(0, 0, 0, 0.1);
    line-height: 1;
  }
  
  .bmc-button img {
    width: 20px;
    height: 20px;
    margin-right: 8px;
  }
  
  .coffee-emoji {
    margin-right: 8px;
    font-size: 1.2rem;
  }
  
  .bmc-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
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
  
  footer a {
    color: var(--primary-button);
    text-decoration: none;
  }
  
  footer a:hover {
    text-decoration: underline;
  }
  
  @media (max-width: 900px) {
    nav {
      flex-wrap: wrap;
    }
    
    .bmc-container {
      margin: 0.5rem auto;
    }
  }
</style>
