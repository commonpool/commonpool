(function (window) {
  window["env"] = window["env"] || {};

  // Environment variables
  window["env"]["apiUrl"] = "${API_URL}";
  window["env"]["apiUrl"] = "${WS_URL}";
  window["env"]["debug"] = "${DEBUG}" === true;

})(this);
