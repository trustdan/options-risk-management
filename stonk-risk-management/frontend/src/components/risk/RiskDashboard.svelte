<script>
  import { onMount } from 'svelte';
  import { GetRiskAssessments, SaveRiskAssessment } from '../../../wailsjs/go/main/App';
  
  let assessment = {
    id: '',
    date: new Date().toISOString().split('T')[0],
    emotionalScore: 0,
    fomoScore: 0,
    biasScore: 0,
    physicalScore: 0,
    plImpactScore: 0,
    otherScore: 0,
    overallScore: 0,
    notes: ''
  };
  
  let positionSize = 50;
  let dateString = new Date().toISOString().split('T')[0];
  let positionAdvice = {
    title: "Standard position sizing recommended with care.",
    tips: [
      "Trade with awareness of your current psychological state",
      "Consider standard position sizing",
      "Focus on quality setups"
    ]
  };
  let showEuphoriaFlag = false;
  let showChasingFlag = false;
  let showFOMOFlag = false;
  
  // For navigation
  let assessments = [];
  let currentIndex = 0;
  
  onMount(async () => {
    try {
      assessments = await GetRiskAssessments();
      if (assessments && assessments.length > 0) {
        // Load the most recent assessment
        currentIndex = assessments.length - 1;
        loadAssessmentAt(currentIndex);
      }
    } catch (error) {
      console.error('Failed to load risk assessments:', error);
    }
  });
  
  function loadAssessmentAt(index) {
    if (index >= 0 && index < assessments.length) {
      const latestAssessment = assessments[index];
      currentIndex = index;
      
      // Make sure we include all required properties with defaults if missing
      assessment = {
        id: latestAssessment.id || '',
        date: new Date(latestAssessment.date).toISOString().split('T')[0],
        emotionalScore: latestAssessment.emotionalScore || 0,
        fomoScore: latestAssessment.fomoScore || 0,
        biasScore: latestAssessment.biasScore || 0,
        physicalScore: latestAssessment.physicalScore || 0,
        plImpactScore: latestAssessment.plImpactScore || 0,
        otherScore: latestAssessment.otherScore || 0,
        overallScore: latestAssessment.overallScore || 0,
        notes: latestAssessment.notes || ''
      };
      
      calculateRecommendedSize();
    }
  }
  
  function navigatePrevious() {
    if (currentIndex > 0) {
      loadAssessmentAt(currentIndex - 1);
    }
  }
  
  function navigateNext() {
    if (currentIndex < assessments.length - 1) {
      loadAssessmentAt(currentIndex + 1);
    }
  }
  
  function isFirstAssessment() {
    return currentIndex === 0;
  }
  
  function isLastAssessment() {
    return currentIndex === assessments.length - 1;
  }
  
  function calculateRecommendedSize() {
    // Calculate individual factor contributions to position size
    
    // Check for warning flags
    showEuphoriaFlag = assessment.emotionalScore >= 3 || assessment.plImpactScore >= 3;
    showChasingFlag = assessment.plImpactScore <= -3;
    showFOMOFlag = assessment.fomoScore >= 3;
    
    // Emotional state: -3 (more caution), peaks at +1/+2 (less caution), +3 (more caution due to euphoria)
    let emotionalFactor;
    if (assessment.emotionalScore <= -3) {
      emotionalFactor = -30; // Max caution for very negative
    } else if (assessment.emotionalScore <= 0) {
      emotionalFactor = assessment.emotionalScore * 10; // Linear increase
    } else if (assessment.emotionalScore <= 2) {
      emotionalFactor = assessment.emotionalScore * 15; // Peak positive effect
    } else {
      emotionalFactor = 30 - (assessment.emotionalScore - 2) * 30; // Decrease for euphoria
    }
    
    // FOMO level: -3 (peak performance, less caution), +3 (bad, more caution)
    let fomoFactor;
    if (assessment.fomoScore <= -3) {
      fomoFactor = 30; // Max benefit at -3 (least FOMO)
    } else if (assessment.fomoScore <= 0) {
      fomoFactor = 30 - Math.abs(assessment.fomoScore + 3) * 10; // Decreasing benefit
    } else {
      fomoFactor = -assessment.fomoScore * 10; // Linear decrease for positive FOMO
    }
    
    // Market bias: doesn't affect position size
    const biasFactor = 0;
    
    // Physical condition: -3 (worst, more caution), +3 (best, less caution)
    const physicalFactor = assessment.physicalScore * 10;
    
    // Recent P&L: Peak at +1, both -3 and +3 warrant more caution
    let plFactor;
    if (assessment.plImpactScore <= -3) {
      plFactor = -30; // Max caution for big losses
    } else if (assessment.plImpactScore < 0) {
      plFactor = assessment.plImpactScore * 10; // Linear increase toward 0
    } else if (assessment.plImpactScore <= 1) {
      plFactor = assessment.plImpactScore * 15; // Peak at +1
    } else if (assessment.plImpactScore <= 2) {
      plFactor = 15 - (assessment.plImpactScore - 1) * 10; // Decrease for +2
    } else {
      plFactor = -15; // Significant caution for +3 (euphoria effect)
    }
    
    // Other (mindset): Linear relationship from -3 (bad vibes) to +3 (zen/in the zone)
    const otherFactor = assessment.otherScore * 10;
    
    // Base position size plus modifying factors
    const baseSize = 50;
    positionSize = baseSize + emotionalFactor + fomoFactor + biasFactor + physicalFactor + plFactor + otherFactor;
    
    // Ensure position size stays within bounds
    if (positionSize < 10) positionSize = 10; // Minimum 10%
    if (positionSize > 100) positionSize = 100; // Maximum 100%
    
    // Round to nearest integer
    positionSize = Math.round(positionSize);
    
    // Update position advice based on calculated size
    updatePositionAdvice();
    
    // Calculate overall score (can still be used for reference)
    assessment.overallScore = Math.round(
      (assessment.emotionalScore + assessment.fomoScore + assessment.biasScore + 
       assessment.physicalScore + assessment.plImpactScore + assessment.otherScore) / 6
    );
  }
  
  function updatePositionAdvice() {
    if (positionSize >= 80) {
      positionAdvice = {
        title: "Optimal Trading Conditions",
        tips: [
          "Trading conditions are excellent",
          "Standard to slightly larger position sizing appropriate",
          "Remain disciplined despite favorable conditions"
        ]
      };
    } else if (positionSize >= 60) {
      positionAdvice = {
        title: "Favorable Trading Conditions",
        tips: [
          "Good psychological and physical state for trading",
          "Standard position sizing appropriate",
          "Continue monitoring for any condition changes"
        ]
      };
    } else if (positionSize >= 40) {
      positionAdvice = {
        title: "Standard Trading Conditions",
        tips: [
          "Moderate caution advised",
          "Consider standard position sizing",
          "Focus on higher-probability setups"
        ]
      };
    } else if (positionSize >= 20) {
      positionAdvice = {
        title: "Caution: Reduce Position Sizing",
        tips: [
          "Current conditions suggest increased risk",
          "Consider reducing position size by 30-50%",
          "Focus only on highest-conviction setups"
        ]
      };
    } else {
      positionAdvice = {
        title: "High Risk: Minimal Trading Recommended",
        tips: [
          "Consider taking a trading break today",
          "If trading, reduce position size by 70-90%",
          "Only take extremely high-probability setups"
        ]
      };
    }
  }
  
  async function saveAssessment() {
    try {
      calculateRecommendedSize();
      const assessmentToSave = {
        ...assessment,
        date: new Date(assessment.date)
      };
      await SaveRiskAssessment(assessmentToSave);
      alert('Assessment saved successfully');
    } catch (error) {
      console.error('Failed to save assessment:', error);
      alert('Failed to save assessment: ' + error.message);
    }
  }
  
  function handleSliderChange() {
    calculateRecommendedSize();
  }
