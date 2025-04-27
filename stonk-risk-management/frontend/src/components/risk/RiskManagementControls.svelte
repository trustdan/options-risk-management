<!-- RiskManagementControls.svelte -->
<script>
  import { maxDollarRiskPerTrade as maxRiskStore } from '../../store/riskStore.js';
  
  // Account info
  let accountValue = 25000;
  
  // Core risk parameters
  let accountRiskPerTrade = 2;
  let maxPortfolioExposure = 20;
  let stopLossPercent = 50;
  let riskRewardRatio = 2;
  
  // Drawdown circuit breakers
  let dailyLossLimit = 3;
  let weeklyLossLimit = 7;
  
  // Advanced parameters
  let positionScaling = 100;
  let correlationAdjustment = 1;
  let volatilityMultiplier = 1;
  let maxDrawdownTolerance = 15;
  
  // Calculated values based on general parameters
  $: {
    let calcMaxDollarRiskPerTrade = accountValue * (accountRiskPerTrade / 100);
    maxRiskStore.set(calcMaxDollarRiskPerTrade);
  }
  $: maxDollarRiskPerTrade = accountValue * (accountRiskPerTrade / 100);
  $: maxPortfolioExposureDollars = accountValue * (maxPortfolioExposure / 100);
  $: stopLossDollarAmount = maxDollarRiskPerTrade * (stopLossPercent / 100);
  $: potentialRewardDollarAmount = stopLossDollarAmount * riskRewardRatio;
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
      <div class="flex items-center space-x-2 mt-1">
        <span class="text-xs text-gray-500 w-8 text-right">0.5%</span>
        <input
          type="range"
          min="0.5"
          max="5"
          step="0.25"
          bind:value={accountRiskPerTrade}
          class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer flex-grow"
        />
        <span class="text-xs text-gray-500 w-8">5%</span>
      </div>
    </div>
    
    <div class="mb-4">
      <div class="flex justify-between">
        <label class="block text-sm font-medium text-gray-700">Maximum Portfolio Exposure: {maxPortfolioExposure}%</label>
        <span class="text-sm text-gray-500">${(accountValue * maxPortfolioExposure / 100).toFixed(2)}</span>
      </div>
      <div class="flex items-center space-x-2 mt-1">
        <span class="text-xs text-gray-500 w-8 text-right">10%</span>
        <input
          type="range"
          min="10"
          max="100"
          step="5"
          bind:value={maxPortfolioExposure}
          class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer flex-grow"
        />
        <span class="text-xs text-gray-500 w-8">100%</span>
      </div>
    </div>
  </div>
  
  <!-- Stop Loss as % of Per-Trade Risk -->
  <div class="bg-white p-4 rounded-md shadow mb-6">
    <h2 class="text-lg font-semibold mb-3 text-gray-700">Stop Loss as % of Per-Trade Risk</h2>
    
    <div class="mb-4">
      <div class="flex justify-between">
        <label class="block text-sm font-medium text-gray-700">Stop Loss as % of Per-Trade Risk: {stopLossPercent}%</label>
        <span class="text-sm text-gray-500">${stopLossDollarAmount.toFixed(2)} Max Loss</span>
      </div>
      <!-- Min 10% of allowed risk -->
      <!-- Max 100% of allowed risk -->
      <div class="flex items-center space-x-2 mt-1">
        <span class="text-xs text-gray-500 w-8 text-right">10%</span>
        <input
          type="range"
          min="10"
          max="100"
          step="5"
          bind:value={stopLossPercent}
          class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer flex-grow"
        />
        <span class="text-xs text-gray-500 w-8">100%</span>
      </div>
    </div>
  </div>
  
  <!-- Risk:Reward Ratio -->
  <div class="bg-white p-4 rounded-md shadow mb-6">
    <h2 class="text-lg font-semibold mb-3 text-gray-700">Target Risk:Reward Ratio</h2>
    
    <div class="mb-4">
      <div class="flex justify-between">
        <label class="block text-sm font-medium text-gray-700">Target Risk:Reward Ratio: 1:{riskRewardRatio}</label>
        <!-- Displaying R:R comparison below -->
      </div>
      <!-- Min 1:1 R:R -->
      <!-- Max 1:5 R:R -->
      <div class="flex items-center space-x-2 mt-1">
        <span class="text-xs text-gray-500 w-8 text-right">1:1</span>
        <input
          type="range"
          min="1"
          max="5"
          step="0.25"
          bind:value={riskRewardRatio}
          class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer flex-grow"
        />
        <span class="text-xs text-gray-500 w-8">1:5</span>
      </div>
      <!-- Dollar R:R display -->
      <div class="mt-2 p-2 bg-gray-50 rounded text-center">
        <span class="text-sm font-medium">
          <span class="text-red-500">${maxDollarRiskPerTrade.toFixed(2)}</span> : <span class="text-green-500">${(maxDollarRiskPerTrade * riskRewardRatio).toFixed(2)}</span>
        </span>
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
      <div class="flex items-center space-x-2 mt-1">
        <span class="text-xs text-gray-500 w-8 text-right">1%</span>
        <input
          type="range"
          min="1"
          max="10"
          step="0.5"
          bind:value={dailyLossLimit}
          class="w-full h-2 bg-red-100 rounded-lg appearance-none cursor-pointer flex-grow"
        />
        <span class="text-xs text-gray-500 w-8">10%</span>
      </div>
    </div>
    
    <div class="mb-4">
      <div class="flex justify-between">
        <label class="block text-sm font-medium text-gray-700">Weekly Loss Limit: {weeklyLossLimit}%</label>
        <span class="text-sm text-gray-500">${(accountValue * weeklyLossLimit / 100).toFixed(2)}</span>
      </div>
      <div class="flex items-center space-x-2 mt-1">
        <span class="text-xs text-gray-500 w-8 text-right">2%</span>
        <input
          type="range"
          min="2"
          max="20"
          step="1"
          bind:value={weeklyLossLimit}
          class="w-full h-2 bg-red-100 rounded-lg appearance-none cursor-pointer flex-grow"
        />
        <span class="text-xs text-gray-500 w-8">20%</span>
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
        <div class="flex items-center space-x-2 mt-1">
          <span class="text-xs text-gray-500 w-8 text-right">25%</span>
          <input
            type="range"
            min="25"
            max="100"
            step="25"
            bind:value={positionScaling}
            class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer flex-grow"
          />
          <span class="text-xs text-gray-500 w-8">100%</span>
        </div>
      </div>
      
      <div class="mb-4">
        <div class="flex justify-between">
          <label class="block text-sm font-medium text-gray-700">Correlation Adjustment: {correlationAdjustment}x</label>
        </div>
        <div class="flex items-center space-x-2 mt-1">
          <span class="text-xs text-gray-500 w-8 text-right">0.5x</span>
          <input
            type="range"
            min="0.5"
            max="1"
            step="0.1"
            bind:value={correlationAdjustment}
            class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer flex-grow"
          />
          <span class="text-xs text-gray-500 w-8">1x</span>
        </div>
      </div>
      
      <div class="mb-4">
        <div class="flex justify-between">
          <label class="block text-sm font-medium text-gray-700">Volatility Multiplier: {volatilityMultiplier}x</label>
        </div>
        <div class="flex items-center space-x-2 mt-1">
          <span class="text-xs text-gray-500 w-8 text-right">0.5x</span>
          <input
            type="range"
            min="0.5"
            max="3"
            step="0.25"
            bind:value={volatilityMultiplier}
            class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer flex-grow"
          />
          <span class="text-xs text-gray-500 w-8">3x</span>
        </div>
      </div>
      
      <div class="mb-4">
        <div class="flex justify-between">
          <label class="block text-sm font-medium text-gray-700">Max Drawdown Tolerance: {maxDrawdownTolerance}%</label>
        </div>
        <div class="flex items-center space-x-2 mt-1">
          <span class="text-xs text-gray-500 w-8 text-right">5%</span>
          <input
            type="range"
            min="5"
            max="25"
            step="1"
            bind:value={maxDrawdownTolerance}
            class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer flex-grow"
          />
          <span class="text-xs text-gray-500 w-8">25%</span>
        </div>
      </div>
    </details>
  </div>
  
  <!-- General Risk Calculation -->
  <div class="bg-blue-50 p-4 rounded-md shadow border-l-4 border-blue-500">
    <h2 class="text-lg font-semibold mb-3 text-blue-800">General Risk Metrics</h2>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="bg-white p-3 rounded shadow-sm">
        <span class="text-sm font-medium text-red-500">Max $ Risk Per Trade</span>
        <p class="text-xl font-bold text-red-600">${maxDollarRiskPerTrade.toFixed(2)}</p>
      </div>
      <div class="bg-white p-3 rounded shadow-sm">
        <span class="text-sm font-medium text-green-500">Target $ Reward Per Trade</span>
        <p class="text-xl font-bold text-green-600">${(maxDollarRiskPerTrade * riskRewardRatio).toFixed(2)}</p>
        <span class="text-xs text-gray-500">(Based on 1:{riskRewardRatio} R:R)</span>
      </div>
      <div class="bg-white p-3 rounded shadow-sm col-span-1 md:col-span-2">
        <span class="text-sm font-medium text-gray-500">Max Portfolio Exposure</span>
        <p class="text-xl font-bold">${maxPortfolioExposureDollars.toFixed(2)}</p>
        <span class="text-xs text-gray-500">({maxPortfolioExposure}% of Account Value)</span>
      </div>
    </div>
  </div>
</div> 