// News4Coder Browser Extension - Background Service Worker

chrome.runtime.onInstalled.addListener(({ reason }) => {
  if (reason === 'install') {
    console.log('[News4Coder] Extension installed');
    // Set default API URL
    chrome.storage.local.set({
      n4cApiUrl: 'http://localhost:8081',
      n4cApiKey: ''
    });
  }
});

// Listen for keyboard shortcut
chrome.commands.onCommand.addListener((command) => {
  if (command === '_execute_action') {
    // Default action opens popup automatically
  }
});

// Optional: Context menu for quick save
chrome.runtime.onInstalled.addListener(() => {
  chrome.contextMenus.create({
    id: 'save-to-news4coder',
    title: '保存到 News4Coder',
    contexts: ['page', 'link']
  });
});

chrome.contextMenus.onClicked.addListener(async (info, tab) => {
  if (info.menuItemId === 'save-to-news4coder') {
    const url = info.linkUrl || tab.url;
    const title = tab.title || url;

    const settings = await chrome.storage.local.get(['n4cApiUrl']);
    const apiUrl = settings.n4cApiUrl || 'http://localhost:8081';

    try {
      await fetch(`${apiUrl}/api/v1/articles`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title,
          url,
          summary: '',
          tags: '',
          source: 'browser-extension',
          source_alias: 'clipper'
        })
      });
    } catch (err) {
      console.error('[News4Coder] Save failed', err);
    }
  }
});
