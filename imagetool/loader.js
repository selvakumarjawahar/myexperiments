"use strict";
exports.__esModule = true;
var Loader = /** @class */ (function () {
    function Loader() {
    }
    Loader.onWindowAllClosed = function () {
        if (process.platform !== 'darwin') {
            Loader.application.quit();
        }
    };
    Loader.onClose = function () {
        // Dereference the window object.
        Loader.mainWindow = null;
    };
    Loader.onReady = function () {
        Loader.mainWindow = new Loader.BrowserWindow({ width: 800, height: 600 });
        Loader.mainWindow.loadFile('index.html');
        Loader.mainWindow.on('closed', Loader.onClose);
    };
    Loader.main = function (app, browserWindow) {
        // we pass the Electron.App object and the
        // Electron.BrowserWindow into this function
        // so this class has no dependencies. This
        // makes the code easier to write tests for
        Loader.BrowserWindow = browserWindow;
        Loader.application = app;
        Loader.application.on('window-all-closed', Loader.onWindowAllClosed);
        Loader.application.on('ready', Loader.onReady);
    };
    return Loader;
}());
exports["default"] = Loader;
