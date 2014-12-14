// Create App namespace on global scope
window.App = {
  VM: null,
  Routes: null
};

// Routes for App
window.App.Routes = {
  '/': function() {
    console.log("Home route");
  }
};
  
