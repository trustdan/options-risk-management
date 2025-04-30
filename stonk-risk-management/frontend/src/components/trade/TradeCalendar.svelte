<script>
  import { onMount, tick } from 'svelte';
  // Import Wails backend functions
  import {
    GetTrades,
    SaveTrade,
    DeleteTrade,
    // Assume these exist or will be created in Go
    GetLatestMarketRating, 
    GetLatestSectorRating,
    GetLatestStockRating,
    SaveStockRating
  } from '../../../wailsjs/go/main/App';
  import { models } from '../../../wailsjs/go/models';
  
  // Import the new TradeJournal component
  import TradeJournal from './TradeJournal.svelte';
  
  // Active tab for navigation
  let activeTab = 'calendar'; // 'calendar', 'journal'
  let selectedTradeForJournal = null;
  
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
    'ETF',
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
    { name: 'Ratio Spreads', color: '#d35400' },
    { name: 'Danger - naked options ahead', color: '#ff0000' }
  ];
  
  // Strategy types with detailed options
  const strategyOptions = [
    {
      category: 'Basic Spreads',
      options: [
        { name: 'Long Call', description: 'Bullish directional bet - Buy a call option, profit from upward price movement' },
        { name: 'Long Put', description: 'Bearish directional bet - Buy a put option, profit from downward price movement' },
        { name: 'Covered Call', description: 'Income strategy - Own stock and sell calls against it for premium income' }
      ]
    },
    {
      category: 'Vertical Spreads',
      options: [
        { name: 'Bull Call Spread', description: 'Bullish with limited risk/reward - Buy lower strike call, sell higher strike call, same expiration (Call Debit Spread)' },
        { name: 'Bear Call Spread', description: 'Bearish with limited risk/reward - Sell lower strike call, buy higher strike call, same expiration (Call Credit Spread)' },
        { name: 'Bull Put Spread', description: 'Bullish with limited risk/reward - Sell higher strike put, buy lower strike put, same expiration (Put Credit Spread)' },
        { name: 'Bear Put Spread', description: 'Bearish with limited risk/reward - Buy higher strike put, sell lower strike put, same expiration (Put Debit Spread)' }
      ]
    },
    {
      category: 'Calendar/Horizontal Spreads',
      options: [
        { name: 'Long Calendar Call Spread', description: 'Neutral to slightly bullish - Sell near-term call, buy longer-term call, same strike' },
        { name: 'Long Calendar Put Spread', description: 'Neutral to slightly bearish - Sell near-term put, buy longer-term put, same strike' }
      ]
    },
    {
      category: 'Diagonal Spreads',
      options: [
        { name: 'Diagonal Call Spread Up', description: 'Moderately bullish - Buy longer-term lower strike call, sell shorter-term higher strike call' },
        { name: 'Diagonal Call Spread Down', description: 'Bearish to neutral - Buy longer-term higher strike call, sell shorter-term lower strike call' },
        { name: 'Diagonal Put Spread Up', description: 'Bearish to neutral - Buy longer-term lower strike put, sell shorter-term higher strike put' },
        { name: 'Diagonal Put Spread Down', description: 'Moderately bullish - Buy longer-term higher strike put, sell shorter-term lower strike put' }
      ]
    },
    {
      category: 'Butterfly Spreads',
      options: [
        { name: 'Long Call Butterfly', description: 'Neutral, expecting little movement - Buy 1 lower strike call, sell 2 middle strike calls, buy 1 higher strike call' },
        { name: 'Long Put Butterfly', description: 'Neutral, expecting little movement - Buy 1 lower strike put, sell 2 middle strike puts, buy 1 higher strike put' },
        { name: 'Broken Wing Butterfly Up', description: 'Bullish, with uneven wings - Like standard butterfly but with wider spread between middle and upper strikes' },
        { name: 'Broken Wing Butterfly Down', description: 'Bearish, with uneven wings - Like standard butterfly but with wider spread between lower and middle strikes' }
      ]
    },
    {
      category: 'Iron Condors/Butterflies',
      options: [
        { name: 'Iron Condor', description: 'Neutral, expecting range-bound movement - Sell OTM put, buy further OTM put, sell OTM call, buy further OTM call' },
        { name: 'Iron Butterfly', description: 'Highly neutral, expecting minimal movement - Buy OTM put, sell ATM put, sell ATM call, buy OTM call' }
      ]
    },
    {
      category: 'Ratio Spreads',
      options: [
        { name: 'Call Ratio Backspread', description: 'Bullish with volatility bias - Buy more calls at higher strike than selling at lower strike (e.g., 2:1 ratio)' },
        { name: 'Put Ratio Backspread', description: 'Bearish with volatility bias - Buy more puts at lower strike than selling at higher strike (e.g., 2:1 ratio)' }
      ]
    },
    {
      category: 'Danger - naked options ahead',
      options: [
        { name: 'Short Call', description: 'Bearish directional bet - Sell a call option, profit from downward or sideways price movement' },
        { name: 'Short Put', description: 'Bullish directional bet - Sell a put option, profit from upward or sideways price movement' },
        { name: 'Cash-Secured Put', description: 'Income strategy - Sell puts with cash to secure the potential stock purchase' },
        { name: 'Call Ratio Spread', description: 'Neutral to moderately bearish - Sell more calls at higher strike than buying at lower strike (e.g., 1:2 ratio)' },
        { name: 'Put Ratio Spread', description: 'Neutral to moderately bullish - Sell more puts at lower strike than buying at higher strike (e.g., 1:2 ratio)' }
      ]
    }
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
      expLabel: `Exp: ${expDate.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: '2-digit' })}`,
      range: `${weekStart.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - ${expDate.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}`
    });
    
    // Move to next week
    currentDate.setDate(currentDate.getDate() + 7);
  }
  
  // Trade data - will be loaded from backend
  let trades = []; // Removed explicit TS type
  let tradesLoaded = false; // Track if trades have been initially loaded
  let componentKey = 0; // Used to force full component re-rendering
  
  // Track whether we're entering a calendar or diagonal strategy
  $: {
    const strategyParts = newTrade.strategy.split(' - ');
    if (strategyParts.length >= 2) {
      const category = strategyParts[0];
      const type = strategyParts[1];
      
      showShortLegField = 
        category.includes('Calendar') || 
        category.includes('Diagonal') ||
        type.includes('Calendar') ||
        type.includes('Diagonal');
    } else {
      showShortLegField = false;
    }
  }
  
  // Update selected strategy description when strategy changes
  $: {
    if (newTrade.strategy) {
      const [category, type] = newTrade.strategy.split(' - ');
      selectedStrategyCategory = category;
      
      const categoryObj = strategyOptions.find(s => s.category === category);
      if (categoryObj) {
        const strategyObj = categoryObj.options.find(o => o.name === type);
        if (strategyObj) {
          selectedStrategyDescription = strategyObj.description;
        } else {
          selectedStrategyDescription = '';
        }
      } else {
        selectedStrategyDescription = '';
      }
    } else {
      selectedStrategyDescription = '';
    }
  }
  
  // Data for the new trade form
  let newTrade = {
    id: '', // Will be set when editing or generated on save
    entryDate: new Date().toISOString().split('T')[0],
    expirationDate: '',
    ticker: '',
    sector: '',
    entryPrice: 0,
    timeframe: '',
    entry: 0,
    stop: 0,
    target: 0,
    strategy: '',
    notes: '',
    shortLegExpiration: '' // New field for short leg expiration info
  };
  
  // For calendar and diagonal spreads
  let showShortLegField = false;
  let selectedStrategyCategory = '';
  let selectedStrategyDescription = '';
  
  // For editing trades
  let isEditing = false;
  let editingTradeId = null;
  
  // For search/filter
  let searchQuery = '';
  let searchType = 'ticker';
  
  // State for displaying ratings in the form
  let marketRating = null;
  let sectorRating = null;
  let stockRating = null;
  let ratingsLoading = { market: false, sector: false, stock: false };
  
  // For rating editing functionality
  let showRatingEditor = false;
  let editingRatingType = ''; // 'sector' or 'stock'
  let editingRating = null;
  
  // Filtered trades for history display
  $: filteredTrades = trades
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
  
  function getStrategyColor(strategyName) {
    const strategy = strategies.find(s => s.name === strategyName);
    return strategy ? strategy.color : '#999';
  }
  
  // Function to determine if a trade spans multiple weeks by checking for other legs with the same ID
  function isMultiWeekTrade(tradeId) {
    if (!trades || trades.length === 0) return false;
    const legsWithSameId = trades.filter(t => t.id === tradeId);
    return legsWithSameId.length > 1;
  }
  
  // Get all weeks that a trade spans
  function getTradeWeeks(tradeId) {
    if (!trades || trades.length === 0) return [];
    const legsWithSameId = trades.filter(t => t.id === tradeId);
    return legsWithSameId.map(leg => leg.week).sort((a, b) => a - b);
  }
  
  // Check if trade has entry in given week
  function hasTradeInWeek(tradeId, weekNumber) {
    return trades.some(t => t.id === tradeId && t.week === weekNumber);
  }
  
  // Calculate which weeks a trade spans based on entry and expiration dates
  function getTradeSpanningWeeks(trade) {
    if (!trade || !trade.entryDate || !trade.expirationDate) {
        console.warn(`Missing date data for trade:`, trade);
        // Fallback to original week if dates are missing or invalid
        return trade && typeof trade.week === 'number' ? [trade.week] : [];
    }

    let entryDate, expirationDate;
    try {
        // Ensure we're working with Date objects
        entryDate = trade.entryDate instanceof Date ? trade.entryDate : new Date(trade.entryDate);
        expirationDate = trade.expirationDate instanceof Date ? trade.expirationDate : new Date(trade.expirationDate);
        
        // Check if dates are valid
        if (isNaN(entryDate.getTime()) || isNaN(expirationDate.getTime())) {
            console.warn(`Invalid dates for trade ${trade.id}:`, {entry: entryDate, expiration: expirationDate});
            return typeof trade.week === 'number' ? [trade.week] : [];
        }
    } catch (e) {
        console.error("Error parsing dates for trade:", trade.id, e);
        return typeof trade.week === 'number' ? [trade.week] : [];
    }

    // Special debug for key trades
    if (trade.symbol === 'MSFT' || trade.symbol === 'META' || trade.symbol === 'DHR') {
      console.log(`${trade.symbol} date range: ${entryDate.toISOString()} to ${expirationDate.toISOString()}`);
    }

    // Find weeks that this trade spans (entry date to expiration date)
    const spannedWeeks = weeks
      .filter(week => {
        const weekStart = new Date(week.startDate);
        const weekEnd = new Date(week.expirationDate);
        // Make weekEnd exclusive for comparison by adding a day
        weekEnd.setDate(weekEnd.getDate() + 1);

        // Check for overlap: week starts before trade ends AND week ends after trade starts
        const overlaps = weekStart < expirationDate && weekEnd > entryDate;
        
        // Debug overlapping for specific trades
        if ((trade.symbol === 'MSFT' || trade.symbol === 'META' || trade.symbol === 'DHR') && overlaps) {
          console.log(`${trade.symbol} overlaps with week ${week.number}: ${week.range}`);
        }
        
        return overlaps;
      })
      .map(week => week.number);
      
    // If we don't find any spans, fall back to the trade's defined week as a safety
    if (spannedWeeks.length === 0 && typeof trade.week === 'number') {
      console.warn(`${trade.symbol} didn't span any weeks by date, using its week property: ${trade.week}`);
      return [trade.week];
    }
    
    return spannedWeeks;
  }
  
  // Check if a trade spans into a specific week
  function tradeSpansIntoWeek(trade, weekNumber) {
    if (!trade) return false;
    
    // Debug log to check trade and parameters
    if (trade.symbol === 'MSFT' || trade.symbol === 'META' || trade.symbol === 'DHR') {
      console.log(`Checking if ${trade.symbol} spans into week ${weekNumber}`);
    }
    
    const spannedWeeks = getTradeSpanningWeeks(trade);
    
    // More detailed logging for key trades
    if (trade.symbol === 'MSFT' || trade.symbol === 'META' || trade.symbol === 'DHR') {
      console.log(`${trade.symbol} spans these weeks:`, spannedWeeks);
    }
    
    return spannedWeeks.includes(weekNumber);
  }
  
  // Determine if this is the entry week, expiration week, or middle week for display styling
  function getTradePositionInSpan(trade, weekNumber) {
    const spannedWeeks = getTradeSpanningWeeks(trade);
    if (spannedWeeks.length <= 1) return 'single'; // Not a spanning trade

    const minWeek = Math.min(...spannedWeeks);
    const maxWeek = Math.max(...spannedWeeks);

    if (weekNumber === minWeek) {
      return 'start';
    } else if (weekNumber === maxWeek) {
      return 'end';
    } else if (weekNumber > minWeek && weekNumber < maxWeek) {
      return 'middle';
    }

    // Fallback for safety, though should be covered by tradeSpansIntoWeek
    return 'single';
  }

  function getTradesForCell(sector, weekNumber) {
    if (!trades || trades.length === 0) return [];

    // Debug log for specific sectors and weeks
    if (sector === 'Technology' || sector === 'Communication' || sector === 'Industrial') {
      console.log(`getTradesForCell(${sector}, ${weekNumber}) - processing ${trades.length} trades`);
    }

    // Find all trades for this sector
    const sectorTrades = trades.filter(trade => trade.sector === sector);
    
    if (sector === 'Technology' || sector === 'Communication' || sector === 'Industrial') {
      console.log(`${sectorTrades.length} trades found for sector ${sector}`);
    }

    // Filter trades that actually span into the current weekNumber
    const relevantTrades = sectorTrades.filter(trade =>
      tradeSpansIntoWeek(trade, weekNumber)
    );
    
    if (sector === 'Technology' || sector === 'Communication' || sector === 'Industrial') {
      console.log(`${relevantTrades.length} relevant trades for ${sector} in week ${weekNumber}`);
      if (relevantTrades.length > 0) {
        console.log("Trade symbols:", relevantTrades.map(t => t.symbol).join(", "));
      }
    }

    // Map these trades to the display format needed for the cell
    const displayTrades = relevantTrades.map(trade => {
      const spannedWeeks = getTradeSpanningWeeks(trade);
      const isMultiWeek = spannedWeeks.length > 1; // Does it actually span more than one week?
      const position = getTradePositionInSpan(trade, weekNumber);

      return {
        ...trade,
        // Mark as spanning if it covers multiple weeks *and* we are rendering it in one of its spanned positions
        isSpanningLeg: isMultiWeek,
        // Determine position within the span for styling
        spanPosition: isMultiWeek ? position : 'single',
        // Keep original trade.week, but ensure we have the correct weekNumber for context if needed elsewhere
        currentDisplayWeek: weekNumber
      };
    });

    return displayTrades;
  }

  function hasTrades(sector, weekNumber) {
    if (!trades || trades.length === 0) return false;

    // Perform more direct check after initial load 
    if (tradesLoaded) {
      // Log for debugging key sectors
      if (sector === 'Technology' || sector === 'Communication' || sector === 'Industrial') {
        console.log(`hasTrades checking ${sector} week ${weekNumber}`);
      }
      
      // Check if *any* trade for this sector should be displayed in the current weekNumber
      return trades.some(trade =>
        trade.sector === sector && tradeSpansIntoWeek(trade, weekNumber)
      );
    }
    
    return false;
  }
  
  // Simplified Function to load trades from backend
  async function loadTrades() {
    try {
      console.log('Loading and processing trades...');
      const backendTrades = await GetTrades();
      console.log('Raw trades received:', backendTrades ? backendTrades.length : 0);
      
      // Create a new array reference to trigger reactivity
      trades = processTrades(backendTrades);
      tradesLoaded = true; // Mark trades as loaded
      
      // Force update all derived values by creating another new reference
      setTimeout(() => {
        console.log('Forcing UI update by creating a new reference...');
        trades = [...trades];
      }, 50);
      
      console.log('Processed trades assigned:', trades.length);
    } catch (error) {
      console.error('Failed to load trades:', error);
      alert(`Error loading trades: ${error.message || error}`);
      trades = [];
      tradesLoaded = false;
    }
  }

  // Function to force a complete refresh of the component
  async function forceRefresh() {
    console.log('Forcing complete component refresh');
    tradesLoaded = false; // Mark as not loaded to prevent flicker with old data
    await loadTrades(); // Reload the data from backend
    await tick(); // Wait for DOM update
    componentKey += 1; // Increment key to force Svelte to recreate the component parts
    console.log('Component refresh complete, new key:', componentKey);
  }

  async function addNewTrade() {
    // Validate required fields
    if (!newTrade.ticker || !newTrade.sector || !newTrade.strategy || !newTrade.expirationDate) {
      alert('Please fill in all required fields: Ticker, Sector, Strategy, and Primary Expiration Date');
      return;
    }
    
    const expDate = new Date(newTrade.expirationDate);
    const entryDate = new Date(newTrade.entryDate);

    // Find the week for the primary expiration date
    const weekIndex = weeks.findIndex(w => {
      const weekStart = new Date(w.startDate);
      const weekEnd = new Date(w.expirationDate);
      weekEnd.setDate(weekEnd.getDate() + 1); // Include Friday
      return expDate >= weekStart && expDate < weekEnd;
    });

    if (weekIndex === -1) {
      alert('Primary expiration date is outside the displayed weeks.');
      return;
    }
    
    const [strategyCategory, strategyType] = newTrade.strategy.split(' - ');
    const tradeId = isEditing ? editingTradeId : Date.now().toString(); 

    if (isEditing) {
      try {
        await DeleteTrade(editingTradeId);
      } catch (error) {
        console.error('Failed to delete old trade during edit:', error);
        alert('Error updating trade. Could not remove old data.');
        return;
      }
    }
    
    // Create a single trade object with all the data
    const tradeData = {
      id: tradeId,
      symbol: newTrade.ticker,
      sector: newTrade.sector,
      strategy: strategyCategory,
      type: strategyType,
      week: weekIndex + 1,
      entryPrice: parseFloat(newTrade.entryPrice.toString()),
      // Add the new fields
      timeframe: newTrade.timeframe,
      entry: parseFloat(newTrade.entry.toString()),
      stop: parseFloat(newTrade.stop.toString()),
      target: parseFloat(newTrade.target.toString()),
      notes: newTrade.notes,
      entryDate: entryDate,
      expirationDate: expDate,
      isMultiLeg: showShortLegField,
      legNumber: 1,
      shortLegExp: newTrade.shortLegExpiration || ''
    };

    // Save to backend
    try {
      const tradeModel = new models.Trade();
      Object.assign(tradeModel, tradeData);
      await SaveTrade(tradeModel);
      
      console.log('Save complete, forcing full refresh...');
      resetForm();
      
      // Force a complete refresh to show the new trade
      await forceRefresh();
      
      alert(`Trade ${isEditing ? 'updated' : 'added'} successfully!`);
    } catch (error) {
      console.error('Failed to save trade:', error);
      alert(`Error saving trade: ${error.message || error}`);
      if (isEditing) {
         console.log('Error during edit, reloading trades...');
         await loadTrades();
      }
    }
  }
  
  function resetForm() {
    newTrade = {
      id: '',
      entryDate: new Date().toISOString().split('T')[0],
      expirationDate: '',
      ticker: '',
      sector: '',
      entryPrice: 0,
      timeframe: '',
      entry: 0,
      stop: 0,
      target: 0,
      strategy: '',
      notes: '',
      shortLegExpiration: ''
    };
    showShortLegField = false;
    selectedStrategyCategory = '';
    selectedStrategyDescription = '';
    isEditing = false;
    editingTradeId = null;
  }
  
  function editTrade(tradeId) {
    const trade = trades.find(t => t.id === tradeId);
    if (!trade) return;
    
    // Ensure dates are proper JS Date objects for formatting
    let entryDate = null;
    let expDate = null;
    
    // Handle potential date format issues safely
    if (trade.entryDate) {
      try {
        entryDate = new Date(trade.entryDate);
      } catch (e) {
        console.error("Failed to parse entry date:", trade.entryDate);
      }
    }
    
    if (trade.expirationDate) {
      try {
        expDate = new Date(trade.expirationDate);
      } catch (e) {
        console.error("Failed to parse expiration date:", trade.expirationDate);
      }
    }

    newTrade = {
      id: trade.id,
      entryDate: entryDate ? entryDate.toISOString().split('T')[0] : '',
      expirationDate: expDate ? expDate.toISOString().split('T')[0] : '',
      ticker: trade.symbol,
      sector: trade.sector,
      entryPrice: trade.entryPrice,
      // For the new fields, use the stored values or defaults
      timeframe: trade.timeframe || '',
      entry: trade.entry || trade.entryPrice || 0,
      stop: trade.stop || 0,
      target: trade.target || 0,
      strategy: `${trade.strategy} - ${trade.type}`,
      notes: trade.notes,
      shortLegExpiration: trade.shortLegExp || ''
    };
    
    // Check if this is a calendar or diagonal spread
    showShortLegField = trade.isMultiLeg || 
                        trade.strategy.includes('Calendar') || 
                        trade.strategy.includes('Diagonal');
    
    isEditing = true;
    editingTradeId = tradeId;
    
    document.querySelector('.add-trade-section').scrollIntoView({ behavior: 'smooth' });
  }
  
  async function deleteTrade(tradeId) {
    if (confirm('Are you sure you want to delete this trade (all legs)?')) {
      try {
        await DeleteTrade(tradeId);
        console.log('Delete complete, reloading trades...');
        await loadTrades();
      } catch (error) {
        console.error('Failed to delete trade:', error);
        alert(`Error deleting trade: ${error.message || error}`);
      }
    }
  }
  
  function handleTradeCardClick(trade) {
    editTrade(trade.id);
  }
  
  function formatDate(date) {
    if (!date) return '';
    // Handle potential Go time object or string
    let dateObj = null;
    
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
  
  // Function to process trades - ensure dates are Date objects
  function processTrades(rawTrades) {
    console.log('Processing raw trades...', rawTrades ? rawTrades.length : 0);
    if (!rawTrades || !Array.isArray(rawTrades)) return [];

    const processed = rawTrades.map(trade => {
      const tradeCopy = { ...trade };
      if (tradeCopy.entryDate) {
        try {
          // @ts-ignore
          const entryDate = new Date(tradeCopy.entryDate);
          tradeCopy.entryDate = !isNaN(entryDate.getTime()) ? entryDate : null;
          console.log(`Processed entry date for ${tradeCopy.symbol}: ${entryDate}`);
        } catch (e) { 
          console.error(`Error parsing entry date for ${tradeCopy.symbol}:`, e);
          tradeCopy.entryDate = null; 
        }
      }
      if (tradeCopy.expirationDate) {
        try {
          // @ts-ignore
          const expirationDate = new Date(tradeCopy.expirationDate);
          tradeCopy.expirationDate = !isNaN(expirationDate.getTime()) ? expirationDate : null;
          console.log(`Processed expiration date for ${tradeCopy.symbol}: ${expirationDate}`);
        } catch (e) { 
          console.error(`Error parsing expiration date for ${tradeCopy.symbol}:`, e);
          tradeCopy.expirationDate = null; 
        }
      }
      return tradeCopy;
    });
    
    const filtered = processed.filter(trade => trade.entryDate && trade.expirationDate);
    console.log(`After filtering: ${filtered.length} valid trades remain`);
    
    // Log first trade spans for debugging
    if (filtered.length > 0) {
      const firstTrade = filtered[0];
      const spans = getTradeSpanningWeeks(firstTrade);
      console.log(`First trade ${firstTrade.symbol} spans weeks:`, spans);
    }
    
    return filtered;
  }

  // Simple reactive logger for the trades array
  $: console.log('Trades array updated, length:', trades.length);

  // Additional reactive debugger for rendered cells
  $: {
    if (trades.length > 0) {
      const sampleSector = trades[0].sector;
      const weekNumber = trades[0].week;
      const cellTrades = getTradesForCell(sampleSector, weekNumber);
      console.log(`getTradesForCell(${sampleSector}, ${weekNumber}) returns ${cellTrades.length} trades`);
    }
  }

  // --- Rating Fetching ---
  async function fetchLatestMarketRating() {
    ratingsLoading.market = true;
    marketRating = null; // Reset
    try {
      marketRating = await GetLatestMarketRating();
      console.log("Market rating fetched:", marketRating);
      calculateSyntheticRating(); // Recalculate synthetic rating
    } catch (error) {
      console.error('Failed to fetch market rating:', error);
      marketRating = { error: 'Failed to load' }; // Indicate error
    } finally {
      ratingsLoading.market = false;
    }
  }

  async function handleSectorChange(event) {
    const selectedSector = event.target.value;
    sectorRating = null; // Reset
    stockRating = null; // Also reset stock rating if sector changes
    if (!selectedSector) {
      calculateSyntheticRating(); // Recalculate even when clearing
      return; // No sector selected
    }

    ratingsLoading.sector = true;
    try {
      sectorRating = await GetLatestSectorRating(selectedSector);
      console.log("Sector rating fetched:", sectorRating);
      // Auto-set sector in the form if it was changed via dropdown
      if (newTrade.sector !== selectedSector) {
          newTrade.sector = selectedSector;
      }
      calculateSyntheticRating(); // Recalculate synthetic rating
    } catch (error) {
      console.error('Failed to fetch sector rating:', error);
      sectorRating = { error: 'Failed to load' };
    } finally {
      ratingsLoading.sector = false;
    }
  }

  // Add a local error state variable instead of trying to add properties to stockRating
  let stockRatingError = null;

  async function handleTickerBlur(event) {
    const currentTicker = event.target.value.toUpperCase().trim(); // Use uppercase ticker and trim whitespace
    if (!currentTicker) {
      stockRating = null; // Reset
      stockRatingError = null; // Reset error state too
      calculateSyntheticRating(); // Recalculate even when clearing
      return; // No ticker entered
    }
    
    // Don't reset stockRating right away to prevent flickering
    ratingsLoading.stock = true;
    stockRatingError = null; // Reset error state
    
    try {
      // Attempt to fetch the stock rating
      const rating = await GetLatestStockRating(currentTicker);
      console.log("Stock rating fetched:", rating);
      
      // Check if we got a valid rating
      if (rating && typeof rating.stockSentiment === 'number') {
        // Valid rating received
        stockRating = rating;
        
        // If a rating is found and sector is empty or doesn't match, update form sector
        if (rating.sector && (!newTrade.sector || newTrade.sector !== rating.sector)) {
          newTrade.sector = rating.sector;
          // Fetch sector rating as well if we auto-filled it
          await handleSectorChange({ target: { value: rating.sector } });
        } else {
          calculateSyntheticRating(); // Recalculate synthetic rating
        }
      } else {
        // No data or invalid data
        stockRating = null;
        stockRatingError = {
          message: 'No rating found',
          ticker: currentTicker
        };
        calculateSyntheticRating();
      }
    } catch (error) {
      console.error('Failed to fetch stock rating:', error);
      stockRating = null;
      stockRatingError = {
        message: 'Failed to load',
        ticker: currentTicker
      };
      calculateSyntheticRating();
    } finally {
      ratingsLoading.stock = false;
    }
  }

  // Add a refresh function for stock ratings
  async function refreshStockRating() {
    if (!newTrade.ticker) return;
    
    const ticker = newTrade.ticker.toUpperCase().trim();
    ratingsLoading.stock = true;
    stockRating = null; // Clear the current rating
    stockRatingError = null; // Reset error state
    
    try {
      // Attempt to fetch the stock rating
      const rating = await GetLatestStockRating(ticker);
      console.log("Stock rating refreshed:", rating);
      
      if (rating && typeof rating.stockSentiment === 'number') {
        stockRating = rating;
        calculateSyntheticRating();
      } else {
        stockRatingError = {
          message: 'No rating found after refresh',
          ticker: ticker
        };
        calculateSyntheticRating();
      }
    } catch (error) {
      console.error('Failed to refresh stock rating:', error);
      stockRatingError = {
        message: 'Failed to refresh',
        ticker: ticker
      };
      calculateSyntheticRating();
    } finally {
      ratingsLoading.stock = false;
    }
  }

  // Synthetic rating calculation
  let syntheticRating = null;
  let recommendedStrategies = null;
  
  function calculateSyntheticRating() {
    // Reset
    syntheticRating = null;
    recommendedStrategies = null;
    
    // Check if we have at least one valid rating
    const hasMarket = marketRating && !marketRating.error && typeof marketRating.stockSentiment === 'number';
    const hasSector = sectorRating && !sectorRating.error && typeof sectorRating.stockSentiment === 'number';
    const hasStock = stockRating && !stockRating.error && typeof stockRating.stockSentiment === 'number';
    
    if (!hasMarket && !hasSector && !hasStock) {
      return; // No data for calculation
    }
    
    // Weights for each component
    const marketWeight = 0.5;
    const sectorWeight = 0.3;
    const stockWeight = 0.2;
    
    // Total weights (adjust if some ratings are missing)
    let totalWeight = 0;
    let weightedSum = 0;
    
    // Add each available rating with its weight
    if (hasMarket) {
      weightedSum += marketRating.stockSentiment * marketWeight;
      totalWeight += marketWeight;
    }
    
    if (hasSector) {
      weightedSum += sectorRating.stockSentiment * sectorWeight;
      totalWeight += sectorWeight;
    }
    
    if (hasStock) {
      weightedSum += stockRating.stockSentiment * stockWeight;
      totalWeight += stockWeight;
    }
    
    // Normalize to account for missing ratings
    if (totalWeight > 0) {
      // Calculate the final weighted value
      const weightedValue = weightedSum / totalWeight;
      
      // Set synthetic rating with additional info
      syntheticRating = {
        value: weightedValue,
        numericValue: Math.round(weightedValue * 10) / 10, // Round to 1 decimal place
        components: {
          market: hasMarket ? marketRating.stockSentiment : null,
          sector: hasSector ? sectorRating.stockSentiment : null,
          stock: hasStock ? stockRating.stockSentiment : null
        },
        weights: {
          market: hasMarket ? marketWeight : 0,
          sector: hasSector ? sectorWeight : 0,
          stock: hasStock ? stockWeight : 0
        }
      };
      
      // After calculating synthetic rating, get recommended strategies
      recommendedStrategies = getRecommendedStrategies(syntheticRating.numericValue);
    }
  }
  
  // Function to get recommended strategies based on synthetic rating
  function getRecommendedStrategies(rating) {
    if (rating === null || rating === undefined) return null;
    
    // Define the strategy groups based on rating ranges
    if (rating > 2.0 && rating <= 3.0) {
      return {
        sentiment: "Very Bullish",
        color: "blue",
        colorCode: "#4299E1", // blue-500
        indicatorEmoji: "🔵",
        strategies: [
          { name: "Long Calls", ratio: "1:3 or better", description: "highly favorable payoff" },
          { name: "Bull Call Spread", ratio: "1:2 to 1:4", description: "capped upside, controlled downside" },
          { name: "Bull Put Spread (Credit Spread)", ratio: "2:1 to 4:1", description: "high probability credit strategy" },
          { name: "Covered Call", ratio: "1:1 to 2:1", description: "income-driven strategy with limited upside" },
          { name: "Cash Secured Puts", ratio: "2:1 to 3:1", description: "higher probability, moderate reward" },
          { name: "Long Synthetic Stock", ratio: "1:3 or better", description: "similar payoff to outright stock" },
          { name: "Call Ratio Backspread", ratio: "1:3 or better", description: "high upside, explosive moves" }
        ]
      };
    }
    else if (rating > 1.0 && rating <= 2.0) {
      return {
        sentiment: "Moderately Bullish",
        color: "green",
        colorCode: "#48BB78", // green-500
        indicatorEmoji: "🟢",
        strategies: [
          { name: "Covered Call", ratio: "1:1 to 2:1", description: "consistent, steady returns" },
          { name: "Cash Secured Puts", ratio: "2:1 to 3:1", description: "moderate return with higher probability" },
          { name: "Bull Call Spread", ratio: "1:1 to 1:3", description: "balanced approach" },
          { name: "Bull Put Spread", ratio: "2:1 to 3:1", description: "credit-driven, moderate reward" },
          { name: "Collar", ratio: "1:1 to 1:2", description: "protective strategy, conservative gains" },
          { name: "Diagonal Call Spreads", ratio: "1:2 to 1:4", description: "moderate directional bet, limited risk" }
        ]
      };
    }
    else if (rating > 0.5 && rating <= 1.0) {
      return {
        sentiment: "Mildly Bullish",
        color: "teal",
        colorCode: "#38B2AC", // teal-500
        indicatorEmoji: "🟢",
        strategies: [
          { name: "Bull Put Spread", ratio: "2:1 to 4:1", description: "consistent small gains" },
          { name: "Covered Calls (conservative strike)", ratio: "1:1 to 2:1", description: "steady, income-oriented" },
          { name: "Cash Secured Puts", ratio: "2:1 to 3:1", description: "probability-focused approach" },
          { name: "Long Stock with Protective Put", ratio: "1:1 to 1:2", description: "safety over high return" },
          { name: "Bull Call Spread (wide)", ratio: "1:2 to 1:3", description: "directional bet with defined risk" }
        ]
      };
    }
    else if (rating >= -0.5 && rating <= 0.5) {
      return {
        sentiment: "Neutral",
        color: "gray",
        colorCode: "#718096", // gray-600
        indicatorEmoji: "⚫",
        strategies: [
          { name: "Iron Condors", ratio: "3:1 to 5:1", description: "high probability, capped reward" },
          { name: "Butterflies", ratio: "1:4 to 1:8", description: "low probability, high reward at specific price" },
          { name: "Short Straddles/Strangles", ratio: "2:1 to 4:1", description: "income-driven, moderate risk" },
          { name: "Calendar Spreads", ratio: "1:2 to 1:5", description: "time-decay and volatility sensitive" }
        ]
      };
    }
    else if (rating >= -1.0 && rating < -0.5) {
      return {
        sentiment: "Mildly Bearish",
        color: "yellow",
        colorCode: "#D69E2E", // yellow-600
        indicatorEmoji: "🟡",
        strategies: [
          { name: "Bear Call Spread", ratio: "2:1 to 4:1", description: "higher probability credit strategy" },
          { name: "Put Calendar Spreads", ratio: "1:2 to 1:4", description: "mild bearish bias, volatility sensitive" },
          { name: "Bear Put Spread (wide)", ratio: "1:2 to 1:3", description: "directional bet with defined risk" },
          { name: "Covered Put", ratio: "1:1 to 2:1", description: "income-oriented bearish play" },
          { name: "Bear Call Ladder", ratio: "2:1 to 3:1", description: "credit strategy with managed risk" }
        ]
      };
    }
    else if (rating >= -2.0 && rating < -1.0) {
      return {
        sentiment: "Moderately Bearish",
        color: "orange",
        colorCode: "#ED8936", // orange-500
        indicatorEmoji: "🟠",
        strategies: [
          { name: "Bear Put Spread", ratio: "1:1 to 1:3", description: "moderate directional play" },
          { name: "Bear Call Spread", ratio: "2:1 to 4:1", description: "higher probability credit strategy" },
          { name: "Long Puts", ratio: "1:3 or better", description: "direct downside exposure, higher payoff" },
          { name: "Diagonal Put Spreads", ratio: "1:2 to 1:4", description: "moderate risk, volatility-structured" },
          { name: "Put Debit Spreads", ratio: "1:2 or better", description: "defined risk with directional bias" }
        ]
      };
    }
    else if (rating >= -3.0 && rating < -2.0) {
      return {
        sentiment: "Very Bearish",
        color: "red",
        colorCode: "#E53E3E", // red-600
        indicatorEmoji: "🔴",
        strategies: [
          { name: "Long Puts", ratio: "1:3 or better", description: "large downside gains" },
          { name: "Bear Put Spreads", ratio: "1:2 to 1:4", description: "aggressively bearish, capped upside" },
          { name: "Bear Call Spread", ratio: "2:1 to 4:1", description: "high probability, credit income" },
          { name: "Short Synthetic Stock", ratio: "1:3 or better", description: "mimics short equity position" },
          { name: "Put Ratio Backspread", ratio: "1:3 or better", description: "significant bearish move, explosive profit potential" }
        ]
      };
    }
    
    // Fallback case for any other values
    return null;
  }
  
  // Helper function to select a strategy when clicking on a recommendation
  function selectStrategy(strategyObj) {
    // Extract the strategy name from the object
    const strategyName = strategyObj.name;
    
    // Find the matching strategy in strategyOptions
    let found = false;
    
    // Search through all categories and options
    for (const category of strategyOptions) {
      for (const option of category.options) {
        // Extract the key part of the strategy name to match
        // For example, "Bull Call Spread (conservative strikes)" -> "Bull Call Spread"
        const simplifiedName = strategyName.split('(')[0].trim();
        
        if (option.name.includes(simplifiedName) || simplifiedName.includes(option.name)) {
          // Found a match, update the trade form
          newTrade.strategy = `${category.category} - ${option.name}`;
          found = true;
          break;
        }
      }
      if (found) break;
    }
    
    // If no exact match, try a more flexible match for common terms
    if (!found) {
      const strategyLower = strategyName.toLowerCase();
      let bestMatchCategory;
      let bestMatchOption;
      
      // Common generic mappings
      if (strategyLower.includes("call spread")) {
        bestMatchCategory = 'Vertical Spreads';
        if (strategyLower.includes("bull")) {
          bestMatchOption = 'Bull Call Spread';
        } else if (strategyLower.includes("bear")) {
          bestMatchOption = 'Bear Call Spread';
        }
      } else if (strategyLower.includes("put spread")) {
        bestMatchCategory = 'Vertical Spreads';
        if (strategyLower.includes("bull")) {
          bestMatchOption = 'Bull Put Spread';
        } else if (strategyLower.includes("bear")) {
          bestMatchOption = 'Bear Put Spread';
        }
      } else if (strategyLower.includes("long call")) {
        bestMatchCategory = 'Basic Spreads';
        bestMatchOption = 'Long Call';
      } else if (strategyLower.includes("long put")) {
        bestMatchCategory = 'Basic Spreads';
        bestMatchOption = 'Long Put';
      } else if (strategyLower.includes("covered call")) {
        bestMatchCategory = 'Basic Spreads';
        bestMatchOption = 'Covered Call';
      } else if (strategyLower.includes("calendar")) {
        bestMatchCategory = 'Calendar/Horizontal Spreads';
        if (strategyLower.includes("call")) {
          bestMatchOption = 'Long Calendar Call Spread';
        } else if (strategyLower.includes("put")) {
          bestMatchOption = 'Long Calendar Put Spread';
        }
      } else if (strategyLower.includes("iron condor")) {
        bestMatchCategory = 'Iron Condors/Butterflies';
        bestMatchOption = 'Iron Condor';
      } else if (strategyLower.includes("butterfly")) {
        bestMatchCategory = 'Butterfly Spreads';
        if (strategyLower.includes("call")) {
          bestMatchOption = 'Long Call Butterfly';
        } else if (strategyLower.includes("put")) {
          bestMatchOption = 'Long Put Butterfly';
        }
      }
      
      // If we found a match in the common mappings
      if (bestMatchCategory && bestMatchOption) {
        newTrade.strategy = `${bestMatchCategory} - ${bestMatchOption}`;
      }
    }
  }

  // --- End Rating Fetching ---

  onMount(async () => {
    console.log('TradeCalendar component mounted');
    await loadTrades();
    await fetchLatestMarketRating(); // Fetch market rating on load
    await tick();
    console.log('After tick, trades length:', trades.length);
    
    // Force DOM update after a short delay as a fallback
    setTimeout(() => {
      console.log('Forcing update after timeout');
      trades = [...trades]; 
    }, 100);
  });

  // Show the rating editor modal
  function showEditRating(type) {
    editingRatingType = type;
    
    if (type === 'sector' && sectorRating) {
      // Create a copy of the rating to edit
      editingRating = { ...sectorRating, stockSentiment: sectorRating.stockSentiment };
    } else if (type === 'stock' && stockRating) {
      // Create a copy of the rating to edit
      editingRating = { ...stockRating, stockSentiment: stockRating.stockSentiment };
    } else {
      // Create a new rating if none exists
      const today = new Date().toISOString().split('T')[0];
      editingRating = {
        id: '',
        date: `${today}T00:00:00Z`,
        symbol: type === 'sector' ? 'SECTOR' : newTrade.ticker.toUpperCase(),
        sector: newTrade.sector,
        stockSentiment: 0,
        notes: ''
      };
    }
    
    showRatingEditor = true;
  }
  
  // Save the edited rating
  async function saveEditedRating() {
    try {
      // Get date in ISO format
      const targetDate = new Date().toISOString().split('T')[0];
      
      // Create a compatible object for saving
      const ratingData = {
        id: editingRating.id || '',
        date: editingRating.date || `${targetDate}T00:00:00Z`,
        symbol: editingRatingType === 'sector' ? 'SECTOR' : newTrade.ticker.toUpperCase(),
        sector: newTrade.sector,
        stockSentiment: editingRating.stockSentiment,
        priceTarget: editingRating.priceTarget || 0,
        confidence: editingRating.confidence || 0,
        enthusiasm: editingRating.enthusiasm || 0,
        chartPattern: editingRating.chartPattern || '',
        notes: editingRating.notes || ''
      };
      
      // Create a proper StockRating object
      const ratingToSave = models.StockRating.createFrom(ratingData);
      console.log('Saving rating with explicit fields:', ratingToSave);
      
      await SaveStockRating(ratingToSave);
      
      // Close the editor
      showRatingEditor = false;
      
      // Refresh the appropriate rating
      if (editingRatingType === 'sector') {
        sectorRating = null;
        await handleSectorChange({ target: { value: newTrade.sector } });
      } else if (editingRatingType === 'stock') {
        stockRating = null;
        if (newTrade.ticker) {
          await refreshStockRating();
        }
      }
      
      // Recalculate synthetic rating
      calculateSyntheticRating();
      
      alert('Rating updated successfully!');
    } catch (error) {
      console.error('Failed to save rating:', error);
      alert('Failed to save rating: ' + error.message);
    }
  }
  
  // Cancel editing
  function cancelEditRating() {
    showRatingEditor = false;
    editingRating = null;
  }
</script>

<div class="calendar-page">
  <!-- Tab navigation -->
  <div class="tab-navigation">
    <button 
      class:active={activeTab === 'calendar'} 
      on:click={() => activeTab = 'calendar'}
    >
      Trades + Options Calendar
    </button>
    <button 
      class:active={activeTab === 'journal'} 
      on:click={() => activeTab = 'journal'}
    >
      Trading Journal
    </button>
    
    <div class="tab-actions">
      <!-- Refresh button removed -->
    </div>
  </div>

  {#if activeTab === 'calendar'}
    <div class="calendar-section">
      <h2>Trades + Options Calendar</h2>
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
        {#key componentKey}
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
                {#each weeks as week}
                  <td class="trade-cell" class:has-trades={hasTrades(sector, week.number)}>
                    {#if tradesLoaded}
                      {#each getTradesForCell(sector, week.number) as trade (trade.id + '_' + week.number)}
                        <div 
                          class="trade-card"
                          class:spanning-trade={trade.isSpanningLeg}
                          class:span-start={trade.spanPosition === 'start'}
                          class:span-middle={trade.spanPosition === 'middle'}
                          class:span-end={trade.spanPosition === 'end'}
                          style="background-color: {getStrategyColor(trade.strategy)}"
                          on:click={() => handleTradeCardClick(trade)}
                        >
                          <div class="trade-symbol">{trade.symbol}</div>
                          <div class="trade-type">{trade.type}</div>
                          {#if trade.isSpanningLeg}
                            <div class="trade-spanning-indicator">
                              {#if trade.spanPosition === 'start'}
                                ►
                              {:else if trade.spanPosition === 'end'}
                                ◄
                              {:else if trade.spanPosition === 'middle'}
                                ↔
                              {/if}
                            </div>
                          {/if}
                        </div>
                      {/each}
                    {/if}
                  </td>
                {/each}
              </tr>
            {/each}
          </tbody>
        </table>
        {/key}
      </div>
      
      <div class="add-trade-section">
        <h3>{isEditing ? 'Edit' : 'Add New'} Options Trade</h3>
        
        <!-- Rating Display Area -->
        <div class="trade-ratings-display">
          <h4>Current Ratings:</h4>
          <div class="rating-pills">
            <span class="rating-pill market" class:loading={ratingsLoading.market}>
              Market: 
              {#if marketRating}
                {#if marketRating.error}
                  <span class="error">{marketRating.error}</span>
                {:else}
                  <strong>{marketRating.stockSentiment ?? 'N/A'}</strong> 
                  ({new Date(marketRating.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })})
                {/if}
              {:else if ratingsLoading.market}
                Loading...
              {:else}
                N/A
              {/if}
            </span>

            {#if newTrade.sector}
            <span class="rating-pill sector" class:loading={ratingsLoading.sector}>
              <span class="pill-content">
                <span class="pill-title">Sector ({newTrade.sector}):</span>
                {#if sectorRating}
                  {#if sectorRating.error}
                    <span class="error">{sectorRating.error}</span>
                  {:else}
                   <strong>{sectorRating.stockSentiment ?? 'N/A'}</strong> 
                   ({new Date(sectorRating.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })})
                   <button class="edit-rating-btn" on:click|preventDefault={() => showEditRating('sector')}>✎</button>
                  {/if}
                {:else if ratingsLoading.sector}
                  Loading...
                {:else}
                  N/A
                  <button class="edit-rating-btn" on:click|preventDefault={() => showEditRating('sector')}>+</button>
                {/if}
              </span>
            </span>
            {/if}

            {#if newTrade.ticker}
            <span class="rating-pill stock" class:loading={ratingsLoading.stock}>
              <span class="pill-content">
                <span class="pill-title">Ticker ({newTrade.ticker.toUpperCase()}):</span>
                {#if stockRating && !stockRatingError}
                  <strong>{stockRating.stockSentiment ?? 'N/A'}</strong> 
                  ({new Date(stockRating.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })})
                  <button class="edit-rating-btn" on:click|preventDefault={() => showEditRating('stock')}>✎</button>
                {:else if stockRatingError}
                  <span class="error">{stockRatingError.message}</span>
                  <button 
                    class="retry-btn" 
                    on:click|preventDefault={refreshStockRating}
                    disabled={ratingsLoading.stock}
                  >
                    Try again
                  </button>
                  <button class="edit-rating-btn" on:click|preventDefault={() => showEditRating('stock')}>+</button>
                {:else if ratingsLoading.stock}
                  Loading...
                {:else}
                  N/A
                  <button 
                    class="retry-btn" 
                    on:click|preventDefault={refreshStockRating}
                    disabled={ratingsLoading.stock}
                  >
                    Try again
                  </button>
                  <button class="edit-rating-btn" on:click|preventDefault={() => showEditRating('stock')}>+</button>
                {/if}
              </span>
            </span>
            {/if}
            
            {#if syntheticRating}
            <span class="rating-pill synthetic">
              <span class="pill-title">Synthetic Outlook:</span>
              <strong>{syntheticRating.numericValue}</strong>
              <span class="tooltip">
                <span class="tooltip-icon">ⓘ</span>
                <span class="tooltip-text">
                  Weighted average: 
                  {#if syntheticRating.components.market !== null}
                    Market ({syntheticRating.components.market}) × 0.5
                    {#if syntheticRating.components.sector !== null || syntheticRating.components.stock !== null} + {/if}
                  {/if}
                  {#if syntheticRating.components.sector !== null}
                    Sector ({syntheticRating.components.sector}) × 0.3
                    {#if syntheticRating.components.stock !== null} + {/if}
                  {/if}
                  {#if syntheticRating.components.stock !== null}
                    Stock ({syntheticRating.components.stock}) × 0.2
                  {/if}
                </span>
              </span>
            </span>
            {/if}
          </div>
        </div>
        <!-- End Rating Display Area -->
        
        <div class="trade-form">
          <div class="form-row">
            <div class="form-group">
              <label>Entry Date:</label>
              <input type="date" bind:value={newTrade.entryDate} />
            </div>
            
            <div class="form-group">
              <label>Primary Expiration Date:</label>
              <input type="date" bind:value={newTrade.expirationDate} />
              <div class="date-guidance">
                <strong>Available expiration weeks:</strong>
                <div class="week-options">
                  {#each weeks as week}
                    <div class="week-option">
                      {week.range}
                    </div>
                  {/each}
                </div>
              </div>
            </div>
            
            <div class="form-group">
              <label>Ticker:</label>
              <div class="input-with-action">
                <input type="text" placeholder="e.g., AAPL" bind:value={newTrade.ticker} on:blur={handleTickerBlur} />
                {#if newTrade.ticker}
                  <button 
                    class="refresh-rating-btn" 
                    on:click|preventDefault={refreshStockRating}
                    title="Refresh stock rating"
                    disabled={ratingsLoading.stock}
                  >
                    {#if ratingsLoading.stock}⟳{:else}↻{/if}
                  </button>
                {/if}
              </div>
            </div>
            
            <div class="form-group">
              <label>Sector:</label>
              <select bind:value={newTrade.sector} on:change={handleSectorChange}>
                <option value="">Select a sector...</option>
                {#each sectors as sector}
                  <option value={sector}>{sector}</option>
                {/each}
              </select>
            </div>
          </div>
          
          <!-- New row for trade details -->
          <div class="form-row">
            <div class="form-group">
              <label>Timeframe:</label>
              <select bind:value={newTrade.timeframe}>
                <option value="">Select timeframe...</option>
                <option value="1 week">1 week</option>
                <option value="2 weeks">2 weeks</option>
                <option value="3 weeks">3 weeks</option>
                <option value="4 weeks">4 weeks</option>
                <option value="5-6 weeks">5-6 weeks</option>
                <option value="7-8 weeks">7-8 weeks</option>
                <option value="9-10 weeks">9-10 weeks</option>
                <option value="11-12 weeks">11-12 weeks</option>
                <option value="3 months">3 months</option>
                <option value="4 months">4 months</option>
                <option value="5 months">5 months</option>
              </select>
            </div>
            
            <div class="form-group">
              <label>Entry Price:</label>
              <input type="number" step="0.01" min="0" bind:value={newTrade.entry} />
            </div>
            
            <div class="form-group">
              <label>Stop Loss:</label>
              <input type="number" step="0.01" min="0" bind:value={newTrade.stop} />
            </div>
            
            <div class="form-group">
              <label>Target Price:</label>
              <input type="number" step="0.01" min="0" bind:value={newTrade.target} />
            </div>
          </div>
          
          <div class="form-row">
            <div class="form-group full-width">
              <label>Options Strategy:</label>
              <select bind:value={newTrade.strategy}>
                <option value="">Select a strategy...</option>
                {#each strategyOptions as category}
                  <optgroup label={category.category}>
                    {#each category.options as option}
                      <option value={`${category.category} - ${option.name}`}>{option.name}</option>
                    {/each}
                  </optgroup>
                {/each}
              </select>
              
              {#if selectedStrategyDescription}
                <div class="strategy-description">{selectedStrategyDescription}</div>
              {/if}
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
          
          {#if showShortLegField}
            <div class="form-row">
              <div class="form-group full-width">
                <label>Short Leg Expiration Date:</label>
                <input type="date" bind:value={newTrade.shortLegExpiration} />
              </div>
            </div>
          {/if}
          
          <div class="form-buttons">
            <button class="reset-btn" on:click={resetForm}>
              {isEditing ? 'Cancel' : 'Reset'}
            </button>
            <button class="add-btn" on:click={addNewTrade}>
              {isEditing ? 'Update Trade' : 'Save Trade'}
            </button>
          </div>
          
          <!-- Recommended Strategies Based on Synthetic Rating -->
          {#if recommendedStrategies}
          <div class="strategy-recommendations">
            <div class="recommendation-header" style="border-color: {recommendedStrategies.colorCode};">
              <span class="sentiment-indicator">{recommendedStrategies.indicatorEmoji}</span>
              <h4>Recommended Strategies for {recommendedStrategies.sentiment} Outlook ({syntheticRating.numericValue})</h4>
            </div>
            <div class="recommended-strategies-list">
              {#each recommendedStrategies.strategies as strategy}
                <div class="strategy-item" 
                     style="border-left-color: {recommendedStrategies.colorCode};"
                     on:click={() => selectStrategy(strategy)}
                     on:keydown={(e) => e.key === 'Enter' && selectStrategy(strategy)}
                     tabindex="0"
                     role="button"
                     title="Click to select this strategy">
                  <div class="strategy-item-left">
                    <div class="strategy-name">{strategy.name}</div>
                    <div class="strategy-ratio">R:R = {strategy.ratio}</div>
                  </div>
                  <div class="strategy-description">{strategy.description}</div>
                </div>
              {/each}
            </div>
          </div>
          {/if}
        </div>
      </div>
      
      <div class="trade-history">
        <h3>Options Trade History</h3>
        
        <div class="search-section">
          <input 
            type="text" 
            placeholder="Search trades..." 
            class="search-input" 
            bind:value={searchQuery}
          />
          
          <div class="search-filters">
            <label>
              <input type="radio" name="searchType" value="ticker" bind:group={searchType} />
              Ticker
            </label>
            <label>
              <input type="radio" name="searchType" value="sector" bind:group={searchType} />
              Sector
            </label>
            <label>
              <input type="radio" name="searchType" value="strategy" bind:group={searchType} />
              Strategy
            </label>
            <label>
              <input type="radio" name="searchType" value="date" bind:group={searchType} />
              Date
            </label>
          </div>
        </div>
        
        {#if filteredTrades.length === 0}
          <div class="no-trades">
            <p>No trades found. Add your first trade above.</p>
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
                {#each filteredTrades as trade (trade.id + '_' + trade.expirationDate.toISOString())}
                  <tr>
                    <td class="ticker-cell">{trade.symbol}</td>
                    <td>{trade.sector}</td>
                    <td>
                      <span class="strategy-pill" style="background-color: {getStrategyColor(trade.strategy)}">
                        {trade.type}
                      </span>
                    </td>
                    <td>{formatDate(trade.entryDate)}</td>
                    <td>{formatDate(trade.expirationDate)}</td>
                    <td>${trade.entryPrice.toFixed(2)}</td>
                    <td class="actions-cell">
                      <button 
                        class="action-btn edit-btn" 
                        on:click={() => editTrade(trade.id)}
                        title="Edit Trade"
                      >
                        ✎
                      </button>
                      <button 
                        class="action-btn journal-btn" 
                        on:click={() => { selectedTradeForJournal = trade; activeTab = 'journal'; }}
                        title="Add to Journal"
                      >
                        📝
                      </button>
                      <button 
                        class="action-btn delete-btn" 
                        on:click={() => deleteTrade(trade.id)}
                        title="Delete Trade"
                      >
                        ✕
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
    </div>
  {:else if activeTab === 'journal'}
    <TradeJournal 
      trades={trades} 
      {strategyOptions} 
      {sectors} 
      editExistingTrade={selectedTradeForJournal}
      refreshTrades={forceRefresh}
    />
  {/if}
  
  <!-- Rating Editor Modal -->
  {#if showRatingEditor}
  <div class="modal-overlay">
    <div class="rating-editor-modal">
      <h3>Edit {editingRatingType === 'sector' ? 'Sector' : 'Stock'} Rating</h3>
      
      <div class="rating-editor-content">
        <div class="rating-info">
          <div class="info-row">
            <strong>{editingRatingType === 'sector' ? 'Sector' : 'Symbol'}:</strong>
            <span>{editingRatingType === 'sector' ? newTrade.sector : newTrade.ticker.toUpperCase()}</span>
          </div>
          <div class="info-row">
            <strong>Date:</strong>
            <span>{new Date().toLocaleDateString()}</span>
          </div>
        </div>
        
        <div class="rating-slider-container">
          <label for="rating-value">Sentiment Rating (-3 to +3):</label>
          <div class="sentiment-control">
            <span class="sentiment-label negative">Bearish</span>
            <input 
              type="range" 
              id="rating-value" 
              min="-3" 
              max="3" 
              step="1"
              bind:value={editingRating.stockSentiment}
            />
            <span class="sentiment-label positive">Bullish</span>
          </div>
          <div class="sentiment-value">
            Current value: <strong>{editingRating.stockSentiment}</strong>
          </div>
        </div>
        
        <div class="rating-notes">
          <label for="rating-notes">Notes:</label>
          <textarea 
            id="rating-notes" 
            bind:value={editingRating.notes} 
            placeholder="Add any notes about this rating..."
          ></textarea>
        </div>
      </div>
      
      <div class="rating-editor-actions">
        <button class="cancel-btn" on:click={cancelEditRating}>Cancel</button>
        <button class="save-btn" on:click={saveEditedRating}>Save Rating</button>
      </div>
    </div>
  </div>
  {/if}
</div>

<style>
  .calendar-page {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
    background-color: var(--content-bg);
    border-radius: 5px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  /* Tab Navigation Styles */
  .tab-navigation {
    display: flex;
    margin-bottom: 1.5rem;
    border-bottom: 1px solid var(--border-color);
  }
  
  .tab-navigation button {
    padding: 0.75rem 1.5rem;
    background: none;
    border: none;
    cursor: pointer;
    font-weight: 500;
    color: inherit;
    border-bottom: 3px solid transparent;
    transition: all 0.2s;
    margin-bottom: -1px;
  }
  
  .tab-navigation button.active {
    font-weight: bold;
    border-bottom: 3px solid var(--primary-button);
    color: var(--primary-button);
  }
  
  .tab-navigation button:hover:not(.active) {
    background-color: rgba(0, 0, 0, 0.05);
  }
  
  .tab-actions {
    margin-left: auto;
    display: flex;
    align-items: center;
  }
  
  .calendar-section {
    /* This wrapper ensures the calendar content layout isn't affected */
  }

  .journal-btn {
    background-color: #805ad5; /* Purple-ish color for journal button */
  }
  
  h2, h3, h4 {
    color: inherit;
    margin-top: 0;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 0;
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
  
  .strategy-guide {
    background-color: var(--bg-color);
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
    max-width: 100%;
  }
  
  .calendar {
    width: 100%;
    border-collapse: collapse;
    border: 1px solid var(--border-color);
    table-layout: fixed; /* Use fixed layout for better space control */
  }
  
  .sector-header {
    width: 120px; /* Reduce width */
    text-align: left;
    background-color: var(--bg-color);
    padding: 0.4rem 0.5rem; /* Reduce padding */
  }
  
  .week-header {
    min-width: 100px; /* Reduce min-width */
    max-width: 120px; /* Add max-width */
    background-color: var(--bg-color);
    text-align: center;
    font-size: 0.85rem; /* Slightly smaller font */
    padding: 0.4rem 0.3rem; /* Reduce padding */
  }
  
  .week-label {
    font-weight: bold;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .exp-date {
    font-size: 0.75rem;
    color: #718096;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .sector-cell {
    padding: 0.5rem;
    font-weight: 500;
    border: 1px solid var(--border-color);
    background-color: var(--bg-color);
    font-size: 0.9rem; /* Slightly smaller font */
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .trade-cell {
    border: 1px solid var(--border-color);
    padding: 0.4rem; /* Reduce padding */
    vertical-align: top;
    height: 70px; /* Reduce height */
    position: relative;
    overflow-y: auto;
    max-height: 110px;
  }
  
  .trade-cell.has-trades {
    background-color: rgba(247, 250, 252, 0.1);
  }
  
  .trade-card {
    padding: 0.35rem 0.25rem; /* Reduce padding */
    border-radius: 3px; /* Smaller radius */
    margin-bottom: 0.35rem; /* Smaller margin */
    color: white;
    font-size: 0.75rem; /* Smaller font */
    text-align: center;
    cursor: pointer;
    transition: transform 0.1s ease-in-out;
    box-shadow: 0 1px 2px rgba(0,0,0,0.1);
  }
  
  .trade-card.spanning-trade {
    border: 2px dashed rgba(255, 255, 255, 0.7);
    opacity: 0.85;
  }
  
  .trade-card.span-start {
    border-right: none;
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
    border-right: 2px dashed rgba(255, 255, 255, 0.7);
    margin-right: -1px;
    position: relative;
    z-index: 1;
  }
  
  .trade-card.span-middle {
    border-radius: 0;
    border-left: none;
    border-right: none;
    border-top: 2px dashed rgba(255, 255, 255, 0.7);
    border-bottom: 2px dashed rgba(255, 255, 255, 0.7);
    position: relative;
    z-index: 1;
  }
  
  .trade-card.span-end {
    border-left: none;
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
    border-left: 2px dashed rgba(255, 255, 255, 0.7);
    margin-left: -1px;
    position: relative;
    z-index: 1;
  }
  
  .trade-spanning-indicator {
    font-size: 0.7rem;
    margin-top: 0.25rem;
    font-weight: bold;
  }
  
  .trade-card:hover {
    transform: scale(1.05);
    z-index: 2;
  }
  
  .trade-symbol {
    font-weight: bold;
    margin-bottom: 0.15rem; /* Smaller margin */
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .trade-type {
    font-size: 0.7rem; /* Smaller font */
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .add-trade-section {
    background-color: var(--card-bg);
    padding: 1.5rem;
    border-radius: 5px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
    margin-bottom: 2rem;
    transition: background-color 0.3s ease, color 0.3s ease;
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
    color: inherit;
  }
  
  .form-group input,
  .form-group select,
  .form-group textarea {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid var(--input-border);
    border-radius: 4px;
    font-size: 0.9rem;
    background-color: var(--input-bg);
    color: var(--text-color);
  }
  
  .form-group textarea {
    resize: vertical;
    min-height: 100px;
  }
  
  .strategy-description {
    margin-top: 0.5rem;
    padding: 0.75rem;
    background-color: var(--bg-color);
    border-radius: 4px;
    font-size: 0.9rem;
    color: inherit;
    border-left: 3px solid var(--active-button);
  }
  
  .date-guidance {
    margin-top: 0.25rem;
    font-size: 0.8rem;
    color: #cbd5e0;
  }
  
  .week-options {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    margin-top: 0.25rem;
    max-height: 80px;
    overflow-y: auto;
    width: 120px; /* Make the dropdown skinnier */
  }
  
  .week-option {
    padding: 0.2rem 0.4rem;
    background-color: var(--bg-color);
    border-radius: 3px;
    font-size: 0.75rem;
    white-space: nowrap;
    text-align: center;
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
    background-color: var(--border-color);
    color: inherit;
  }
  
  .add-btn {
    background-color: var(--active-button);
    color: white;
  }
  
  .reset-btn:hover {
    background-color: #cbd5e0;
  }
  
  .add-btn:hover {
    background-color: #3182ce;
  }
  
  .trade-history {
    background-color: var(--card-bg);
    padding: 1.5rem;
    border-radius: 5px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
    transition: background-color 0.3s ease, color 0.3s ease;
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
  }
  
  .search-filters label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }
  
  .no-trades {
    text-align: center;
    color: #718096;
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
    color: white;
    font-size: 0.8rem;
    font-weight: 500;
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
    background-color: #4299e1;
  }
  
  .delete-btn {
    background-color: #e53e3e;
  }
  
  .notes-row {
    background-color: #f8fafc;
  }
  
  .notes-content {
    padding: 0.5rem 0;
    font-size: 0.9rem;
    color: #4a5568;
  }
  
  .header-container {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 1rem;
  }
  
  .refresh-btn {
    background-color: var(--active-button);
    color: white;
    border: none;
    border-radius: 4px;
    padding: 0.5rem 1rem;
    font-weight: 500;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.25rem;
  }
  
  .refresh-btn:hover {
    background-color: #3182ce;
  }
  
  :global(body.dark-mode) .sector-cell,
  :global(body.dark-mode) .week-header,
  :global(body.dark-mode) .trades-table th {
    background-color: #2d3748;
    color: #e2e8f0;
  }
  
  :global(body.dark-mode) .trade-cell.has-trades {
    background-color: rgba(45, 55, 72, 0.3);
  }
  
  :global(body.dark-mode) .refresh-btn,
  :global(body.dark-mode) .add-btn {
    background-color: #4fd1c5;
  }
  
  :global(body.dark-mode) .reset-btn {
    background-color: #4a5568;
    color: #e2e8f0;
  }
  
  :global(body.dark-mode) .add-trade-section {
    background-color: var(--card-bg);
    color: var(--text-color);
  }
  
  :global(body.dark-mode) .date-guidance,
  :global(body.dark-mode) .additional-note {
    color: #cbd5e0;
  }
  
  :global(body.dark-mode) .form-group input,
  :global(body.dark-mode) .form-group select,
  :global(body.dark-mode) .form-group textarea {
    border-color: #4a5568;
  }
  
  :global(body.dark-mode) .week-option {
    background-color: #4a5568;
    color: #e2e8f0;
  }

  .calendar-section {
    background-color: var(--card-bg);
    padding: 1rem;
    border-radius: 5px;
    box-shadow: 0 1px 2px rgba(0,0,0,0.05);
    transition: background-color 0.3s ease, color 0.3s ease;
  }

  button {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s ease;
    background-color: var(--primary-button);
    color: white;
  }

  button.secondary {
    background-color: var(--secondary-button);
    color: var(--text-color);
  }

  button:hover {
    opacity: 0.9;
  }

  button:disabled {
    background-color: #a0aec0;
    cursor: not-allowed;
  }

  :global(body.dark-mode) .calendar-section {
    background-color: var(--card-bg);
    color: var(--text-color);
  }

  :global(body.dark-mode) button.secondary {
    background-color: #4a5568;
    color: #e2e8f0;
  }

  /* Styles for Rating Display */
  .trade-ratings-display {
    background-color: var(--input-bg);
    padding: 0.75rem 1rem;
    margin-bottom: 1.5rem;
    border-radius: 4px;
    border: 1px solid var(--border-color);
  }

  .trade-ratings-display h4 {
    margin: 0 0 0.75rem 0;
    font-size: 0.9rem;
    font-weight: 600;
    color: inherit;
  }

  .rating-pills {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
  }

  .rating-pill {
    display: inline-flex;
    align-items: center;
    padding: 0.3rem 0.8rem;
    border-radius: 15px;
    font-size: 0.85rem;
    background-color: var(--bg-color);
    border: 1px solid transparent;
    transition: background-color 0.2s, border-color 0.2s;
  }
  
  .rating-pill.loading {
    opacity: 0.7;
    cursor: default;
  }

  .rating-pill strong {
    margin: 0 0.3rem;
    font-weight: bold;
    font-size: 0.9rem;
  }
  
  .rating-pill .error {
    color: #e53e3e; /* Red for errors */
    font-style: italic;
    margin-left: 0.3rem;
  }
  
  .rating-pill.market {
    border-color: #ECC94B; /* Gold border */
  }
  .rating-pill.sector {
    border-color: #63B3ED; /* Blue border */
  }
  .rating-pill.stock {
    border-color: #68D391; /* Green border */
  }
  
  .rating-pill.synthetic {
    border-color: #805AD5; /* Purple border */
    background-color: rgba(128, 90, 213, 0.1); /* Light purple background */
  }
  
  .pill-title {
    font-weight: 500;
    margin-right: 0.3rem;
  }
  
  .tooltip {
    position: relative;
    display: inline-block;
    margin-left: 0.5rem;
  }
  
  .tooltip-icon {
    color: #718096;
    cursor: help;
    font-size: 0.8rem;
  }
  
  .tooltip-text {
    visibility: hidden;
    width: 250px;
    background-color: #4A5568;
    color: #fff;
    text-align: center;
    border-radius: 6px;
    padding: 8px;
    position: absolute;
    z-index: 1;
    bottom: 125%; /* Position above the icon */
    left: 50%;
    margin-left: -125px; /* Center horizontally */
    opacity: 0;
    transition: opacity 0.3s;
    font-size: 0.75rem;
    line-height: 1.4;
    pointer-events: none;
  }
  
  /* Arrow for tooltip */
  .tooltip-text::after {
    content: "";
    position: absolute;
    top: 100%;
    left: 50%;
    margin-left: -5px;
    border-width: 5px;
    border-style: solid;
    border-color: #4A5568 transparent transparent transparent;
  }
  
  /* Show tooltip on hover */
  .tooltip:hover .tooltip-text {
    visibility: visible;
    opacity: 1;
  }
  
  /* Strategy Recommendations Styles */
  .strategy-recommendations {
    background-color: var(--input-bg);
    border-radius: 4px;
    margin-top: 1.5rem;
    overflow: hidden;
    border: 1px solid var(--border-color);
  }
  
  .recommendation-header {
    padding: 0.75rem 1rem;
    background-color: rgba(0, 0, 0, 0.05);
    border-bottom: 1px solid var(--border-color);
    display: flex;
    align-items: center;
    gap: 0.5rem;
    border-left: 4px solid; /* Color set via inline style */
  }
  
  .recommendation-header h4 {
    margin: 0;
    font-size: 0.95rem;
    font-weight: 600;
  }
  
  .sentiment-indicator {
    font-size: 1.25rem;
  }
  
  .recommended-strategies-list {
    padding: 0.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .strategy-item {
    display: flex;
    align-items: center;
    padding: 0.75rem;
    background-color: var(--bg-color);
    border-radius: 4px;
    border-left: 4px solid;
    cursor: pointer;
    transition: background-color 0.2s ease;
  }
  
  .strategy-item:hover {
    background-color: var(--input-bg);
  }
  
  .strategy-item-left {
    flex: 0 0 30%;
    margin-right: 1rem;
  }
  
  .strategy-name {
    font-weight: 600;
    margin-bottom: 0.25rem;
  }
  
  .strategy-ratio {
    font-size: 0.8rem;
    color: var(--active-button);
  }
  
  .strategy-description {
    flex: 1;
    font-size: 0.9rem;
    color: var(--text-secondary);
    border-left: 1px solid var(--border-color);
    padding-left: 1rem;
  }
  
  @media (max-width: 768px) {
    .strategy-item {
      flex-direction: column;
      align-items: flex-start;
    }
    
    .strategy-item-left {
      flex: 1 0 100%;
      margin-bottom: 0.5rem;
      margin-right: 0;
    }
    
    .strategy-description {
      flex: 1 0 100%;
      border-left: none;
      padding-left: 0;
      border-top: 1px solid var(--border-color);
      padding-top: 0.5rem;
      margin-top: 0.5rem;
    }
  }
  
  /* Dark mode adjustments */
  :global(body.dark-mode) .strategy-recommendations {
    background-color: rgba(45, 55, 72, 0.3);
  }
  
  :global(body.dark-mode) .recommendation-header {
    background-color: rgba(45, 55, 72, 0.5);
  }
  
  :global(body.dark-mode) .strategy-item {
    background-color: rgba(45, 55, 72, 0.8);
  }

  /* Additional styles for the refresh button */
  .input-with-action {
    position: relative;
    display: flex;
  }

  .input-with-action input {
    flex: 1;
  }

  .refresh-rating-btn {
    position: absolute;
    right: 5px;
    top: 50%;
    transform: translateY(-50%);
    background: transparent;
    color: var(--text-secondary);
    border: none;
    font-size: 1rem;
    width: 28px;
    height: 28px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    padding: 0;
    cursor: pointer;
    transition: all 0.2s;
  }

  .refresh-rating-btn:hover {
    background-color: rgba(0, 0, 0, 0.1);
    color: var(--active-button);
  }

  .refresh-rating-btn:disabled {
    opacity: 0.5;
    cursor: wait;
  }

  .retry-btn {
    background: transparent;
    color: var(--active-button);
    border: none;
    padding: 0 0.25rem;
    margin-left: 0.25rem;
    font-size: 0.75rem;
    cursor: pointer;
    text-decoration: underline;
  }

  .retry-btn:hover {
    opacity: 0.8;
  }

  .retry-btn:disabled {
    opacity: 0.5;
    cursor: wait;
  }

  .pill-content {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .rating-editor-modal {
    background-color: white;
    padding: 2rem;
    border-radius: 5px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    width: 300px;
  }

  .rating-editor-content {
    margin-bottom: 2rem;
  }

  .rating-info {
    margin-bottom: 1rem;
  }

  .info-row {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.5rem;
  }

  .rating-slider-container {
    margin-bottom: 1rem;
  }

  .sentiment-control {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .sentiment-label {
    font-weight: bold;
  }

  .sentiment-label.negative {
    color: #e53e3e;
  }

  .sentiment-label.positive {
    color: #48BB78;
  }

  .sentiment-value {
    font-size: 0.8rem;
    color: #718096;
  }

  .rating-notes {
    margin-top: 1rem;
  }

  .rating-editor-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
  }

  .cancel-btn, .save-btn {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-weight: 500;
    cursor: pointer;
  }

  .cancel-btn {
    background-color: #e53e3e;
    color: white;
  }

  .save-btn {
    background-color: #48BB78;
    color: white;
  }

  .edit-rating-btn {
    background: transparent;
    color: var(--primary-button);
    border: none;
    font-size: 0.9rem;
    width: 20px;
    height: 20px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    padding: 0;
    margin-left: 0.25rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .edit-rating-btn:hover {
    background-color: rgba(0, 0, 0, 0.1);
  }

  /* Dark mode styles for the rating editor */
  :global(body.dark-mode) .rating-editor-modal {
    background-color: var(--card-bg);
    color: var(--text-color);
  }

  :global(body.dark-mode) .rating-notes textarea {
    background-color: var(--input-bg);
    color: var(--text-color);
    border-color: var(--border-color);
  }
  
  :global(body.dark-mode) .sentiment-value {
    color: var(--text-secondary);
  }

  :global(body.dark-mode) input[type="range"] {
    background-color: var(--input-bg);
  }
</style> 