<script>
  import { onMount } from 'svelte';
  
  export let ratings = [];
  let chartCanvas;
  let chart;
  let sectors = [
    'Basic Materials',
    'Communication Services',
    'Consumer Cyclical',
    'Consumer Defensive',
    'Energy',
    'Financial',
    'Healthcare',
    'Industrial',
    'Real Estate',
    'Technology',
    'Utilities'
  ];
  let selectedSectors = [];
  let sectorToggleFilters = {};
  
  $: {
    if (ratings && ratings.length) {
      initializeSectorToggles();
      if (chartCanvas) {
        renderChart();
      }
    }
  }
  
  onMount(() => {
    loadChartLibrary().then(() => {
      renderChart();
    });
  });
  
  function loadChartLibrary() {
    return new Promise((resolve) => {
      if (window.Chart) {
        resolve();
        return;
      }
      
      const script = document.createElement('script');
      script.src = 'https://cdn.jsdelivr.net/npm/chart.js';
      script.onload = () => resolve();
      document.head.appendChild(script);
    });
  }
  
  function initializeSectorToggles() {
    // Initialize sector toggle filters if not done
    if (Object.keys(sectorToggleFilters).length === 0) {
      sectors.forEach(sector => {
        sectorToggleFilters[sector] = false; // Initially all off
      });
      // Always show Market rating by default
      sectorToggleFilters['Market'] = true;
    }
  }
  
  function toggleSector(sector) {
    sectorToggleFilters[sector] = !sectorToggleFilters[sector];
    selectedSectors = Object.keys(sectorToggleFilters).filter(s => sectorToggleFilters[s]);
    renderChart();
  }
  
  function computeTrends() {
    // Group ratings by date
    const dateGroups = {};
    
    // Find all market ratings (symbol: 'MARKET')
    const marketRatings = ratings.filter(r => r.symbol === 'MARKET');
    
    // Find sector ratings (symbol: 'SECTOR')
    const sectorRatingsAll = ratings.filter(r => r.symbol === 'SECTOR');
    
    // Create a map of all dates
    const allDates = new Set();
    
    // Add all dates from market ratings
    marketRatings.forEach(rating => {
      if (!rating.date) return;
      const date = new Date(rating.date);
      const dateStr = date.toISOString().split('T')[0];
      allDates.add(dateStr);
    });
    
    // Add all dates from sector ratings
    sectorRatingsAll.forEach(rating => {
      if (!rating.date) return;
      const date = new Date(rating.date);
      const dateStr = date.toISOString().split('T')[0];
      allDates.add(dateStr);
    });
    
    // Sort dates chronologically
    const sortedDates = Array.from(allDates).sort();
    
    // Prepare market data
    const marketData = {
      label: "Overall Market",
      data: [],
      borderColor: "#FFD700", // Gold for market
      borderWidth: 3,
      tension: 0.1,
      fill: false
    };
    
    // Calculate market sentiment for each date
    sortedDates.forEach(date => {
      const marketRating = marketRatings.find(r => {
        const ratingDate = new Date(r.date);
        return ratingDate.toISOString().split('T')[0] === date;
      });
      
      if (marketRating) {
        const sentiment = typeof marketRating.stockSentiment === 'number' ? marketRating.stockSentiment : 0;
        marketData.data.push({
          x: date,
          y: sentiment + 3.5 // Convert from -3..3 to 0.5..6.5 for chart
        });
      } else {
        // No data for this date
        marketData.data.push({
          x: date,
          y: null
        });
      }
    });
    
    // Prepare sector data
    const sectorData = {};
    sectors.forEach(sector => {
      sectorData[sector] = {
        label: sector,
        data: [],
        borderColor: getRandomColor(),
        tension: 0.1,
        fill: false
      };
      
      // Find ratings for this sector
      const sectorRatings = sectorRatingsAll.filter(r => r.sector === sector);
      
      // Calculate sector sentiment for each date
      sortedDates.forEach(date => {
        const sectorRating = sectorRatings.find(r => {
          const ratingDate = new Date(r.date);
          return ratingDate.toISOString().split('T')[0] === date;
        });
        
        if (sectorRating) {
          const sentiment = typeof sectorRating.stockSentiment === 'number' ? sectorRating.stockSentiment : 0;
          sectorData[sector].data.push({
            x: date,
            y: sentiment + 3.5 // Convert from -3..3 to 0.5..6.5 for chart
          });
        } else {
          // No data for this date
          sectorData[sector].data.push({
            x: date,
            y: null
          });
        }
      });
    });
    
    return {
      labels: sortedDates.map(d => {
        const date = new Date(d);
        return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
      }),
      market: marketData,
      sectors: Object.values(sectorData),
      dates: sortedDates
    };
  }
  
  function getRandomColor() {
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
  }
  
  function renderChart() {
    if (!window.Chart || !chartCanvas || ratings.length === 0) return;
    
    const trendData = computeTrends();
    
    // Filter only selected sectors
    const selectedSectorData = trendData.sectors.filter(s => 
      selectedSectors.includes(s.label)
    );
    
    // Always include market data if it's selected
    const datasets = sectorToggleFilters['Market'] ? [trendData.market] : [];
    datasets.push(...selectedSectorData);
    
    // Clean up previous chart if it exists
    if (chart) {
      chart.destroy();
    }
    
    // Create the chart
    chart = new window.Chart(chartCanvas, {
      type: 'line',
      data: {
        labels: trendData.labels,
        datasets: datasets
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: {
            min: 0.5,
            max: 6.5,
            ticks: {
              callback: function(value) {
                // Convert the 0.5-6.5 scale to -3 to +3 display
                return Math.round(value - 3.5);
              }
            },
            title: {
              display: true,
              text: 'Sentiment Rating (-3 to +3)'
            }
          },
          x: {
            title: {
              display: true,
              text: 'Date'
            }
          }
        },
        plugins: {
          title: {
            display: true,
            text: 'Market & Sector Rating Trends'
          },
          tooltip: {
            callbacks: {
              label: function(context) {
                const label = context.dataset.label || '';
                const value = context.parsed.y !== null ? 
                  (Math.round((context.parsed.y - 3.5) * 10) / 10) : 'No data';
                return `${label}: ${value}`;
              }
            }
          },
          legend: {
            position: 'bottom'
          }
        }
      }
    });
  }
