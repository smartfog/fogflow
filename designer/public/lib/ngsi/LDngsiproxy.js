(function() {
    
var LDNGSIProxy = (function() {
    var LDNGSIProxy = function() {
        this.subscriptions = [];
        this.notifyHandler = null;
        this.socketURL = window.location.protocol + "//" + window.location.hostname + ":" + window.location.port;
        this.socket = io.connect(this.socketURL);
        
        var self = this;
        this.socket.on('connect', function() {
            console.log('connected to the remote NGSI agent');
            self.socket.emit('subscriptions', self.subscriptions);
        });
        this.socket.on('notify', function(data) {
            entities = data.entities;
            if(self.notifyHandler) {
                self.notifyHandler(entities);
            }
        });
    };      
    
    LDNGSIProxy.prototype.reportSubID = function reportSubID(sid) {
        if(this.socket.connected) {
            console.log('connected');
            var newSubscriptions = [];
            newSubscriptions.push(sid);
            this.socket.emit('subscriptions', newSubscriptions)
        } else {
            console.log('not connected');
        	this.subscriptions.push(sid);			
        }
    };    
    
    LDNGSIProxy.prototype.setNotifyHandler = function setNotifyHandler(cb) {
        this.notifyHandler = cb;
    };

    return LDNGSIProxy;
})();

    window.LDNGSIProxy = LDNGSIProxy;
})();
