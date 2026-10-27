function matches(elt, selector) {
  // noinspection JSUnresolvedVariable
  var matchesFunction = elt.matches ||
    elt.matchesSelector || elt.msMatchesSelector || elt.mozMatchesSelector
    || elt.webkitMatchesSelector || elt.oMatchesSelector;
  return matchesFunction && matchesFunction.call(elt, selector);
}

function firstChild(elt, selector) {
  for (var i = 0; i < elt.childNodes.length; i++) {
    if (matches(elt.childNodes[i], selector)) return elt.childNodes[i];
  }

  return null
}

function overrideObject(obj, key, val) {
  obj[key] = val
  return obj;
}

htmx.defineExtension('submit-cm-selection', {
  onEvent : function(name, evt) {
    if (name === "htmx:configRequest") {
      const textArea = firstChild(evt.detail.elt, "[data-submit-cm-selection]");
      const name = textArea.getAttribute("name");

      const editor = firstChild(evt.detail.elt, ".CodeMirror");

      const selected = editor.CodeMirror.getSelectedText();

      overrideObject(evt.detail.parameters, name, selected.value);
    }
  },
});
