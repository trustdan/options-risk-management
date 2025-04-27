<!-- RiskManagementControls.svelte -->
<script>
  // Account info
  let accountValue = 10000;
  let symbol = '';
  let entryPrice = 0;
  
  // Core risk parameters
  let accountRiskPerTrade = 1;
  let maxPortfolioExposure = 60;
  let stopLossPercent = 5;
  let profitTargetRatio = 2;
  
  // Drawdown circuit breakers
  let dailyLossLimit = 3;
  let weeklyLossLimit = 7;
  
  // Advanced parameters
  let positionScaling = 100;
  let correlationAdjustment = 1;
  let volatilityMultiplier = 1;
  let maxDrawdownTolerance = 15;
  
  // Calculated values
  let positionSize = 0;
  let dollarRisk = 0;
  let potentialProfit = 0;
  let stopLossPrice = 0;
  let targetPrice = 0;
  
  // Calculate position size whenever inputs change
  $: if (entryPrice > 0 && stopLossPercent > 0) {
    const rawDollarRisk = accountValue * (accountRiskPerTrade / 100);
    const adjustedDollarRisk = rawDollarRisk * correlationAdjustment * volatilityMultiplier;
    const dollarPerShare = entryPrice * (stopLossPercent / 100);
    const rawPositionSize = adjustedDollarRisk / dollarPerShare;
    const finalPositionSize = Math.floor(rawPositionSize * (positionScaling / 100));
    
    const calculatedStopPrice = entryPrice * (1 - stopLossPercent/100);
    const calculatedTargetPrice = entryPrice + (entryPrice - calculatedStopPrice) * profitTargetRatio;
    
    positionSize = finalPositionSize;
    dollarRisk = finalPositionSize * (entryPrice - calculatedStopPrice);
    potentialProfit = finalPositionSize * (calculatedTargetPrice - entryPrice);
    stopLossPrice = calculatedStopPrice.toFixed(2);
    targetPrice = calculatedTargetPrice.toFixed(2);
  }
</script>

