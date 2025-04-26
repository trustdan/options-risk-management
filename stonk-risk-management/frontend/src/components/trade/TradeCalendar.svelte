<script>
  import { onMount } from 'svelte';
  
  // Sample data for calendar
  const sectors = [
    'Technology',
    'Healthcare',
    'Financial',
    'Consumer Defensive',
    'Consumer Cyclical',
    'Industrial',
    'Energy',
    'Materials',
    'Utilities',
    'Real Estate',
    'Communication',
    'Others'
  ];
  
  // Strategy types with colors
  const strategies = [
    { name: 'Basic Spreads', color: '#3498db' },
    { name: 'Vertical Spreads', color: '#9b59b6' },
    { name: 'Calendar/Horizontal Spreads', color: '#2ecc71' },
    { name: 'Diagonal Spreads', color: '#f39c12' },
    { name: 'Butterfly Spreads', color: '#e74c3c' },
    { name: 'Iron Condors/Butterflies', color: '#1abc9c' },
    { name: 'Ratio Spreads', color: '#d35400' }
  ];
  
  // Generate 8 weeks for the calendar
  const weeks = [];
  const today = new Date();
  let currentDate = new Date(today);
  
  // Move to the beginning of the week (adjust as needed)
  const dayOfWeek = currentDate.getDay();
  currentDate.setDate(currentDate.getDate() - dayOfWeek);
  
  for (let i = 0; i < 8; i++) {
    const weekStart = new Date(currentDate);
    const weekNumber = i + 1;
    const expDate = new Date(weekStart);
    expDate.setDate(expDate.getDate() + 5); // Friday expiration
    
    weeks.push({
      number: weekNumber,
      startDate: weekStart,
      expirationDate: expDate,
      label: `Week ${weekNumber}`,
      expLabel: `Exp: ${expDate.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: '2-digit' })}`
    });
    
    // Move to next week
    currentDate.setDate(currentDate.getDate() + 7);
  }
  
  // Sample trade data
  let trades = [
    {
      symbol: 'AAPL',
      sector: 'Technology',
      strategy: 'Basic Spreads',
      type: 'Long Call',
      week: 1,
      entryDate: new Date(),
      expirationDate: weeks[0].expirationDate,
      entryPrice: 3.45
    },
    {
      symbol: 'AAPL',
      sector: 'Technology',
      strategy: 'Basic Spreads',
      type: 'Long Call',
      week: 2,
      entryDate: new Date(),
      expirationDate: weeks[1].expirationDate,
      entryPrice: 2.80
    },
    {
      symbol: 'GE',
      sector: 'Industrial',
      strategy: 'Diagonal Spreads',
      type: 'Diagonal Call Spread Up',
      week: 4,
      entryDate: new Date(),
      expirationDate: weeks[3].expirationDate,
      entryPrice: 1.25
    },
    {
      symbol: 'GE',
      sector: 'Industrial',
      strategy: 'Diagonal Spreads',
      type: 'Diagonal Call Spread Up',
      week: 5,
      entryDate: new Date(),
      expirationDate: weeks[4].expirationDate,
      entryPrice: 1.35
    },
    {
      symbol: 'GE',
      sector: 'Industrial',
      strategy: 'Diagonal Spreads',
      type: 'Diagonal Call Spread Up',
      week: 6,
      entryDate: new Date(),
      expirationDate: weeks[5].expirationDate,
      entryPrice: 1.40
    },
    {
      symbol: 'SPY',
      sector: 'Others',
      strategy: 'Diagonal Spreads',
      type: 'Diagonal Put Spread Down',
      week: 3,
      entryDate: new Date(),
      expirationDate: weeks[2].expirationDate,
      entryPrice: 2.15
    },
    {
      symbol: 'SPY',
      sector: 'Others',
      strategy: 'Diagonal Spreads',
      type: 'Diagonal Put Spread Down',
      week: 4,
      entryDate: new Date(),
      expirationDate: weeks[3].expirationDate,
      entryPrice: 2.25
    },
    {
      symbol: 'SPY',
      sector: 'Others',
      strategy: 'Diagonal Spreads',
      type: 'Diagonal Put Spread Down',
      week: 5,
      entryDate: new Date(),
      expirationDate: weeks[4].expirationDate,
      entryPrice: 2.20
    }
  ];
  
  // Data for the new trade form
  let newTrade = {
    entryDate: new Date().toISOString().split('T')[0],
    expirationDate: '',
    ticker: '',
    sector: '',
    entryPrice: 0,
    strategy: '',
    notes: ''
  };
  
  function getStrategyColor(strategyName) {
    const strategy = strategies.find(s => s.name === strategyName);
    return strategy ? strategy.color : '#999';
  }
  
  function getTradesForCell(sector, week) {
    return trades.filter(t => t.sector === sector && t.week === week);
  }
  
  function addNewTrade() {
    // Find the week for the expiration date
    const expDate = new Date(newTrade.expirationDate);
    const weekIndex = weeks.findIndex(w => 
      w.expirationDate.toDateString() === expDate.toDateString()
    );
    
    if (weekIndex === -1) {
      alert('Please select a valid expiration date matching one of the weeks shown');
      return;
    }
    
    const trade = {
      symbol: newTrade.ticker,
      sector: newTrade.sector,
      strategy: newTrade.strategy.split(' - ')[0],
      type: newTrade.strategy.includes(' - ') ? newTrade.strategy.split(' - ')[1] : newTrade.strategy,
      week: weekIndex + 1,
      entryDate: new Date(newTrade.entryDate),
      expirationDate: expDate,
      entryPrice: parseFloat(newTrade.entryPrice),
      notes: newTrade.notes
    };
    
    trades = [...trades, trade];
    
    // Reset form
    resetForm();
    
    alert('Trade added successfully!');
  }
  
  function resetForm() {
    newTrade = {
      entryDate: new Date().toISOString().split('T')[0],
      expirationDate: '',
      ticker: '',
      sector: '',
      entryPrice: 0,
      strategy: '',
      notes: ''
    };
  }
  
  onMount(() => {
    // Could load trades from backend here
  });
