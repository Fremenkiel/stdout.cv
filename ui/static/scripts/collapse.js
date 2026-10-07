function toggleSchemaList(tableName) {
  document.getElementById("js-list-schema-" + tableName)
    .classList.toggle("schema__list--closed");
}
