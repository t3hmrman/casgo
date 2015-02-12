/* global Router */
/* global location */
/* global _ */
'use strict';

// Routes for App
RegisterRoute('ServicesRoute', {
  controllers: ['ServicesCtrl'],
  url: '/services'
});
RegisterRoute('ManageRoute', {
  controllers: [],
  url: '/manage'
});
RegisterRoute('ManageUsersRoute', {
  controllers: ['ManageUsersCtrl'],
  url: '/manage/users'
});
RegisterRoute('ManageServicesRoute', {
  controllers: ['ManageServicesCtrl'],
  url: '/manage/services'
});
RegisterRoute('StatisticsRoute', {
  controllers: [],
  url: '/statistics'
});

/**
 * Create and initialize Director's router with the internal list of functions for each route
 */
window.App.Router = Router(
  _.zipObject(_.map(window.App.Routes, 'url'),
              _.map(window.App.Routes, '_fns'))
);
window.App.Router.init();

// Go to default route on initial page load
if (location.hash === ''){ window.App.VM.gotoRoute(window.App.VM.ServicesRoute); }


/**
 * Generate pre-route utility function
 * The pre-route utility function runs before the users provided route function, and does some seutp/cleanup tasks
 *
 * @param {object} route - The route object being modified
 */
function generatePreRouteHelper(route) {
  if (!_.isObject(route)) {
    throw new Error("Only string routeUrls are allowed, invalid route object: [" + route + "]");
  }

  return function() {
    // Go to a routeUrl when it is navigated to
    window.App.VM.gotoRoute(route);
  };
}

/**
 * Utility function for registering routes
 *
 * @param {string} - Name of the route
 * @param {object} - The route object
 */
function RegisterRoute(name, route) {
  if (!_.isString(name)) { throw new Error("Expect name of registered route to be a string"); }
  if (!_.isObject(route)) { throw new Error("Expect route to be an object"); }

  route._fns = [generatePreRouteHelper(route), route.fn || null];
  window.App.Routes[route.url] = window.App.VM[name] = route;
}