</script>

<div class="calendar-page">
  <h2>Options Trading Calendar</h2>
  <div class="divider"></div>
  
  <p class="description">Track and visualize your options trades across sectors and expiration weeks.</p>
  
  <div class="strategy-guide">
    <h3>Strategy Guide</h3>
    <div class="strategy-legend">
      {#each strategies as strategy}
        <div class="strategy-item">
          <span class="color-box" style="background-color: {strategy.color}"></span>
          <span class="strategy-name">{strategy.name}</span>
        </div>
      {/each}
    </div>
  </div>
  
  <div class="calendar-container">
    <table class="calendar">
      <thead>
        <tr>
          <th class="sector-header">Sector</th>
          {#each weeks as week}
            <th class="week-header">
              <div class="week-label">{week.label}</div>
              <div class="exp-date">{week.expLabel}</div>
            </th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each sectors as sector}
          <tr>
            <td class="sector-cell">{sector}</td>
            {#each weeks as week, weekIndex}
              <td class="trade-cell">
                {#each getTradesForCell(sector, week.number) as trade}
                  <div class="trade-card" style="background-color: {getStrategyColor(trade.strategy)}">
                    <div class="trade-symbol">{trade.symbol}</div>
                    <div class="trade-type">{trade.type}</div>
                  </div>
                {/each}
              </td>
            {/each}
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
  
  <div class="add-trade-section">
    <h3>Add New Options Trade</h3>
    
    <div class="trade-form">
      <div class="form-row">
        <div class="form-group">
          <label>Entry Date:</label>
          <input type="date" bind:value={newTrade.entryDate} />
        </div>
        
        <div class="form-group">
          <label>Expiration Date:</label>
          <input type="date" bind:value={newTrade.expirationDate} />
        </div>
        
        <div class="form-group">
          <label>Ticker:</label>
          <input type="text" placeholder="e.g., AAPL" bind:value={newTrade.ticker} />
        </div>
        
        <div class="form-group">
          <label>Sector:</label>
          <select bind:value={newTrade.sector}>
            <option value="">Select a sector...</option>
            {#each sectors as sector}
              <option value={sector}>{sector}</option>
            {/each}
          </select>
        </div>
        
        <div class="form-group">
          <label>Entry Price:</label>
          <input type="number" step="0.01" min="0" bind:value={newTrade.entryPrice} />
        </div>
      </div>
      
      <div class="form-row">
        <div class="form-group full-width">
          <label>Options Strategy:</label>
          <select bind:value={newTrade.strategy}>
            <option value="">Select a strategy...</option>
            <option value="Basic Spreads - Long Call">Basic Spreads - Long Call</option>
            <option value="Basic Spreads - Long Put">Basic Spreads - Long Put</option>
            <option value="Vertical Spreads - Bull Call Spread">Vertical Spreads - Bull Call Spread</option>
            <option value="Vertical Spreads - Bear Put Spread">Vertical Spreads - Bear Put Spread</option>
            <option value="Calendar/Horizontal Spreads - Call Calendar">Calendar Spreads - Call Calendar</option>
            <option value="Calendar/Horizontal Spreads - Put Calendar">Calendar Spreads - Put Calendar</option>
            <option value="Diagonal Spreads - Diagonal Call Spread Up">Diagonal Spreads - Diagonal Call Spread Up</option>
            <option value="Diagonal Spreads - Diagonal Put Spread Down">Diagonal Spreads - Diagonal Put Spread Down</option>
          </select>
        </div>
      </div>
      
      <div class="form-row">
        <div class="form-group full-width">
          <label>Trade Notes:</label>
          <textarea 
            placeholder="Add any relevant trade notes, strategy details, or observations..." 
            bind:value={newTrade.notes}
          ></textarea>
        </div>
      </div>
      
      <div class="form-buttons">
        <button class="reset-btn" on:click={resetForm}>Reset</button>
        <button class="add-btn" on:click={addNewTrade}>Save Trade</button>
      </div>
    </div>
  </div>
  
  <div class="trade-history">
    <h3>Options Trade History</h3>
    
    <div class="search-section">
      <input type="text" placeholder="Search trades..." class="search-input" />
      
      <div class="search-filters">
        <label>
          <input type="radio" name="searchType" value="ticker" checked />
          Ticker
        </label>
        <label>
          <input type="radio" name="searchType" value="sector" />
          Sector
        </label>
        <label>
          <input type="radio" name="searchType" value="strategy" />
          Strategy
        </label>
        <label>
          <input type="radio" name="searchType" value="date" />
          Date
        </label>
      </div>
    </div>
  </div>
</div>

<style>
  .calendar-page {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
    background-color: white;
    border-radius: 5px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  }
  
  h2, h3 {
    color: #2d3748;
    margin-top: 0;
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
  
  .strategy-guide {
    background-color: #f7fafc;
    padding: 1rem;
    border-radius: 5px;
    margin-bottom: 2rem;
  }
  
  .strategy-legend {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-top: 0.5rem;
  }
  
  .strategy-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  
  .color-box {
    width: 16px;
    height: 16px;
    border-radius: 4px;
  }
  
  .strategy-name {
    font-size: 0.9rem;
  }
  
  .calendar-container {
    overflow-x: auto;
    margin-bottom: 2rem;
  }
  
  .calendar {
    width: 100%;
    border-collapse: collapse;
    border: 1px solid #e2e8f0;
  }
  
  .sector-header {
    width: 150px;
    text-align: left;
    background-color: #f7fafc;
  }
  
  .week-header {
    min-width: 120px;
    background-color: #f7fafc;
    text-align: center;
    font-size: 0.9rem;
    padding: 0.5rem;
  }
  
  .week-label {
    font-weight: bold;
  }
  
  .exp-date {
    font-size: 0.8rem;
    color: #718096;
  }
  
  .sector-cell {
    padding: 0.75rem;
    font-weight: 500;
    border: 1px solid #e2e8f0;
    background-color: #f7fafc;
  }
  
  .trade-cell {
    border: 1px solid #e2e8f0;
    padding: 0.5rem;
    vertical-align: top;
    height: 80px;
  }
  
  .trade-card {
    padding: 0.5rem;
    border-radius: 4px;
    margin-bottom: 0.5rem;
    color: white;
    font-size: 0.8rem;
    text-align: center;
  }
  
  .trade-symbol {
    font-weight: bold;
    margin-bottom: 0.25rem;
  }
  
  .trade-type {
    font-size: 0.75rem;
  }
  
  .add-trade-section {
    background-color: white;
    padding: 1.5rem;
    border-radius: 5px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
    margin-bottom: 2rem;
  }
  
  .trade-form {
    margin-top: 1rem;
  }
  
  .form-row {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-bottom: 1rem;
  }
  
  .form-group {
    flex: 1;
    min-width: 150px;
  }
  
  .form-group.full-width {
    width: 100%;
    flex-basis: 100%;
  }
  
  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: bold;
    color: #4a5568;
  }
  
  .form-group input,
  .form-group select,
  .form-group textarea {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #cbd5e0;
    border-radius: 4px;
    font-size: 0.9rem;
  }
  
  .form-group textarea {
    resize: vertical;
    min-height: 100px;
  }
  
  .form-buttons {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 1rem;
  }
  
  .reset-btn, .add-btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-weight: 500;
    cursor: pointer;
  }
  
  .reset-btn {
    background-color: #e2e8f0;
    color: #4a5568;
  }
  
  .add-btn {
    background-color: #4299e1;
    color: white;
  }
  
  .reset-btn:hover {
    background-color: #cbd5e0;
  }
  
  .add-btn:hover {
    background-color: #3182ce;
  }
  
  .trade-history {
    background-color: white;
    padding: 1.5rem;
    border-radius: 5px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  }
  
  .search-section {
    margin-top: 1rem;
  }
  
  .search-input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #cbd5e0;
    border-radius: 4px;
    margin-bottom: 1rem;
  }
  
  .search-filters {
    display: flex;
    gap: 1.5rem;
  }
  
  .search-filters label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }
</style> 