<script>
  import { onMount } from 'svelte';
  import StockForm from './StockForm.svelte';
  import MarketAnalytics from './MarketAnalytics.svelte';
  
  let ratings = [];
  let marketAndSectorRatings = []; // Filtered ratings for market analytics
  let dates = [];
  let dateIndex = 0;
  let currentDay = [];
  let activeTab = 'form'; // 'form' or 'market'
  let newRating = {};
  
  onMount(async () => {
    try {
      await loadAllRatings();
    } catch (error) {
      console.error("Failed to load stock ratings:", error);
    }
  });
  
  async function loadAllRatings() {
    const result = await window.go.main.App.GetStockRatings();
    ratings = result || [];
    
    // Filter market and sector-specific ratings
    marketAndSectorRatings = ratings.filter(r => 
      r.symbol === 'MARKET' || 
      r.symbol === 'SECTOR' || 
      (r.notes && (r.notes.includes('marketRating') || r.notes.includes('sectorRating')))
    );
    
    // Group by date
    const byDate = groupByDate(ratings);
    dates = Object.keys(byDate).sort();
    
    if (dates.length) {
      loadDate(dates.length - 1); // Load most recent by default
    } else {
      resetNewRating();
      currentDay = [];
    }
  }
  
  function groupByDate(allRatings) {
    const grouped = {};
    
    allRatings.forEach(rating => {
      // Format the date to YYYY-MM-DD
      const date = new Date(rating.date);
      const dateStr = date.toISOString().split('T')[0];
      
      if (!grouped[dateStr]) {
        grouped[dateStr] = [];
      }
      
      grouped[dateStr].push(rating);
    });
    
    return grouped;
  }
  
  function loadDate(idx) {
    if (idx >= 0 && idx < dates.length) {
      dateIndex = idx;
      const dateStr = dates[idx];
      const byDate = groupByDate(ratings);
      currentDay = byDate[dateStr] || [];
      resetNewRating(dateStr);
    }
  }
  
  function resetNewRating(dateStr = null) {
    newRating = {
      id: '',
      date: dateStr || new Date().toISOString().split('T')[0],
      symbol: '',
      sector: '',
      stockSentiment: 5,
      priceTarget: 0,
      confidence: 5,
      notes: ''
    };
  }
  
  function nextDay() {
    if (dateIndex < dates.length - 1) loadDate(dateIndex + 1);
  }
  
  function prevDay() {
    if (dateIndex > 0) loadDate(dateIndex - 1);
  }
  
  // Refresh ratings after saving
  async function refreshRatings() {
    try {
      const result = await window.go.main.App.GetStockRatings();
      ratings = result || [];
      
      // Update market and sector filtered ratings
      marketAndSectorRatings = ratings.filter(r => 
        r.symbol === 'MARKET' || 
        r.symbol === 'SECTOR' || 
        (r.notes && (r.notes.includes('marketRating') || r.notes.includes('sectorRating')))
      );
      
      // Group by date
      const byDate = groupByDate(ratings);
      dates = Object.keys(byDate).sort();
      
      if (dates.length) {
        // Stay on current date if possible
        const currentDateStr = newRating.date;
        const currentDateIndex = dates.indexOf(currentDateStr);
        
        if (currentDateIndex >= 0) {
          loadDate(currentDateIndex);
        } else {
          loadDate(dates.length - 1); // Load most recent
        }
      } else {
        resetNewRating();
        currentDay = [];
      }
    } catch (error) {
      console.error('Failed to refresh ratings:', error);
    }
  }
  
  async function saveStockRating() {
    try {
      // Ensure enthusiasm is a number
      const enthusiasmValue = parseInt(newRating.enthusiasm?.toString() || '0', 10) || 0;
      
      // Create a compatible object for saving
      const ratingToSave = {
        id: newRating.id || '',
        date: new Date(newRating.date),
        symbol: newRating.symbol,
        sector: newRating.sector,
        stockSentiment: newRating.stockSentiment,
        priceTarget: newRating.priceTarget || 0,
        confidence: newRating.confidence || 0,
        // Explicitly set enthusiasm
        enthusiasm: enthusiasmValue,
        // Convert chartPatterns array to a string if they exist
        chartPattern: newRating.chartPatterns ? newRating.chartPatterns.join(', ') : '',
        // Store original data in notes if needed
        notes: newRating.notes
      };

      console.log('Saving rating with explicit fields:', ratingToSave);
      await window.go.main.App.SaveStockRating(ratingToSave);
      
      // Refresh ratings
      await refreshRatings();
      
      // Clear form for next entry
      resetNewRating(newRating.date);
      
      showToast('Stock rating saved successfully');
    } catch (error) {
      console.error('Failed to save rating:', error);
      showToast('Failed to save rating: ' + error.message, true);
    }
  }
  
  async function deleteRating(id) {
    if (!confirm('Are you sure you want to delete this rating?')) return;
    
    try {
      console.log('Deleting rating with ID:', id);
      if (!id) {
        console.error('No ID provided for deletion');
        showToast("Error: No ID provided for deletion", true);
        return;
      }
      await window.go.main.App.DeleteStockRating(id);
      await loadAllRatings();
      showToast("Rating deleted successfully!");
    } catch (error) {
      console.error("Failed to delete rating:", error);
      showToast("Error deleting rating", true);
    }
  }
  
  function showToast(message, isError = false) {
    const toast = document.createElement('div');
    toast.className = `toast ${isError ? 'error' : 'success'}`;
    toast.textContent = message;
    document.body.appendChild(toast);
    
    setTimeout(() => {
      toast.classList.add('show');
      setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => {
          document.body.removeChild(toast);
        }, 300);
      }, 3000);
    }, 100);
  }
  
  // Enhance the logRecentRatings function for better debugging
  function logRecentRatings() {
    console.log("Recent ratings:", recentRatings);
    if (recentRatings.length > 0) {
      recentRatings.forEach(r => {
        console.log(`Rating ID: ${r.id}`);
        console.log(`Symbol: ${r.symbol}, Sector: ${r.sector}`);
        console.log(`Sentiment: ${r.stockSentiment}, Enthusiasm: ${r.enthusiasm} (type: ${typeof r.enthusiasm})`);
        console.log(`Chart Pattern: ${r.chartPattern} (type: ${typeof r.chartPattern})`);
        console.log('Full object:', r);
        console.log('---');
      });
    }
  }
