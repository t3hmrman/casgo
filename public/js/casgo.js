/* global ko */
/* global fetch */
/* global q */
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
	var vm = this;

	/**
	 * Routing utility functions
	 */
	vm.currentRoute =  ko.observable(window.location.hash.slice(1));
	vm.currentRouteIs = function(route) {return vm.currentRoute() === route; };
	vm.gotoRoute =  function(route) {
		window.location.href = '#' + route;
		vm.currentRoute(window.location.hash.slice(1));

		// Run setup functions for all specified controllers, if any
		if (_.has(window.App.Routes, route) &&
				_.has(window.App.Routes[route], 'controllers') &&
				_.isArray(window.App.Routes[route].controllers)) {

			_.forEach(window.App.Routes[route].controllers, function(c) {
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

		getSession: function() {
			var svc = vm.SessionService;
			fetch('/api/sessions')
				.then(function(resp){
					return resp.json();
				}).then(function(json) {
					if (json.status === "success")
						svc.currentUser(json.data);
					else
						throw new Error("API call failed", json.message);
				}).catch(function(err) {
					console.log("An error occurred retrieving user session", err);
				});
		}

	},

	// Services
	vm.ServicesService = {
		services: ko.observableArray([]),

		/**
		 * Get services for given user
     *
     * @param {string} userEmail - The username for which to retrieve services
		 */
		getServices: function(userEmail) {
			var svc = vm.ServicesService;
			userEmail = userEmail || vm.SessionService.currentUser().email;

      // Get user's services
			fetch('/api/sessions/' + userEmail + "/services")
				.then(function(resp) {
					return resp.json();
				}).then(function(json) {
					if (json.status === "success")
						svc.services(json.data);
					else
						throw new Error("API call failed", json.message);
				}).catch(function(err) {
					console.log("An Error occurred retrieving services", err);
				});
		}

	},

	// Controllers
	vm.ServicesCtrl = {
		pageSize: ko.observable(10),
		currentPage: ko.observable(0),
		numPages: ko.pureComputed(function() {
			var ctrl = vm.ServicesCtrl;
			return Math.Ceil(vm.ServicesService.services() / ctrl.pageSize()) + 1;
		}),
		pagedServices: ko.pureComputed(function() {
			var ctrl = vm.ServicesCtrl;
			var startingIndex = ctrl.currentPage() * ctrl.pageSize();
			var services = vm.ServicesService.services();

			// Return early if not enough data
			if (services.length == 0 || startingIndex > services.length) { return []; }

			return vm.ServicesService.services().slice(ctrl.currentPage() * ctrl.pageSize(), ctrl.pageSize());
		}),

		/**
		 * Setup function for the ServicesCtrl (only run once)
		 */
		setup: function() {
			vm.ServicesService.getServices();
		}
	};

	vm.ManageCtrl = {};
	vm.ManageUsersCtrl = {users: []};
	vm.ManageServicesCtrl = {services: []};
	vm.StatisticsCtrl = {};

	/**
	 * App initialization function, to be run once, when the app starts
	 */
	vm.init = function() {
		vm.SessionService.getSession();
	};

	vm.init();
}

// Attach VM
window.App.VM = new CasgoViewModel();
ko.applyBindings(App.VM);
