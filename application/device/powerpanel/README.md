To use the simulated PowerPanel devices, follow these steps:

1. Download the fogflow code repository:
        git clone https://github.com/smartfog/fogflow.git

2. Install:
        python2, pip for python2, nodejs and npm in order to run the simulated devices.

3. Start the simulated powerpanel device for "anomaly detection":
        cd  application/device/powerpanel
	      npm install 
        node powerpanel profile1.json
        node powerpanel profile2.json
        node powerpanel profile3.json

  [Note: Please change the device configuration in profile.json files.]