</script>

<div class="dashboard">
  <div class="tabs">
    <button 
      class:active={activeTab === 'form'} 
      on:click={() => activeTab = 'form'}
    >
      Stock Ratings
    </button>
    <button 
      class:active={activeTab === 'market'} 
      on:click={() => activeTab = 'market'}
    >
      Market Analytics
    </button>
  </div>
  
  {#if activeTab === 'form'}
    <div class="date-navigation">
      <button 
        on:click={prevDay} 
        disabled={dateIndex === 0 && dates.length > 0}
        class="nav-button"
      >
        ‹ Previous Day
      </button>
      
      <span class="current-date">
        {dates.length ? new Date(dates[dateIndex]).toLocaleDateString('en-US', {
          weekday: 'short',
          month: 'short', 
          day: 'numeric',
          year: 'numeric'
        }) : 'New Entry'}
      </span>
      
      <button 
        on:click={nextDay} 
        disabled={dateIndex === dates.length - 1 && dates.length > 0}
        class="nav-button"
      >
        Next Day ›
      </button>
    </div>
    
    <div class="day-ratings">
      {#if currentDay.length > 0}
        <h3>Existing Ratings ({currentDay.length})</h3>
        <div class="ratings-list">
          {#each currentDay as rating}
            <div class="rating-card">
              <div class="rating-header">
                <h4>{rating.symbol}</h4>
                <span class="sector">{rating.sector}</span>
              </div>
              <div class="rating-body">
                <div class="rating-stats">
                  <div class="stat">
                    <span class="label">Sentiment:</span>
                    <span class="value">{rating.stockSentiment}</span>
                  </div>
                  <div class="stat">
                    <span class="label">Confidence:</span>
                    <span class="value">{rating.confidence}</span>
                  </div>
                  {#if rating.priceTarget > 0}
                    <div class="stat">
                      <span class="label">Target:</span>
                      <span class="value">${rating.priceTarget.toFixed(2)}</span>
                    </div>
                  {/if}
                </div>
                {#if rating.notes}
                  <div class="notes">{rating.notes}</div>
                {/if}
              </div>
              <button 
                class="delete-btn" 
                on:click={() => {
                  console.log('Delete clicked for rating:', rating);
                  console.log('ID to delete:', rating.id);
                  deleteRating(rating.id);
                }}
              >
                Delete
              </button>
            </div>
          {/each}
        </div>
      {/if}
      
      <h3>Add New Rating</h3>
      <StockForm 
        bind:data={newRating} 
        on:save={saveStockRating} 
      />
    </div>
  {:else if activeTab === 'market'}
    <MarketAnalytics ratings={marketAndSectorRatings} />
  {/if}
</div>

<style>
  .dashboard {
    margin: 0 auto;
    padding: 1rem;
  }
  
  .tabs {
    display: flex;
    margin-bottom: 1rem;
    border-bottom: 1px solid #ddd;
  }
  
  .tabs button {
    padding: 0.75rem 1.5rem;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 1rem;
    border-bottom: 3px solid transparent;
    transition: all 0.2s;
  }
  
  .tabs button.active {
    font-weight: bold;
    border-bottom: 3px solid #4caf50;
  }
  
  .date-navigation {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    padding: 0.5rem;
    background-color: #f9f9f9;
    border-radius: 4px;
  }
  
  .nav-button {
    background-color: #e0e0e0;
    border: none;
    padding: 0.5rem 1rem;
    cursor: pointer;
    border-radius: 4px;
    transition: background-color 0.2s;
  }
  
  .nav-button:hover:not(:disabled) {
    background-color: #d0d0d0;
  }
  
  .nav-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .current-date {
    font-weight: bold;
    font-size: 1.1rem;
  }
  
  h3 {
    margin-top: 1.5rem;
    margin-bottom: 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 1px solid #eee;
  }
  
  .ratings-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
  }
  
  .rating-card {
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 1rem;
    position: relative;
  }
  
  .rating-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }
  
  .rating-header h4 {
    margin: 0;
    font-size: 1.2rem;
  }
  
  .sector {
    font-size: 0.8rem;
    background-color: #f0f0f0;
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
  }
  
  .rating-stats {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }
  
  .stat {
    display: flex;
    align-items: center;
    gap: 0.25rem;
  }
  
  .label {
    font-weight: bold;
    font-size: 0.9rem;
  }
  
  .value {
    font-size: 1rem;
  }
  
  .notes {
    font-size: 0.9rem;
    font-style: italic;
    color: #666;
    margin-top: 0.5rem;
    border-top: 1px solid #eee;
    padding-top: 0.5rem;
  }
  
  .delete-btn {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    background: none;
    border: none;
    font-size: 0.8rem;
    color: #f44336;
    cursor: pointer;
    opacity: 0.5;
    transition: opacity 0.2s;
  }
  
  .delete-btn:hover {
    opacity: 1;
  }
</style> 