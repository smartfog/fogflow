exports.handler = function(contextEntity, publish, query, subscribe) {
    console.log("entered into the user-defined fog function");

    var entityID = contextEntity.entityId.id;

    if (contextEntity == null) {
        return;
    }
    if (contextEntity.attributes == null) {
        return;
    }

        var FogflowIP = contextEntity.attributes.fogflowIP.value;
        FogflowIP = FogflowIP.toString();
        console.log('FogflowIP:  ',FogflowIP);

        var NGBIP = contextEntity.attributes.ngbIP.value;
        NGBIP = NGBIP.toString();
        console.log('NGBIP:  ',NGBIP);

        var command = "../transformer-config.sh " + FogflowIP + " " + NGBIP;
        console.log(command);
        const shell = require('shelljs');
        shell.exec(command);
//        command = "../transformer-config.sh " + FogflowIP + " " + NGBIP;
//        console.log(command);        
//        const shell = require('shelljs');
//        shell.exec(command);


};
