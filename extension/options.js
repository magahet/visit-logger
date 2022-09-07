function save_options() {
  let name = document.getElementById('name').value;
  let server = document.getElementById('server').value;
  chrome.storage.sync.set({
    name,
    server
  }, function () {
    // Update status to let user know options were saved.
    let status = document.getElementById('status');
    status.textContent = 'Options saved.';
    setTimeout(function () {
      status.textContent = '';
    }, 750);
  });
}

// Restores select box and checkbox state using the preferences
// stored in chrome.storage.
function restore_options() {
  // Use default value color = 'red' and likesColor = true.
  chrome.storage.sync.get({
    name: '',
    server: ''
  }, function (items) {
    document.getElementById('name').value = items.name
    document.getElementById('server').value = items.server
  });
}
document.addEventListener('DOMContentLoaded', restore_options);
document.getElementById('save').addEventListener('click',
  save_options);