<div class="p-6 bg-gray-100 rounded-lg shadow-lg max-w-4xl mx-auto">
  <h1 class="text-2xl font-bold mb-6 text-blue-800">Risk Management Controls</h1>
  
  <!-- Account Information -->
  <div class="bg-white p-4 rounded-md shadow mb-6">
    <h2 class="text-lg font-semibold mb-3 text-gray-700">Account Information</h2>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div>
        <label class="block text-sm font-medium text-gray-700">Account Value ($)</label>
        <input
          type="number"
          bind:value={accountValue}
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700">Symbol</label>
        <input
          type="text"
          bind:value={symbol}
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
        />
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700">Entry Price ($)</label>
        <input
          type="number"
          bind:value={entryPrice}
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
        />
      </div>
    </div>
  </div>
  
  <!-- Core Risk Parameters -->
  <div class="bg-white p-4 rounded-md shadow mb-6">
    <h2 class="text-lg font-semibold mb-3 text-gray-700">Core Risk Parameters</h2>
    
    <div class="mb-4">
      <div class="flex justify-between">
        <label class="block text-sm font-medium text-gray-700">Account Risk Per Trade: {accountRiskPerTrade}%</label>
        <span class="text-sm text-gray-500">${(accountValue * accountRiskPerTrade / 100).toFixed(2)}</span>
      </div>
      <input
        type="range"
        min="0.5"
        max="5"
        step="0.25"
        bind:value={accountRiskPerTrade}
        class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
      />
      <div class="flex justify-between text-xs text-gray-500">
        <span>0.5%</span>
        <span>5%</span>
      </div>
    </div>
    
    <div class="mb-4">
      <div class="flex justify-between">
        <label class="block text-sm font-medium text-gray-700">Maximum Portfolio Exposure: {maxPortfolioExposure}%</label>
        <span class="text-sm text-gray-500">${(accountValue * maxPortfolioExposure / 100).toFixed(2)}</span>
      </div>
      <input
        type="range"
        min="10"
        max="100"
        step="5"
        bind:value={maxPortfolioExposure}
        class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
      />
      <div class="flex justify-between text-xs text-gray-500">
        <span>10%</span>
        <span>100%</span>
      </div>
    </div>
    
    <div class="mb-4">
      <div class="flex justify-between">
        <label class="block text-sm font-medium text-gray-700">Stop Loss Distance: {stopLossPercent}%</label>
        <span class="text-sm text-gray-500">${stopLossPrice}</span>
      </div>
      <input
        type="range"
        min="1"
        max="20"
        step="0.5"
        bind:value={stopLossPercent}
        class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
      />
      <div class="flex justify-between text-xs text-gray-500">
        <span>1%</span>
        <span>20%</span>
      </div>
    </div>
    
    <div class="mb-4">
      <div class="flex justify-between">
        <label class="block text-sm font-medium text-gray-700">Profit Target Ratio: {profitTargetRatio}:1</label>
        <span class="text-sm text-gray-500">${targetPrice}</span>
      </div>
      <input
        type="range"
        min="1"
        max="5"
        step="0.25"
        bind:value={profitTargetRatio}
        class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
      />
      <div class="flex justify-between text-xs text-gray-500">
        <span>1:1</span>
        <span>5:1</span>
      </div>
    </div>
  </div>
  
  <!-- Drawdown Circuit Breakers -->
  <div class="bg-white p-4 rounded-md shadow mb-6 border-l-4 border-red-500">
    <h2 class="text-lg font-semibold mb-3 text-gray-700">Drawdown Circuit Breakers</h2>
    
    <div class="mb-4">
      <div class="flex justify-between">
        <label class="block text-sm font-medium text-gray-700">Daily Loss Limit: {dailyLossLimit}%</label>
        <span class="text-sm text-gray-500">${(accountValue * dailyLossLimit / 100).toFixed(2)}</span>
      </div>
      <input
        type="range"
        min="1"
        max="10"
        step="0.5"
        bind:value={dailyLossLimit}
        class="w-full h-2 bg-red-100 rounded-lg appearance-none cursor-pointer"
      />
      <div class="flex justify-between text-xs text-gray-500">
        <span>1%</span>
        <span>10%</span>
      </div>
    </div>
    
    <div class="mb-4">
      <div class="flex justify-between">
        <label class="block text-sm font-medium text-gray-700">Weekly Loss Limit: {weeklyLossLimit}%</label>
        <span class="text-sm text-gray-500">${(accountValue * weeklyLossLimit / 100).toFixed(2)}</span>
      </div>
      <input
        type="range"
        min="2"
        max="20"
        step="1"
        bind:value={weeklyLossLimit}
        class="w-full h-2 bg-red-100 rounded-lg appearance-none cursor-pointer"
      />
      <div class="flex justify-between text-xs text-gray-500">
        <span>2%</span>
        <span>20%</span>
      </div>
    </div>
  </div>
  
  <!-- Advanced Parameters -->
  <div class="bg-white p-4 rounded-md shadow mb-6">
    <details>
      <summary class="text-lg font-semibold mb-3 text-gray-700 cursor-pointer">Advanced Parameters</summary>
      
      <div class="mb-4">
        <div class="flex justify-between">
          <label class="block text-sm font-medium text-gray-700">Position Scaling: {positionScaling}%</label>
        </div>
        <input
          type="range"
          min="25"
          max="100"
          step="25"
          bind:value={positionScaling}
          class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
        />
        <div class="flex justify-between text-xs text-gray-500">
          <span>25%</span>
          <span>100%</span>
        </div>
      </div>
      
      <div class="mb-4">
        <div class="flex justify-between">
          <label class="block text-sm font-medium text-gray-700">Correlation Adjustment: {correlationAdjustment}x</label>
        </div>
        <input
          type="range"
          min="0.5"
          max="1"
          step="0.1"
          bind:value={correlationAdjustment}
          class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
        />
        <div class="flex justify-between text-xs text-gray-500">
          <span>0.5x</span>
          <span>1x</span>
        </div>
      </div>
      
      <div class="mb-4">
        <div class="flex justify-between">
          <label class="block text-sm font-medium text-gray-700">Volatility Multiplier: {volatilityMultiplier}x</label>
        </div>
        <input
          type="range"
          min="0.5"
          max="3"
          step="0.25"
          bind:value={volatilityMultiplier}
          class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
        />
        <div class="flex justify-between text-xs text-gray-500">
          <span>0.5x</span>
          <span>3x</span>
        </div>
      </div>
      
      <div class="mb-4">
        <div class="flex justify-between">
          <label class="block text-sm font-medium text-gray-700">Max Drawdown Tolerance: {maxDrawdownTolerance}%</label>
        </div>
        <input
          type="range"
          min="5"
          max="25"
          step="1"
          bind:value={maxDrawdownTolerance}
          class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
        />
        <div class="flex justify-between text-xs text-gray-500">
          <span>5%</span>
          <span>25%</span>
        </div>
      </div>
    </details>
  </div>
  
  <!-- Position Size Calculation -->
  <div class="bg-blue-50 p-4 rounded-md shadow border-l-4 border-blue-500">
    <h2 class="text-lg font-semibold mb-3 text-blue-800">Position Calculation</h2>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="bg-white p-3 rounded shadow-sm">
        <span class="text-sm font-medium text-gray-500">Position Size</span>
        <p class="text-xl font-bold">{positionSize} shares</p>
      </div>
      <div class="bg-white p-3 rounded shadow-sm">
        <span class="text-sm font-medium text-gray-500">Position Value</span>
        <p class="text-xl font-bold">${(positionSize * entryPrice).toFixed(2)}</p>
      </div>
      <div class="bg-white p-3 rounded shadow-sm">
        <span class="text-sm font-medium text-red-500">Dollar Risk</span>
        <p class="text-xl font-bold text-red-600">${dollarRisk.toFixed(2)}</p>
      </div>
      <div class="bg-white p-3 rounded shadow-sm">
        <span class="text-sm font-medium text-green-500">Potential Profit</span>
        <p class="text-xl font-bold text-green-600">${potentialProfit.toFixed(2)}</p>
      </div>
    </div>
  </div>
</div> 