/* global ko */

// Create App namespace on global scope
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
  var self = this;

  // Routes
  self.ServicesRoute = function() { console.log("Services route!"); };
  self.ManageRoute = function() { console.log("Manage route!"); };
  self.ManageUsersRoute = function() { console.log("Manage users route!"); };
  self.ManageServicesRoute = function() { console.log("Manage services route!"); };
  self.StatisticsRoute = function() { console.log("Statistics route!"); };
}

// Attach VM
window.App.VM = new CasgoViewModel();
ko.applyBindings(App.VM);
