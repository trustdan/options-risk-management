<script>
  import { onMount } from 'svelte';
  
  // Define journal entry properties
  let journalEntries = [];
  let currentEntry = {
    id: '',
    date: new Date().toISOString().split('T')[0],
    title: '',
    tradeSymbol: '',
    strategy: '',
    outcome: 'win', // win, loss, breakeven
    pnlAmount: 0,
    emotionalState: '',
    lessonsLearned: '',
    whatWentWell: '',
    whatWentPoorly: '',
    improvementPlan: ''
  };
  
  // Strategy types with detailed options
  const strategyOptions = [
    {
      category: 'Basic Spreads',
      options: [
        { name: 'Long Call', description: 'Bullish directional bet - Buy a call option' },
        { name: 'Long Put', description: 'Bearish directional bet - Buy a put option' },
        { name: 'Covered Call', description: 'Income strategy - Own stock and sell calls against it' }
      ]
    },
    {
      category: 'Vertical Spreads',
      options: [
        { name: 'Bull Call Spread', description: 'Bullish with limited risk/reward' },
        { name: 'Bear Call Spread', description: 'Bearish with limited risk/reward' },
        { name: 'Bull Put Spread', description: 'Bullish with limited risk/reward' },
        { name: 'Bear Put Spread', description: 'Bearish with limited risk/reward' }
      ]
    },
    {
      category: 'Calendar/Horizontal Spreads',
      options: [
        { name: 'Long Calendar Call Spread', description: 'Neutral to slightly bullish' },
        { name: 'Long Calendar Put Spread', description: 'Neutral to slightly bearish' }
      ]
    },
    {
      category: 'Diagonal Spreads',
      options: [
        { name: 'Diagonal Call Spread Up', description: 'Moderately bullish' },
        { name: 'Diagonal Call Spread Down', description: 'Bearish to neutral' },
        { name: 'Diagonal Put Spread Up', description: 'Bearish to neutral' },
        { name: 'Diagonal Put Spread Down', description: 'Moderately bullish' }
      ]
    },
    {
      category: 'Butterfly Spreads',
      options: [
        { name: 'Long Call Butterfly', description: 'Neutral, expecting little movement' },
        { name: 'Long Put Butterfly', description: 'Neutral, expecting little movement' },
        { name: 'Broken Wing Butterfly Up', description: 'Bullish, with uneven wings' },
        { name: 'Broken Wing Butterfly Down', description: 'Bearish, with uneven wings' }
      ]
    },
    {
      category: 'Iron Condors/Butterflies',
      options: [
        { name: 'Iron Condor', description: 'Neutral, expecting range-bound movement' },
        { name: 'Iron Butterfly', description: 'Highly neutral, expecting minimal movement' }
      ]
    },
    {
      category: 'Ratio Spreads',
      options: [
        { name: 'Call Ratio Backspread', description: 'Bullish with volatility bias' },
        { name: 'Put Ratio Backspread', description: 'Bearish with volatility bias' }
      ]
    },
    {
      category: 'Danger - naked options ahead',
      options: [
        { name: 'Short Call', description: 'Bearish directional bet - Sell a call option' },
        { name: 'Short Put', description: 'Bullish directional bet - Sell a put option' },
        { name: 'Cash-Secured Put', description: 'Income strategy - Sell puts with cash to secure' },
        { name: 'Call Ratio Spread', description: 'Neutral to moderately bearish - Short 2, long 1' },
        { name: 'Put Ratio Spread', description: 'Neutral to moderately bullish - Short 2, long 1' }
      ]
    },
    {
      category: 'Other',
      options: [
        { name: 'Long Stock', description: 'Direct stock ownership' }
      ]
    }
  ];
  
  // Function to get strategy color based on type
  function getStrategyColor(strategy) {
    // Match colors from the TradeJournal component
    const strategyLower = strategy ? strategy.toLowerCase() : '';
    
    // Danger - naked options ahead category
    if (strategyLower.includes('short call') || 
        strategyLower.includes('short put') ||
        strategyLower.includes('cash-secured put') ||
        strategyLower.includes('call ratio spread') ||
        strategyLower.includes('put ratio spread')) {
      return '#ff0000'; // Bright red for dangerous strategies
    }
    
    // Other strategies
    if (strategyLower.includes('covered call')) {
      return '#4CAF50'; // Green
    } else if (strategyLower.includes('put spread')) {
      return '#9b59b6'; // Purple
    } else if (strategyLower.includes('call spread')) {
      return '#3498db'; // Blue
    } else if (strategyLower.includes('iron condor')) {
      return '#1abc9c'; // Teal
    } else if (strategyLower.includes('long call')) {
      return '#2196F3'; // Blue
    } else if (strategyLower.includes('long put')) {
      return '#f39c12'; // Orange
    } else if (strategyLower.includes('calendar')) {
      return '#2ecc71'; // Green
    } else if (strategyLower.includes('long stock')) {
      return '#4CAF50'; // Green
    } else if (strategyLower.includes('bull') || strategyLower.includes('bear')) {
      return '#9C27B0'; // Purple
    } else if (strategyLower.includes('butterfly')) {
      return '#F44336'; // Red
    } else {
      return '#607D8B'; // Default blue-grey
    }
  }
  
  onMount(() => {
    // Placeholder for API call to get journal entries
    // In a real implementation, this would call the Go backend
    journalEntries = [
      {
        id: '1',
        date: '2023-05-10',
        title: 'AAPL Earnings Play',
        tradeSymbol: 'AAPL',
        strategy: 'Iron Condor',
        outcome: 'win',
        pnlAmount: 250,
        emotionalState: 'Confident but cautious',
        lessonsLearned: 'Patience pays off in earnings plays',
        whatWentWell: 'Proper position sizing and timing',
        whatWentPoorly: 'Could have sized up a bit more',
        improvementPlan: 'Develop a more specific earnings play checklist'
      },
      {
        id: '2',
        date: '2023-05-15',
        title: 'TSLA Overextended',
        tradeSymbol: 'TSLA',
        strategy: 'Bear Call Spread',
        outcome: 'loss',
        pnlAmount: -320,
        emotionalState: 'Overconfident',
        lessonsLearned: 'Respect the momentum in high-volatility stocks',
        whatWentWell: 'Proper stop loss prevented bigger losses',
        whatWentPoorly: 'Entered too early, ignored technical signals',
        improvementPlan: 'Wait for confirmation before entering counter-trend trades'
      }
    ];
  });
  
  function saveEntry() {
    // Placeholder for save functionality
    const isNew = !currentEntry.id;
    
    if (isNew) {
      // Generate a dummy ID for new entries
      currentEntry.id = Date.now().toString();
      journalEntries = [...journalEntries, {...currentEntry}];
    } else {
      // Update existing entry
      journalEntries = journalEntries.map(entry => 
        entry.id === currentEntry.id ? {...currentEntry} : entry
      );
    }
    
    // Reset the form
    resetForm();
  }
  
  function editEntry(entry) {
    currentEntry = {...entry};
  }
  
  function deleteEntry(id) {
    if (confirm('Are you sure you want to delete this journal entry?')) {
      journalEntries = journalEntries.filter(entry => entry.id !== id);
    }
  }
  
  function resetForm() {
    currentEntry = {
      id: '',
      date: new Date().toISOString().split('T')[0],
      title: '',
      tradeSymbol: '',
      strategy: '',
      outcome: 'win',
      pnlAmount: 0,
      emotionalState: '',
      lessonsLearned: '',
      whatWentWell: '',
      whatWentPoorly: '',
      improvementPlan: ''
    };
  }
