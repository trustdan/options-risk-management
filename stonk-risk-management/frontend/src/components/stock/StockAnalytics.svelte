<script>
  import { onMount } from 'svelte';
  
  export let ratings = [];
  let chartCanvas;
  let chart;
  let sectors = [];
  let symbols = [];
  let selectedSymbols = [];
  let symbolToggleFilters = {};
  
  $: {
    if (ratings && ratings.length) {
      extractSectorsAndSymbols();
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
  
  function extractSectorsAndSymbols() {
    // Extract unique sectors and symbols
    sectors = [...new Set(ratings.map(r => r.sector))].filter(Boolean);
    symbols = [...new Set(ratings.map(r => r.symbol))].filter(Boolean);
    
    // Initialize symbol toggle filters if not done
    if (Object.keys(symbolToggleFilters).length === 0 && symbols.length > 0) {
      symbols.forEach(symbol => {
        symbolToggleFilters[symbol] = false; // Initially all off
      });
    }
  }
  
  function toggleSymbol(symbol) {
    symbolToggleFilters[symbol] = !symbolToggleFilters[symbol];
    selectedSymbols = Object.keys(symbolToggleFilters).filter(s => symbolToggleFilters[s]);
    renderChart();
  }
  
  function computeSectorTrends() {
    // Group ratings by date
    const dateGroups = {};
    ratings.forEach(rating => {
      if (!rating.date || !rating.sector) return;
      
      const date = new Date(rating.date);
      const dateStr = date.toISOString().split('T')[0];
      
      if (!dateGroups[dateStr]) {
        dateGroups[dateStr] = {};
      }
      
      if (!dateGroups[dateStr][rating.sector]) {
        dateGroups[dateStr][rating.sector] = [];
      }
      
      dateGroups[dateStr][rating.sector].push(rating);
    });
    
    // Sort dates
    const sortedDates = Object.keys(dateGroups).sort();
    
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
      
      // Calculate average sentiment for each date
      sortedDates.forEach(date => {
        if (dateGroups[date][sector]) {
          const sectorRatings = dateGroups[date][sector];
          const avgSentiment = sectorRatings.reduce((sum, r) => sum + r.stockSentiment, 0) / sectorRatings.length;
          sectorData[sector].data.push({
            x: date,
            y: avgSentiment
          });
        } else {
          // No data for this sector on this date
          sectorData[sector].data.push({
            x: date,
            y: null
          });
        }
      });
    });
    
    // Prepare symbol data if any symbols are selected
    const symbolData = {};
    if (selectedSymbols.length > 0) {
      selectedSymbols.forEach(symbol => {
        symbolData[symbol] = {
          label: symbol,
          data: [],
          borderColor: getRandomColor(),
          borderWidth: 2,
          borderDash: [5, 5],
          pointRadius: 4,
          tension: 0,
          fill: false
        };
        
        // Get all ratings for this symbol
        const symbolRatings = ratings.filter(r => r.symbol === symbol);
        
        // Group by date and calculate
        const symbolByDate = {};
        symbolRatings.forEach(rating => {
          const date = new Date(rating.date);
          const dateStr = date.toISOString().split('T')[0];
          symbolByDate[dateStr] = rating.stockSentiment;
        });
        
        // Fill in data points
        sortedDates.forEach(date => {
          symbolData[symbol].data.push({
            x: date,
            y: symbolByDate[date] || null
          });
        });
      });
    }
    
    return {
      labels: sortedDates.map(d => {
        const date = new Date(d);
        return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
      }),
      sectors: Object.values(sectorData),
      symbols: Object.values(symbolData),
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
    
    const trendData = computeSectorTrends();
    
    // Clean up previous chart if it exists
    if (chart) {
      chart.destroy();
    }
    
    // Create the chart
    chart = new window.Chart(chartCanvas, {
      type: 'line',
      data: {
        labels: trendData.labels,
        datasets: [
          ...trendData.sectors,
          ...trendData.symbols
        ]
      },
      options: {
        responsive: true,
        scales: {
          y: {
            beginAtZero: true,
            max: 10,
            title: {
              display: true,
              text: 'Average Sentiment'
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
            text: 'Sector & Symbol Sentiment Trends'
          },
          tooltip: {
            callbacks: {
              label: function(context) {
                const label = context.dataset.label || '';
                const value = context.parsed.y !== null ? context.parsed.y.toFixed(1) : 'No data';
                return `${label}: ${value}`;
              }
            }
          }
        }
      }
    });
  }
</script>

<div class="analytics">
  <h2>Stock Rating Analytics</h2>
  
  {#if ratings.length === 0}
    <p class="no-data">No stock ratings yet. Add some ratings to see analytics.</p>
  {:else}
    <div class="controls">
      <div class="symbol-toggles">
        <h3>Symbols</h3>
        <p class="toggle-instructions">Toggle symbols to overlay individual trendlines</p>
        <div class="toggle-container">
          {#each symbols as symbol}
            <button 
              class="toggle-btn" 
              class:active={symbolToggleFilters[symbol]} 
              on:click={() => toggleSymbol(symbol)}
            >
              {symbol}
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
        <h3>Sector Overview</h3>
        <div class="sector-grid">
          {#each sectors as sector}
            <div class="sector-item">
              <span class="sector-name">{sector}</span>
              <span class="sector-count">
                {ratings.filter(r => r.sector === sector).length} ratings
              </span>
            </div>
          {/each}
        </div>
      </div>
      
      <div class="stats-section">
        <h3>Symbol Count</h3>
        <div class="symbol-grid">
          {#each symbols as symbol}
            <div class="symbol-item">
              <span class="symbol-name">{symbol}</span>
              <span class="symbol-count">
                {ratings.filter(r => r.symbol === symbol).length} ratings
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
  }
  
  h2, h3 {
    text-align: center;
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
    color: #666;
    padding: 2rem;
  }
  
  .controls {
    margin-bottom: 1rem;
  }
  
  .toggle-instructions {
    font-size: 0.9rem;
    color: #666;
    text-align: center;
    margin-bottom: 0.5rem;
  }
  
  .toggle-container {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    justify-content: center;
  }
  
  .toggle-btn {
    padding: 0.25rem 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    background: white;
    cursor: pointer;
    font-size: 0.9rem;
    transition: all 0.2s;
  }
  
  .toggle-btn.active {
    background-color: #4caf50;
    color: white;
    border-color: #4caf50;
  }
  
  .chart-container {
    width: 100%;
    height: 400px;
    margin: 2rem 0;
  }
  
  .stats {
    display: flex;
    flex-wrap: wrap;
    gap: 2rem;
  }
  
  .stats-section {
    flex: 1;
    min-width: 250px;
    background-color: #f9f9f9;
    padding: 1rem;
    border-radius: 8px;
  }
  
  .sector-grid, .symbol-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 0.75rem;
  }
  
  .sector-item, .symbol-item {
    background-color: white;
    padding: 0.5rem;
    border-radius: 4px;
    border: 1px solid #eee;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  
  .sector-name, .symbol-name {
    font-weight: bold;
    margin-bottom: 0.25rem;
  }
  
  .sector-count, .symbol-count {
    font-size: 0.8rem;
    color: #666;
  }
</style> 