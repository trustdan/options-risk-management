<script>
  import { onMount } from 'svelte';
  import StockForm from './StockForm.svelte';
  import StockAnalytics from './StockAnalytics.svelte';
  
  let ratings = [];
  let dates = [];
  let dateIndex = 0;
  let currentDay = [];
  let activeTab = 'form'; // 'form' or 'analytics'
  let newRating = {};
  
  onMount(async () => {
    try {
      loadAllRatings();
    } catch (error) {
      console.error("Failed to load stock ratings:", error);
    }
  });
  
  async function loadAllRatings() {
    const result = await window.go.main.App.GetStockRatings();
    ratings = result || [];
    
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
  
  async function saveRating(event) {
    try {
      const rating = event.detail;
      await window.go.main.App.SaveStockRating(rating);
      await loadAllRatings();
      
      // Reset the form
      resetNewRating(dates[dateIndex]);
      
      showToast("Stock rating saved successfully!");
    } catch (error) {
      console.error("Failed to save stock rating:", error);
      showToast("Error saving stock rating", true);
    }
  }
  
  async function deleteRating(id) {
    if (!confirm('Are you sure you want to delete this rating?')) return;
    
    try {
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
      class:active={activeTab === 'analytics'} 
      on:click={() => activeTab = 'analytics'}
    >
      Analytics
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
                on:click={() => deleteRating(rating.id)}
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
        on:save={saveRating} 
      />
    </div>
  {:else if activeTab === 'analytics'}
    <StockAnalytics {ratings} />
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