</script>

<div class="analytics">
  <h2>Market & Sector Analytics</h2>
  
  {#if ratings.length === 0}
    <p class="no-data">No market or sector ratings yet. Add some ratings to see analytics.</p>
  {:else}
    <div class="controls">
      <div class="sector-toggles">
        <h3>Display Options</h3>
        <p class="toggle-instructions">Toggle sectors to show their historical ratings</p>
        
        <div class="toggle-group">
          <button 
            class="toggle-btn market-btn" 
            class:active={sectorToggleFilters['Market']} 
            on:click={() => toggleSector('Market')}
          >
            Overall Market
          </button>
        </div>
        
        <div class="toggle-container">
          {#each sectors as sector}
            <button 
              class="toggle-btn" 
              class:active={sectorToggleFilters[sector]} 
              on:click={() => toggleSector(sector)}
            >
              {sector}
            </button>
          {/each}
        </div>
      </div>
    </div>
    
    <div class="chart-container">
      <canvas bind:this={chartCanvas}></canvas>
    </div>
    
    <div class="stats">
      <div class="stats-section">
        <h3>Latest Ratings</h3>
        <div class="rating-grid">
          <div class="rating-item market-item">
            <span class="rating-name">Overall Market</span>
            <span class="rating-value">
              {#if ratings.find(r => r.symbol === 'MARKET')}
                {ratings.filter(r => r.symbol === 'MARKET').sort((a, b) => new Date(b.date) - new Date(a.date))[0].stockSentiment}
              {:else}
                N/A
              {/if}
            </span>
          </div>
          
          {#each sectors as sector}
            <div class="rating-item">
              <span class="rating-name">{sector}</span>
              <span class="rating-value">
                {#if ratings.find(r => r.symbol === 'SECTOR' && r.sector === sector)}
                  {ratings.filter(r => r.symbol === 'SECTOR' && r.sector === sector).sort((a, b) => new Date(b.date) - new Date(a.date))[0].stockSentiment}
                {:else}
                  N/A
                {/if}
              </span>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .analytics {
    max-width: 800px;
    margin: 0 auto;
    padding: 1rem;
    color: var(--text-color);
  }
  
  h2, h3 {
    text-align: center;
    color: var(--text-color);
  }
  
  h2 {
    margin-bottom: 1.5rem;
  }
  
  h3 {
    margin-bottom: 0.5rem;
  }
  
  .no-data {
    text-align: center;
    font-style: italic;
    color: var(--text-color-muted);
    padding: 2rem;
  }
  
  .controls {
    margin-bottom: 1rem;
  }
  
  .toggle-instructions {
    font-size: 0.9rem;
    color: var(--text-color-muted);
    text-align: center;
    margin-bottom: 0.5rem;
  }
  
  .toggle-group {
    display: flex;
    justify-content: center;
    margin-bottom: 1rem;
  }
  
  .toggle-container {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    justify-content: center;
  }
  
  .toggle-btn {
    padding: 0.25rem 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    background: var(--input-bg);
    color: var(--text-color);
    cursor: pointer;
    font-size: 0.9rem;
    transition: all 0.2s;
  }
  
  .toggle-btn.active {
    background-color: var(--active-button);
    color: white;
    border-color: var(--active-button);
  }
  
  .market-btn {
    background-color: var(--border-color);
    color: var(--text-color);
    font-weight: bold;
    padding: 0.4rem 1rem;
    margin-bottom: 0.5rem;
  }
  
  .market-btn.active {
    background-color: #FFD700;
    color: #333;
    border-color: #FFD700;
  }
  
  .chart-container {
    width: 100%;
    height: 400px;
    margin: 2rem 0;
    background-color: var(--card-bg);
    border-radius: 8px;
    padding: 1rem;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  }
  
  .stats {
    display: flex;
    flex-wrap: wrap;
    gap: 2rem;
    justify-content: center;
  }
  
  .stats-section {
    flex: 1;
    min-width: 300px;
    background-color: var(--card-bg);
    padding: 1rem;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  }
  
  .rating-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 0.75rem;
  }
  
  .rating-item {
    background-color: var(--bg-color);
    padding: 0.5rem;
    border-radius: 4px;
    border: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  
  .market-item {
    background-color: rgba(255, 215, 0, 0.1);
    border-color: rgba(255, 215, 0, 0.4);
    grid-column: 1 / -1;
    padding: 0.75rem;
  }
  
  .rating-name {
    font-weight: bold;
    margin-bottom: 0.25rem;
  }
  
  .rating-value {
    font-size: 1.5rem;
    font-weight: bold;
  }
</style> 