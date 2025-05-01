<script>
  import { onMount } from 'svelte';
  import { GetRiskAssessments, SaveRiskAssessment } from '../../../wailsjs/go/main/App';
  import { maxDollarRiskPerTrade } from '../../store/riskStore.js';
  
  // Import components for the tabbed interface
  import RiskAnalytics from './RiskAnalytics.svelte';
  import RiskManagementControls from './RiskManagementControls.svelte';
  
  // Active tab state
  let activeTab = 'assessment'; // 'assessment', 'analytics', or 'controls'
  
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
  let showStayOutWarning = false; // Flag for the stay out of market popup
  
  // For navigation
  let assessments = [];
  let currentIndex = 0;
  
  // For tracking scroll position to show/hide sticky mini-bar
  let showStickyBar = false;
  let mainRecommendationElement;
  
  onMount(async () => {
    try {
      assessments = await GetRiskAssessments();
      if (assessments && assessments.length > 0) {
        // Load the most recent assessment
        currentIndex = assessments.length - 1;
        loadAssessmentAt(currentIndex);
      }
      
      // Add scroll event listener to check when to show sticky bar
      window.addEventListener('scroll', handleScroll);
      
      // Update global position bar with initial value
      updateGlobalPositionBar(positionSize);
    } catch (error) {
      console.error('Failed to load risk assessments:', error);
    }
    
    return () => {
      // Clean up event listener
      window.removeEventListener('scroll', handleScroll);
    };
  });
  
  function loadAssessmentAt(index) {
    if (index >= 0 && index < assessments.length) {
      const latestAssessment = assessments[index];
      currentIndex = index;
      
      // Default values for the extended fields
      let physicalScore = 0;
      let plImpactScore = 0;
      let otherScore = 0;
      let userNotes = '';
      
      // Try to parse extra data from notes
      if (latestAssessment.notes) {
        try {
          const parsedData = JSON.parse(latestAssessment.notes);
          if (parsedData && typeof parsedData === 'object') {
            // Extract extra data if available
            if (parsedData.extraData) {
              physicalScore = parsedData.extraData.physicalScore || 0;
              plImpactScore = parsedData.extraData.plImpactScore || 0;
              otherScore = parsedData.extraData.otherScore || 0;
            }
            
            // Extract user notes if available
            userNotes = parsedData.userNotes || '';
          } else {
            // If not in our expected format, use as plain notes
            userNotes = latestAssessment.notes;
          }
        } catch (e) {
          // Not valid JSON, use as plain notes
          console.log('Notes is not valid JSON, using as plain text');
          userNotes = latestAssessment.notes;
        }
      }
      
      // Update the assessment object
      assessment = {
        id: latestAssessment.id || '',
        date: new Date(latestAssessment.date).toISOString().split('T')[0],
        emotionalScore: latestAssessment.emotionalScore || 0,
        fomoScore: latestAssessment.fomoScore || 0,
        biasScore: latestAssessment.biasScore || 0,
        physicalScore: physicalScore,
        plImpactScore: plImpactScore,
        otherScore: otherScore,
        overallScore: latestAssessment.overallScore || 0,
        notes: userNotes
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
    
    // Always keep bias score at 0 (neutral) since we removed the slider
    assessment.biasScore = 0;
    
    // Check for warning flags
    showEuphoriaFlag = assessment.emotionalScore >= 3 || assessment.plImpactScore >= 3;
    showChasingFlag = assessment.plImpactScore <= -3;
    showFOMOFlag = assessment.fomoScore >= 3;
    
    // Count extreme ratings to apply additional multiplier if multiple extremes exist
    let extremeRatingsCount = 0;
    
    if (assessment.emotionalScore <= -3 || assessment.emotionalScore >= 3) extremeRatingsCount++;
    if (assessment.fomoScore >= 3) extremeRatingsCount++;
    if (assessment.physicalScore <= -3) extremeRatingsCount++;
    if (assessment.plImpactScore <= -3 || assessment.plImpactScore >= 3) extremeRatingsCount++;
    if (assessment.otherScore <= -3) extremeRatingsCount++;
    
    // Define multiplier for extreme conditions - grows with multiple extreme ratings
    const extremeMultiplier = extremeRatingsCount >= 2 ? 3 : extremeRatingsCount >= 1 ? 2 : 1;
    
    // Count moderate-high ratings (-2 or +2 levels)
    let moderateHighCount = 0;
    
    if (assessment.emotionalScore === -2 || assessment.emotionalScore === 2) moderateHighCount++;
    if (assessment.fomoScore === 2) moderateHighCount++; 
    if (assessment.physicalScore === -2) moderateHighCount++;
    if (assessment.plImpactScore === -2 || assessment.plImpactScore === 2) moderateHighCount++;
    if (assessment.otherScore === -2) moderateHighCount++;
    
    // Define moderate multiplier - only applies to ratings of -2 or +2
    const moderateMultiplier = moderateHighCount >= 2 ? 1.7 : moderateHighCount >= 1 ? 1.4 : 1;
    
    // Emotional state: -3 (more caution), peaks at +1/+2 (less caution), +3 (more caution due to euphoria)
    let emotionalFactor;
    if (assessment.emotionalScore <= -3) {
      // Significantly increased impact for -3 emotional state to ensure visible changes
      emotionalFactor = -45 * extremeMultiplier; // Increased from -30 to -45
    } else if (assessment.emotionalScore === -2) {
      emotionalFactor = -15 * moderateMultiplier; // Enhanced impact for -2
    } else if (assessment.emotionalScore <= 0) {
      emotionalFactor = assessment.emotionalScore * 5; // Normal impact for non-extreme negative
    } else if (assessment.emotionalScore === 1) {
      emotionalFactor = 15; // Normal impact for +1
    } else if (assessment.emotionalScore === 2) {
      emotionalFactor = 20; // Slightly reduced enhanced impact for +2 (still positive)
    } else {
      emotionalFactor = -30 * extremeMultiplier; // Triple impact for +3 (euphoria) with multiplier
    }
    
    // FOMO level: -3 (peak performance, less caution), +3 (bad, more caution)
    let fomoFactor;
    if (assessment.fomoScore <= -3) {
      fomoFactor = 30; // Normal impact for positive extreme
    } else if (assessment.fomoScore <= 0) {
      fomoFactor = 30 - Math.abs(assessment.fomoScore + 3) * 10; // Normal curve
    } else if (assessment.fomoScore === 1) {
      fomoFactor = -5; // Normal impact for +1
    } else if (assessment.fomoScore === 2) {
      fomoFactor = -15 * moderateMultiplier; // Enhanced impact for +2
    } else {
      fomoFactor = -30 * extremeMultiplier; // Triple impact for +3 FOMO with multiplier
    }
    
    // Market bias: doesn't affect position size
    const biasFactor = 0;
    
    // Physical condition: -3 (worst, more caution), +3 (best, less caution)
    let physicalFactor;
    if (assessment.physicalScore <= -3) {
      // Significantly increased impact for -3 physical state to ensure visible changes
      physicalFactor = -45 * extremeMultiplier; // Increased from -30 to -45
    } else if (assessment.physicalScore === -2) {
      physicalFactor = -15 * moderateMultiplier; // Enhanced impact for -2
    } else if (assessment.physicalScore < 0) {
      physicalFactor = assessment.physicalScore * 5; // Normal impact for -1
    } else {
      physicalFactor = assessment.physicalScore * 10; // Normal positive impact
    }
    
    // Recent P&L: Peak at +1, both -3 and +3 warrant more caution
    let plFactor;
    if (assessment.plImpactScore <= -3) {
      plFactor = -45 * extremeMultiplier; // Triple impact of -15, plus extremeMultiplier
    } else if (assessment.plImpactScore === -2) {
      plFactor = -20 * moderateMultiplier; // Enhanced impact for -2
    } else if (assessment.plImpactScore < 0) {
      plFactor = assessment.plImpactScore * 7.5; // Normal impact for -1
    } else if (assessment.plImpactScore === 1) {
      plFactor = 22.5; // Normal for +1 (peak positive)
    } else if (assessment.plImpactScore === 2) {
      plFactor = 5 * moderateMultiplier; // Enhanced but still positive for +2 (starting to be cautious)
    } else {
      plFactor = -22.5 * extremeMultiplier; // Triple impact for +3 with multiplier
    }
    
    // Other (mindset): Linear relationship from -3 (bad vibes) to +3 (zen/in the zone)
    let otherFactor;
    if (assessment.otherScore <= -3) {
      otherFactor = -30 * extremeMultiplier; // Triple impact for extremely bad vibes with multiplier
    } else if (assessment.otherScore === -2) {
      otherFactor = -15 * moderateMultiplier; // Enhanced impact for -2
    } else if (assessment.otherScore < 0) {
      otherFactor = assessment.otherScore * 5; // Normal impact for -1
    } else {
      otherFactor = assessment.otherScore * 10; // Normal positive impact
    }
    
    // Log the extreme and moderate ratings and multipliers for debugging
    console.log(`Extreme ratings: ${extremeRatingsCount}, Multiplier: ${extremeMultiplier}`);
    console.log(`Moderate-high ratings: ${moderateHighCount}, Multiplier: ${moderateMultiplier.toFixed(2)}`);
    console.log(`Factors - Emotional: ${emotionalFactor.toFixed(1)}, FOMO: ${fomoFactor.toFixed(1)}, Physical: ${physicalFactor.toFixed(1)}, P&L: ${plFactor.toFixed(1)}, Other: ${otherFactor.toFixed(1)}`);
    
    // Increased base size from 50 to 80 to make average position size higher
    const baseSize = 80;
    positionSize = baseSize + emotionalFactor + fomoFactor + biasFactor + physicalFactor + plFactor + otherFactor;
    
    // Allow position size to drop below 30% only in extreme scenarios
    const isExtremeScenario = 
      ((assessment.emotionalScore <= -2 || assessment.physicalScore <= -2) && assessment.fomoScore >= 2) ||
      (assessment.plImpactScore <= -3 && assessment.fomoScore >= 2) ||
      extremeRatingsCount >= 2 || // Also consider it extreme if multiple extreme ratings exist
      assessment.emotionalScore <= -3 || // Always consider -3 emotional state as extreme scenario
      assessment.physicalScore <= -3;    // Always consider -3 physical state as extreme scenario
    
    // If not in an extreme scenario, enforce minimum of 30%
    if (!isExtremeScenario && positionSize < 30) {
      positionSize = 30;
    }
    
    // Lower minimum floor to better visualize differences in extreme cases
    if (positionSize < 3) positionSize = 3; // Changed from 5% to 3% minimum
    
    // Maximum is still 100%
    if (positionSize > 100) positionSize = 100;
    
    // Round to nearest integer
    positionSize = Math.round(positionSize);
    
    console.log(`Final position size: ${positionSize}%`);
    
    // Update global position bar if available
    updateGlobalPositionBar(positionSize);
    
    // Set stay out warning flag when position size is below 30%
    showStayOutWarning = positionSize < 30;
    
    // Update position advice based on calculated size
    updatePositionAdvice();
    
    // Calculate overall score (can still be used for reference)
    assessment.overallScore = Math.round(
      (assessment.emotionalScore + assessment.fomoScore + 
       assessment.physicalScore + assessment.plImpactScore + assessment.otherScore) / 5
    );
  }
  
  // Function to update the global position bar
  function updateGlobalPositionBar(size) {
    // Method 1: Use the register function if available
    // @ts-ignore - Custom property on window object
    if (typeof window.__positionBarRegister === 'function') {
      // @ts-ignore - Custom property on window object
      window.__positionBarRegister(size);
    }
    
    // Method 2: Dispatch a custom event for the global bar to listen to
    const event = new CustomEvent('position-size-updated', {
      detail: { positionSize: size }
    });
    window.dispatchEvent(event);
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
    } else if (positionSize >= 45) {
      positionAdvice = {
        title: "Standard Trading Conditions",
        tips: [
          "Moderate caution advised",
          "Consider standard position sizing",
          "Focus on higher-probability setups"
        ]
      };
    } else if (positionSize >= 35) {
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
          "If trading, reduce position size by 70%",
          "Only take extremely high-probability setups"
        ]
      };
    }
  }
  
  async function saveAssessment() {
    try {
      calculateRecommendedSize();
      
      // Import the RiskAssessment model
      const { models } = await import('../../../wailsjs/go/models');
      
      // Format date as RFC3339 string (ISO format that Go expects)
      const formattedDate = assessment.date + 'T00:00:00Z';
      
      // Store the extra fields and user notes in a structured format
      const combinedData = {
        userNotes: assessment.notes,
        extraData: {
          physicalScore: assessment.physicalScore,
          plImpactScore: assessment.plImpactScore,
          otherScore: assessment.otherScore
        }
      };
      
      // Create a new assessment using the model's createFrom method
      const assessmentData = models.RiskAssessment.createFrom({
        id: isEditingSameDay() ? assessment.id : '',
        date: formattedDate,
        emotionalScore: assessment.emotionalScore,
        fomoScore: assessment.fomoScore,
        biasScore: assessment.biasScore,
        overallScore: assessment.overallScore,
        notes: JSON.stringify(combinedData)
      });
      
      // Use the proper API method to save the assessment
      await SaveRiskAssessment(assessmentData);
      
      // Reload all assessments to get the updated list
      assessments = await GetRiskAssessments();
      
      // Find where our new assessment is in the list
      if (assessmentData.id) {
        // If updating, find by ID
        const idx = assessments.findIndex(a => a.id === assessmentData.id);
        if (idx >= 0) {
          currentIndex = idx;
        } else {
          currentIndex = assessments.length - 1; // Fallback to the most recent
        }
      } else {
        // New assessment will be the most recent
        currentIndex = assessments.length - 1;
      }
      
      // Load the assessment data
      loadAssessmentAt(currentIndex);
      
      alert('Assessment saved successfully');
    } catch (error) {
      console.error('Failed to save assessment:', error);
      alert('Failed to save assessment: ' + error);
    }
  }
  
  // Helper function to determine if we're editing the same day
  function isEditingSameDay() {
    if (!assessment.id) return false;
    
    // Find the current assessment in the list
    const existingAssessment = assessments.find(a => a.id === assessment.id);
    if (!existingAssessment) return false;
    
    // Check if the date is the same
    const existingDate = new Date(existingAssessment.date).toISOString().split('T')[0];
    const currentDate = assessment.date;
    
    return existingDate === currentDate;
  }
  
  function handleSliderChange() {
    calculateRecommendedSize();
  }
  
  // Function to handle scroll events and toggle sticky bar
  function handleScroll() {
    if (!mainRecommendationElement) return;
    
    // Get position of main recommendation bar
    const rect = mainRecommendationElement.getBoundingClientRect();
    
    // Show sticky bar when main bar is scrolled out of view
    showStickyBar = rect.bottom < 0;
  }
