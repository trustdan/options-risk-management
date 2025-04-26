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
    overallScore: 0,
    notes: ''
  };
  
  let positionSize = 50;
  let dateString = new Date().toISOString().split('T')[0];
  
  onMount(async () => {
    try {
      const assessments = await GetRiskAssessments();
      if (assessments && assessments.length > 0) {
        // Load the most recent assessment
        const latestAssessment = assessments[assessments.length - 1];
        assessment = {
          ...latestAssessment,
          date: new Date(latestAssessment.date).toISOString().split('T')[0]
        };
        calculateRecommendedSize();
      }
    } catch (error) {
      console.error('Failed to load risk assessments:', error);
    }
  });
  
  function calculateRecommendedSize() {
    // Calculate overall score as average of all factors
    assessment.overallScore = Math.round(
      (assessment.emotionalScore + assessment.fomoScore + assessment.biasScore + 
       assessment.physicalScore + assessment.plImpactScore) / 5
    );
    
    // Position size calculation - decrease position size as risk increases
    positionSize = 100 - (assessment.overallScore * 10);
    if (positionSize < 0) positionSize = 0;
    if (positionSize > 100) positionSize = 100;
  }
  
  async function saveAssessment() {
    try {
      calculateRecommendedSize();
      await SaveRiskAssessment({
        ...assessment,
        date: new Date(assessment.date)
      });
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
    <div class="date-picker">
      <label>Assessment Date:</label>
      <input type="date" bind:value={assessment.date} />
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
    
    <div class="button-row">
      <button class="save-btn" on:click={saveAssessment}>Save Assessment</button>
    </div>
  </div>
  
  <div class="recommendation">
    <h2>Recommended Position Size</h2>
    <div class="position-bar">
      <div class="position-indicator" style="width: {positionSize}%"></div>
    </div>
    <div class="position-value">{positionSize}%</div>
    
    <div class="position-advice">
      <h3>Caution: Standard position sizing recommended with care.</h3>
      <ul>
        <li>Trade with additional awareness of your current limitations</li>
        <li>Consider reducing position size by 30-50%</li>
        <li>Focus on high-conviction setups only</li>
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
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
    background-color: white;
    border-radius: 5px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  }
  
  h2 {
    color: #2d3748;
    text-align: center;
    margin-bottom: 0;
  }
  
  .divider {
    height: 2px;
    background-color: #e2e8f0;
    margin: 1rem 0 2rem;
  }
  
  .description {
    text-align: center;
    color: #4a5568;
    margin-bottom: 2rem;
  }
  
  .assessment-form {
    background-color: white;
    padding: 1.5rem;
    border-radius: 5px;
    margin-bottom: 2rem;
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
  
  .slider-group {
    margin-bottom: 2rem;
  }
  
  .slider-group h3 {
    margin-bottom: 0.5rem;
    color: #2d3748;
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
  
  .button-row {
    display: flex;
    justify-content: center;
    margin-top: 2rem;
  }
  
  .save-btn {
    background-color: #4299e1;
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
    background-color: white;
    padding: 1.5rem;
    border-radius: 5px;
    margin-bottom: 2rem;
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
    background-color: #f6ad55;
    border-radius: 10px;
  }
  
  .position-value {
    text-align: center;
    font-weight: bold;
    margin-bottom: 1rem;
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
    background-color: white;
    padding: 1.5rem;
    border-radius: 5px;
  }
  
  .guidelines-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1.5rem;
    margin-top: 1.5rem;
  }
  
  .guideline-item {
    padding: 1rem;
    background-color: #f7fafc;
    border-radius: 5px;
  }
  
  .guideline-item h3 {
    color: #38a169;
    margin-bottom: 0.5rem;
  }
  
  .guideline-item p {
    color: #4a5568;
    font-size: 0.95rem;
  }
</style> 