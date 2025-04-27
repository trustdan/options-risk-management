# Stonk Risk Management

A comprehensive desktop application for options traders to manage risk, track stock sentiments, and visualize upcoming trades.

## Features

### Risk Management Dashboard
- Daily emotional and psychological state assessment
- FOMO (Fear of Missing Out) level tracking
- Market bias monitoring
- Physical condition assessment
- P&L impact analysis
- Position sizing recommendations based on psychological state

### Trading Journal
- Record and reflect on completed trades
- Track emotional state during trades
- Document lessons learned and improvement plans
- Color-coded strategy visualization
- Win/loss tracking with P&L recording
- Consistent strategy categorization with Options Trading Calendar

### Stock Rating Dashboard
- Market and sector sentiment tracking
- Individual stock ratings
- Chart pattern identification
- Sector-specific analysis
- Historical ratings view

### Options Trading Calendar
- Visual calendar of options trades by expiration week
- Color-coded strategy visualization
- Multi-sector trade organization
- Trade entry and management
- Strategy guide with visual indicators
- Complete options strategy selection including:
  - Basic strategies (Long/Short Calls and Puts)
  - Income strategies (Covered Calls and Cash-Secured Puts)
  - Advanced strategies (Spreads, Iron Condors, Butterflies)

## Installation

### Windows
1. Download the latest installer (`stonk-risk-management-setup-x.x.x.exe`) from the releases page
2. Run the installer and follow the on-screen instructions
3. Launch the application from the Start menu or desktop shortcut

### Building from Source
To build the application from source, you'll need:
- Go 1.18 or higher
- Node.js 16 or higher
- Wails (https://wails.io)

Clone the repository and run:
```
wails build
```

## Usage

### Getting Started
1. Start with the Risk Management Dashboard to assess your daily trading psychology
2. Use the Stock Rating Dashboard to identify trading opportunities
3. Plan and visualize your trades with the Options Trading Calendar

### Data Storage
All your data is stored locally in `~/.stonk-risk-management` directory.

### Supported Options Strategies
The application supports a comprehensive range of options trading strategies, including:

#### Basic Spreads
- Long Call/Put - Directional bets on market movement
- Covered Call - Owning stock and selling calls against it for premium income

#### Vertical Spreads
- Bull Call Spread - Bullish strategy with limited risk/reward
- Bear Call Spread - Bearish strategy with limited risk/reward
- Bull Put Spread - Bullish strategy with limited risk/reward
- Bear Put Spread - Bearish strategy with limited risk/reward

#### Advanced Strategies
- Calendar/Horizontal Spreads - Time-based spreads with different expiration dates
- Diagonal Spreads - Combining different strikes and expiration dates
- Butterfly Spreads - Neutral strategies for limited price movement
- Iron Condors/Butterflies - Range-bound strategies for selling premium
- Ratio Spreads (Backspreads) - Unbalanced spreads with more long than short options

#### Danger - naked options ahead
- Short Call - Bearish directional bet with unlimited risk potential
- Short Put - Bullish directional bet with significant downside risk
- Cash-Secured Put - Selling puts with cash to secure the potential stock purchase (risk of capital misallocation)
- Call Ratio Spread - Neutral to moderately bearish with naked call exposure (more short than long options)
- Put Ratio Spread - Neutral to moderately bullish with naked put exposure (more short than long options)

## Development

### Prerequisites
- Go 1.18+
- Node.js 16+
- Wails CLI
- Svelte knowledge

### Dev Environment
Run the development server with:
```
wails dev
```

### Common Issues & Solutions

#### Date Handling in JavaScript/Go Applications
When working with dates between JavaScript frontend and Go backend:

1. **ISO String Format**: Always convert JavaScript Date objects to ISO strings before sending to Go:
   ```javascript
   // Correct way to send dates to Go backend
   const dateStr = new Date(date).toISOString().split('T')[0];
   const formattedDate = `${dateStr}T00:00:00Z`; // Use complete ISO format
   ```

2. **Date Comparison**: When comparing dates, convert both to the same format first:
   ```javascript
   // Convert both dates to YYYY-MM-DD format before comparing
   const date1 = new Date(firstDate).toISOString().split('T')[0];
   const date2 = new Date(secondDate).toISOString().split('T')[0];
   const datesMatch = date1 === date2;
   ```

3. **Using Wails Models**: When using generated model classes, create objects with `createFrom()`:
   ```javascript
   import { models } from '../wailsjs/go/models';
   
   // Create data object
   const ratingData = {
     id: existingRating ? existingRating.id : '',
     date: formattedDate, // ISO string, not Date object
     // other properties...
   };
   
   // Create proper model instance
   const ratingToSave = models.StockRating.createFrom(ratingData);
   
   // Send to backend
   await SaveStockRating(ratingToSave);
   ```

These solutions address common issues with time.Time unmarshalling errors in Go and date comparison mismatches in JavaScript.

## License
MIT

## Credits
Created by [Your Name]
