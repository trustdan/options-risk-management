<script>
  import { createEventDispatcher } from 'svelte';
  
  export let data = {
    id: '',
    date: new Date().toISOString().split('T')[0],
    symbol: '',
    sector: '',
    stockSentiment: 5,
    priceTarget: 0,
    confidence: 5,
    notes: ''
  };
  
  // Common stock sectors for dropdown
  const sectors = [
    'Technology',
    'Healthcare',
    'Financials',
    'Consumer Discretionary',
    'Communication Services',
    'Industrials',
    'Consumer Staples',
    'Energy',
    'Utilities',
    'Real Estate',
    'Materials'
  ];
  
  const dispatch = createEventDispatcher();
  
  function handleSubmit() {
    dispatch('save', data);
  }
</script>

<div class="stock-form">
  <h2>Stock Rating</h2>
  
  <div class="form-row">
    <div class="form-group">
      <label for="date">Date</label>
      <input 
        type="date" 
        id="date" 
        bind:value={data.date}
      />
    </div>
    
    <div class="form-group">
      <label for="symbol">Symbol</label>
      <input 
        type="text" 
        id="symbol" 
        bind:value={data.symbol}
        placeholder="e.g. AAPL"
        maxlength="10"
      />
    </div>
  </div>
  
  <div class="form-group">
    <label for="sector">Sector</label>
    <select id="sector" bind:value={data.sector}>
      <option value="">Select Sector</option>
      {#each sectors as sector}
        <option value={sector}>{sector}</option>
      {/each}
    </select>
  </div>
  
  <div class="form-group">
    <label for="sentiment">Stock Sentiment (1-10)</label>
    <div class="slider-with-value">
      <input 
        type="range" 
        id="sentiment" 
        min="1" 
        max="10" 
        bind:value={data.stockSentiment}
      />
      <span class="slider-value">{data.stockSentiment}</span>
    </div>
  </div>
  
  <div class="form-row">
    <div class="form-group">
      <label for="price-target">Price Target</label>
      <input 
        type="number" 
        id="price-target" 
        bind:value={data.priceTarget}
        step="0.01"
        min="0"
      />
    </div>
    
    <div class="form-group">
      <label for="confidence">Confidence (1-10)</label>
      <div class="slider-with-value">
        <input 
          type="range" 
          id="confidence" 
          min="1" 
          max="10" 
          bind:value={data.confidence}
        />
        <span class="slider-value">{data.confidence}</span>
      </div>
    </div>
  </div>
  
  <div class="form-group">
    <label for="notes">Notes</label>
    <textarea 
      id="notes" 
      bind:value={data.notes}
      placeholder="Why do you feel this way about the stock?"
    ></textarea>
  </div>
  
  <button on:click={handleSubmit}>Save Rating</button>
</div>

<style>
  .stock-form {
    max-width: 600px;
    margin: 0 auto;
    padding: 1rem;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 1.5rem;
  }
  
  .form-row {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
  }
  
  .form-row .form-group {
    flex: 1;
  }
  
  .form-group {
    margin-bottom: 1rem;
    display: flex;
    flex-direction: column;
  }
  
  label {
    margin-bottom: 0.5rem;
    font-weight: bold;
  }
  
  input, select, textarea {
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
  }
  
  .slider-with-value {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  
  input[type="range"] {
    flex-grow: 1;
  }
  
  .slider-value {
    font-weight: bold;
    font-size: 1.2rem;
    min-width: 2rem;
    text-align: center;
  }
  
  textarea {
    min-height: 100px;
    resize: vertical;
  }
  
  button {
    background-color: #4caf50;
    color: white;
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
    margin-top: 1rem;
    display: block;
    width: 100%;
  }
  
  button:hover {
    background-color: #45a049;
  }
</style> 