var editor = CodeMirror.fromTextArea(document.getElementById('textarea-query'), {
  lineNumbers: true,
  mode: 'text/x-sql',
  keyMap: 'vim',
});

