<script>
  import { onMount } from 'svelte';
  import { DeleteStockRating } from '../../../wailsjs/go/main/App';
  
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
  let marketAndSectorRatings = [];
  let ratingsByDate = [];
  
  $: {
    if (ratings && ratings.length) {
      // Filter out market and sector ratings for the table
      marketAndSectorRatings = ratings.filter(r => 
        r.symbol === 'MARKET' || r.symbol === 'SECTOR'
      ).sort((a, b) => {
        const dateA = new Date(a.date).getTime();
        const dateB = new Date(b.date).getTime();
        return dateB - dateA;
      });
      
      // Group ratings by date for summary table
      ratingsByDate = groupRatingsByDate(marketAndSectorRatings);
      
      initializeSectorToggles();
      if (chartCanvas) {
        renderChart();
      }
    } else {
      // Initialize to empty array when no ratings available
      ratingsByDate = [];
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
  
  // Function to group ratings by date
  function groupRatingsByDate(ratings) {
    // Create object to hold grouped data
    const grouped = {};
    
    ratings.forEach(rating => {
      const date = new Date(rating.date);
      // Get just the date part in ISO format (YYYY-MM-DD)
      const dateStr = date.toISOString().split('T')[0];
      
      if (!grouped[dateStr]) {
        grouped[dateStr] = {
          date: dateStr,
          formattedDate: formatDateShort(date),
          ratings: [],
          marketCount: 0,
          sectorCount: 0
        };
      }
      
      grouped[dateStr].ratings.push(rating);
      
      if (rating.symbol === 'MARKET') {
        grouped[dateStr].marketCount++;
      } else if (rating.symbol === 'SECTOR') {
        grouped[dateStr].sectorCount++;
      }
    });
    
    // Convert to array and sort by date (descending)
    const result = [];
    for (const dateStr in grouped) {
      result.push(grouped[dateStr]);
    }
    
    return result.sort((a, b) => {
      const dateA = new Date(a.date).getTime();
      const dateB = new Date(b.date).getTime();
      return dateB - dateA;
    });
  }
  
  // Handle bulk deletion of ratings for a specific date
  async function deleteRatingsByDate(dateStr, ratings) {
    if (!confirm(`Are you sure you want to delete ALL ratings for ${formatDateShort(new Date(dateStr))}?`)) {
      return;
    }
    
    try {
      const idsToDelete = ratings.map(r => r.id);
      console.log(`Deleting ${idsToDelete.length} ratings for date ${dateStr}`);
      
      let successCount = 0;
      let errorCount = 0;
      
      // Delete each rating individually
      for (const id of idsToDelete) {
        try {
          await DeleteStockRating(id);
          successCount++;
        } catch (error) {
          console.error(`Failed to delete rating ${id}:`, error);
          errorCount++;
        }
      }
      
      // Update the arrays to remove the deleted ratings
      if (successCount > 0) {
        // Remove all ratings for this date from our ratings array
        const deletedIds = new Set(idsToDelete);
        ratings = ratings.filter(r => !deletedIds.has(r.id));
        marketAndSectorRatings = marketAndSectorRatings.filter(r => !deletedIds.has(r.id));
        
        // Regenerate the grouped ratings
        ratingsByDate = groupRatingsByDate(marketAndSectorRatings);
        
        // Refresh the chart
        renderChart();
      }
      
      // Show results
      if (errorCount === 0) {
        alert(`Successfully deleted all ${successCount} ratings for ${formatDateShort(new Date(dateStr))}`);
      } else {
        alert(`Deleted ${successCount} ratings, but failed to delete ${errorCount} ratings for ${formatDateShort(new Date(dateStr))}`);
      }
    } catch (error) {
      console.error('Failed to delete ratings:', error);
      alert('An error occurred while deleting ratings: ' + error.message);
    }
  }
  
  // Format date for display (full format)
  function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
  
  // Format date for display (shorter format)
  function formatDateShort(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
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
    
    <div class="ratings-table-container">
      <h3>All Rating Entries</h3>
      <p class="table-instructions">View and manage all market and sector ratings by date. Delete all ratings for a given day.</p>
      
      <div class="ratings-table">
        <table>
          <thead>
            <tr>
              <th>Date</th>
              <th>Market Ratings</th>
              <th>Sector Ratings</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {#if ratingsByDate.length > 0}
              {#each ratingsByDate as dateGroup}
                <tr>
                  <td>{dateGroup.formattedDate}</td>
                  <td>{dateGroup.marketCount}</td>
                  <td>{dateGroup.sectorCount}</td>
                  <td>
                    <button class="delete-btn" on:click={() => deleteRatingsByDate(dateGroup.date, dateGroup.ratings)}>
                      Delete All
                    </button>
                  </td>
                </tr>
              {/each}
            {:else}
              <tr>
                <td colspan="4">No market or sector ratings available</td>
              </tr>
            {/if}
          </tbody>
        </table>
      </div>
    </div>
  {/if}
</div>

<style>
  .analytics {
    max-width: 1000px;
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
  
  .table-instructions {
    font-size: 0.9rem;
    color: var(--text-color-muted);
    text-align: center;
    margin-bottom: 1rem;
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
  
  .ratings-table-container {
    background-color: var(--card-bg);
    border-radius: 8px;
    padding: 1rem;
    margin-top: 2rem;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
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
    background-color: var(--bg-color);
  }
  
  th, td {
    padding: 0.75rem 1rem;
    text-align: left;
    border-bottom: 1px solid var(--border-color);
    color: inherit;
  }
  
  .positive {
    color: #38a169;
  }
  
  .negative {
    color: #e53e3e;
  }
  
  .neutral {
    color: inherit;
    opacity: 0.8;
  }
  
  .delete-btn {
    background-color: #e53e3e;
    color: white;
    border: none;
    border-radius: 4px;
    padding: 0.25rem 0.5rem;
    font-size: 0.8rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .delete-btn:hover {
    opacity: 0.9;
  }
  
  /* Keep existing styles... */
</style> 