/* global ko */
'use strict';

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
  self.ManageRoute = function() { console.log("Manage route!"); };
  self.ManageUsersRoute = function() { console.log("Manage users route!"); };
  self.ManageServicesRoute = function() { console.log("Manage services route!"); };
  self.StatisticsRoute = function() { console.log("Statistics route!"); };

  // Routing utility functions
  self.currentRoute =  ko.observable(window.location.hash.slice(1));
  self.currentRouteIs = function(route) {return self.currentRoute() === route; };
  self.gotoRoute =  function(route) {
    console.log("Going to route:", route);
    window.location.href = '#' + route;
    self.currentRoute(window.location.hash.slice(1));
  };


  // Controllers
  self.ServicesCtrl = {
    setup: function() {
      console.log("Would have loaded services!");
    },
    test: function() { console.log("TEST!"); },
    services: ko.observableArray([])
  };

  self.ManageCtrl = {};
  self.ManageUsersCtrl = {users: []};
  self.ManageServicesCtrl = {services: []};
  self.StatisticsCtrl = {};



}

// Attach VM
window.App.VM = new CasgoViewModel();
ko.applyBindings(App.VM);
