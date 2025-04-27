<script>
  import { onMount } from 'svelte';
  
  export let assessments = [];
  let chartCanvas;
  let trendChartCanvas;
  let correlationChartCanvas;
  let chart;
  let trendChart;
  let correlationChart;
  let activeTimeFrame = 'all'; // 'all', 'month', 'week'
  let showCorrelation = false;
  
  let summary = {
    avgEmotional: 0,
    avgFOMO: 0,
    avgBias: 0,
    avgOverall: 0,
    lastWeek: [],
    lastMonth: [],
    performance: {
      goodDays: 0,
      neutralDays: 0,
      badDays: 0
    },
    trend: {
      emotional: 'neutral',
      fomo: 'neutral',
      bias: 'neutral'
    }
  };
  
  // Watch assessments and recalculate
  $: {
    if (assessments && assessments.length) {
      calculateSummary();
      if (chartCanvas) {
        renderCharts();
      }
    }
  }
  
  // Watch timeframe changes
  $: if (assessments && assessments.length && chartCanvas) {
    renderCharts();
  }
  
  onMount(() => {
    loadChartLibrary().then(() => {
      renderCharts();
    });
  });
  
  function loadChartLibrary() {
    return new Promise((resolve) => {
      if (window.Chart) {
        resolve();
        return;
      }
      
      const script = document.createElement('script');
      script.src = 'https://cdn.jsdelivr.net/npm/chart.js';
      script.onload = () => {
        // Load Chart.js plugins
        const annotationScript = document.createElement('script');
        annotationScript.src = 'https://cdn.jsdelivr.net/npm/chartjs-plugin-annotation';
        annotationScript.onload = () => resolve();
        document.head.appendChild(annotationScript);
      };
      document.head.appendChild(script);
    });
  }
  
  function getFilteredAssessments() {
    const now = new Date();
    
    if (activeTimeFrame === 'week') {
      const oneWeekAgo = new Date();
      oneWeekAgo.setDate(now.getDate() - 7);
      return assessments.filter(a => {
        const assessmentDate = new Date(a.date);
        return assessmentDate >= oneWeekAgo;
      });
    } else if (activeTimeFrame === 'month') {
      const oneMonthAgo = new Date();
      oneMonthAgo.setMonth(now.getMonth() - 1);
      return assessments.filter(a => {
        const assessmentDate = new Date(a.date);
        return assessmentDate >= oneMonthAgo;
      });
    } else {
      return [...assessments];
    }
  }
  
  function calculateTrend(scores) {
    if (scores.length < 3) return 'neutral';
    
    // Take last 5 scores or all if less than 5
    const lastScores = scores.slice(-Math.min(5, scores.length));
    const trend = lastScores[lastScores.length - 1] - lastScores[0];
    
    if (trend > 1) return 'improving';
    if (trend < -1) return 'declining';
    return 'neutral';
  }
  
  function calculateSummary() {
    // Calculate overall averages
    let totalEmotional = 0;
    let totalFOMO = 0;
    let totalBias = 0;
    let totalOverall = 0;
    
    assessments.forEach(a => {
      totalEmotional += a.emotionalScore;
      totalFOMO += a.fomoScore;
      totalBias += a.biasScore;
      totalOverall += a.overallScore;
      
      // Categorize performance
      if (a.overallScore >= 7) {
        summary.performance.goodDays++;
      } else if (a.overallScore <= 4) {
        summary.performance.badDays++;
      } else {
        summary.performance.neutralDays++;
      }
    });
    
    const count = assessments.length;
    summary.avgEmotional = Math.round((totalEmotional / count) * 10) / 10;
    summary.avgFOMO = Math.round((totalFOMO / count) * 10) / 10;
    summary.avgBias = Math.round((totalBias / count) * 10) / 10;
    summary.avgOverall = Math.round((totalOverall / count) * 10) / 10;
    
    // Calculate last 7 days summary
    const oneWeekAgo = new Date();
    oneWeekAgo.setDate(oneWeekAgo.getDate() - 7);
    
    summary.lastWeek = assessments.filter(a => {
      const assessmentDate = new Date(a.date);
      return assessmentDate >= oneWeekAgo;
    });
    
    // Calculate last month summary
    const oneMonthAgo = new Date();
    oneMonthAgo.setMonth(oneMonthAgo.getMonth() - 1);
    
    summary.lastMonth = assessments.filter(a => {
      const assessmentDate = new Date(a.date);
      return assessmentDate >= oneMonthAgo;
    });
    
    // Set default last week stats if no entries
    if (summary.lastWeek.length === 0) {
      summary.lastWeekEmotional = 'N/A';
      summary.lastWeekFOMO = 'N/A';
      summary.lastWeekBias = 'N/A';
    } else {
      let lwEmotional = 0, lwFOMO = 0, lwBias = 0;
      
      summary.lastWeek.forEach(a => {
        lwEmotional += a.emotionalScore;
        lwFOMO += a.fomoScore;
        lwBias += a.biasScore;
      });
      
      const lwCount = summary.lastWeek.length;
      summary.lastWeekEmotional = Math.round((lwEmotional / lwCount) * 10) / 10;
      summary.lastWeekFOMO = Math.round((lwFOMO / lwCount) * 10) / 10;
      summary.lastWeekBias = Math.round((lwBias / lwCount) * 10) / 10;
    }
    
    // Calculate trend direction
    summary.trend = {
      emotional: calculateTrend(assessments.map(a => a.emotionalScore)),
      fomo: calculateTrend(assessments.map(a => a.fomoScore)),
      bias: calculateTrend(assessments.map(a => a.biasScore))
    };
  }
  
  function renderCharts() {
    if (!window.Chart || !chartCanvas || assessments.length === 0) return;
    
    const filteredAssessments = getFilteredAssessments();
    if (filteredAssessments.length === 0) return;
    
    // Sort assessments by date (oldest to newest)
    filteredAssessments.sort((a, b) => {
      return new Date(a.date) - new Date(b.date);
    });
    
    // Prepare data for chart
    const dates = filteredAssessments.map(a => {
      const date = new Date(a.date);
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    });
    
    const overallScores = filteredAssessments.map(a => a.overallScore);
    const emotionalScores = filteredAssessments.map(a => a.emotionalScore);
    const fomoScores = filteredAssessments.map(a => a.fomoScore);
    const biasScores = filteredAssessments.map(a => a.biasScore);
    
    // Clean up previous chart if it exists
    if (chart) {
      chart.destroy();
    }
    
    // Create the main chart
    chart = new window.Chart(chartCanvas, {
      type: 'line',
      data: {
        labels: dates,
        datasets: [
          {
            label: 'Overall',
            data: overallScores,
            borderColor: 'rgb(54, 162, 235)',
            backgroundColor: 'rgba(54, 162, 235, 0.2)',
            tension: 0.1,
            fill: false
          },
          {
            label: 'Emotional',
            data: emotionalScores,
            borderColor: 'rgb(255, 99, 132)',
            backgroundColor: 'rgba(255, 99, 132, 0.2)',
            tension: 0.1,
            fill: false
          },
          {
            label: 'FOMO',
            data: fomoScores,
            borderColor: 'rgb(75, 192, 192)',
            backgroundColor: 'rgba(75, 192, 192, 0.2)',
            tension: 0.1,
            fill: false
          },
          {
            label: 'Bias',
            data: biasScores,
            borderColor: 'rgb(255, 159, 64)',
            backgroundColor: 'rgba(255, 159, 64, 0.2)',
            tension: 0.1,
            fill: false
          }
        ]
      },
      options: {
        responsive: true,
        scales: {
          y: {
            beginAtZero: true,
            max: 10,
            title: {
              display: true,
              text: 'Score (1-10)'
            }
          },
          x: {
            title: {
              display: true,
              text: 'Date'
            }
          }
        },
        plugins: {
          title: {
            display: true,
            text: 'Risk Assessment Trends'
          },
          annotation: {
            annotations: {
              line1: {
                type: 'line',
                yMin: 7,
                yMax: 7,
                borderColor: 'rgba(0, 200, 0, 0.3)',
                borderWidth: 2,
                borderDash: [6, 6],
                label: {
                  content: 'Good Zone',
                  enabled: true,
                  position: 'end'
                }
              },
              line2: {
                type: 'line',
                yMin: 4,
                yMax: 4,
                borderColor: 'rgba(255, 0, 0, 0.3)',
                borderWidth: 2,
                borderDash: [6, 6],
                label: {
                  content: 'Warning Zone',
                  enabled: true,
                  position: 'end'
                }
              }
            }
          }
        }
      }
    });
    
    // Create the moving average chart
    if (trendChartCanvas) {
      if (trendChart) {
        trendChart.destroy();
      }
      
      // Calculate 5-day moving averages
      const movingAvgWindow = Math.min(5, filteredAssessments.length);
      let emotionalMA = [];
      let fomoMA = [];
      let biasMA = [];
      
      if (filteredAssessments.length >= movingAvgWindow) {
        for (let i = movingAvgWindow - 1; i < filteredAssessments.length; i++) {
          let sumEmotional = 0;
          let sumFOMO = 0;
          let sumBias = 0;
          
          for (let j = 0; j < movingAvgWindow; j++) {
            sumEmotional += filteredAssessments[i - j].emotionalScore;
            sumFOMO += filteredAssessments[i - j].fomoScore;
            sumBias += filteredAssessments[i - j].biasScore;
          }
          
          emotionalMA.push(sumEmotional / movingAvgWindow);
          fomoMA.push(sumFOMO / movingAvgWindow);
          biasMA.push(sumBias / movingAvgWindow);
        }
      }
      
      // Create labels for moving average (need to remove first n-1 dates)
      const maLabels = dates.slice(movingAvgWindow - 1);
      
      trendChart = new window.Chart(trendChartCanvas, {
        type: 'line',
        data: {
          labels: maLabels,
          datasets: [
            {
              label: 'Emotional (5-day avg)',
              data: emotionalMA,
              borderColor: 'rgb(255, 99, 132)',
              backgroundColor: 'rgba(255, 99, 132, 0.5)',
              tension: 0.3,
              fill: true
            },
            {
              label: 'FOMO (5-day avg)',
              data: fomoMA,
              borderColor: 'rgb(75, 192, 192)',
              backgroundColor: 'rgba(75, 192, 192, 0.5)',
              tension: 0.3,
              fill: true
            },
            {
              label: 'Bias (5-day avg)',
              data: biasMA,
              borderColor: 'rgb(255, 159, 64)',
              backgroundColor: 'rgba(255, 159, 64, 0.5)',
              tension: 0.3,
              fill: true
            }
          ]
        },
        options: {
          responsive: true,
          scales: {
            y: {
              beginAtZero: true,
              max: 10,
              title: {
                display: true,
                text: 'Score (1-10)'
              }
            },
            x: {
              title: {
                display: true,
                text: 'Date'
              }
            }
          },
          plugins: {
            title: {
              display: true,
              text: 'Trend Analysis (5-day Moving Average)'
            }
          }
        }
      });
    }
    
    // Create correlation chart
    if (showCorrelation && correlationChartCanvas && filteredAssessments.length > 3) {
      if (correlationChart) {
        correlationChart.destroy();
      }
      
      correlationChart = new window.Chart(correlationChartCanvas, {
        type: 'scatter',
        data: {
          datasets: [
            {
              label: 'Emotional vs. FOMO',
              data: filteredAssessments.map(a => ({
                x: a.emotionalScore,
                y: a.fomoScore
              })),
              backgroundColor: 'rgba(255, 99, 132, 0.7)',
            },
            {
              label: 'Emotional vs. Bias',
              data: filteredAssessments.map(a => ({
                x: a.emotionalScore,
                y: a.biasScore
              })),
              backgroundColor: 'rgba(54, 162, 235, 0.7)',
            }
          ]
        },
        options: {
          responsive: true,
          scales: {
            x: {
              title: {
                display: true,
                text: 'Emotional Score'
              },
              min: 1,
              max: 10
            },
            y: {
              title: {
                display: true,
                text: 'FOMO/Bias Score'
              },
              min: 1,
              max: 10
            }
          },
          plugins: {
            title: {
              display: true,
              text: 'Correlation Analysis'
            }
          }
        }
      });
    }
  }
  
  function toggleCorrelation() {
    showCorrelation = !showCorrelation;
    setTimeout(() => {
      if (showCorrelation && correlationChartCanvas) {
        renderCharts();
      }
    }, 0);
  }
</script>

<div class="analytics">
  <h2>Risk Analytics Dashboard</h2>
  
  {#if assessments.length === 0}
    <p class="no-data">No risk assessments yet. Add some assessments to see analytics.</p>
  {:else}
    <div class="controls">
      <div class="timeframe-selector">
        <button class:active={activeTimeFrame === 'all'} on:click={() => activeTimeFrame = 'all'}>All Time</button>
        <button class:active={activeTimeFrame === 'month'} on:click={() => activeTimeFrame = 'month'}>Last Month</button>
        <button class:active={activeTimeFrame === 'week'} on:click={() => activeTimeFrame = 'week'}>Last Week</button>
      </div>
    </div>
    
    <!-- Main Dashboard Stats -->
    <div class="dashboard-header">
      <div class="quick-stats">
        <div class="stat-card">
          <div class="stat-title">Emotional State</div>
          <div class="stat-value">{summary.avgEmotional}</div>
          <div class="stat-trend {summary.trend.emotional}">
            {#if summary.trend.emotional === 'improving'}
              ↗ Improving
            {:else if summary.trend.emotional === 'declining'}
              ↘ Declining
            {:else}
              → Stable
            {/if}
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-title">FOMO Level</div>
          <div class="stat-value">{summary.avgFOMO}</div>
          <div class="stat-trend {summary.trend.fomo}">
            {#if summary.trend.fomo === 'improving'}
              ↗ Improving
            {:else if summary.trend.fomo === 'declining'}
              ↘ Declining
            {:else}
              → Stable
            {/if}
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-title">Bias Control</div>
          <div class="stat-value">{summary.avgBias}</div>
          <div class="stat-trend {summary.trend.bias}">
            {#if summary.trend.bias === 'improving'}
              ↗ Improving
            {:else if summary.trend.bias === 'declining'}
              ↘ Declining
            {:else}
              → Stable
            {/if}
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-title">Overall Score</div>
          <div class="stat-value">{summary.avgOverall}</div>
          <div class="stat-info">
            <span>{summary.performance.goodDays} good day{summary.performance.goodDays !== 1 ? 's' : ''}</span>
            <span>{summary.performance.badDays} bad day{summary.performance.badDays !== 1 ? 's' : ''}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Main Chart -->
    <div class="chart-container">
      <canvas bind:this={chartCanvas}></canvas>
    </div>
    
    <!-- Trend Chart -->
    <div class="chart-container">
      <canvas bind:this={trendChartCanvas}></canvas>
    </div>
    
    <!-- Correlation Analysis (Toggle) -->
    <div class="correlation-section">
      <button class="correlation-toggle" on:click={toggleCorrelation}>
        {showCorrelation ? 'Hide Correlation Analysis' : 'Show Correlation Analysis'}
      </button>
      
      {#if showCorrelation}
        <div class="chart-container correlation-chart">
          <canvas bind:this={correlationChartCanvas}></canvas>
        </div>
        <div class="correlation-explanation">
          <h4>Understanding Correlations</h4>
          <p>This chart helps visualize how your emotional state correlates with FOMO and bias levels. Clusters of points indicate potential patterns in your trading psychology.</p>
        </div>
      {/if}
    </div>
    
    <!-- Timeframe Stats -->
    <div class="stats">
      <div class="stats-section">
        <h3>Overall Averages</h3>
        <div class="stat-grid">
          <div class="stat-item">
            <span class="stat-label">Emotional:</span>
            <span class="stat-value">{summary.avgEmotional}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">FOMO:</span>
            <span class="stat-value">{summary.avgFOMO}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Bias:</span>
            <span class="stat-value">{summary.avgBias}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Overall:</span>
            <span class="stat-value">{summary.avgOverall}</span>
          </div>
        </div>
      </div>
      
      <div class="stats-section">
        <h3>Last 7 Days</h3>
        <div class="stat-grid">
          <div class="stat-item">
            <span class="stat-label">Emotional:</span>
            <span class="stat-value">{summary.lastWeekEmotional}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">FOMO:</span>
            <span class="stat-value">{summary.lastWeekFOMO}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Bias:</span>
            <span class="stat-value">{summary.lastWeekBias}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Entries:</span>
            <span class="stat-value">{summary.lastWeek.length}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Insights Section -->
    <div class="insights-section">
      <h3>Trading Psychology Insights</h3>
      <div class="insights-content">
        {#if summary.trend.emotional === 'declining' || summary.trend.fomo === 'declining' || summary.trend.bias === 'declining'}
          <div class="insight-card warning">
            <h4>Warning Signs</h4>
            <ul>
              {#if summary.trend.emotional === 'declining'}
                <li>Your emotional control is trending downward - consider taking a break if this continues.</li>
              {/if}
              {#if summary.trend.fomo === 'declining'}
                <li>Increasing FOMO levels detected - review recent impulsive decisions.</li>
              {/if}
              {#if summary.trend.bias === 'declining'}
                <li>Market bias control is weakening - be aware of confirmation bias in your analysis.</li>
              {/if}
            </ul>
          </div>
        {/if}
        
        {#if summary.trend.emotional === 'improving' || summary.trend.fomo === 'improving' || summary.trend.bias === 'improving'}
          <div class="insight-card positive">
            <h4>Positive Developments</h4>
            <ul>
              {#if summary.trend.emotional === 'improving'}
                <li>Your emotional control is improving - continue your current practices.</li>
              {/if}
              {#if summary.trend.fomo === 'improving'}
                <li>FOMO levels are decreasing - your patience is paying off.</li>
              {/if}
              {#if summary.trend.bias === 'improving'}
                <li>Better market bias control - your objective analysis is strengthening.</li>
              {/if}
            </ul>
          </div>
        {/if}
        
        <div class="insight-card neutral">
          <h4>Recommendations</h4>
          <ul>
            {#if summary.avgEmotional < 5}
              <li>Consider implementing emotional check-ins before each trade.</li>
            {/if}
            {#if summary.avgFOMO > 6}
              <li>Your FOMO score is elevated - try enforcing a 24-hour rule before entering trades.</li>
            {/if}
            {#if summary.avgBias < 5}
              <li>Practice reviewing both confirming and contradicting evidence for each trade thesis.</li>
            {/if}
            {#if summary.performance.badDays > summary.performance.goodDays}
              <li>Review your strategy during periods of sub-optimal psychological state.</li>
            {/if}
          </ul>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .analytics {
    max-width: 900px;
    margin: 0 auto;
    padding: 1rem;
  }
  
  h2 {
    text-align: center;
    margin-bottom: 1.5rem;
    color: var(--text-color);
  }
  
  .controls {
    margin-bottom: 1.5rem;
  }
  
  .timeframe-selector {
    display: flex;
    justify-content: center;
    gap: 0.5rem;
  }
  
  .timeframe-selector button {
    background-color: var(--button-bg);
    color: var(--button-text);
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .timeframe-selector button.active {
    background-color: #4caf50;
    color: white;
  }
  
  .dashboard-header {
    margin-bottom: 2rem;
  }
  
  .quick-stats {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
  }
  
  .stat-card {
    background-color: var(--card-bg);
    color: var(--text-color);
    padding: 1rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    text-align: center;
  }
  
  .stat-title {
    font-size: 0.9rem;
    color: var(--text-secondary);
    margin-bottom: 0.5rem;
  }
  
  .stat-value {
    font-size: 2rem;
    font-weight: bold;
    margin-bottom: 0.5rem;
    color: var(--text-color);
  }
  
  .stat-trend {
    font-size: 0.8rem;
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
  }
  
  .stat-trend.improving {
    background-color: rgba(0, 200, 0, 0.1);
    color: var(--success-color);
  }
  
  .stat-trend.declining {
    background-color: rgba(255, 0, 0, 0.1);
    color: var(--danger-color);
  }
  
  .stat-trend.neutral {
    background-color: rgba(200, 200, 200, 0.2);
    color: var(--text-secondary);
  }
  
  .stat-info {
    display: flex;
    justify-content: space-around;
    font-size: 0.8rem;
    color: var(--text-secondary);
  }
  
  .no-data {
    text-align: center;
    font-style: italic;
    color: var(--text-secondary);
    padding: 2rem;
  }
  
  .chart-container {
    width: 100%;
    height: 300px;
    margin-bottom: 2rem;
    background-color: var(--card-bg);
    color: var(--text-color);
    padding: 1rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }
  
  .correlation-section {
    margin-bottom: 2rem;
    text-align: center;
  }
  
  .correlation-toggle {
    background-color: var(--button-bg);
    color: var(--button-text);
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    margin-bottom: 1rem;
  }
  
  .correlation-chart {
    height: 400px;
  }
  
  .correlation-explanation {
    background-color: var(--card-bg-light);
    color: var(--text-color);
    padding: 1rem;
    border-radius: 8px;
    margin-top: 1rem;
    font-size: 0.9rem;
  }
  
  .correlation-explanation h4 {
    margin-top: 0;
    margin-bottom: 0.5rem;
    color: var(--text-color);
  }
  
  .stats {
    display: flex;
    flex-wrap: wrap;
    gap: 2rem;
    justify-content: space-around;
    margin-bottom: 2rem;
  }
  
  .stats-section {
    flex: 1;
    min-width: 250px;
    padding: 1rem;
    background-color: var(--card-bg);
    color: var(--text-color);
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }
  
  .stats-section h3 {
    margin-top: 0;
    margin-bottom: 1rem;
    text-align: center;
    font-size: 1.2rem;
    color: var(--text-color);
  }
  
  .stats-section .stat-value,
  .stats-section .stat-label {
    color: var(--text-color);
  }
  
  .stat-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }
  
  .stat-item {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .stat-label {
    font-weight: bold;
    color: var(--text-secondary);
  }
  
  .stat-value {
    font-size: 1.5rem;
    font-weight: bold;
    color: var(--text-color);
  }
  
  .insights-section {
    background-color: var(--card-bg);
    color: var(--text-color);
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }
  
  .insights-content {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1rem;
  }
  
  .insight-card {
    padding: 1rem;
    border-radius: 8px;
    color: var(--text-color);
  }
  
  .insight-card h4 {
    margin-top: 0;
    margin-bottom: 0.75rem;
    color: var(--text-color);
  }
  
  .insight-card ul {
    margin: 0;
    padding-left: 1.5rem;
  }
  
  .insight-card li {
    margin-bottom: 0.5rem;
  }
  
  .insight-card.warning {
    background-color: var(--warning-bg);
  }
  
  .insight-card.positive {
    background-color: var(--success-bg);
  }
  
  .insight-card.neutral {
    background-color: var(--neutral-bg);
  }
  
  /* Responsive adjustments */
  @media (max-width: 768px) {
    .quick-stats {
      grid-template-columns: 1fr 1fr;
    }
    
    .insights-content {
      grid-template-columns: 1fr;
    }
  }
  
  @media (max-width: 480px) {
    .quick-stats {
      grid-template-columns: 1fr;
    }
  }
</style> 