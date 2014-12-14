/* global ko */

/**
 * View model for top-level casgo app
 *
 * @exports CasgoViewModel
 */
function CasgoViewModel() {

}

// Attach VM
window.App.VM = new CasgoViewModel();
ko.applyBindings(App.VM);
