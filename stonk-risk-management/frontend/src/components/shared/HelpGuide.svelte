<script>
  import { onMount } from 'svelte';
  
  export let visible = false;
  
  function closeHelp() {
    visible = false;
  }
  
  function handleKeydown(event) {
    if (event.key === 'Escape' && visible) {
      closeHelp();
    }
  }
  
  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    
    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  });
</script>

{#if visible}
  <div class="help-guide">
    <button class="close-button" on:click={closeHelp}>×</button>
    <div class="help-content">
      <h2>Options Trading Dashboard User Guide</h2>
      
      <p class="help-intro">
        This guide explains the features and functionality of the Options Trading Dashboard, designed to help you manage risk, 
        assess market conditions, track trades, and maintain a healthy trading mindset.
      </p>
      
      <div class="help-section">
        <h3>Risk Management Dashboard</h3>
        <p>The Risk Management Dashboard helps you assess your psychological state and guides your position sizing.</p>
        
        <h4>Key Features:</h4>
        <ul>
          <li><strong>Daily Assessment</strong> - Rate your emotional state, FOMO level, market bias, physical condition, P&L impact, and general mindset on a scale from -3 to +3.</li>
          <li><strong>Position Size Recommendation</strong> - Receive guidance on appropriate position sizing based on your psychological assessment.</li>
          <li><strong>Warning Flags</strong> - Get alerted to potential euphoria, chasing behavior, or high FOMO levels.</li>
          <li><strong>Analytics</strong> - View trends in your psychological states over time.</li>
          <li><strong>Position Sizing Controls</strong> - Configure max position size and risk parameters.</li>
        </ul>
        
        <p><strong>Best Practices:</strong> Complete your assessment at the beginning of each trading day before placing any trades.</p>
      </div>
      
      <div class="help-section">
        <h3>Market, Sector, Stock Rating</h3>
        <p>This dashboard helps you track and analyze market conditions, sector performance, and individual stocks.</p>
        
        <h4>Key Features:</h4>
        <ul>
          <li><strong>Market Rating</strong> - Rate overall market conditions on a scale from -3 (bearish) to +3 (bullish).</li>
          <li><strong>Sector Ratings</strong> - Assign sentiment scores to different market sectors.</li>
          <li><strong>Stock Rating</strong> - Track individual stocks with ratings, price targets, and confidence levels.</li>
          <li><strong>Chart Patterns</strong> - Document chart patterns you've identified for specific stocks.</li>
          <li><strong>Historical Ratings</strong> - Navigate through past ratings to track your market analysis over time.</li>
        </ul>
        
        <p><strong>Best Practices:</strong> Update market and sector ratings regularly, especially after significant market events or data releases.</p>
      </div>
      
      <div class="help-section">
        <h3>Trades, Calendar, Journal</h3>
        <p>This feature allows you to plan, document, and reflect on your trades across a multi-week calendar view.</p>
        
        <h4>Key Features:</h4>
        <ul>
          <li><strong>Trade Calendar</strong> - Visualize your trades across an 8-week view, organized by sector.</li>
          <li><strong>Add Trade</strong> - Document trades with detailed information including strategy, entry/exit prices, and expiration dates.</li>
          <li><strong>Strategy Selection</strong> - Choose from a comprehensive list of options strategies with descriptions.</li>
          <li><strong>Trade Journal</strong> - Keep detailed notes about your trades, including emotional context and lessons learned.</li>
          <li><strong>Multi-leg Strategies</strong> - Document calendar and diagonal spreads with different expiration dates.</li>
        </ul>
        
        <p><strong>Best Practices:</strong> Enter trades as you place them, and review the calendar regularly to ensure balanced exposure across sectors and expiration dates.</p>
      </div>
      
      <div class="help-section">
        <h3>Trading Koans</h3>
        <p>A collection of trading wisdom to help improve your trading psychology and decision-making.</p>
        
        <h4>Key Features:</h4>
        <ul>
          <li><strong>Mindset and Psychology</strong> - Wisdom focused on developing a healthy trading mindset.</li>
          <li><strong>Risk Management</strong> - Principles for preserving capital and managing risk effectively.</li>
          <li><strong>Market Behavior</strong> - Insights into how markets typically function.</li>
          <li><strong>Trading Process</strong> - Guidance on developing and following a systematic approach.</li>
          <li><strong>Self-Discipline</strong> - Reminders about the importance of discipline in trading.</li>
          <li><strong>Managing Tilt</strong> - Strategies for recognizing and addressing emotional reactions.</li>
        </ul>
        
        <p><strong>Best Practices:</strong> Review these koans regularly, especially after difficult trading days or losses.</p>
      </div>
      
      <div class="help-section">
        <h3>Other Features</h3>
        
        <h4>Data Management:</h4>
        <ul>
          <li><strong>Export Data</strong> - Back up your trading data to your local filesystem.</li>
          <li><strong>Import Data</strong> - Restore previously exported data.</li>
          <li><strong>Database Cleanup</strong> - Perform maintenance to optimize database performance.</li>
        </ul>
        
        <h4>Navigation and Settings:</h4>
        <ul>
          <li><strong>Theme Toggle</strong> - Switch between light and dark modes (☀️/🌙).</li>
          <li><strong>Keyboard Shortcuts</strong> - Access VIM-style keyboard navigation with the ⌨️ button or by pressing ?.</li>
          <li><strong>Refresh</strong> - Update data displays with the latest information.</li>
          <li><strong>Feedback</strong> - Provide feedback to the developers.</li>
        </ul>
      </div>
      
      <div class="help-note">
        <p>Remember: This dashboard is designed to help you take control of your trading psychology, manage risk effectively, and make more informed trading decisions. Use it consistently for the best results.</p>
      </div>
    </div>
  </div>
{/if}

<style>
  .help-guide {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: var(--card-bg, white);
    border: 1px solid var(--border-color, #e2e8f0);
    border-radius: 6px;
    padding: 25px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 9999;
    width: 90%;
    max-width: 800px;
    max-height: 85vh;
    overflow-y: auto;
    animation: fadeIn 0.3s ease-out;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
  
  .close-button {
    position: absolute;
    top: 15px;
    right: 15px;
    background: transparent;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: var(--text-muted, #718096);
    padding: 0;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
  }
  
  .close-button:hover {
    background: rgba(0, 0, 0, 0.05);
    color: var(--text-color, #333);
  }
  
  h2 {
    margin: 0 0 15px 0;
    font-size: 24px;
    color: var(--heading-color, #2d3748);
    text-align: center;
    border-bottom: 1px solid var(--border-color, #e2e8f0);
    padding-bottom: 15px;
  }
  
  h3 {
    margin: 20px 0 10px 0;
    font-size: 20px;
    color: var(--heading-color, #2d3748);
  }
  
  h4 {
    margin: 15px 0 8px 0;
    font-size: 16px;
    color: var(--heading-color, #2d3748);
  }
  
  .help-intro {
    font-size: 16px;
    line-height: 1.5;
    margin-bottom: 20px;
  }
  
  .help-section {
    margin-bottom: 25px;
    padding-bottom: 15px;
    border-bottom: 1px solid var(--border-color, #e2e8f0);
  }
  
  .help-section:last-child {
    border-bottom: none;
  }
  
  ul {
    margin: 8px 0;
    padding-left: 25px;
  }
  
  li {
    margin-bottom: 8px;
    line-height: 1.4;
  }
  
  p {
    margin: 10px 0;
    line-height: 1.5;
  }
  
  .help-note {
    margin-top: 20px;
    padding: 15px;
    background-color: var(--hover-color, rgba(0,0,0,0.05));
    border-radius: 6px;
    font-style: italic;
  }
  
  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
  
  /* Add responsive styles */
  @media (max-width: 600px) {
    .help-guide {
      padding: 15px;
      width: 95%;
    }
    
    h2 {
      font-size: 20px;
    }
    
    h3 {
      font-size: 18px;
    }
    
    .help-intro, p, li {
      font-size: 14px;
    }
  }
</style> 