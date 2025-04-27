// Create a new file to store risk management values
import { writable } from 'svelte/store';

// Default risk values
export const maxDollarRiskPerTrade = writable(500);

// Can add more shared values here as needed 