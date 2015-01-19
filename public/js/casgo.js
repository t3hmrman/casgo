/* global ko */
/* global fetch */
/* global q */
/* global _ */
'use strict';

// Create App namespace on global scope
window.App = {
  VM: null,
  Routes: {}
};

/**
 * View model for top-level casgo app
 *
 * @exports CasgoViewModel
 */
function CasgoViewModel() {
  var vm = this;

  vm.currentRouteUrl =  ko.observable(window.location.hash.slice(1));
  /**
   * Routing utility functions
   */

  /**
   * Check if the current {@link currentRouteUrl} against a given URL
   *
   * @param {string} url - The URL to check
   * @returns whether the current route url is exactly the given value
   */
  vm.currentRouteUrlIs = function(url) {return vm.currentRouteUrl() === url; };

  /**
   * Check whether the current route URL starts with a given URL
   * This usually indicates that it is a sub-route
   *
   * @param {string} prefix - The prefix to search for
   * @returns whether the current route URL contains the prefix
   */
  vm.currentRouteUrlHasPrefix = function(prefix) { return vm.currentRouteUrl().match('^' + prefix); };

  /**
   * Lookup route by it's url
   *
   * @param {string} url - The URL to use to look up the route
   * @returns the first route that matche the URL exactly
   */
  vm.getRouteByUrl = function(url) {
    return _.first(_.values(window.App.Routes), function(v) { return v.url === url; });
  };

  /**
   * Navigate to a route
   *
   * @param {object} route - The route to navigate to
   */
  vm.gotoRoute =  function(route) {
    window.location.href = '#' + route.url;
    vm.currentRouteUrl(window.location.hash.slice(1));

    // Run setup functions for all specified controllers, if any
    if (_.has(route, 'controllers') && _.isArray(route.controllers)) {

      // Run setup relevant all controllers
      _.forEach(route.controllers, function(c) {
        if (_.isString(c) &&
            _.has(vm, c) &&
            _.has(vm[c], 'setup') &&
            _.isFunction(vm[c].setup)) {
          vm[c].setup();
          // TODO: Pass route args
        }

      });
    }
  };


  /**
   * Service that maintains user sessions and information retrieved from the backend
   */
  vm.SessionService = {
    currentUser: ko.observable({}),

    /**
     * Get the current user
     *
     * @returns A Promise that will resolve to the user object
     */
    getCurrentUser: function() {
      var svc = vm.SessionService;
      return new Promise(function(resolve, reject) {
        if (_.isEmpty(svc.currentUser())) {
          svc
            .fetchCurrentUser()
            .then(resolve);
        } else {
          resolve(svc.currentUser());
        }
      });
    },

    /**
     * Fetch the current user in the session from the backend, updating the service
     *
     * @returns a promise that will resolve to the user object
     */
    fetchCurrentUser: function() {
      var svc = vm.SessionService;
      return new Promise(function(resolve, reject) {
        fetch('/api/sessions')
          .then(function(resp){
            return resp.json();
          }).then(function(json) {
            if (json.status === "success") {
              svc.currentUser(json.data);
              resolve(svc.currentUser());
            } else {
              reject(new Error("API call failed", json.message));
            }
        }).catch(function(err) {
          reject(err);
        });

      });
    }

  },

  // Services
  vm.ServicesService = {
    currentUserServices: ko.observableArray([]),
    allServices: ko.observableArray([]),

    /**
     * Get services for given user
     *
     * @param {string} userEmail - The username for which to retrieve services
     */
    getServices: function(userEmail) {
      var svc = vm.ServicesService;

      // Get the user
      var getUserFn;
      if (_.isUndefined(userEmail)) {
        getUserFn = vm.SessionService.getCurrentUser();
      } else {
        new Promise(function(resolve, reject) {resolve(userEmail);});
      }

      getUserFn
        .then(svc.fetchServicesForUser)
        .then(function(services) {
          svc.currentUserServices(services);
        });
    },

    /**
     * Get all services (admins only)
     */
    getAllServices: function() {
      var svc = vm.ServicesService;
      return new Promise(function(resolve, reject) {
        fetch('/api/services')
        .then(function(resp) { return resp.json(); })
        .then(function(json) {
          if (json.status === "success") {
            svc.allServices(json.data);
            resolve(svc.allServices());
          } else {
            reject(json.message);
          }
        }).catch(function(err) {
          reject(err);
        });
      });
    },

    /**
     * Fetch user's services and update observables
     *
     * @param {string} user - current user
     * @returns A Promise which evaluates to the users' services
     */
    fetchServicesForUser: function(user) {
      var svc = vm.ServicesService;
      // Get user's services
      return new Promise(function(resolve, reject) {
        fetch('/api/sessions/' + user.email + "/services")
          .then(function(resp) { return resp.json();})
          .then(function(json) {
          if (json.status === "success") {
            svc.currentUserServices(json.data);
            resolve(svc.currentUserServices());
          } else {
            reject(new Error("API call failed", json.message));
          }
        }).catch(function(err) {
          reject(err);
        });
      });
    }

  },

  // Controllers
  vm.ServicesCtrl = {
    pageSize: ko.observable(10),
    currentPage: ko.observable(0),
    numPages: ko.pureComputed(function() {
      var ctrl = vm.ServicesCtrl;
      return Math.Ceil(vm.ServicesService.currentUserServices() / ctrl.pageSize()) + 1;
    }),
    pagedServices: ko.pureComputed(function() {
      var ctrl = vm.ServicesCtrl;
      var startingIndex = ctrl.currentPage() * ctrl.pageSize();
      var services = vm.ServicesService.currentUserServices();

      // Return early if not enough data
      if (services.length == 0 || startingIndex > services.length) { return []; }

      return vm.ServicesService
        .currentUserServices()
        .slice(ctrl.currentPage() * ctrl.pageSize(), ctrl.pageSize());
    }),

    /**
     * Setup function for the ServicesCtrl (only run once)
     */
    setup: function() {
      vm.ServicesService.getServices();
    }
  };

  /**
   * Manage services controller
   */
  vm.ManageServicesCtrl = {
    pageSize: ko.observable(10),
    currentPage: ko.observable(0),
    numPages: ko.pureComputed(function() {
      var ctrl = vm.ServicesCtrl;
      return Math.Ceil(vm.ServicesService.allServices() / ctrl.pageSize()) + 1;
    }),
    pagedServices: ko.pureComputed(function() {
      var ctrl = vm.ServicesCtrl;
      var startingIndex = ctrl.currentPage() * ctrl.pageSize();
      var services = vm.ServicesService.allServices();

      // Return early if not enough data
      if (services.length == 0 || startingIndex > services.length) { return []; }

      return vm.ServicesService
        .allServices()
        .slice(ctrl.currentPage() * ctrl.pageSize(), ctrl.pageSize());
    }),

    /**
     * Setup function for the ManageServicesCtrl
     */
    setup: function() {
      vm.ServicesService.getAllServices();
    }
  };


  vm.ManageCtrl = {
    /**
     * The showSidebar observable controls toggling, getSidebarCSSLeft helps by returning the CSS style value to trigger animation
     */
    showSidebar: ko.observable(false),
    hideSidebar: function() { vm.ManageCtrl.showSidebar(false); },
    getSidebarCSSRight: ko.pureComputed(function() {
      var ctrl = vm.ManageCtrl;
      if (ctrl.showSidebar()) {
        return "0%";
      } else {
        return "-100%";
      }
    }),

    /**
     * Show service information for modification/saving in sidebar
     *
     * @param {object} svc - The service to show
     */
    showServiceInSidebar: function(svc) {
      var ctrl = vm.ManageCtrl;

      // Update (and set) the controller that will be attached to the template with the service it should be editing
      vm.EditServiceCtrl.currentSvc(svc);
      ctrl.sidebarController(vm.EditServiceCtrl);

      // Change the contents of the sidebar to the appropriate template and controller for services form
      ctrl.sidebarTemplateName('EditServiceFormTemplate');

      // Show sidebar
      if (!ctrl.showSidebar()) { ctrl.showSidebar(true); }
    },

    /**
     * Observables for controlling the sidebar (that could be used by any sub views in manage controller (ex. services/users)
     */
    sidebarTemplateName: ko.observable(null),
    sidebarController: ko.observable(null)
  };

  /**
   * Controller for the EditServiceFormTemplate. It is used to both create and update (edit), and contains page state.
   * Before showing this controller, calling context must initialize the controller with the service being modified (if there is one)
   */
  vm.EditServiceCtrl = {
    // Value for monitoring whether the template is create or edit mode
    create: ko.observable(false),

    // Current service being modified (empty if new)
    currentSvc: ko.observable({}),

    /**
     * Get action text for the title/other elements, depends on {@link create}'s value.
     */
    actionText: ko.pureComputed(function() {
      return vm.EditServiceCtrl.create() ? "Add service" : "Update service";
    }),

    /**
     * Remove a service, mostly a proxy call to the ServicesService, and some alerting behavior.
     */
    removeService: function() {
    },

    /**
     * Create or update a service, dispatches to {@link vm.EditServiceCtrl.createService} or {@link vm.EditServiceCtrl.updateService}
     */
    createOrUpdateService: function() {
      if (vm.EditServiceCtrl.create)
        vm.EditServiceCtrl.createService();
      else
        vm.EditServiceCtrl.updateService();
    },

    /**
     * Methods for creating and updating services, mostly proxies to ServicesService, and some alerting behavior.
     */
    createService: function() { },
    updateService: function() { }
  };

  vm.ManageUsersCtrl = {users: []};
  vm.StatisticsCtrl = {};

  /**
   * App initialization function, to be run once, when the app starts
   */
  vm.init = function() {
    vm.SessionService.fetchCurrentUser();
  };

  vm.init();
}

// Attach VM
window.App.VM = new CasgoViewModel();
ko.applyBindings(App.VM);
