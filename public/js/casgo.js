/* global ko */

// Create app global
window.App = {
  VM: null,
  Routes: null
};

/**
 * View model for top-level casgo app
 *
 * @exports CasgoViewModel
 */
function CasgoViewModel() {

}
window.App.VM = new CasgoViewModel();
ko.applyBindings(App.VM);
