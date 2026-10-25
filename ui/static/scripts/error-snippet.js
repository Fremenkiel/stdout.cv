function RenderErrorSnippet(editor, parentId, lineNumber, startChar, endChar) {
  const firstLineNumber = lineNumber == 1 ? 1 : lineNumber - 1;

  // line index start at 0
  const value = editor.CodeMirror.doc.getRange({line: firstLineNumber - 1, ch: 0}, {line: lineNumber});

  var parent = document.getElementById(parentId);

  var node = document.createElement("div");
  parent.appendChild(node);

  var snippet = CodeMirror(node, {
    value: value,
    firstLineNumber: firstLineNumber,
    lineNumbers: true,
    mode: 'text/x-sql',
    readOnly: true,
    showCursorWhenSelecting: false,
  });

  snippet.doc.markText(
    { line: 1, ch: startChar},
    { line: 1, ch: endChar},
    { css: "background-color: var(--color-background-error); color: var(--color-text-error); text-decoration: underline wavy var(--color-text-error-decoration)" }
  );

  snippet.doc.addLineClass(1, "wrap", "error-line")
}

function ResetErrorSnippet(parentId) {
  var parent = document.getElementById(parentId);

  for(var node in parent.childNodes) {
    node.remove();
  }
}

