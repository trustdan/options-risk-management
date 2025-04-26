<script>
  import { onMount } from 'svelte';
  import { GetStockRatings, SaveStockRating } from '../../../wailsjs/go/main/App';
  
  // Sample sectors
  const sectors = [
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
  
  let rating = {
    id: '',
    date: new Date().toISOString().split('T')[0],
    symbol: '',
    sector: '',
    stockSentiment: 0,
    priceTarget: 0,
    confidence: 0,
    chartPattern: '',
    enthusiasm: 0,
    notes: ''
  };
  
  let marketRatings = {
    overall: 0,
    sectorRatings: {}
  };
  
  // Initialize sector ratings with zeros
  sectors.forEach(sector => {
    marketRatings.sectorRatings[sector] = 0;
  });
  
  let recentRatings = [];
  let dateString = new Date().toISOString().split('T')[0];
  
  onMount(async () => {
    try {
      const ratings = await GetStockRatings();
      if (ratings && ratings.length > 0) {
        // Group by date to show the most recent
        const byDate = groupByDate(ratings);
        const latestDate = Object.keys(byDate).sort().pop();
        
        if (latestDate) {
          recentRatings = byDate[latestDate];
          dateString = new Date(latestDate).toISOString().split('T')[0];
        }
      }
    } catch (error) {
      console.error('Failed to load stock ratings:', error);
    }
  });
  
  function groupByDate(ratings) {
    const byDate = {};
    ratings.forEach(r => {
      const dateStr = new Date(r.date).toISOString().split('T')[0];
      if (!byDate[dateStr]) byDate[dateStr] = [];
      byDate[dateStr].push(r);
    });
    return byDate;
  }
  
  // Helper to get sector rating value safely
  function getSectorRating(sector) {
    return marketRatings.sectorRatings[sector] || 0;
  }
  
  // Helper to set sector rating value safely
  function setSectorRating(sector, value) {
    marketRatings.sectorRatings[sector] = value;
  }
  
  async function saveRating() {
    try {
      await SaveStockRating({
        ...rating,
        date: new Date(rating.date)
      });
      
      // Refresh ratings
      const ratings = await GetStockRatings();
      const byDate = groupByDate(ratings);
      const dateStr = rating.date;
      
      if (byDate[dateStr]) {
        recentRatings = byDate[dateStr];
      }
      
      // Clear form for next entry
      rating = {
        id: '',
        date: rating.date,
        symbol: '',
        sector: '',
        stockSentiment: 0,
        priceTarget: 0,
        confidence: 0,
        chartPattern: '',
        enthusiasm: 0,
        notes: ''
      };
      
      alert('Rating saved successfully');
    } catch (error) {
      console.error('Failed to save rating:', error);
      alert('Failed to save rating: ' + error.message);
    }
  }
  
  function resetForm() {
    rating = {
      id: '',
      date: rating.date,
      symbol: '',
      sector: '',
      stockSentiment: 0,
      priceTarget: 0,
      confidence: 0,
      chartPattern: '',
      enthusiasm: 0,
      notes: ''
    };
  }
</script>

<div class="dashboard">
  <h2>Stock Rating Dashboard</h2>
  <div class="divider"></div>
  
  <p class="description">Rate market conditions, sectors, and individual stocks to identify the best trading opportunities.</p>
  
  <div class="rating-container">
    <div class="rating-section market-section">
      <h3>Market & Sectors</h3>
      
      <div class="date-picker">
        <label>Rating Date:</label>
        <input type="date" bind:value={rating.date} />
      </div>
      
      <div class="slider-group">
        <h4>Overall Market: {marketRatings.overall}</h4>
        <input 
          type="range" 
          min="-3" 
          max="3" 
          step="1" 
          bind:value={marketRatings.overall}
        />
        <div class="scale-labels">
          <span>Bearish (-3)</span>
          <span>Bullish (+3)</span>
        </div>
      </div>
      
      <div class="sector-ratings">
        <h4>Sector Ratings</h4>
        <div class="primary-sector">
          <label>Primary Sector for This Stock:</label>
          <select bind:value={rating.sector}>
            <option value="">Select primary sector...</option>
            {#each sectors as sector}
              <option value={sector}>{sector}</option>
            {/each}
          </select>
        </div>
        
        <div class="scrollable-sectors">
          {#each sectors as sector}
            <div class="sector-slider">
              <h5>{sector}: {getSectorRating(sector)}</h5>
              <input 
                type="range" 
                min="-3" 
                max="3" 
                step="1" 
                value={getSectorRating(sector)}
                on:input={(e) => setSectorRating(sector, parseInt(e.target.value))}
              />
              <div class="scale-labels small">
                <span>Bearish (-3)</span>
                <span>Bullish (+3)</span>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>
    
    <div class="rating-section stock-section">
      <h3>Stock Rating</h3>
      
      <div class="input-group">
        <label>Ticker Symbol:</label>
        <input 
          type="text" 
          placeholder="e.g., AAPL" 
          bind:value={rating.symbol}
        />
      </div>
      
      <div class="slider-group">
        <h4>Stock Sentiment: {rating.stockSentiment}</h4>
        <input 
          type="range" 
          min="-3" 
          max="3" 
          step="1" 
          bind:value={rating.stockSentiment}
        />
        <div class="scale-labels">
          <span>Bearish (-3)</span>
          <span>Bullish (+3)</span>
        </div>
      </div>
      
      <div class="dual-inputs">
        <div class="input-group">
          <label>Chart Pattern:</label>
          <select bind:value={rating.chartPattern}>
            <option value="">Select a pattern...</option>
            <option value="Bull Flag">Bull Flag</option>
            <option value="Cup & Handle">Cup & Handle</option>
            <option value="Double Bottom">Double Bottom</option>
            <option value="Ascending Triangle">Ascending Triangle</option>
            <option value="Bear Flag">Bear Flag</option>
            <option value="Head & Shoulders">Head & Shoulders</option>
            <option value="Double Top">Double Top</option>
            <option value="Descending Triangle">Descending Triangle</option>
          </select>
          <button class="info-btn">Pattern Info</button>
        </div>
      </div>
      
      <div class="slider-group">
        <h4>Estimated Enthusiasm: {rating.enthusiasm}</h4>
        <input 
          type="range" 
          min="0" 
          max="10" 
          step="1" 
          bind:value={rating.enthusiasm}
        />
      </div>
      
      <div class="buttons-row">
        <button class="reset-btn" on:click={resetForm}>Reset</button>
        <button class="save-btn" on:click={saveRating}>Save Rating</button>
      </div>
    </div>
  </div>
  
  <div class="recent-ratings">
    <h3>Recent Stock Ratings</h3>
    <div class="ratings-table">
      <table>
        <thead>
          <tr>
            <th>Symbol</th>
            <th>Sector</th>
            <th>Sentiment</th>
            <th>Pattern</th>
            <th>Enthusiasm</th>
          </tr>
        </thead>
        <tbody>
          {#if recentRatings.length > 0}
            {#each recentRatings as item}
              <tr>
                <td>{item.symbol}</td>
                <td>{item.sector}</td>
                <td class={item.stockSentiment > 0 ? 'positive' : item.stockSentiment < 0 ? 'negative' : 'neutral'}>
                  {item.stockSentiment > 0 ? '+' : ''}{item.stockSentiment}
                </td>
                <td>{item.chartPattern || '-'}</td>
                <td>{item.enthusiasm}/10</td>
              </tr>
            {/each}
          {:else}
            <tr>
              <td colspan="5">No ratings available for this date</td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>

<style>
  .dashboard {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
    background-color: white;
    border-radius: 5px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  }
  
  h2, h3, h4, h5 {
    color: #2d3748;
    margin-top: 0;
    margin-bottom: 0.5rem;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 0;
  }
  
  .divider {
    height: 2px;
    background-color: #e2e8f0;
    margin: 1rem 0 2rem;
  }
  
  .description {
    text-align: center;
    color: #4a5568;
    margin-bottom: 2rem;
  }
  
  .rating-container {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    margin-bottom: 2rem;
  }
  
  .rating-section {
    background-color: white;
    padding: 1.5rem;
    border-radius: 5px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  }
  
  .date-picker {
    margin-bottom: 1.5rem;
    display: flex;
    align-items: center;
  }
  
  .date-picker label {
    margin-right: 1rem;
    font-weight: bold;
  }
  
  .date-picker input {
    padding: 0.5rem;
    border: 1px solid #cbd5e0;
    border-radius: 4px;
  }
  
  .input-group {
    margin-bottom: 1rem;
  }
  
  .input-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: bold;
  }
  
  .input-group input, .input-group select {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #cbd5e0;
    border-radius: 4px;
  }
  
  .dual-inputs {
    display: grid;
    grid-template-columns: 1fr;
    gap: 1rem;
    margin-bottom: 1rem;
  }
  
  .slider-group {
    margin-bottom: 1.5rem;
  }
  
  .slider-group input[type="range"] {
    width: 100%;
    margin: 0.5rem 0;
  }
  
  .scale-labels {
    display: flex;
    justify-content: space-between;
    color: #718096;
    font-size: 0.85rem;
  }
  
  .scale-labels.small {
    font-size: 0.75rem;
  }
  
  .sector-ratings {
    margin-top: 1.5rem;
  }
  
  .primary-sector {
    margin-bottom: 1rem;
  }
  
  .scrollable-sectors {
    max-height: 300px;
    overflow-y: auto;
    border: 1px solid #e2e8f0;
    padding: 1rem;
    border-radius: 4px;
  }
  
  .sector-slider {
    margin-bottom: 1rem;
  }
  
  .buttons-row {
    display: flex;
    justify-content: space-between;
    margin-top: 2rem;
  }
  
  .reset-btn, .save-btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-weight: bold;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .reset-btn {
    background-color: #e2e8f0;
    color: #4a5568;
  }
  
  .save-btn {
    background-color: #4299e1;
    color: white;
  }
  
  .reset-btn:hover {
    background-color: #cbd5e0;
  }
  
  .save-btn:hover {
    background-color: #3182ce;
  }
  
  .info-btn {
    margin-top: 0.5rem;
    padding: 0.25rem 0.5rem;
    background-color: #e2e8f0;
    border: none;
    border-radius: 4px;
    color: #4a5568;
    font-size: 0.8rem;
    cursor: pointer;
  }
  
  .recent-ratings {
    background-color: white;
    padding: 1.5rem;
    border-radius: 5px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  }
  
  .ratings-table {
    overflow-x: auto;
  }
  
  table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 1rem;
  }
  
  thead {
    background-color: #f7fafc;
  }
  
  th, td {
    padding: 0.75rem 1rem;
    text-align: left;
    border-bottom: 1px solid #e2e8f0;
  }
  
  .positive {
    color: #38a169;
  }
  
  .negative {
    color: #e53e3e;
  }
  
  .neutral {
    color: #718096;
  }
</style> 