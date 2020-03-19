import { BrowserWindow } from 'electron';

export default class Loader {
    static mainWindow: Electron.BrowserWindow;
    static application: Electron.App;
    static BrowserWindow;

    private static onWindowAllClosed() {
        if (process.platform !== 'darwin') {
            Loader.application.quit();
        }
    }

    private static onClose() {
        // Dereference the window object.
        Loader.mainWindow = null;
    }

    private static onReady() {
        Loader.mainWindow = new Loader.BrowserWindow({ width: 800, height: 600 });
        Loader.mainWindow.loadFile('index.html');
        Loader.mainWindow.on('closed', Loader.onClose);
    }

    static main(app: Electron.App, browserWindow: typeof BrowserWindow) {
        // we pass the Electron.App object and the
        // Electron.BrowserWindow into this function
        // so this class has no dependencies. This
        // makes the code easier to write tests for
        Loader.BrowserWindow = browserWindow;
        Loader.application = app;
        Loader.application.on('window-all-closed', Loader.onWindowAllClosed);
        Loader.application.on('ready', Loader.onReady);
    }
}
