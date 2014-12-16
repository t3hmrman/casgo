/* global Router */

// Routes for App
window.App.Routes = {
  '/services': window.App.VM.ServicesRoute,
  '/manage': window.App.VM.ManageRoute,
  '/manage/users': window.App.VM.ManageUsersRoute,
  '/manage/services': window.App.VM.ManageServicesRoute,
  '/statistics': window.App.VM.StatisticsRoute
};

// Create and initialize Director's router
window.App.Router = Router(window.App.Routes);
window.App.Router.init();
