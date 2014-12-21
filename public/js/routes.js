/* global Router */
/* global location */

'use strict';

// Routes for App
window.App.VM.ServicesRoute = {
  controllers: ['ServicesController'],
  fn: function() {
    window.App.VM.ServicesService.loadServices();
  }
};

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
    window.App.VM.gotoRoute(route);
  };
}


// Route registration
window.App.Routes = {
  default: '/services',
  '/services': window.App.VM.ServicesRoute.fn,
  '/manage': window.App.VM.ManageRoute,
  '/manage/users': window.App.VM.ManageUsersRoute,
  '/manage/services': window.App.VM.ManageServicesRoute,
  '/statistics': window.App.VM.StatisticsRoute
};

// Add PreRouteHelper to all routes
_.forEach(window.App.Routes, function(v, k) {
  if (k !== "default") window.App.Routes[k] = [generatePreRouteHelper(k), window.App.Routes[k]];
});


// Create and initialize Director's router
window.App.Router = Router(window.App.Routes);
window.App.Router.init();

// Go to default route
if (location.hash === '' && _.has(window.App.Routes, 'default')) {
  window.App.VM.gotoRoute(window.App.Routes.default);
}
