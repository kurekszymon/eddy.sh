package javascript

// Manual Installation
// If you need to install manually, there is a separate download called nvm-noinstall.zip. You should also uninstall any existing versions of node.js.

// Download nvm-noinstall.zip. Extract this to the directory where NVM should be "installed". The default directory used by the installer is C:\Users\<username>\AppData\Roaming\nvm, but you can use whatever you like.
// The zip archive contains three files, including nvm.exe, elevate.vbs, and elevate.cmd. All three of these are required for NVM to function properly. The "elevate" scripts help elevate administrative permissions for actions that require it. This is a critical component for switching between versions of node.js.

// NVM for Windows "switches" versions of node.js by updating a symlink, using the mklink command. The symlink is recreated to point to whichever version of node.js should run. This process requires elevated administrative permissions.

// Update the system environment variables.
// There are two system environment variables that need to be created, and one that needs to be modified.

// First, add a new environment variable called NVM_HOME. This should be set to the directory from step 1. If you used the default, this would be C:\Users\<username>\AppData\Roaming\nvm.

// Second, add a new environment variable called NVM_SYMLINK. This should be set to the path that will be used to identify which version of node.js is running. THIS DIRECTORY SHOULD NOT EXIST. It will automatically be created and maintained by NVM.

// Finally, update the system path by appending %NVM_HOME%;%NVM_SYMLINK% to the end. The result should look something like:
