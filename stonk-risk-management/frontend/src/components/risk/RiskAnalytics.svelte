<script>
  import { onMount, createEventDispatcher } from 'svelte';
  import { DeleteRiskAssessment } from '../../../wailsjs/go/main/App';
  
  export let assessments = [];
  let chartCanvas;
  let chart;
  let isDeleting = false;
  let sortedEntries = [];
  let sortField = 'date';
  let sortDirection = 'desc';
  
  const dispatch = createEventDispatcher();
  
  let summary = {
    avgEmotional: 0,
    avgFOMO: 0,
    avgBias: 0,
    avgPhysical: 0,
    avgPLImpact: 0,
    avgOther: 0,
    avgOverall: 0,
    lastWeek: []
  };
  
  $: {
    if (assessments && assessments.length) {
      calculateSummary();
      updateSortedEntries();
      if (chartCanvas) {
        renderChart();
      }
    }
  }
  
  function updateSortedEntries() {
    sortedEntries = [...assessments].sort((a, b) => {
      let aValue = a[sortField];
      let bValue = b[sortField];
      
      // Handle date specially
      if (sortField === 'date') {
        aValue = new Date(aValue).getTime();
        bValue = new Date(bValue).getTime();
      }
      
      // For numerical values
      if (typeof aValue === 'number' && typeof bValue === 'number') {
        return sortDirection === 'asc' ? aValue - bValue : bValue - aValue;
      }
      
      // For string values
      return sortDirection === 'asc' 
        ? String(aValue).localeCompare(String(bValue))
        : String(bValue).localeCompare(String(aValue));
    });
  }
  
  function sortBy(field) {
    if (sortField === field) {
      // Toggle direction if clicking the same field
      sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
    } else {
      // New field, default to descending for dates, ascending for others
      sortField = field;
      sortDirection = field === 'date' ? 'desc' : 'asc';
    }
    updateSortedEntries();
  }
  
  function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
  
  async function deleteEntry(id) {
    if (isDeleting) return;
    
    if (confirm('Are you sure you want to delete this assessment?')) {
      isDeleting = true;
      try {
        await DeleteRiskAssessment(id);
        // Signal the parent component to refresh
        dispatch('refresh-assessments');
      } catch (error) {
        console.error('Failed to delete assessment:', error);
        alert('Failed to delete assessment. Please try again.');
      } finally {
        isDeleting = false;
      }
    }
  }
  
  onMount(() => {
    loadChartLibrary().then(() => {
      renderChart();
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
      script.onload = () => resolve();
      document.head.appendChild(script);
    });
  }
  
  function calculateSummary() {
    // Calculate overall averages
    let totalEmotional = 0;
    let totalFOMO = 0;
    let totalBias = 0;
    let totalPhysical = 0;
    let totalPLImpact = 0;
    let totalOther = 0;
    let totalOverall = 0;
    
    assessments.forEach(a => {
      totalEmotional += a.emotionalScore;
      totalFOMO += a.fomoScore;
      totalBias += a.biasScore;
      totalPhysical += a.physicalScore || 0;
      totalPLImpact += a.plImpactScore || 0;
      totalOther += a.otherScore || 0;
      totalOverall += a.overallScore;
    });
    
    const count = assessments.length;
    summary.avgEmotional = Math.round((totalEmotional / count) * 10) / 10;
    summary.avgFOMO = Math.round((totalFOMO / count) * 10) / 10;
    summary.avgBias = Math.round((totalBias / count) * 10) / 10;
    summary.avgPhysical = Math.round((totalPhysical / count) * 10) / 10;
    summary.avgPLImpact = Math.round((totalPLImpact / count) * 10) / 10;
    summary.avgOther = Math.round((totalOther / count) * 10) / 10;
    summary.avgOverall = Math.round((totalOverall / count) * 10) / 10;
    
    // Calculate last 7 days summary
    const oneWeekAgo = new Date();
    oneWeekAgo.setDate(oneWeekAgo.getDate() - 7);
    
    summary.lastWeek = assessments.filter(a => {
      const assessmentDate = new Date(a.date);
      return assessmentDate >= oneWeekAgo;
    });
    
    // Set default last week stats if no entries
    if (summary.lastWeek.length === 0) {
      summary.lastWeekEmotional = 'N/A';
      summary.lastWeekFOMO = 'N/A';
      summary.lastWeekBias = 'N/A';
      summary.lastWeekPhysical = 'N/A';
      summary.lastWeekPLImpact = 'N/A';
      summary.lastWeekOther = 'N/A';
    } else {
      let lwEmotional = 0, lwFOMO = 0, lwBias = 0, lwPhysical = 0, lwPLImpact = 0, lwOther = 0;
      
      summary.lastWeek.forEach(a => {
        lwEmotional += a.emotionalScore;
        lwFOMO += a.fomoScore;
        lwBias += a.biasScore;
        lwPhysical += a.physicalScore || 0;
        lwPLImpact += a.plImpactScore || 0;
        lwOther += a.otherScore || 0;
      });
      
      const lwCount = summary.lastWeek.length;
      summary.lastWeekEmotional = Math.round((lwEmotional / lwCount) * 10) / 10;
      summary.lastWeekFOMO = Math.round((lwFOMO / lwCount) * 10) / 10;
      summary.lastWeekBias = Math.round((lwBias / lwCount) * 10) / 10;
      summary.lastWeekPhysical = Math.round((lwPhysical / lwCount) * 10) / 10;
      summary.lastWeekPLImpact = Math.round((lwPLImpact / lwCount) * 10) / 10;
      summary.lastWeekOther = Math.round((lwOther / lwCount) * 10) / 10;
    }
  }
  
  function renderChart() {
    if (!window.Chart || !chartCanvas || assessments.length === 0) return;
    
    // Create a fixed 21-day date range: 19 days back + today + 1 day ahead
    const today = new Date();
    today.setHours(0, 0, 0, 0); // Normalize to start of day
    
    const startDate = new Date(today);
    startDate.setDate(today.getDate() - 19);
    
    const endDate = new Date(today);
    endDate.setDate(today.getDate() + 1);
    
    // Generate all dates in the range
    const allDatesInRange = [];
    const currentDate = new Date(startDate);
    
    while (currentDate <= endDate) {
      allDatesInRange.push(new Date(currentDate));
      currentDate.setDate(currentDate.getDate() + 1);
    }
    
    // Sort assessments by date
    const sortedAssessments = [...assessments].sort((a, b) => {
      return new Date(a.date).getTime() - new Date(b.date).getTime();
    });
    
    // Map the assessments to the date range - we'll keep the latest assessment for each date
    const assessmentsByDate = new Map();
    
    sortedAssessments.forEach(assessment => {
      const assessmentDate = new Date(assessment.date);
      // Normalize to start of day for comparing
      assessmentDate.setHours(0, 0, 0, 0);
      const dateKey = assessmentDate.toISOString().split('T')[0];
      
      // If we already have an assessment for this date, only replace if this one is newer
      if (!assessmentsByDate.has(dateKey) || 
          new Date(assessment.date) > new Date(assessmentsByDate.get(dateKey).date)) {
        assessmentsByDate.set(dateKey, assessment);
      }
    });
    
    // Prepare data for chart - for all dates in range
    const dates = allDatesInRange.map(date => 
      date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
    );
    
    // Create arrays for each metric, with null for dates without data
    // Important: use 0 (not null) for dates with 0 values to show the trendline
    const emotionalScores = [];
    const fomoScores = [];
    const biasScores = [];
    const physicalScores = [];
    const plImpactScores = [];
    const otherScores = [];
    
    // Fill in data for each date
    allDatesInRange.forEach(date => {
      const dateKey = date.toISOString().split('T')[0];
      const assessment = assessmentsByDate.get(dateKey);
      
      if (assessment) {
        // Convert scores from -3 to +3 scale to 0-7 scale for display
        // (add 3 to each value to shift from -3..+3 to 0..6, then add 0.5 for center alignment)
        // Use Number() to ensure we're working with numbers, and the || 0 ensures null/undefined values become 0
        // The Number check is essential to ensure the values are treated as numbers for addition
        emotionalScores.push(Number(assessment.emotionalScore !== undefined ? assessment.emotionalScore : 0) + 3.5);
        fomoScores.push(Number(assessment.fomoScore !== undefined ? assessment.fomoScore : 0) + 3.5);
        biasScores.push(Number(assessment.biasScore !== undefined ? assessment.biasScore : 0) + 3.5);
        physicalScores.push(Number(assessment.physicalScore !== undefined ? assessment.physicalScore : 0) + 3.5);
        plImpactScores.push(Number(assessment.plImpactScore !== undefined ? assessment.plImpactScore : 0) + 3.5);
        otherScores.push(Number(assessment.otherScore !== undefined ? assessment.otherScore : 0) + 3.5);
      } else {
        // No data for this date - use null to show gaps in chart
        emotionalScores.push(null);
        fomoScores.push(null);
        biasScores.push(null);
        physicalScores.push(null);
        plImpactScores.push(null);
        otherScores.push(null);
      }
    });
    
    // Clean up previous chart if it exists
    if (chart) {
      chart.destroy();
    }
    
    // Create the chart with updated options 
    chart = new window.Chart(chartCanvas, {
      type: 'line',
      data: {
        labels: dates,
        datasets: [
          {
            label: 'Emotional State',
            data: emotionalScores,
            borderColor: 'rgb(255, 99, 132)',
            backgroundColor: 'rgba(255, 99, 132, 0.2)',
            tension: 0.1,
            fill: false,
            spanGaps: true, // Connect points across gaps
            pointRadius: 4, // Make points a bit larger
            pointHoverRadius: 6
          },
          {
            label: 'FOMO Level',
            data: fomoScores,
            borderColor: 'rgb(75, 192, 192)',
            backgroundColor: 'rgba(75, 192, 192, 0.2)',
            tension: 0.1,
            fill: false,
            spanGaps: true,
            pointRadius: 4,
            pointHoverRadius: 6
          },
          {
            label: 'Market Bias',
            data: biasScores,
            borderColor: 'rgb(255, 159, 64)',
            backgroundColor: 'rgba(255, 159, 64, 0.2)',
            tension: 0.1,
            fill: false,
            spanGaps: true,
            pointRadius: 4,
            pointHoverRadius: 6
          },
          {
            label: 'Physical Condition',
            data: physicalScores,
            borderColor: 'rgb(54, 162, 235)',
            backgroundColor: 'rgba(54, 162, 235, 0.2)',
            tension: 0.1,
            fill: false,
            spanGaps: true,
            pointRadius: 4,
            pointHoverRadius: 6
          },
          {
            label: 'Recent P&L Impact',
            data: plImpactScores,
            borderColor: 'rgb(153, 102, 255)',
            backgroundColor: 'rgba(153, 102, 255, 0.2)',
            tension: 0.1,
            fill: false,
            spanGaps: true,
            pointRadius: 4,
            pointHoverRadius: 6
          },
          {
            label: 'Other Feelings',
            data: otherScores,
            borderColor: 'rgb(255, 205, 86)',
            backgroundColor: 'rgba(255, 205, 86, 0.2)',
            tension: 0.1,
            fill: false,
            spanGaps: true,
            pointRadius: 4,
            pointHoverRadius: 6
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          y: {
            beginAtZero: true,
            min: 0,
            max: 7,
            grid: {
              color: 'rgba(255, 255, 255, 0.1)' // Lighter grid lines for dark mode
            },
            ticks: {
              stepSize: 1,
              callback: function(value) {
                // Convert y-axis labels back to -3 to +3 scale
                return value - 3.5;
              },
              color: 'rgba(255, 255, 255, 0.8)' // Lighter text for dark mode
            },
            title: {
              display: true,
              text: 'Score (-3 to +3)',
              color: 'rgba(255, 255, 255, 0.9)' // Lighter text for dark mode
            }
          },
          x: {
            grid: {
              color: 'rgba(255, 255, 255, 0.1)' // Lighter grid lines for dark mode
            },
            ticks: {
              color: 'rgba(255, 255, 255, 0.8)' // Lighter text for dark mode
            },
            title: {
              display: true,
              text: 'Date',
              color: 'rgba(255, 255, 255, 0.9)' // Lighter text for dark mode
            }
          }
        },
        plugins: {
          legend: {
            labels: {
              color: 'rgba(255, 255, 255, 0.9)', // Lighter legend text for dark mode
              boxWidth: 12,  // Make legend icons smaller
              padding: 15    // Add more padding between items
            }
          },
          title: {
            display: true,
            text: 'Risk Assessment Trends (Last 21 Days)',
            color: 'rgba(255, 255, 255, 0.9)', // Lighter title for dark mode
            font: {
              size: 16
            }
          },
          tooltip: {
            backgroundColor: 'rgba(0, 0, 0, 0.8)',
            titleColor: 'rgba(255, 255, 255, 1)',
            bodyColor: 'rgba(255, 255, 255, 1)',
            callbacks: {
              label: function(context) {
                const rawValue = context.raw;
                const actualValue = typeof rawValue === 'number' ? rawValue - 3.5 : 0;
                return `${context.dataset.label}: ${actualValue.toFixed(1)}`;
              }
            }
          }
        }
      }
    });
  }
</script>

<div class="analytics">
  <h2>Risk Analytics</h2>
  
  {#if assessments.length === 0}
    <p class="no-data">No risk assessments yet. Add some assessments to see analytics.</p>
  {:else}
    <div class="chart-container">
      <canvas bind:this={chartCanvas}></canvas>
    </div>
    
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
            <span class="stat-label">Physical:</span>
            <span class="stat-value">{summary.avgPhysical}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">P&L Impact:</span>
            <span class="stat-value">{summary.avgPLImpact}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Other:</span>
            <span class="stat-value">{summary.avgOther}</span>
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
            <span class="stat-label">Physical:</span>
            <span class="stat-value">{summary.lastWeekPhysical}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">P&L Impact:</span>
            <span class="stat-value">{summary.lastWeekPLImpact}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Other:</span>
            <span class="stat-value">{summary.lastWeekOther}</span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="entries-section">
      <h3>All Assessment Entries</h3>
      <div class="entries-table-container">
        <table class="entries-table">
          <thead>
            <tr>
              <th on:click={() => sortBy('date')}>
                Date {sortField === 'date' ? (sortDirection === 'asc' ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('emotionalScore')}>
                Emotional {sortField === 'emotionalScore' ? (sortDirection === 'asc' ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('fomoScore')}>
                FOMO {sortField === 'fomoScore' ? (sortDirection === 'asc' ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('biasScore')}>
                Bias {sortField === 'biasScore' ? (sortDirection === 'asc' ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('physicalScore')}>
                Physical {sortField === 'physicalScore' ? (sortDirection === 'asc' ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('plImpactScore')}>
                P&L {sortField === 'plImpactScore' ? (sortDirection === 'asc' ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('otherScore')}>
                Other {sortField === 'otherScore' ? (sortDirection === 'asc' ? '↑' : '↓') : ''}
              </th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each sortedEntries as entry}
              <tr>
                <td>{formatDate(entry.date)}</td>
                <td>{entry.emotionalScore}</td>
                <td>{entry.fomoScore}</td>
                <td>{entry.biasScore}</td>
                <td>{entry.physicalScore || 0}</td>
                <td>{entry.plImpactScore || 0}</td>
                <td>{entry.otherScore || 0}</td>
                <td>
                  <button 
                    class="delete-btn" 
                    on:click={() => deleteEntry(entry.id)}
                    disabled={isDeleting}
                  >
                    Delete
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  {/if}
</div>

<style>
  .analytics {
    max-width: 800px;
    margin: 0 auto;
    padding: 1rem;
    color: rgba(255, 255, 255, 0.9); /* Add text color for dark mode */
  }
  
  h2 {
    text-align: center;
    margin-bottom: 1.5rem;
    color: rgba(255, 255, 255, 0.9); /* Add text color for dark mode */
  }
  
  .no-data {
    text-align: center;
    font-style: italic;
    color: rgba(255, 255, 255, 0.6);
    padding: 2rem;
  }
  
  .chart-container {
    width: 100%;
    height: 450px; /* Increase height for better visibility */
    margin-bottom: 2rem;
    background-color: rgba(29, 36, 53, 0.6); /* Dark background for the chart */
    border-radius: 8px;
    padding: 1rem;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
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
    padding: 1.5rem;
    background-color: rgba(29, 36, 53, 0.8); /* Darker background for dark mode */
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  }
  
  h3 {
    margin-top: 0;
    margin-bottom: 1rem;
    text-align: center;
    font-size: 1.2rem;
    color: rgba(255, 255, 255, 0.9); /* Add text color for dark mode */
  }
  
  .stat-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
  }
  
  .stat-item {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .stat-label {
    font-weight: bold;
    color: rgba(255, 255, 255, 0.7);
    margin-bottom: 0.3rem;
  }
  
  .stat-value {
    font-size: 1.8rem;
    font-weight: bold;
    color: rgba(255, 255, 255, 0.9);
  }
  
  .entries-section {
    margin-top: 2rem;
    padding: 1.5rem;
    background-color: rgba(29, 36, 53, 0.8); /* Darker background for dark mode */
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
  }
  
  .entries-table-container {
    overflow-x: auto;
    max-height: 400px;
    overflow-y: auto;
    border-radius: 4px;
  }
  
  .entries-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.9rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  }
  
  .entries-table th {
    background-color: #1e293b; /* Dark header for dark mode */
    color: rgba(255, 255, 255, 0.9);
    padding: 0.7rem 0.5rem;
    text-align: left;
    cursor: pointer;
    position: sticky;
    top: 0;
    z-index: 10;
  }
  
  .entries-table th:hover {
    background-color: #334155; /* Slightly lighter on hover */
  }
  
  .entries-table td {
    padding: 0.7rem 0.5rem;
    border-bottom: 1px solid #334155;
    color: rgba(255, 255, 255, 0.9);
  }
  
  .entries-table tr:nth-child(even) {
    background-color: rgba(30, 41, 59, 0.5); /* Slightly darker rows */
  }
  
  .entries-table tr:hover {
    background-color: rgba(51, 65, 85, 0.8); /* Highlight on hover */
  }
  
  .delete-btn {
    background-color: #ef4444; /* Red button */
    color: white;
    border: none;
    padding: 0.4rem 0.8rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
    transition: background-color 0.2s;
  }
  
  .delete-btn:hover {
    background-color: #dc2626; /* Darker red on hover */
  }
  
  .delete-btn:disabled {
    background-color: #475569;
    cursor: not-allowed;
  }
</style> 