</script>

<div class="dashboard">
  <!-- Tab navigation -->
  <div class="tab-navigation">
    <button 
      class:active={activeTab === 'assessment'} 
      on:click={() => activeTab = 'assessment'}
    >
      Daily Assessment
    </button>
    <button 
      class:active={activeTab === 'analytics'} 
      on:click={() => activeTab = 'analytics'}
    >
      Analytics
    </button>
    <button 
      class:active={activeTab === 'controls'} 
      on:click={() => activeTab = 'controls'}
    >
      Position Sizing
    </button>
  </div>

  <!-- Daily Assessment Tab Content -->
  {#if activeTab === 'assessment'}
    <h2>Risk Management Dashboard</h2>
    <div class="divider"></div>
    
    <p class="description">Assess your daily emotional and psychological state to determine optimal position sizing.</p>
    
    <!-- Sticky Mini Position Bar (visible when scrolling) -->
    {#if showStickyBar}
      <div class="sticky-mini-bar">
        <div class="mini-position-label">Position Size:</div>
        <div class="mini-position-bar">
          <div 
            class="mini-position-indicator" 
            style="width: {positionSize}%; background-color: {positionSize >= 70 ? '#68D391' : positionSize <= 30 ? '#FC8181' : '#F6AD55'}"
          ></div>
        </div>
        <div class="mini-position-value">{positionSize}%</div>
      </div>
    {/if}
    
    <div class="recommendation" bind:this={mainRecommendationElement}>
      <h2>Recommended Position Size</h2>
      <div class="position-bar">
        <div 
          class="position-indicator" 
          style="width: {positionSize}%; background-color: {positionSize >= 70 ? '#68D391' : positionSize <= 30 ? '#FC8181' : '#F6AD55'}"
        ></div>
      </div>
      <div class="position-value">{positionSize}% of max position size is recommended</div>
      <div class="position-value-dollars">
        Position sizes of ${Math.round((positionSize/100) * $maxDollarRiskPerTrade).toFixed(2)} out of ${$maxDollarRiskPerTrade.toFixed(2)} are recommended for today
      </div>
      
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
        
        <div class="guideline-item">
          <h3>✓ Your First Loss Is Your Best Loss</h3>
          <p>Taking small initial losses prevents the much larger losses that come from hoping positions will recover.</p>
        </div>
        
        <div class="guideline-item">
          <h3>✓ Process Over Outcome</h3>
          <p>Focus on execution quality rather than profit/loss, as proper process leads to long-term results.</p>
        </div>
        
        <div class="guideline-item">
          <h3>✓ No Single Trade Should Matter</h3>
          <p>If any individual trade can substantially impact your account, your position sizing is excessive.</p>
        </div>
        
        <div class="guideline-item">
          <h3>✓ When Confused, Reduce Size</h3>
          <p>Uncertainty should trigger immediate risk reduction rather than hoping for clarity.</p>
        </div>
        
        <div class="guideline-item">
          <h3>✓ The Trade Journal Is Your Map</h3>
          <p>Systematic recording of your trades provides the objective feedback necessary for improvement.</p>
        </div>
        
        <div class="guideline-item">
          <h3>✓ Winners Take Small Losses By Choice</h3>
          <p>Proactive risk management prevents the catastrophic losses that destroy accounts.</p>
        </div>
        
        <div class="guideline-item">
          <h3>✓ Follow Rules Absolutely On Bad Days</h3>
          <p>Adherence to your trading plan is most critical precisely when emotions are strongest.</p>
        </div>
        
        <div class="guideline-item">
          <h3>✓ Daily Loss Limit As Circuit Breaker</h3>
          <p>Automatic trading cessation after reaching a specific drawdown protects you from poor decisions.</p>
        </div>
        
        <div class="guideline-item">
          <h3>✓ The Market Doesn't Know You Exist</h3>
          <p>Markets respond to supply and demand, not to your hopes, fears, or needs.</p>
        </div>
        
        <div class="guideline-item">
          <h3>✓ Avoid The Revenge Trade</h3>
          <p>Trades taken to recover losses typically compound those losses by bypassing risk management.</p>
        </div>
        
        <div class="guideline-item">
          <h3>✓ A Great Trader Knows When Not To Trade</h3>
          <p>The ability to stay on the sidelines during unfavorable conditions preserves capital.</p>
        </div>
      </div>
    </div>
  <!-- Analytics Tab Content -->
  {:else if activeTab === 'analytics'}
    <RiskAnalytics 
      {assessments} 
      on:refresh-assessments={async () => {
        try {
          assessments = await GetRiskAssessments();
        } catch (error) {
          console.error('Failed to refresh risk assessments:', error);
        }
      }}
    />
  <!-- Position Sizing Controls Tab Content -->
  {:else if activeTab === 'controls'}
    <RiskManagementControls />
  {/if}
</div>

<!-- Stay Out Warning Modal Popup -->
{#if showStayOutWarning}
  <div class="stay-out-modal">
    <button class="close-button" on:click={() => showStayOutWarning = false}>×</button>
    <div class="modal-content">
      <h3>⚠️ Market Warning</h3>
      <p>
        Your current psychological state suggests <strong>elevated risk</strong>. 
        Consider staying out of the market today or trading with extreme caution.
      </p>
    </div>
  </div>
{/if}

<style>
  .dashboard {
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
    justify-content: space-between;
    margin-bottom: 2rem;
    border-bottom: 1px solid var(--border-color);
    flex-wrap: wrap;
  }
  
  .tab-navigation button {
    padding: 0.75rem 1rem;
    background: none;
    border: none;
    cursor: pointer;
    color: inherit;
    font-weight: 500;
    border-bottom: 3px solid transparent;
    transition: all 0.3s ease;
    margin-bottom: -1px;
    white-space: nowrap;
  }
  
  .tab-navigation button.active {
    font-weight: bold;
    border-bottom: 3px solid var(--primary-button);
    color: var(--primary-button);
  }
  
  .tab-navigation button:hover:not(.active) {
    border-bottom-color: var(--border-color);
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
  
  .date-nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }
  
  .date-picker {
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
    border: 1px solid var(--input-border);
    background-color: var(--input-bg);
    color: inherit;
    border-radius: 4px;
  }
  
  .nav-btn {
    background-color: var(--secondary-button);
    color: inherit;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 500;
    opacity: 0.9;
    transition: all 0.2s;
  }
  
  .nav-btn:hover:not(:disabled) {
    opacity: 1;
  }
  
  .nav-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  .slider-group {
    margin-bottom: 1.5rem;
  }
  
  .slider-group h3 {
    margin-top: 0;
    margin-bottom: 0.5rem;
  }
  
  .slider-group input[type=range] {
    width: 100%;
    margin: 0.5rem 0;
    -webkit-appearance: none;
    background: var(--secondary-button);
    height: 8px;
    border-radius: 4px;
  }
  
  .slider-group input[type=range]::-webkit-slider-thumb {
    -webkit-appearance: none;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: var(--primary-button);
    cursor: pointer;
  }
  
  .scale-labels {
    display: flex;
    justify-content: space-between;
    font-size: 0.8rem;
    color: inherit;
    opacity: 0.8;
    margin-bottom: 0.5rem;
  }
  
  .slider-desc {
    font-size: 0.9rem;
    margin-top: 0.5rem;
    color: inherit;
    opacity: 0.9;
  }
  
  .bias-caption {
    font-style: italic;
    font-size: 0.8rem;
    margin-top: 0.5rem;
    text-align: center;
    color: inherit;
    opacity: 0.8;
  }
  
  .recommendation {
    background-color: var(--card-bg);
    padding: 1.5rem;
    border-radius: 5px;
    margin-bottom: 2rem;
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  .position-bar {
    width: 100%;
    height: 30px;
    background-color: rgba(0, 0, 0, 0.1);
    border-radius: 15px;
    margin: 1rem 0;
    overflow: hidden;
  }
  
  .position-indicator {
    height: 100%;
    border-radius: 15px;
    transition: width 0.5s, background-color 0.5s;
  }
  
  .position-value {
    text-align: center;
    font-size: 1.5rem;
    font-weight: bold;
    margin-bottom: 1rem;
  }
  
  .position-value-dollars {
    text-align: center;
    font-size: 1.1rem;
    color: inherit;
    margin-bottom: 1.5rem;
    font-weight: 500;
  }
  
  .warning-flag {
    margin: 1rem 0;
    padding: 1rem;
    border-radius: 5px;
    font-size: 0.9rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    position: relative;
  }
  
  .euphoria-flag {
    background-color: rgba(255, 193, 7, 0.2);
  }
  
  .chasing-flag {
    background-color: rgba(220, 53, 69, 0.2);
  }
  
  .fomo-flag {
    background-color: rgba(111, 66, 193, 0.2);
  }
  
  .stay-out-flag {
    background-color: rgba(220, 53, 69, 0.3);
    border: 1px solid rgba(220, 53, 69, 0.5);
    animation: pulse 2s infinite;
  }
  
  .close-warning-btn {
    position: absolute;
    top: 5px;
    right: 5px;
    background: transparent;
    border: none;
    color: rgba(0, 0, 0, 0.5);
    font-size: 1.2rem;
    font-weight: bold;
    cursor: pointer;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: all 0.2s;
  }
  
  .close-warning-btn:hover {
    background-color: rgba(0, 0, 0, 0.1);
    color: rgba(0, 0, 0, 0.8);
  }
  
  @keyframes pulse {
    0% {
      box-shadow: 0 0 0 0 rgba(220, 53, 69, 0.4);
    }
    70% {
      box-shadow: 0 0 0 8px rgba(220, 53, 69, 0);
    }
    100% {
      box-shadow: 0 0 0 0 rgba(220, 53, 69, 0);
    }
  }
  
  .flag-icon {
    font-size: 1.5rem;
    margin-bottom: 0.5rem;
  }
  
  .flag-text {
    font-weight: bold;
    margin-bottom: 0.5rem;
  }
  
  .flag-desc {
    text-align: center;
    margin: 0;
  }
  
  .position-advice {
    margin-top: 2rem;
    padding: 1rem;
    background-color: rgba(0, 0, 0, 0.05);
    border-radius: 5px;
  }
  
  .position-advice h3 {
    text-align: center;
    margin-top: 0;
    margin-bottom: 1rem;
  }
  
  .position-advice ul {
    padding-left: 1.5rem;
    margin: 0;
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
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 1rem;
    margin-top: 1.5rem;
  }
  
  .guideline-item {
    background-color: rgba(0, 0, 0, 0.03);
    padding: 1rem;
    border-radius: 5px;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }
  
  .guideline-item:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }
  
  .guideline-item h3 {
    margin-top: 0;
    font-size: 1rem;
    color: var(--primary-button);
  }
  
  .guideline-item p {
    margin: 0;
    font-size: 0.9rem;
    line-height: 1.4;
  }
  
  @media (max-width: 768px) {
    .dashboard {
      padding: 1rem;
    }
    
    .tab-navigation {
      flex-wrap: wrap;
      justify-content: center;
    }
    
    .tab-navigation button {
      flex: 1 1 auto;
      font-size: 0.9rem;
      padding: 0.5rem;
    }
    
    .date-nav {
      flex-direction: column;
      gap: 1rem;
    }
    
    .guidelines-grid {
      grid-template-columns: 1fr;
    }
  }
  
  .button-row {
    display: flex;
    justify-content: center;
    margin-top: 2rem;
  }
  
  .save-btn {
    background-color: var(--primary-button);
    color: white;
    border: none;
    padding: 0.75rem 2rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    transition: all 0.2s;
  }
  
  .save-btn:hover {
    opacity: 0.9;
  }
  
  /* Stay Out Modal Popup Styles */
  .stay-out-modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: var(--card-bg, white);
    border: 2px solid #e53e3e; /* Red border for warning */
    border-radius: 6px;
    padding: 15px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 9999;
    min-width: 300px;
    text-align: center;
    animation: fadeIn 0.3s ease-out, pulse 2s infinite;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
  
  .close-button {
    position: absolute;
    top: 5px;
    right: 8px;
    background: transparent;
    border: none;
    font-size: 18px;
    cursor: pointer;
    color: var(--text-muted, #718096);
    padding: 0;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
  }
  
  .close-button:hover {
    background: rgba(0, 0, 0, 0.05);
    color: var(--text-color, #333);
  }
  
  .modal-content {
    margin: 0 auto;
  }
  
  .modal-content h3 {
    margin: 0 0 10px 0;
    font-size: 18px;
    color: #e53e3e; /* Red text for warning */
  }
  
  .modal-content p {
    margin: 0;
    font-size: 14px;
    line-height: 1.5;
  }
  
  .modal-content strong {
    color: #e53e3e; /* Red text for emphasis */
  }
  
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  
  @keyframes pulse {
    0% { box-shadow: 0 0 0 0 rgba(229, 62, 62, 0.4); }
    70% { box-shadow: 0 0 0 10px rgba(229, 62, 62, 0); }
    100% { box-shadow: 0 0 0 0 rgba(229, 62, 62, 0); }
  }
  
  /* Sticky Mini Position Bar Styles */
  .sticky-mini-bar {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    background-color: var(--card-bg);
    padding: 8px 16px;
    display: flex;
    align-items: center;
    gap: 10px;
    z-index: 100;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
    animation: slideDown 0.3s forwards;
  }
  
  .mini-position-label {
    white-space: nowrap;
    font-weight: bold;
  }
  
  .mini-position-bar {
    flex: 1;
    height: 20px;
    background-color: rgba(0, 0, 0, 0.1);
    border-radius: 10px;
    overflow: hidden;
  }
  
  .mini-position-indicator {
    height: 100%;
    border-radius: 10px;
    transition: width 0.5s, background-color 0.5s;
  }
  
  .mini-position-value {
    font-weight: bold;
    min-width: 50px;
    text-align: right;
  }
  
  @keyframes slideDown {
    from { transform: translateY(-100%); }
    to { transform: translateY(0); }
  }
</style> 