var editor = CodeMirror.fromTextArea(document.getElementById('textarea-query'), {
  lineNumbers: true,
  mode: 'text/x-sql',
  keyMap: 'vim',
});

document.addEventListener('keydown', (e) => {
  if (e.key === "Enter" &&
    (e.composed || e.ctrlKey)) {
    document.getElementById('submit-query').click();
  }
});

