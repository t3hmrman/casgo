/* global Router */
/* global location */

'use strict';

var AppVM = window.App.VM;

// Routes for App
AppVM.ServicesRoute = { controllers: ['ServicesCtrl'] };
AppVM.ManageRoute = { controllers: [] };
AppVM.ManageUsersRoute = { controllers: [] };
AppVM.ManageServicesRoute = { controllers: [] };
AppVM.StatisticsRoute = { controllers: [] };

/**
 * Generate pre-route utility function
 * The pre-route utility function runs before the users provided route function, and does some seutp/cleanup tasks
 *
 * @param {string} route - The path of the route
 */
function generatePreRouteHelper(route) {
  if (!_.isString(route)) { throw new Error("Only string routes are allowed, invalid route: [" + route + "]"); }
  return function() {
    // Go to a route when it is navigated to
    AppVM.gotoRoute(route);
  };
}


// Route registration
window.App.Routes = {
  default: '/services',
  '/services': AppVM.ServicesRoute,
  '/manage': AppVM.ManageRoute,
  '/manage/users': AppVM.ManageUsersRoute,
  '/manage/services': AppVM.ManageServicesRoute,
  '/statistics': AppVM.StatisticsRoute
};

// Add PreRouteHelper to all routes
_.forEach(window.App.Routes, function(v, k) {
  if (k !== "default") {
    // Create an internal list of functions to call when the route happens
    window.App.Routes[k]._fns = [generatePreRouteHelper(k), window.App.Routes[k].fn || null];
  }
});

/**
 * Create and initialize Director's router with the internal list of functions for each route
 */
window.App.Router = Router(
  _.zipObject(_.keys(window.App.Routes),
              _.map(window.App.Routes, '_fns'))
);
window.App.Router.init();

// Go to default route
if (location.hash === '' && _.has(window.App.Routes, 'default')) {
  AppVM.gotoRoute(window.App.Routes.default);
}
