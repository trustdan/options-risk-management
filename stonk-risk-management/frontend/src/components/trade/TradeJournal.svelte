<script>
  import { onMount } from 'svelte';
  import { GetTrades, SaveTrade, DeleteTrade } from '../../../wailsjs/go/main/App';
  import { models } from '../../../wailsjs/go/models';
  
  // Props to receive from parent
  export let trades = [];
  export let strategyOptions = [];
  export let sectors = [];
  export let editExistingTrade = null;
  export let refreshTrades = () => {};

  // Journal entry properties (enhanced trades)
  let journalEntries = [];
  let currentEntry = {
    id: '',
    date: new Date().toISOString().split('T')[0],
    title: '',
    symbol: '',
    sector: '',
    strategy: '',
    type: '',
    outcome: 'win', // win, loss, breakeven
    pnlAmount: 0,
    emotionalState: '',
    lessonsLearned: '',
    whatWentWell: '',
    whatWentPoorly: '',
    improvementPlan: ''
  };
  
  $: {
    // When editExistingTrade changes, populate the form
    if (editExistingTrade) {
      prepopulateFromTrade(editExistingTrade);
    }
  }
  
  // Function to convert Trade to Journal Entry format
  function prepopulateFromTrade(trade) {
    if (!trade) return;
    
    let entryDate = trade.entryDate;
    if (!(entryDate instanceof Date)) {
      try {
        entryDate = new Date(trade.entryDate);
      } catch (e) {
        console.error("Failed to parse entry date:", trade.entryDate);
        entryDate = new Date();
      }
    }
    
    currentEntry = {
      id: trade.id || '',
      date: entryDate.toISOString().split('T')[0],
      title: trade.notes || '',
      symbol: trade.symbol || '',
      sector: trade.sector || '',
      strategy: trade.strategy || '',
      type: trade.type || '',
      outcome: 'win', // Default, as trade doesn't have this yet
      pnlAmount: trade.entryPrice || 0,
      emotionalState: '',
      lessonsLearned: '',
      whatWentWell: '',
      whatWentPoorly: '',
      improvementPlan: ''
    };
  }
  
  onMount(async () => {
    // Load journal entries from backend - enhanced trades with journal data
    await loadJournalEntries();
  });
  
  async function loadJournalEntries() {
    try {
      // In a real implementation, you might have a separate backend function for journal entries
      // For now, we'll just use the trades as a base
      const allTrades = await GetTrades();
      
      // Filter only trades that have journal data
      journalEntries = allTrades.filter(trade => 
        trade.hasJournalEntry || 
        trade.emotionalState || 
        trade.whatWentWell || 
        trade.lessonsLearned
      );
    } catch (error) {
      console.error("Failed to load journal entries:", error);
    }
  }
  
  async function saveEntry() {
    try {
      // Prepare the data to save
      const trade = new models.Trade();
      
      // If editing an existing trade, preserve its core properties
      if (currentEntry.id) {
        // Find the existing trade
        const existingTrade = trades.find(t => t.id === currentEntry.id);
        if (existingTrade) {
          Object.assign(trade, existingTrade);
        }
      }
      
      // Set or update trade properties
      trade.id = currentEntry.id || Date.now().toString();
      trade.symbol = currentEntry.symbol;
      trade.sector = currentEntry.sector;
      
      // Handle strategy and type
      if (currentEntry.strategy.includes(' - ')) {
        const [strategy, type] = currentEntry.strategy.split(' - ');
        trade.strategy = strategy;
        trade.type = type;
      } else {
        trade.strategy = currentEntry.strategy;
        trade.type = currentEntry.type;
      }
      
      // Entry date from journal date
      trade.entryDate = new Date(currentEntry.date);
      
      // P&L amount
      trade.entryPrice = parseFloat(currentEntry.pnlAmount);
      
      // Notes field can store the trade title
      trade.notes = currentEntry.title;
      
      // Add journal-specific fields
      trade.hasJournalEntry = true;
      trade.outcome = currentEntry.outcome;
      trade.emotionalState = currentEntry.emotionalState;
      trade.lessonsLearned = currentEntry.lessonsLearned;
      trade.whatWentWell = currentEntry.whatWentWell;
      trade.whatWentPoorly = currentEntry.whatWentPoorly;
      trade.improvementPlan = currentEntry.improvementPlan;
      
      // Save to backend
      await SaveTrade(trade);
      
      // Refresh journal entries and trades
      await loadJournalEntries();
      if (refreshTrades) refreshTrades();
      
      // Reset the form
      resetForm();
      
      alert('Journal entry saved successfully');
    } catch (error) {
      console.error('Failed to save journal entry:', error);
      alert('Error saving journal entry: ' + error.message);
    }
  }
  
  function editEntry(entry) {
    currentEntry = {...entry};
  }
  
  async function deleteEntry(id) {
    if (confirm('Are you sure you want to delete this journal entry?')) {
      try {
        await DeleteTrade(id);
        await loadJournalEntries();
        if (refreshTrades) refreshTrades();
      } catch (error) {
        console.error('Failed to delete journal entry:', error);
        alert('Error deleting journal entry: ' + error.message);
      }
    }
  }
  
  function resetForm() {
    currentEntry = {
      id: '',
      date: new Date().toISOString().split('T')[0],
      title: '',
      symbol: '',
      sector: '',
      strategy: '',
      type: '',
      outcome: 'win',
      pnlAmount: 0,
      emotionalState: '',
      lessonsLearned: '',
      whatWentWell: '',
      whatWentPoorly: '',
      improvementPlan: ''
    };
  }
  
  // Format date helper
  function formatDate(date) {
    if (!date) return '';
    let dateObj;
    
    try {
      if (typeof date === 'string' || date instanceof Date) {
        dateObj = new Date(date);
      }
    } catch (e) {
      console.error("Failed to format date:", date);
      return 'Invalid Date';
    }
    
    if (!dateObj || isNaN(dateObj.getTime())) return 'Invalid Date';
    
    return dateObj.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  }

  // Additional properties for trade import section
  let searchQuery = '';
  let searchType = 'ticker';
  let showImportSection = true;
  
  // Function to get strategy color based on type
  function getStrategyColor(type) {
    // Match colors from the screenshots
    const type_lower = type ? type.toLowerCase() : '';
    
    // Danger - naked options ahead category
    if (type_lower.includes('short call') || 
        type_lower.includes('short put') ||
        type_lower.includes('cash-secured put') ||
        (type_lower.includes('call ratio') && !type_lower.includes('backspread')) ||
        (type_lower.includes('put ratio') && !type_lower.includes('backspread'))) {
      return '#ff0000'; // Bright red for dangerous strategies
    }
    
    // Other categories
    if (type_lower.includes('diagonal call')) {
      return '#F9A826'; // Orange
    } else if (type_lower.includes('butterfly')) {
      return '#F44336'; // Red
    } else if (type_lower.includes('long call') && !type_lower.includes('butterfly')) {
      return '#2196F3'; // Blue
    } else if (type_lower.includes('bull') || type_lower.includes('bear')) {
      return '#9C27B0'; // Purple
    } else if (type_lower.includes('covered call')) {
      return '#4CAF50'; // Green
    } else {
      return '#607D8B'; // Default blue-grey
    }
  }
  
  // Filtered trades for import display
  $: filteredTrades = trades
    .filter(trade => !trade.hasJournalEntry) // Only show trades that aren't already journals
    .map(trade => {
      // Ensure we have valid date objects for calculations/display
      let entryDate = null;
      let expirationDate = null;
      
      // Handle different potential date formats safely
      if (trade.entryDate) {
        try {
          entryDate = new Date(trade.entryDate);
        } catch (e) {
          console.error("Failed to parse entry date:", trade.entryDate);
        }
      }
      
      if (trade.expirationDate) {
        try {
          expirationDate = new Date(trade.expirationDate);
        } catch (e) {
          console.error("Failed to parse expiration date:", trade.expirationDate);
        }
      }
      
      // Return an object with the proper date formats
      return {
        ...trade,
        entryDate: entryDate,
        expirationDate: expirationDate
      };
    })
    .filter(trade => {
      if (!searchQuery) return true;
      if (!trade.entryDate || !trade.expirationDate) return false; // Filter out invalid dates
      
      const query = searchQuery.toLowerCase();
      
      switch (searchType) {
        case 'ticker':
          return trade.symbol.toLowerCase().includes(query);
        case 'sector':
          return trade.sector.toLowerCase().includes(query);
        case 'strategy':
          return (
            trade.strategy.toLowerCase().includes(query) || 
            trade.type.toLowerCase().includes(query)
          );
        case 'date':
          return (
            trade.entryDate.toLocaleDateString().includes(query) ||
            trade.expirationDate.toLocaleDateString().includes(query)
          );
        default:
          return true;
      }
    })
    .sort((a, b) => {
      // Ensure dates are valid before comparing
      const timeA = a.entryDate ? a.entryDate.getTime() : 0;
      const timeB = b.entryDate ? b.entryDate.getTime() : 0;
      return timeB - timeA;
    });
  
  function importTradeToJournal(trade) {
    // Prepopulate the form with the selected trade
    prepopulateFromTrade(trade);
    
    // Scroll to the top of the form
    document.querySelector('.journal-form').scrollIntoView({ behavior: 'smooth' });
  }