</script>

<div class="journal-container p-6 bg-gray-100 rounded-lg shadow-lg max-w-5xl mx-auto">
  <h1 class="text-2xl font-bold mb-6 text-blue-800">Trading Journal</h1>
  
  <!-- Journal Entry Form -->
  <div class="bg-white p-5 rounded-md shadow mb-6">
    <h2 class="text-lg font-semibold mb-4 text-gray-700">
      {currentEntry.id ? 'Edit Entry' : 'New Journal Entry'}
    </h2>
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
      <div>
        <label class="block text-sm font-medium text-gray-700">Date</label>
        <input 
          type="date" 
          bind:value={currentEntry.date}
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
        />
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700">Symbol</label>
        <input 
          type="text" 
          bind:value={currentEntry.tradeSymbol}
          placeholder="AAPL, SPY, etc."
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
        />
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700">Strategy</label>
        <select 
          bind:value={currentEntry.strategy}
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
        >
          <option value="">Select a strategy</option>
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
    
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700">Trade Title</label>
      <input 
        type="text" 
        bind:value={currentEntry.title}
        placeholder="Brief description of the trade"
        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
      />
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
      <div>
        <label class="block text-sm font-medium text-gray-700">Outcome</label>
        <div class="mt-2 flex gap-4">
          <label class="inline-flex items-center">
            <input type="radio" bind:group={currentEntry.outcome} value="win" class="form-radio text-green-600">
            <span class="ml-2 text-green-600 font-medium">Win</span>
          </label>
          <label class="inline-flex items-center">
            <input type="radio" bind:group={currentEntry.outcome} value="loss" class="form-radio text-red-600">
            <span class="ml-2 text-red-600 font-medium">Loss</span>
          </label>
          <label class="inline-flex items-center">
            <input type="radio" bind:group={currentEntry.outcome} value="breakeven" class="form-radio text-yellow-600">
            <span class="ml-2 text-yellow-600 font-medium">Breakeven</span>
          </label>
        </div>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700">P&L Amount ($)</label>
        <input 
          type="number" 
          bind:value={currentEntry.pnlAmount}
          placeholder="$ Profit or Loss"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
        />
      </div>
    </div>
    
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700">Emotional State</label>
      <input 
        type="text" 
        bind:value={currentEntry.emotionalState}
        placeholder="How were you feeling before/during the trade?"
        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
      />
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
      <div>
        <label class="block text-sm font-medium text-gray-700">What Went Well</label>
        <textarea 
          bind:value={currentEntry.whatWentWell}
          rows="3"
          placeholder="Aspects of the trade that were executed properly"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
        ></textarea>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700">What Went Poorly</label>
        <textarea 
          bind:value={currentEntry.whatWentPoorly}
          rows="3"
          placeholder="Mistakes or areas for improvement"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
        ></textarea>
      </div>
    </div>
    
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700">Lessons Learned</label>
      <textarea 
        bind:value={currentEntry.lessonsLearned}
        rows="3"
        placeholder="Key takeaways from this trade"
        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
      ></textarea>
    </div>
    
    <div class="mb-4">
      <label class="block text-sm font-medium text-gray-700">Improvement Plan</label>
      <textarea 
        bind:value={currentEntry.improvementPlan}
        rows="3"
        placeholder="Specific actions to improve future trades"
        class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
      ></textarea>
    </div>
    
    <div class="flex justify-end gap-2 mt-6">
      <button 
        on:click={resetForm}
        class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300"
      >
        Clear
      </button>
      <button 
        on:click={saveEntry}
        class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
      >
        {currentEntry.id ? 'Update Entry' : 'Save Entry'}
      </button>
    </div>
  </div>
  
  <!-- Journal Entries List -->
  <div class="bg-white p-5 rounded-md shadow">
    <h2 class="text-lg font-semibold mb-4 text-gray-700">Recent Journal Entries</h2>
    
    {#if journalEntries.length === 0}
      <p class="text-gray-500 italic">No journal entries found. Create your first entry above.</p>
    {:else}
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Symbol</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Strategy</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">P&L</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            {#each journalEntries as entry}
              <tr>
                <td class="px-6 py-4 whitespace-nowrap">{entry.date}</td>
                <td class="px-6 py-4 whitespace-nowrap font-medium">{entry.tradeSymbol}</td>
                <td class="px-6 py-4">{entry.title}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="strategy-pill" style="background-color: {getStrategyColor(entry.strategy)}">
                    {entry.strategy}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class={entry.outcome === 'win' ? 'text-green-600' : entry.outcome === 'loss' ? 'text-red-600' : 'text-yellow-600'}>
                    ${entry.pnlAmount}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button 
                    on:click={() => editEntry(entry)}
                    class="text-indigo-600 hover:text-indigo-900 mr-3"
                  >
                    Edit
                  </button>
                  <button 
                    on:click={() => deleteEntry(entry.id)}
                    class="text-red-600 hover:text-red-900"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
</div>

<style>
  /* Additional styling if needed beyond the Tailwind classes */
  .journal-container {
    min-height: calc(100vh - 180px);
  }
  
  .strategy-pill {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: 9999px;
    font-size: 0.8rem;
    font-weight: 500;
    color: white;
  }
</style> 