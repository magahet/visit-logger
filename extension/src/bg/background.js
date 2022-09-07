// if you checked "fancy-settings" in extensionizr.com, uncomment this lines

// var settings = new Store("settings", {
//     "sample_setting": "This is how you use Store.js to remember values"
// });

// Where we will expose all the data we retrieve from storage.sync.
const storageCache = {};
// Asynchronously retrieve data from storage.sync, then cache it.
const initStorageCache = getAllStorageSyncData().then(items => {
  // Copy the data retrieved from storage into storageCache.
  Object.assign(storageCache, items);
});

function getAllStorageSyncData() {
  // Immediately return a promise and start asynchronous work
  return new Promise((resolve, reject) => {
    // Asynchronously fetch all data from storage.sync.
    chrome.storage.sync.get(null, (items) => {
      // Pass any observed errors down the promise chain.
      if (chrome.runtime.lastError) {
        return reject(chrome.runtime.lastError);
      }
      // Pass the data retrieved from storage down the promise chain.
      resolve(items);
    });
  });
}


chrome.tabs.onUpdated.addListener(async (tabId, changeInfo, tab) => {
  try {
    await initStorageCache
  } catch (e) {
    console.error(e)
    return
  }

  const name = storageCache?.name
  const server = storageCache?.server

  console.log(name, server)

  if (!name || !server) {
    return
  }

  const entity = {
    'title': tab.title,
    'url': tab.url,
    'name': name
  }
  console.log(entity)

  //https://us-west1-gars-cloud.cloudfunctions.net/visit-logger

  fetch(server, {
    method: 'POST',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(entity)
  })
    .then(resp => {
      if (!resp.ok) {
        return resp.text().then(text => { throw new Error(text) })
      }
    })
    .catch(error => console.log('Error:', error));
});