</script>

<div class="journal-container">
  <h2>Trading Journal</h2>
  <p class="description">
    Record and reflect on your trades to build your trading psychology and improve over time.
  </p>
  
  <!-- Journal Entry Form -->
  <div class="journal-form">
    <h3>
      {currentEntry.id ? 'Edit Journal Entry' : 'New Journal Entry'}
    </h3>
    
    <div class="form-row">
      <div class="form-group">
        <label>Date</label>
        <input 
          type="date" 
          bind:value={currentEntry.date}
          class="input-field"
        />
      </div>
      
      <div class="form-group">
        <label>Symbol</label>
        <input 
          type="text" 
          bind:value={currentEntry.symbol}
          placeholder="AAPL, SPY, etc."
          class="input-field"
        />
      </div>
      
      <div class="form-group">
        <label>Sector</label>
        <select bind:value={currentEntry.sector} class="input-field">
          <option value="">Select sector...</option>
          {#each sectors as sector}
            <option value={sector}>{sector}</option>
          {/each}
        </select>
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group full-width">
        <label>Trade Title</label>
        <input 
          type="text" 
          bind:value={currentEntry.title}
          placeholder="Brief description of the trade"
          class="input-field"
        />
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group full-width">
        <label>Strategy</label>
        <select bind:value={currentEntry.strategy} class="input-field">
          <option value="">Select a strategy...</option>
          {#each strategyOptions as category}
            <optgroup label={category.category}>
              {#each category.options as option}
                <option value={`${category.category} - ${option.name}`}>{option.name}</option>
              {/each}
            </optgroup>
          {/each}
        </select>
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group">
        <label>Outcome</label>
        <div class="radio-group">
          <label class="radio-label">
            <input type="radio" bind:group={currentEntry.outcome} value="win">
            <span class="text-win">Win</span>
          </label>
          <label class="radio-label">
            <input type="radio" bind:group={currentEntry.outcome} value="loss">
            <span class="text-loss">Loss</span>
          </label>
          <label class="radio-label">
            <input type="radio" bind:group={currentEntry.outcome} value="breakeven">
            <span class="text-breakeven">Breakeven</span>
          </label>
        </div>
      </div>
      
      <div class="form-group">
        <label>P&L Amount ($)</label>
        <input 
          type="number" 
          step="0.01"
          bind:value={currentEntry.pnlAmount}
          placeholder="$ Profit or Loss"
          class="input-field"
        />
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group full-width">
        <label>Emotional State</label>
        <input 
          type="text" 
          bind:value={currentEntry.emotionalState}
          placeholder="How were you feeling before/during the trade"
          class="input-field"
        />
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group">
        <label>What Went Well</label>
        <textarea 
          bind:value={currentEntry.whatWentWell}
          rows="3"
          placeholder="Aspects of the trade that were executed properly"
          class="input-field"
        ></textarea>
      </div>
      
      <div class="form-group">
        <label>What Went Poorly</label>
        <textarea 
          bind:value={currentEntry.whatWentPoorly}
          rows="3"
          placeholder="Mistakes or areas for improvement"
          class="input-field"
        ></textarea>
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group full-width">
        <label>Lessons Learned</label>
        <textarea 
          bind:value={currentEntry.lessonsLearned}
          rows="3"
          placeholder="Key takeaways from this trade"
          class="input-field"
        ></textarea>
      </div>
    </div>
    
    <div class="form-row">
      <div class="form-group full-width">
        <label>Improvement Plan</label>
        <textarea 
          bind:value={currentEntry.improvementPlan}
          rows="3"
          placeholder="Specific actions to improve future trades"
          class="input-field"
        ></textarea>
      </div>
    </div>
    
    <div class="form-buttons">
      <button class="reset-btn" on:click={resetForm}>
        {currentEntry.id ? 'Cancel' : 'Clear'}
      </button>
      <button class="add-btn" on:click={saveEntry}>
        {currentEntry.id ? 'Update Entry' : 'Save Journal Entry'}
      </button>
    </div>
  </div>
  
  <!-- Journal Entries List -->
  <div class="entries-table">
    <h3>Recent Journal Entries</h3>
    
    {#if journalEntries.length === 0}
      <div class="no-trades">
        <p>No journal entries found. Create your first entry above.</p>
      </div>
    {:else}
      <div class="trades-table-container">
        <table class="trades-table">
          <thead>
            <tr>
              <th>Date</th>
              <th>Symbol</th>
              <th>Strategy</th>
              <th>Title</th>
              <th>Outcome</th>
              <th>P&L</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each journalEntries as entry}
              <tr>
                <td>{formatDate(entry.entryDate)}</td>
                <td class="ticker-cell">{entry.symbol}</td>
                <td>
                  <span class="strategy-pill" style="background-color: {getStrategyColor(entry.type)}">
                    {entry.type}
                  </span>
                </td>
                <td>{entry.notes || 'Untitled'}</td>
                <td>
                  <span class={entry.outcome === 'win' ? 'text-win' : entry.outcome === 'loss' ? 'text-loss' : 'text-breakeven'}>
                    {entry.outcome || 'Unknown'}
                  </span>
                </td>
                <td>${entry.entryPrice?.toFixed(2) || '0.00'}</td>
                <td class="actions-cell">
                  <button 
                    class="action-btn edit-btn"
                    on:click={() => editEntry(entry)}
                    title="Edit Entry"
                  >
                    ✎
                  </button>
                  <button 
                    class="action-btn delete-btn"
                    on:click={() => deleteEntry(entry.id)}
                    title="Delete Entry"
                  >
                    ✕
                  </button>
                </td>
              </tr>
              {#if entry.emotionalState || entry.lessonsLearned}
                <tr class="notes-row">
                  <td colspan="7">
                    <div class="notes-content">
                      {#if entry.emotionalState}
                        <strong>Emotional State:</strong> {entry.emotionalState}<br>
                      {/if}
                      {#if entry.lessonsLearned}
                        <strong>Lessons:</strong> {entry.lessonsLearned}
                      {/if}
                    </div>
                  </td>
                </tr>
              {/if}
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
  
  <!-- Trade Import Section -->
  <div class="import-section">
    <div class="import-header" on:click={() => showImportSection = !showImportSection}>
      <h3>Import Trades to Journal</h3>
      <button class="toggle-btn">
        {showImportSection ? '▲ Hide' : '▼ Show'}
      </button>
    </div>
    
    {#if showImportSection}
      <div class="import-content">
        <p>Select trades from your history to add to your trading journal</p>
        
        <div class="search-section">
          <input 
            type="text" 
            placeholder="Search trades..." 
            class="search-input" 
            bind:value={searchQuery}
          />
          
          <div class="search-filters">
            <label>
              <input type="radio" name="importSearchType" value="ticker" bind:group={searchType} />
              Ticker
            </label>
            <label>
              <input type="radio" name="importSearchType" value="sector" bind:group={searchType} />
              Sector
            </label>
            <label>
              <input type="radio" name="importSearchType" value="strategy" bind:group={searchType} />
              Strategy
            </label>
            <label>
              <input type="radio" name="importSearchType" value="date" bind:group={searchType} />
              Date
            </label>
          </div>
        </div>
        
        {#if filteredTrades.length === 0}
          <div class="no-trades">
            <p>No uninventoried trades found. All your trades are already in your journal or there are no trades yet.</p>
          </div>
        {:else}
          <div class="trades-table-container">
            <table class="trades-table">
              <thead>
                <tr>
                  <th>Ticker</th>
                  <th>Sector</th>
                  <th>Strategy</th>
                  <th>Entry Date</th>
                  <th>Expiration</th>
                  <th>Price</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {#each filteredTrades as trade (trade.id + '_' + (trade.expirationDate?.toISOString() || ''))}
                  <tr>
                    <td class="ticker-cell">{trade.symbol}</td>
                    <td>{trade.sector}</td>
                    <td>
                      <span class="strategy-pill" style="background-color: {getStrategyColor(trade.type)}">
                        {trade.type}
                      </span>
                    </td>
                    <td>{formatDate(trade.entryDate)}</td>
                    <td>{formatDate(trade.expirationDate)}</td>
                    <td>${trade.entryPrice?.toFixed(2) || '0.00'}</td>
                    <td class="actions-cell">
                      <button 
                        class="action-btn import-btn" 
                        on:click={() => importTradeToJournal(trade)}
                        title="Add to Journal"
                      >
                        📝 Journal
                      </button>
                    </td>
                  </tr>
                  {#if trade.notes}
                    <tr class="notes-row">
                      <td colspan="7">
                        <div class="notes-content">
                          <strong>Notes:</strong> {trade.notes}
                        </div>
                      </td>
                    </tr>
                  {/if}
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
      </div>
    {/if}
  </div>
</div>

<style>
  .journal-container {
    margin-bottom: 2rem;
  }
  
  .description {
    text-align: center;
    color: inherit;
    margin-bottom: 1.5rem;
  }

  h2, h3 {
    color: inherit;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 0.5rem;
  }
  
  h3 {
    margin-top: 0;
    margin-bottom: 1rem;
  }
  
  .journal-form, .entries-table {
    background-color: var(--card-bg);
    padding: 1.5rem;
    border-radius: 5px;
    margin-bottom: 1.5rem;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  }

  .form-row {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-bottom: 1rem;
  }
  
  .form-group {
    flex: 1;
    min-width: 200px;
  }
  
  .form-group.full-width {
    width: 100%;
    flex-basis: 100%;
  }
  
  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
    color: inherit;
  }
  
  .input-field {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid var(--input-border);
    background-color: var(--input-bg);
    color: var(--text-color);
    border-radius: 4px;
    transition: border-color 0.3s, background-color 0.3s;
  }
  
  .input-field:focus {
    border-color: var(--primary-button);
    outline: none;
  }
  
  textarea.input-field {
    resize: vertical;
    min-height: 80px;
  }
  
  .radio-group {
    display: flex;
    gap: 1rem;
    margin-top: 0.5rem;
  }
  
  .radio-label {
    display: inline-flex;
    align-items: center;
    cursor: pointer;
    gap: 0.25rem;
  }
  
  .text-win {
    color: var(--success-color);
    font-weight: 500;
  }
  
  .text-loss {
    color: var(--danger-color);
    font-weight: 500;
  }
  
  .text-breakeven {
    color: var(--warning-color);
    font-weight: 500;
  }
  
  .form-buttons {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 1.5rem;
  }
  
  .reset-btn, .add-btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-weight: 500;
    cursor: pointer;
  }
  
  .reset-btn {
    background-color: var(--border-color);
    color: inherit;
  }
  
  .add-btn {
    background-color: var(--active-button);
    color: white;
  }
  
  .reset-btn:hover {
    background-color: var(--secondary-button-hover);
  }
  
  .add-btn:hover {
    background-color: var(--primary-button-hover);
  }
  
  .no-trades {
    text-align: center;
    color: var(--text-muted);
    padding: 2rem 0;
  }
  
  .trades-table-container {
    overflow-x: auto;
  }
  
  .trades-table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 1rem;
  }
  
  .trades-table th {
    text-align: left;
    padding: 0.75rem;
    background-color: var(--bg-color);
    color: inherit;
    font-weight: 600;
    border-bottom: 2px solid var(--border-color);
  }
  
  .trades-table td {
    padding: 0.75rem;
    border-bottom: 1px solid var(--border-color);
    color: inherit;
  }
  
  .ticker-cell {
    font-weight: 600;
  }
  
  .strategy-pill {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: 9999px;
    font-size: 0.8rem;
    font-weight: 500;
    color: white;
  }
  
  .actions-cell {
    display: flex;
    gap: 0.5rem;
  }
  
  .action-btn {
    width: 28px;
    height: 28px;
    border: none;
    border-radius: 4px;
    color: white;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .edit-btn {
    background-color: var(--primary-button);
  }
  
  .delete-btn {
    background-color: var(--danger-color);
  }
  
  .notes-row {
    background-color: var(--card-bg-light);
  }
  
  .notes-content {
    padding: 0.5rem 0;
    font-size: 0.9rem;
    color: inherit;
  }
  
  /* Import Section Styles */
  .import-section {
    background-color: var(--card-bg);
    padding: 1rem;
    border-radius: 5px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
    margin-top: 2rem;
  }
  
  .import-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 4px;
    transition: background-color 0.2s;
  }
  
  .import-header:hover {
    background-color: var(--bg-color);
  }
  
  .import-header h3 {
    margin: 0;
  }
  
  .toggle-btn {
    border: none;
    background: transparent;
    color: var(--text-color);
    cursor: pointer;
    padding: 0.25rem 0.5rem;
    border-radius: 3px;
    transition: background-color 0.2s;
  }
  
  .toggle-btn:hover {
    background-color: var(--border-color);
  }
  
  .import-content {
    padding: 1rem;
    border-top: 1px solid var(--border-color);
    margin-top: 1rem;
  }
  
  .import-content p {
    margin-top: 0;
    margin-bottom: 1rem;
    text-align: center;
    font-style: italic;
    color: var(--text-secondary);
  }
  
  .search-section {
    margin-top: 1rem;
    margin-bottom: 1.5rem;
  }
  
  .search-input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid var(--input-border);
    border-radius: 4px;
    margin-bottom: 1rem;
    background-color: var(--input-bg);
    color: var(--text-color);
  }
  
  .search-filters {
    display: flex;
    gap: 1.5rem;
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .search-filters label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }
  
  .import-btn {
    background-color: var(--primary-button);
    font-size: 0.8rem;
    padding: 0.25rem 0.75rem;
    height: auto;
    width: auto;
    white-space: nowrap;
  }
</style> 