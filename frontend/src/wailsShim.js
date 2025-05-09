// Wails shim for development mode
// This provides mock implementations of window.go and window.runtime 
// to avoid errors when running the frontend separately with Vite

// Detect if we're running in development mode
const isDev = import.meta.env.MODE === 'development' || !window.go || !window.runtime;

// Mock Go bindings
if (isDev && !window.go) {
  console.log("Initializing mock Go runtime for development");
  
  window.go = {
    main: {
      App: {
        // Mock methods with dummy responses
        GetStockRatings: async () => {
          console.log("Mock: GetStockRatings called");
          return { 
            overall: 0, 
            sectorRatings: {} 
          };
        },
        SaveStockRating: async (rating) => {
          console.log("Mock: SaveStockRating called with", rating);
          return true;
        },
        DeleteStockRating: async (id) => {
          console.log("Mock: DeleteStockRating called with", id);
          return true;
        },
        GetRiskAssessments: async () => {
          console.log("Mock: GetRiskAssessments called");
          return [];
        },
        SaveRiskAssessment: async (assessment) => {
          console.log("Mock: SaveRiskAssessment called with", assessment);
          return true;
        },
        GetPositionSettings: async () => {
          console.log("Mock: GetPositionSettings called");
          return {
            maxRiskPerTrade: 2,
            maxPortfolioRisk: 10,
            accountSize: 10000
          };
        },
        SavePositionSettings: async (settings) => {
          console.log("Mock: SavePositionSettings called with", settings);
          return true;
        }
      }
    }
  };
}

// Mock runtime bindings
if (isDev && !window.runtime) {
  console.log("Initializing mock Wails runtime for development");
  
  window.runtime = {
    WindowSetDarkTheme: () => console.log("Mock: WindowSetDarkTheme called"),
    WindowSetLightTheme: () => console.log("Mock: WindowSetLightTheme called"),
    LogPrint: (msg) => console.log("Wails Log:", msg),
    LogTrace: (msg) => console.log("Wails Trace:", msg),
    LogDebug: (msg) => console.log("Wails Debug:", msg),
    LogInfo: (msg) => console.log("Wails Info:", msg),
    LogWarning: (msg) => console.warn("Wails Warning:", msg),
    LogError: (msg) => console.error("Wails Error:", msg),
    LogFatal: (msg) => console.error("Wails Fatal:", msg),
    EventsOnMultiple: (evt, callback, max) => {
      console.log(`Mock: EventsOnMultiple registered for ${evt}`);
      return () => console.log(`Mock: EventsOnMultiple unregistered for ${evt}`);
    },
    EventsOff: (evt) => console.log(`Mock: EventsOff called for ${evt}`),
    EventsEmit: (evt, ...args) => console.log(`Mock: EventsEmit called for ${evt} with`, args),
    WindowReload: () => window.location.reload(),
    BrowserOpenURL: (url) => window.open(url, '_blank')
  };
}

// Export a function to check if we're running in development mode with shims
export const isUsingShims = () => isDev; 