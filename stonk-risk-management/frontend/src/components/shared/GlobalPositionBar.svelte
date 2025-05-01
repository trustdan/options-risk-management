<script>
  import { onMount, onDestroy } from 'svelte';
  import { maxDollarRiskPerTrade } from '../../store/riskStore.js';
  
  // Default position size (will be updated by RiskDashboard)
  export let positionSize = 80;
  
  // Create a custom event system to communicate between components
  let positionUpdateCallback;
  
  onMount(() => {
    // Set up listener for position updates
    window.addEventListener('position-size-updated', handlePositionUpdate);
    
    // Register this component to receive updates
    // @ts-ignore - Custom property on window object
    if (typeof window.__positionBarRegister === 'function') {
      // @ts-ignore - Custom property on window object
      window.__positionBarRegister(updatePositionSize);
    } else {
      // Create the registration function if it doesn't exist
      // @ts-ignore - Custom property on window object
      window.__positionBarRegister = (callback) => {
        positionUpdateCallback = callback;
      };
    }
  });
  
  onDestroy(() => {
    window.removeEventListener('position-size-updated', handlePositionUpdate);
  });
  
  // Handler for position update events
  function handlePositionUpdate(event) {
    if (event.detail && typeof event.detail.positionSize === 'number') {
      positionSize = event.detail.positionSize;
    }
  }
  
  // Function that will be called by RiskDashboard to update position size
  function updatePositionSize(newSize) {
    positionSize = newSize;
  }
</script>

<div class="global-position-bar">
  <div class="position-label">Recommended Position Size:</div>
  <div class="position-bar">
    <div 
      class="position-indicator" 
      style="width: {positionSize}%; background-color: {positionSize >= 70 ? '#68D391' : positionSize <= 30 ? '#FC8181' : '#F6AD55'}"
    ></div>
  </div>
  <div class="position-value">{positionSize}% of max</div>
  <div class="position-value-dollars">
    ${Math.round((positionSize/100) * $maxDollarRiskPerTrade).toFixed(2)} recommended
  </div>
</div>

<style>
  .global-position-bar {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    background-color: var(--card-bg);
    padding: 4px 80px 4px 16px; /* Added right padding to avoid GitHub octocat */
    display: flex;
    align-items: center;
    gap: 10px;
    z-index: 9999; /* Ensure it's above everything else */
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    font-size: 0.9rem;
    height: 32px;
  }
  
  .position-label {
    white-space: nowrap;
    font-weight: bold;
    flex-shrink: 0;
  }
  
  .position-bar {
    flex: 1;
    height: 20px;
    background-color: rgba(0, 0, 0, 0.1);
    border-radius: 10px;
    overflow: hidden;
    max-width: 300px;
  }
  
  .position-indicator {
    height: 100%;
    border-radius: 10px;
    transition: width 0.5s, background-color 0.5s;
  }
  
  .position-value {
    font-weight: bold;
    white-space: nowrap;
    flex-shrink: 0;
  }
  
  .position-value-dollars {
    margin-left: auto;
    font-weight: bold;
    white-space: nowrap;
    flex-shrink: 0;
  }
  
  @media (max-width: 768px) {
    .global-position-bar {
      flex-wrap: wrap;
      height: auto;
      padding: 4px 80px 4px 10px;
      font-size: 0.8rem;
    }
    
    .position-value-dollars {
      margin-left: 0;
    }
  }
</style> 