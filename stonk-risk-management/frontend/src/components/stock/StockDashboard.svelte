<script>
  import { onMount } from 'svelte';
  import { GetStockRatings, SaveStockRating, DeleteStockRating } from '../../../wailsjs/go/main/App';
  import { models } from '../../../wailsjs/go/models';
  
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
    chartPatterns: [],
    selectedPattern: '',
    enthusiasm: 0,
    notes: ''
  };
  
  // Initialize market ratings object
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
  
  // For date navigation
  let availableDates = [];
  let currentDateIndex = 0;
  let allRatingsByDate = {};
  
  // For tracking save status
  let saveFeedback = {
    market: '',
    sectors: {}
  };
  
  // Initialize feedback for each sector
  sectors.forEach(sector => {
    saveFeedback.sectors[sector] = '';
  });
  
  // For chart pattern cheat sheet
  let showPatternCheatSheet = false;
  
  onMount(async () => {
    try {
      const ratings = await GetStockRatings();
      console.log("Loaded ratings:", ratings);
      
      if (ratings && ratings.length > 0) {
        // Group by date to show the most recent
        allRatingsByDate = groupByDate(ratings);
        availableDates = Object.keys(allRatingsByDate).sort();
        console.log("Available dates:", availableDates);
        console.log("Ratings by date:", allRatingsByDate);
        
        if (availableDates.length > 0) {
          currentDateIndex = availableDates.length - 1; // Start with most recent date
          const latestDate = availableDates[currentDateIndex];
          console.log("Current date index:", currentDateIndex, "of", availableDates.length);
          console.log("Current date:", latestDate);
          
          loadDateData(latestDate);
        } else {
          console.log("No dates available to load");
        }
      } else {
        console.log("No ratings available");
      }
    } catch (error) {
      console.error('Failed to load stock ratings:', error);
    }
  });
  
  function loadDateData(dateStr) {
    recentRatings = allRatingsByDate[dateStr] || [];
    dateString = dateStr;
    rating.date = dateStr;
    
    // Reset market ratings
    marketRatings.overall = 0;
    sectors.forEach(sector => {
      marketRatings.sectorRatings[sector] = 0;
    });
    
    console.log('Loading ratings for date:', dateStr);
    console.log('Found ratings:', recentRatings);
    
    // First pass: filter out only the Market and Sector specific ratings for easier processing
    const marketRatingRecords = recentRatings.filter(r => r.symbol === 'MARKET' || (r.notes && r.notes.includes('marketRating')));
    const sectorRatingRecords = {};
    
    // Group sector ratings by sector name for easier processing
    recentRatings.forEach(r => {
      // Handle records with SECTOR symbol (new format)
      if (r.symbol === 'SECTOR' && sectors.includes(r.sector)) {
        if (!sectorRatingRecords[r.sector]) {
          sectorRatingRecords[r.sector] = [];
        }
        sectorRatingRecords[r.sector].push(r);
      }
      // Handle legacy format with sectorRating in notes
      else if (r.sector && sectors.includes(r.sector) && r.notes && r.notes.includes('sectorRating')) {
        if (!sectorRatingRecords[r.sector]) {
          sectorRatingRecords[r.sector] = [];
        }
        sectorRatingRecords[r.sector].push(r);
      }
    });
    
    console.log('Market rating records:', marketRatingRecords);
    console.log('Sector rating records:', sectorRatingRecords);
    
    // Process market ratings - prioritize MARKET symbol records
    if (marketRatingRecords.length > 0) {
      // First try to find a record with MARKET symbol
      const marketRecord = marketRatingRecords.find(r => r.symbol === 'MARKET');
      if (marketRecord) {
        console.log('Using MARKET symbol record for market rating:', marketRecord);
        marketRatings.overall = marketRecord.stockSentiment;
      } else {
        // Otherwise, try to parse from notes
        for (const record of marketRatingRecords) {
          try {
            const parsedNotes = JSON.parse(record.notes);
            if (parsedNotes.marketRating !== undefined) {
              console.log('Using notes field for market rating:', parsedNotes.marketRating);
              marketRatings.overall = parsedNotes.marketRating;
              break; // Use the first valid record
            }
          } catch (e) {
            console.log('Error parsing notes for market rating:', e);
          }
        }
      }
    }
    
    // Process sector ratings
    sectors.forEach(sector => {
      const sectorRecords = sectorRatingRecords[sector] || [];
      
      // Debug log all records for this sector
      console.log(`Processing records for sector ${sector}:`, sectorRecords);
      
      if (sectorRecords.length > 0) {
        // First try to find a record with SECTOR symbol and stockSentiment value
        const sectorRecord = sectorRecords.find(r => r.symbol === 'SECTOR' && r.stockSentiment !== undefined);
        
        if (sectorRecord) {
          console.log(`Using SECTOR symbol record for ${sector} rating:`, sectorRecord);
          console.log(`Setting ${sector} rating to:`, sectorRecord.stockSentiment);
          marketRatings.sectorRatings[sector] = sectorRecord.stockSentiment;
        } else {
          // Otherwise, try to parse from notes
          for (const record of sectorRecords) {
            try {
              if (record.notes) {
                const parsedNotes = JSON.parse(record.notes);
                if (parsedNotes.type === 'sector_rating' && parsedNotes.value !== undefined) {
                  console.log(`Using type/value notes field for ${sector} rating:`, parsedNotes.value);
                  marketRatings.sectorRatings[sector] = parsedNotes.value;
                  break;
                } else if (parsedNotes.sectorRating !== undefined) {
                  console.log(`Using sectorRating notes field for ${sector} rating:`, parsedNotes.sectorRating);
                  marketRatings.sectorRatings[sector] = parsedNotes.sectorRating;
                  break;
                }
              }
            } catch (e) {
              console.log(`Error parsing notes for ${sector} rating:`, e);
            }
          }
        }
      }
    });
    
    // Debug log the loaded ratings
    logRecentRatings();
    console.log('Market and sector ratings loaded:', marketRatings);
    
    // Force Svelte to recognize the changes
    marketRatings = {...marketRatings};
  }
  
  function navigateDate(direction) {
    if (availableDates.length === 0) return;
    
    const newIndex = currentDateIndex + direction;
    if (newIndex >= 0 && newIndex < availableDates.length) {
      currentDateIndex = newIndex;
      loadDateData(availableDates[currentDateIndex]);
    }
  }
  
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
  
  // Save an individual sector rating
  async function saveSectorRating(sector, noRefresh) {
    try {
      const ratingValue = marketRatings.sectorRatings[sector];
      console.log(`Saving ${sector} rating with value ${ratingValue}`);
      
      // Find existing sector rating for this date and sector
      const targetDate = new Date(rating.date).toISOString().split('T')[0];
      const existingSectorRating = recentRatings.find(r => {
        const ratingDate = new Date(r.date).toISOString().split('T')[0];
        return ratingDate === targetDate && 
              r.symbol === 'SECTOR' && 
              r.sector === sector;
      });
      
      console.log(`Existing sector rating for ${sector} on ${targetDate}:`, existingSectorRating);
      
      // Create sector rating with correct format - use ISO string for date
      const sectorRatingData = {
        id: existingSectorRating ? existingSectorRating.id : '', // Use existing ID if found
        date: `${targetDate}T00:00:00Z`, // Use ISO string format
        symbol: 'SECTOR',
        sector: sector,
        stockSentiment: ratingValue,
        priceTarget: 0,
        confidence: 0,
        enthusiasm: 0,
        chartPattern: '',
        notes: JSON.stringify({
          type: 'sector_rating',
          value: ratingValue
        })
      };
      
      // Create a proper StockRating object
      const sectorRatingToSave = models.StockRating.createFrom(sectorRatingData);
      console.log(`Saving sector rating:`, sectorRatingToSave);
      
      // Call API to save the rating
      try {
        await SaveStockRating(sectorRatingToSave);
        console.log(`${sector} rating saved successfully`);
      } catch (saveError) {
        console.error(`API error saving ${sector} rating:`, saveError);
        throw saveError; // Re-throw to be caught by outer try/catch
      }
      
      // Show feedback
      saveFeedback.sectors[sector] = 'Saved!';
      setTimeout(() => { saveFeedback.sectors[sector] = ''; }, 2000);
      
      // Refresh the data (unless this was called from saveAllSectorRatings)
      if (noRefresh !== 'noRefresh') {
        await refreshRatings();
      }
      
      return true; // Return success
    } catch (error) {
      console.error(`Failed to save ${sector} rating:`, error);
      saveFeedback.sectors[sector] = 'Error!';
      setTimeout(() => { saveFeedback.sectors[sector] = ''; }, 2000);
      throw error; // Re-throw so saveAllSectorRatings can catch it
    }
  }
  
  // Save overall market rating
  async function saveMarketRating(noRefresh) {
    try {
      if (!marketRatings.overall && marketRatings.overall !== 0) {
        console.error("Market rating is undefined");
        return;
      }

      // Overall rating is just a number, not an object
      const ratingValue = marketRatings.overall;
      
      // Find existing market rating for this date
      const targetDate = new Date(rating.date).toISOString().split('T')[0];
      const existingMarketRating = recentRatings.find(r => {
        const ratingDate = new Date(r.date).toISOString().split('T')[0];
        return ratingDate === targetDate && 
              r.symbol === 'MARKET' && 
              r.sector === 'MARKET';
      });
      
      console.log(`Existing market rating on ${targetDate}:`, existingMarketRating);
      
      // Prepare the rating object with the existing ID if found - use ISO string for date
      const marketRatingData = {
        id: existingMarketRating ? existingMarketRating.id : '',
        date: `${targetDate}T00:00:00Z`, // Use ISO string format
        symbol: 'MARKET',
        sector: 'MARKET',
        stockSentiment: ratingValue,
        priceTarget: 0,
        confidence: 0,
        enthusiasm: 0,
        chartPattern: '',
        notes: JSON.stringify({
          type: 'market_rating',
          value: ratingValue
        })
      };
      
      // Create a proper StockRating object
      const marketRatingToSave = models.StockRating.createFrom(marketRatingData);
      console.log('Saving market rating:', marketRatingToSave);
      
      // Call API to save the rating
      try {
        await SaveStockRating(marketRatingToSave);
        console.log('Market rating saved successfully');
      } catch (saveError) {
        console.error('API error saving market rating:', saveError);
        throw saveError; // Re-throw to be caught by outer try/catch
      }
      
      // Show feedback
      saveFeedback.market = 'Saved!';
      setTimeout(() => { saveFeedback.market = ''; }, 2000);
      
      // Refresh when saving individually
      if (noRefresh !== 'noRefresh') {
        await refreshRatings();
      }
      
      return true; // Return success
    } catch (error) {
      console.error('Failed to save market rating:', error);
      saveFeedback.market = 'Error!';
      setTimeout(() => { saveFeedback.market = ''; }, 2000);
      throw error; // Re-throw so saveAllSectorRatings can catch it
    }
  }
  
  // Save all market & sector ratings at once
  async function saveAllSectorRatings() {
    try {
      console.log('Starting to save all market and sector ratings...');
      
      // First save the market rating (with noRefresh parameter)
      console.log('Saving market rating...');
      await saveMarketRating('noRefresh');
      
      // Then save each sector rating individually (sequentially)
      console.log('Saving all sector ratings one by one...');
      for (const sector of sectors) {
        try {
          console.log(`Saving ${sector} rating: ${marketRatings.sectorRatings[sector]}`);
          await saveSectorRating(sector, 'noRefresh'); // Pass noRefresh parameter
          console.log(`Successfully saved ${sector} rating`);
        } catch (sectorError) {
          console.error(`Error saving ${sector} rating:`, sectorError);
          // Continue with other sectors even if one fails
        }
      }
      
      // Final refresh to ensure all data is up-to-date
      console.log('Finished saving all ratings, refreshing data...');
      await refreshRatings();
      
      console.log('All market and sector ratings save completed!');
      alert('All market and sector ratings saved successfully!');
    } catch (error) {
      console.error('Error during saving all ratings:', error);
      alert('Error saving market/sector ratings: ' + error.message);
    }
  }
  
  // Refresh ratings after saving
  async function refreshRatings() {
    console.log('Refreshing ratings data...');
    try {
      const ratings = await GetStockRatings();
      console.log('Refreshed ratings data from API:', ratings);
      
      allRatingsByDate = groupByDate(ratings);
      availableDates = Object.keys(allRatingsByDate).sort();
      
      // Find index of current date
      const currentDateStr = rating.date;
      currentDateIndex = availableDates.indexOf(currentDateStr);
      if (currentDateIndex === -1 && availableDates.length > 0) {
        currentDateIndex = availableDates.length - 1;
      }
      
      if (availableDates.length > 0) {
        // Use the date string we found or keep using the current one
        const dateToLoad = currentDateIndex !== -1 ? availableDates[currentDateIndex] : currentDateStr;
        console.log(`Loading date data for: ${dateToLoad} (index: ${currentDateIndex})`);
        loadDateData(dateToLoad);
      } else {
        console.warn('No dates available after refresh');
      }
    } catch (error) {
      console.error('Error refreshing ratings data:', error);
    }
  }
  
  // Add a function to edit an existing rating
  function editRating(item) {
    console.log('Editing rating:', item);
    rating = {
      id: item.id || item.idd || '',
      date: item.date instanceof Date ? item.date.toISOString().split('T')[0] : new Date(item.date).toISOString().split('T')[0],
      symbol: item.symbol || '',
      sector: item.sector || '',
      stockSentiment: item.stockSentiment || 0,
      priceTarget: item.priceTarget || 0,
      confidence: item.confidence || 0,
      enthusiasm: item.enthusiasm || 0,
      notes: item.notes || '',
      chartPatterns: item.chartPattern ? item.chartPattern.split(', ').filter(p => p.trim() !== '') : [],
      selectedPattern: ''
    };
    
    console.log('Populated form with rating:', rating);
    document.querySelector('.stock-rating-container').scrollIntoView({ behavior: 'smooth' });
  }
  
  // Add a function to delete a rating
  async function deleteRating(id) {
    if (confirm('Are you sure you want to delete this rating?')) {
      try {
        console.log('Deleting rating with ID:', id);
        if (!id) {
          console.error('No ID provided for deletion');
          alert('Error: No ID provided for deletion');
          return;
        }
        
        // Actually call the DeleteStockRating API function
        await DeleteStockRating(id);
        console.log('Successfully deleted rating with ID:', id);
        
        // Refresh data after deletion
        await refreshRatings();
        alert('Rating deleted successfully');
      } catch (error) {
        console.error('Failed to delete rating:', error);
        alert('Failed to delete rating: ' + error.message);
      }
    }
  }
  
  // Improve the saveStockRating function to better handle enthusiasm and patterns
  async function saveStockRating() {
    try {
      // Ensure enthusiasm is a number
      const enthusiasmValue = parseInt(rating.enthusiasm.toString(), 10) || 0;
      
      // Get date in ISO format
      const targetDate = new Date(rating.date).toISOString().split('T')[0];
      
      // Create a compatible object for saving - use ISO string for date
      const ratingData = {
        id: rating.id || '',
        date: `${targetDate}T00:00:00Z`, // Use ISO string format
        symbol: rating.symbol,
        sector: rating.sector,
        stockSentiment: rating.stockSentiment,
        priceTarget: rating.priceTarget || 0,
        confidence: rating.confidence || 0,
        // Explicitly set enthusiasm
        enthusiasm: enthusiasmValue,
        // Convert chartPatterns array to a string
        chartPattern: rating.chartPatterns.join(', '),
        // Store original data in notes if needed
        notes: rating.notes
      };

      // Create a proper StockRating object
      const ratingToSave = models.StockRating.createFrom(ratingData);
      console.log('Saving rating with explicit fields:', ratingToSave);
      await SaveStockRating(ratingToSave);
      
      // Refresh ratings
      await refreshRatings();
      
      // Clear form for next entry
      rating = {
        id: '',
        date: rating.date,
        symbol: '',
        sector: '',
        stockSentiment: 0,
        priceTarget: 0,
        confidence: 0,
        chartPatterns: [],
        selectedPattern: '',
        enthusiasm: 0,
        notes: ''
      };
      
      alert('Stock rating saved successfully');
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
      chartPatterns: [],
      selectedPattern: '',
      enthusiasm: 0,
      notes: ''
    };
  }
  
  // Create a dummy rating to enable testing if no ratings exist
  function createDummyRating() {
    const today = new Date().toISOString().split('T')[0];
    rating.date = today;
    rating.symbol = "TEST";
    rating.sector = "Technology";
    rating.stockSentiment = 1;
    rating.enthusiasm = 5;
    
    marketRatings.overall = 2;
    marketRatings.sectorRatings["Technology"] = 3;
    
    saveStockRating();
  }
  
  // Formatted date for display
  $: formattedDate = new Date(rating.date).toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short', 
    day: 'numeric',
    year: 'numeric'
  });
  
  // Make marketRatings reactive
  $: {
    if (marketRatings) {
      console.log('Reactive update of marketRatings:', {...marketRatings});
    }
  }
  
  function togglePatternCheatSheet() {
    showPatternCheatSheet = !showPatternCheatSheet;
  }
  
  function addChartPattern() {
    if (rating.selectedPattern && !rating.chartPatterns.includes(rating.selectedPattern)) {
      rating.chartPatterns = [...rating.chartPatterns, rating.selectedPattern];
      rating.selectedPattern = '';
    }
  }
  
  function removeChartPattern(pattern) {
    rating.chartPatterns = rating.chartPatterns.filter(p => p !== pattern);
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
  <h2>Stock Rating Dashboard</h2>
  <div class="divider"></div>
  
  <p class="description">Rate market conditions, sectors, and individual stocks to identify the best trading opportunities.</p>
  
  <div class="date-navigation">
    <button 
      class="nav-button" 
      on:click={() => navigateDate(-1)}
      disabled={currentDateIndex <= 0}
    >
      ← Previous Day
    </button>
    <div class="current-date">{formattedDate}</div>
    <button 
      class="nav-button" 
      on:click={() => navigateDate(1)}
      disabled={currentDateIndex >= availableDates.length - 1}
    >
      Next Day →
    </button>
  </div>
  
  {#if availableDates.length <= 1}
    <div class="help-message">
      <p>No historical data to navigate. Save some ratings to enable navigation.</p>
      <button class="debug-btn" on:click={createDummyRating}>Create Test Rating</button>
    </div>
  {/if}
  
  <!-- Market & Sectors Section (Top) -->
  <div class="market-sectors-container">
    <div class="date-picker">
      <label>Rating Date:</label>
      <input type="date" bind:value={rating.date} />
    </div>
    
    <div class="section-header">
      <h3>Market & Sectors</h3>
      <button class="save-all-btn" on:click={saveAllSectorRatings}>Save All Ratings</button>
    </div>
    
    <!-- Overall Market Rating -->
    <div class="market-rating-box">
      <div class="market-header">
        <h4>Overall Market: {marketRatings.overall}</h4>
        <button class="save-sector-btn" on:click={() => saveMarketRating()}>
          Save {#if saveFeedback.market}<span class="feedback">{saveFeedback.market}</span>{/if}
        </button>
      </div>
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
    
    <!-- Sector Ratings Grid -->
    <div class="sectors-grid">
      {#each sectors as sector}
        <div class="sector-box">
          <div class="sector-header">
            <h4>{sector}: {marketRatings.sectorRatings[sector]}</h4>
            <button class="save-sector-btn" on:click={() => saveSectorRating(sector, null)}>
              Save {#if saveFeedback.sectors[sector]}<span class="feedback">{saveFeedback.sectors[sector]}</span>{/if}
            </button>
          </div>
          <input 
            type="range" 
            min="-3" 
            max="3" 
            step="1" 
            bind:value={marketRatings.sectorRatings[sector]}
          />
          <div class="scale-labels small">
            <span>Bearish (-3)</span>
            <span>Bullish (+3)</span>
          </div>
        </div>
      {/each}
    </div>
  </div>
  
  <!-- Individual Stock Rating Section (Bottom) -->
  <div class="stock-rating-container">
    <h3>Individual Stock Rating</h3>
    
    <div class="stock-rating-content">
      <div class="stock-rating-left">
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
      </div>
      
      <div class="stock-rating-right">
        <div class="input-group">
          <label>Primary Sector:</label>
          <select bind:value={rating.sector}>
            <option value="">Select primary sector...</option>
            {#each sectors as sector}
              <option value={sector}>{sector}</option>
            {/each}
          </select>
        </div>
        
        <div class="input-group">
          <label>Chart Patterns:</label>
          <div class="pattern-selection">
            <select bind:value={rating.selectedPattern}>
              <option value="">Select a pattern...</option>
              <optgroup label="Candlestick Patterns">
                <option value="Bullish Engulfing">Bullish Engulfing</option>
                <option value="Bearish Engulfing">Bearish Engulfing</option>
                <option value="Morning Star">Morning Star</option>
                <option value="Evening Star">Evening Star</option>
                <option value="Hammer">Hammer</option>
                <option value="Shooting Star">Shooting Star</option>
                <option value="Doji">Doji</option>
                <option value="Harami">Harami</option>
                <option value="Three White Soldiers">Three White Soldiers</option>
                <option value="Three Black Crows">Three Black Crows</option>
                <option value="Piercing Line">Piercing Line</option>
                <option value="Dark Cloud Cover">Dark Cloud Cover</option>
              </optgroup>
              <optgroup label="Chart Formations">
                <option value="Bull Flag">Bull Flag</option>
                <option value="Bear Flag">Bear Flag</option>
                <option value="Cup & Handle">Cup & Handle</option>
                <option value="Double Bottom">Double Bottom</option>
                <option value="Double Top">Double Top</option>
                <option value="Ascending Triangle">Ascending Triangle</option>
                <option value="Descending Triangle">Descending Triangle</option>
                <option value="Rounding Bottom">Rounding Bottom</option>
                <option value="Falling Wedge">Falling Wedge</option>
                <option value="Rising Wedge">Rising Wedge</option>
                <option value="Pennant">Pennant</option>
                <option value="Rectangle">Rectangle</option>
                <option value="Diamond">Diamond</option>
                <option value="Triple Top">Triple Top</option>
                <option value="Triple Bottom">Triple Bottom</option>
                <option value="Island Reversal">Island Reversal</option>
              </optgroup>
              <optgroup label="Moving Average Signals">
                <option value="Golden Cross">Golden Cross</option>
                <option value="Death Cross">Death Cross</option>
                <option value="20-day MA cross of 50-day MA">20-day MA cross of 50-day MA</option>
                <option value="Price above 200-day MA">Price above 200-day MA</option>
                <option value="Price below 200-day MA">Price below 200-day MA</option>
                <option value="Price above 50-day MA">Price above 50-day MA</option>
                <option value="Price below 50-day MA">Price below 50-day MA</option>
                <option value="Triple Moving Average Cross">Triple Moving Average Cross</option>
              </optgroup>
            </select>
            <button class="add-pattern-btn" on:click={addChartPattern}>Add Pattern</button>
            <button class="info-btn" on:click={() => togglePatternCheatSheet()}>Pattern Reference</button>
          </div>
          
          {#if rating.chartPatterns.length > 0}
            <div class="selected-patterns">
              <h5>Selected Patterns:</h5>
              <div class="pattern-tags">
                {#each rating.chartPatterns as pattern}
                  <div class="pattern-tag">
                    {pattern}
                    <button class="remove-tag" on:click={() => removeChartPattern(pattern)}>×</button>
                  </div>
                {/each}
              </div>
            </div>
          {/if}
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
      </div>
    </div>
    
    <div class="buttons-row">
      <button class="reset-btn" on:click={resetForm}>Reset</button>
      <button class="save-btn" on:click={saveStockRating}>Save Stock Rating</button>
    </div>
  </div>
  
  <!-- Recent Ratings Section (Optional, can be commented out if not needed) -->
  <div class="recent-ratings">
    <h3>Recent Stock Ratings</h3>
    <div class="ratings-table">
      <table>
        <thead>
          <tr>
            <th>Symbol</th>
            <th>Sector</th>
            <th>Sentiment</th>
            <th>Patterns</th>
            <th>Enthusiasm</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if recentRatings.length > 0}
            {#each recentRatings.filter(r => r.symbol && r.symbol !== 'SECTOR' && r.symbol !== 'MARKET') as item}
              <tr>
                <td>{item.symbol}</td>
                <td>{item.sector}</td>
                <td class={item.stockSentiment > 0 ? 'positive' : item.stockSentiment < 0 ? 'negative' : 'neutral'}>
                  {item.stockSentiment > 0 ? '+' : ''}{item.stockSentiment}
                </td>
                <td>
                  {#if item.chartPattern}
                    <!-- Handle data format with comma-separated patterns -->
                    <div class="pattern-list">
                      {#each item.chartPattern.split(', ').filter(p => p.trim() !== '') as pattern}
                        <span class="pattern-pill">{pattern}</span>
                      {/each}
                    </div>
                  {:else}
                    -
                  {/if}
                </td>
                <td>{item.enthusiasm !== undefined ? item.enthusiasm : 0}/10</td>
                <td class="action-buttons">
                  <button class="edit-btn" on:click={() => editRating(item)}>Edit</button>
                  <button class="delete-btn" on:click={() => {
                    console.log('Delete clicked for item:', item);
                    console.log('ID to delete:', item.id);
                    if (item.idd && !item.id) {
                      console.log('Found idd property instead of id:', item.idd);
                      deleteRating(item.idd);
                    } else {
                      deleteRating(item.id);
                    }
                  }}>Delete</button>
                </td>
              </tr>
            {/each}
          {:else}
            <tr>
              <td colspan="6">No ratings available for this date</td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>

<div class="chart-patterns-cheatsheet" class:visible={showPatternCheatSheet}>
  <div class="cheatsheet-header">
    <h3>Chart Pattern Reference</h3>
    <button class="close-btn" on:click={() => togglePatternCheatSheet()}>×</button>
  </div>
  
  <div class="cheatsheet-content">
    <div class="pattern-category">
      <h4>Candlestick Patterns</h4>
      <div class="pattern-grid">
        <div class="pattern-item">
          <h5>Bullish Engulfing</h5>
          <p>Reversal pattern where a bullish candle completely engulfs the previous bearish candle</p>
        </div>
        <div class="pattern-item">
          <h5>Bearish Engulfing</h5>
          <p>Reversal pattern where a bearish candle completely engulfs the previous bullish candle</p>
        </div>
        <div class="pattern-item">
          <h5>Morning Star</h5>
          <p>Bullish reversal pattern with three candles forming a bottom</p>
        </div>
        <div class="pattern-item">
          <h5>Evening Star</h5>
          <p>Bearish reversal pattern with three candles forming a top</p>
        </div>
        <div class="pattern-item">
          <h5>Hammer</h5>
          <p>Bullish reversal pattern with small body and long lower shadow</p>
        </div>
        <div class="pattern-item">
          <h5>Shooting Star</h5>
          <p>Bearish reversal pattern with small body and long upper shadow</p>
        </div>
        <div class="pattern-item">
          <h5>Doji</h5>
          <p>Indecision pattern with nearly equal open and close</p>
        </div>
        <div class="pattern-item">
          <h5>Harami</h5>
          <p>Reversal pattern where a small candle is contained within previous candle's body</p>
        </div>
        <div class="pattern-item">
          <h5>Three White Soldiers</h5>
          <p>Bullish continuation pattern with three consecutive rising candles</p>
        </div>
        <div class="pattern-item">
          <h5>Three Black Crows</h5>
          <p>Bearish continuation pattern with three consecutive falling candles</p>
        </div>
        <div class="pattern-item">
          <h5>Piercing Line</h5>
          <p>Bullish reversal pattern where a bullish candle opens below but closes above the midpoint of previous bearish candle</p>
        </div>
        <div class="pattern-item">
          <h5>Dark Cloud Cover</h5>
          <p>Bearish reversal pattern where a bearish candle opens above but closes below the midpoint of previous bullish candle</p>
        </div>
      </div>
    </div>
    
    <div class="pattern-category">
      <h4>Chart Formations</h4>
      <div class="pattern-grid">
        <div class="pattern-item">
          <h5>Bull Flag</h5>
          <p>Bullish continuation pattern with parallel downward trendlines following an uptrend</p>
        </div>
        <div class="pattern-item">
          <h5>Bear Flag</h5>
          <p>Bearish continuation pattern with parallel upward trendlines following a downtrend</p>
        </div>
        <div class="pattern-item">
          <h5>Cup & Handle</h5>
          <p>Bullish continuation pattern resembling a cup with a handle, showing consolidation then breakout</p>
        </div>
        <div class="pattern-item">
          <h5>Double Bottom</h5>
          <p>Bullish reversal pattern with two lows at approximately the same level</p>
        </div>
        <div class="pattern-item">
          <h5>Double Top</h5>
          <p>Bearish reversal pattern with two peaks at approximately the same level</p>
        </div>
        <div class="pattern-item">
          <h5>Ascending Triangle</h5>
          <p>Bullish pattern with horizontal resistance and rising support line</p>
        </div>
        <div class="pattern-item">
          <h5>Descending Triangle</h5>
          <p>Bearish pattern with horizontal support and falling resistance line</p>
        </div>
        <div class="pattern-item">
          <h5>Rounding Bottom</h5>
          <p>Long-term bullish reversal pattern shaped like a "U"</p>
        </div>
        <div class="pattern-item">
          <h5>Falling Wedge</h5>
          <p>Bullish pattern with converging downward trendlines</p>
        </div>
        <div class="pattern-item">
          <h5>Rising Wedge</h5>
          <p>Bearish pattern with converging upward trendlines</p>
        </div>
        <div class="pattern-item">
          <h5>Pennant</h5>
          <p>Short-term continuation pattern following a sharp move</p>
        </div>
        <div class="pattern-item">
          <h5>Rectangle</h5>
          <p>Continuation pattern with horizontal support and resistance</p>
        </div>
        <div class="pattern-item">
          <h5>Diamond</h5>
          <p>Rare reversal pattern formed by expanding then contracting price action</p>
        </div>
        <div class="pattern-item">
          <h5>Triple Top</h5>
          <p>Reversal pattern with three peaks at similar levels</p>
        </div>
        <div class="pattern-item">
          <h5>Triple Bottom</h5>
          <p>Reversal pattern with three troughs at similar levels</p>
        </div>
        <div class="pattern-item">
          <h5>Island Reversal</h5>
          <p>Pattern formed by gaps in both directions isolating price action</p>
        </div>
      </div>
    </div>
    
    <div class="pattern-category">
      <h4>Moving Average Signals</h4>
      <div class="pattern-grid">
        <div class="pattern-item">
          <h5>Golden Cross</h5>
          <p>Bullish signal when shorter-term MA crosses above longer-term MA (often 50-day crossing above 200-day)</p>
        </div>
        <div class="pattern-item">
          <h5>Death Cross</h5>
          <p>Bearish signal when shorter-term MA crosses below longer-term MA (often 50-day crossing below 200-day)</p>
        </div>
        <div class="pattern-item">
          <h5>20-day MA cross of 50-day MA</h5>
          <p>Shorter-term trend signal</p>
        </div>
        <div class="pattern-item">
          <h5>Price above/below 200-day MA</h5>
          <p>Major trend change signal</p>
        </div>
        <div class="pattern-item">
          <h5>Price above/below 50-day MA</h5>
          <p>Intermediate trend change signal</p>
        </div>
        <div class="pattern-item">
          <h5>Triple Moving Average Cross</h5>
          <p>Using three MAs (e.g., 5, 20, 50) for confirmation</p>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .dashboard {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
    background-color: var(--content-bg);
    border-radius: 5px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  h2, h3, h4, h5 {
    color: inherit;
    margin-top: 0;
    margin-bottom: 0.5rem;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 0;
  }
  
  h3 {
    margin-bottom: 1rem;
  }
  
  .divider {
    height: 2px;
    background-color: var(--border-color);
    margin: 1rem 0 2rem;
  }
  
  .description {
    text-align: center;
    color: inherit;
    margin-bottom: 2rem;
  }
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  /* Date navigation styles */
  .date-navigation {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    padding: 0.75rem;
    background-color: var(--card-bg);
    border-radius: 5px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
    position: relative;
  }
  
  .nav-button {
    background-color: var(--border-color);
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
    color: inherit;
    transition: all 0.2s;
  }
  
  .nav-button:hover:not(:disabled) {
    background-color: var(--active-button);
    color: white;
  }
  
  .nav-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    position: relative;
  }
  
  .nav-button:disabled:hover {
    cursor: not-allowed;
  }
  
  .current-date {
    font-weight: bold;
    font-size: 1.1rem;
  }
  
  /* Help message styles */
  .help-message {
    margin: 1rem 0;
    padding: 0.75rem;
    background-color: var(--bg-color);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    text-align: center;
    font-style: italic;
  }
  
  .debug-btn {
    background-color: var(--active-button);
    color: white;
    border: none;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    margin-top: 0.5rem;
    cursor: pointer;
  }
  
  /* Market & Sectors container styles */
  .market-sectors-container {
    background-color: var(--card-bg);
    padding: 1.5rem;
    border-radius: 5px;
    margin-bottom: 2rem;
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
    border: 1px solid var(--input-border);
    border-radius: 4px;
    background-color: var(--input-bg);
    color: var(--text-color);
  }
  
  .save-all-btn {
    background-color: var(--active-button);
    color: white;
    border: none;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    font-weight: bold;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .save-all-btn:hover {
    background-color: var(--active-button);
    opacity: 0.9;
  }
  
  /* Overall market rating box */
  .market-rating-box {
    background-color: var(--bg-color);
    padding: 1rem;
    border-radius: 5px;
    margin-bottom: 1.5rem;
    border: 1px solid var(--border-color);
  }
  
  .market-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }
  
  /* Sectors grid layout */
  .sectors-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
  }
  
  .sector-box {
    background-color: var(--bg-color);
    padding: 1rem;
    border-radius: 5px;
    border: 1px solid var(--border-color);
  }
  
  .sector-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }
  
  .save-sector-btn {
    background-color: var(--border-color);
    color: inherit;
    border: none;
    border-radius: 4px;
    padding: 0.25rem 0.5rem;
    font-size: 0.9rem;
    cursor: pointer;
    transition: background-color 0.2s;
    display: flex;
    align-items: center;
  }
  
  .save-sector-btn:hover {
    background-color: var(--active-button);
    color: white;
  }
  
  .feedback {
    margin-left: 0.5rem;
    font-style: italic;
    font-size: 0.8rem;
  }
  
  /* Stock Rating container */
  .stock-rating-container {
    background-color: var(--card-bg);
    padding: 1.5rem;
    border-radius: 5px;
    margin-bottom: 2rem;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  }
  
  .stock-rating-content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
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
    border: 1px solid var(--input-border);
    border-radius: 4px;
    background-color: var(--input-bg);
    color: var(--text-color);
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
    color: inherit;
    font-size: 0.85rem;
    opacity: 0.8;
  }
  
  .scale-labels.small {
    font-size: 0.75rem;
  }
  
  .buttons-row {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
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
    background-color: var(--border-color);
    color: inherit;
  }
  
  .save-btn {
    background-color: var(--active-button);
    color: white;
  }
  
  .reset-btn:hover {
    background-color: var(--border-color);
    opacity: 0.8;
  }
  
  .save-btn:hover {
    background-color: var(--active-button);
    opacity: 0.9;
  }
  
  /* Recent ratings section */
  .recent-ratings {
    background-color: var(--card-bg);
    padding: 1.5rem;
    border-radius: 5px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
    transition: background-color 0.3s ease, color 0.3s ease;
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
  
  .info-btn {
    margin-top: 0.5rem;
    padding: 0.5rem 1rem;
    background-color: var(--border-color);
    border: none;
    border-radius: 4px;
    color: inherit;
    font-size: 0.9rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .info-btn:hover {
    background-color: var(--active-button);
    color: white;
  }
  
  .chart-patterns-cheatsheet {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 90%;
    max-width: 1000px;
    max-height: 80vh;
    background-color: var(--card-bg);
    border-radius: 8px;
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
    z-index: 1000;
    display: none;
    overflow: hidden;
    transition: all 0.3s ease;
  }
  
  .chart-patterns-cheatsheet.visible {
    display: block;
  }
  
  .cheatsheet-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--border-color);
  }
  
  .cheatsheet-header h3 {
    margin: 0;
  }
  
  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: inherit;
  }
  
  .cheatsheet-content {
    padding: 1.5rem;
    overflow-y: auto;
    max-height: calc(80vh - 60px);
  }
  
  .pattern-category {
    margin-bottom: 2rem;
  }
  
  .pattern-category h4 {
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 0.5rem;
    margin-bottom: 1rem;
  }
  
  .pattern-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
  }
  
  .pattern-item {
    background-color: var(--bg-color);
    padding: 1rem;
    border-radius: 4px;
    border: 1px solid var(--border-color);
  }
  
  .pattern-item h5 {
    margin-top: 0;
    margin-bottom: 0.5rem;
    color: var(--active-button);
  }
  
  .pattern-item p {
    margin: 0;
    font-size: 0.9rem;
    line-height: 1.4;
  }
  
  /* Responsive fixes */
  @media (max-width: 768px) {
    .pattern-grid {
      grid-template-columns: 1fr;
    }
    
    .chart-patterns-cheatsheet {
      width: 95%;
      max-height: 90vh;
    }
  }
  
  .pattern-selection {
    display: flex;
    gap: 0.5rem;
    align-items: flex-start;
    flex-wrap: wrap;
    margin-bottom: 0.5rem;
  }
  
  .pattern-selection select {
    flex: 1;
    min-width: 200px;
  }
  
  .add-pattern-btn {
    background-color: var(--active-button);
    color: white;
    border: none;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .add-pattern-btn:hover {
    opacity: 0.9;
  }
  
  .selected-patterns {
    margin-top: 0.75rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    padding: 0.75rem;
    background-color: var(--bg-color);
  }
  
  .selected-patterns h5 {
    margin-top: 0;
    margin-bottom: 0.5rem;
    font-size: 0.9rem;
  }
  
  .pattern-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }
  
  .pattern-tag {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    background-color: var(--active-button);
    color: white;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.8rem;
  }
  
  .remove-tag {
    background: none;
    border: none;
    color: white;
    cursor: pointer;
    padding: 0;
    font-size: 1rem;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-left: 0.25rem;
  }
  
  .pattern-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
  }
  
  .pattern-pill {
    font-size: 0.8rem;
    background-color: var(--active-button);
    color: white;
    padding: 0.15rem 0.4rem;
    border-radius: 4px;
    white-space: nowrap;
  }
  
  /* New action buttons */
  .action-buttons {
    display: flex;
    gap: 0.5rem;
  }
  
  .edit-btn, .delete-btn {
    padding: 0.25rem 0.5rem;
    border: none;
    border-radius: 4px;
    font-size: 0.8rem;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .edit-btn {
    background-color: var(--active-button);
    color: white;
  }
  
  .delete-btn {
    background-color: #e53e3e;
    color: white;
  }
  
  .edit-btn:hover, .delete-btn:hover {
    opacity: 0.9;
  }
</style> 