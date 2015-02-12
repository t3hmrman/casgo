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

  //////////////
  // Services //
  //////////////

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
      var self = vm.ServicesService;
      // Get user's services
      return new Promise(function(resolve, reject) {
        fetch('/api/sessions/' + user.email + "/services")
          .then(function(resp) { return resp.json();})
          .then(function(json) {
            if (json.status === "success") {
              self.currentUserServices(json.data);
              resolve(self.currentUserServices());
            } else {
              reject(new Error("API call failed", json.message));
            }
          }).catch(function(err) {
            reject(err);
          });
      });
    },

    /**
     * Create a service
     *
     * @param {object} svc - Service to be created
     * @returns A Promise for the ajax request
     */
    createService: function(svc) {
      var self = vm.ServicesService;
      if (_.isUndefined(svc) || !self.isValidService(svc)) throw new Error("Invalid service:", svc);

      return fetch('/api/services', {
        method: 'post',
        headers: { 'Accept': 'application/json', 'Content-Type': 'application/json'},
        body: JSON.stringify(svc)
      });
    },

    /**
     * Update a service
     *
     * @param {object} svc - Service to be updated
     * @returns A Promise for the ajax request
     */
    updateService: function(svc) {
      var self = vm.ServicesService;
      if (_.isUndefined(svc) || !self.isValidService(svc)) throw new Error("Invalid service:", svc);

      return fetch('/api/services/' + svc.name, {
        method: 'put',
        headers: { 'Accept': 'application/json', 'Content-Type': 'application/json'},
        body: JSON.stringify(svc)
      });
    },

    /**
     * Check if a given service is valid
     */
    isValidService: function(svc) {
      return _(['url', 'adminEmail', 'name'])
        .map(function(k) { return k in svc; })
        .every();
    },

    /**
     * Delete a service (by name)
     *
     * @param {object} svc - Service to be deleted
     * @returns A Promise for the ajax request
     */
    deleteService: function(svc) {
      var self = vm.ServicesService;
      if (_.isUndefined(svc) || !self.isValidService(svc)) throw new Error("Invalid service:", svc);
      return fetch('/api/services/' + svc.name, {
        method: 'delete',
        headers: { 'Accept': 'application/json', 'Content-Type': 'application/json'}
      });
    }

  },

  vm.UsersService = {
    currentUser: ko.observable(null),
    allUsers: ko.observableArray([]),

    /**
     * Get all users (admins only)
     */
    getAllUsers: function() {
      var svc = vm.UsersService;
      return new Promise(function(resolve, reject) {
        fetch('/api/users')
          .then(function(resp) { return resp.json(); })
          .then(function(json) {
            if (json.status === "success") {
              svc.allUsers(json.data);
              resolve(svc.allUsers());
            } else {
              reject(json.message);
            }
          }).catch(function(err) {
            reject(err);
          });
      });
    }

  },


  /**
   * Generate a paginated list controller
   *
   * @param {function} totalFn - The function to be called to retrieve all the objects (ex. from a service)
   * @param {function} setupFn - The function to be called once to setup the controller
   */
  vm.generatePaginatedListController = function(totalFn, setupFn) {
    var ctrl = {};
    ctrl.pageSize = ko.observable(10);
    ctrl.currentPage = ko.observable(0);

    ctrl.numPages =  ko.pureComputed(function() {
      return Math.Ceil(totalFn() / ctrl.pageSize()) + 1;
    });

    ctrl.pagedObjects = ko.pureComputed(function() {
      var startingIndex = ctrl.currentPage() * ctrl.pageSize();
      var objects = totalFn();

      // Return early if not enough data
      if (objects.length == 0 || startingIndex > objects.length) { return []; }

      return totalFn().slice(ctrl.currentPage() * ctrl.pageSize(), ctrl.pageSize());
    });

    ctrl.setup = function() {
      if (!_.isUndefined(setupFn)) setupFn();
    };
    return ctrl;
  }


  /////////////////
  // Controllers //
  /////////////////

  /**
   * Basic paginated list controllers for manage services/users tabs, which contain not much more than a basic list
   */
  vm.ServicesCtrl = vm.generatePaginatedListController(vm.ServicesService.currentUserServices, vm.ServicesService.getAllServices);
  vm.ManageServicesCtrl = vm.generatePaginatedListController(vm.ServicesService.allServices, vm.ServicesService.getAllServices);
  vm.ManageUsersCtrl = vm.generatePaginatedListController(vm.UsersService.allUsers, function() {
    vm.UsersService.getAllUsers();
    // Set focus to username filter when controller is shown
    document.getElementById('usernameFilter').focus();
  });

  vm.ManageUsersCtrl.usernameFilter = ko.observable('');
  vm.ManageUsersCtrl.filteredUsers = ko.computed(function() {
    var ctrl = vm.ManageUsersCtrl;
    var filterTerm = ctrl.usernameFilter();
    var pagedObjects = vm.ManageUsersCtrl.pagedObjects();
    if (_.isEmpty(filterTerm)) { return pagedObjects; };

    // Check the user's email to see if it contains the filter term
    return _.filter(pagedObjects, function(o) { return o.email.contains(filterTerm); });
  });

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
     * Show new service form in sidebar by clearing the {@link vm.EditServiceCtrl}, then
     * showing the sidebar using {@link vm.ManageCtrl.showSidebarEditServiceForm}
     */
    showAddServiceInSidebar: function() {
      var self = vm.ManageCtrl;
      var ctrl = vm.EditServiceCtrl;

      // Set up the edit service controller (also used for creating services)
      vm.EditServiceCtrl.create(true);
      vm.EditServiceCtrl.loadSvc({});

      // Set the service controller as the controller in the sidebar, show the sidebar
      self.sidebarController(ctrl);
      vm.ManageCtrl.showSidebarEditServiceForm();
    },

    /**
     * Show edit service form in sidebar by loading the appropriate service into
     * {@link vm.EditServiceCtrl} and showing the sidebar using {@link vm.ManageCtrl.showSidebarEditServiceForm}
     *
     * @param {object} svc - The service to edit
     */
    showEditServiceInSidebar: function(svc) {
      var self = vm.ManageCtrl;
      var ctrl = vm.EditServiceCtrl;
      ctrl.create(false);
      ctrl.loadSvc(svc);
      self.sidebarController(ctrl);
      vm.ManageCtrl.showSidebarEditServiceForm();
    },

    /**
     * Show sidebar edit form
     */
    showSidebarEditServiceForm: function() {
      var self = vm.ManageCtrl;
      // Change to the appropriate template and controller for edit service form
      self.sidebarTemplateName('EditServiceFormTemplate');

      // Show sidebar (if not already visible)
      if (!self.showSidebar()) { self.showSidebar(true); }
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
    // Alerts
    alerts: ko.observableArray([]),

    // Utility function for deleting alerts
    dismissAlert: function(alert) {
      var ctrl = vm.EditServiceCtrl;
      ctrl.alerts(_.reject(ctrl.alerts(), function(a) { return _.isEqual(alert, a); }));
    },

    // Value for monitoring whether the template is create or edit mode
    create: ko.observable(true),

    // Current service being modified (empty if new)
    svcId: ko.observable(undefined),
    svcName: ko.observable(""),
    svcUrl: ko.observable(""),
    svcAdminEmail: ko.observable(""),

    /**
     * Load the controller with an existing service's data
     */
    loadSvc: function(svc) {
      var ctrl = vm.EditServiceCtrl;
      ctrl.svcName(svc.name || "");
      ctrl.svcUrl(svc.url || "");
      ctrl.svcAdminEmail(svc.adminEmail || "");
    },

    /**
     * Create a service out of the information currently being held by the controller
     */
    makeService: ko.pureComputed(function() {
      var ctrl = vm.EditServiceCtrl;

      // Create service object from current state
      var svc =  {
        name: ctrl.svcName(),
        url: ctrl.svcUrl(),
        adminEmail: ctrl.svcAdminEmail()
      };

      // Add ID if present
      var svcId = ctrl.svcId();
      if (!_.isUndefined(svcId)) { svc.id = svcId; }

      return svc;
    }),

    /**
     * Get action text for the title/other elements, depends on {@link create}'s value.
     */
    actionText: ko.pureComputed(function() {
      return vm.EditServiceCtrl.create() ? "Add service" : "Update service";
    }),

    /**
     * Remove a service, mostly a proxy call to the {@link vm.ServicesService}, and some alerting behavior.
     */
    removeService: function() {
      var servicesService = vm.ServicesService;
      var ctrl = vm.EditServiceCtrl;
      servicesService.deleteService(ctrl.makeService())
        .then(function(resp) {
          return resp.json();
        }).then(function(json) {
          var msg = json.status === "success" ? "Successfully removed service" : "Failed to remove service";
          ctrl.alerts.push({type: json.status, msg: msg});
        });
    },

    /**
     * Create or update a service, mostly a proxy call to the {@link vm.ServicesService}, and some alerting behavior"
     */
    createOrUpdateService: function() {
      var serviceService = vm.ServicesService;
      var ctrl = vm.EditServiceCtrl;
      var service = ctrl.makeService();
      var createdService = ctrl.create();
      var createOrUpdateFn = createdService ? serviceService.createService : serviceService.updateService;
      createOrUpdateFn(service)
        .then(function(resp) {
          return resp.json();
        }).then(function(json) {
          var verb = createdService ? "created" : "updated";
          var msg = json.status === "success" ? "Successfully " + verb + " service" : "Service not successfully " + verb + ".";

          // Append error message from server if provided
          if (_.has(json, 'message')) { msg += " " + json.message;}

          ctrl.alerts.push({type: json.status, msg: msg});
        });
    }
  };

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
