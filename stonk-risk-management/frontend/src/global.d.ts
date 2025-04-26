interface Window {
  go: {
    main: {
      App: {
        GetRiskAssessments(): Promise<any[]>;
        SaveRiskAssessment(assessment: any): Promise<void>;
        GetStockRatings(): Promise<any[]>;
        GetStockRatingsByDate(date: Date): Promise<any[]>;
        SaveStockRating(rating: any): Promise<void>;
        DeleteStockRating(id: string): Promise<void>;
        DeleteRiskAssessment(id: string): Promise<void>;
      }
    }
  },
  Chart: any;
} 