</script>

<div class="dashboard">
  <h2>Risk Management Dashboard</h2>
  <div class="divider"></div>
  
  <p class="description">Assess your daily emotional and psychological state to determine optimal position sizing.</p>
  
  <div class="assessment-form">
    <div class="date-nav">
      <button 
        class="nav-btn" 
        on:click={navigatePrevious} 
        disabled={isFirstAssessment() || assessments.length === 0}
      >
        &larr; Yesterday
      </button>
      
      <div class="date-picker">
        <label>Assessment Date:</label>
        <input type="date" bind:value={assessment.date} />
      </div>
      
      <button 
        class="nav-btn" 
        on:click={navigateNext} 
        disabled={isLastAssessment() || assessments.length === 0}
      >
        Tomorrow &rarr;
      </button>
    </div>
    
    <div class="slider-group">
      <h3>Emotional State: {assessment.emotionalScore}</h3>
      <input 
        type="range" 
        min="-3" 
        max="3" 
        step="1" 
        bind:value={assessment.emotionalScore} 
        on:change={handleSliderChange}
      />
      <div class="scale-labels">
        <span>Negative (-3)</span>
        <span>Neutral (0)</span>
        <span>Positive (+3)</span>
      </div>
      <p class="slider-desc">How are you feeling emotionally today? Negative values indicate stress, anxiety, or emotional disturbance.</p>
    </div>
    
    <div class="slider-group">
      <h3>FOMO Level: {assessment.fomoScore}</h3>
      <input 
        type="range" 
        min="-3" 
        max="3" 
        step="1" 
        bind:value={assessment.fomoScore} 
        on:change={handleSliderChange}
      />
      <div class="scale-labels">
        <span>Low (-3)</span>
        <span>Medium (0)</span>
        <span>High (+3)</span>
      </div>
      <p class="slider-desc">How strongly are you feeling the "Fear Of Missing Out"? Higher values indicate stronger FOMO which can lead to impulsive decisions.</p>
    </div>
    
    <div class="slider-group">
      <h3>Market Bias: {assessment.biasScore}</h3>
      <input 
        type="range" 
        min="-3" 
        max="3" 
        step="1" 
        bind:value={assessment.biasScore} 
        on:change={handleSliderChange}
      />
      <div class="scale-labels">
        <span>Bearish (-3)</span>
        <span>Neutral (0)</span>
        <span>Bullish (+3)</span>
      </div>
      <p class="slider-desc">Are you biased toward a bearish or bullish perspective? Strong biases can affect your trading decisions.</p>
      <p class="bias-caption">One must always be flexible. Go with the flow and check your unreasonable biases at the door.</p>
    </div>
    
    <div class="slider-group">
      <h3>Physical Condition: {assessment.physicalScore}</h3>
      <input 
        type="range" 
        min="-3" 
        max="3" 
        step="1" 
        bind:value={assessment.physicalScore} 
        on:change={handleSliderChange}
      />
      <div class="scale-labels">
        <span>Poor (-3)</span>
        <span>Average (0)</span>
        <span>Excellent (+3)</span>
      </div>
      <p class="slider-desc">How is your physical well-being today? Fatigue, illness, or poor sleep can impair decision-making.</p>
    </div>
    
    <div class="slider-group">
      <h3>Recent P&L Impact: {assessment.plImpactScore}</h3>
      <input 
        type="range" 
        min="-3" 
        max="3" 
        step="1" 
        bind:value={assessment.plImpactScore} 
        on:change={handleSliderChange}
      />
      <div class="scale-labels">
        <span>Losses (-3)</span>
        <span>Breakeven (0)</span>
        <span>Profits (+3)</span>
      </div>
      <p class="slider-desc">How have your recent trading results affected your mindset? Recent losses can lead to revenge trading.</p>
    </div>
    
    <div class="slider-group">
      <h3>Other Feelings: {assessment.otherScore}</h3>
      <input 
        type="range" 
        min="-3" 
        max="3" 
        step="1" 
        bind:value={assessment.otherScore} 
        on:change={handleSliderChange}
      />
      <div class="scale-labels">
        <span>Bad Vibes (-3)</span>
        <span>Normal (0)</span>
        <span>Zen/In the Zone (+3)</span>
      </div>
      <p class="slider-desc">How is your general trading mindset today? Are you feeling "in the zone" or experiencing bad vibes?</p>
    </div>
    
    <div class="button-row">
      <button class="save-btn" on:click={saveAssessment}>Save Assessment</button>
    </div>
  </div>
  
  <div class="recommendation">
    <h2>Recommended Position Size</h2>
    <div class="position-bar">
      <div 
        class="position-indicator" 
        style="width: {positionSize}%; background-color: {positionSize >= 70 ? '#68D391' : positionSize <= 30 ? '#FC8181' : '#F6AD55'}"
      ></div>
    </div>
    <div class="position-value">{positionSize}%</div>
    
    {#if showEuphoriaFlag}
      <div class="warning-flag euphoria-flag">
        <span class="flag-icon">⚠️</span>
        <span class="flag-text">Caution: Potential Euphoria Detected</span>
        <p class="flag-desc">High positive emotional state or recent profits may lead to overconfidence. Consider scaling back position sizes and being extra cautious with entries.</p>
      </div>
    {/if}
    
    {#if showChasingFlag}
      <div class="warning-flag chasing-flag">
        <span class="flag-icon">⚠️</span>
        <span class="flag-text">Warning: Potential Chasing Risk</span>
        <p class="flag-desc">Recent significant losses may lead to revenge trading or "chasing" lost money. Consider taking a break or trading minimal size with extreme caution.</p>
      </div>
    {/if}
    
    {#if showFOMOFlag}
      <div class="warning-flag fomo-flag">
        <span class="flag-icon">⚠️</span>
        <span class="flag-text">Warning: High FOMO Detected</span>
        <p class="flag-desc">Significant fear of missing out can lead to taking unreasonable risks and impulsive decisions. Step back, reassess, and only enter high-probability setups.</p>
      </div>
    {/if}
    
    <div class="position-advice">
      <h3>{positionAdvice.title}</h3>
      <ul>
        {#each positionAdvice.tips as tip}
          <li>{tip}</li>
        {/each}
      </ul>
    </div>
  </div>
  
  <div class="guidelines">
    <h2>Daily Trading Psychology Guidelines</h2>
    
    <div class="guidelines-grid">
      <div class="guideline-item">
        <h3>✓ Start with Self-Assessment</h3>
        <p>Begin each trading day by honestly evaluating your mental and physical state.</p>
      </div>
      
      <div class="guideline-item">
        <h3>✓ Adjust Position Sizing</h3>
        <p>Use your risk assessment score to modify position sizing - trade smaller on difficult days.</p>
      </div>
      
      <div class="guideline-item">
        <h3>✓ Recognize Emotional Triggers</h3>
        <p>Be aware of market events or conditions that might trigger emotional responses.</p>
      </div>
      
      <div class="guideline-item">
        <h3>✓ Implement Circuit Breakers</h3>
        <p>Have predefined rules for when to stop trading (e.g., after 2-3 consecutive losses).</p>
      </div>
    </div>
  </div>
</div>

<style>
  .dashboard {
    margin: 0 auto;
    padding: 2rem;
    background-color: var(--content-bg);
    border-radius: 5px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  h2 {
    color: inherit;
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
  
  .assessment-form {
    background-color: var(--card-bg);
    padding: 1.5rem;
    border-radius: 5px;
    margin-bottom: 2rem;
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  .date-picker {
    margin-bottom: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .date-picker label {
    margin-right: 1rem;
    font-weight: bold;
  }
  
  .date-picker input {
    padding: 0.5rem;
    border: 1px solid #cbd5e0;
    border-radius: 4px;
  }
  
  .date-nav {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.5rem;
  }
  
  .nav-btn {
    background-color: var(--bg-color);
    border: 1px solid var(--border-color);
    color: inherit;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
    transition: background-color 0.2s, color 0.2s;
  }
  
  .nav-btn:hover:not(:disabled) {
    background-color: var(--active-button);
    color: white;
  }
  
  .nav-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .slider-group {
    margin-bottom: 2rem;
  }
  
  .slider-group h3 {
    margin-bottom: 0.5rem;
    color: inherit;
  }
  
  .slider-group input[type="range"] {
    width: 100%;
    margin: 0.5rem 0;
  }
  
  .scale-labels {
    display: flex;
    justify-content: space-between;
    margin-top: 0.25rem;
    color: #718096;
    font-size: 0.85rem;
  }
  
  .slider-desc {
    color: #718096;
    font-style: italic;
    margin-top: 0.5rem;
    font-size: 0.9rem;
  }
  
  .bias-caption {
    color: #4299e1;
    font-style: italic;
    margin-top: 0.5rem;
    font-weight: 500;
    font-size: 0.9rem;
    text-align: center;
    padding: 5px;
    border-radius: 4px;
    background-color: rgba(66, 153, 225, 0.1);
  }
  
  .button-row {
    display: flex;
    justify-content: center;
    margin-top: 2rem;
  }
  
  .save-btn {
    background-color: var(--active-button);
    color: white;
    padding: 0.75rem 2rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    transition: background-color 0.2s;
  }
  
  .save-btn:hover {
    background-color: #3182ce;
  }
  
  .recommendation {
    background-color: var(--card-bg);
    padding: 1.5rem;
    border-radius: 5px;
    margin-bottom: 2rem;
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  .position-bar {
    height: 20px;
    background-color: #edf2f7;
    border-radius: 10px;
    margin: 1rem 0;
    overflow: hidden;
  }
  
  .position-indicator {
    height: 100%;
    border-radius: 10px;
    transition: width 0.3s ease, background-color 0.3s ease;
  }
  
  .position-value {
    text-align: center;
    font-weight: bold;
    margin-bottom: 1rem;
  }
  
  .warning-flag {
    padding: 1rem;
    border-radius: 5px;
    margin-bottom: 1rem;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .euphoria-flag {
    background-color: rgba(237, 137, 54, 0.15);
    border-left: 4px solid #ed8936;
  }
  
  .chasing-flag {
    background-color: rgba(229, 62, 62, 0.15);
    border-left: 4px solid #e53e3e;
  }
  
  .fomo-flag {
    background-color: rgba(102, 126, 234, 0.15);
    border-left: 4px solid #667eea;
  }
  
  .flag-icon {
    font-size: 1.5rem;
    margin-bottom: 0.5rem;
  }
  
  .flag-text {
    font-weight: bold;
    font-size: 1.1rem;
    margin-bottom: 0.5rem;
  }
  
  .flag-desc {
    text-align: center;
    font-size: 0.9rem;
    color: inherit;
    margin: 0;
  }
  
  .position-advice {
    border-top: 1px solid #e2e8f0;
    padding-top: 1rem;
  }
  
  .position-advice h3 {
    margin-bottom: 1rem;
    color: #ed8936;
  }
  
  .position-advice ul {
    padding-left: 1.5rem;
  }
  
  .position-advice li {
    margin-bottom: 0.5rem;
  }
  
  .guidelines {
    background-color: var(--card-bg);
    padding: 1.5rem;
    border-radius: 5px;
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  .guidelines-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1.5rem;
    margin-top: 1.5rem;
  }
  
  .guideline-item {
    padding: 1rem;
    background-color: var(--bg-color);
    border-radius: 5px;
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  .guideline-item h3 {
    color: var(--active-button);
    margin-bottom: 0.5rem;
  }
  
  .guideline-item p {
    color: inherit;
    font-size: 0.95rem;
  }
</style> 