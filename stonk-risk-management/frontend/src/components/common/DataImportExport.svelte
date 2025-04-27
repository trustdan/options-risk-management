<script>
  import { onMount } from 'svelte';
  import { ExportData, ImportData, ShowSaveDialog, ShowOpenDialog, GetCurrentDirectory, GetUserDirectory } from '../../../wailsjs/go/main/App';
  
  // Props
  export let position = 'header'; // 'header' or 'footer'
  
  let isExporting = false;
  let isImporting = false;
  let exportError = '';
  let importError = '';
  let exportSuccess = '';
  let importSuccess = '';
  
  // Clear status messages after some time
  function clearStatus() {
    setTimeout(() => {
      exportError = '';
      importError = '';
      exportSuccess = '';
      importSuccess = '';
    }, 5000);
  }

  // Import Wails runtime for file dialogs
  let runtime;
  onMount(async () => {
    // @ts-ignore
    runtime = window.runtime;
    console.log("Available runtime methods:", runtime ? Object.keys(runtime) : "Runtime not available");
    console.log("Runtime object:", runtime);
  });
  
  // Handle export button click
  async function handleExport() {
    isExporting = true;
    exportError = '';
    exportSuccess = '';
    
    try {
      // Get today's date for the filename
      const today = new Date();
      const dateStr = today.toISOString().split('T')[0]; // YYYY-MM-DD format
      
      const defaultPath = `options-trading-backup-${dateStr}.json`;
      
      // Use our Go backend method to show the save dialog
      let filePath = '';
      try {
        // @ts-ignore - Ignore all type checking for this block
        const result = await ShowSaveDialog({
          title: "Save data backup file",
          defaultFilename: defaultPath,
          filters: [{ name: "JSON Files", extensions: ["json"] }]
        });
        
        // @ts-ignore - Ignore type checking for this conditional block
        if (typeof result === 'string') {
          filePath = result;
        } 
        // @ts-ignore - Ignore type checking for this conditional block
        else if (Array.isArray(result) && result.length > 0) {
          // @ts-ignore
          filePath = result[0];
        }
      } catch (dialogError) {
        console.error("Save dialog error:", dialogError);
        
        // Get the current directory to show in the prompt
        let currentDir;
        try {
          currentDir = await GetCurrentDirectory();
        } catch (dirError) {
          console.error("Failed to get current directory:", dirError);
          currentDir = "current directory";
        }
        
        // Fallback for users to input the path manually
        const promptResult = prompt(
          `Please enter the full path to save your backup file. If you enter just a filename, it will be saved in: ${currentDir}`,
          defaultPath
        );
        
        if (promptResult) {
          filePath = promptResult;
        }
      }
      
      if (!filePath) {
        isExporting = false;
        return; // User cancelled
      }
      
      // Call backend export function
      await ExportData(filePath);
      
      exportSuccess = `Data successfully exported to ${filePath}`;
      console.log('Export successful:', filePath);
    } catch (error) {
      exportError = `Export failed: ${error.message || error}`;
      console.error('Export error:', error);
    } finally {
      isExporting = false;
      clearStatus();
    }
  }
  
  // Handle import button click
  async function handleImport() {
    isImporting = true;
    importError = '';
    importSuccess = '';
    
    try {
      // Use our Go backend method to show the open dialog
      let filePath = '';
      try {
        // @ts-ignore - Ignore all type checking for this block
        const result = await ShowOpenDialog({
          title: "Select backup file to import",
          filters: [{ name: "JSON Files", extensions: ["json"] }]
        });
        
        // @ts-ignore - Ignore type checking for this conditional block
        if (typeof result === 'string') {
          filePath = result;
        } 
        // @ts-ignore - Ignore type checking for this conditional block
        else if (Array.isArray(result) && result.length > 0) {
          // @ts-ignore
          filePath = result[0];
        }
      } catch (dialogError) {
        console.error("Open dialog error:", dialogError);
        
        // Get the current directory to show in the prompt
        let currentDir;
        try {
          currentDir = await GetCurrentDirectory();
        } catch (dirError) {
          console.error("Failed to get current directory:", dirError);
          currentDir = "current directory";
        }
        
        // Fallback for users to input the path manually
        const promptResult = prompt(
          `Please enter the full path to your backup file. If you enter just a filename, it will be loaded from: ${currentDir}`
        );
        
        if (promptResult) {
          filePath = promptResult;
        }
      }
      
      if (!filePath) {
        isImporting = false;
        return; // User cancelled
      }
      
      // Call backend import function
      await ImportData(filePath);
      
      importSuccess = 'Data successfully imported';
      console.log('Import successful');
      
      // Refresh the application to show imported data
      setTimeout(() => {
        window.location.reload();
      }, 1500);
    } catch (error) {
      importError = `Import failed: ${error.message || error}`;
      console.error('Import error:', error);
    } finally {
      isImporting = false;
      clearStatus();
    }
  }
</script>

<div class="data-import-export {position}">
  <button 
    class="export-btn" 
    on:click={handleExport} 
    disabled={isExporting || isImporting}
    title="Export all data to a backup file"
  >
    {#if isExporting}
      ⏳
    {:else}
      📤 Export Data
    {/if}
  </button>
  
  <button 
    class="import-btn"
    on:click={handleImport} 
    disabled={isExporting || isImporting}
    title="Import data from a backup file"
  >
    {#if isImporting}
      ⏳
    {:else}
      📥 Import Data
    {/if}
  </button>
  
  {#if exportError}
    <span class="error">{exportError}</span>
  {/if}
  
  {#if importError}
    <span class="error">{importError}</span>
  {/if}
  
  {#if exportSuccess}
    <span class="success">{exportSuccess}</span>
  {/if}
  
  {#if importSuccess}
    <span class="success">{importSuccess}</span>
  {/if}
</div>

<style>
  .data-import-export {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }
  
  .data-import-export.header {
    margin-right: 1rem;
  }
  
  .data-import-export.footer {
    justify-content: center;
    margin: 0.5rem 0;
  }
  
  button {
    background-color: var(--secondary-button);
    color: inherit;
    border: none;
    padding: 0.5rem 0.75rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
  }
  
  button:hover {
    background-color: var(--secondary-button-hover);
  }
  
  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  .export-btn {
    background-color: var(--secondary-button);
  }
  
  .import-btn {
    background-color: var(--secondary-button);
  }
  
  .error, .success {
    font-size: 0.8rem;
    padding: 0.2rem 0.5rem;
    border-radius: 4px;
    animation: fadeIn 0.3s ease;
  }
  
  .error {
    background-color: var(--warning-bg);
    color: var(--danger-color);
  }
  
  .success {
    background-color: var(--success-bg);
    color: var(--success-color);
  }
  
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
</style> 