<script>
  import RiskDashboard from './components/risk/RiskDashboard.svelte';
  import StockDashboard from './components/stock/StockDashboard.svelte';
  import TradeCalendar from './components/trade/TradeCalendar.svelte';
  
  let activeView = 'risk'; // 'risk', 'stock', or 'trade'
  
  function switchView(view) {
    activeView = view;
  }
</script>

<main>
  <header>
    <h1>Trading Dashboard</h1>
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
  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    margin: 0;
    padding: 0;
    color: #333;
    background-color: #f7f9fc;
  }
  
  main {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
  }
  
  header {
    background-color: #2d3748;
    color: white;
    padding: 1rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  h1 {
    margin: 0;
    font-size: 1.5rem;
    margin-bottom: 1rem;
  }
  
  nav {
    display: flex;
    gap: 0;
  }
  
  nav button {
    background-color: transparent;
    border: none;
    color: white;
    padding: 0.75rem 1.5rem;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s;
    border-radius: 0;
  }
  
  nav button:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
  
  nav button.active {
    background-color: #38b2ac;
    color: white;
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
    background-color: #f5f5f5;
    font-size: 0.8rem;
    color: #666;
  }
</style>
