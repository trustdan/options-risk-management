<script>
  import { onMount, tick } from 'svelte';
  // Import Wails backend functions
  import {
    GetTrades,
    SaveTrade,
    DeleteTrade
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
  
  // Data for the new trade form
  let newTrade = {
    id: '', // Will be set when editing or generated on save
    entryDate: new Date().toISOString().split('T')[0],
    expirationDate: '',
    ticker: '',
    sector: '',
    entryPrice: 0,
    strategy: '',
    notes: ''
  };
  
  // For multi-leg strategies that span multiple weeks
  let additionalExpirations = [];
  let showAdditionalExpiration = false;
  let selectedStrategyCategory = '';
  let selectedStrategyDescription = '';
  
  // For editing trades
  let isEditing = false;
  let editingTradeId = null;
  
  // For search/filter
  let searchQuery = '';
  let searchType = 'ticker';
  
  // Track whether we're entering a multi-week strategy
  $: {
    const strategyParts = newTrade.strategy.split(' - ');
    if (strategyParts.length >= 2) {
      const category = strategyParts[0];
      const type = strategyParts[1];
      
      showAdditionalExpiration = 
        category.includes('Calendar') || 
        category.includes('Diagonal') ||
        type.includes('Calendar') ||
        type.includes('Diagonal');
    } else {
      showAdditionalExpiration = false;
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
        console.error('Failed to delete old trade legs during edit:', error);
        alert('Error updating trade. Could not remove old data.');
        return;
      }
    }
    
    // Base trade data (plain JS object first)
    const baseTradeData = {
      id: tradeId,
      symbol: newTrade.ticker,
      sector: newTrade.sector,
      strategy: strategyCategory,
      type: strategyType,
      entryPrice: parseFloat(newTrade.entryPrice.toString()),
      notes: newTrade.notes,
      entryDate: entryDate, 
      // week and expirationDate added per leg
    };

    let tradesToSaveData = [];

    // Primary Leg
    tradesToSaveData.push({
      ...baseTradeData,
      week: weekIndex + 1,
      expirationDate: expDate,
    });

    // Additional Legs
    if (showAdditionalExpiration && additionalExpirations.length > 0) {
      for (const additionalExpStr of additionalExpirations) {
        if (!additionalExpStr) continue;
        const addExpDate = new Date(additionalExpStr);
        const addWeekIndex = weeks.findIndex(w => {
          const weekStart = new Date(w.startDate);
          const weekEnd = new Date(w.expirationDate);
          weekEnd.setDate(weekEnd.getDate() + 1);
          return addExpDate >= weekStart && addExpDate < weekEnd;
        });

        if (addWeekIndex !== -1) {
          tradesToSaveData.push({
            ...baseTradeData,
            week: addWeekIndex + 1,
            expirationDate: addExpDate,
          });
        } else {
          alert(`Additional expiration date ${formatDate(addExpDate)} is outside the displayed weeks and was ignored.`);
        }
      }
    }

    // Save all trade legs to backend
    try {
      for (const legData of tradesToSaveData) {
        const tradeModel = new models.Trade();
        Object.assign(tradeModel, legData);
        await SaveTrade(tradeModel);
      } 
      
      console.log('Save complete, forcing full refresh...');
      resetForm();
      
      // Force a complete refresh to show the new trade
      await forceRefresh();
      
      alert(`Trade ${isEditing ? 'updated' : 'added'} successfully!`);

    } catch (error) {
      console.error('Failed to save trade(s):', error);
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
      strategy: '',
      notes: ''
    };
    additionalExpirations = [];
    showAdditionalExpiration = false;
    selectedStrategyCategory = '';
    selectedStrategyDescription = '';
    isEditing = false;
    editingTradeId = null;
  }
  
  function editTrade(tradeId) {
    const representativeLeg = trades.find(t => t.id === tradeId);
    if (!representativeLeg) return;
    const allLegs = trades.filter(t => t.id === tradeId);
    
    // Ensure dates are proper JS Date objects for formatting
    let repEntryDate = null;
    let repExpDate = null;
    
    // Handle potential date format issues safely
    if (representativeLeg.entryDate) {
      try {
        repEntryDate = new Date(representativeLeg.entryDate);
      } catch (e) {
        console.error("Failed to parse entry date:", representativeLeg.entryDate);
      }
    }
    
    if (representativeLeg.expirationDate) {
      try {
        repExpDate = new Date(representativeLeg.expirationDate);
      } catch (e) {
        console.error("Failed to parse expiration date:", representativeLeg.expirationDate);
      }
    }

    newTrade = {
      id: representativeLeg.id,
      entryDate: repEntryDate ? repEntryDate.toISOString().split('T')[0] : '',
      expirationDate: repExpDate ? repExpDate.toISOString().split('T')[0] : '',
      ticker: representativeLeg.symbol,
      sector: representativeLeg.sector,
      entryPrice: representativeLeg.entryPrice,
      strategy: `${representativeLeg.strategy} - ${representativeLeg.type}`,
      notes: representativeLeg.notes
    };
    
    // Filter and map additional expirations safely
    additionalExpirations = allLegs
      .map(leg => {
        if (!leg.expirationDate) return null;
        try {
          const date = new Date(leg.expirationDate);
          return date.toISOString().split('T')[0];
        } catch (e) {
          console.error("Failed to parse leg expiration date:", leg.expirationDate);
          return null;
        }
      })
      .filter(date => date && date !== newTrade.expirationDate)
      .filter((value, index, self) => self.indexOf(value) === index);
      
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
  
  function addExpirationField() {
    additionalExpirations = [...additionalExpirations, ''];
  }
  
  function removeExpirationField(index) {
    additionalExpirations = additionalExpirations.filter((_, i) => i !== index);
  }
  
  function updateAdditionalExpiration(index, value) {
    additionalExpirations[index] = value;
    additionalExpirations = [...additionalExpirations];
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

  onMount(async () => {
    console.log('TradeCalendar component mounted');
    await loadTrades();
    await tick();
    console.log('After tick, trades length:', trades.length);
    
    // Force DOM update after a short delay as a fallback
    setTimeout(() => {
      console.log('Forcing update after timeout');
      trades = [...trades]; 
    }, 100);
  });
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
          
          {#if showAdditionalExpiration}
            <div class="additional-expirations">
              <div class="additional-header">
                <h4>Additional Expiration Dates</h4>
                <button type="button" class="add-btn" on:click={addExpirationField}>Add Date</button>
              </div>
              
              {#each additionalExpirations as expiration, i}
                <div class="additional-exp-row">
                  <input 
                    type="date" 
                    value={expiration} 
                    on:change={(e) => updateAdditionalExpiration(i, e.currentTarget.value)}
                  />
                  <button type="button" class="remove-btn" on:click={() => removeExpirationField(i)}>
                    ✕
                  </button>
                </div>
              {/each}
              
              {#if additionalExpirations.length === 0}
                <p class="additional-note">For multi-leg strategies spanning different expiration dates</p>
              {/if}
            </div>
          {/if}
          
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
            <button class="reset-btn" on:click={resetForm}>
              {isEditing ? 'Cancel' : 'Reset'}
            </button>
            <button class="add-btn" on:click={addNewTrade}>
              {isEditing ? 'Update Trade' : 'Save Trade'}
            </button>
          </div>
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
  }
  
  .calendar {
    width: 100%;
    border-collapse: collapse;
    border: 1px solid var(--border-color);
  }
  
  .sector-header {
    width: 150px;
    text-align: left;
    background-color: var(--bg-color);
  }
  
  .week-header {
    min-width: 120px;
    background-color: var(--bg-color);
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
    border: 1px solid var(--border-color);
    background-color: var(--bg-color);
  }
  
  .trade-cell {
    border: 1px solid var(--border-color);
    padding: 0.5rem;
    vertical-align: top;
    height: 80px;
    position: relative;
    overflow-y: auto;
    max-height: 120px;
  }
  
  .trade-cell.has-trades {
    background-color: rgba(247, 250, 252, 0.1);
  }
  
  .trade-card {
    padding: 0.5rem;
    border-radius: 4px;
    margin-bottom: 0.5rem;
    color: white;
    font-size: 0.8rem;
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
    margin-bottom: 0.25rem;
  }
  
  .trade-type {
    font-size: 0.75rem;
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
  }
  
  .week-option {
    padding: 0.2rem 0.4rem;
    background-color: var(--bg-color);
    border-radius: 3px;
    font-size: 0.75rem;
  }
  
  .additional-expirations {
    margin-bottom: 1rem;
    padding: 1rem;
    background-color: var(--bg-color);
    border-radius: 4px;
  }
  
  .additional-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
  }
  
  .additional-header h4 {
    margin: 0;
  }
  
  .additional-exp-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 0.5rem;
  }
  
  .additional-exp-row input {
    flex: 1;
  }
  
  .remove-btn {
    background: #e53e3e;
    color: white;
    border: none;
    border-radius: 50%;
    width: 24px;
    height: 24px;
    cursor: pointer;
    font-size: 0.75rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .additional-note {
    color: #cbd5e0;
    font-size: 0.9rem;
    font-style: italic;
    margin: 0.5rem 0 0;
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
</style> 